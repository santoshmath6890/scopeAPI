package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"strings"

	"github.com/jmoiron/sqlx"
	"scopeapi.local/backend/services/api-discovery/internal/models"
)

type InventoryRepositoryInterface interface {
	GetAPIs(ctx context.Context, page, limit int) (*models.APIInventory, error)
	GetAPI(ctx context.Context, apiID string) (*models.API, error)
	CreateAPI(ctx context.Context, api *models.API) error
	UpdateAPI(ctx context.Context, api *models.API) error
	DeleteAPI(ctx context.Context, apiID string) error
	GetAPIDetails(ctx context.Context, apiID string) (*models.APIDetails, error)
	GetAPIStatistics(ctx context.Context) (*models.APIStatistics, error)
	SearchAPIs(ctx context.Context, query string, page, limit int) (*models.APIInventory, error)
	GetAPIsByTags(ctx context.Context, tags []string, page, limit int) (*models.APIInventory, error)
	GetAPIsByStatus(ctx context.Context, status string, page, limit int) (*models.APIInventory, error)
	CreateEndpoint(ctx context.Context, endpoint *models.Endpoint) error
	UpdateEndpoint(ctx context.Context, endpoint *models.Endpoint) error
	GetEndpoint(ctx context.Context, endpointID string) (*models.Endpoint, error)
	GetEndpoints(ctx context.Context, page, limit int) ([]models.Endpoint, error)
	DeleteEndpoint(ctx context.Context, endpointID string) error
}

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) InventoryRepositoryInterface {
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) GetAPIs(ctx context.Context, page, limit int) (*models.APIInventory, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := `SELECT COUNT(*) FROM scopeapi.apis`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get APIs
	query := `
		SELECT id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at
		FROM scopeapi.apis
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query APIs: %w", err)
	}
	defer rows.Close()

	var apis []models.API
	for rows.Next() {
		var api models.API
		var tagsJSON []byte

		err := rows.Scan(
			&api.ID,
			&api.Name,
			&api.URL,
			&api.BaseURL,
			&api.Version,
			&api.Protocol,
			&api.Status,
			&api.Description,
			&tagsJSON,
			&api.CreatedAt,
			&api.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &api.Tags)
		}

		apis = append(apis, api)
	}

	return &models.APIInventory{
		Total: total,
		Page:  page,
		Limit: limit,
		APIs:  apis,
	}, nil
}

func (r *InventoryRepository) GetAPI(ctx context.Context, apiID string) (*models.API, error) {
	query := `
		SELECT id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at
		FROM scopeapi.apis
		WHERE id = $1
	`

	var api models.API
	var tagsJSON []byte

	err := r.db.QueryRowContext(ctx, query, apiID).Scan(
		&api.ID,
		&api.Name,
		&api.URL,
		&api.BaseURL,
		&api.Version,
		&api.Protocol,
		&api.Status,
		&api.Description,
		&tagsJSON,
		&api.CreatedAt,
		&api.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API not found: %s", apiID)
		}
		return nil, fmt.Errorf("failed to get API: %w", err)
	}

	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &api.Tags)
	}

	return &api, nil
}

func (r *InventoryRepository) CreateAPI(ctx context.Context, api *models.API) error {
	tagsJSON, err := json.Marshal(api.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO scopeapi.apis (id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.db.ExecContext(ctx, query,
		api.ID,
		api.Name,
		api.URL,
		api.BaseURL,
		api.Version,
		api.Protocol,
		api.Status,
		api.Description,
		tagsJSON,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create API: %w", err)
	}

	return nil
}

func (r *InventoryRepository) UpdateAPI(ctx context.Context, api *models.API) error {
	tagsJSON, err := json.Marshal(api.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE scopeapi.apis 
		SET name = $1, url = $2, base_url = $3, version = $4, protocol = $5, status = $6, description = $7, tags = $8, updated_at = $9
		WHERE id = $10
	`

	_, err = r.db.ExecContext(ctx, query,
		api.Name,
		api.URL,
		api.BaseURL,
		api.Version,
		api.Protocol,
		api.Status,
		api.Description,
		tagsJSON,
		time.Now(),
		api.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update API: %w", err)
	}

	return nil
}

