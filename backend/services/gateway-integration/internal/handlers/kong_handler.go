package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
)

// KongHandler handles HTTP requests for Kong gateway integration
type KongHandler struct {
	kongService *services.KongIntegrationService
}

// NewKongHandler creates a new KongHandler instance
func NewKongHandler(kongService *services.KongIntegrationService) *KongHandler {
	return &KongHandler{
		kongService: kongService,
	}
}

// GetStatus retrieves the status of Kong gateway
func (h *KongHandler) GetStatus(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	status, err := h.kongService.GetStatus(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetServices retrieves all services from Kong
func (h *KongHandler) GetServices(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	services, err := h.kongService.GetServices(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"count":    len(services),
	})
}

// GetRoutes retrieves all routes from Kong
func (h *KongHandler) GetRoutes(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	routes, err := h.kongService.GetRoutes(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"routes": routes,
		"count":  len(routes),
	})
}

// GetPlugins retrieves all plugins from Kong
func (h *KongHandler) GetPlugins(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	plugins, err := h.kongService.GetPlugins(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"count":   len(plugins),
	})
}

// CreatePlugin creates a new plugin in Kong
func (h *KongHandler) CreatePlugin(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	var plugin models.KongPlugin
	if err := c.ShouldBindJSON(&plugin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	createdPlugin, err := h.kongService.CreatePlugin(c.Request.Context(), integrationID, &plugin)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plugin created successfully",
		"plugin":  createdPlugin,
	})
}

// UpdatePlugin updates an existing plugin in Kong
func (h *KongHandler) UpdatePlugin(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	// Extract plugin ID from URL parameter
	pluginID := c.Param("id")
	if pluginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plugin ID is required"})
		return
	}

	var plugin models.KongPlugin
	if err := c.ShouldBindJSON(&plugin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	updatedPlugin, err := h.kongService.UpdatePlugin(c.Request.Context(), integrationID, pluginID, &plugin)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin updated successfully",
		"plugin":  updatedPlugin,
	})
}

// DeletePlugin deletes a plugin from Kong
func (h *KongHandler) DeletePlugin(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	// Extract plugin ID from URL parameter
	pluginID := c.Param("id")
	if pluginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plugin ID is required"})
		return
	}

	if err := h.kongService.DeletePlugin(c.Request.Context(), integrationID, pluginID); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin deleted successfully"})
}

// SyncConfiguration synchronizes Kong configuration
func (h *KongHandler) SyncConfiguration(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	result, err := h.kongService.SyncConfiguration(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kong configuration synchronized successfully",
		"result":  result,
	})
}
