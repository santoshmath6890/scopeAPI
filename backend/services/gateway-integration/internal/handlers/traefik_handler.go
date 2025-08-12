package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
	"scopeapi.local/backend/shared/monitoring/metrics"
)

// TraefikHandler handles HTTP requests for Traefik gateway integration
type TraefikHandler struct {
	traefikService services.TraefikIntegrationService
	
}

// NewTraefikHandler creates a new TraefikHandler instance
func NewTraefikHandler(traefikService services.TraefikIntegrationService, ) *TraefikHandler {
	return &TraefikHandler{
		traefikService: traefikService,
		metrics:        metrics,
	}
}

// GetStatus retrieves the status of Traefik gateway
func (h *TraefikHandler) GetStatus(c *gin.Context) {
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

	status, err := h.traefikService.GetStatus(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetProviders retrieves all providers from Traefik
func (h *TraefikHandler) GetProviders(c *gin.Context) {
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

	providers, err := h.traefikService.GetProviders(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
		"count":     len(providers),
	})
}

// GetMiddlewares retrieves all middlewares from Traefik
func (h *TraefikHandler) GetMiddlewares(c *gin.Context) {
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

	middlewares, err := h.traefikService.GetMiddlewares(c.Request.Context(), integrationID)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"middlewares": middlewares,
		"count":       len(middlewares),
	})
}

// CreateMiddleware creates a new middleware in Traefik
func (h *TraefikHandler) CreateMiddleware(c *gin.Context) {
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

	var middlewareData map[string]interface{}
	if err := c.ShouldBindJSON(&middlewareData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	middleware, err := h.traefikService.CreateMiddleware(c.Request.Context(), integrationID, middlewareData)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Middleware created successfully",
		"middleware": middleware,
	})
}

// UpdateMiddleware updates an existing middleware in Traefik
func (h *TraefikHandler) UpdateMiddleware(c *gin.Context) {
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

	// Extract middleware ID from URL parameter
	middlewareID := c.Param("id")
	if middlewareID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "middleware ID is required"})
		return
	}

	var middlewareData map[string]interface{}
	if err := c.ShouldBindJSON(&middlewareData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	middleware, err := h.traefikService.UpdateMiddleware(c.Request.Context(), integrationID, middlewareID, middlewareData)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Middleware updated successfully",
		"middleware": middleware,
	})
}

// DeleteMiddleware deletes a middleware from Traefik
func (h *TraefikHandler) DeleteMiddleware(c *gin.Context) {
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

	// Extract middleware ID from URL parameter
	middlewareID := c.Param("id")
	if middlewareID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "middleware ID is required"})
		return
	}

	if err := h.traefikService.DeleteMiddleware(c.Request.Context(), integrationID, middlewareID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Middleware deleted successfully"})
}

// SyncConfiguration synchronizes Traefik configuration
func (h *TraefikHandler) SyncConfiguration(c *gin.Context) {
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

	if err := h.traefikService.SyncConfiguration(c.Request.Context(), integrationID); err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"message": "Traefik configuration synchronized successfully"})
} 