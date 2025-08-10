package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Messaging MessagingConfig `mapstructure:"messaging"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Metrics   MetricsConfig   `mapstructure:"metrics"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `mapstructure:"postgresql"`
}

type PostgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MessagingConfig struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Brokers     []string `mapstructure:"brokers"`
	TopicPrefix string   `mapstructure:"topic_prefix"`
}

type AuthConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("server.port", "8082")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.postgresql.host", "localhost")
	viper.SetDefault("database.postgresql.port", "5432")
	viper.SetDefault("database.postgresql.user", "postgres")
	viper.SetDefault("database.postgresql.password", "password")
	viper.SetDefault("database.postgresql.database", "scopeapi")
	viper.SetDefault("database.postgresql.sslmode", "disable")
	viper.SetDefault("messaging.kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("messaging.kafka.topic_prefix", "scopeapi")
	viper.SetDefault("auth.jwt.secret", "your-secret-key")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")

	// Read from environment variables
	viper.AutomaticEnv()

	// Read from config file if it exists
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/threat-detection.yaml"
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults and environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables
	loadFromEnv(&config)

	return &config, nil
}

func loadFromEnv(config *Config) {
	// Server configuration
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if port := os.Getenv("THREAT_DETECTION_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	
	// Database configuration
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.PostgreSQL.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.PostgreSQL.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.PostgreSQL.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.PostgreSQL.Password = password
	}
	if database := os.Getenv("DB_NAME"); database != "" {
		config.Database.PostgreSQL.Database = database
	}
	
	// Kafka configuration
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		config.Messaging.Kafka.Brokers = []string{brokers}
	}
	if topicPrefix := os.Getenv("KAFKA_TOPIC_PREFIX"); topicPrefix != "" {
		config.Messaging.Kafka.TopicPrefix = topicPrefix
	}
} 