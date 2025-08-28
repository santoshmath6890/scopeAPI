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

type RiskHandler struct {
	riskService services.RiskScoringServiceInterface
	logger      logging.Logger
}

func NewRiskHandler(service services.RiskScoringServiceInterface, logger logging.Logger) *RiskHandler {
	return &RiskHandler{
		riskService: service,
		logger:      logger,
	}
}

// AssessRisk handles risk assessment requests
func (h *RiskHandler) AssessRisk(c *gin.Context) {
	var request models.RiskAssessmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind risk assessment request", "error", err)
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
	if request.DataSource == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_DATA_SOURCE",
				"message": "Data source is required",
			},
		})
		return
	}

	// Process risk assessment
	result, err := h.riskService.AssessRisk(c.Request.Context(), &request)
	if err != nil {
		h.logger.Error("Failed to assess risk", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RISK_ASSESSMENT_FAILED",
				"message": "Failed to assess risk",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "Risk assessment completed successfully",
	})
}

// GetRiskScores retrieves risk scores with filtering
func (h *RiskHandler) GetRiskScores(c *gin.Context) {
	filter := &models.RiskScoreFilter{}

	// Parse query parameters
	if dataSource := c.Query("data_source"); dataSource != "" {
		filter.DataSource = dataSource
	}
	if riskLevel := c.Query("risk_level"); riskLevel != "" {
		filter.RiskLevel = riskLevel
	}
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

	scores, err := h.riskService.GetRiskScores(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get risk scores", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RISK_SCORES_FETCH_FAILED",
				"message": "Failed to fetch risk scores",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"scores": scores,
			"count":  len(scores),
		},
		"message": "Risk scores retrieved successfully",
	})
}

// GetRiskScore retrieves a specific risk score
func (h *RiskHandler) GetRiskScore(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Risk score ID is required",
			},
		})
		return
	}

	score, err := h.riskService.GetRiskScore(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get risk score", "error", err, "score_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RISK_SCORE_NOT_FOUND",
				"message": "Risk score not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    score,
		"message": "Risk score retrieved successfully",
	})
}

// CreateMitigationPlan creates a new mitigation plan
func (h *RiskHandler) CreateMitigationPlan(c *gin.Context) {
	var plan models.MitigationPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		h.logger.Error("Failed to bind mitigation plan", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_PLAN",
				"message": "Invalid plan format",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if plan.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TITLE",
				"message": "Plan title is required",
			},
		})
		return
	}

	if plan.RiskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_RISK_ID",
				"message": "Risk ID is required",
			},
		})
		return
	}

	if len(plan.MitigationActions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ACTIONS",
				"message": "At least one mitigation action is required",
			},
		})
		return
	}

	// Create mitigation plan
	err := h.riskService.CreateMitigationPlan(c.Request.Context(), &plan)
	if err != nil {
		h.logger.Error("Failed to create mitigation plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PLAN_CREATION_FAILED",
				"message": "Failed to create mitigation plan",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    plan,
		"message": "Mitigation plan created successfully",
	})
}

// GetMitigationPlan retrieves a specific mitigation plan
func (h *RiskHandler) GetMitigationPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Plan ID is required",
			},
		})
		return
	}

	plan, err := h.riskService.GetMitigationPlan(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get mitigation plan", "error", err, "plan_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PLAN_NOT_FOUND",
				"message": "Mitigation plan not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
		"message": "Mitigation plan retrieved successfully",
	})
}

// UpdateMitigationPlan updates an existing mitigation plan
func (h *RiskHandler) UpdateMitigationPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Plan ID is required",
			},
		})
		return
	}

	var plan models.MitigationPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		h.logger.Error("Failed to bind mitigation plan update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_PLAN",
				"message": "Invalid plan format",
				"details": err.Error(),
			},
		})
		return
	}

	plan.ID = id
	plan.UpdatedAt = time.Now()

	err := h.riskService.UpdateMitigationPlan(c.Request.Context(), &plan)
	if err != nil {
		h.logger.Error("Failed to update mitigation plan", "error", err, "plan_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PLAN_UPDATE_FAILED",
				"message": "Failed to update mitigation plan",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
		"message": "Mitigation plan updated successfully",
	})
}

// DeleteMitigationPlan deletes a mitigation plan
func (h *RiskHandler) DeleteMitigationPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Plan ID is required",
			},
		})
		return
	}

	err := h.riskService.DeleteMitigationPlan(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete mitigation plan", "error", err, "plan_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "PLAN_DELETION_FAILED",
				"message": "Failed to delete mitigation plan",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mitigation plan deleted successfully",
	})
}
