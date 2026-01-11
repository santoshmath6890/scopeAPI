package config

import (
	"os"

	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/messaging/kafka"
)

// Config represents the application configuration
type Config struct {
	Database  DatabaseConfig
	Messaging MessagingConfig
	Server    ServerConfig
	Auth      AuthConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	PostgreSQL postgresql.Config
}

// MessagingConfig holds messaging configuration
type MessagingConfig struct {
	Kafka kafka.Config
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Simple mock implementation for now
	return &Config{
		Database: DatabaseConfig{
			PostgreSQL: postgresql.Config{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				DBName:   getEnv("DB_NAME", "scopeapi"),
				SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			},
		},
		Messaging: MessagingConfig{
			Kafka: kafka.Config{
				Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			},
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
