package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

// IntegrationServiceInterface defines the interface for integration service
type IntegrationServiceInterface interface {
	CreateIntegration(ctx context.Context, integration *models.Integration) error
	GetIntegration(ctx context.Context, id string) (*models.Integration, error)
	GetIntegrations(ctx context.Context, filters map[string]interface{}) ([]*models.Integration, error)
	UpdateIntegration(ctx context.Context, integration *models.Integration) error
	DeleteIntegration(ctx context.Context, id string) error
	TestIntegration(ctx context.Context, id string) (*models.HealthStatus, error)
	SyncIntegration(ctx context.Context, id string) (*models.SyncResult, error)
	ProcessGatewayEvent(ctx context.Context, eventData []byte) error
	ProcessSecurityEvent(ctx context.Context, eventData []byte) error
}

// IntegrationService implements the integration service
type IntegrationService struct {
	repo           repository.IntegrationRepository
	kafkaProducer  kafka.Producer
	logger         logging.Logger
	gatewayClients map[models.GatewayType]GatewayClient
}

// GatewayClient interface for different gateway implementations
type GatewayClient interface {
	GetStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error)
	SyncConfiguration(ctx context.Context, integration *models.Integration) (*models.SyncResult, error)
	TestConnection(ctx context.Context, integration *models.Integration) error
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(repo repository.IntegrationRepository, kafkaProducer kafka.Producer, logger logging.Logger) *IntegrationService {
	service := &IntegrationService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		gatewayClients: make(map[models.GatewayType]GatewayClient),
	}

	// Initialize gateway clients
	service.gatewayClients[models.GatewayTypeKong] = NewKongClient(logger)
	service.gatewayClients[models.GatewayTypeNginx] = NewNginxClient(logger)
	service.gatewayClients[models.GatewayTypeTraefik] = NewTraefikClient(logger)
	service.gatewayClients[models.GatewayTypeEnvoy] = NewEnvoyClient(logger)
	service.gatewayClients[models.GatewayTypeHAProxy] = NewHAProxyClient(logger)

	return service
}

