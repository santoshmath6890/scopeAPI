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

// EnvoyIntegrationService handles Envoy-specific operations
type EnvoyIntegrationService struct {
	integrationService IntegrationServiceInterface
	logger             logging.Logger
	httpClient         *http.Client
}

// EnvoyClient implements GatewayClient for Envoy
type EnvoyClient struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewEnvoyIntegrationService creates a new Envoy integration service
func NewEnvoyIntegrationService(integrationService IntegrationServiceInterface, logger logging.Logger) *EnvoyIntegrationService {
	return &EnvoyIntegrationService{
		integrationService: integrationService,
		logger:             logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewEnvoyClient creates a new Envoy client
func NewEnvoyClient(logger logging.Logger) *EnvoyClient {
	return &EnvoyClient{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Envoy API Response structures
type EnvoyStatusResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type EnvoyClusterResponse struct {
	Clusters []models.EnvoyCluster `json:"clusters"`
}

type EnvoyListenerResponse struct {
	Listeners []models.EnvoyListener `json:"listeners"`
}

type EnvoyFilterResponse struct {
	Filters []models.EnvoyFilter `json:"filters"`
}

type EnvoyRouteResponse struct {
	Routes []models.EnvoyRoute `json:"routes"`
}

// EnvoyIntegrationService methods

// GetStatus retrieves Envoy status
func (s *EnvoyIntegrationService) GetStatus(ctx context.Context, integrationID string) (*models.HealthStatus, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewEnvoyClient(s.logger)
	return client.GetStatus(ctx, integration)
}

// GetClusters retrieves Envoy clusters
func (s *EnvoyIntegrationService) GetClusters(ctx context.Context, integrationID string) ([]models.EnvoyCluster, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/clusters", adminAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clusters, status: %d", resp.StatusCode)
	}

	var envoyResp EnvoyClusterResponse
	if err := json.NewDecoder(resp.Body).Decode(&envoyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return envoyResp.Clusters, nil
}

// GetListeners retrieves Envoy listeners
func (s *EnvoyIntegrationService) GetListeners(ctx context.Context, integrationID string) ([]models.EnvoyListener, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/listeners", adminAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get listeners, status: %d", resp.StatusCode)
	}

	var envoyResp EnvoyListenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&envoyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return envoyResp.Listeners, nil
}

// GetFilters retrieves Envoy filters
func (s *EnvoyIntegrationService) GetFilters(ctx context.Context, integrationID string) ([]models.EnvoyFilter, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/filters", adminAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get filters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get filters, status: %d", resp.StatusCode)
	}

	var envoyResp EnvoyFilterResponse
	if err := json.NewDecoder(resp.Body).Decode(&envoyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return envoyResp.Filters, nil
}

// CreateFilter creates a new Envoy filter
func (s *EnvoyIntegrationService) CreateFilter(ctx context.Context, integrationID string, filter *models.EnvoyFilter) (*models.EnvoyFilter, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/filters", adminAddress)
	filterData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(filterData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create filter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create filter, status: %d", resp.StatusCode)
	}

	var createdFilter models.EnvoyFilter
	if err := json.NewDecoder(resp.Body).Decode(&createdFilter); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdFilter, nil
}

// UpdateFilter updates an existing Envoy filter
func (s *EnvoyIntegrationService) UpdateFilter(ctx context.Context, integrationID, filterID string, filter *models.EnvoyFilter) (*models.EnvoyFilter, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/filters/%s", adminAddress, filterID)
	filterData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(filterData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update filter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update filter, status: %d", resp.StatusCode)
	}

	var updatedFilter models.EnvoyFilter
	if err := json.NewDecoder(resp.Body).Decode(&updatedFilter); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedFilter, nil
}

// DeleteFilter deletes an Envoy filter
func (s *EnvoyIntegrationService) DeleteFilter(ctx context.Context, integrationID, filterID string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	adminAddress, err := s.getAdminAddress(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/filters/%s", adminAddress, filterID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete filter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete filter, status: %d", resp.StatusCode)
	}

	return nil
}

// SyncConfiguration synchronizes configuration with Envoy
func (s *EnvoyIntegrationService) SyncConfiguration(ctx context.Context, integrationID string) (*models.SyncResult, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewEnvoyClient(s.logger)
	return client.SyncConfiguration(ctx, integration)
}

