package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/data-protection/internal/config"
	"scopeapi.local/backend/services/data-protection/internal/handlers"
	"scopeapi.local/backend/services/data-protection/internal/repository"
	"scopeapi.local/backend/services/data-protection/internal/services"
	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.NewStructuredLogger("data-protection")

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

	// Initialize Kafka producer
	kafkaConfig := kafka.Config{
		Brokers: cfg.Messaging.Kafka.Brokers,
	}
	kafkaProducer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", "error", err)
	}
	defer kafkaProducer.Close()

	// Initialize metrics
	metricsCollector := metrics.NewPrometheusCollector("data_protection")

	// Initialize repositories
	classificationRepo := repository.NewClassificationRepository(db)
	piiRepo := repository.NewPIIRepository(db)
	complianceRepo := repository.NewComplianceRepository(db)

	// Initialize services
	dataClassificationService := services.NewDataClassificationService(classificationRepo, kafkaProducer, logger)
	piiDetectionService := services.NewPIIDetectionService(piiRepo, kafkaProducer, logger)
	complianceService := services.NewComplianceService(complianceRepo, kafkaProducer, logger)
	riskScoringService := services.NewRiskScoringService(classificationRepo, piiRepo, complianceRepo, kafkaProducer, logger)

	// Initialize handlers
	classificationHandler := handlers.NewClassificationHandler(dataClassificationService, logger)
	piiHandler := handlers.NewPIIHandler(piiDetectionService, logger)
	complianceHandler := handlers.NewComplianceHandler(complianceService, logger)
	riskHandler := handlers.NewRiskHandler(riskScoringService, logger)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "data-protection",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		metricsCollector.ServeHTTP(c.Writer, c.Request)
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Data Classification routes
		classification := api.Group("/classification")
		{
			classification.POST("/classify", classificationHandler.ClassifyData)
			classification.GET("/rules", classificationHandler.GetClassificationRules)
			classification.POST("/rules", classificationHandler.CreateClassificationRule)
			classification.PUT("/rules/:id", classificationHandler.UpdateClassificationRule)
			classification.DELETE("/rules/:id", classificationHandler.DeleteClassificationRule)
			classification.GET("/rules/:id", classificationHandler.GetClassificationRule)
			classification.POST("/rules/:id/enable", classificationHandler.EnableClassificationRule)
			classification.POST("/rules/:id/disable", classificationHandler.DisableClassificationRule)
			classification.GET("/report", classificationHandler.GetClassificationReport)
		}

		// PII Detection routes
		pii := api.Group("/pii")
		{
			pii.POST("/detect", piiHandler.DetectPII)
			pii.GET("/patterns", piiHandler.GetPIIPatterns)
			pii.POST("/patterns", piiHandler.CreatePIIPattern)
			pii.PUT("/patterns/:id", piiHandler.UpdatePIIPattern)
			pii.DELETE("/patterns/:id", piiHandler.DeletePIIPattern)
			pii.GET("/scan", piiHandler.ScanForPII)
			pii.GET("/report", piiHandler.GetPIIReport)
		}

		// Compliance routes
		compliance := api.Group("/compliance")
		{
			compliance.GET("/frameworks", complianceHandler.GetComplianceFrameworks)
			compliance.GET("/frameworks/:id", complianceHandler.GetComplianceFramework)
			compliance.POST("/frameworks", complianceHandler.CreateComplianceFramework)
			compliance.PUT("/frameworks/:id", complianceHandler.UpdateComplianceFramework)
			compliance.DELETE("/frameworks/:id", complianceHandler.DeleteComplianceFramework)
			compliance.GET("/reports", complianceHandler.GetComplianceReports)
			compliance.POST("/reports", complianceHandler.CreateComplianceReport)
			compliance.GET("/reports/:id", complianceHandler.GetComplianceReport)
			compliance.PUT("/reports/:id", complianceHandler.UpdateComplianceReport)
			compliance.DELETE("/reports/:id", complianceHandler.DeleteComplianceReport)
			compliance.GET("/audit", complianceHandler.GetAuditLog)
		}

		// Risk Assessment routes
		risk := api.Group("/risk")
		{
			risk.POST("/assess", riskHandler.AssessRisk)
			risk.GET("/scores", riskHandler.GetRiskScores)
			risk.GET("/scores/:id", riskHandler.GetRiskScore)
			risk.POST("/mitigation", riskHandler.CreateMitigationPlan)
			risk.PUT("/mitigation/:id", riskHandler.UpdateMitigationPlan)
			risk.GET("/mitigation/:id", riskHandler.GetMitigationPlan)
			risk.DELETE("/mitigation/:id", riskHandler.DeleteMitigationPlan)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	logger.Info("Data Protection Service started", "port", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Data Protection Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Data Protection Service exited")
}
