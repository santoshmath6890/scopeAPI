package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	Path    string
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// RunMigrations executes all pending migrations
func (mr *MigrationRunner) RunMigrations(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := mr.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := mr.getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := mr.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if mr.isMigrationApplied(migration.Version, appliedMigrations) {
			log.Printf("Migration %s already applied, skipping", migration.Version)
			continue
		}

		log.Printf("Applying migration %s: %s", migration.Version, migration.Name)
		if err := mr.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		log.Printf("Migration %s applied successfully", migration.Version)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (mr *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`
	_, err := mr.db.Exec(query)
	return err
}

// getMigrationFiles gets all migration files from the directory
func (mr *MigrationRunner) getMigrationFiles(migrationsDir string) ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Extract version and name from filename
		// Format: 001_create_table_name.sql
		filename := d.Name()
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename format: %s", filename)
		}

		version := parts[0]
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Path:    path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations gets the list of applied migrations
func (mr *MigrationRunner) getAppliedMigrations() ([]string, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := mr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied = append(applied, version)
	}

	return applied, nil
}

// isMigrationApplied checks if a migration has been applied
func (mr *MigrationRunner) isMigrationApplied(version string, appliedMigrations []string) bool {
	for _, applied := range appliedMigrations {
		if applied == version {
			return true
		}
	}
	return false
}

// applyMigration applies a single migration
func (mr *MigrationRunner) applyMigration(migration Migration) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := mr.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration as applied
	recordQuery := `INSERT INTO schema_migrations (version) VALUES ($1)`
	if _, err := tx.Exec(recordQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackMigration rolls back a specific migration (not implemented in this simple version)
func (mr *MigrationRunner) RollbackMigration(version string) error {
	return fmt.Errorf("rollback not implemented in this version")
}

// GetMigrationStatus returns the status of all migrations
func (mr *MigrationRunner) GetMigrationStatus(migrationsDir string) (map[string]bool, error) {
	migrations, err := mr.getMigrationFiles(migrationsDir)
	if err != nil {
		return nil, err
	}

	appliedMigrations, err := mr.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	status := make(map[string]bool)
	for _, migration := range migrations {
		status[migration.Version] = mr.isMigrationApplied(migration.Version, appliedMigrations)
	}

	return status, nil
}

func main() {
	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost/threat_detection?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Get migrations directory
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	// Create migration runner
	runner := NewMigrationRunner(db)

	// Check if we should show status
	if len(os.Args) > 1 && os.Args[1] == "status" {
		status, err := runner.GetMigrationStatus(migrationsDir)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		fmt.Println("Migration Status:")
		for version, applied := range status {
			statusStr := "PENDING"
			if applied {
				statusStr = "APPLIED"
			}
			fmt.Printf("  %s: %s\n", version, statusStr)
		}
		return
	}

	// Run migrations
	log.Println("Starting database migrations...")
	if err := runner.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("All migrations completed successfully!")
}
