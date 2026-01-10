package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/api-discovery/internal/models"
	"scopeapi.local/backend/services/api-discovery/internal/repository"
	"scopeapi.local/backend/shared/logging"
)

type MetadataServiceInterface interface {
	GetEndpointMetadata(ctx context.Context, endpointID string) (*models.Metadata, error)
	UpdateEndpointMetadata(ctx context.Context, endpointID string, metadata *models.Metadata) error
	ExtractMetadata(ctx context.Context, endpoint *models.Endpoint) (*models.Metadata, error)
	GetAPISpecification(ctx context.Context, apiID string) (*models.APISpec, error)
	GenerateAPISpec(ctx context.Context, apiID string) (*models.APISpec, error)
	AnalyzeEndpointSecurity(ctx context.Context, endpoint *models.Endpoint) (*models.SecurityMetadata, error)
	CalculateQualityMetrics(ctx context.Context, endpointID string) (*models.QualityMetrics, error)
	UpdateUsageMetrics(ctx context.Context, endpointID string, requestCount int64) error
	EnrichMetadata(ctx context.Context, metadata *models.Metadata) error
}

type MetadataService struct {
	repo   repository.DiscoveryRepositoryInterface
	logger logging.Logger
}

func NewMetadataService(repo repository.DiscoveryRepositoryInterface, logger logging.Logger) MetadataServiceInterface {
	return &MetadataService{
		repo:   repo,
		logger: logger,
	}
}

func (s *MetadataService) GetEndpointMetadata(ctx context.Context, endpointID string) (*models.Metadata, error) {
	metadata, err := s.repo.GetEndpointMetadata(ctx, endpointID)
	if err != nil {
		s.logger.Error("Failed to get endpoint metadata", "error", err, "endpoint_id", endpointID)
		return nil, fmt.Errorf("failed to get endpoint metadata: %w", err)
	}

	return metadata, nil
}

func (s *MetadataService) UpdateEndpointMetadata(ctx context.Context, endpointID string, metadata *models.Metadata) error {
	metadata.UpdatedAt = time.Now()
	
	err := s.repo.UpdateEndpointMetadata(ctx, endpointID, metadata)
	if err != nil {
		s.logger.Error("Failed to update endpoint metadata", "error", err, "endpoint_id", endpointID)
		return fmt.Errorf("failed to update endpoint metadata: %w", err)
	}

	s.logger.Info("Endpoint metadata updated", "endpoint_id", endpointID)
	return nil
}

