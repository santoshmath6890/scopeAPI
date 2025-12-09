package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ServiceDiscovery handles inter-service communication
type ServiceDiscovery struct {
	client *http.Client
}

// NewServiceDiscovery creates a new service discovery instance
func NewServiceDiscovery() *ServiceDiscovery {
	return &ServiceDiscovery{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetServiceURL returns the URL for a given service
func (sd *ServiceDiscovery) GetServiceURL(serviceName string) string {
	// Try to get from config first
	configKey := fmt.Sprintf("services.%s.url", serviceName)
	if viper.IsSet(configKey) {
		return viper.GetString(configKey)
	}

	// Fallback to environment variable
	envKey := fmt.Sprintf("SERVICE_%s_URL", serviceName)
	if url := viper.GetString(envKey); url != "" {
		return url
	}

	// Default fallback based on service name
	switch serviceName {
	case "api_discovery":
		return "http://api-discovery:8080"
	case "threat_detection":
		return "http://threat-detection:8080"
	case "data_protection":
		return "http://data-protection:8080"
	case "gateway_integration":
		return "http://gateway-integration:8080"
	case "attack_blocking":
		return "http://attack-blocking:8080"
	default:
		return fmt.Sprintf("http://%s:8080", serviceName)
	}
}

// HealthCheck checks if a service is healthy
func (sd *ServiceDiscovery) HealthCheck(serviceName string) (bool, error) {
	url := sd.GetServiceURL(serviceName) + "/health"
	
	resp, err := sd.client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK, nil
}

// GetServiceStatus returns the status of all services
func (sd *ServiceDiscovery) GetServiceStatus() map[string]interface{} {
	services := []string{
		"api_discovery",
		"threat_detection", 
		"data_protection",
		"gateway_integration",
		"attack_blocking",
	}

	status := make(map[string]interface{})
	
	for _, service := range services {
		healthy, err := sd.HealthCheck(service)
		status[service] = map[string]interface{}{
			"healthy": healthy,
			"url":     sd.GetServiceURL(service),
			"error":   err,
		}
	}
	
	return status
}

// CallService makes an HTTP request to another service
func (sd *ServiceDiscovery) CallService(serviceName, endpoint, method string, body interface{}) (*http.Response, error) {
	url := sd.GetServiceURL(serviceName) + endpoint
	
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "admin-console/1.0")
	
	logrus.Debugf("Calling service %s at %s", serviceName, url)
	
	return sd.client.Do(req)
}

// ServiceStatusHandler returns a Gin handler for service status
func (sd *ServiceDiscovery) ServiceStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := sd.GetServiceStatus()
		c.JSON(http.StatusOK, gin.H{
			"services": status,
			"timestamp": time.Now().UTC(),
		})
	}
} 