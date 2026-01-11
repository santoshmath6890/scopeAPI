package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/services"
)

// IntegrationHandler handles HTTP requests for integration management
type IntegrationHandler struct {
	integrationService services.IntegrationServiceInterface
}

// NewIntegrationHandler creates a new integration handler
func NewIntegrationHandler(integrationService services.IntegrationServiceInterface) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService: integrationService,
	}
}

// GetIntegrations handles GET /api/v1/integrations
func (h *IntegrationHandler) GetIntegrations(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters for filtering
	filters := make(map[string]interface{})
	if gatewayType := c.Query("type"); gatewayType != "" {
		filters["type"] = gatewayType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	// Parse pagination parameters
	page := 1
	limit := 50

	integrations, err := h.integrationService.GetIntegrations(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve integrations",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"integrations": integrations,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(integrations),
		},
	})
}

// GetIntegration handles GET /api/v1/integrations/:id
func (h *IntegrationHandler) GetIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	integrationID := c.Param("id")

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Integration ID is required",
		})
		return
	}

	integration, err := h.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Integration not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"integration": integration,
	})
}

// CreateIntegration handles POST /api/v1/integrations
func (h *IntegrationHandler) CreateIntegration(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Name        string                 `json:"name" binding:"required"`
		Type        models.GatewayType     `json:"type" binding:"required"`
		Config      map[string]interface{} `json:"config"`
		Credentials *models.Credentials    `json:"credentials,omitempty"`
		Endpoints   []models.Endpoint      `json:"endpoints" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate gateway type
	if !isValidGatewayType(request.Type) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gateway type",
		})
		return
	}

	// Create integration model
	integration := &models.Integration{
		Name:        request.Name,
		Type:        request.Type,
		Status:      models.IntegrationStatusPending,
		Config:      request.Config,
		Credentials: request.Credentials,
		Endpoints:   request.Endpoints,
	}

	// Validate integration
	if err := h.validateIntegration(integration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	if err := h.integrationService.CreateIntegration(ctx, integration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create integration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"integration": integration,
		"message":     "Integration created successfully",
	})
}

// UpdateIntegration handles PUT /api/v1/integrations/:id
func (h *IntegrationHandler) UpdateIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	integrationID := c.Param("id")

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Integration ID is required",
		})
		return
	}

	var request struct {
		Name        string                   `json:"name,omitempty"`
		Type        models.GatewayType       `json:"type,omitempty"`
		Status      models.IntegrationStatus `json:"status,omitempty"`
		Config      map[string]interface{}   `json:"config,omitempty"`
		Credentials *models.Credentials      `json:"credentials,omitempty"`
		Endpoints   []models.Endpoint        `json:"endpoints,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Get existing integration
	existingIntegration, err := h.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Integration not found",
			"message": err.Error(),
		})
		return
	}

	// Update fields if provided
	if request.Name != "" {
		existingIntegration.Name = request.Name
	}
	if request.Type != "" {
		if !isValidGatewayType(request.Type) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid gateway type",
			})
			return
		}
		existingIntegration.Type = request.Type
	}
	if request.Status != "" {
		existingIntegration.Status = request.Status
	}
	if request.Config != nil {
		existingIntegration.Config = request.Config
	}
	if request.Credentials != nil {
		existingIntegration.Credentials = request.Credentials
	}
	if request.Endpoints != nil {
		existingIntegration.Endpoints = request.Endpoints
	}

	// Validate integration
	if err := h.validateIntegration(existingIntegration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	if err := h.integrationService.UpdateIntegration(ctx, existingIntegration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update integration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"integration": existingIntegration,
		"message":     "Integration updated successfully",
	})
}

// DeleteIntegration handles DELETE /api/v1/integrations/:id
func (h *IntegrationHandler) DeleteIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	integrationID := c.Param("id")

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Integration ID is required",
		})
		return
	}

	err := h.integrationService.DeleteIntegration(ctx, integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete integration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Integration deleted successfully",
	})
}

// TestIntegration handles POST /api/v1/integrations/:id/test
func (h *IntegrationHandler) TestIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	integrationID := c.Param("id")

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Integration ID is required",
		})
		return
	}

	health, err := h.integrationService.TestIntegration(ctx, integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to test integration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"health":  health,
		"message": "Integration test completed",
	})
}

// SyncIntegration handles POST /api/v1/integrations/:id/sync
func (h *IntegrationHandler) SyncIntegration(c *gin.Context) {
	ctx := c.Request.Context()
	integrationID := c.Param("id")

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Integration ID is required",
		})
		return
	}

	result, err := h.integrationService.SyncIntegration(ctx, integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to sync integration",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  result,
		"message": "Integration sync completed",
	})
}

// GetIntegrationStats handles GET /api/v1/integrations/stats
func (h *IntegrationHandler) GetIntegrationStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.integrationService.GetIntegrationStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve integration statistics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// Helper functions

func isValidGatewayType(gatewayType models.GatewayType) bool {
	validTypes := []models.GatewayType{
		models.GatewayTypeKong,
		models.GatewayTypeNginx,
		models.GatewayTypeTraefik,
		models.GatewayTypeEnvoy,
		models.GatewayTypeHAProxy,
	}

	for _, validType := range validTypes {
		if gatewayType == validType {
			return true
		}
	}
	return false
}

func (h *IntegrationHandler) validateIntegration(integration *models.Integration) error {
	// Validate name
	if integration.Name == "" {
		return fmt.Errorf("integration name is required")
	}

	// Validate type
	if !isValidGatewayType(integration.Type) {
		return fmt.Errorf("invalid gateway type: %s", integration.Type)
	}

	// Validate endpoints
	if len(integration.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint is required")
	}

	for i, endpoint := range integration.Endpoints {
		if endpoint.Name == "" {
			return fmt.Errorf("endpoint %d name is required", i+1)
		}
		if endpoint.URL == "" {
			return fmt.Errorf("endpoint %d URL is required", i+1)
		}
		if endpoint.Protocol == "" {
			return fmt.Errorf("endpoint %d protocol is required", i+1)
		}
		if endpoint.Port <= 0 {
			return fmt.Errorf("endpoint %d port must be positive", i+1)
		}
	}

	return nil
}
