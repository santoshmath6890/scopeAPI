package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

// EnvoyHandler handles HTTP requests for Envoy gateway integration
type EnvoyHandler struct {
	envoyService services.EnvoyIntegrationService
}

// NewEnvoyHandler creates a new EnvoyHandler instance
func NewEnvoyHandler(envoyService services.EnvoyIntegrationService) *EnvoyHandler {
	return &EnvoyHandler{
		envoyService: envoyService,
	}
}

// GetStatus retrieves the status of Envoy gateway
func (h *EnvoyHandler) GetStatus(c *gin.Context) {
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

	status, err := h.envoyService.GetStatus(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetClusters retrieves all clusters from Envoy
func (h *EnvoyHandler) GetClusters(c *gin.Context) {
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

	clusters, err := h.envoyService.GetClusters(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"clusters": clusters,
		"count":    len(clusters),
	})
}

// GetListeners retrieves all listeners from Envoy
func (h *EnvoyHandler) GetListeners(c *gin.Context) {
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

	listeners, err := h.envoyService.GetListeners(c.Request.Context(), integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"listeners": listeners,
		"count":     len(listeners),
	})
}

// GetFilters retrieves all filters from Envoy
func (h *EnvoyHandler) GetFilters(c *gin.Context) {
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

	filters, err := h.envoyService.GetFilters(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(filters),
	})
}

// CreateFilter creates a new filter in Envoy
func (h *EnvoyHandler) CreateFilter(c *gin.Context) {
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

	var filterData map[string]interface{}
	if err := c.ShouldBindJSON(&filterData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	filter, err := h.envoyService.CreateFilter(c.Request.Context(), integrationID, filterData)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Filter created successfully",
		"filter":  filter,
	})
}

// UpdateFilter updates an existing filter in Envoy
func (h *EnvoyHandler) UpdateFilter(c *gin.Context) {
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

	// Extract filter ID from URL parameter
	filterID := c.Param("id")
	if filterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter ID is required"})
		return
	}

	var filterData map[string]interface{}
	if err := c.ShouldBindJSON(&filterData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	filter, err := h.envoyService.UpdateFilter(c.Request.Context(), integrationID, filterID, filterData)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"message": "Filter updated successfully",
		"filter":  filter,
	})
}

// DeleteFilter deletes a filter from Envoy
func (h *EnvoyHandler) DeleteFilter(c *gin.Context) {
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

	// Extract filter ID from URL parameter
	filterID := c.Param("id")
	if filterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter ID is required"})
		return
	}

	if err := h.envoyService.DeleteFilter(c.Request.Context(), integrationID, filterID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Filter deleted successfully"})
}

// SyncConfiguration synchronizes Envoy configuration
func (h *EnvoyHandler) SyncConfiguration(c *gin.Context) {
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

	if err := h.envoyService.SyncConfiguration(c.Request.Context(), integrationID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Envoy configuration synchronized successfully"})
} 