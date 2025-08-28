package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"scopeapi.local/backend/services/api-discovery/internal/models"
)

type DiscoveryRepositoryInterface interface {
	CreateDiscovery(ctx context.Context, discovery *models.Discovery) error
	GetDiscovery(ctx context.Context, discoveryID string) (*models.Discovery, error)
	UpdateDiscoveryStatus(ctx context.Context, discoveryID string, status string) error
	UpdateDiscoveryProgress(ctx context.Context, discoveryID string, progress int) error
	UpdateDiscoveryEndpointsFound(ctx context.Context, discoveryID string, count int) error
	GetDiscoveryResults(ctx context.Context, discoveryID string, page, limit int) (*models.DiscoveryResults, error)
	SaveEndpoint(ctx context.Context, endpoint *models.Endpoint) error
	UpdateEndpoint(ctx context.Context, endpoint *models.Endpoint) error
	SaveEndpointAnalysis(ctx context.Context, analysis *models.EndpointAnalysis) error
	GetEndpointMetadata(ctx context.Context, endpointID string) (*models.Metadata, error)
	UpdateEndpointMetadata(ctx context.Context, endpointID string, metadata *models.Metadata) error
	SaveEndpointMetadata(ctx context.Context, metadata *models.Metadata) error
	GetAPISpecification(ctx context.Context, apiID string) (*models.APISpec, error)
	SaveAPISpecification(ctx context.Context, spec *models.APISpec) error
	GetAPIEndpoints(ctx context.Context, apiID string) ([]models.Endpoint, error)
}

type DiscoveryRepository struct {
	db *sqlx.DB
}

func NewDiscoveryRepository(db *sqlx.DB) DiscoveryRepositoryInterface {
	return &DiscoveryRepository{
		db: db,
	}
}

