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

// KongIntegrationService handles Kong-specific operations
type KongIntegrationService struct {
	integrationService IntegrationServiceInterface
	logger             logging.Logger
	httpClient         *http.Client
}

// KongClient implements GatewayClient for Kong
type KongClient struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewKongIntegrationService creates a new Kong integration service
func NewKongIntegrationService(integrationService IntegrationServiceInterface, logger logging.Logger) *KongIntegrationService {
	return &KongIntegrationService{
		integrationService: integrationService,
		logger:             logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewKongClient creates a new Kong client
func NewKongClient(logger logging.Logger) *KongClient {
	return &KongClient{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Kong API Response structures
type KongStatusResponse struct {
	Server string `json:"server"`
	Database struct {
		Reachable bool `json:"reachable"`
	} `json:"database"`
}

type KongServiceResponse struct {
	Data []models.KongService `json:"data"`
	Total int `json:"total"`
}

type KongRouteResponse struct {
	Data []models.KongRoute `json:"data"`
	Total int `json:"total"`
}

type KongPluginResponse struct {
	Data []models.KongPlugin `json:"data"`
	Total int `json:"total"`
}

// KongService methods

// GetStatus retrieves Kong status
func (s *KongIntegrationService) GetStatus(ctx context.Context, integrationID string) (*models.HealthStatus, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewKongClient(s.logger)
	return client.GetStatus(ctx, integration)
}

// GetServices retrieves Kong services
func (s *KongIntegrationService) GetServices(ctx context.Context, integrationID string) ([]models.KongService, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/services", adminURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get services, status: %d", resp.StatusCode)
	}

	var kongResp KongServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&kongResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return kongResp.Data, nil
}

// GetRoutes retrieves Kong routes
func (s *KongIntegrationService) GetRoutes(ctx context.Context, integrationID string) ([]models.KongRoute, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/routes", adminURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get routes, status: %d", resp.StatusCode)
	}

	var kongResp KongRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&kongResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return kongResp.Data, nil
}

// GetPlugins retrieves Kong plugins
func (s *KongIntegrationService) GetPlugins(ctx context.Context, integrationID string) ([]models.KongPlugin, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/plugins", adminURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get plugins, status: %d", resp.StatusCode)
	}

	var kongResp KongPluginResponse
	if err := json.NewDecoder(resp.Body).Decode(&kongResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return kongResp.Data, nil
}

// CreatePlugin creates a new Kong plugin
func (s *KongIntegrationService) CreatePlugin(ctx context.Context, integrationID string, plugin *models.KongPlugin) (*models.KongPlugin, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/plugins", adminURL)
	pluginData, err := json.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plugin: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(pluginData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create plugin, status: %d", resp.StatusCode)
	}

	var createdPlugin models.KongPlugin
	if err := json.NewDecoder(resp.Body).Decode(&createdPlugin); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdPlugin, nil
}

// UpdatePlugin updates an existing Kong plugin
func (s *KongIntegrationService) UpdatePlugin(ctx context.Context, integrationID, pluginID string, plugin *models.KongPlugin) (*models.KongPlugin, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/plugins/%s", adminURL, pluginID)
	pluginData, err := json.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plugin: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(pluginData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update plugin, status: %d", resp.StatusCode)
	}

	var updatedPlugin models.KongPlugin
	if err := json.NewDecoder(resp.Body).Decode(&updatedPlugin); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedPlugin, nil
}

// DeletePlugin deletes a Kong plugin
func (s *KongIntegrationService) DeletePlugin(ctx context.Context, integrationID, pluginID string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	adminURL, err := s.getAdminURL(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/plugins/%s", adminURL, pluginID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete plugin, status: %d", resp.StatusCode)
	}

	return nil
}

// SyncConfiguration synchronizes configuration with Kong
func (s *KongIntegrationService) SyncConfiguration(ctx context.Context, integrationID string) (*models.SyncResult, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewKongClient(s.logger)
	return client.SyncConfiguration(ctx, integration)
}

// KongClient methods