func (r *InventoryRepository) DeleteAPI(ctx context.Context, apiID string) error {
	query := `DELETE FROM scopeapi.apis WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, apiID)
	if err != nil {
		return fmt.Errorf("failed to delete API: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetAPIDetails(ctx context.Context, apiID string) (*models.APIDetails, error) {
	// Get API
	api, err := r.GetAPI(ctx, apiID)
	if err != nil {
		return nil, err
	}

	// Get endpoints
	endpoints, err := r.getAPIEndpoints(ctx, apiID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API endpoints: %w", err)
	}

	// Calculate statistics
	stats := models.APIStats{
		TotalEndpoints: len(endpoints),
	}

	for _, endpoint := range endpoints {
		if endpoint.IsActive {
			stats.ActiveEndpoints++
		} else {
			stats.InactiveEndpoints++
		}
	}

	return &models.APIDetails{
		API:        *api,
		Endpoints:  endpoints,
		Statistics: stats,
	}, nil
}

func (r *InventoryRepository) GetAPIStatistics(ctx context.Context) (*models.APIStatistics, error) {
	stats := &models.APIStatistics{
		ProtocolBreakdown: make(map[string]int),
		StatusBreakdown:   make(map[string]int),
	}

	// Get total APIs
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM scopeapi.apis").Scan(&stats.TotalAPIs)
	if err != nil {
		return nil, fmt.Errorf("failed to get total APIs: %w", err)
	}

	// Get active/inactive APIs
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM scopeapi.apis WHERE status = 'active'").Scan(&stats.ActiveAPIs)
	if err != nil {
		return nil, fmt.Errorf("failed to get active APIs: %w", err)
	}

	stats.InactiveAPIs = stats.TotalAPIs - stats.ActiveAPIs

	// Get total endpoints
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM scopeapi.endpoints").Scan(&stats.TotalEndpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get total endpoints: %w", err)
	}

	// Get recent discoveries (last 7 days)
	err = r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM scopeapi.discoveries WHERE created_at > $1", 
		time.Now().AddDate(0, 0, -7)).Scan(&stats.RecentDiscoveries)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent discoveries: %w", err)
	}

	// Get protocol breakdown
	protocolRows, err := r.db.QueryContext(ctx, "SELECT protocol, COUNT(*) FROM scopeapi.apis GROUP BY protocol")
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol breakdown: %w", err)
	}
	defer protocolRows.Close()

	for protocolRows.Next() {
		var protocol string
		var count int
		err := protocolRows.Scan(&protocol, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan protocol breakdown: %w", err)
		}
		stats.ProtocolBreakdown[protocol] = count
	}

	// Get status breakdown
	statusRows, err := r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM scopeapi.apis GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("failed to get status breakdown: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int
		err := statusRows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status breakdown: %w", err)
		}
		stats.StatusBreakdown[status] = count
	}

	return stats, nil
}

func (r *InventoryRepository) SearchAPIs(ctx context.Context, query string, page, limit int) (*models.APIInventory, error) {
	offset := (page - 1) * limit
	searchPattern := "%" + query + "%"

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM scopeapi.apis 
		WHERE name ILIKE $1 OR description ILIKE $1 OR url ILIKE $1
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get search count: %w", err)
	}

	// Get APIs
	searchQuery := `
		SELECT id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at
		FROM scopeapi.apis
		WHERE name ILIKE $1 OR description ILIKE $1 OR url ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search APIs: %w", err)
	}
	defer rows.Close()

	var apis []models.API
	for rows.Next() {
		var api models.API
		var tagsJSON []byte

		err := rows.Scan(
			&api.ID,
			&api.Name,
			&api.URL,
			&api.BaseURL,
			&api.Version,
			&api.Protocol,
			&api.Status,
			&api.Description,
			&tagsJSON,
			&api.CreatedAt,
			&api.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &api.Tags)
		}

		apis = append(apis, api)
	}

	return &models.APIInventory{
		Total: total,
		Page:  page,
		Limit: limit,
		APIs:  apis,
	}, nil
}

