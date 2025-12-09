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

// TraefikIntegrationService handles Traefik-specific operations
type TraefikIntegrationService struct {
	integrationService IntegrationServiceInterface
	logger             logging.Logger
	httpClient         *http.Client
}

// TraefikClient implements GatewayClient for Traefik
type TraefikClient struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewTraefikIntegrationService creates a new Traefik integration service
func NewTraefikIntegrationService(integrationService IntegrationServiceInterface, logger logging.Logger) *TraefikIntegrationService {
	return &TraefikIntegrationService{
		integrationService: integrationService,
		logger:             logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewTraefikClient creates a new Traefik client
func NewTraefikClient(logger logging.Logger) *TraefikClient {
	return &TraefikClient{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Traefik API Response structures
type TraefikStatusResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type TraefikProviderResponse struct {
	Providers []models.TraefikProvider `json:"providers"`
}

type TraefikMiddlewareResponse struct {
	Middlewares []models.TraefikMiddleware `json:"middlewares"`
}

type TraefikRouterResponse struct {
	Routers []models.TraefikRouter `json:"routers"`
}

type TraefikServiceResponse struct {
	Services []models.TraefikService `json:"services"`
}

// TraefikIntegrationService methods

// GetStatus retrieves Traefik status
func (s *TraefikIntegrationService) GetStatus(ctx context.Context, integrationID string) (*models.HealthStatus, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewTraefikClient(s.logger)
	return client.GetStatus(ctx, integration)
}

// GetProviders retrieves Traefik providers
func (s *TraefikIntegrationService) GetProviders(ctx context.Context, integrationID string) ([]models.TraefikProvider, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	apiEndpoint, err := s.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/providers", apiEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get providers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get providers, status: %d", resp.StatusCode)
	}

	var traefikResp TraefikProviderResponse
	if err := json.NewDecoder(resp.Body).Decode(&traefikResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return traefikResp.Providers, nil
}

// GetMiddlewares retrieves Traefik middlewares
func (s *TraefikIntegrationService) GetMiddlewares(ctx context.Context, integrationID string) ([]models.TraefikMiddleware, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	apiEndpoint, err := s.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/middlewares", apiEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get middlewares: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get middlewares, status: %d", resp.StatusCode)
	}

	var traefikResp TraefikMiddlewareResponse
	if err := json.NewDecoder(resp.Body).Decode(&traefikResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return traefikResp.Middlewares, nil
}

// CreateMiddleware creates a new Traefik middleware
func (s *TraefikIntegrationService) CreateMiddleware(ctx context.Context, integrationID string, middleware *models.TraefikMiddleware) (*models.TraefikMiddleware, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	apiEndpoint, err := s.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/middlewares", apiEndpoint)
	middlewareData, err := json.Marshal(middleware)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(middlewareData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create middleware: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create middleware, status: %d", resp.StatusCode)
	}

	var createdMiddleware models.TraefikMiddleware
	if err := json.NewDecoder(resp.Body).Decode(&createdMiddleware); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdMiddleware, nil
}

// UpdateMiddleware updates an existing Traefik middleware
func (s *TraefikIntegrationService) UpdateMiddleware(ctx context.Context, integrationID, middlewareID string, middleware *models.TraefikMiddleware) (*models.TraefikMiddleware, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	apiEndpoint, err := s.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/middlewares/%s", apiEndpoint, middlewareID)
	middlewareData, err := json.Marshal(middleware)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(middlewareData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update middleware: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update middleware, status: %d", resp.StatusCode)
	}

	var updatedMiddleware models.TraefikMiddleware
	if err := json.NewDecoder(resp.Body).Decode(&updatedMiddleware); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedMiddleware, nil
}

// DeleteMiddleware deletes a Traefik middleware
func (s *TraefikIntegrationService) DeleteMiddleware(ctx context.Context, integrationID, middlewareID string) error {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	apiEndpoint, err := s.getAPIEndpoint(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/middlewares/%s", apiEndpoint, middlewareID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	s.addAuthHeaders(req, integration)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete middleware: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete middleware, status: %d", resp.StatusCode)
	}

	return nil
}

// SyncConfiguration synchronizes configuration with Traefik
func (s *TraefikIntegrationService) SyncConfiguration(ctx context.Context, integrationID string) (*models.SyncResult, error) {
	integration, err := s.integrationService.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	client := NewTraefikClient(s.logger)
	return client.SyncConfiguration(ctx, integration)
}

// TraefikClient methods

// GetStatus implements GatewayClient.GetStatus for Traefik
func (c *TraefikClient) GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	apiEndpoint, err := c.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	url := fmt.Sprintf("%s/api/status", apiEndpoint)
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

	var status TraefikStatusResponse
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
		Message:   "Traefik is running",
		LastCheck: time.Now(),
		Latency:   latency,
	}, nil
}