// GetStatus implements GatewayClient.GetStatus for Kong
func (c *KongClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	adminURL, err := c.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	url := fmt.Sprintf("%s/status", adminURL)
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

	var status KongStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return &models.HealthStatus{
			Status:    "unhealthy",
			Message:   "Failed to decode status response",
			LastCheck: time.Now(),
			Latency:   latency,
		}, nil
	}

	healthStatus := "healthy"
	message := "Kong is running"
	if !status.Database.Reachable {
		healthStatus = "degraded"
		message = "Database is not reachable"
	}

	return &models.HealthStatus{
		Status:    healthStatus,
		Message:   message,
		LastCheck: time.Now(),
		Latency:   latency,
	}, nil
}

// SyncConfiguration implements GatewayClient.SyncConfiguration for Kong
func (c *KongClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	start := time.Now()
	changes := []models.Change{}

	// Parse Kong configuration from integration config
	var kongConfig models.KongConfig
	if configData, exists := integration.Config["kong_config"]; exists {
		if configBytes, err := json.Marshal(configData); err == nil {
			if err := json.Unmarshal(configBytes, &kongConfig); err != nil {
				return nil, fmt.Errorf("failed to parse Kong config: %w", err)
			}
		}
	}

	adminURL, err := c.getAdminURL(integration)
	if err != nil {
		return nil, err
	}

	// Sync services
	for _, service := range kongConfig.Services {
		if err := c.syncService(ctx, adminURL, integration, &service); err != nil {
			c.logger.Error("Failed to sync service", "service_name", service.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "service",
				Resource: service.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "service",
				Resource: service.Name,
				Action:   "synced",
				Details:  "Service configuration synchronized",
			})
		}
	}

	// Sync routes
	for _, route := range kongConfig.Routes {
		if err := c.syncRoute(ctx, adminURL, integration, &route); err != nil {
			c.logger.Error("Failed to sync route", "route_name", route.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "route",
				Resource: route.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "route",
				Resource: route.Name,
				Action:   "synced",
				Details:  "Route configuration synchronized",
			})
		}
	}

	// Sync plugins
	for _, plugin := range kongConfig.Plugins {
		if err := c.syncPlugin(ctx, adminURL, integration, &plugin); err != nil {
			c.logger.Error("Failed to sync plugin", "plugin_name", plugin.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "plugin",
				Resource: plugin.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "plugin",
				Resource: plugin.Name,
				Action:   "synced",
				Details:  "Plugin configuration synchronized",
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

// TestConnection implements GatewayClient.TestConnection for Kong
func (c *KongClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	adminURL, err := c.getAdminURL(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/status", adminURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req, integration)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Kong returned status %d", resp.StatusCode)
	}

	return nil
}

// Private helper methods

func (s *KongIntegrationService) getAdminURL(integration *models.Integration) (string, error) {
	if adminURL, exists := integration.Config["admin_url"]; exists {
		if url, ok := adminURL.(string); ok {
			return url, nil
		}
	}
	return "", fmt.Errorf("admin_url not found in integration config")
}

func (s *KongIntegrationService) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *KongClient) getAdminURL(integration *models.Integration) (string, error) {
	if adminURL, exists := integration.Config["admin_url"]; exists {
		if url, ok := adminURL.(string); ok {
			return url, nil
		}
	}
	return "", fmt.Errorf("admin_url not found in integration config")
}

func (c *KongClient) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *KongClient) syncService(ctx context.Context, adminURL string, integration *models.Integration, service *models.KongService) error {
	// Implementation for syncing Kong service
	// This would check if the service exists and create/update it accordingly
	c.logger.Info("Syncing Kong service", "service_name", service.Name)
	return nil
}

func (c *KongClient) syncRoute(ctx context.Context, adminURL string, integration *models.Integration, route *models.KongRoute) error {
	// Implementation for syncing Kong route
	// This would check if the route exists and create/update it accordingly
	c.logger.Info("Syncing Kong route", "route_name", route.Name)
	return nil
}

func (c *KongClient) syncPlugin(ctx context.Context, adminURL string, integration *models.Integration, plugin *models.KongPlugin) error {
	// Implementation for syncing Kong plugin
	// This would check if the plugin exists and create/update it accordingly
	c.logger.Info("Syncing Kong plugin", "plugin_name", plugin.Name)
	return nil
} 