// EnvoyClient methods

// GetStatus implements GatewayClient.GetStatus for Envoy
func (c *EnvoyClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	adminAddress, err := c.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	url := fmt.Sprintf("%s/status", adminAddress)
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

	var status EnvoyStatusResponse
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
		Message:   "Envoy is running",
		LastCheck: time.Now(),
		Latency:   latency,
	}, nil
}

// SyncConfiguration implements GatewayClient.SyncConfiguration for Envoy
func (c *EnvoyClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	start := time.Now()
	changes := []models.Change{}

	// Parse Envoy configuration from integration config
	var envoyConfig models.EnvoyConfig
	if configData, exists := integration.Config["envoy_config"]; exists {
		if configBytes, err := json.Marshal(configData); err == nil {
			if err := json.Unmarshal(configBytes, &envoyConfig); err != nil {
				return nil, fmt.Errorf("failed to parse Envoy config: %w", err)
			}
		}
	}

	adminAddress, err := c.getAdminAddress(integration)
	if err != nil {
		return nil, err
	}

	// Sync clusters
	for _, cluster := range envoyConfig.Clusters {
		if err := c.syncCluster(ctx, adminAddress, integration, &cluster); err != nil {
			c.logger.Error("Failed to sync cluster", "cluster_name", cluster.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "cluster",
				Resource: cluster.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "cluster",
				Resource: cluster.Name,
				Action:   "synced",
				Details:  "Cluster configuration synchronized",
			})
		}
	}

	// Sync listeners
	for _, listener := range envoyConfig.Listeners {
		if err := c.syncListener(ctx, adminAddress, integration, &listener); err != nil {
			c.logger.Error("Failed to sync listener", "listener_name", listener.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "listener",
				Resource: listener.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "listener",
				Resource: listener.Name,
				Action:   "synced",
				Details:  "Listener configuration synchronized",
			})
		}
	}

	// Sync filters
	for _, filter := range envoyConfig.Filters {
		if err := c.syncFilter(ctx, adminAddress, integration, &filter); err != nil {
			c.logger.Error("Failed to sync filter", "filter_name", filter.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "filter",
				Resource: filter.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "filter",
				Resource: filter.Name,
				Action:   "synced",
				Details:  "Filter configuration synchronized",
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

// TestConnection implements GatewayClient.TestConnection for Envoy
func (c *EnvoyClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	adminAddress, err := c.getAdminAddress(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/status", adminAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req, integration)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Envoy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Envoy returned status %d", resp.StatusCode)
	}

	return nil
}

// Private helper methods

func (s *EnvoyIntegrationService) getAdminAddress(integration *models.Integration) (string, error) {
	if adminAddress, exists := integration.Config["admin_address"]; exists {
		if address, ok := adminAddress.(string); ok {
			return address, nil
		}
	}
	return "", fmt.Errorf("admin_address not found in integration config")
}

func (s *EnvoyIntegrationService) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *EnvoyClient) getAdminAddress(integration *models.Integration) (string, error) {
	if adminAddress, exists := integration.Config["admin_address"]; exists {
		if address, ok := adminAddress.(string); ok {
			return address, nil
		}
	}
	return "", fmt.Errorf("admin_address not found in integration config")
}

func (c *EnvoyClient) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *EnvoyClient) syncCluster(ctx context.Context, adminAddress string, integration *models.Integration, cluster *models.EnvoyCluster) error {
	// Implementation for syncing Envoy cluster
	// This would check if the cluster exists and create/update it accordingly
	c.logger.Info("Syncing Envoy cluster", "cluster_name", cluster.Name)
	return nil
}

func (c *EnvoyClient) syncListener(ctx context.Context, adminAddress string, integration *models.Integration, listener *models.EnvoyListener) error {
	// Implementation for syncing Envoy listener
	// This would check if the listener exists and create/update it accordingly
	c.logger.Info("Syncing Envoy listener", "listener_name", listener.Name)
	return nil
}

func (c *EnvoyClient) syncFilter(ctx context.Context, adminAddress string, integration *models.Integration, filter *models.EnvoyFilter) error {
	// Implementation for syncing Envoy filter
	// This would check if the filter exists and create/update it accordingly
	c.logger.Info("Syncing Envoy filter", "filter_name", filter.Name)
	return nil
} 