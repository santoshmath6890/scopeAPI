package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/api-discovery/internal/models"
	"scopeapi.local/backend/services/api-discovery/internal/services"
	"scopeapi.local/backend/shared/logging"
)

type DiscoveryHandler struct {
	discoveryService services.DiscoveryServiceInterface
	logger          logging.Logger
}

type DiscoveryRequest struct {
	Target      string            `json:"target" binding:"required"`
	Method      string            `json:"method" binding:"required,oneof=passive active"`
	Options     map[string]string `json:"options"`
	Credentials *models.Credentials `json:"credentials,omitempty"`
}

type DiscoveryResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func NewDiscoveryHandler(discoveryService services.DiscoveryServiceInterface, logger logging.Logger) *DiscoveryHandler {
	return &DiscoveryHandler{
		discoveryService: discoveryService,
		logger:          logger,
	}
}

func (h *DiscoveryHandler) StartDiscovery(c *gin.Context) {
	var req DiscoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start discovery process
	discoveryID, err := h.discoveryService.StartDiscovery(c.Request.Context(), &models.DiscoveryConfig{
		Target:      req.Target,
		Method:      req.Method,
		Options:     req.Options,
		Credentials: req.Credentials,
	})

	if err != nil {
		h.logger.Error("Failed to start discovery", "error", err, "target", req.Target)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start discovery"})
		return
	}

	h.logger.Info("Discovery started", "id", discoveryID, "target", req.Target)
	
	c.JSON(http.StatusAccepted, DiscoveryResponse{
		ID:      discoveryID,
		Status:  "started",
		Message: "Discovery process initiated successfully",
	})
}

func (h *DiscoveryHandler) GetDiscoveryStatus(c *gin.Context) {
	discoveryID := c.Param("id")
	if discoveryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discovery ID is required"})
		return
	}

	status, err := h.discoveryService.GetDiscoveryStatus(c.Request.Context(), discoveryID)
	if err != nil {
		h.logger.Error("Failed to get discovery status", "error", err, "id", discoveryID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discovery status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *DiscoveryHandler) GetDiscoveryResults(c *gin.Context) {
	discoveryID := c.Param("id")
	if discoveryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discovery ID is required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	results, err := h.discoveryService.GetDiscoveryResults(c.Request.Context(), discoveryID, page, limit)
	if err != nil {
		h.logger.Error("Failed to get discovery results", "error", err, "id", discoveryID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get discovery results"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DiscoveryHandler) StopDiscovery(c *gin.Context) {
	discoveryID := c.Param("id")
	if discoveryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Discovery ID is required"})
		return
	}

	err := h.discoveryService.StopDiscovery(c.Request.Context(), discoveryID)
	if err != nil {
		h.logger.Error("Failed to stop discovery", "error", err, "id", discoveryID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop discovery"})
		return
	}

	h.logger.Info("Discovery stopped", "id", discoveryID)
	c.JSON(http.StatusOK, gin.H{"message": "Discovery stopped successfully"})
}
