package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
)

// HAProxyHandler handles HTTP requests for HAProxy gateway integration
type HAProxyHandler struct {
	haproxyService *services.HAProxyIntegrationService
}

// NewHAProxyHandler creates a new HAProxyHandler instance
func NewHAProxyHandler(haproxyService *services.HAProxyIntegrationService) *HAProxyHandler {
	return &HAProxyHandler{
		haproxyService: haproxyService,
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

	var request struct {
		Config string `json:"config"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.haproxyService.UpdateConfig(c.Request.Context(), integrationID, request.Config); err != nil {

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

	var backend models.HAProxyBackend
	if err := c.ShouldBindJSON(&backend); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	updatedBackend, err := h.haproxyService.UpdateBackend(c.Request.Context(), integrationID, &backend)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backend updated successfully",
		"backend": updatedBackend,
	})
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

	result, err := h.haproxyService.SyncConfiguration(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "HAProxy configuration synchronized successfully",
		"result":  result,
	})
}
