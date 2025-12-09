package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/data-ingestion/internal/config"
	"scopeapi.local/backend/services/data-ingestion/internal/handlers"
	"scopeapi.local/backend/services/data-ingestion/internal/services"
	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
	"scopeapi.local/backend/shared/monitoring/health"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.NewStructuredLogger("data-ingestion")

	// Initialize database connection
	var db *postgresql.Connection
	dbConfig := postgresql.Config{
		Host:     cfg.Database.PostgreSQL.Host,
		Port:     fmt.Sprintf("%d", cfg.Database.PostgreSQL.Port),
		User:     cfg.Database.PostgreSQL.User,
		Password: cfg.Database.PostgreSQL.Password,
		DBName:   cfg.Database.PostgreSQL.Database,
		SSLMode:  cfg.Database.PostgreSQL.SSLMode,
	}
	db, err = postgresql.NewConnection(dbConfig)
	if err != nil {
		logger.Warn("Failed to connect to database, continuing without database", "error", err)
	} else {
		defer db.Close()
		logger.Info("Database connected successfully")
	}

	// Initialize Kafka producer (optional for now)
	var kafkaProducer *kafka.Producer
	if len(cfg.Messaging.Kafka.Brokers) > 0 {
		// Convert config.KafkaConfig to kafka.Config
		kafkaConfig := kafka.Config{
			Brokers:        cfg.Messaging.Kafka.Brokers,
			TopicPrefix:    cfg.Messaging.Kafka.TopicPrefix,
			ProducerConfig: kafka.ProducerConfig{
				Acks:           cfg.Messaging.Kafka.ProducerConfig.Acks,
				Retries:        cfg.Messaging.Kafka.ProducerConfig.Retries,
				BatchSize:      cfg.Messaging.Kafka.ProducerConfig.BatchSize,
				BatchTimeout:   cfg.Messaging.Kafka.ProducerConfig.BatchTimeout,
				Compression:    cfg.Messaging.Kafka.ProducerConfig.Compression,
				MaxMessageSize: cfg.Messaging.Kafka.ProducerConfig.MaxMessageSize,
			},
			ConsumerConfig: kafka.ConsumerConfig{
				GroupID:           cfg.Messaging.Kafka.ConsumerConfig.GroupID,
				AutoOffsetReset:   cfg.Messaging.Kafka.ConsumerConfig.AutoOffsetReset,
				SessionTimeout:    cfg.Messaging.Kafka.ConsumerConfig.SessionTimeout,
				HeartbeatInterval: cfg.Messaging.Kafka.ConsumerConfig.HeartbeatInterval,
				MaxPollRecords:    cfg.Messaging.Kafka.ConsumerConfig.MaxPollRecords,
				MaxPollInterval:   cfg.Messaging.Kafka.ConsumerConfig.MaxPollInterval,
			},
		}
		kafkaProducer, err = kafka.NewProducer(kafkaConfig)
		if err != nil {
			logger.Warn("Failed to initialize Kafka producer, continuing without Kafka", "error", err)
		} else {
			defer kafkaProducer.Close()
		}
	}

	// Initialize Kafka consumer for configuration updates (optional for now)
	var kafkaConsumer *kafka.Consumer
	if kafkaProducer != nil {
		// Convert config.KafkaConfig to kafka.Config
		kafkaConfig := kafka.Config{
			Brokers:        cfg.Messaging.Kafka.Brokers,
			TopicPrefix:    cfg.Messaging.Kafka.TopicPrefix,
			ProducerConfig: kafka.ProducerConfig{
				Acks:           cfg.Messaging.Kafka.ProducerConfig.Acks,
				Retries:        cfg.Messaging.Kafka.ProducerConfig.Retries,
				BatchSize:      cfg.Messaging.Kafka.ProducerConfig.BatchSize,
				BatchTimeout:   cfg.Messaging.Kafka.ProducerConfig.BatchTimeout,
				Compression:    cfg.Messaging.Kafka.ProducerConfig.Compression,
				MaxMessageSize: cfg.Messaging.Kafka.ProducerConfig.MaxMessageSize,
			},
			ConsumerConfig: kafka.ConsumerConfig{
				GroupID:           cfg.Messaging.Kafka.ConsumerConfig.GroupID,
				AutoOffsetReset:   cfg.Messaging.Kafka.ConsumerConfig.AutoOffsetReset,
				SessionTimeout:    cfg.Messaging.Kafka.ConsumerConfig.SessionTimeout,
				HeartbeatInterval: cfg.Messaging.Kafka.ConsumerConfig.HeartbeatInterval,
				MaxPollRecords:    cfg.Messaging.Kafka.ConsumerConfig.MaxPollRecords,
				MaxPollInterval:   cfg.Messaging.Kafka.ConsumerConfig.MaxPollInterval,
			},
		}
		kafkaConsumer, err = kafka.NewConsumer(kafkaConfig, []string{"config_updates"})
		if err != nil {
			logger.Warn("Failed to initialize Kafka consumer, continuing without consumer", "error", err)
		} else {
			defer kafkaConsumer.Close()
		}
	}

	// Initialize metrics collector
	metricsCollector := metrics.NewPrometheusCollector("data_ingestion")

	// Initialize services
	var ingestionService services.DataIngestionServiceInterface
	if kafkaProducer != nil {
		ingestionService = services.NewDataIngestionService(kafkaProducer, logger, cfg)
	} else {
		// Create a mock producer for now
		ingestionService = services.NewDataIngestionService(nil, logger, cfg)
	}
	parserService := services.NewDataParserService(logger, cfg)
	normalizerService := services.NewDataNormalizerService(logger, cfg)
	
	var queueService services.QueueServiceInterface
	if kafkaProducer != nil {
		queueService = services.NewQueueService(kafkaProducer, logger, cfg)
	} else {
		// Create a mock queue service for now
		queueService = services.NewQueueService(nil, logger, cfg)
	}

	// Initialize handlers
	ingestionHandler := handlers.NewIngestionHandler(ingestionService, logger)
	parserHandler := handlers.NewParserHandler(parserService, logger)
	normalizerHandler := handlers.NewNormalizerHandler(normalizerService, logger)
	queueHandler := handlers.NewQueueHandler(queueService, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", health.HealthCheckHandler)
	router.GET("/ready", health.ReadinessCheckHandler)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metricsCollector.Handler()))

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Data ingestion routes
		ingestion := v1.Group("/ingestion")
		{
			ingestion.POST("/traffic", ingestionHandler.IngestTraffic)
			ingestion.POST("/batch", ingestionHandler.IngestBatch)
			ingestion.GET("/status/:id", ingestionHandler.GetIngestionStatus)
			ingestion.GET("/stats", ingestionHandler.GetIngestionStats)
		}

		// Parser routes
		parser := v1.Group("/parser")
		{
			parser.POST("/parse", parserHandler.ParseData)
			parser.GET("/formats", parserHandler.GetSupportedFormats)
			parser.POST("/validate", parserHandler.ValidateFormat)
		}

		// Normalizer routes
		normalizer := v1.Group("/normalizer")
		{
			normalizer.POST("/normalize", normalizerHandler.NormalizeData)
			normalizer.GET("/schemas", normalizerHandler.GetSchemas)
			normalizer.POST("/schema", normalizerHandler.CreateSchema)
		}

		// Queue management routes
		queue := v1.Group("/queue")
		{
			queue.GET("/status", queueHandler.GetQueueStatus)
			queue.POST("/flush", queueHandler.FlushQueue)
			queue.GET("/topics", queueHandler.GetTopics)
		}
	}

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer for configuration updates (if available)
	if kafkaConsumer != nil {
		go func() {
			backoff := time.Second
			maxBackoff := 30 * time.Second
			
			for {
				select {
				case <-ctx.Done():
					return
				default:
					messages, err := kafkaConsumer.Consume(ctx, 10)
					if err != nil {
						logger.Error("Failed to consume Kafka messages", "error", err)
						// Exponential backoff to reduce log spam
						time.Sleep(backoff)
						if backoff < maxBackoff {
							backoff *= 2
						}
						continue
					}
					
					// Reset backoff on successful connection
					backoff = time.Second

					for _, message := range messages {
						go processConfigUpdate(message, ingestionService, logger)
					}
				}
			}
		}()
	}

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting data ingestion service", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down data ingestion service...")

	// Cancel context to stop background services
	cancel()

	// Shutdown HTTP server
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Data ingestion service stopped")
}

func processConfigUpdate(message kafka.Message, ingestionService services.DataIngestionServiceInterface, logger logging.Logger) {
	var configUpdate struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}

	if err := json.Unmarshal([]byte(message.Value), &configUpdate); err != nil {
		logger.Error("Failed to unmarshal config update", "error", err)
		return
	}

	switch configUpdate.Type {
	case "parser_config":
		logger.Info("Received parser configuration update")
		// Update parser configuration
	case "normalizer_config":
		logger.Info("Received normalizer configuration update")
		// Update normalizer configuration
	case "queue_config":
		logger.Info("Received queue configuration update")
		// Update queue configuration
	default:
		logger.Warn("Unknown configuration update type", "type", configUpdate.Type)
	}
} 