func (s *MetadataService) ExtractMetadata(ctx context.Context, endpoint *models.Endpoint) (*models.Metadata, error) {
	metadata := &models.Metadata{
		ID:            uuid.New().String(),
		EndpointID:    endpoint.ID,
		APIID:         endpoint.APIID,
		URL:           endpoint.URL,
		Method:        endpoint.Method,
		Title:         s.generateTitle(endpoint),
		Description:   s.generateDescription(endpoint),
		Tags:          s.extractTags(endpoint),
		Category:      s.determineCategory(endpoint),
		Parameters:    s.extractParameters(endpoint),
		ResponseSchema: s.extractResponseSchema(endpoint),
		RequestSchema:  s.extractRequestSchema(endpoint),
		Examples:      s.generateExamples(endpoint),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Extract security metadata
	securityMetadata, err := s.AnalyzeEndpointSecurity(ctx, endpoint)
	if err != nil {
		s.logger.Warn("Failed to analyze endpoint security", "error", err, "endpoint_id", endpoint.ID)
	} else {
		metadata.Security = securityMetadata
	}

	// Calculate quality metrics
	qualityMetrics, err := s.CalculateQualityMetrics(ctx, endpoint.ID)
	if err != nil {
		s.logger.Warn("Failed to calculate quality metrics", "error", err, "endpoint_id", endpoint.ID)
	} else {
		metadata.Quality = qualityMetrics
	}

	// Initialize performance metrics
	metadata.Performance = &models.PerformanceMetrics{
		LastMeasured: time.Now(),
	}

	// Initialize usage metrics
	metadata.Usage = &models.UsageMetrics{
		LastUsed: time.Now(),
	}

	err = s.repo.SaveEndpointMetadata(ctx, metadata)
	if err != nil {
		s.logger.Error("Failed to save endpoint metadata", "error", err, "endpoint_id", endpoint.ID)
		return nil, fmt.Errorf("failed to save endpoint metadata: %w", err)
	}

	s.logger.Info("Endpoint metadata extracted and saved", "endpoint_id", endpoint.ID)
	return metadata, nil
}

func (s *MetadataService) GetAPISpecification(ctx context.Context, apiID string) (*models.APISpec, error) {
	spec, err := s.repo.GetAPISpecification(ctx, apiID)
	if err != nil {
		s.logger.Error("Failed to get API specification", "error", err, "api_id", apiID)
		return nil, fmt.Errorf("failed to get API specification: %w", err)
	}

	return spec, nil
}

func (s *MetadataService) GenerateAPISpec(ctx context.Context, apiID string) (*models.APISpec, error) {
	// Get all endpoints for the API
	endpoints, err := s.repo.GetAPIEndpoints(ctx, apiID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API endpoints: %w", err)
	}

	spec := &models.APISpec{
		ID:             uuid.New().String(),
		APIID:          apiID,
		Version:        "1.0.0",
		OpenAPIVersion: "3.0.3",
		Title:          "Generated API Specification",
		Description:    "Auto-generated API specification from discovered endpoints",
		Info: &models.SpecInfo{
			Title:       "Generated API",
			Description: "Auto-generated API specification",
			Version:     "1.0.0",
		},
		Paths:     make(map[string]*models.PathItem),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate paths from endpoints
	for _, endpoint := range endpoints {
		pathItem := s.convertEndpointToPathItem(&endpoint)
		spec.Paths[endpoint.Path] = pathItem
	}

	// Generate components
	spec.Components = s.generateComponents(endpoints)

	// Generate tags
	spec.Tags = s.generateSpecTags(endpoints)

	err = s.repo.SaveAPISpecification(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to save API specification: %w", err)
	}

	s.logger.Info("API specification generated", "api_id", apiID, "endpoints_count", len(endpoints))
	return spec, nil
}

func (s *MetadataService) AnalyzeEndpointSecurity(ctx context.Context, endpoint *models.Endpoint) (*models.SecurityMetadata, error) {
	security := &models.SecurityMetadata{
		HasHTTPS:              strings.HasPrefix(endpoint.URL, "https://"),
		SecurityHeaders:       make(map[string]string),
		VulnerabilityScans:    []models.VulnerabilityScan{},
		ComplianceStatus:      make(map[string]string),
		LastSecurityScan:      nil,
	}

	// Check for HTTPS
	security.HasHTTPS = strings.HasPrefix(endpoint.URL, "https://")

	// Analyze headers for security
	security.HasSecurityHeaders = s.hasSecurityHeaders(endpoint.Headers)
	security.SecurityHeaders = s.extractSecurityHeaders(endpoint.Headers)

	// Check authentication requirements
	security.AuthenticationRequired = s.requiresAuthentication(endpoint)
	security.AuthenticationMethods = s.detectAuthMethods(endpoint)

	// Determine data classification
	security.DataClassification = s.classifyDataSensitivity(endpoint)

	// Check encryption
	security.EncryptionInTransit = security.HasHTTPS
	security.EncryptionAtRest = false // Would need additional analysis

	return security, nil
}

func (s *MetadataService) CalculateQualityMetrics(ctx context.Context, endpointID string) (*models.QualityMetrics, error) {
	quality := &models.QualityMetrics{
		LastQualityCheck: time.Now(),
		QualityIssues:    []string{},
	}

	// Get endpoint metadata for analysis
	metadata, err := s.repo.GetEndpointMetadata(ctx, endpointID)
	if err != nil {
		// If no metadata exists, return basic quality metrics
		quality.DocumentationScore = 0.0
		quality.APIDesignScore = 50.0
		quality.ConsistencyScore = 50.0
		return quality, nil
	}

	// Calculate documentation score
	quality.DocumentationScore = s.calculateDocumentationScore(metadata)

	// Calculate API design score
	quality.APIDesignScore = s.calculateAPIDesignScore(metadata)

	// Calculate consistency score
	quality.ConsistencyScore = s.calculateConsistencyScore(metadata)

	// Identify quality issues
	quality.QualityIssues = s.identifyQualityIssues(metadata)

	return quality, nil
}

func (s *MetadataService) UpdateUsageMetrics(ctx context.Context, endpointID string, requestCount int64) error {
	metadata, err := s.repo.GetEndpointMetadata(ctx, endpointID)
	if err != nil {
		return fmt.Errorf("failed to get endpoint metadata: %w", err)
	}

	if metadata.Usage == nil {
		metadata.Usage = &models.UsageMetrics{}
	}

	// Update usage metrics
	metadata.Usage.TotalRequests += requestCount
	metadata.Usage.LastUsed = time.Now()

	// Update time-based metrics (simplified)
	now := time.Now()
	if now.Sub(metadata.UpdatedAt) < 24*time.Hour {
		metadata.Usage.RequestsLast24h += requestCount
	}
	if now.Sub(metadata.UpdatedAt) < 7*24*time.Hour {
		metadata.Usage.RequestsLastWeek += requestCount
	}
	if now.Sub(metadata.UpdatedAt) < 30*24*time.Hour {
		metadata.Usage.RequestsLastMonth += requestCount
	}

	err = s.repo.UpdateEndpointMetadata(ctx, endpointID, metadata)
	if err != nil {
		return fmt.Errorf("failed to update usage metrics: %w", err)
	}

	return nil
}

func (s *MetadataService) EnrichMetadata(ctx context.Context, metadata *models.Metadata) error {
	// Enrich with additional data sources
	s.enrichWithBusinessContext(metadata)
	s.enrichWithTechnicalContext(metadata)
	s.enrichWithComplianceInfo(metadata)

	return nil
}

// Helper methods

func (s *MetadataService) generateTitle(endpoint *models.Endpoint) string {
	if endpoint.Summary != "" {
		return endpoint.Summary
	}
	
	// Generate title from path and method
	path := strings.TrimPrefix(endpoint.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return fmt.Sprintf("%s %s", strings.ToUpper(endpoint.Method), strings.Title(parts[len(parts)-1]))
	}
	
	return fmt.Sprintf("%s %s", strings.ToUpper(endpoint.Method), endpoint.Path)
}

func (s *MetadataService) generateDescription(endpoint *models.Endpoint) string {
	if endpoint.Description != "" {
		return endpoint.Description
	}
	
	return fmt.Sprintf("API endpoint for %s %s", endpoint.Method, endpoint.Path)
}

func (s *MetadataService) extractTags(endpoint *models.Endpoint) []string {
	if len(endpoint.Tags) > 0 {
		return endpoint.Tags
	}
	
	// Generate tags from path
	var tags []string
	path := strings.TrimPrefix(endpoint.Path, "/")
	parts := strings.Split(path, "/")
	
	for _, part := range parts {
		if part != "" && !strings.Contains(part, "{") {
			tags = append(tags, part)
		}
	}
	
	return tags
}

func (s *MetadataService) determineCategory(endpoint *models.Endpoint) string {
	path := strings.ToLower(endpoint.Path)
	method := strings.ToUpper(endpoint.Method)
	
	// Categorize based on common patterns
	if strings.Contains(path, "/auth") || strings.Contains(path, "/login") {
		return "Authentication"
	}
	if strings.Contains(path, "/user") || strings.Contains(path, "/profile") {
		return "User Management"
	}
	if strings.Contains(path, "/admin") {
		return "Administration"
	}
	if method == "GET" && strings.Contains(path, "/search") {
		return "Search"
	}
	if method == "POST" && !strings.Contains(path, "/{") {
		return "Creation"
	}
	if method == "PUT" || method == "PATCH" {
		return "Update"
	}
	if method == "DELETE" {
		return "Deletion"
	}
	if method == "GET" {
		return "Retrieval"
	}
	
	return "General"
}

func (s *MetadataService) extractParameters(endpoint *models.Endpoint) []models.Parameter {
	return endpoint.Parameters
}

func (s *MetadataService) extractResponseSchema(endpoint *models.Endpoint) map[string]interface{} {
	schema := make(map[string]interface{})
	
	for statusCode, response := range endpoint.Responses {
		if response.Schema != nil {
			schema[statusCode] = response.Schema
		}
	}
	
	return schema
}

func (s *MetadataService) extractRequestSchema(endpoint *models.Endpoint) map[string]interface{} {
	schema := make(map[string]interface{})
	
	// Extract from parameters
	if len(endpoint.Parameters) > 0 {
		properties := make(map[string]interface{})
		required := []string{}
		
		for _, param := range endpoint.Parameters {
			properties[param.Name] = map[string]interface{}{
				"type":        param.Type,
				"description": param.Description,
			}
			if param.Required {
				required = append(required, param.Name)
			}
		}
		
		schema["type"] = "object"
		schema["properties"] = properties
		if len(required) > 0 {
			schema["required"] = required
		}
	}
	
	return schema
}

func (s *MetadataService) generateExamples(endpoint *models.Endpoint) []models.MetadataExample {
	var examples []models.MetadataExample
	
	// Generate basic example
	example := models.MetadataExample{
		Name:        "Basic Example",
		Description: fmt.Sprintf("Basic example for %s %s", endpoint.Method, endpoint.Path),
		Request:     make(map[string]interface{}),
		Response:    make(map[string]interface{}),
		StatusCode:  endpoint.StatusCode,
	}
	
	// Add parameter examples
	if len(endpoint.Parameters) > 0 {
		for _, param := range endpoint.Parameters {
			if param.Example != nil {
				example.Request[param.Name] = param.Example
			}
		}
	}
	
	examples = append(examples, example)
	return examples
}

func (s *MetadataService) convertEndpointToPathItem(endpoint *models.Endpoint) *models.PathItem {
	pathItem := &models.PathItem{
		Summary:     endpoint.Summary,
		Description: endpoint.Description,
		Parameters:  endpoint.Parameters,
	}
	
	operation := &models.Operation{
		Tags:        endpoint.Tags,
		Summary:     endpoint.Summary,
		Description: endpoint.Description,
		Parameters:  endpoint.Parameters,
		Responses:   s.convertResponses(endpoint.Responses),
	}
	
	switch strings.ToUpper(endpoint.Method) {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	case "OPTIONS":
		pathItem.Options = operation
	case "HEAD":
		pathItem.Head = operation
	}
	
	return pathItem
}

func (s *MetadataService) convertResponses(responses map[string]models.Response) map[string]models.Response {
	converted := make(map[string]models.Response)
	
	for statusCode, response := range responses {
		converted[statusCode] = models.Response{
			StatusCode:  response.StatusCode,
			Description: response.Description,
			Headers:     response.Headers,
			Schema:      response.Schema,
			Examples:    response.Examples,
		}
	}
	
	return converted
}

func (s *MetadataService) generateComponents(endpoints []models.Endpoint) *models.SpecComponents {
	components := &models.SpecComponents{
		Schemas:         make(map[string]*models.Schema),
		Responses:       make(map[string]models.Response),
		Parameters:      make(map[string]models.Parameter),
		SecuritySchemes: make(map[string]models.SecurityScheme),
	}
	
	// Extract common schemas from endpoints
	schemaMap := make(map[string]*models.Schema)
	
	for _, endpoint := range endpoints {
		for _, response := range endpoint.Responses {
			if response.Schema != nil {
				schemaName := fmt.Sprintf("%s%sResponse", 
					strings.Title(strings.ToLower(endpoint.Method)),
					s.pathToSchemaName(endpoint.Path))
				
				schema := s.convertToSchema(response.Schema)
				schemaMap[schemaName] = schema
			}
		}
	}
	
	components.Schemas = schemaMap
	return components
}

func (s *MetadataService) generateSpecTags(endpoints []models.Endpoint) []models.SpecTag {
	tagMap := make(map[string]models.SpecTag)
	
	for _, endpoint := range endpoints {
		for _, tag := range endpoint.Tags {
			if _, exists := tagMap[tag]; !exists {
				tagMap[tag] = models.SpecTag{
					Name:        tag,
					Description: fmt.Sprintf("Operations related to %s", tag),
				}
			}
		}
	}
	
	var tags []models.SpecTag
	for _, tag := range tagMap {
		tags = append(tags, tag)
	}
	
	return tags
}

func (s *MetadataService) hasSecurityHeaders(headers map[string]string) bool {
	securityHeaders := []string{
		"authorization",
		"x-api-key",
		"x-auth-token",
		"x-access-token",
		"bearer",
	}
	
	for key := range headers {
		for _, secHeader := range securityHeaders {
			if strings.Contains(strings.ToLower(key), secHeader) {
				return true
			}
		}
	}
	
	return false
}

func (s *MetadataService) extractSecurityHeaders(headers map[string]string) map[string]string {
	securityHeaders := make(map[string]string)
	
	for key, value := range headers {
		lowerKey := strings.ToLower(key)
		if strings.Contains(lowerKey, "auth") || 
		   strings.Contains(lowerKey, "security") ||
		   strings.Contains(lowerKey, "token") ||
		   strings.Contains(lowerKey, "api-key") {
			securityHeaders[key] = value
		}
	}
	
	return securityHeaders
}

func (s *MetadataService) requiresAuthentication(endpoint *models.Endpoint) bool {
	// Check headers for authentication indicators
	for key := range endpoint.Headers {
		lowerKey := strings.ToLower(key)
		if strings.Contains(lowerKey, "authorization") ||
		   strings.Contains(lowerKey, "x-api-key") ||
		   strings.Contains(lowerKey, "x-auth-token") {
			return true
		}
	}
	
	// Check status codes for auth-related responses
	for _, response := range endpoint.Responses {
		if response.StatusCode == 401 || response.StatusCode == 403 {
			return true
		}
	}
	
	return false
}

func (s *MetadataService) detectAuthMethods(endpoint *models.Endpoint) []string {
	var methods []string
	
	for key := range endpoint.Headers {
		lowerKey := strings.ToLower(key)
		if strings.Contains(lowerKey, "authorization") {
			methods = append(methods, "Bearer Token")
		}
		if strings.Contains(lowerKey, "x-api-key") {
			methods = append(methods, "API Key")
		}
		if strings.Contains(lowerKey, "basic") {
			methods = append(methods, "Basic Auth")
		}
	}
	
	return methods
}

func (s *MetadataService) classifyDataSensitivity(endpoint *models.Endpoint) string {
	path := strings.ToLower(endpoint.Path)
	
	// High sensitivity indicators
	if strings.Contains(path, "password") ||
	   strings.Contains(path, "credit") ||
	   strings.Contains(path, "payment") ||
	   strings.Contains(path, "ssn") ||
	   strings.Contains(path, "personal") {
		return "High"
	}
	
	// Medium sensitivity indicators
	if strings.Contains(path, "user") ||
	   strings.Contains(path, "profile") ||
	   strings.Contains(path, "account") ||
	   strings.Contains(path, "email") {
		return "Medium"
	}
	
	return "Low"
}

func (s *MetadataService) calculateDocumentationScore(metadata *models.Metadata) float64 {
	score := 0.0
	maxScore := 100.0
	
	// Title and description
	if metadata.Title != "" {
		score += 20
	}
	if metadata.Description != "" {
		score += 20
	}
	
	// Parameters documentation
	if len(metadata.Parameters) > 0 {
		documented := 0
		for _, param := range metadata.Parameters {
			if param.Description != "" {
				documented++
			}
		}
		score += (float64(documented) / float64(len(metadata.Parameters))) * 20
	} else {
		score += 20 // No parameters to document
	}
	
	// Examples
	if len(metadata.Examples) > 0 {
		score += 20
	}
	
	// Documentation section
	if metadata.Documentation != nil {
		if metadata.Documentation.Summary != "" {
			score += 10
		}
		if len(metadata.Documentation.ExternalDocs) > 0 {
			score += 10
		}
	}
	
	return (score / maxScore) * 100
}

func (s *MetadataService) calculateAPIDesignScore(metadata *models.Metadata) float64 {
	score := 0.0
	maxScore := 100.0
	
	// RESTful design principles
	if s.followsRESTfulNaming(metadata.URL, metadata.Method) {
		score += 30
	}
	
	// Consistent parameter naming
	if s.hasConsistentParameterNaming(metadata.Parameters) {
		score += 20
	}
	
	// Proper HTTP methods
	if s.usesProperHTTPMethod(metadata.URL, metadata.Method) {
		score += 25
	}
	
	// Response structure
	if len(metadata.ResponseSchema) > 0 {
		score += 25
	}
	
	return (score / maxScore) * 100
}

func (s *MetadataService) calculateConsistencyScore(metadata *models.Metadata) float64 {
	// This would require comparing with other endpoints in the same API
	// For now, return a base score
	return 75.0
}

func (s *MetadataService) identifyQualityIssues(metadata *models.Metadata) []string {
	var issues []string
	
	if metadata.Title == "" {
		issues = append(issues, "Missing endpoint title")
	}
	if metadata.Description == "" {
		issues = append(issues, "Missing endpoint description")
	}
	
	// Check parameter documentation
	for _, param := range metadata.Parameters {
		if param.Description == "" {
			issues = append(issues, fmt.Sprintf("Parameter '%s' lacks description", param.Name))
		}
	}
	
	// Check for examples
	if len(metadata.Examples) == 0 {
		issues = append(issues, "No examples provided")
	}
	
	// Check response schema
	if len(metadata.ResponseSchema) == 0 {
		issues = append(issues, "Missing response schema")
	}
	
	return issues
}

func (s *MetadataService) followsRESTfulNaming(url, method string) bool {
	// Basic RESTful naming checks
	path := strings.ToLower(url)
	
	// Check for proper resource naming (plural nouns)
	if method == "GET" && strings.Contains(path, "/") {
		return true // Simplified check
	}
	
	return true
}

func (s *MetadataService) hasConsistentParameterNaming(parameters []models.Parameter) bool {
	// Check for consistent naming conventions (snake_case, camelCase, etc.)
	if len(parameters) == 0 {
		return true
	}
	
	// Simplified consistency check
	return true
}

func (s *MetadataService) usesProperHTTPMethod(url, method string) bool {
	path := strings.ToLower(url)
	
	switch method {
	case "GET":
		return !strings.Contains(path, "create") && !strings.Contains(path, "update")
	case "POST":
		return strings.Contains(path, "create") || !strings.Contains(path, "/{id}")
	case "PUT", "PATCH":
		return strings.Contains(path, "/{") || strings.Contains(path, "update")
	case "DELETE":
		return strings.Contains(path, "/{") || strings.Contains(path, "delete")
	}
	
	return true
}

func (s *MetadataService) pathToSchemaName(path string) string {
	// Convert path to schema name
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var name string
	
	for _, part := range parts {
		if !strings.Contains(part, "{") {
			name += strings.Title(part)
		}
	}
	
	return name
}

func (s *MetadataService) convertToSchema(schemaData map[string]interface{}) *models.Schema {
	schema := &models.Schema{
		Type:       "object",
		Properties: make(map[string]*models.Schema),
	}
	
	if schemaType, ok := schemaData["type"].(string); ok {
		schema.Type = schemaType
	}
	
	if properties, ok := schemaData["properties"].(map[string]interface{}); ok {
		for key, value := range properties {
			if propMap, ok := value.(map[string]interface{}); ok {
				schema.Properties[key] = s.convertToSchema(propMap)
			}
		}
	}
	
	return schema
}

func (s *MetadataService) enrichWithBusinessContext(metadata *models.Metadata) {
	// Enrich with business context based on endpoint patterns
	path := strings.ToLower(metadata.URL)
	
	if strings.Contains(path, "/payment") || strings.Contains(path, "/billing") {
		metadata.BusinessOwner = "Finance Team"
		metadata.DataSensitivity = "High"
	} else if strings.Contains(path, "/user") || strings.Contains(path, "/customer") {
		metadata.BusinessOwner = "Customer Success Team"
		metadata.DataSensitivity = "Medium"
	} else if strings.Contains(path, "/admin") {
		metadata.BusinessOwner = "IT Team"
		metadata.DataSensitivity = "High"
	}
}

func (s *MetadataService) enrichWithTechnicalContext(metadata *models.Metadata) {
	// Enrich with technical context
	if metadata.Security != nil && metadata.Security.HasHTTPS {
		if metadata.Performance == nil {
			metadata.Performance = &models.PerformanceMetrics{}
		}
		metadata.Performance.SecurityOverhead = 10 // ms
	}
}

func (s *MetadataService) enrichWithComplianceInfo(metadata *models.Metadata) {
	// Enrich with compliance information
	if metadata.DataSensitivity == "High" {
		metadata.ComplianceReqs = []string{"GDPR", "CCPA", "SOX"}
	} else if metadata.DataSensitivity == "Medium" {
		metadata.ComplianceReqs = []string{"GDPR", "CCPA"}
	}
}
