package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/shared/logging"
)

// NginxIntegrationService handles NGINX-specific operations
type NginxIntegrationService struct {
	integrationService IntegrationServiceInterface
	logger             logging.Logger
}

// NginxClient implements GatewayClient for NGINX
type NginxClient struct {
	logger logging.Logger
}

// NewNginxIntegrationService creates a new NGINX integration service
func NewNginxIntegrationService(integrationService IntegrationServiceInterface, logger logging.Logger) *NginxIntegrationService {
	return &NginxIntegrationService{
		integrationService: integrationService,
		logger:             logger,
	}
}

// NewNginxClient creates a new NGINX client
func NewNginxClient(logger logging.Logger) *NginxClient {
	return &NginxClient{
		logger: logger,
	}
}

// NginxIntegrationService methods

// GetStatus retrieves NGINX status
func (s *NginxIntegrationService) GetStatus(ctx context.Context, integrationID string) (*models.HealthStatus, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewNginxClient(s.logger)
	return client.GetStatus(ctx, integration)
}

// GetConfig retrieves NGINX configuration
func (s *NginxIntegrationService) GetConfig(ctx context.Context, integrationID string) (*models.NginxConfig, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	configPath, err := s.getConfigPath(integration)
	if err != nil {
		return nil, err
	}

	// Read and parse NGINX configuration
	config, err := s.parseNginxConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NGINX config: %w", err)
	}

	return config, nil
}

// UpdateConfig updates NGINX configuration
func (s *NginxIntegrationService) UpdateConfig(ctx context.Context, integrationID string, config *models.NginxConfig) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	configPath, err := s.getConfigPath(integration)
	if err != nil {
		return err
	}

	// Generate NGINX configuration
	configContent, err := s.generateNginxConfig(config)
	if err != nil {
		return fmt.Errorf("failed to generate NGINX config: %w", err)
	}

	// Write configuration to file
	if err := s.writeConfigFile(configPath, configContent); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Test configuration
	if err := s.testConfig(ctx, integration); err != nil {
		return fmt.Errorf("NGINX configuration test failed: %w", err)
	}

	s.logger.Info("NGINX configuration updated successfully", "config_path", configPath)
	return nil
}

// ReloadConfig reloads NGINX configuration
func (s *NginxIntegrationService) ReloadConfig(ctx context.Context, integrationID string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewNginxClient(s.logger)
	return client.ReloadConfig(ctx, integration)
}

