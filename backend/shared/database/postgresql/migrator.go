package postgresql

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt time.Time
}

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	config Config
}

// NewMigrator creates a new migration manager
func NewMigrator(config Config) (*Migrator, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Migrator{
		db:     db,
		config: config,
	}, nil
}

// Close closes the database connection
func (m *Migrator) Close() error {
	return m.db.Close()
}

// Init creates the migrations table if it doesn't exist
func (m *Migrator) Init() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns all applied migrations
func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
	query := `SELECT version, name, applied_at FROM schema_migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		err := rows.Scan(&m.Version, &m.Name, &m.AppliedAt)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}
	return migrations, nil
}

// LoadMigrations loads migration files from the migrations directory
func (m *Migrator) LoadMigrations(migrationsDir string) ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		// Parse version and name from filename (e.g., "001_initial_schema.up.sql")
		filename := filepath.Base(path)
		parts := strings.Split(strings.TrimSuffix(filename, ".up.sql"), "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename: %s", filename)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid migration version in filename: %s", filename)
		}

		name := strings.Join(parts[1:], "_")

		// Read up migration
		upSQL, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read up migration %s: %w", path, err)
		}

		// Read down migration
		downPath := strings.Replace(path, ".up.sql", ".down.sql", 1)
		downSQL, err := os.ReadFile(downPath)
		if err != nil {
			// Down migration is optional
			downSQL = []byte("")
		}

		migration := Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		}

		migrations = append(migrations, migration)
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

// Migrate applies all pending migrations
func (m *Migrator) Migrate(migrationsDir string) error {
	// Initialize migrations table
	if err := m.Init(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Load available migrations
	availableMigrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[int]bool)
	for _, m := range appliedMigrations {
		appliedMap[m.Version] = true
	}

	// Apply pending migrations
	for _, migration := range availableMigrations {
		if !appliedMap[migration.Version] {
			fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)
			
			// Start transaction
			tx, err := m.db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction: %w", err)
			}

			// Execute migration
			_, err = tx.Exec(migration.UpSQL)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
			}

			// Record migration as applied
			_, err = tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", 
				migration.Version, migration.Name)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
			}

			// Commit transaction
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
			}

			fmt.Printf("Successfully applied migration %d: %s\n", migration.Version, migration.Name)
		}
	}

	return nil
}

// Rollback rolls back the last migration
func (m *Migrator) Rollback() error {
	// Get applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedMigrations) == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Get the last applied migration
	lastMigration := appliedMigrations[len(appliedMigrations)-1]

	// Load the migration file to get the down SQL
	migrationsDir := "backend/shared/database/postgresql/migrations"
	downPath := filepath.Join(migrationsDir, fmt.Sprintf("%03d_%s.down.sql", lastMigration.Version, lastMigration.Name))
	
	downSQL, err := os.ReadFile(downPath)
	if err != nil {
		return fmt.Errorf("failed to read down migration for %d: %w", lastMigration.Version, err)
	}

	fmt.Printf("Rolling back migration %d: %s\n", lastMigration.Version, lastMigration.Name)

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute down migration
	_, err = tx.Exec(string(downSQL))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute down migration %d: %w", lastMigration.Version, err)
	}

	// Remove migration record
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", lastMigration.Version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", lastMigration.Version, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback for migration %d: %w", lastMigration.Version, err)
	}

	fmt.Printf("Successfully rolled back migration %d: %s\n", lastMigration.Version, lastMigration.Name)
	return nil
}

// Status shows the status of all migrations
func (m *Migrator) Status(migrationsDir string) error {
	// Load available migrations
	availableMigrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[int]Migration)
	for _, m := range appliedMigrations {
		appliedMap[m.Version] = m
	}

	fmt.Printf("\nMigration Status:\n")
	fmt.Printf("%-10s %-30s %-20s %s\n", "Version", "Name", "Applied At", "Status")
	fmt.Printf("%-10s %-30s %-20s %s\n", "-------", "----", "-----------", "------")

	for _, migration := range availableMigrations {
		applied, exists := appliedMap[migration.Version]
		status := "Pending"
		appliedAt := ""

		if exists {
			status = "Applied"
			appliedAt = applied.AppliedAt.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("%-10d %-30s %-20s %s\n", migration.Version, migration.Name, appliedAt, status)
	}

	return nil
} 