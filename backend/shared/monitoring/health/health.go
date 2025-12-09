package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "scopeapi",
		Details: map[string]interface{}{
			"uptime": time.Since(time.Now()).String(),
		},
	}

	c.JSON(http.StatusOK, status)
}

// ReadinessCheckHandler handles readiness check requests
func ReadinessCheckHandler(c *gin.Context) {
	status := HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Service:   "scopeapi",
		Details: map[string]interface{}{
			"ready": true,
		},
	}

	c.JSON(http.StatusOK, status)
}
