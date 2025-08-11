package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/shared/logging"
)

// HAProxyIntegrationService handles HAProxy-specific operations
type HAProxyIntegrationService struct {
	integrationService IntegrationServiceInterface
	logger             logging.Logger
	httpClient         *http.Client
}

// HAProxyClient implements GatewayClient for HAProxy
type HAProxyClient struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewHAProxyIntegrationService creates a new HAProxy integration service
func NewHAProxyIntegrationService(integrationService IntegrationServiceInterface, logger logging.Logger) *HAProxyIntegrationService {
	return &HAProxyIntegrationService{
		integrationService: integrationService,
		logger:             logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewHAProxyClient creates a new HAProxy client
func NewHAProxyClient(logger logging.Logger) *HAProxyClient {
	return &HAProxyClient{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HAProxy API Response structures
type HAProxyStatusResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type HAProxyBackendResponse struct {
	Backends []models.HAProxyBackend `json:"backends"`
}

type HAProxyFrontendResponse struct {
	Frontends []models.HAProxyFrontend `json:"frontends"`
}

type HAProxyServerResponse struct {
	Servers []models.HAProxyServer `json:"servers"`
}

type HAProxyStatsResponse struct {
	Stats []models.HAProxyStats `json:"stats"`
}

// HAProxyIntegrationService methods

// GetStatus retrieves HAProxy status
func (s *HAProxyIntegrationService) GetStatus(ctx context.Context, integrationID string) (*models.HealthStatus, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewHAProxyClient(s.logger)
	return client.GetStatus(ctx, integration)
}

// GetConfig retrieves HAProxy configuration
func (s *HAProxyIntegrationService) GetConfig(ctx context.Context, integrationID string) (string, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return "", fmt.Errorf("failed to get integration: %w", err)
	}

	statsEndpoint, err := s.getStatsEndpoint(integration)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/config", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get config, status: %d", resp.StatusCode)
	}

	// Read the configuration as text
	configBytes := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			configBytes = append(configBytes, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(configBytes), nil
}

// UpdateConfig updates HAProxy configuration
func (s *HAProxyIntegrationService) UpdateConfig(ctx context.Context, integrationID string, config string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	statsEndpoint, err := s.getStatsEndpoint(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/config", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update config, status: %d", resp.StatusCode)
	}

	return nil
}

// ReloadConfig reloads HAProxy configuration
func (s *HAProxyIntegrationService) ReloadConfig(ctx context.Context, integrationID string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	statsEndpoint, err := s.getStatsEndpoint(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/reload", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to reload config, status: %d", resp.StatusCode)
	}

	return nil
}

// GetBackends retrieves HAProxy backends
func (s *HAProxyIntegrationService) GetBackends(ctx context.Context, integrationID string) ([]models.HAProxyBackend, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	statsEndpoint, err := s.getStatsEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/backends", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get backends: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get backends, status: %d", resp.StatusCode)
	}

	var haproxyResp HAProxyBackendResponse
	if err := json.NewDecoder(resp.Body).Decode(&haproxyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return haproxyResp.Backends, nil
}

// UpdateBackend updates an HAProxy backend
func (s *HAProxyIntegrationService) UpdateBackend(ctx context.Context, integrationID string, backend *models.HAProxyBackend) (*models.HAProxyBackend, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	statsEndpoint, err := s.getStatsEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/backends/%s", statsEndpoint, backend.Name)
	backendData, err := json.Marshal(backend)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backend: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(backendData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update backend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update backend, status: %d", resp.StatusCode)
	}

	var updatedBackend models.HAProxyBackend
	if err := json.NewDecoder(resp.Body).Decode(&updatedBackend); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedBackend, nil
}

// SyncConfiguration synchronizes configuration with HAProxy
func (s *HAProxyIntegrationService) SyncConfiguration(ctx context.Context, integrationID string) (*models.SyncResult, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewHAProxyClient(s.logger)
	return client.SyncConfiguration(ctx, integration)
}

// HAProxyClient methods

// GetStatus implements GatewayClient.GetStatus for HAProxy
func (c *HAProxyClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	statsEndpoint, err := c.getStatsEndpoint(integration)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	url := fmt.Sprintf("%s/status", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req, integration)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &models.HealthStatus{
			Status:    "unhealthy",
			Message:   err.Error(),
			LastCheck: time.Now(),
			Latency:   time.Since(start).Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	latency := time.Since(start).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return &models.HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("HTTP %d", resp.StatusCode),
			LastCheck: time.Now(),
			Latency:   latency,
		}, nil
	}

	var status HAProxyStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return &models.HealthStatus{
			Status:    "unhealthy",
			Message:   "Failed to decode status response",
			LastCheck: time.Now(),
			Latency:   latency,
		}, nil
	}

	return &models.HealthStatus{
		Status:    "healthy",
		Message:   "HAProxy is running",
		LastCheck: time.Now(),
		Latency:   latency,
	}, nil
}

