package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"scopeapi.local/backend/services/gateway-integration/internal/models"
	"scopeapi.local/backend/services/gateway-integration/internal/repository"
	"scopeapi.local/backend/shared/logging"
)

// ConfigService handles business logic for gateway configurations
type ConfigService struct {
	configRepo repository.ConfigRepository
	logger     logging.Logger
}

// NewConfigService creates a new ConfigService instance
func NewConfigService(configRepo repository.ConfigRepository, logger logging.Logger) *ConfigService {
	return &ConfigService{
		configRepo: configRepo,
		logger:     logger,
	}
}

// CreateConfig creates a new gateway configuration
func (s *ConfigService) CreateConfig(ctx context.Context, config *models.GatewayConfig) error {
	s.logger.Info("Creating new gateway configuration", 
		"integration_id", config.IntegrationID, 
		"config_type", config.ConfigType)

	// Validate configuration data
	if err := s.validateConfigData(config); err != nil {
		return fmt.Errorf("invalid configuration data: %w", err)
	}

	// Set initial values
	config.Version = 1
	config.Status = "draft"
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Create configuration in database
	if err := s.configRepo.CreateConfig(ctx, config); err != nil {
		s.logger.Error("Failed to create configuration", "error", err)
		return fmt.Errorf("failed to create configuration: %w", err)
	}

	s.logger.Info("Successfully created gateway configuration", 
		"config_id", config.ID, 
		"integration_id", config.IntegrationID)

	return nil
}

