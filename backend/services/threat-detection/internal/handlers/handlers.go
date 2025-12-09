package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"scopeapi.local/backend/services/threat-detection/internal/models"
	"scopeapi.local/backend/services/threat-detection/internal/services"
	"scopeapi.local/backend/shared/logging"
)

// ThreatHandler handles threat-related HTTP requests
type ThreatHandler struct {
	threatService services.ThreatDetectionServiceInterface
	logger        logging.Logger
}

// AnomalyHandler handles anomaly-related HTTP requests
type AnomalyHandler struct {
	anomalyService services.AnomalyDetectionServiceInterface
	logger         logging.Logger
}

// BehavioralHandler handles behavioral analysis HTTP requests
type BehavioralHandler struct {
	behavioralService services.BehavioralAnalysisServiceInterface
	logger           logging.Logger
}

// SignatureHandler handles signature detection HTTP requests
type SignatureHandler struct {
	signatureService services.SignatureDetectionServiceInterface
	logger          logging.Logger
}

// Constructor functions
func NewThreatHandler(threatService services.ThreatDetectionServiceInterface, logger logging.Logger) *ThreatHandler {
	return &ThreatHandler{
		threatService: threatService,
		logger:        logger,
	}
}

func NewAnomalyHandler(anomalyService services.AnomalyDetectionServiceInterface, logger logging.Logger) *AnomalyHandler {
	return &AnomalyHandler{
		anomalyService: anomalyService,
		logger:         logger,
	}
}

func NewBehavioralHandler(behavioralService services.BehavioralAnalysisServiceInterface, logger logging.Logger) *BehavioralHandler {
	return &BehavioralHandler{
		behavioralService: behavioralService,
		logger:           logger,
	}
}

func NewSignatureHandler(signatureService services.SignatureDetectionServiceInterface, logger logging.Logger) *SignatureHandler {
	return &SignatureHandler{
		signatureService: signatureService,
		logger:          logger,
	}
}

// =============================================================================
// THREAT HANDLER METHODS
// =============================================================================

// GetThreats retrieves a list of threats with filtering and pagination
func (h *ThreatHandler) GetThreats(c *gin.Context) {
	// Parse query parameters
	filter := &models.ThreatFilter{}
	
	// Parse attack type filter
	if attackType := c.Query("threat_type"); attackType != "" {
		filter.AttackType = attackType
	}
	
	// Parse status filter
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}
	
	// Parse date filter
	if sinceStr := c.Query("since"); sinceStr != "" {
		if since, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			filter.Since = since
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_DATE_FORMAT",
					"message": "Invalid date format. Use RFC3339 format (e.g., 2024-01-15T10:30:00Z)",
				},
			})
			return
		}
	}
	
	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			// Note: ThreatFilter doesn't have pagination fields, but we can add them
		}
	}
	
	// Get threats from service
	threats, err := h.threatService.GetThreats(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get threats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve threats",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": threats,
			"total": len(threats),
		},
		"message": "Threats retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetThreat retrieves a specific threat by ID
