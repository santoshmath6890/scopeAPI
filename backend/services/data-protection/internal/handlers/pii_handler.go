package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"data-protection/internal/models"
	"data-protection/internal/services"
	"shared/logging"
)

type PIIHandler struct {
	piiService services.PIIDetectionServiceInterface
	logger     logging.Logger
}

func NewPIIHandler(service services.PIIDetectionServiceInterface, logger logging.Logger) *PIIHandler {
	return &PIIHandler{
		piiService: service,
		logger:     logger,
	}
}

// DetectPII handles PII detection requests
func (h *PIIHandler) DetectPII(c *gin.Context) {
	var request models.PIIDetectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind PII detection request", "error", err)
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

	// Validate request
	if request.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_CONTENT",
				"message": "Content field is required",
			},
		})
		return
	}

	// Process PII detection
	result, err := h.piiService.DetectPII(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to detect PII", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PII_DETECTION_FAILED",
				"message": "Failed to detect PII",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PII detection completed successfully",
	})
}

// GetPIIPatterns retrieves PII patterns with filtering
func (h *PIIHandler) GetPIIPatterns(c *gin.Context) {
	filter := &models.PIIPatternFilter{}

	// Parse query parameters
	if piiType := c.Query("pii_type"); piiType != "" {
		filter.PIIType = piiType
	}
	if patternType := c.Query("type"); patternType != "" {
		filter.Type = patternType
	}
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			filter.Enabled = &enabled
		}
	}

	patterns, err := h.piiService.GetPIIPatterns(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get PII patterns", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PATTERNS_FETCH_FAILED",
				"message": "Failed to fetch PII patterns",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"patterns": patterns,
			"count":    len(patterns),
		},
		"message": "PII patterns retrieved successfully",
	})
}

// CreatePIIPattern creates a new PII pattern
func (h *PIIHandler) CreatePIIPattern(c *gin.Context) {
	var pattern models.PIIPattern
	if err := c.ShouldBindJSON(&pattern); err != nil {
		h.logger.Error("Failed to bind PII pattern", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_PATTERN",
				"message": "Invalid pattern format",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if pattern.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_NAME",
				"message": "Pattern name is required",
			},
		})
		return
	}

	if pattern.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_PATTERN",
				"message": "Pattern regex/expression is required",
			},
		})
		return
	}

	if pattern.PIIType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_PII_TYPE",
				"message": "PII type is required",
			},
		})
		return
	}

	// Create pattern
	err := h.piiService.CreatePIIPattern(c.Request.Context(), &pattern)
	if err != nil {
		h.logger.Error("Failed to create PII pattern", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PATTERN_CREATION_FAILED",
				"message": "Failed to create PII pattern",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    pattern,
		"message": "PII pattern created successfully",
	})
}

// GetPIIPattern retrieves a specific PII pattern
func (h *PIIHandler) GetPIIPattern(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Pattern ID is required",
			},
		})
		return
	}

	pattern, err := h.piiService.GetPIIPattern(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get PII pattern", "error", err, "pattern_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PATTERN_NOT_FOUND",
				"message": "PII pattern not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pattern,
		"message": "PII pattern retrieved successfully",
	})
}

// UpdatePIIPattern updates an existing PII pattern
func (h *PIIHandler) UpdatePIIPattern(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Pattern ID is required",
			},
		})
		return
	}

	var pattern models.PIIPattern
	if err := c.ShouldBindJSON(&pattern); err != nil {
		h.logger.Error("Failed to bind PII pattern update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_PATTERN",
				"message": "Invalid pattern format",
				"details": err.Error(),
			},
		})
		return
	}

	pattern.ID = id
	pattern.UpdatedAt = time.Now()

	err := h.piiService.UpdatePIIPattern(c.Request.Context(), &pattern)
	if err != nil {
		h.logger.Error("Failed to update PII pattern", "error", err, "pattern_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PATTERN_UPDATE_FAILED",
				"message": "Failed to update PII pattern",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pattern,
		"message": "PII pattern updated successfully",
	})
}

// DeletePIIPattern deletes a PII pattern
func (h *PIIHandler) DeletePIIPattern(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Pattern ID is required",
			},
		})
		return
	}

	err := h.piiService.DeletePIIPattern(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete PII pattern", "error", err, "pattern_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PATTERN_DELETION_FAILED",
				"message": "Failed to delete PII pattern",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PII pattern deleted successfully",
	})
}

// ScanForPII scans content for PII using all available patterns
func (h *PIIHandler) ScanForPII(c *gin.Context) {
	var request models.PIIScanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind PII scan request", "error", err)
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

	// Validate request
	if request.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_CONTENT",
				"message": "Content field is required",
			},
		})
		return
	}

	// Process PII scan
	result, err := h.piiService.ScanForPII(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to scan for PII", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PII_SCAN_FAILED",
				"message": "Failed to scan for PII",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PII scan completed successfully",
	})
}

// GetPIIReport generates a PII detection report
func (h *PIIHandler) GetPIIReport(c *gin.Context) {
	filter := &models.PIIReportFilter{}

	// Parse query parameters
	if sinceStr := c.Query("since"); sinceStr != "" {
		if since, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			filter.Since = &since
		}
	}
	if untilStr := c.Query("until"); untilStr != "" {
		if until, err := time.Parse(time.RFC3339, untilStr); err == nil {
			filter.Until = &until
		}
	}
	if piiType := c.Query("pii_type"); piiType != "" {
		filter.PIIType = piiType
	}
	if confidenceStr := c.Query("min_confidence"); confidenceStr != "" {
		if confidence, err := strconv.ParseFloat(confidenceStr, 64); err == nil {
			filter.MinConfidence = &confidence
		}
	}

	report, err := h.piiService.GetPIIReport(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to generate PII report", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_GENERATION_FAILED",
				"message": "Failed to generate PII report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
		"message": "PII report generated successfully",
	})
}
