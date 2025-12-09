package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"scopeapi.local/backend/services/gateway-integration/internal/models"
)

// IntegrationRepository handles database operations for integrations
type IntegrationRepository struct {
	db *sqlx.DB
}

// NewIntegrationRepository creates a new integration repository
func NewIntegrationRepository(db *sqlx.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

// CreateIntegration creates a new integration
func (r *IntegrationRepository) CreateIntegration(ctx context.Context, integration *models.Integration) error {
	query := `
		INSERT INTO integrations (integration_id, name, type, status, config, credentials, endpoints, health)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING integration_id, created_at, updated_at
	`

	// Generate UUID if not provided
	if integration.ID == "" {
		integration.ID = uuid.New().String()
	}

	// Convert credentials to JSONB if present
	var credentialsJSON []byte
	if integration.Credentials != nil {
		var err error
		credentialsJSON, err = json.Marshal(integration.Credentials)
		if err != nil {
			return fmt.Errorf("failed to marshal credentials: %w", err)
		}
	}

	// Convert endpoints to JSONB
	endpointsJSON, err := json.Marshal(integration.Endpoints)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoints: %w", err)
	}

	// Convert config to JSONB
	configJSON, err := json.Marshal(integration.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Convert health to JSONB if present
	var healthJSON []byte
	if integration.Health != nil {
		healthJSON, err = json.Marshal(integration.Health)
		if err != nil {
			return fmt.Errorf("failed to marshal health: %w", err)
		}
	}

	var createdAt, updatedAt time.Time
	err = r.db.QueryRowContext(ctx, query,
		integration.ID,
		integration.Name,
		integration.Type,
		integration.Status,
		configJSON,
		credentialsJSON,
		endpointsJSON,
		healthJSON,
	).Scan(&integration.ID, &createdAt, &updatedAt)

	if err != nil {
		return fmt.Errorf("failed to create integration: %w", err)
	}

	integration.CreatedAt = createdAt
	integration.UpdatedAt = updatedAt

	return nil
}

// GetIntegration retrieves an integration by ID
func (r *IntegrationRepository) GetIntegration(ctx context.Context, id string) (*models.Integration, error) {
	query := `
		SELECT integration_id, name, type, status, config, credentials, endpoints, health, 
		       created_at, updated_at, last_sync
		FROM integrations
		WHERE integration_id = $1
	`

	var integration models.Integration
	var configJSON, credentialsJSON, endpointsJSON, healthJSON []byte
	var lastSync sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&integration.ID,
		&integration.Name,
		&integration.Type,
		&integration.Status,
		&configJSON,
		&credentialsJSON,
		&endpointsJSON,
		&healthJSON,
		&integration.CreatedAt,
		&integration.UpdatedAt,
		&lastSync,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("integration not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	// Parse JSONB fields
	if err := json.Unmarshal(configJSON, &integration.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if credentialsJSON != nil {
		if err := json.Unmarshal(credentialsJSON, &integration.Credentials); err != nil {
			return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
		}
	}

	if err := json.Unmarshal(endpointsJSON, &integration.Endpoints); err != nil {
		return nil, fmt.Errorf("failed to unmarshal endpoints: %w", err)
	}

	if healthJSON != nil {
		if err := json.Unmarshal(healthJSON, &integration.Health); err != nil {
			return nil, fmt.Errorf("failed to unmarshal health: %w", err)
		}
	}

	if lastSync.Valid {
		integration.LastSync = &lastSync.Time
	}

	return &integration, nil
}

// GetIntegrations retrieves all integrations with optional filters
func (r *IntegrationRepository) GetIntegrations(ctx context.Context, filters map[string]interface{}) ([]*models.Integration, error) {
	query := `
		SELECT integration_id, name, type, status, config, credentials, endpoints, health, 
		       created_at, updated_at, last_sync
		FROM integrations
	`
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if len(filters) > 0 {
		query += " WHERE "
		conditions := []string{}

		if gatewayType, exists := filters["type"]; exists {
			conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
			args = append(args, gatewayType)
			argIndex++
		}

		if status, exists := filters["status"]; exists {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, status)
			argIndex++
		}

		if name, exists := filters["name"]; exists {
			conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
			args = append(args, "%"+name.(string)+"%")
			argIndex++
		}

		query += fmt.Sprintf("(%s)", conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += fmt.Sprintf(" AND (%s)", conditions[i])
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query integrations: %w", err)
	}
	defer rows.Close()

	var integrations []*models.Integration
	for rows.Next() {
		var integration models.Integration
		var configJSON, credentialsJSON, endpointsJSON, healthJSON []byte
		var lastSync sql.NullTime

		err := rows.Scan(
			&integration.ID,
			&integration.Name,
			&integration.Type,
			&integration.Status,
			&configJSON,
			&credentialsJSON,
			&endpointsJSON,
			&healthJSON,
			&integration.CreatedAt,
			&integration.UpdatedAt,
			&lastSync,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan integration: %w", err)
		}

		// Parse JSONB fields
		if err := json.Unmarshal(configJSON, &integration.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if credentialsJSON != nil {
			if err := json.Unmarshal(credentialsJSON, &integration.Credentials); err != nil {
				return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
			}
		}

		if err := json.Unmarshal(endpointsJSON, &integration.Endpoints); err != nil {
			return nil, fmt.Errorf("failed to unmarshal endpoints: %w", err)
		}

		if healthJSON != nil {
			if err := json.Unmarshal(healthJSON, &integration.Health); err != nil {
				return nil, fmt.Errorf("failed to unmarshal health: %w", err)
			}
		}

		if lastSync.Valid {
			integration.LastSync = &lastSync.Time
		}

		integrations = append(integrations, &integration)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating integrations: %w", err)
	}

	return integrations, nil
}

// UpdateIntegration updates an existing integration
func (r *IntegrationRepository) UpdateIntegration(ctx context.Context, integration *models.Integration) error {
	query := `
		UPDATE integrations
		SET name = $2, type = $3, status = $4, config = $5, credentials = $6, 
		    endpoints = $7, health = $8, updated_at = NOW()
		WHERE integration_id = $1
		RETURNING updated_at
	`

	// Convert credentials to JSONB if present
	var credentialsJSON []byte
	if integration.Credentials != nil {
		var err error
		credentialsJSON, err = json.Marshal(integration.Credentials)
		if err != nil {
			return fmt.Errorf("failed to marshal credentials: %w", err)
		}
	}

	// Convert endpoints to JSONB
	endpointsJSON, err := json.Marshal(integration.Endpoints)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoints: %w", err)
	}

	// Convert config to JSONB
	configJSON, err := json.Marshal(integration.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Convert health to JSONB if present
	var healthJSON []byte
	if integration.Health != nil {
		healthJSON, err = json.Marshal(integration.Health)
		if err != nil {
			return fmt.Errorf("failed to marshal health: %w", err)
		}
	}

	var updatedAt time.Time
	err = r.db.QueryRowContext(ctx, query,
		integration.ID,
		integration.Name,
		integration.Type,
		integration.Status,
		configJSON,
		credentialsJSON,
		endpointsJSON,
		healthJSON,
	).Scan(&updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("integration not found: %s", integration.ID)
		}
		return fmt.Errorf("failed to update integration: %w", err)
	}

	integration.UpdatedAt = updatedAt
	return nil
}

// DeleteIntegration deletes an integration
func (r *IntegrationRepository) DeleteIntegration(ctx context.Context, id string) error {
	query := `DELETE FROM integrations WHERE integration_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete integration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("integration not found: %s", id)
	}

	return nil
}

// UpdateIntegrationHealth updates the health status of an integration
func (r *IntegrationRepository) UpdateIntegrationHealth(ctx context.Context, id string, health *models.HealthStatus) error {
	query := `
		UPDATE integrations
		SET health = $2, updated_at = NOW()
		WHERE integration_id = $1
	`

	healthJSON, err := json.Marshal(health)
	if err != nil {
		return fmt.Errorf("failed to marshal health: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, id, healthJSON)
	if err != nil {
		return fmt.Errorf("failed to update integration health: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("integration not found: %s", id)
	}

	return nil
}

// UpdateIntegrationLastSync updates the last sync time of an integration
func (r *IntegrationRepository) UpdateIntegrationLastSync(ctx context.Context, id string, lastSync time.Time) error {
	query := `
		UPDATE integrations
		SET last_sync = $2, updated_at = NOW()
		WHERE integration_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, lastSync)
	if err != nil {
		return fmt.Errorf("failed to update integration last sync: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("integration not found: %s", id)
	}

	return nil
}

// GetIntegrationStats retrieves statistics about integrations
func (r *IntegrationRepository) GetIntegrationStats(ctx context.Context) (*models.IntegrationStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
			COUNT(CASE WHEN status = 'error' THEN 1 END) as error,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN type = 'kong' THEN 1 END) as kong_count,
			COUNT(CASE WHEN type = 'nginx' THEN 1 END) as nginx_count,
			COUNT(CASE WHEN type = 'traefik' THEN 1 END) as traefik_count,
			COUNT(CASE WHEN type = 'envoy' THEN 1 END) as envoy_count,
			COUNT(CASE WHEN type = 'haproxy' THEN 1 END) as haproxy_count
		FROM integrations
	`

	var stats models.IntegrationStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.Active,
		&stats.Error,
		&stats.Pending,
		&stats.KongCount,
		&stats.NginxCount,
		&stats.TraefikCount,
		&stats.EnvoyCount,
		&stats.HAProxyCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get integration stats: %w", err)
	}

	return &stats, nil
} 