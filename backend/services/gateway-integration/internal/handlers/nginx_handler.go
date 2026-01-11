package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
)

// NginxHandler handles HTTP requests for NGINX gateway integration
type NginxHandler struct {
	nginxService *services.NginxIntegrationService
}

// NewNginxHandler creates a new NginxHandler instance
func NewNginxHandler(nginxService *services.NginxIntegrationService) *NginxHandler {
	return &NginxHandler{
		nginxService: nginxService,
	}
}

// GetStatus retrieves the status of NGINX gateway
func (h *NginxHandler) GetStatus(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	status, err := h.nginxService.GetStatus(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetConfig retrieves the NGINX configuration
func (h *NginxHandler) GetConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	config, err := h.nginxService.GetConfig(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}

// UpdateConfig updates the NGINX configuration
func (h *NginxHandler) UpdateConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	var config models.NginxConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.nginxService.UpdateConfig(c.Request.Context(), integrationID, &config); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NGINX configuration updated successfully"})
}

// ReloadConfig reloads the NGINX configuration
func (h *NginxHandler) ReloadConfig(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	if err := h.nginxService.ReloadConfig(c.Request.Context(), integrationID); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NGINX configuration reloaded successfully"})
}

// GetUpstreams retrieves all upstream configurations from NGINX
func (h *NginxHandler) GetUpstreams(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	upstreams, err := h.nginxService.GetUpstreams(c.Request.Context(), integrationID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upstreams": upstreams,
		"count":     len(upstreams),
	})
}

// UpdateUpstream updates an upstream configuration in NGINX
func (h *NginxHandler) UpdateUpstream(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	var upstream models.NginxUpstream
	if err := c.ShouldBindJSON(&upstream); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.nginxService.UpdateUpstream(c.Request.Context(), integrationID, &upstream); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upstream configuration updated successfully"})
}

// SyncConfiguration synchronizes NGINX configuration
func (h *NginxHandler) SyncConfiguration(c *gin.Context) {
	// Extract integration ID from query parameter
	integrationIDStr := c.Query("integration_id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration_id parameter is required"})
		return
	}

	integrationID := integrationIDStr

	result, err := h.nginxService.SyncConfiguration(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "NGINX configuration synchronized successfully",
		"result":  result,
	})
}
