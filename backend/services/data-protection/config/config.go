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
	Security  SecurityConfig  `mapstructure:"security"`
	Features  FeaturesConfig  `mapstructure:"features"`
}

type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Environment  string `mapstructure:"environment"`
	LogLevel     string `mapstructure:"log_level"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `mapstructure:"postgresql"`
	Redis      RedisConfig      `mapstructure:"redis"`
}

type PostgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConns int    `mapstructure:"max_connections"`
	Timeout  int    `mapstructure:"timeout"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type MessagingConfig struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Brokers        []string           `mapstructure:"brokers"`
	TopicPrefix    string             `mapstructure:"topic_prefix"`
	ProducerConfig ProducerConfig     `mapstructure:"producer"`
	ConsumerConfig ConsumerConfig     `mapstructure:"consumer"`
}

type ProducerConfig struct {
	Acks           string `mapstructure:"acks"`
	Retries        int    `mapstructure:"retries"`
	BatchSize      int    `mapstructure:"batch_size"`
	BatchTimeout   int    `mapstructure:"batch_timeout"`
	Compression    string `mapstructure:"compression"`
	MaxMessageSize int    `mapstructure:"max_message_size"`
}

type ConsumerConfig struct {
	GroupID           string `mapstructure:"group_id"`
	AutoOffsetReset   string `mapstructure:"auto_offset_reset"`
	SessionTimeout    int    `mapstructure:"session_timeout"`
	HeartbeatInterval int    `mapstructure:"heartbeat_interval"`
	MaxPollRecords    int    `mapstructure:"max_poll_records"`
	MaxPollInterval   int    `mapstructure:"max_poll_interval"`
}

type SecurityConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
	Issuer     string `mapstructure:"issuer"`
}

type FeaturesConfig struct {
	DataClassification DataClassificationConfig `mapstructure:"data_classification"`
	PIIDetection       PIIDetectionConfig       `mapstructure:"pii_detection"`
	Compliance         ComplianceConfig         `mapstructure:"compliance"`
	RiskAssessment    RiskAssessmentConfig     `mapstructure:"risk_assessment"`
}

type DataClassificationConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	MLEnabled         bool    `mapstructure:"ml_enabled"`
	ConfidenceThreshold float64 `mapstructure:"confidence_threshold"`
	BatchSize         int     `mapstructure:"batch_size"`
	CacheEnabled      bool    `mapstructure:"cache_enabled"`
}

type PIIDetectionConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	PatternMatching   bool    `mapstructure:"pattern_matching"`
	MLDetection       bool    `mapstructure:"ml_detection"`
	ConfidenceThreshold float64 `mapstructure:"confidence_threshold"`
	ScanDepth         int     `mapstructure:"scan_depth"`
}

type ComplianceConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	GDPR             bool `mapstructure:"gdpr"`
	HIPAA            bool `mapstructure:"hipaa"`
	PCI_DSS          bool `mapstructure:"pci_dss"`
	SOX              bool `mapstructure:"sox"`
	AuditLogging     bool `mapstructure:"audit_logging"`
	AutoReporting    bool `mapstructure:"auto_reporting"`
}

type RiskAssessmentConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	ScoringAlgorithm string  `mapstructure:"scoring_algorithm"`
	RiskThreshold    float64 `mapstructure:"risk_threshold"`
	UpdateFrequency  int     `mapstructure:"update_frequency"`
	MLEnabled        bool    `mapstructure:"ml_enabled"`
}

