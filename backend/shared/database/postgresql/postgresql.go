package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config holds PostgreSQL connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN returns the database connection string
func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Connection represents a PostgreSQL database connection
type Connection struct {
	db *sql.DB
}

// NewConnection creates a new PostgreSQL connection
func NewConnection(config Config) (*Connection, error) {
	dsn := config.DSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set search path to scopeapi schema
	if _, err := db.Exec("SET search_path TO scopeapi,public"); err != nil {
		return nil, fmt.Errorf("failed to set search path: %w", err)
	}

	return &Connection{db: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// DB returns the underlying sql.DB instance
func (c *Connection) DB() *sql.DB {
	return c.db
}

// Ping tests the database connection
func (c *Connection) Ping() error {
	return c.db.Ping()
}