// CreateIntegration creates a new gateway integration
func (s *IntegrationService) CreateIntegration(ctx context.Context, integration *models.Integration) error {
	// Validate integration
	if err := s.validateIntegration(integration); err != nil {
		return fmt.Errorf("invalid integration: %w", err)
	}

	// Generate ID if not provided
	if integration.ID == "" {
		integration.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	integration.CreatedAt = now
	integration.UpdatedAt = now

	// Set initial status
	if integration.Status == "" {
		integration.Status = models.IntegrationStatusPending
	}

	// Test connection before saving
	if err := s.testConnection(ctx, integration); err != nil {
		integration.Status = models.IntegrationStatusError
		s.logger.Warn("Integration connection test failed", "integration_id", integration.ID, "error", err)
	}

	// Save to database
	if err := s.repo.CreateIntegration(ctx, integration); err != nil {
		return fmt.Errorf("failed to create integration: %w", err)
	}

	// Publish integration event
	s.publishIntegrationEvent(ctx, "integration_created", integration, nil)

	s.logger.Info("Integration created successfully", "integration_id", integration.ID, "type", integration.Type)
	return nil
}

// GetIntegration retrieves an integration by ID
func (s *IntegrationService) GetIntegration(ctx context.Context, id string) (*models.Integration, error) {
	integration, err := s.repo.GetIntegration(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	// Update health status
	if integration.Status == models.IntegrationStatusActive {
		if health, err := s.getHealthStatus(ctx, integration); err == nil {
			integration.Health = health
		}
	}

	return integration, nil
}

// GetIntegrations retrieves integrations with optional filters
func (s *IntegrationService) GetIntegrations(ctx context.Context, filters map[string]interface{}) ([]*models.Integration, error) {
	integrations, err := s.repo.GetIntegrations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get integrations: %w", err)
	}

	// Update health status for active integrations
	for _, integration := range integrations {
		if integration.Status == models.IntegrationStatusActive {
			if health, err := s.getHealthStatus(ctx, integration); err == nil {
				integration.Health = health
			}
		}
	}

	return integrations, nil
}

// UpdateIntegration updates an existing integration
func (s *IntegrationService) UpdateIntegration(ctx context.Context, integration *models.Integration) error {
	// Validate integration
	if err := s.validateIntegration(integration); err != nil {
		return fmt.Errorf("invalid integration: %w", err)
	}

	// Get existing integration
	existing, err := s.repo.GetIntegration(ctx, integration.ID)
	if err != nil {
		return fmt.Errorf("integration not found: %w", err)
	}

	// Update timestamp
	integration.UpdatedAt = time.Now()
	integration.CreatedAt = existing.CreatedAt

	// Test connection if configuration changed
	if s.configurationChanged(existing, integration) {
		if err := s.testConnection(ctx, integration); err != nil {
			integration.Status = models.IntegrationStatusError
			s.logger.Warn("Integration connection test failed after update", "integration_id", integration.ID, "error", err)
		} else {
			integration.Status = models.IntegrationStatusActive
		}
	}

	// Save to database
	if err := s.repo.UpdateIntegration(ctx, integration); err != nil {
		return fmt.Errorf("failed to update integration: %w", err)
	}

	// Publish integration event
	s.publishIntegrationEvent(ctx, "integration_updated", integration, nil)

	s.logger.Info("Integration updated successfully", "integration_id", integration.ID)
	return nil
}

// DeleteIntegration deletes an integration
func (s *IntegrationService) DeleteIntegration(ctx context.Context, id string) error {
	// Get integration before deletion for event
	integration, err := s.repo.GetIntegration(ctx, id)
	if err != nil {
		return fmt.Errorf("integration not found: %w", err)
	}

	// Delete from database
	if err := s.repo.DeleteIntegration(ctx, id); err != nil {
		return fmt.Errorf("failed to delete integration: %w", err)
	}

	// Publish integration event
	s.publishIntegrationEvent(ctx, "integration_deleted", integration, nil)

	s.logger.Info("Integration deleted successfully", "integration_id", id)
	return nil
}

// TestIntegration tests the connection to a gateway
func (s *IntegrationService) TestIntegration(ctx context.Context, id string) (*models.HealthStatus, error) {
	integration, err := s.repo.GetIntegration(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("integration not found: %w", err)
	}

	// Test connection
	if err := s.testConnection(ctx, integration); err != nil {
		health := &models.HealthStatus{
			Status:    "unhealthy",
			Message:   err.Error(),
			LastCheck: time.Now(),
		}

		// Update integration status
		integration.Status = models.IntegrationStatusError
		integration.Health = health
		s.repo.UpdateIntegration(ctx, integration)

		return health, nil
	}

	// Get health status
	health, err := s.getHealthStatus(ctx, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to get health status: %w", err)
	}

	// Update integration status and health
	integration.Status = models.IntegrationStatusActive
	integration.Health = health
	s.repo.UpdateIntegration(ctx, integration)

	return health, nil
}

// SyncIntegration synchronizes configuration with the gateway
func (s *IntegrationService) SyncIntegration(ctx context.Context, id string) (*models.SyncResult, error) {
	integration, err := s.repo.GetIntegration(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("integration not found: %w", err)
	}

	// Get gateway client
	client, exists := s.gatewayClients[integration.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported gateway type: %s", integration.Type)
	}

	// Sync configuration
	result, err := client.SyncConfiguration(ctx, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to sync configuration: %w", err)
	}

	// Update last sync time
	now := time.Now()
	integration.LastSync = &now
	s.repo.UpdateIntegration(ctx, integration)

	// Publish sync event
	s.publishIntegrationEvent(ctx, "integration_synced", integration, result)

	s.logger.Info("Integration synced successfully", "integration_id", id, "changes", len(result.Changes))
	return result, nil
}

// ProcessGatewayEvent processes gateway events from Kafka
func (s *IntegrationService) ProcessGatewayEvent(ctx context.Context, eventData []byte) error {
	var event models.IntegrationEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal gateway event: %w", err)
	}

	s.logger.Info("Processing gateway event", "event_id", event.ID, "type", event.Type)

	// Process event based on type
	switch event.Type {
	case "configuration_change":
		return s.handleConfigurationChange(ctx, &event)
	case "health_check":
		return s.handleHealthCheck(ctx, &event)
	case "sync_request":
		return s.handleSyncRequest(ctx, &event)
	default:
		s.logger.Warn("Unknown gateway event type", "event_type", event.Type)
		return nil
	}
}

