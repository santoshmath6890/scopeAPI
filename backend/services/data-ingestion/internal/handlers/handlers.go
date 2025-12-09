package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/services/data-ingestion/internal/services"
	"scopeapi.local/backend/services/data-ingestion/internal/models"
	"github.com/google/uuid"
)

type IngestionHandler struct {
	ingestionService services.DataIngestionServiceInterface
	logger           Logger
}

type ParserHandler struct {
	parserService services.DataParserServiceInterface
	logger        Logger
}

type NormalizerHandler struct {
	normalizerService services.DataNormalizerServiceInterface
	logger           Logger
}

type QueueHandler struct {
	queueService services.QueueServiceInterface
	logger       Logger
}

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

func NewIngestionHandler(ingestionService services.DataIngestionServiceInterface, logger Logger) *IngestionHandler {
	return &IngestionHandler{
		ingestionService: ingestionService,
		logger:           logger,
	}
}

func NewParserHandler(parserService services.DataParserServiceInterface, logger Logger) *ParserHandler {
	return &ParserHandler{
		parserService: parserService,
		logger:        logger,
	}
}

func NewNormalizerHandler(normalizerService services.DataNormalizerServiceInterface, logger Logger) *NormalizerHandler {
	return &NormalizerHandler{
		normalizerService: normalizerService,
		logger:           logger,
	}
}

func NewQueueHandler(queueService services.QueueServiceInterface, logger Logger) *QueueHandler {
	return &QueueHandler{
		queueService: queueService,
		logger:       logger,
	}
}

// IngestTraffic handles incoming API traffic data
func (h *IngestionHandler) IngestTraffic(c *gin.Context) {
	var request models.IngestionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Set request ID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	// Process the traffic data
	response, err := h.ingestionService.IngestTraffic(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to process traffic", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process traffic"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// IngestBatch handles batch ingestion of API traffic data
func (h *IngestionHandler) IngestBatch(c *gin.Context) {
	var batch models.BatchTrafficData
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Set batch ID if not provided
	if batch.ID == "" {
		batch.ID = uuid.New().String()
	}

	// Process batch
	response, err := h.ingestionService.IngestBatch(c.Request.Context(), &batch)
	if err != nil {
		h.logger.Error("Failed to process batch", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process batch"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetIngestionStatus returns the status of a specific ingestion job
func (h *IngestionHandler) GetIngestionStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job ID"})
		return
	}

	status, err := h.ingestionService.GetIngestionStatus(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get ingestion status", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetIngestionStats returns ingestion statistics
func (h *IngestionHandler) GetIngestionStats(c *gin.Context) {
	// Create a time range for the last 24 hours
	timeRange := &models.TimeRange{
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now(),
	}

	stats, err := h.ingestionService.GetIngestionStats(c.Request.Context(), timeRange)
	if err != nil {
		h.logger.Error("Failed to get ingestion stats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ParseData handles data parsing requests
func (h *ParserHandler) ParseData(c *gin.Context) {
	var request struct {
		Data   []byte `json:"data" binding:"required"`
		Format string `json:"format"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	result, err := h.parserService.ParseData(c.Request.Context(), request.Data, request.Format, "")
	if err != nil {
		h.logger.Error("Failed to parse data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSupportedFormats returns list of supported data formats
func (h *ParserHandler) GetSupportedFormats(c *gin.Context) {
	formats, err := h.parserService.GetSupportedFormats(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get supported formats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get formats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"formats": formats})
}

// ValidateFormat validates if a format is supported
func (h *ParserHandler) ValidateFormat(c *gin.Context) {
	var request struct {
		Data   []byte `json:"data" binding:"required"`
		Format string `json:"format" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	result, err := h.parserService.ValidateFormat(c.Request.Context(), request.Data, request.Format)
	if err != nil {
		h.logger.Error("Failed to validate format", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate format"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// NormalizeData handles data normalization requests
func (h *NormalizerHandler) NormalizeData(c *gin.Context) {
	var request struct {
		ParsedData *models.ParsedData `json:"parsed_data" binding:"required"`
		Schema     string             `json:"schema"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	result, err := h.normalizerService.NormalizeData(c.Request.Context(), request.ParsedData, request.Schema)
	if err != nil {
		h.logger.Error("Failed to normalize data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to normalize data"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSchemas returns available normalization schemas
func (h *NormalizerHandler) GetSchemas(c *gin.Context) {
	schemas, err := h.normalizerService.GetSchemas(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get schemas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schemas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"schemas": schemas})
}

// CreateSchema creates a new normalization schema
func (h *NormalizerHandler) CreateSchema(c *gin.Context) {
	var schema models.SchemaInfo
	if err := c.ShouldBindJSON(&schema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.normalizerService.CreateSchema(c.Request.Context(), schema)
	if err != nil {
		h.logger.Error("Failed to create schema", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schema"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schema created successfully"})
}

// GetQueueStatus returns the current queue status
func (h *QueueHandler) GetQueueStatus(c *gin.Context) {
	status, err := h.queueService.GetQueueStatus(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get queue status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue status"})
		return
	}
	c.JSON(http.StatusOK, status)
}

// FlushQueue flushes the current queue
func (h *QueueHandler) FlushQueue(c *gin.Context) {
	err := h.queueService.Flush(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to flush queue", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flush queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Queue flushed successfully"})
}

// GetTopics returns available Kafka topics
func (h *QueueHandler) GetTopics(c *gin.Context) {
	topics, err := h.queueService.GetTopics(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get topics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get topics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"topics": topics})
} 