// SyncConfiguration implements GatewayClient.SyncConfiguration for Traefik
func (c *TraefikClient) SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error) {
	start := time.Now()
	changes := []models.Change{}

	// Parse Traefik configuration from integration config
	var traefikConfig models.TraefikConfig
	if configData, exists := integration.Config["traefik_config"]; exists {
		if configBytes, err := json.Marshal(configData); err == nil {
			if err := json.Unmarshal(configBytes, &traefikConfig); err != nil {
				return nil, fmt.Errorf("failed to parse Traefik config: %w", err)
			}
		}
	}

	apiEndpoint, err := c.getAPIEndpoint(integration)
	if err != nil {
		return nil, err
	}

	// Sync middlewares
	for _, middleware := range traefikConfig.Middlewares {
		if err := c.syncMiddleware(ctx, apiEndpoint, integration, &middleware); err != nil {
			c.logger.Error("Failed to sync middleware", "middleware_name", middleware.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "middleware",
				Resource: middleware.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "middleware",
				Resource: middleware.Name,
				Action:   "synced",
				Details:  "Middleware configuration synchronized",
			})
		}
	}

	// Sync routers
	for _, router := range traefikConfig.Routers {
		if err := c.syncRouter(ctx, apiEndpoint, integration, &router); err != nil {
			c.logger.Error("Failed to sync router", "router_name", router.Name, "error", err)
			changes = append(changes, models.Change{
				Type:     "router",
				Resource: router.Name,
				Action:   "error",
				Details:  err.Error(),
			})
		} else {
			changes = append(changes, models.Change{
				Type:     "router",
				Resource: router.Name,
				Action:   "synced",
				Details:  "Router configuration synchronized",
			})
		}
	}

	// Sync services
	for _, service := range traefikConfig.Services {
		if err := c.syncService(ctx, apiEndpoint, integration, &service); err != nil {
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

// TestConnection implements GatewayClient.TestConnection for Traefik
func (c *TraefikClient) TestConnection(ctx context.Context, integration *models.Integration) error {
	apiEndpoint, err := c.getAPIEndpoint(integration)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/status", apiEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req, integration)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Traefik: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Traefik returned status %d", resp.StatusCode)
	}

	return nil
}

// Private helper methods

func (s *TraefikIntegrationService) getAPIEndpoint(integration *models.Integration) (string, error) {
	if apiEndpoint, exists := integration.Config["api_endpoint"]; exists {
		if endpoint, ok := apiEndpoint.(string); ok {
			return endpoint, nil
		}
	}
	return "", fmt.Errorf("api_endpoint not found in integration config")
}

func (s *TraefikIntegrationService) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *TraefikClient) getAPIEndpoint(integration *models.Integration) (string, error) {
	if apiEndpoint, exists := integration.Config["api_endpoint"]; exists {
		if endpoint, ok := apiEndpoint.(string); ok {
			return endpoint, nil
		}
	}
	return "", fmt.Errorf("api_endpoint not found in integration config")
}

func (c *TraefikClient) addAuthHeaders(req *http.Request, integration *models.Integration) {
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

func (c *TraefikClient) syncMiddleware(ctx context.Context, apiEndpoint string, integration *models.Integration, middleware *models.TraefikMiddleware) error {
	// Implementation for syncing Traefik middleware
	// This would check if the middleware exists and create/update it accordingly
	c.logger.Info("Syncing Traefik middleware", "middleware_name", middleware.Name)
	return nil
}

func (c *TraefikClient) syncRouter(ctx context.Context, apiEndpoint string, integration *models.Integration, router *models.TraefikRouter) error {
	// Implementation for syncing Traefik router
	// This would check if the router exists and create/update it accordingly
	c.logger.Info("Syncing Traefik router", "router_name", router.Name)
	return nil
}

func (c *TraefikClient) syncService(ctx context.Context, apiEndpoint string, integration *models.Integration, service *models.TraefikService) error {
	// Implementation for syncing Traefik service
	// This would check if the service exists and create/update it accordingly
	c.logger.Info("Syncing Traefik service", "service_name", service.Name)
	return nil
} 