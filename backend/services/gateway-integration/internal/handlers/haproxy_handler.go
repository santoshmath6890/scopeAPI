package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

// HAProxyHandler handles HTTP requests for HAProxy gateway integration
type HAProxyHandler struct {
	haproxyService services.HAProxyIntegrationService
	
}

// NewHAProxyHandler creates a new HAProxyHandler instance
func NewHAProxyHandler(haproxyService services.HAProxyIntegrationService, ) *HAProxyHandler {
	return &HAProxyHandler{
		haproxyService: haproxyService,
		metrics:        metrics,
	}
}

// GetStatus retrieves the status of HAProxy gateway
func (h *HAProxyHandler) GetStatus(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	status, err := h.haproxyService.GetStatus(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetConfig retrieves the HAProxy configuration
func (h *HAProxyHandler) GetConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	config, err := h.haproxyService.GetConfig(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"config": config})
}

// UpdateConfig updates the HAProxy configuration
func (h *HAProxyHandler) UpdateConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.haproxyService.UpdateConfig(c.Request.Context(), integrationID, configData); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "HAProxy configuration updated successfully"})
}

// ReloadConfig reloads the HAProxy configuration
func (h *HAProxyHandler) ReloadConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	if err := h.haproxyService.ReloadConfig(c.Request.Context(), integrationID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "HAProxy configuration reloaded successfully"})
}

// GetBackends retrieves all backend configurations from HAProxy
func (h *HAProxyHandler) GetBackends(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	backends, err := h.haproxyService.GetBackends(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"backends": backends,
		"count":    len(backends),
	})
}

// UpdateBackend updates a backend configuration in HAProxy
func (h *HAProxyHandler) UpdateBackend(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	var backendData map[string]interface{}
	if err := c.ShouldBindJSON(&backendData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.haproxyService.UpdateBackend(c.Request.Context(), integrationID, backendData); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Backend configuration updated successfully"})
}

// SyncConfiguration synchronizes HAProxy configuration
func (h *HAProxyHandler) SyncConfiguration(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration_id parameter"})
		return
	}

	if err := h.haproxyService.SyncConfiguration(c.Request.Context(), integrationID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "HAProxy configuration synchronized successfully"})
} 