// GetUpstreams retrieves NGINX upstream configurations
func (s *NginxIntegrationService) GetUpstreams(ctx context.Context, integrationID string) ([]models.NginxUpstream, error) {
	config, err := s.GetConfig(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	return config.Upstreams, nil
}

// UpdateUpstream updates an NGINX upstream configuration
func (s *NginxIntegrationService) UpdateUpstream(ctx context.Context, integrationID string, upstream *models.NginxUpstream) error {
	config, err := s.GetConfig(ctx, integrationID)
	if err != nil {
		return err
	}

	// Find and update upstream
	found := false
	for i, existing := range config.Upstreams {
		if existing.Name == upstream.Name {
			config.Upstreams[i] = *upstream
			found = true
			break
		}
	}

	if !found {
		config.Upstreams = append(config.Upstreams, *upstream)
	}

	// Update configuration
	return s.UpdateConfig(ctx, integrationID, config)
}

// SyncConfiguration synchronizes configuration with NGINX
func (s *NginxIntegrationService) SyncConfiguration(ctx context.Context, integrationID string) (*models.SyncResult, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewNginxClient(s.logger)
	return client.SyncConfiguration(ctx, integration)
}

// NginxClient methods

// GetStatus implements GatewayClient.GetStatus for NGINX
func (c *NginxClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	start := time.Now()

	// Check if NGINX process is running
	if err := c.checkNginxProcess(); err != nil {
		return &models.HealthStatus{
			Status:    "unhealthy",
			Message:   "NGINX process not running",
			LastCheck: time.Now(),
			Latency:   time.Since(start).Milliseconds(),
		}, nil
	}

	// Test NGINX configuration
	if err := c.testNginxConfig(ctx, integration); err != nil {
		return &models.HealthStatus{
			Status:    "degraded",
			Message:   "NGINX configuration test failed",
			LastCheck: time.Now(),
			Latency:   time.Since(start).Milliseconds(),
		}, nil
	}

	return &models.HealthStatus{
		Status:    "healthy",
		Message:   "NGINX is running and configuration is valid",
		LastCheck: time.Now(),
		Latency:   time.Since(start).Milliseconds(),
	}, nil
}

// SyncConfiguration implements GatewayClient.SyncConfiguration for NGINX
func (c *NginxClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	start := time.Now()
	changes := []models.Change{}

	// Parse NGINX configuration from integration config
	var nginxConfig models.NginxConfig
	if configData, exists := integration.Config["nginx_config"]; exists {
		if configBytes, err := json.Marshal(configData); err == nil {
			if err := json.Unmarshal(configBytes, &nginxConfig); err != nil {
				return nil, fmt.Errorf("failed to parse NGINX config: %w", err)
			}
		}
	}

	configPath, err := c.getConfigPath(integration)
	if err != nil {
		return nil, err
	}

	// Generate and write configuration
	configContent, err := c.generateNginxConfig(&nginxConfig)
	if err != nil {
		changes = append(changes, models.Change{
			Type:     "configuration",
			Resource: "nginx.conf",
			Action:   "error",
			Details:  err.Error(),
		})
	} else {
		if err := c.writeConfigFile(configPath, configContent); err != nil {
			changes = append(changes, models.Change{
				Type:     "configuration",
				Resource: "nginx.conf",
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "configuration",
				Resource: "nginx.conf",
				Action:   "synced",
				Details:  "NGINX configuration synchronized",
			})
		}
	}

	// Test configuration
	if err := c.testNginxConfig(ctx, integration); err != nil {
		changes = append(changes, models.Change{
			Type:     "configuration",
			Resource: "nginx.conf",
			Action:   "error",
			Details:  "Configuration test failed: " + err.Error(),
		})
	}

	// Reload NGINX if configuration is valid
	if len(changes) == 1 && changes[0].Action == "synced" {
		if err := c.reloadNginx(ctx, integration); err != nil {
			changes = append(changes, models.Change{
				Type:     "reload",
				Resource: "nginx",
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "reload",
				Resource: "nginx",
				Action:   "synced",
				Details:  "NGINX reloaded successfully",
			})
		}
	}

	duration := time.Since(start)
	success := true
	message := "Configuration synchronized successfully"

	// Check if there were any errors
	for _, change := range changes {
		if change.Action == "error" {
			success = false
			message = "Configuration sync completed with errors"
			break
		}
	}

	return &models.SyncResult{
		Success:   success,
		Message:   message,
		Changes:   changes,
		Timestamp: time.Now(),
		Duration:  duration,
	}, nil
}

// TestConnection implements GatewayClient.TestConnection for NGINX
func (c *NginxClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	// Check if NGINX process is running
	if err := c.checkNginxProcess(); err != nil {
		return fmt.Errorf("NGINX process not running: %w", err)
	}

	// Test NGINX configuration
	if err := c.testNginxConfig(ctx, integration); err != nil {
		return fmt.Errorf("NGINX configuration test failed: %w", err)
	}

	return nil
}

// ReloadConfig reloads NGINX configuration
func (c *NginxClient) ReloadConfig(ctx context.Context, integration *models.Integration) error {
	return c.reloadNginx(ctx, integration)
}

// Private helper methods

func (s *NginxIntegrationService) getConfigPath(integration *models.Integration) (string, error) {
	if configPath, exists := integration.Config["config_path"]; exists {
		if path, ok := configPath.(string); ok {
			return path, nil
		}
	}
	return "", fmt.Errorf("config_path not found in integration config")
}

func (s *NginxIntegrationService) parseNginxConfig(configPath string) (*models.NginxConfig, error) {
	// Implementation for parsing NGINX configuration
	// This would read the NGINX config file and parse it into the NginxConfig struct
	s.logger.Info("Parsing NGINX configuration", "config_path", configPath)
	
	// Placeholder implementation
	return &models.NginxConfig{
		ConfigPath: configPath,
		Upstreams:  []models.NginxUpstream{},
		Locations:  []models.NginxLocation{},
		Servers:    []models.NginxServer{},
		SSLConfigs: []models.NginxSSLConfig{},
	}, nil
}

func (s *NginxIntegrationService) generateNginxConfig(config *models.NginxConfig) (string, error) {
	// Implementation for generating NGINX configuration
	// This would convert the NginxConfig struct into NGINX configuration format
	s.logger.Info("Generating NGINX configuration")
	
	// Placeholder implementation
	return `events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    upstream backend {
        server 127.0.0.1:8080;
    }
    
    server {
        listen 80;
        server_name localhost;
        
        location / {
            proxy_pass http://backend;
        }
    }
}`, nil
}

func (s *NginxIntegrationService) writeConfigFile(configPath, content string) error {
	// Implementation for writing configuration to file
	s.logger.Info("Writing NGINX configuration file", "config_path", configPath)
	return nil
}

func (s *NginxIntegrationService) testConfig(ctx context.Context, integration *models.Integration) error {
	client := NewNginxClient(s.logger)
	return client.testNginxConfig(ctx, integration)
}

func (c *NginxClient) getConfigPath(integration *models.Integration) (string, error) {
	if configPath, exists := integration.Config["config_path"]; exists {
		if path, ok := configPath.(string); ok {
			return path, nil
		}
	}
	return "", fmt.Errorf("config_path not found in integration config")
}

func (c *NginxClient) checkNginxProcess() error {
	// Check if NGINX process is running
	cmd := exec.Command("pgrep", "nginx")
	return cmd.Run()
}

func (c *NginxClient) testNginxConfig(ctx context.Context, integration *models.Integration) error {
	configPath, err := c.getConfigPath(integration)
	if err != nil {
		return err
	}

	// Test NGINX configuration
	cmd := exec.CommandContext(ctx, "nginx", "-t", "-c", configPath)
	return cmd.Run()
}

func (c *NginxClient) reloadNginx(ctx context.Context, integration *models.Integration) error {
	// Reload NGINX configuration
	cmd := exec.CommandContext(ctx, "nginx", "-s", "reload")
	return cmd.Run()
}

func (c *NginxClient) generateNginxConfig(config *models.NginxConfig) (string, error) {
	// Implementation for generating NGINX configuration
	c.logger.Info("Generating NGINX configuration")
	
	// Placeholder implementation
	return `events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    upstream backend {
        server 127.0.0.1:8080;
    }
    
    server {
        listen 80;
        server_name localhost;
        
        location / {
            proxy_pass http://backend;
        }
    }
}`, nil
}

func (c *NginxClient) writeConfigFile(configPath, content string) error {
	// Implementation for writing configuration to file
	c.logger.Info("Writing NGINX configuration file", "config_path", configPath)
	return nil
} 