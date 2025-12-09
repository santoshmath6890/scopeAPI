package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	
	"scopeapi/admin-console/internal/services"
)

func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Load configuration based on environment
	configEnv := os.Getenv("ENVIRONMENT")
	if configEnv == "" {
		configEnv = "development"
	}
	
	viper.SetConfigName("config." + configEnv)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"*"})
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)

	// Override with environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("No config file found for environment %s, using defaults", configEnv)
	}

	// Set Gin mode
	if viper.GetString("environment") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = viper.GetStringSlice("cors.allowed_origins")
	corsConfig.AllowMethods = viper.GetStringSlice("cors.allowed_methods")
	corsConfig.AllowHeaders = viper.GetStringSlice("cors.allowed_headers")
	router.Use(cors.New(corsConfig))

	// Initialize service discovery
	serviceDiscovery := services.NewServiceDiscovery()

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "admin-console",
				"version": "1.0.0",
			})
		})

		// Service status
		api.GET("/services/status", serviceDiscovery.ServiceStatusHandler())

		// Admin console specific APIs
		api.GET("/dashboard/stats", getDashboardStats)
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)
	}

	// Serve static files (Angular app)
	staticPath := viper.GetString("static.path")
	if staticPath == "" {
		staticPath = "./dist"
	}

	// Check if static files exist
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		logrus.Warnf("Static files directory %s does not exist", staticPath)
	} else {
		router.Use(static.Serve("/", static.LocalFile(staticPath, false)))
		
		// Handle Angular routing - serve index.html for all non-API routes
		router.NoRoute(func(c *gin.Context) {
			// Don't serve index.html for API routes
			if filepath.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
				return
			}
			
			indexPath := filepath.Join(staticPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Application not built"})
			}
		})
	}

	// Start server
	port := viper.GetString("server.port")
	host := viper.GetString("server.host")
	addr := host + ":" + port

	logrus.Infof("Starting admin console service on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// API handlers
func getDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_users": 150,
		"active_sessions": 45,
		"api_endpoints": 23,
		"threats_blocked": 156,
		"system_health": "good",
	})
}

func getUsers(c *gin.Context) {
	users := []gin.H{
		{"id": 1, "name": "Admin User", "email": "admin@scopeapi.com", "role": "admin"},
		{"id": 2, "name": "John Doe", "email": "john@example.com", "role": "user"},
		{"id": 3, "name": "Jane Smith", "email": "jane@example.com", "role": "user"},
	}
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var user struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock user creation
	newUser := gin.H{
		"id":    4,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}

	c.JSON(http.StatusCreated, newUser)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock user update
	updatedUser := gin.H{
		"id":    id,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}

	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	
	// Mock user deletion
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
} 