func (h *ThreatHandler) GetThreat(c *gin.Context) {
	threatID := c.Param("id")
	if threatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_THREAT_ID",
				"message": "Threat ID is required",
			},
		})
		return
	}
	
	// Get threat from service
	threat, err := h.threatService.GetThreat(c.Request.Context(), threatID)
	if err != nil {
		h.logger.Error("Failed to get threat", "threat_id", threatID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "THREAT_NOT_FOUND",
				"message": "Threat not found",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    threat,
		"message": "Threat retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// AnalyzeThreat analyzes traffic data for potential threats
func (h *ThreatHandler) AnalyzeThreat(c *gin.Context) {
	var request models.ThreatAnalysisRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.TrafficData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TRAFFIC_DATA",
				"message": "Traffic data is required",
			},
		})
		return
	}
	
	// Set request ID if not provided
	if request.RequestID == "" {
		request.RequestID = uuid.New().String()
	}
	
	// Set timestamp if not provided
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}
	
	// Analyze threat using service
	result, err := h.threatService.AnalyzeThreat(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to analyze threat", "request_id", request.RequestID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ANALYSIS_FAILED",
				"message": "Failed to analyze threat",
			},
		})
		return
	}
	
	// Return analysis result
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Threat analysis completed",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// UpdateThreatStatus updates the status of a threat
func (h *ThreatHandler) UpdateThreatStatus(c *gin.Context) {
	threatID := c.Param("id")
	if threatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_THREAT_ID",
				"message": "Threat ID is required",
			},
		})
		return
	}
	
	var updateRequest models.ThreatUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Update threat status using service
	err := h.threatService.UpdateThreatStatus(c.Request.Context(), threatID, &updateRequest)
	if err != nil {
		h.logger.Error("Failed to update threat status", "threat_id", threatID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update threat status",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Threat status updated successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// DeleteThreat deletes a threat
func (h *ThreatHandler) DeleteThreat(c *gin.Context) {
	threatID := c.Param("id")
	if threatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_THREAT_ID",
				"message": "Threat ID is required",
			},
		})
		return
	}
	
	// Delete threat using service
	err := h.threatService.DeleteThreat(c.Request.Context(), threatID)
	if err != nil {
		h.logger.Error("Failed to delete threat", "threat_id", threatID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "DELETE_FAILED",
				"message": "Failed to delete threat",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Threat deleted successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetThreatStatistics retrieves threat statistics
func (h *ThreatHandler) GetThreatStatistics(c *gin.Context) {
	// Parse time range parameter
	timeRangeStr := c.DefaultQuery("time_range", "24h")
	var timeRange time.Duration
	
	switch timeRangeStr {
	case "1h":
		timeRange = time.Hour
	case "24h":
		timeRange = 24 * time.Hour
	case "7d":
		timeRange = 7 * 24 * time.Hour
	case "30d":
		timeRange = 30 * 24 * time.Hour
	default:
		// Try to parse as duration string
		if parsed, err := time.ParseDuration(timeRangeStr); err == nil {
			timeRange = parsed
		} else {
			timeRange = 24 * time.Hour // Default to 24 hours
		}
	}
	
	// Get statistics from service
	stats, err := h.threatService.GetThreatStatistics(c.Request.Context(), timeRange)
	if err != nil {
		h.logger.Error("Failed to get threat statistics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "STATS_FAILED",
				"message": "Failed to retrieve threat statistics",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Threat statistics retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// =============================================================================
// ANOMALY HANDLER METHODS
// =============================================================================

// GetAnomalies retrieves a list of anomalies with filtering and pagination
func (h *AnomalyHandler) GetAnomalies(c *gin.Context) {
	// Parse query parameters
	filter := &models.AnomalyFilter{}
	
	// Parse type filter
	if types := c.QueryArray("type"); len(types) > 0 {
		filter.Type = types
	}
	
	// Parse severity filter
	if severities := c.QueryArray("severity"); len(severities) > 0 {
		filter.Severity = severities
	}
	
	// Parse status filter
	if statuses := c.QueryArray("status"); len(statuses) > 0 {
		filter.Status = statuses
	}
	
	// Parse API ID filter
	if apiID := c.Query("api_id"); apiID != "" {
		filter.APIID = apiID
	}
	
	// Parse endpoint ID filter
	if endpointID := c.Query("endpoint_id"); endpointID != "" {
		filter.EndpointID = endpointID
	}
	
	// Parse IP address filter
	if ipAddress := c.Query("ip_address"); ipAddress != "" {
		filter.IPAddress = ipAddress
	}
	
	// Parse user ID filter
	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = userID
	}
	
	// Parse score range filters
	if minScoreStr := c.Query("min_score"); minScoreStr != "" {
		if minScore, err := strconv.ParseFloat(minScoreStr, 64); err == nil {
			filter.MinScore = minScore
		}
	}
	
	if maxScoreStr := c.Query("max_score"); maxScoreStr != "" {
		if maxScore, err := strconv.ParseFloat(maxScoreStr, 64); err == nil {
			filter.MaxScore = maxScore
		}
	}
	
	// Parse false positive filter
	if falsePositiveStr := c.Query("false_positive"); falsePositiveStr != "" {
		if falsePositive, err := strconv.ParseBool(falsePositiveStr); err == nil {
			filter.FalsePositive = &falsePositive
		}
	}
	
	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			filter.DateFrom = dateFrom
		}
	}
	
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			filter.DateTo = dateTo
		}
	}
	
	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	
	// Set defaults
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	
	// Parse sorting
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}
	
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}
	
	// Get anomalies from service
	anomalies, err := h.anomalyService.GetAnomalies(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get anomalies", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve anomalies",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": anomalies,
			"pagination": gin.H{
				"page":  filter.Page,
				"limit": filter.Limit,
				"total": len(anomalies),
			},
		},
		"message": "Anomalies retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetAnomaly retrieves a specific anomaly by ID
