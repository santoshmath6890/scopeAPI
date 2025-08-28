package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/api-discovery/internal/services"
	"shared/logging"
)

type InventoryHandler struct {
	inventoryService services.InventoryServiceInterface
	logger          logging.Logger
}

func NewInventoryHandler(inventoryService services.InventoryServiceInterface, logger logging.Logger) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		logger:          logger,
	}
}

func (h *InventoryHandler) GetAPIInventory(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Parse filters
	var filters services.InventoryFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		h.logger.Error("Invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.inventoryService.GetAPIInventory(c.Request.Context(), page, limit, filters)
	if err != nil {
		h.logger.Error("Failed to get API inventory", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API inventory"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *InventoryHandler) GetAPIDetails(c *gin.Context) {
	apiID := c.Param("id")
	if apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	}

	apiDetails, err := h.inventoryService.GetAPIDetails(c.Request.Context(), apiID)
	if err != nil {
		h.logger.Error("Failed to get API details", "error", err, "api_id", apiID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API details"})
		return
	}

	c.JSON(http.StatusOK, apiDetails)
}

func (h *InventoryHandler) UpdateAPITags(c *gin.Context) {
	apiID := c.Param("id")
	if apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	}

	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.inventoryService.UpdateAPITags(c.Request.Context(), apiID, req.Tags)
	if err != nil {
		h.logger.Error("Failed to update API tags", "error", err, "api_id", apiID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API tags"})
		return
	}

	h.logger.Info("API tags updated", "api_id", apiID, "tags", req.Tags)
	c.JSON(http.StatusOK, gin.H{"message": "API tags updated successfully"})
}

func (h *InventoryHandler) DeleteAPI(c *gin.Context) {
	apiID := c.Param("id")
	if apiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API ID is required"})
		return
	}

	err := h.inventoryService.DeleteAPI(c.Request.Context(), apiID)
	if err != nil {
		h.logger.Error("Failed to delete API", "error", err, "api_id", apiID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API"})
		return
	}

	h.logger.Info("API deleted", "api_id", apiID)
	c.JSON(http.StatusOK, gin.H{"message": "API deleted successfully"})
}

func (h *InventoryHandler) GetAPIStatistics(c *gin.Context) {
	stats, err := h.inventoryService.GetAPIStatistics(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get API statistics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