// GetConfig retrieves a gateway configuration by ID
func (s *ConfigService) GetConfig(ctx context.Context, id int64) (*models.GatewayConfig, error) {
	s.logger.Info("Retrieving gateway configuration", "config_id", id)

	config, err := s.configRepo.GetConfig(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve configuration", "config_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	if config == nil {
		s.logger.Warn("Configuration not found", "config_id", id)
		return nil, nil
	}

	s.logger.Info("Successfully retrieved gateway configuration", 
		"config_id", config.ID, 
		"integration_id", config.IntegrationID)

	return config, nil
}

// GetConfigs retrieves all configurations for a given integration
func (s *ConfigService) GetConfigs(ctx context.Context, integrationID int64, configType string) ([]*models.GatewayConfig, error) {
	s.logger.Info("Retrieving gateway configurations", 
		"integration_id", integrationID, 
		"config_type", configType)

	configs, err := s.configRepo.GetConfigs(ctx, integrationID, configType)
	if err != nil {
		s.logger.Error("Failed to retrieve configurations", 
			"integration_id", integrationID, "error", err)
		return nil, fmt.Errorf("failed to retrieve configurations: %w", err)
	}

	s.logger.Info("Successfully retrieved gateway configurations", 
		"integration_id", integrationID, 
		"count", len(configs))

	return configs, nil
}

// UpdateConfig updates an existing gateway configuration
func (s *ConfigService) UpdateConfig(ctx context.Context, config *models.GatewayConfig) error {
	s.logger.Info("Updating gateway configuration", 
		"config_id", config.ID, 
		"integration_id", config.IntegrationID)

	// Validate configuration data
	if err := s.validateConfigData(config); err != nil {
		return fmt.Errorf("invalid configuration data: %w", err)
	}

	// Get existing configuration to check version
	existingConfig, err := s.configRepo.GetConfig(ctx, config.ID)
	if err != nil {
		s.logger.Error("Failed to retrieve existing configuration", "config_id", config.ID, "error", err)
		return fmt.Errorf("failed to retrieve existing configuration: %w", err)
	}

	if existingConfig == nil {
		s.logger.Warn("Configuration not found for update", "config_id", config.ID)
		return fmt.Errorf("configuration not found")
	}

	// Increment version
	config.Version = existingConfig.Version + 1
	config.UpdatedAt = time.Now()

	// Update configuration in database
	if err := s.configRepo.UpdateConfig(ctx, config); err != nil {
		s.logger.Error("Failed to update configuration", "config_id", config.ID, "error", err)
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	s.logger.Info("Successfully updated gateway configuration", 
		"config_id", config.ID, 
		"version", config.Version)

	return nil
}

// DeleteConfig deletes a gateway configuration
func (s *ConfigService) DeleteConfig(ctx context.Context, id int64) error {
	s.logger.Info("Deleting gateway configuration", "config_id", id)

	// Check if configuration exists
	existingConfig, err := s.configRepo.GetConfig(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve configuration for deletion", "config_id", id, "error", err)
		return fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	if existingConfig == nil {
		s.logger.Warn("Configuration not found for deletion", "config_id", id)
		return fmt.Errorf("configuration not found")
	}

	// Delete configuration from database
	if err := s.configRepo.DeleteConfig(ctx, id); err != nil {
		s.logger.Error("Failed to delete configuration", "config_id", id, "error", err)
		return fmt.Errorf("failed to delete configuration: %w", err)
	}

	s.logger.Info("Successfully deleted gateway configuration", 
		"config_id", id, 
		"integration_id", existingConfig.IntegrationID)

	return nil
}

// ValidateConfig validates a configuration without saving it
func (s *ConfigService) ValidateConfig(ctx context.Context, config *models.GatewayConfig) error {
	s.logger.Info("Validating gateway configuration", 
		"integration_id", config.IntegrationID, 
		"config_type", config.ConfigType)

	if err := s.validateConfigData(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	s.logger.Info("Configuration validation successful", 
		"integration_id", config.IntegrationID, 
		"config_type", config.ConfigType)

	return nil
}

// DeployConfig deploys a configuration by updating its status
func (s *ConfigService) DeployConfig(ctx context.Context, id int64) error {
	s.logger.Info("Deploying gateway configuration", "config_id", id)

	// Update configuration status to deployed
	if err := s.configRepo.UpdateConfigStatus(ctx, id, "deployed"); err != nil {
		s.logger.Error("Failed to deploy configuration", "config_id", id, "error", err)
		return fmt.Errorf("failed to deploy configuration: %w", err)
	}

	s.logger.Info("Successfully deployed gateway configuration", "config_id", id)
	return nil
}

// RollbackConfig rolls back a configuration to a previous version
func (s *ConfigService) RollbackConfig(ctx context.Context, id int64, targetVersion int) error {
	s.logger.Info("Rolling back gateway configuration", 
		"config_id", id, 
		"target_version", targetVersion)

	// Get the target version configuration
	configs, err := s.configRepo.GetConfigs(ctx, 0, "") // Get all configs for this integration
	if err != nil {
		s.logger.Error("Failed to retrieve configurations for rollback", "config_id", id, "error", err)
		return fmt.Errorf("failed to retrieve configurations: %w", err)
	}

	// Find the target version
	var targetConfig *models.GatewayConfig
	for _, config := range configs {
		if config.Version == targetVersion {
			targetConfig = config
			break
		}
	}

	if targetConfig == nil {
		s.logger.Warn("Target version not found for rollback", 
			"config_id", id, 
			"target_version", targetVersion)
		return fmt.Errorf("target version not found")
	}

	// Create a new configuration with the target version data
	newConfig := &models.GatewayConfig{
		IntegrationID: targetConfig.IntegrationID,
		ConfigType:    targetConfig.ConfigType,
		ConfigData:    targetConfig.ConfigData,
		Status:        "draft",
	}

	if err := s.CreateConfig(ctx, newConfig); err != nil {
		s.logger.Error("Failed to create rollback configuration", "config_id", id, "error", err)
		return fmt.Errorf("failed to create rollback configuration: %w", err)
	}

	s.logger.Info("Successfully rolled back gateway configuration", 
		"config_id", id, 
		"new_config_id", newConfig.ID, 
		"target_version", targetVersion)

	return nil
}

// validateConfigData validates the configuration data structure
func (s *ConfigService) validateConfigData(config *models.GatewayConfig) error {
	if config.IntegrationID <= 0 {
		return fmt.Errorf("invalid integration ID")
	}

	if config.ConfigType == "" {
		return fmt.Errorf("config type is required")
	}

	if config.ConfigData == nil {
		return fmt.Errorf("config data is required")
	}

	// Validate that config data can be marshaled to JSON
	if _, err := json.Marshal(config.ConfigData); err != nil {
		return fmt.Errorf("invalid config data format: %w", err)
	}

	return nil
} 