func (h *AnomalyHandler) GetAnomaly(c *gin.Context) {
	anomalyID := c.Param("id")
	if anomalyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ANOMALY_ID",
				"message": "Anomaly ID is required",
			},
		})
		return
	}
	
	// Get anomaly from service
	anomaly, err := h.anomalyService.GetAnomaly(c.Request.Context(), anomalyID)
	if err != nil {
		h.logger.Error("Failed to get anomaly", "anomaly_id", anomalyID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ANOMALY_NOT_FOUND",
				"message": "Anomaly not found",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    anomaly,
		"message": "Anomaly retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// DetectAnomalies analyzes traffic data for anomalies
func (h *AnomalyHandler) DetectAnomalies(c *gin.Context) {
	var request models.AnomalyDetectionRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.TrafficData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TRAFFIC_DATA",
				"message": "Traffic data is required",
			},
		})
		return
	}
	
	// Set request ID if not provided
	if request.RequestID == "" {
		request.RequestID = uuid.New().String()
	}
	
	// Set timestamp if not provided
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}
	
	// Set default sensitivity if not provided
	if request.Sensitivity == 0 {
		request.Sensitivity = 0.7
	}
	
	// Detect anomalies using service
	result, err := h.anomalyService.DetectAnomalies(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to detect anomalies", "request_id", request.RequestID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "DETECTION_FAILED",
				"message": "Failed to detect anomalies",
			},
		})
		return
	}
	
	// Return detection result
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Anomaly detection completed",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ProvideFeedback provides feedback on an anomaly
func (h *AnomalyHandler) ProvideFeedback(c *gin.Context) {
	anomalyID := c.Param("id")
	if anomalyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ANOMALY_ID",
				"message": "Anomaly ID is required",
			},
		})
		return
	}
	
	var feedback models.AnomalyFeedback
	
	// Parse request body
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Set anomaly ID from URL parameter
	feedback.AnomalyID = anomalyID
	
	// Set timestamp if not provided
	if feedback.Timestamp.IsZero() {
		feedback.Timestamp = time.Now()
	}
	
	// Update anomaly feedback using service
	err := h.anomalyService.UpdateAnomalyFeedback(c.Request.Context(), &feedback)
	if err != nil {
		h.logger.Error("Failed to update anomaly feedback", "anomaly_id", anomalyID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FEEDBACK_UPDATE_FAILED",
				"message": "Failed to update anomaly feedback",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Anomaly feedback updated successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetAnomalyStatistics retrieves anomaly statistics
func (h *AnomalyHandler) GetAnomalyStatistics(c *gin.Context) {
	// Parse time range parameter
	timeRangeStr := c.DefaultQuery("time_range", "24h")
	var timeRange time.Duration
	
	switch timeRangeStr {
	case "1h":
		timeRange = time.Hour
	case "24h":
		timeRange = 24 * time.Hour
	case "7d":
		timeRange = 7 * 24 * time.Hour
	case "30d":
		timeRange = 30 * 24 * time.Hour
	default:
		// Try to parse as duration string
		if parsed, err := time.ParseDuration(timeRangeStr); err == nil {
			timeRange = parsed
		} else {
			timeRange = 24 * time.Hour // Default to 24 hours
		}
	}
	
	// Get statistics from service
	stats, err := h.anomalyService.GetAnomalyStatistics(c.Request.Context(), timeRange)
	if err != nil {
		h.logger.Error("Failed to get anomaly statistics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "STATS_FAILED",
				"message": "Failed to retrieve anomaly statistics",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Anomaly statistics retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// =============================================================================
// BEHAVIORAL HANDLER METHODS
// =============================================================================

// GetBehaviorPatterns retrieves behavior patterns with filtering
func (h *BehavioralHandler) GetBehaviorPatterns(c *gin.Context) {
	// Parse query parameters
	filter := &models.BehaviorPatternFilter{}
	
	// Parse type filter
	if patternType := c.Query("type"); patternType != "" {
		filter.Type = patternType
	}
	
	// Parse category filter
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}
	
	// Parse user ID filter
	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = userID
	}
	
	// Parse IP address filter
	if ipAddress := c.Query("ip_address"); ipAddress != "" {
		filter.IPAddress = ipAddress
	}
	
	// Parse API ID filter
	if apiID := c.Query("api_id"); apiID != "" {
		filter.APIID = apiID
	}
	
	// Parse risk score range
	if minRiskStr := c.Query("min_risk_score"); minRiskStr != "" {
		if minRisk, err := strconv.ParseFloat(minRiskStr, 64); err == nil {
			filter.MinRiskScore = minRisk
		}
	}
	
	if maxRiskStr := c.Query("max_risk_score"); maxRiskStr != "" {
		if maxRisk, err := strconv.ParseFloat(maxRiskStr, 64); err == nil {
			filter.MaxRiskScore = maxRisk
		}
	}
	
	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	
	// Set defaults
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	
	// Get behavior patterns from service
	patterns, err := h.behavioralService.GetBehaviorPatternsWithFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get behavior patterns", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve behavior patterns",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": patterns,
			"pagination": gin.H{
				"page":  filter.Page,
				"limit": filter.Limit,
				"total": len(patterns),
			},
		},
		"message": "Behavior patterns retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// AnalyzeBehavior analyzes behavior patterns
func (h *BehavioralHandler) AnalyzeBehavior(c *gin.Context) {
	var request models.BehaviorAnalysisRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.EntityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ENTITY_ID",
				"message": "Entity ID is required",
			},
		})
		return
	}
	
	if request.EntityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ENTITY_TYPE",
				"message": "Entity type is required",
			},
		})
		return
	}
	
	// Set timestamp if not provided
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}
	
	// Analyze behavior using service
	patterns, err := h.behavioralService.AnalyzeBehavior(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to analyze behavior", "entity_id", request.EntityID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ANALYSIS_FAILED",
				"message": "Failed to analyze behavior",
			},
		})
		return
	}
	
	// Return analysis result
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"patterns":     patterns,
			"entity_id":    request.EntityID,
			"entity_type":  request.EntityType,
			"analyzed_at":  time.Now(),
		},
		"message": "Behavior analysis completed",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetBaselines retrieves baseline profiles
