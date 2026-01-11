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
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"scopeapi.local/backend/services/gateway-integration/internal/handlers"
	"scopeapi.local/backend/services/gateway-integration/internal/repository"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
	"scopeapi.local/backend/shared/auth/jwt"
	"scopeapi.local/backend/shared/database/postgresql"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
	"scopeapi.local/backend/shared/monitoring/metrics"
	"scopeapi.local/backend/shared/utils/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.NewStructuredLogger("gateway-integration")

	// Initialize database connection
	dbConn, err := postgresql.NewConnection(cfg.Database.PostgreSQL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer dbConn.Close()

	// Initialize sqlx.DB
	db := sqlx.NewDb(dbConn.DB(), "postgres")

	// Initialize Kafka producer/consumer
	kafkaProducer, err := kafka.NewProducer(cfg.Messaging.Kafka)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", "error", err)
	}
	defer kafkaProducer.Close()

	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging.Kafka, []string{"gateway_events", "security_events"})
	if err != nil {
		logger.Fatal("Failed to initialize Kafka consumer", "error", err)
	}
	defer kafkaConsumer.Close()

	// Initialize metrics
	metricsCollector := metrics.NewPrometheusCollector("gateway_integration")

	// Initialize repositories
	integrationRepo := repository.NewIntegrationRepository(db)
	configRepo := repository.NewConfigRepository(db)

	// Initialize services
	integrationService := services.NewIntegrationService(integrationRepo, kafkaProducer, logger)
	configService := services.NewConfigService(configRepo, logger)
	kongService := services.NewKongIntegrationService(integrationService, logger)
	nginxService := services.NewNginxIntegrationService(integrationService, logger)
	traefikService := services.NewTraefikIntegrationService(integrationService, logger)
	envoyService := services.NewEnvoyIntegrationService(integrationService, logger)
	haproxyService := services.NewHAProxyIntegrationService(integrationService, logger)

	// Initialize handlers
	integrationHandler := handlers.NewIntegrationHandler(integrationService)
	configHandler := handlers.NewConfigHandler(configService)
	kongHandler := handlers.NewKongHandler(kongService)
	nginxHandler := handlers.NewNginxHandler(nginxService)
	traefikHandler := handlers.NewTraefikHandler(traefikService)
	envoyHandler := handlers.NewEnvoyHandler(envoyService)
	haproxyHandler := handlers.NewHAProxyHandler(haproxyService)

	// Initialize JWT middleware
	jwtMiddleware := jwt.NewJWTMiddleware(cfg.Auth.JWTSecret)

	// Setup router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "gateway-integration"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metricsCollector.Handler()))

	// API routes
	v1 := router.Group("/api/v1")
	v1.Use(jwtMiddleware.AuthMiddleware())
	{
		// Integration management routes
		integrations := v1.Group("/integrations")
		{
			integrations.GET("", integrationHandler.GetIntegrations)
			integrations.GET("/stats", integrationHandler.GetIntegrationStats)
			integrations.GET("/:id", integrationHandler.GetIntegration)
			integrations.POST("", integrationHandler.CreateIntegration)
			integrations.PUT("/:id", integrationHandler.UpdateIntegration)
			integrations.DELETE("/:id", integrationHandler.DeleteIntegration)
			integrations.POST("/:id/test", integrationHandler.TestIntegration)
			integrations.POST("/:id/sync", integrationHandler.SyncIntegration)
		}

		// Configuration management routes
		configs := v1.Group("/configs")
		{
			configs.GET("", configHandler.GetConfigs)
			configs.GET("/:id", configHandler.GetConfig)
			configs.POST("", configHandler.CreateConfig)
			configs.PUT("/:id", configHandler.UpdateConfig)
			configs.DELETE("/:id", configHandler.DeleteConfig)
			configs.POST("/:id/validate", configHandler.ValidateConfig)
			configs.POST("/:id/deploy", configHandler.DeployConfig)
		}

		// Kong integration routes
		kong := v1.Group("/kong")
		{
			kong.GET("/status", kongHandler.GetStatus)
			kong.GET("/services", kongHandler.GetServices)
			kong.GET("/routes", kongHandler.GetRoutes)
			kong.GET("/plugins", kongHandler.GetPlugins)
			kong.POST("/plugins", kongHandler.CreatePlugin)
			kong.PUT("/plugins/:id", kongHandler.UpdatePlugin)
			kong.DELETE("/plugins/:id", kongHandler.DeletePlugin)
			kong.POST("/sync", kongHandler.SyncConfiguration)
		}

		// NGINX integration routes
		nginx := v1.Group("/nginx")
		{
			nginx.GET("/status", nginxHandler.GetStatus)
			nginx.GET("/config", nginxHandler.GetConfig)
			nginx.POST("/config", nginxHandler.UpdateConfig)
			nginx.POST("/reload", nginxHandler.ReloadConfig)
			nginx.GET("/upstreams", nginxHandler.GetUpstreams)
			nginx.POST("/upstreams", nginxHandler.UpdateUpstream)
			nginx.POST("/sync", nginxHandler.SyncConfiguration)
		}

		// Traefik integration routes
		traefik := v1.Group("/traefik")
		{
			traefik.GET("/status", traefikHandler.GetStatus)
			traefik.GET("/providers", traefikHandler.GetProviders)
			traefik.GET("/middlewares", traefikHandler.GetMiddlewares)
			traefik.POST("/middlewares", traefikHandler.CreateMiddleware)
			traefik.PUT("/middlewares/:id", traefikHandler.UpdateMiddleware)
			traefik.DELETE("/middlewares/:id", traefikHandler.DeleteMiddleware)
			traefik.POST("/sync", traefikHandler.SyncConfiguration)
		}

		// Envoy integration routes
		envoy := v1.Group("/envoy")
		{
			envoy.GET("/status", envoyHandler.GetStatus)
			envoy.GET("/clusters", envoyHandler.GetClusters)
			envoy.GET("/listeners", envoyHandler.GetListeners)
			envoy.GET("/filters", envoyHandler.GetFilters)
			envoy.POST("/filters", envoyHandler.CreateFilter)
			envoy.PUT("/filters/:id", envoyHandler.UpdateFilter)
			envoy.DELETE("/filters/:id", envoyHandler.DeleteFilter)
			envoy.POST("/sync", envoyHandler.SyncConfiguration)
		}

		// HAProxy integration routes
		haproxy := v1.Group("/haproxy")
		{
			haproxy.GET("/status", haproxyHandler.GetStatus)
			haproxy.GET("/config", haproxyHandler.GetConfig)
			haproxy.POST("/config", haproxyHandler.UpdateConfig)
			haproxy.POST("/reload", haproxyHandler.ReloadConfig)
			haproxy.GET("/backends", haproxyHandler.GetBackends)
			haproxy.POST("/backends", haproxyHandler.UpdateBackend)
			haproxy.POST("/sync", haproxyHandler.SyncConfiguration)
		}
	}

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer for gateway events
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
					go processMessage(message, integrationService, logger)
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
		logger.Info("Starting gateway integration service", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gateway integration service...")

	// Cancel context to stop background services
	cancel()

	// Shutdown HTTP server
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Gateway integration service stopped")
}

func processMessage(message kafka.Message, integrationService services.IntegrationServiceInterface, logger logging.Logger) {
	ctx := context.Background()

	switch message.Topic {
	case "gateway_events":
		// Process gateway events
		if err := integrationService.ProcessGatewayEvent(ctx, message.Value); err != nil {
			logger.Error("Failed to process gateway event", "error", err)
		}

	case "security_events":
		// Process security events for gateway integration
		if err := integrationService.ProcessSecurityEvent(ctx, message.Value); err != nil {
			logger.Error("Failed to process security event", "error", err)
		}
	}
}