// ProcessSecurityEvent processes security events for gateway integration
func (s *IntegrationService) ProcessSecurityEvent(ctx context.Context, eventData []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal security event: %w", err)
	}

	s.logger.Info("Processing security event for gateway integration", "event_type", event["type"])

	// Handle security events that affect gateway configuration
	// For example, blocking rules, rate limiting updates, etc.
	return s.handleSecurityEvent(ctx, event)
}

// Private helper methods

func (s *IntegrationService) validateIntegration(integration *models.Integration) error {
	if integration.Name == "" {
		return fmt.Errorf("integration name is required")
	}

	if integration.Type == "" {
		return fmt.Errorf("gateway type is required")
	}

	// Validate gateway type
	switch integration.Type {
	case models.GatewayTypeKong, models.GatewayTypeNginx, models.GatewayTypeTraefik, models.GatewayTypeEnvoy, models.GatewayTypeHAProxy:
		// Valid types
	default:
		return fmt.Errorf("unsupported gateway type: %s", integration.Type)
	}

	if len(integration.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint is required")
	}

	return nil
}

func (s *IntegrationService) testConnection(ctx context.Context, integration *models.Integration) error {
	client, exists := s.gatewayClients[integration.Type]
	if !exists {
		return fmt.Errorf("unsupported gateway type: %s", integration.Type)
	}

	return client.TestConnection(ctx, integration)
}

func (s *IntegrationService) getHealthStatus(ctx context.Context, integration *models.Integration) (*models.HealthStatus, error) {
	client, exists := s.gatewayClients[integration.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported gateway type: %s", integration.Type)
	}

	return client.GetStatus(ctx, integration)
}

func (s *IntegrationService) configurationChanged(existing, updated *models.Integration) bool {
	// Compare relevant fields to determine if configuration changed
	if existing.Type != updated.Type {
		return true
	}

	// Compare endpoints
	if len(existing.Endpoints) != len(updated.Endpoints) {
		return true
	}

	// Compare config (simplified comparison)
	existingConfig, _ := json.Marshal(existing.Config)
	updatedConfig, _ := json.Marshal(updated.Config)
	return string(existingConfig) != string(updatedConfig)
}

func (s *IntegrationService) publishIntegrationEvent(ctx context.Context, eventType string, integration *models.Integration, data interface{}) {
	event := models.IntegrationEvent{
		ID:            uuid.New().String(),
		IntegrationID: integration.ID,
		Type:          eventType,
		Timestamp:     time.Now(),
		Status:        "success",
		Data:          make(map[string]interface{}),
	}

	if data != nil {
		if dataBytes, err := json.Marshal(data); err == nil {
			var dataMap map[string]interface{}
			if json.Unmarshal(dataBytes, &dataMap) == nil {
				event.Data = dataMap
			}
		}
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("Failed to marshal integration event", "error", err)
		return
	}

	message := kafka.Message{
		Topic: "gateway_events",
		Key:   []byte(integration.ID),
		Value: eventBytes,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to publish integration event", "error", err)
	}
}

func (s *IntegrationService) handleConfigurationChange(ctx context.Context, event *models.IntegrationEvent) error {
	// Handle configuration change events
	// This could trigger a sync or update of the gateway configuration
	s.logger.Info("Handling configuration change event", "integration_id", event.IntegrationID)
	return nil
}

func (s *IntegrationService) handleHealthCheck(ctx context.Context, event *models.IntegrationEvent) error {
	// Handle health check events
	// This could update the health status of the integration
	s.logger.Info("Handling health check event", "integration_id", event.IntegrationID)
	return nil
}

func (s *IntegrationService) handleSyncRequest(ctx context.Context, event *models.IntegrationEvent) error {
	// Handle sync request events
	// This could trigger a configuration sync
	s.logger.Info("Handling sync request event", "integration_id", event.IntegrationID)
	return nil
}

func (s *IntegrationService) handleSecurityEvent(ctx context.Context, event map[string]interface{}) error {
	// Handle security events that affect gateway configuration
	// For example, updating blocking rules, rate limits, etc.
	eventType, _ := event["type"].(string)
	s.logger.Info("Handling security event for gateway integration", "event_type", eventType)
	return nil
} 