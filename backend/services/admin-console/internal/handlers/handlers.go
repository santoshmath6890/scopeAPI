package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct{}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// GetStats returns dashboard statistics
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats := gin.H{
		"total_users":      150,
		"active_sessions":  45,
		"api_endpoints":    23,
		"threats_blocked":  156,
		"system_health":    "good",
		"uptime":           "99.9%",
		"response_time":    "45ms",
		"data_processed":   "2.3TB",
		"alerts_today":     12,
		"incidents_resolved": 8,
	}
	c.JSON(http.StatusOK, stats)
}

// UserHandler handles user management requests
type UserHandler struct{}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users := []gin.H{
		{
			"id":        1,
			"name":      "Admin User",
			"email":     "admin@scopeapi.com",
			"role":      "admin",
			"status":    "active",
			"created_at": "2024-01-15T10:30:00Z",
			"last_login": "2024-01-20T14:45:00Z",
		},
		{
			"id":        2,
			"name":      "John Doe",
			"email":     "john@example.com",
			"role":      "user",
			"status":    "active",
			"created_at": "2024-01-16T09:15:00Z",
			"last_login": "2024-01-20T11:20:00Z",
		},
		{
			"id":        3,
			"name":      "Jane Smith",
			"email":     "jane@example.com",
			"role":      "user",
			"status":    "inactive",
			"created_at": "2024-01-17T16:45:00Z",
			"last_login": "2024-01-19T08:30:00Z",
		},
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock user creation
	newUser := gin.H{
		"id":        4,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"status":    "active",
		"created_at": "2024-01-20T15:00:00Z",
		"last_login": nil,
	}

	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock user update
	updatedUser := gin.H{
		"id":         userID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"updated_at": "2024-01-20T15:30:00Z",
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User deleted successfully",
		"id":        userID,
		"deleted_at": "2024-01-20T16:00:00Z",
	})
}

// SystemHandler handles system-related requests
type SystemHandler struct{}

// NewSystemHandler creates a new system handler
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetSystemInfo returns system information
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	systemInfo := gin.H{
		"version":     "1.0.0",
		"build_date":  "2024-01-20T10:00:00Z",
		"go_version":  "1.21.0",
		"uptime":      "72h 15m 30s",
		"memory_usage": "45%",
		"cpu_usage":   "23%",
		"disk_usage":  "67%",
		"services": []gin.H{
			{"name": "api-discovery", "status": "healthy", "port": 8081},
			{"name": "threat-detection", "status": "healthy", "port": 8082},
			{"name": "data-protection", "status": "healthy", "port": 8083},
			{"name": "gateway-integration", "status": "healthy", "port": 8084},
			{"name": "attack-blocking", "status": "healthy", "port": 8085},
		},
	}
	c.JSON(http.StatusOK, systemInfo)
}

// GetLogs returns system logs
func (h *SystemHandler) GetLogs(c *gin.Context) {
	logs := []gin.H{
		{
			"timestamp": "2024-01-20T15:45:00Z",
			"level":     "INFO",
			"message":   "Admin console service started successfully",
			"service":   "admin-console",
		},
		{
			"timestamp": "2024-01-20T15:44:30Z",
			"level":     "INFO",
			"message":   "Connected to database",
			"service":   "admin-console",
		},
		{
			"timestamp": "2024-01-20T15:44:00Z",
			"level":     "WARN",
			"message":   "High memory usage detected",
			"service":   "system",
		},
	}
	c.JSON(http.StatusOK, logs)
} 