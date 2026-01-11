package handlers

import (
	"net/http"
	"strconv"
	"time"

	"data-protection/internal/models"
	"data-protection/internal/services"

	"github.com/gin-gonic/gin"
	"scopeapi.local/backend/shared/logging"
)

type ClassificationHandler struct {
	classificationService services.DataClassificationServiceInterface
	logger                logging.Logger
}

func NewClassificationHandler(service services.DataClassificationServiceInterface, logger logging.Logger) *ClassificationHandler {
	return &ClassificationHandler{
		classificationService: service,
		logger:                logger,
	}
}

// ClassifyData handles data classification requests
func (h *ClassificationHandler) ClassifyData(c *gin.Context) {
	var request models.DataClassificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind classification request", "error", err)
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

	// Process classification
	result, err := h.classificationService.ClassifyData(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to classify data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CLASSIFICATION_FAILED",
				"message": "Failed to classify data",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Data classified successfully",
	})
}

// GetClassificationRules retrieves classification rules with filtering
func (h *ClassificationHandler) GetClassificationRules(c *gin.Context) {
	filter := &models.ClassificationRuleFilter{}

	// Parse query parameters
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}
	if method := c.Query("method"); method != "" {
		filter.Method = method
	}
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			filter.Enabled = &enabled
		}
	}

	rules, err := h.classificationService.GetClassificationRules(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get classification rules", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULES_FETCH_FAILED",
				"message": "Failed to fetch classification rules",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"rules": rules,
			"count": len(rules),
		},
		"message": "Classification rules retrieved successfully",
	})
}

// CreateClassificationRule creates a new classification rule
func (h *ClassificationHandler) CreateClassificationRule(c *gin.Context) {
	var rule models.ClassificationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		h.logger.Error("Failed to bind classification rule", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RULE",
				"message": "Invalid rule format",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if rule.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_NAME",
				"message": "Rule name is required",
			},
		})
		return
	}

	if rule.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_CATEGORY",
				"message": "Rule category is required",
			},
		})
		return
	}

	// Create rule
	err := h.classificationService.CreateClassificationRule(c.Request.Context(), &rule)
	if err != nil {
		h.logger.Error("Failed to create classification rule", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_CREATION_FAILED",
				"message": "Failed to create classification rule",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    rule,
		"message": "Classification rule created successfully",
	})
}

// GetClassificationRule retrieves a specific classification rule
func (h *ClassificationHandler) GetClassificationRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Rule ID is required",
			},
		})
		return
	}

	rule, err := h.classificationService.GetClassificationRule(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get classification rule", "error", err, "rule_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_NOT_FOUND",
				"message": "Classification rule not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rule,
		"message": "Classification rule retrieved successfully",
	})
}

// UpdateClassificationRule updates an existing classification rule
func (h *ClassificationHandler) UpdateClassificationRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Rule ID is required",
			},
		})
		return
	}

	var rule models.ClassificationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		h.logger.Error("Failed to bind classification rule update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RULE",
				"message": "Invalid rule format",
				"details": err.Error(),
			},
		})
		return
	}

	rule.ID = id
	rule.UpdatedAt = time.Now()

	err := h.classificationService.UpdateClassificationRule(c.Request.Context(), id, &rule)
	if err != nil {
		h.logger.Error("Failed to update classification rule", "error", err, "rule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_UPDATE_FAILED",
				"message": "Failed to update classification rule",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rule,
		"message": "Classification rule updated successfully",
	})
}

// DeleteClassificationRule deletes a classification rule
func (h *ClassificationHandler) DeleteClassificationRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Rule ID is required",
			},
		})
		return
	}

	err := h.classificationService.DeleteClassificationRule(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete classification rule", "error", err, "rule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_DELETION_FAILED",
				"message": "Failed to delete classification rule",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Classification rule deleted successfully",
	})
}

// EnableClassificationRule enables a classification rule
func (h *ClassificationHandler) EnableClassificationRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Rule ID is required",
			},
		})
		return
	}

	err := h.classificationService.EnableClassificationRule(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to enable classification rule", "error", err, "rule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_ENABLE_FAILED",
				"message": "Failed to enable classification rule",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Classification rule enabled successfully",
	})
}

// DisableClassificationRule disables a classification rule
func (h *ClassificationHandler) DisableClassificationRule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Rule ID is required",
			},
		})
		return
	}

	err := h.classificationService.DisableClassificationRule(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to disable classification rule", "error", err, "rule_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RULE_DISABLE_FAILED",
				"message": "Failed to disable classification rule",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Classification rule disabled successfully",
	})
}

// GetClassificationReport generates a classification report
func (h *ClassificationHandler) GetClassificationReport(c *gin.Context) {
	filter := &models.ClassificationReportFilter{}

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
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}

	report, err := h.classificationService.GetDataClassificationReport(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to generate classification report", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_GENERATION_FAILED",
				"message": "Failed to generate classification report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
		"message": "Classification report generated successfully",
	})
}
