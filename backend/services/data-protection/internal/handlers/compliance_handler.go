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

type ComplianceHandler struct {
	complianceService services.ComplianceServiceInterface
	logger            logging.Logger
}

func NewComplianceHandler(service services.ComplianceServiceInterface, logger logging.Logger) *ComplianceHandler {
	return &ComplianceHandler{
		complianceService: service,
		logger:            logger,
	}
}

// GetComplianceFrameworks retrieves compliance frameworks with filtering
func (h *ComplianceHandler) GetComplianceFrameworks(c *gin.Context) {
	filter := &models.ComplianceFrameworkFilter{}

	// Parse query parameters
	if frameworkType := c.Query("type"); frameworkType != "" {
		filter.Type = frameworkType
	}
	if region := c.Query("region"); region != "" {
		filter.Region = region
	}
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			filter.Enabled = &enabled
		}
	}

	frameworks, err := h.complianceService.GetComplianceFrameworks(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get compliance frameworks", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FRAMEWORKS_FETCH_FAILED",
				"message": "Failed to fetch compliance frameworks",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"frameworks": frameworks,
			"count":      len(frameworks),
		},
		"message": "Compliance frameworks retrieved successfully",
	})
}

// GetComplianceFramework retrieves a specific compliance framework
func (h *ComplianceHandler) GetComplianceFramework(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Framework ID is required",
			},
		})
		return
	}

	framework, err := h.complianceService.GetComplianceFramework(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get compliance framework", "error", err, "framework_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FRAMEWORK_NOT_FOUND",
				"message": "Compliance framework not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    framework,
		"message": "Compliance framework retrieved successfully",
	})
}

// CreateComplianceFramework creates a new compliance framework
func (h *ComplianceHandler) CreateComplianceFramework(c *gin.Context) {
	var framework models.ComplianceFramework
	if err := c.ShouldBindJSON(&framework); err != nil {
		h.logger.Error("Failed to bind compliance framework", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_FRAMEWORK",
				"message": "Invalid framework format",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if framework.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_NAME",
				"message": "Framework name is required",
			},
		})
		return
	}

	if framework.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TYPE",
				"message": "Framework type is required",
			},
		})
		return
	}

	if framework.Region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_REGION",
				"message": "Framework region is required",
			},
		})
		return
	}

	// Create framework
	err := h.complianceService.CreateComplianceFramework(c.Request.Context(), &framework)
	if err != nil {
		h.logger.Error("Failed to create compliance framework", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FRAMEWORK_CREATION_FAILED",
				"message": "Failed to create compliance framework",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    framework,
		"message": "Compliance framework created successfully",
	})
}

// UpdateComplianceFramework updates an existing compliance framework
func (h *ComplianceHandler) UpdateComplianceFramework(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Framework ID is required",
			},
		})
		return
	}

	var framework models.ComplianceFramework
	if err := c.ShouldBindJSON(&framework); err != nil {
		h.logger.Error("Failed to bind compliance framework update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_FRAMEWORK",
				"message": "Invalid framework format",
				"details": err.Error(),
			},
		})
		return
	}

	framework.ID = id
	framework.UpdatedAt = time.Now()

	err := h.complianceService.UpdateComplianceFramework(c.Request.Context(), &framework)
	if err != nil {
		h.logger.Error("Failed to update compliance framework", "error", err, "framework_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FRAMEWORK_UPDATE_FAILED",
				"message": "Failed to update compliance framework",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    framework,
		"message": "Compliance framework updated successfully",
	})
}

// DeleteComplianceFramework deletes a compliance framework
func (h *ComplianceHandler) DeleteComplianceFramework(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Framework ID is required",
			},
		})
		return
	}

	err := h.complianceService.DeleteComplianceFramework(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete compliance framework", "error", err, "framework_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FRAMEWORK_DELETION_FAILED",
				"message": "Failed to delete compliance framework",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Compliance framework deleted successfully",
	})
}

// GetComplianceReports retrieves compliance reports with filtering
func (h *ComplianceHandler) GetComplianceReports(c *gin.Context) {
	filter := &models.ComplianceReportFilter{}

	// Parse query parameters
	if frameworkID := c.Query("framework_id"); frameworkID != "" {
		filter.FrameworkID = frameworkID
	}
	if status := c.Query("status"); status != "" {
		filter.Status = status
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

	reports, err := h.complianceService.GetComplianceReports(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get compliance reports", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORTS_FETCH_FAILED",
				"message": "Failed to fetch compliance reports",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"reports": reports,
			"count":   len(reports),
		},
		"message": "Compliance reports retrieved successfully",
	})
}

// CreateComplianceReport creates a new compliance report
func (h *ComplianceHandler) CreateComplianceReport(c *gin.Context) {
	var report models.ComplianceReport
	if err := c.ShouldBindJSON(&report); err != nil {
		h.logger.Error("Failed to bind compliance report", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REPORT",
				"message": "Invalid report format",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if report.FrameworkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_FRAMEWORK_ID",
				"message": "Framework ID is required",
			},
		})
		return
	}

	if report.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_STATUS",
				"message": "Report status is required",
			},
		})
		return
	}

	// Create report
	err := h.complianceService.CreateComplianceReport(c.Request.Context(), &report)
	if err != nil {
		h.logger.Error("Failed to create compliance report", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_CREATION_FAILED",
				"message": "Failed to create compliance report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    report,
		"message": "Compliance report created successfully",
	})
}

// GetComplianceReport retrieves a specific compliance report
func (h *ComplianceHandler) GetComplianceReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Report ID is required",
			},
		})
		return
	}

	report, err := h.complianceService.GetComplianceReport(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get compliance report", "error", err, "report_id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_NOT_FOUND",
				"message": "Compliance report not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
		"message": "Compliance report retrieved successfully",
	})
}

// UpdateComplianceReport updates an existing compliance report
func (h *ComplianceHandler) UpdateComplianceReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Report ID is required",
			},
		})
		return
	}

	var report models.ComplianceReport
	if err := c.ShouldBindJSON(&report); err != nil {
		h.logger.Error("Failed to bind compliance report update", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REPORT",
				"message": "Invalid report format",
				"details": err.Error(),
			},
		})
		return
	}

	report.ID = id
	report.UpdatedAt = time.Now()

	err := h.complianceService.UpdateComplianceReport(c.Request.Context(), &report)
	if err != nil {
		h.logger.Error("Failed to update compliance report", "error", err, "report_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_UPDATE_FAILED",
				"message": "Failed to update compliance report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
		"message": "Compliance report updated successfully",
	})
}

// DeleteComplianceReport deletes a compliance report
func (h *ComplianceHandler) DeleteComplianceReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Report ID is required",
			},
		})
		return
	}

	err := h.complianceService.DeleteComplianceReport(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete compliance report", "error", err, "report_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "REPORT_DELETION_FAILED",
				"message": "Failed to delete compliance report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Compliance report deleted successfully",
	})
}

// GetAuditLog retrieves audit log entries with filtering
func (h *ComplianceHandler) GetAuditLog(c *gin.Context) {
	filter := &models.AuditLogFilter{}

	// Parse query parameters
	if action := c.Query("action"); action != "" {
		filter.Action = action
	}
	if resourceType := c.Query("resource_type"); resourceType != "" {
		filter.ResourceType = resourceType
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

	entries, err := h.complianceService.GetAuditLog(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get audit log", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUDIT_LOG_FETCH_FAILED",
				"message": "Failed to fetch audit log",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"entries": entries,
			"count":   len(entries),
		},
		"message": "Audit log retrieved successfully",
	})
}
