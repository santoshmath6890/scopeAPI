package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/api-discovery/internal/models"
	"scopeapi.local/backend/services/api-discovery/internal/services"
	"scopeapi.local/backend/shared/logging"
)

type EndpointHandler struct {
	discoveryService services.DiscoveryServiceInterface
	metadataService  services.MetadataServiceInterface
	logger          logging.Logger
}

type EndpointAnalysisRequest struct {
	URL     string            `json:"url" binding:"required"`
	Method  string            `json:"method" binding:"required"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func NewEndpointHandler(discoveryService services.DiscoveryServiceInterface, metadataService services.MetadataServiceInterface, logger logging.Logger) *EndpointHandler {
	return &EndpointHandler{
		discoveryService: discoveryService,
		metadataService:  metadataService,
		logger:          logger,
	}
}

func (h *EndpointHandler) AnalyzeEndpoint(c *gin.Context) {
	var req EndpointAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpoint := &models.Endpoint{
		URL:     req.URL,
		Method:  req.Method,
		Headers: req.Headers,
		Body:    req.Body,
	}

	analysis, err := h.discoveryService.AnalyzeEndpoint(c.Request.Context(), endpoint)
	if err != nil {
		h.logger.Error("Failed to analyze endpoint", "error", err, "url", req.URL)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze endpoint"})
		return
	}

	h.logger.Info("Endpoint analyzed", "url", req.URL, "method", req.Method)
	c.JSON(http.StatusOK, analysis)
}

func (h *EndpointHandler) GetEndpointMetadata(c *gin.Context) {
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint ID is required"})
		return
	}

	metadata, err := h.metadataService.GetEndpointMetadata(c.Request.Context(), endpointID)
	if err != nil {
		h.logger.Error("Failed to get endpoint metadata", "error", err, "endpoint_id", endpointID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get endpoint metadata"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

func (h *EndpointHandler) UpdateEndpointMetadata(c *gin.Context) {
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint ID is required"})
		return
	}

	var metadata models.Metadata
	if err := c.ShouldBindJSON(&metadata); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.metadataService.UpdateEndpointMetadata(c.Request.Context(), endpointID, &metadata)
	if err != nil {
		h.logger.Error("Failed to update endpoint metadata", "error", err, "endpoint_id", endpointID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update endpoint metadata"})
		return
	}

	h.logger.Info("Endpoint metadata updated", "endpoint_id", endpointID)
	c.JSON(http.StatusOK, gin.H{"message": "Endpoint metadata updated successfully"})
}
