package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
)

// ConfigRepository handles database operations for gateway configurations
type ConfigRepository struct {
	db *sqlx.DB
}

// NewConfigRepository creates a new ConfigRepository instance
func NewConfigRepository(db *sqlx.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// CreateConfig creates a new gateway configuration
func (r *ConfigRepository) CreateConfig(ctx context.Context, config *models.GatewayConfig) error {
	query := `
		INSERT INTO gateway_configs (
			integration_id, config_type, config_data, version, 
			status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id`

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	configData, err := json.Marshal(config.ConfigData)
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		config.IntegrationID, config.ConfigType, configData, config.Version,
		config.Status, config.CreatedAt, config.UpdatedAt,
	).Scan(&config.ID)

	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	return nil
}

// GetConfig retrieves a gateway configuration by ID
func (r *ConfigRepository) GetConfig(ctx context.Context, id int64) (*models.GatewayConfig, error) {
	query := `
		SELECT id, integration_id, config_type, config_data, version,
			   status, created_at, updated_at
		FROM gateway_configs WHERE id = $1`

	var config models.GatewayConfig
	var configData []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID, &config.IntegrationID, &config.ConfigType, &configData,
		&config.Version, &config.Status, &config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Unmarshal config data
	if err := json.Unmarshal(configData, &config.ConfigData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return &config, nil
}

// GetConfigs retrieves all configurations for a given integration
func (r *ConfigRepository) GetConfigs(ctx context.Context, integrationID int64, configType string) ([]*models.GatewayConfig, error) {
	query := `
		SELECT id, integration_id, config_type, config_data, version,
			   status, created_at, updated_at
		FROM gateway_configs 
		WHERE integration_id = $1 AND ($2 = '' OR config_type = $2)
		ORDER BY version DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, integrationID, configType)
	if err != nil {
		return nil, fmt.Errorf("failed to get configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.GatewayConfig
	for rows.Next() {
		var config models.GatewayConfig
		var configData []byte

		err := rows.Scan(
			&config.ID, &config.IntegrationID, &config.ConfigType, &configData,
			&config.Version, &config.Status, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}

		// Unmarshal config data
		if err := json.Unmarshal(configData, &config.ConfigData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

// UpdateConfig updates an existing gateway configuration
func (r *ConfigRepository) UpdateConfig(ctx context.Context, config *models.GatewayConfig) error {
	query := `
		UPDATE gateway_configs SET
			config_data = $1, version = $2, status = $3, updated_at = $4
		WHERE id = $5`

	config.UpdatedAt = time.Now()

	configData, err := json.Marshal(config.ConfigData)
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		configData, config.Version, config.Status, config.UpdatedAt, config.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("config not found")
	}

	return nil
}

// DeleteConfig deletes a gateway configuration
func (r *ConfigRepository) DeleteConfig(ctx context.Context, id int64) error {
	query := `DELETE FROM gateway_configs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("config not found")
	}

	return nil
}

// GetLatestConfig retrieves the latest configuration for a given integration and type
func (r *ConfigRepository) GetLatestConfig(ctx context.Context, integrationID int64, configType string) (*models.GatewayConfig, error) {
	query := `
		SELECT id, integration_id, config_type, config_data, version,
			   status, created_at, updated_at
		FROM gateway_configs 
		WHERE integration_id = $1 AND config_type = $2
		ORDER BY version DESC, created_at DESC
		LIMIT 1`

	var config models.GatewayConfig
	var configData []byte

	err := r.db.QueryRowContext(ctx, query, integrationID, configType).Scan(
		&config.ID, &config.IntegrationID, &config.ConfigType, &configData,
		&config.Version, &config.Status, &config.CreatedAt, &config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest config: %w", err)
	}

	// Unmarshal config data
	if err := json.Unmarshal(configData, &config.ConfigData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return &config, nil
}

// UpdateConfigStatus updates the status of a configuration
func (r *ConfigRepository) UpdateConfigStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE gateway_configs SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update config status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("config not found")
	}

	return nil
} 