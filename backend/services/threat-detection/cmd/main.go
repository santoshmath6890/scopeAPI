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
	"scopeapi.local/backend/services/threat-detection/internal/handlers"
	"scopeapi.local/backend/services/threat-detection/internal/models"
	"scopeapi.local/backend/services/threat-detection/internal/repository"
	"scopeapi.local/backend/services/threat-detection/internal/services"
	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
	"scopeapi.local/backend/shared/monitoring/metrics"
	"scopeapi.local/backend/services/threat-detection/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.NewStructuredLogger("threat-detection")

	// Initialize database connection
	dbConfig := postgresql.Config{
		Host:     cfg.Database.PostgreSQL.Host,
		Port:     fmt.Sprintf("%d", cfg.Database.PostgreSQL.Port),
		User:     cfg.Database.PostgreSQL.User,
		Password: cfg.Database.PostgreSQL.Password,
		DBName:   cfg.Database.PostgreSQL.Database,
		SSLMode:  cfg.Database.PostgreSQL.SSLMode,
	}
	db, err := postgresql.NewConnection(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize Kafka producer/consumer
	kafkaConfig := kafka.Config{
		Brokers: cfg.Messaging.Kafka.Brokers,
	}
	kafkaProducer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", "error", err)
	}
	defer kafkaProducer.Close()

	kafkaConsumer, err := kafka.NewConsumer(kafkaConfig, []string{"api_traffic", "security_events"})
	if err != nil {
		logger.Fatal("Failed to initialize Kafka consumer", "error", err)
	}
	defer kafkaConsumer.Close()

	// Initialize metrics
	metricsCollector := metrics.NewPrometheusCollector("threat_detection")

	// Initialize repositories
	threatRepo := repository.NewThreatRepository(db)
	patternRepo := repository.NewPatternRepository(db)
	anomalyRepo := repository.NewAnomalyRepository(db)

	// Initialize services
	threatDetectionService := services.NewThreatDetectionService(threatRepo, kafkaProducer, logger)
	anomalyDetectionService := services.NewAnomalyDetectionService(anomalyRepo, kafkaProducer, logger)
	behavioralAnalysisService := services.NewBehavioralAnalysisService(patternRepo, kafkaProducer, logger)
	_ = services.NewSignatureDetectionService(threatRepo, kafkaProducer, logger)

	// Initialize JWT middleware (placeholder for now)
	// jwtMiddleware := jwt.NewMiddleware(cfg.Auth.JWT.Secret)

	// Initialize handlers
	threatHandler := handlers.NewThreatHandler(threatDetectionService, logger)
	anomalyHandler := handlers.NewAnomalyHandler(anomalyDetectionService, logger)
	behavioralHandler := handlers.NewBehavioralHandler(behavioralAnalysisService, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "threat-detection"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metricsCollector.Handler()))

	// API routes (JWT authentication disabled for now)
	v1 := router.Group("/api/v1")
	// v1.Use(jwtMiddleware.AuthMiddleware())
	{
		// Threat detection routes
		threats := v1.Group("/threats")
		{
			threats.GET("", threatHandler.GetThreats)
			threats.GET("/:id", threatHandler.GetThreat)
			threats.POST("/analyze", threatHandler.AnalyzeThreat)
			threats.PUT("/:id/status", threatHandler.UpdateThreatStatus)
			threats.DELETE("/:id", threatHandler.DeleteThreat)
		}

		// Anomaly detection routes
		anomalies := v1.Group("/anomalies")
		{
			anomalies.GET("", anomalyHandler.GetAnomalies)
			anomalies.GET("/:id", anomalyHandler.GetAnomaly)
			anomalies.POST("/detect", anomalyHandler.DetectAnomalies)
			anomalies.PUT("/:id/feedback", anomalyHandler.ProvideFeedback)
		}

		// Behavioral analysis routes
		behavioral := v1.Group("/behavioral")
		{
			behavioral.GET("/patterns", behavioralHandler.GetBehaviorPatterns)
			behavioral.POST("/analyze", behavioralHandler.AnalyzeBehavior)
			behavioral.GET("/baselines", behavioralHandler.GetBaselines)
			behavioral.POST("/baselines", behavioralHandler.CreateBaseline)
		}
	}

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer for real-time threat detection
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				messages, err := kafkaConsumer.Consume(ctx, 100)
				if err != nil {
					logger.Error("Failed to consume Kafka messages", "error", err)
					continue
				}

				for _, message := range messages {
					go processMessage(message, threatDetectionService, anomalyDetectionService, logger)
				}
			}
		}
	}()

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting threat detection service", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down threat detection service...")

	// Cancel context to stop background services
	cancel()

	// Shutdown HTTP server
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Threat detection service stopped")
}

func processMessage(message kafka.Message, threatService services.ThreatDetectionServiceInterface, anomalyService services.AnomalyDetectionServiceInterface, logger logging.Logger) {
	ctx := context.Background()

	switch message.Topic {
	case "api_traffic":
		// Process API traffic for threat detection
		result, err := threatService.AnalyzeTraffic(ctx, message.Value)
		if err != nil {
			logger.Error("Failed to analyze traffic for threats", "error", err)
			return
		}

		if result.ThreatDetected {
			logger.Warn("Threat detected in API traffic", "threat_type", result.ThreatType, "severity", result.Severity)
		}

	case "security_events":
		// Process security events for anomaly detection
		var request models.AnomalyDetectionRequest
		if err := json.Unmarshal(message.Value, &request); err != nil {
			logger.Error("Failed to parse anomaly detection request", "error", err)
			return
		}
		
		result, err := anomalyService.DetectAnomalies(ctx, &request)
		if err != nil {
			logger.Error("Failed to detect anomalies", "error", err)
			return
		}

		for _, anomaly := range result.Anomalies {
			logger.Warn("Anomaly detected", "type", anomaly.Type, "score", anomaly.Score)
		}
	}
}