func (r *InventoryRepository) GetAPIsByTags(ctx context.Context, tags []string, page, limit int) (*models.APIInventory, error) {
	offset := (page - 1) * limit

	// Build query for tag matching
	tagConditions := make([]string, len(tags))
	args := make([]interface{}, len(tags)+2)
	
	for i, tag := range tags {
		tagConditions[i] = fmt.Sprintf("tags::text ILIKE $%d", i+1)
		args[i] = "%" + tag + "%"
	}
	args[len(tags)] = limit
	args[len(tags)+1] = offset

	whereClause := strings.Join(tagConditions, " OR ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM scopeapi.apis WHERE %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(tags)]...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag count: %w", err)
	}

	// Get APIs
	query := fmt.Sprintf(`
		SELECT id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at
		FROM scopeapi.apis
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(tags)+1, len(tags)+2)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query APIs by tags: %w", err)
	}
	defer rows.Close()

	var apis []models.API
	for rows.Next() {
		var api models.API
		var tagsJSON []byte

		err := rows.Scan(
			&api.ID,
			&api.Name,
			&api.URL,
			&api.BaseURL,
			&api.Version,
			&api.Protocol,
			&api.Status,
			&api.Description,
			&tagsJSON,
			&api.CreatedAt,
			&api.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &api.Tags)
		}

		apis = append(apis, api)
	}

	return &models.APIInventory{
		Total: total,
		Page:  page,
		Limit: limit,
		APIs:  apis,
	}, nil
}

func (r *InventoryRepository) GetAPIsByStatus(ctx context.Context, status string, page, limit int) (*models.APIInventory, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := `SELECT COUNT(*) FROM scopeapi.apis WHERE status = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get status count: %w", err)
	}

	// Get APIs
	query := `
		SELECT id, name, url, base_url, version, protocol, status, description, tags, created_at, updated_at
		FROM scopeapi.apis
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query APIs by status: %w", err)
	}
	defer rows.Close()

	var apis []models.API
	for rows.Next() {
		var api models.API
		var tagsJSON []byte

		err := rows.Scan(
			&api.ID,
			&api.Name,
			&api.URL,
			&api.BaseURL,
			&api.Version,
			&api.Protocol,
			&api.Status,
			&api.Description,
			&tagsJSON,
			&api.CreatedAt,
			&api.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &api.Tags)
		}

		apis = append(apis, api)
	}

	return &models.APIInventory{
		Total: total,
		Page:  page,
		Limit: limit,
		APIs:  apis,
	}, nil
}

func (r *InventoryRepository) CreateEndpoint(ctx context.Context, endpoint *models.Endpoint) error {
	headersJSON, _ := json.Marshal(endpoint.Headers)
	parametersJSON, _ := json.Marshal(endpoint.Parameters)
	responsesJSON, _ := json.Marshal(endpoint.Responses)
	tagsJSON, _ := json.Marshal(endpoint.Tags)

	query := `
		INSERT INTO scopeapi.endpoints (id, api_id, url, path, method, headers, body, status_code, content_type, 
		                      summary, description, parameters, responses, tags, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.ExecContext(ctx, query,
		endpoint.ID,
		endpoint.APIID,
		endpoint.URL,
		endpoint.Path,
		endpoint.Method,
		headersJSON,
		endpoint.Body,
		endpoint.StatusCode,
		endpoint.ContentType,
		endpoint.Summary,
		endpoint.Description,
		parametersJSON,
		responsesJSON,
		tagsJSON,
		endpoint.IsActive,
		endpoint.CreatedAt,
		endpoint.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	return nil
}

func (r *InventoryRepository) UpdateEndpoint(ctx context.Context, endpoint *models.Endpoint) error {
	headersJSON, _ := json.Marshal(endpoint.Headers)
	parametersJSON, _ := json.Marshal(endpoint.Parameters)
	responsesJSON, _ := json.Marshal(endpoint.Responses)
	tagsJSON, _ := json.Marshal(endpoint.Tags)

	query := `
		UPDATE scopeapi.endpoints 
		SET api_id = $1, url = $2, path = $3, method = $4, headers = $5, body = $6, status_code = $7,
		    content_type = $8, summary = $9, description = $10, parameters = $11, responses = $12,
		    tags = $13, is_active = $14, updated_at = $15
		WHERE id = $16
	`

	_, err := r.db.ExecContext(ctx, query,
		endpoint.APIID,
		endpoint.URL,
		endpoint.Path,
		endpoint.Method,
		headersJSON,
		endpoint.Body,
		endpoint.StatusCode,
		endpoint.ContentType,
		endpoint.Summary,
		endpoint.Description,
		parametersJSON,
		responsesJSON,
		tagsJSON,
		endpoint.IsActive,
		time.Now(),
		endpoint.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update endpoint: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetEndpoint(ctx context.Context, endpointID string) (*models.Endpoint, error) {
	query := `
		SELECT id, api_id, url, path, method, headers, body, status_code, content_type, 
		       summary, description, parameters, responses, tags, is_active, created_at, updated_at
		FROM scopeapi.endpoints
		WHERE id = $1
	`

	var endpoint models.Endpoint
	var headersJSON, parametersJSON, responsesJSON, tagsJSON []byte

	err := r.db.QueryRowContext(ctx, query, endpointID).Scan(
		&endpoint.ID,
		&endpoint.APIID,
		&endpoint.URL,
		&endpoint.Path,
		&endpoint.Method,
		&headersJSON,
		&endpoint.Body,
		&endpoint.StatusCode,
		&endpoint.ContentType,
		&endpoint.Summary,
		&endpoint.Description,
		&parametersJSON,
		&responsesJSON,
		&tagsJSON,
		&endpoint.IsActive,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("endpoint not found: %s", endpointID)
		}
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Unmarshal JSON fields
	if len(headersJSON) > 0 {
		json.Unmarshal(headersJSON, &endpoint.Headers)
	}
	if len(parametersJSON) > 0 {
		json.Unmarshal(parametersJSON, &endpoint.Parameters)
	}
	if len(responsesJSON) > 0 {
		json.Unmarshal(responsesJSON, &endpoint.Responses)
	}
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &endpoint.Tags)
	}

	return &endpoint, nil
}

func (r *InventoryRepository) GetEndpoints(ctx context.Context, page, limit int) ([]models.Endpoint, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, api_id, url, path, method, headers, body, status_code, content_type, 
		       summary, description, parameters, responses, tags, is_active, created_at, updated_at
		FROM scopeapi.endpoints
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query endpoints: %w", err)
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var endpoint models.Endpoint
		var headersJSON, parametersJSON, responsesJSON, tagsJSON []byte

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.APIID,
			&endpoint.URL,
			&endpoint.Path,
			&endpoint.Method,
			&headersJSON,
			&endpoint.Body,
			&endpoint.StatusCode,
			&endpoint.ContentType,
			&endpoint.Summary,
			&endpoint.Description,
			&parametersJSON,
			&responsesJSON,
			&tagsJSON,
			&endpoint.IsActive,
			&endpoint.CreatedAt,
			&endpoint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}

		// Unmarshal JSON fields
		if len(headersJSON) > 0 {
			json.Unmarshal(headersJSON, &endpoint.Headers)
		}
		if len(parametersJSON) > 0 {
			json.Unmarshal(parametersJSON, &endpoint.Parameters)
		}
		if len(responsesJSON) > 0 {
			json.Unmarshal(responsesJSON, &endpoint.Responses)
		}
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &endpoint.Tags)
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func (r *InventoryRepository) DeleteEndpoint(ctx context.Context, endpointID string) error {
	query := `DELETE FROM scopeapi.endpoints WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, endpointID)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %w", err)
	}

	return nil
}

// Helper method to get API endpoints
func (r *InventoryRepository) getAPIEndpoints(ctx context.Context, apiID string) ([]models.Endpoint, error) {
	query := `
		SELECT id, api_id, url, path, method, headers, body, status_code, content_type, 
		       summary, description, parameters, responses, tags, is_active, created_at, updated_at
		FROM scopeapi.endpoints
		WHERE api_id = $1
		ORDER BY path, method
	`

	rows, err := r.db.QueryContext(ctx, query, apiID)
	if err != nil {
		return nil, fmt.Errorf("failed to query API endpoints: %w", err)
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var endpoint models.Endpoint
		var headersJSON, parametersJSON, responsesJSON, tagsJSON []byte

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.APIID,
			&endpoint.URL,
			&endpoint.Path,
			&endpoint.Method,
			&headersJSON,
			&endpoint.Body,
			&endpoint.StatusCode,
			&endpoint.ContentType,
			&endpoint.Summary,
			&endpoint.Description,
			&parametersJSON,
			&responsesJSON,
			&tagsJSON,
			&endpoint.IsActive,
			&endpoint.CreatedAt,
			&endpoint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}

		// Unmarshal JSON fields
		if len(headersJSON) > 0 {
			json.Unmarshal(headersJSON, &endpoint.Headers)
		}
		if len(parametersJSON) > 0 {
			json.Unmarshal(parametersJSON, &endpoint.Parameters)
		}
		if len(responsesJSON) > 0 {
			json.Unmarshal(responsesJSON, &endpoint.Responses)
		}
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &endpoint.Tags)
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}
