package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Messaging  MessagingConfig  `yaml:"messaging"`
	Ingestion  IngestionConfig  `yaml:"ingestion"`
	Parser     ParserConfig     `yaml:"parser"`
	Normalizer NormalizerConfig `yaml:"normalizer"`
	Queue      QueueConfig      `yaml:"queue"`
	Security   SecurityConfig   `yaml:"security"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
	MaxConns int    `yaml:"max_conns"`
}

type MessagingConfig struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers        []string      `yaml:"brokers"`
	TopicPrefix    string        `yaml:"topic_prefix"`
	ProducerConfig ProducerConfig `yaml:"producer"`
	ConsumerConfig ConsumerConfig `yaml:"consumer"`
}

type ProducerConfig struct {
	Acks           string        `yaml:"acks"`
	Retries        int           `yaml:"retries"`
	BatchSize      int           `yaml:"batch_size"`
	BatchTimeout   time.Duration `yaml:"batch_timeout"`
	Compression    string        `yaml:"compression"`
	MaxMessageSize int           `yaml:"max_message_size"`
}

type ConsumerConfig struct {
	GroupID           string        `yaml:"group_id"`
	AutoOffsetReset   string        `yaml:"auto_offset_reset"`
	SessionTimeout    time.Duration `yaml:"session_timeout"`
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
	MaxPollRecords    int           `yaml:"max_poll_records"`
	MaxPollInterval   time.Duration `yaml:"max_poll_interval"`
}

type IngestionConfig struct {
	BatchSize        int           `yaml:"batch_size"`
	BatchTimeout     time.Duration `yaml:"batch_timeout"`
	MaxConcurrency   int           `yaml:"max_concurrency"`
	BufferSize       int           `yaml:"buffer_size"`
	RetryAttempts    int           `yaml:"retry_attempts"`
	RetryDelay       time.Duration `yaml:"retry_delay"`
	Validation       bool          `yaml:"validation"`
	Compression      bool          `yaml:"compression"`
	Topics           TopicsConfig  `yaml:"topics"`
	Formats          FormatsConfig `yaml:"formats"`
}

type TopicsConfig struct {
	APITraffic     string `yaml:"api_traffic"`
	SecurityEvents string `yaml:"security_events"`
	ConfigUpdates  string `yaml:"config_updates"`
	Alerts         string `yaml:"alerts"`
	Metrics        string `yaml:"metrics"`
}

type FormatsConfig struct {
	Supported []string `yaml:"supported"`
	Default   string   `yaml:"default"`
}

type ParserConfig struct {
	MaxPayloadSize int           `yaml:"max_payload_size"`
	Timeout        time.Duration `yaml:"timeout"`
	Formats        []FormatConfig `yaml:"formats"`
}

type FormatConfig struct {
	Name        string   `yaml:"name"`
	Extensions  []string `yaml:"extensions"`
	MimeTypes   []string `yaml:"mime_types"`
	Enabled     bool     `yaml:"enabled"`
	Priority    int      `yaml:"priority"`
	Config      map[string]interface{} `yaml:"config"`
}

type NormalizerConfig struct {
	DefaultSchema string                    `yaml:"default_schema"`
	Schemas       map[string]SchemaConfig   `yaml:"schemas"`
	Transformers  map[string]TransformerConfig `yaml:"transformers"`
}

type SchemaConfig struct {
	Name       string                 `yaml:"name"`
	Version    string                 `yaml:"version"`
	Fields     map[string]FieldConfig `yaml:"fields"`
	Required   []string               `yaml:"required"`
	Validators []ValidatorConfig      `yaml:"validators"`
}

type FieldConfig struct {
	Type     string      `yaml:"type"`
	Required bool        `yaml:"required"`
	Default  interface{} `yaml:"default"`
	Format   string      `yaml:"format"`
	Pattern  string      `yaml:"pattern"`
}

type ValidatorConfig struct {
	Type    string                 `yaml:"type"`
	Config  map[string]interface{} `yaml:"config"`
	Message string                 `yaml:"message"`
}

type TransformerConfig struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

type QueueConfig struct {
	MaxSize       int           `yaml:"max_size"`
	FlushInterval time.Duration `yaml:"flush_interval"`
	RetryPolicy   RetryPolicy   `yaml:"retry_policy"`
}

type RetryPolicy struct {
	MaxAttempts int           `yaml:"max_attempts"`
	Backoff     time.Duration `yaml:"backoff"`
	MaxBackoff  time.Duration `yaml:"max_backoff"`
}

type SecurityConfig struct {
	EnableTLS     bool   `yaml:"enable_tls"`
	CertFile      string `yaml:"cert_file"`
	KeyFile       string `yaml:"key_file"`
	AllowedHosts  []string `yaml:"allowed_hosts"`
	RateLimit     RateLimitConfig `yaml:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled  bool  `yaml:"enabled"`
	Requests int   `yaml:"requests"`
	Window   time.Duration `yaml:"window"`
}

type MonitoringConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Health  HealthConfig  `yaml:"health"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type HealthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.postgresql.host", "localhost")
	viper.SetDefault("database.postgresql.port", "5432")
	viper.SetDefault("database.postgresql.user", "postgres")
	viper.SetDefault("database.postgresql.password", "password")
	viper.SetDefault("database.postgresql.dbname", "scopeapi")
	viper.SetDefault("database.postgresql.sslmode", "disable")
	viper.SetDefault("messaging.kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("messaging.kafka.topic", "api-traffic")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Read from environment variables
	viper.AutomaticEnv()

	// Read from config file if it exists
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/data-ingestion.yaml"
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "config/data-ingestion.yaml"
}

func loadFromEnv(config *Config) {
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
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

func setDefaults(config *Config) {
	// Server defaults
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 60 * time.Second
	}
	
	// Database defaults
	if config.Database.PostgreSQL.Host == "" {
		config.Database.PostgreSQL.Host = "localhost"
	}
	if config.Database.PostgreSQL.Port == 0 {
		config.Database.PostgreSQL.Port = 5432
	}
	if config.Database.PostgreSQL.SSLMode == "" {
		config.Database.PostgreSQL.SSLMode = "disable"
	}
	if config.Database.PostgreSQL.MaxConns == 0 {
		config.Database.PostgreSQL.MaxConns = 10
	}
	
	// Kafka defaults
	if len(config.Messaging.Kafka.Brokers) == 0 {
		config.Messaging.Kafka.Brokers = []string{"localhost:9092"}
	}
	if config.Messaging.Kafka.TopicPrefix == "" {
		config.Messaging.Kafka.TopicPrefix = "scopeapi"
	}
	
	// Ingestion defaults
	if config.Ingestion.BatchSize == 0 {
		config.Ingestion.BatchSize = 1000
	}
	if config.Ingestion.BatchTimeout == 0 {
		config.Ingestion.BatchTimeout = 5 * time.Second
	}
	if config.Ingestion.MaxConcurrency == 0 {
		config.Ingestion.MaxConcurrency = 10
	}
	if config.Ingestion.BufferSize == 0 {
		config.Ingestion.BufferSize = 10000
	}
	if config.Ingestion.RetryAttempts == 0 {
		config.Ingestion.RetryAttempts = 3
	}
	if config.Ingestion.RetryDelay == 0 {
		config.Ingestion.RetryDelay = 1 * time.Second
	}
	
	// Topics defaults
	if config.Ingestion.Topics.APITraffic == "" {
		config.Ingestion.Topics.APITraffic = "api_traffic"
	}
	if config.Ingestion.Topics.SecurityEvents == "" {
		config.Ingestion.Topics.SecurityEvents = "security_events"
	}
	if config.Ingestion.Topics.ConfigUpdates == "" {
		config.Ingestion.Topics.ConfigUpdates = "config_updates"
	}
	if config.Ingestion.Topics.Alerts == "" {
		config.Ingestion.Topics.Alerts = "alerts"
	}
	if config.Ingestion.Topics.Metrics == "" {
		config.Ingestion.Topics.Metrics = "metrics"
	}
	
	// Formats defaults
	if len(config.Ingestion.Formats.Supported) == 0 {
		config.Ingestion.Formats.Supported = []string{"json", "xml", "yaml", "protobuf"}
	}
	if config.Ingestion.Formats.Default == "" {
		config.Ingestion.Formats.Default = "json"
	}
	
	// Parser defaults
	if config.Parser.MaxPayloadSize == 0 {
		config.Parser.MaxPayloadSize = 10 * 1024 * 1024 // 10MB
	}
	if config.Parser.Timeout == 0 {
		config.Parser.Timeout = 30 * time.Second
	}
	
	// Queue defaults
	if config.Queue.MaxSize == 0 {
		config.Queue.MaxSize = 10000
	}
	if config.Queue.FlushInterval == 0 {
		config.Queue.FlushInterval = 1 * time.Second
	}
	if config.Queue.RetryPolicy.MaxAttempts == 0 {
		config.Queue.RetryPolicy.MaxAttempts = 3
	}
	if config.Queue.RetryPolicy.Backoff == 0 {
		config.Queue.RetryPolicy.Backoff = 1 * time.Second
	}
	if config.Queue.RetryPolicy.MaxBackoff == 0 {
		config.Queue.RetryPolicy.MaxBackoff = 60 * time.Second
	}
	
	// Monitoring defaults
	if config.Monitoring.Metrics.Path == "" {
		config.Monitoring.Metrics.Path = "/metrics"
	}
	if config.Monitoring.Health.Path == "" {
		config.Monitoring.Health.Path = "/health"
	}
}

func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	
	if config.Database.PostgreSQL.Host == "" {
		return fmt.Errorf("database host is required")
	}
	
	if len(config.Messaging.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers are required")
	}
	
	if config.Ingestion.BatchSize <= 0 {
		return fmt.Errorf("batch size must be greater than 0")
	}
	
	if config.Ingestion.MaxConcurrency <= 0 {
		return fmt.Errorf("max concurrency must be greater than 0")
	}
	
	return nil
} 