func (h *BehavioralHandler) GetBaselines(c *gin.Context) {
	// Parse query parameters
	entityID := c.Query("entity_id")
	entityType := c.Query("entity_type")
	
	if entityID == "" || entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_PARAMETERS",
				"message": "Both entity_id and entity_type are required",
			},
		})
		return
	}
	
	// Get baseline from service
	baseline, err := h.behavioralService.GetBaselineProfile(c.Request.Context(), entityID, entityType)
	if err != nil {
		h.logger.Error("Failed to get baseline", "entity_id", entityID, "entity_type", entityType, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "BASELINE_NOT_FOUND",
				"message": "Baseline profile not found",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    baseline,
		"message": "Baseline profile retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// CreateBaseline creates a new baseline profile
func (h *BehavioralHandler) CreateBaseline(c *gin.Context) {
	var request models.BaselineCreationRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.EntityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ENTITY_ID",
				"message": "Entity ID is required",
			},
		})
		return
	}
	
	if request.EntityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ENTITY_TYPE",
				"message": "Entity type is required",
			},
		})
		return
	}
	
	// Create baseline using service
	err := h.behavioralService.CreateBaselineProfile(c.Request.Context(), request.EntityID, request.EntityType, request.TrainingData)
	if err != nil {
		h.logger.Error("Failed to create baseline", "entity_id", request.EntityID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CREATION_FAILED",
				"message": "Failed to create baseline profile",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Baseline profile created successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// =============================================================================
// SIGNATURE HANDLER METHODS
// =============================================================================

// GetSignatures retrieves a list of threat signatures with filtering
func (h *SignatureHandler) GetSignatures(c *gin.Context) {
	// Parse query parameters
	filter := &models.SignatureFilter{}
	
	// Parse signature type filter
	if sigType := c.Query("type"); sigType != "" {
		filter.Type = sigType
	}
	
	// Parse severity filter
	if severity := c.Query("severity"); severity != "" {
		filter.Severity = severity
	}
	
	// Parse category filter
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}
	
	// Get signatures from service
	signatures, err := h.signatureService.GetSignatures(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get signatures", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve signatures",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": signatures,
			"total": len(signatures),
		},
		"message": "Signatures retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetSignature retrieves a specific signature by ID
func (h *SignatureHandler) GetSignature(c *gin.Context) {
	signatureID := c.Param("id")
	if signatureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_SIGNATURE_ID",
				"message": "Signature ID is required",
			},
		})
		return
	}
	
	// Get signature from service
	signatures, err := h.signatureService.GetSignatures(c.Request.Context(), &models.SignatureFilter{})
	if err != nil {
		h.logger.Error("Failed to get signatures", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve signatures",
			},
		})
		return
	}
	
	// Find specific signature
	var signature *models.ThreatSignature
	for _, sig := range signatures {
		if sig.ID == signatureID {
			signature = &sig
			break
		}
	}
	
	if signature == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "SIGNATURE_NOT_FOUND",
				"message": "Signature not found",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    signature,
		"message": "Signature retrieved successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// DetectSignatures analyzes traffic data for signature matches