func (r *DiscoveryRepository) CreateDiscovery(ctx context.Context, discovery *models.Discovery) error {
	configJSON, err := json.Marshal(discovery.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		INSERT INTO discoveries (id, target, method, status, progress, start_time, endpoints_found, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = r.db.ExecContext(ctx, query,
		discovery.ID,
		discovery.Target,
		discovery.Method,
		discovery.Status,
		discovery.Progress,
		discovery.StartTime,
		discovery.EndpointsFound,
		configJSON,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create discovery: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) GetDiscovery(ctx context.Context, discoveryID string) (*models.Discovery, error) {
	query := `
		SELECT id, target, method, status, progress, start_time, end_time, endpoints_found, error_message, config, created_at, updated_at
		FROM scopeapi.discoveries
		WHERE id = $1
	`

	var discovery models.Discovery
	var configJSON []byte
	var endTime sql.NullTime

	err := r.db.QueryRowContext(ctx, query, discoveryID).Scan(
		&discovery.ID,
		&discovery.Target,
		&discovery.Method,
		&discovery.Status,
		&discovery.Progress,
		&discovery.StartTime,
		&endTime,
		&discovery.EndpointsFound,
		&discovery.ErrorMessage,
		&configJSON,
		&discovery.CreatedAt,
		&discovery.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("discovery not found: %s", discoveryID)
		}
		return nil, fmt.Errorf("failed to get discovery: %w", err)
	}

	if endTime.Valid {
		discovery.EndTime = &endTime.Time
	}

	if len(configJSON) > 0 {
		var config models.DiscoveryConfig
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
		discovery.Config = &config
	}

	return &discovery, nil
}

func (r *DiscoveryRepository) UpdateDiscoveryStatus(ctx context.Context, discoveryID string, status string) error {
	query := `
		UPDATE scopeapi.discoveries 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), discoveryID)
	if err != nil {
		return fmt.Errorf("failed to update discovery status: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) UpdateDiscoveryProgress(ctx context.Context, discoveryID string, progress int) error {
	query := `
		UPDATE scopeapi.discoveries 
		SET progress = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, progress, time.Now(), discoveryID)
	if err != nil {
		return fmt.Errorf("failed to update discovery progress: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) UpdateDiscoveryEndpointsFound(ctx context.Context, discoveryID string, count int) error {
	query := `
		UPDATE scopeapi.discoveries 
		SET endpoints_found = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, count, time.Now(), discoveryID)
	if err != nil {
		return fmt.Errorf("failed to update discovery endpoints found: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) GetDiscoveryResults(ctx context.Context, discoveryID string, page, limit int) (*models.DiscoveryResults, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := `SELECT COUNT(*) FROM scopeapi.endpoints WHERE discovery_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, discoveryID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get endpoints
	query := `
		SELECT id, api_id, url, path, method, status_code, content_type, summary, description, tags, is_active, created_at, updated_at
		FROM scopeapi.endpoints 
		WHERE discovery_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, discoveryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query endpoints: %w", err)
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var endpoint models.Endpoint
		var tagsJSON []byte

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.APIID,
			&endpoint.URL,
			&endpoint.Path,
			&endpoint.Method,
			&endpoint.StatusCode,
			&endpoint.ContentType,
			&endpoint.Summary,
			&endpoint.Description,
			&tagsJSON,
			&endpoint.IsActive,
			&endpoint.CreatedAt,
			&endpoint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &endpoint.Tags)
		}

		endpoints = append(endpoints, endpoint)
	}

	return &models.DiscoveryResults{
		DiscoveryID: discoveryID,
		Total:       total,
		Page:        page,
		Limit:       limit,
		Endpoints:   endpoints,
	}, nil
}

func (r *DiscoveryRepository) SaveEndpoint(ctx context.Context, endpoint *models.Endpoint) error {
	headersJSON, _ := json.Marshal(endpoint.Headers)
	parametersJSON, _ := json.Marshal(endpoint.Parameters)
	responsesJSON, _ := json.Marshal(endpoint.Responses)
	tagsJSON, _ := json.Marshal(endpoint.Tags)

	query := `
		INSERT INTO scopeapi.endpoints (id, api_id, url, path, method, headers, body, status_code, content_type, summary, description, parameters, responses, tags, is_active, created_at, updated_at)
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
		return fmt.Errorf("failed to save endpoint: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) UpdateEndpoint(ctx context.Context, endpoint *models.Endpoint) error {
	headersJSON, _ := json.Marshal(endpoint.Headers)
	parametersJSON, _ := json.Marshal(endpoint.Parameters)
	responsesJSON, _ := json.Marshal(endpoint.Responses)
	tagsJSON, _ := json.Marshal(endpoint.Tags)

	query := `
		UPDATE scopeapi.endpoints 
		SET headers = $1, body = $2, status_code = $3, content_type = $4, summary = $5, description = $6, parameters = $7, responses = $8, tags = $9, is_active = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.ExecContext(ctx, query,
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

func (r *DiscoveryRepository) SaveEndpointAnalysis(ctx context.Context, analysis *models.EndpointAnalysis) error {
	securityJSON, _ := json.Marshal(analysis.Security)
	headersJSON, _ := json.Marshal(analysis.Headers)
	parametersJSON, _ := json.Marshal(analysis.Parameters)

	query := `
		INSERT INTO scopeapi.endpoint_analyses (endpoint_id, url, method, response_time, status_code, content_type, parameters, headers, security, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		analysis.EndpointID,
		analysis.URL,
		analysis.Method,
		analysis.ResponseTime,
		analysis.StatusCode,
		analysis.ContentType,
		parametersJSON,
		headersJSON,
		securityJSON,
		analysis.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save endpoint analysis: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) GetEndpointMetadata(ctx context.Context, endpointID string) (*models.Metadata, error) {
	query := `
		SELECT id, endpoint_id, api_id, url, method, title, description, tags, category, business_owner, technical_owner, 
		       data_sensitivity, compliance_requirements, parameters, response_schema, request_schema, examples, 
		       documentation, versioning, performance, security, quality, usage, created_at, updated_at
		FROM scopeapi.endpoint_metadata
		WHERE endpoint_id = $1
	`

	var metadata models.Metadata
	var tagsJSON, complianceJSON, parametersJSON, responseSchemaJSON, requestSchemaJSON, examplesJSON []byte
	var documentationJSON, versioningJSON, performanceJSON, securityJSON, qualityJSON, usageJSON []byte

	err := r.db.QueryRowContext(ctx, query, endpointID).Scan(
		&metadata.ID,
		&metadata.EndpointID,
		&metadata.APIID,
		&metadata.URL,
		&metadata.Method,
		&metadata.Title,
		&metadata.Description,
		&tagsJSON,
		&metadata.Category,
		&metadata.BusinessOwner,
		&metadata.TechnicalOwner,
		&metadata.DataSensitivity,
		&complianceJSON,
		&parametersJSON,
		&responseSchemaJSON,
		&requestSchemaJSON,
		&examplesJSON,
		&documentationJSON,
		&versioningJSON,
		&performanceJSON,
		&securityJSON,
		&qualityJSON,
		&usageJSON,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("metadata not found for endpoint: %s", endpointID)
		}
		return nil, fmt.Errorf("failed to get endpoint metadata: %w", err)
	}

	// Unmarshal JSON fields
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &metadata.Tags)
	}
	if len(complianceJSON) > 0 {
		json.Unmarshal(complianceJSON, &metadata.ComplianceReqs)
	}
	if len(parametersJSON) > 0 {
		json.Unmarshal(parametersJSON, &metadata.Parameters)
	}
	if len(responseSchemaJSON) > 0 {
		json.Unmarshal(responseSchemaJSON, &metadata.ResponseSchema)
	}
	if len(requestSchemaJSON) > 0 {
		json.Unmarshal(requestSchemaJSON, &metadata.RequestSchema)
	}
	if len(examplesJSON) > 0 {
		json.Unmarshal(examplesJSON, &metadata.Examples)
	}
	if len(documentationJSON) > 0 {
		json.Unmarshal(documentationJSON, &metadata.Documentation)
	}
	if len(versioningJSON) > 0 {
		json.Unmarshal(versioningJSON, &metadata.Versioning)
	}
	if len(performanceJSON) > 0 {
		json.Unmarshal(performanceJSON, &metadata.Performance)
	}
	if len(securityJSON) > 0 {
		json.Unmarshal(securityJSON, &metadata.Security)
	}
	if len(qualityJSON) > 0 {
		json.Unmarshal(qualityJSON, &metadata.Quality)
	}
	if len(usageJSON) > 0 {
		json.Unmarshal(usageJSON, &metadata.Usage)
	}

	return &metadata, nil
}

func (r *DiscoveryRepository) UpdateEndpointMetadata(ctx context.Context, endpointID string, metadata *models.Metadata) error {
	tagsJSON, _ := json.Marshal(metadata.Tags)
	complianceJSON, _ := json.Marshal(metadata.ComplianceReqs)
	parametersJSON, _ := json.Marshal(metadata.Parameters)
	responseSchemaJSON, _ := json.Marshal(metadata.ResponseSchema)
	requestSchemaJSON, _ := json.Marshal(metadata.RequestSchema)
	examplesJSON, _ := json.Marshal(metadata.Examples)
	documentationJSON, _ := json.Marshal(metadata.Documentation)
	versioningJSON, _ := json.Marshal(metadata.Versioning)
	performanceJSON, _ := json.Marshal(metadata.Performance)
	securityJSON, _ := json.Marshal(metadata.Security)
	qualityJSON, _ := json.Marshal(metadata.Quality)
	usageJSON, _ := json.Marshal(metadata.Usage)

	query := `
		UPDATE endpoint_metadata 
		SET title = $1, description = $2, tags = $3, category = $4, business_owner = $5, technical_owner = $6,
		    data_sensitivity = $7, compliance_requirements = $8, parameters = $9, response_schema = $10,
		    request_schema = $11, examples = $12, documentation = $13, versioning = $14, performance = $15,
		    security = $16, quality = $17, usage = $18, updated_at = $19
		WHERE endpoint_id = $20
	`

	_, err := r.db.ExecContext(ctx, query,
		metadata.Title,
		metadata.Description,
		tagsJSON,
		metadata.Category,
		metadata.BusinessOwner,
		metadata.TechnicalOwner,
		metadata.DataSensitivity,
		complianceJSON,
		parametersJSON,
		responseSchemaJSON,
		requestSchemaJSON,
		examplesJSON,
		documentationJSON,
		versioningJSON,
		performanceJSON,
		securityJSON,
		qualityJSON,
		usageJSON,
		time.Now(),
		endpointID,
	)

	if err != nil {
		return fmt.Errorf("failed to update endpoint metadata: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) SaveEndpointMetadata(ctx context.Context, metadata *models.Metadata) error {
	tagsJSON, _ := json.Marshal(metadata.Tags)
	complianceJSON, _ := json.Marshal(metadata.ComplianceReqs)
	parametersJSON, _ := json.Marshal(metadata.Parameters)
	responseSchemaJSON, _ := json.Marshal(metadata.ResponseSchema)
	requestSchemaJSON, _ := json.Marshal(metadata.RequestSchema)
	examplesJSON, _ := json.Marshal(metadata.Examples)
	documentationJSON, _ := json.Marshal(metadata.Documentation)
	versioningJSON, _ := json.Marshal(metadata.Versioning)
	performanceJSON, _ := json.Marshal(metadata.Performance)
	securityJSON, _ := json.Marshal(metadata.Security)
	qualityJSON, _ := json.Marshal(metadata.Quality)
	usageJSON, _ := json.Marshal(metadata.Usage)

	query := `
		INSERT INTO endpoint_metadata (id, endpoint_id, api_id, url, method, title, description, tags, category, 
		                              business_owner, technical_owner, data_sensitivity, compliance_requirements, 
		                              parameters, response_schema, request_schema, examples, documentation, 
		                              versioning, performance, security, quality, usage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`

	_, err := r.db.ExecContext(ctx, query,
		metadata.ID,
		metadata.EndpointID,
		metadata.APIID,
		metadata.URL,
		metadata.Method,
		metadata.Title,
		metadata.Description,
		tagsJSON,
		metadata.Category,
		metadata.BusinessOwner,
		metadata.TechnicalOwner,
		metadata.DataSensitivity,
		complianceJSON,
		parametersJSON,
		responseSchemaJSON,
		requestSchemaJSON,
		examplesJSON,
		documentationJSON,
		versioningJSON,
		performanceJSON,
		securityJSON,
		qualityJSON,
		usageJSON,
		metadata.CreatedAt,
		metadata.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save endpoint metadata: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) GetAPISpecification(ctx context.Context, apiID string) (*models.APISpec, error) {
	query := `
		SELECT id, api_id, version, title, description, openapi_version, info, servers, paths, components, security, tags, external_docs, created_at, updated_at
		FROM scopeapi.api_specifications
		WHERE api_id = $1
	`

	var spec models.APISpec
	var infoJSON, serversJSON, pathsJSON, componentsJSON, securityJSON, tagsJSON, externalDocsJSON []byte

	err := r.db.QueryRowContext(ctx, query, apiID).Scan(
		&spec.ID,
		&spec.APIID,
		&spec.Version,
		&spec.Title,
		&spec.Description,
		&spec.OpenAPIVersion,
		&infoJSON,
		&serversJSON,
		&pathsJSON,
		&componentsJSON,
		&securityJSON,
		&tagsJSON,
		&externalDocsJSON,
		&spec.CreatedAt,
		&spec.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API specification not found for API: %s", apiID)
		}
		return nil, fmt.Errorf("failed to get API specification: %w", err)
	}

	// Unmarshal JSON fields
	if len(infoJSON) > 0 {
		json.Unmarshal(infoJSON, &spec.Info)
	}
	if len(serversJSON) > 0 {
		json.Unmarshal(serversJSON, &spec.Servers)
	}
	if len(pathsJSON) > 0 {
		json.Unmarshal(pathsJSON, &spec.Paths)
	}
	if len(componentsJSON) > 0 {
		json.Unmarshal(componentsJSON, &spec.Components)
	}
	if len(securityJSON) > 0 {
		json.Unmarshal(securityJSON, &spec.Security)
	}
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &spec.Tags)
	}
	if len(externalDocsJSON) > 0 {
		json.Unmarshal(externalDocsJSON, &spec.ExternalDocs)
	}

	return &spec, nil
}

func (r *DiscoveryRepository) SaveAPISpecification(ctx context.Context, spec *models.APISpec) error {
	infoJSON, _ := json.Marshal(spec.Info)
	serversJSON, _ := json.Marshal(spec.Servers)
	pathsJSON, _ := json.Marshal(spec.Paths)
	componentsJSON, _ := json.Marshal(spec.Components)
	securityJSON, _ := json.Marshal(spec.Security)
	tagsJSON, _ := json.Marshal(spec.Tags)
	externalDocsJSON, _ := json.Marshal(spec.ExternalDocs)

	query := `
		INSERT INTO scopeapi.api_specifications (id, api_id, version, title, description, openapi_version, info, servers, paths, components, security, tags, external_docs, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (api_id) DO UPDATE SET
		    version = EXCLUDED.version,
		    title = EXCLUDED.title,
		    description = EXCLUDED.description,
		    openapi_version = EXCLUDED.openapi_version,
		    info = EXCLUDED.info,
		    servers = EXCLUDED.servers,
		    paths = EXCLUDED.paths,
		    components = EXCLUDED.components,
		    security = EXCLUDED.security,
		    tags = EXCLUDED.tags,
		    external_docs = EXCLUDED.external_docs,
		    updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		spec.ID,
		spec.APIID,
		spec.Version,
		spec.Title,
		spec.Description,
		spec.OpenAPIVersion,
		infoJSON,
		serversJSON,
		pathsJSON,
		componentsJSON,
		securityJSON,
		tagsJSON,
		externalDocsJSON,
		spec.CreatedAt,
		spec.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save API specification: %w", err)
	}

	return nil
}

func (r *DiscoveryRepository) GetAPIEndpoints(ctx context.Context, apiID string) ([]models.Endpoint, error) {
	query := `
		SELECT id, api_id, url, path, method, headers, body, status_code, content_type, summary, description, parameters, responses, tags, is_active, created_at, updated_at
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