func LoadConfig() (*Config, error) {
	// Set defaults
	setDefaults()

	// Load from environment variables
	loadFromEnv()

	// Load from config file if exists
	if err := loadFromFile(); err != nil {
		// Continue with environment variables if config file not found
	}

	// Create config struct
	config := &Config{}

	// Bind environment variables
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8084)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.log_level", "info")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// Database defaults
	viper.SetDefault("database.postgresql.host", "localhost")
	viper.SetDefault("database.postgresql.port", 5432)
	viper.SetDefault("database.postgresql.user", "scopeapi")
	viper.SetDefault("database.postgresql.database", "scopeapi")
	viper.SetDefault("database.postgresql.ssl_mode", "disable")
	viper.SetDefault("database.postgresql.max_connections", 10)
	viper.SetDefault("database.postgresql.timeout", 30)

	viper.SetDefault("database.redis.host", "localhost")
	viper.SetDefault("database.redis.port", 6379)
	viper.SetDefault("database.redis.database", 0)
	viper.SetDefault("database.redis.timeout", 30)

	// Kafka defaults
	viper.SetDefault("messaging.kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("messaging.kafka.topic_prefix", "scopeapi.data-protection")
	viper.SetDefault("messaging.kafka.producer.acks", "all")
	viper.SetDefault("messaging.kafka.producer.retries", 3)
	viper.SetDefault("messaging.kafka.producer.batch_size", 16384)
	viper.SetDefault("messaging.kafka.producer.batch_timeout", 100)
	viper.SetDefault("messaging.kafka.producer.compression", "snappy")
	viper.SetDefault("messaging.kafka.producer.max_message_size", 1000000)

	viper.SetDefault("messaging.kafka.consumer.group_id", "data-protection-group")
	viper.SetDefault("messaging.kafka.consumer.auto_offset_reset", "earliest")
	viper.SetDefault("messaging.kafka.consumer.session_timeout", 30000)
	viper.SetDefault("messaging.kafka.consumer.heartbeat_interval", 3000)
	viper.SetDefault("messaging.kafka.consumer.max_poll_records", 500)
	viper.SetDefault("messaging.kafka.consumer.max_poll_interval", 300000)

	// Security defaults
	viper.SetDefault("security.jwt.expiration", 3600)
	viper.SetDefault("security.jwt.issuer", "scopeapi")

	// Features defaults
	viper.SetDefault("features.data_classification.enabled", true)
	viper.SetDefault("features.data_classification.ml_enabled", true)
	viper.SetDefault("features.data_classification.confidence_threshold", 0.8)
	viper.SetDefault("features.data_classification.batch_size", 1000)
	viper.SetDefault("features.data_classification.cache_enabled", true)

	viper.SetDefault("features.pii_detection.enabled", true)
	viper.SetDefault("features.pii_detection.pattern_matching", true)
	viper.SetDefault("features.pii_detection.ml_detection", true)
	viper.SetDefault("features.pii_detection.confidence_threshold", 0.85)
	viper.SetDefault("features.pii_detection.scan_depth", 3)

	viper.SetDefault("features.compliance.enabled", true)
	viper.SetDefault("features.compliance.gdpr", true)
	viper.SetDefault("features.compliance.hipaa", true)
	viper.SetDefault("features.compliance.pci_dss", true)
	viper.SetDefault("features.compliance.sox", true)
	viper.SetDefault("features.compliance.audit_logging", true)
	viper.SetDefault("features.compliance.auto_reporting", true)

	viper.SetDefault("features.risk_assessment.enabled", true)
	viper.SetDefault("features.risk_assessment.scoring_algorithm", "weighted")
	viper.SetDefault("features.risk_assessment.risk_threshold", 0.7)
	viper.SetDefault("features.risk_assessment.update_frequency", 3600)
	viper.SetDefault("features.risk_assessment.ml_enabled", true)
}

func loadFromEnv() {
	// Server
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			viper.Set("server.port", p)
		}
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		viper.Set("server.host", host)
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		viper.Set("server.environment", env)
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		viper.Set("server.log_level", logLevel)
	}

	// Database
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.postgresql.host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			viper.Set("database.postgresql.port", p)
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database.postgresql.user", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.postgresql.password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.postgresql.database", dbName)
	}

	// Kafka
	if kafkaBrokers := os.Getenv("KAFKA_BROKERS"); kafkaBrokers != "" {
		brokers := []string{kafkaBrokers}
		viper.Set("messaging.kafka.brokers", brokers)
	}

	// JWT
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("security.jwt.secret", jwtSecret)
	}
}

func loadFromFile() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

func validateConfig(config *Config) error {
	// Validate server config
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate database config
	if config.Database.PostgreSQL.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if config.Database.PostgreSQL.Port <= 0 || config.Database.PostgreSQL.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.PostgreSQL.Port)
	}
	if config.Database.PostgreSQL.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}
	if config.Database.PostgreSQL.Database == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Validate Kafka config
	if len(config.Messaging.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers cannot be empty")
	}

	// Validate JWT config
	if config.Security.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	return nil
}
