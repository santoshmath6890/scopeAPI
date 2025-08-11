package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"scopeapi.local/backend/services/api-discovery/internal/handlers"
	"scopeapi.local/backend/services/api-discovery/internal/repository"
	"scopeapi.local/backend/services/api-discovery/internal/services"
	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/monitoring/health"
)

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize logger
	logger := logging.NewStructuredLogger("api-discovery")

	// Initialize database connection using postgresql package
	config := postgresql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "scopeapi"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	conn, err := postgresql.NewConnection(config)
	var db *sqlx.DB
	if err != nil {
		logger.Warn("Failed to connect to database, starting without database", "error", err)
	} else if conn != nil {
		db = sqlx.NewDb(conn.DB(), "postgres")
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	var discoveryRepo repository.DiscoveryRepositoryInterface
	var inventoryRepo repository.InventoryRepositoryInterface

	if db != nil {
		if err := db.Ping(); err != nil {
			logger.Warn("Database ping failed, starting without database", "error", err)
			discoveryRepo = repository.NewDiscoveryRepository(nil)
			inventoryRepo = repository.NewInventoryRepository(nil)
		} else {
			logger.Info("Database connected successfully")
			discoveryRepo = repository.NewDiscoveryRepository(db)
			inventoryRepo = repository.NewInventoryRepository(db)
		}
	} else {
		discoveryRepo = repository.NewDiscoveryRepository(nil)
		inventoryRepo = repository.NewInventoryRepository(nil)
	}

	// Initialize services
	discoveryService := services.NewDiscoveryService(discoveryRepo, logger)
	inventoryService := services.NewInventoryService(inventoryRepo, logger)
	metadataService := services.NewMetadataService(discoveryRepo, logger)

	// Initialize handlers
	discoveryHandler := handlers.NewDiscoveryHandler(discoveryService, logger)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, logger)
	endpointHandler := handlers.NewEndpointHandler(discoveryService, metadataService, logger)

	// Setup router
	router := gin.Default()
	
	// Health check endpoints
	router.GET("/health", health.HealthCheckHandler)
	router.GET("/ready", health.ReadinessCheckHandler)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/discovery/scan", discoveryHandler.StartDiscovery)
		v1.GET("/discovery/status/:id", discoveryHandler.GetDiscoveryStatus)
		v1.GET("/inventory/apis", inventoryHandler.GetAPIInventory)
		v1.GET("/inventory/apis/:id", inventoryHandler.GetAPIDetails)
		v1.POST("/endpoints/analyze", endpointHandler.AnalyzeEndpoint)
		v1.GET("/endpoints/:id/metadata", endpointHandler.GetEndpointMetadata)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	logger.Info("API Discovery Service started on port " + os.Getenv("SERVER_PORT"))

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