func (h *SignatureHandler) DetectSignatures(c *gin.Context) {
	var request models.SignatureDetectionRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.TrafficData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TRAFFIC_DATA",
				"message": "Traffic data is required",
			},
		})
		return
	}
	
	// Set request ID if not provided
	if request.RequestID == "" {
		request.RequestID = uuid.New().String()
	}
	
	// Set timestamp if not provided
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}
	
	// Detect signatures using service
	result, err := h.signatureService.DetectSignatures(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to detect signatures", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "DETECTION_FAILED",
				"message": "Failed to detect signatures",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Signature detection completed successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// TestSignature tests a signature against sample data
func (h *SignatureHandler) TestSignature(c *gin.Context) {
	var request models.SignatureTestRequest
	
	// Parse request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request format",
				"details": err.Error(),
			},
		})
		return
	}
	
	// Validate required fields
	if request.SignatureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_SIGNATURE_ID",
				"message": "Signature ID is required",
			},
		})
		return
	}
	
	if len(request.TestData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TEST_DATA",
				"message": "Test data is required",
			},
		})
		return
	}
	
	// Test signature using service
	result, err := h.signatureService.TestSignature(c.Request.Context(), request.SignatureID, request.TestData)
	if err != nil {
		h.logger.Error("Failed to test signature", "signature_id", request.SignatureID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TEST_FAILED",
				"message": "Failed to test signature",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Signature test completed successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ImportSignatureSet imports a set of signatures
func (h *SignatureHandler) ImportSignatureSet(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_FORM",
				"message": "Invalid form data",
			},
		})
		return
	}
	
	// Get signature set name
	signatureSet := c.PostForm("signature_set")
	if signatureSet == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_SIGNATURE_SET",
				"message": "Signature set name is required",
			},
		})
		return
	}
	
	// Get file data
	file, _, err := c.Request.FormFile("signature_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_FILE",
				"message": "Signature file is required",
			},
		})
		return
	}
	defer file.Close()
	
	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FILE_READ_ERROR",
				"message": "Failed to read signature file",
			},
		})
		return
	}
	
	// Import signatures using service
	err = h.signatureService.ImportSignatureSet(c.Request.Context(), fileData, signatureSet)
	if err != nil {
		h.logger.Error("Failed to import signature set", "signature_set", signatureSet, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "IMPORT_FAILED",
				"message": "Failed to import signature set",
			},
		})
		return
	}
	
	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signature set imported successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ExportSignatureSet exports a set of signatures
func (h *SignatureHandler) ExportSignatureSet(c *gin.Context) {
	signatureSet := c.Param("set")
	if signatureSet == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_SIGNATURE_SET",
				"message": "Signature set name is required",
			},
		})
		return
	}
	
	// Export signatures using service
	exportData, err := h.signatureService.ExportSignatureSet(c.Request.Context(), signatureSet)
	if err != nil {
		h.logger.Error("Failed to export signature set", "signature_set", signatureSet, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "EXPORT_FAILED",
				"message": "Failed to export signature set",
			},
		})
		return
	}
	
	// Return file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", signatureSet))
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", exportData)
}
