package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

// ConfigHandler handles HTTP requests for configuration management
type ConfigHandler struct {
	configService services.ConfigService
}

// NewConfigHandler creates a new ConfigHandler instance
func NewConfigHandler(configService services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetConfigs retrieves all configurations for a given integration
func (h *ConfigHandler) GetConfigs(c *gin.Context) {
	// Extract query parameters
	integrationIDStr := c.Query("integration_id")
	configType := c.Query("config_type")

	var integrationID int64
	var err error
	if integrationIDStr != "" {
		}
	}

	// Get configurations
	configs, err := h.configService.GetConfigs(c.Request.Context(), integrationID, configType)
	}

	
	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
		"count":   len(configs),
	})
}

// GetConfig retrieves a specific configuration by ID
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	// Extract ID parameter
	id := c.Param("id")

	// Get configuration
	config, err := h.configService.GetConfig(c.Request.Context(), id)
	}

	if config == nil {
		
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"config": config})
}

// CreateConfig creates a new configuration
func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var config models.GatewayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Create configuration
	if err := h.configService.CreateConfig(c.Request.Context(), &config); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Configuration created successfully",
		"config":  config,
	})
}

// UpdateConfig updates an existing configuration
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	// Extract ID parameter
	id := c.Param("id")
	}

	var config models.GatewayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Set the ID from the URL parameter
	config.ID = id

	// Update configuration
	if err := h.configService.UpdateConfig(c.Request.Context(), &config); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config":  config,
	})
}

// DeleteConfig deletes a configuration
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	// Extract ID parameter
	id := c.Param("id")
	}

	// Delete configuration
	if err := h.configService.DeleteConfig(c.Request.Context(), id); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}

// ValidateConfig validates a configuration without saving it
func (h *ConfigHandler) ValidateConfig(c *gin.Context) {
	var config models.GatewayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Validate configuration
	if err := h.configService.ValidateConfig(c.Request.Context(), &config); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Configuration validation successful"})
}

// DeployConfig deploys a configuration
func (h *ConfigHandler) DeployConfig(c *gin.Context) {
	// Extract ID parameter
	id := c.Param("id")
	}

	// Deploy configuration
	if err := h.configService.DeployConfig(c.Request.Context(), id); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Configuration deployed successfully"})
} 