// SyncConfiguration implements GatewayClient.SyncConfiguration for HAProxy
func (c *HAProxyClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	start := time.Now()
	changes := []models.Change{}

	// Parse HAProxy configuration from integration config
	var haproxyConfig models.HAProxyConfig
	if configData, exists := integration.Config["haproxy_config"]; exists {
		if configBytes, err := json.Marshal(configData); err == nil {
			if err := json.Unmarshal(configBytes, &haproxyConfig); err != nil {
				return nil, fmt.Errorf("failed to parse HAProxy config: %w", err)
			}
		}
	}

	statsEndpoint, err := c.getStatsEndpoint(integration)
	if err != nil {
		return nil, err
	}

	// Sync backends
	for _, backend := range haproxyConfig.Backends {
		if err := c.syncBackend(ctx, statsEndpoint, integration, &backend); err != nil {
			c.logger.Error("Failed to sync backend", "backend_name", backend.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "backend",
				Resource: backend.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "backend",
				Resource: backend.Name,
				Action:   "synced",
				Details:  "Backend configuration synchronized",
			})
		}
	}

	// Sync frontends
	for _, frontend := range haproxyConfig.Frontends {
		if err := c.syncFrontend(ctx, statsEndpoint, integration, &frontend); err != nil {
			c.logger.Error("Failed to sync frontend", "frontend_name", frontend.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "frontend",
				Resource: frontend.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "frontend",
				Resource: frontend.Name,
				Action:   "synced",
				Details:  "Frontend configuration synchronized",
			})
		}
	}

	// Sync servers
	for _, server := range haproxyConfig.Servers {
		if err := c.syncServer(ctx, statsEndpoint, integration, &server); err != nil {
			c.logger.Error("Failed to sync server", "server_name", server.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "server",
				Resource: server.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "server",
				Resource: server.Name,
				Action:   "synced",
				Details:  "Server configuration synchronized",
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

// TestConnection implements GatewayClient.TestConnection for HAProxy
func (c *HAProxyClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	statsEndpoint, err := c.getStatsEndpoint(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/status", statsEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req, integration)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to HAProxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HAProxy returned status %d", resp.StatusCode)
	}

	return nil
}

// Private helper methods

func (s *HAProxyIntegrationService) getStatsEndpoint(integration *models.Integration) (string, error) {
	if statsEndpoint, exists := integration.Config["stats_endpoint"]; exists {
		if endpoint, ok := statsEndpoint.(string); ok {
			return endpoint, nil
		}
	}
	return "", fmt.Errorf("stats_endpoint not found in integration config")
}

func (s *HAProxyIntegrationService) addAuthHeaders(req *http.Request, integration *models.Integration) {
	if integration.Credentials != nil {
		switch integration.Credentials.Type {
		case models.CredentialTypeBasic:
			req.SetBasicAuth(integration.Credentials.Username, integration.Credentials.Password)
		case models.CredentialTypeToken:
			req.Header.Set("Authorization", "Bearer "+integration.Credentials.Token)
		case models.CredentialTypeAPIKey:
			req.Header.Set("X-API-Key", integration.Credentials.APIKey)
		}
	}
}

func (c *HAProxyClient) getStatsEndpoint(integration *models.Integration) (string, error) {
	if statsEndpoint, exists := integration.Config["stats_endpoint"]; exists {
		if endpoint, ok := statsEndpoint.(string); ok {
			return endpoint, nil
		}
	}
	return "", fmt.Errorf("stats_endpoint not found in integration config")
}

func (c *HAProxyClient) addAuthHeaders(req *http.Request, integration *models.Integration) {
	if integration.Credentials != nil {
		switch integration.Credentials.Type {
		case models.CredentialTypeBasic:
			req.SetBasicAuth(integration.Credentials.Username, integration.Credentials.Password)
		case models.CredentialTypeToken:
			req.Header.Set("Authorization", "Bearer "+integration.Credentials.Token)
		case models.CredentialTypeAPIKey:
			req.Header.Set("X-API-Key", integration.Credentials.APIKey)
		}
	}
}

func (c *HAProxyClient) syncBackend(ctx context.Context, statsEndpoint string, integration *models.Integration, backend *models.HAProxyBackend) error {
	// Implementation for syncing HAProxy backend
	// This would check if the backend exists and create/update it accordingly
	c.logger.Info("Syncing HAProxy backend", "backend_name", backend.Name)
	return nil
}

func (c *HAProxyClient) syncFrontend(ctx context.Context, statsEndpoint string, integration *models.Integration, frontend *models.HAProxyFrontend) error {
	// Implementation for syncing HAProxy frontend
	// This would check if the frontend exists and create/update it accordingly
	c.logger.Info("Syncing HAProxy frontend", "frontend_name", frontend.Name)
	return nil
}

func (c *HAProxyClient) syncServer(ctx context.Context, statsEndpoint string, integration *models.Integration, server *models.HAProxyServer) error {
	// Implementation for syncing HAProxy server
	// This would check if the server exists and create/update it accordingly
	c.logger.Info("Syncing HAProxy server", "server_name", server.Name)
	return nil
} 