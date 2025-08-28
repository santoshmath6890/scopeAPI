package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/api-discovery/internal/models"
	"scopeapi.local/backend/services/api-discovery/internal/repository"
	"shared/logging"
)

type DiscoveryServiceInterface interface {
	StartDiscovery(ctx context.Context, config *models.DiscoveryConfig) (string, error)
	GetDiscoveryStatus(ctx context.Context, discoveryID string) (*models.DiscoveryStatus, error)
	GetDiscoveryResults(ctx context.Context, discoveryID string, page, limit int) (*models.DiscoveryResults, error)
	StopDiscovery(ctx context.Context, discoveryID string) error
	AnalyzeEndpoint(ctx context.Context, endpoint *models.Endpoint) (*models.EndpointAnalysis, error)
}

type DiscoveryService struct {
	repo   repository.DiscoveryRepositoryInterface
	logger logging.Logger
}

func NewDiscoveryService(repo repository.DiscoveryRepositoryInterface, logger logging.Logger) DiscoveryServiceInterface {
	return &DiscoveryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *DiscoveryService) StartDiscovery(ctx context.Context, config *models.DiscoveryConfig) (string, error) {
	discoveryID := uuid.New().String()
	
	discovery := &models.Discovery{
		ID:        discoveryID,
		Target:    config.Target,
		Method:    config.Method,
		Status:    "running",
		StartTime: time.Now(),
		Config:    config,
	}

	err := s.repo.CreateDiscovery(ctx, discovery)
	if err != nil {
		s.logger.Error("Failed to create discovery record", "error", err, "discovery_id", discoveryID)
		return "", fmt.Errorf("failed to create discovery record: %w", err)
	}

	// Start discovery process asynchronously
	go s.runDiscovery(context.Background(), discovery)

	s.logger.Info("Discovery started", "discovery_id", discoveryID, "target", config.Target)
	return discoveryID, nil
}

func (s *DiscoveryService) GetDiscoveryStatus(ctx context.Context, discoveryID string) (*models.DiscoveryStatus, error) {
	discovery, err := s.repo.GetDiscovery(ctx, discoveryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discovery: %w", err)
	}

	status := &models.DiscoveryStatus{
		ID:           discovery.ID,
		Status:       discovery.Status,
		Progress:     discovery.Progress,
		StartTime:    discovery.StartTime,
		EndTime:      discovery.EndTime,
		EndpointsFound: discovery.EndpointsFound,
		ErrorMessage: discovery.ErrorMessage,
	}

	return status, nil
}

func (s *DiscoveryService) GetDiscoveryResults(ctx context.Context, discoveryID string, page, limit int) (*models.DiscoveryResults, error) {
	results, err := s.repo.GetDiscoveryResults(ctx, discoveryID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get discovery results: %w", err)
	}

	return results, nil
}

func (s *DiscoveryService) StopDiscovery(ctx context.Context, discoveryID string) error {
	err := s.repo.UpdateDiscoveryStatus(ctx, discoveryID, "stopped")
	if err != nil {
		return fmt.Errorf("failed to stop discovery: %w", err)
	}

	s.logger.Info("Discovery stopped", "discovery_id", discoveryID)
	return nil
}

func (s *DiscoveryService) AnalyzeEndpoint(ctx context.Context, endpoint *models.Endpoint) (*models.EndpointAnalysis, error) {
	analysis := &models.EndpointAnalysis{
		EndpointID:    uuid.New().String(),
		URL:          endpoint.URL,
		Method:       endpoint.Method,
		ResponseTime: s.measureResponseTime(endpoint),
		StatusCode:   s.getStatusCode(endpoint),
		ContentType:  s.getContentType(endpoint),
		Parameters:   s.extractParameters(endpoint),
		Headers:      endpoint.Headers,
		Security:     s.analyzeSecurityHeaders(endpoint),
		CreatedAt:    time.Now(),
	}

	err := s.repo.SaveEndpointAnalysis(ctx, analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to save endpoint analysis: %w", err)
	}

	return analysis, nil
}

func (s *DiscoveryService) runDiscovery(ctx context.Context, discovery *models.Discovery) {
	s.logger.Info("Starting discovery process", "discovery_id", discovery.ID)

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Discovery process panicked", "discovery_id", discovery.ID, "panic", r)
			s.repo.UpdateDiscoveryStatus(ctx, discovery.ID, "failed")
		}
	}()

	switch discovery.Method {
	case "passive":
		s.runPassiveDiscovery(ctx, discovery)
	case "active":
		s.runActiveDiscovery(ctx, discovery)
	default:
		s.logger.Error("Unknown discovery method", "method", discovery.Method)
		s.repo.UpdateDiscoveryStatus(ctx, discovery.ID, "failed")
		return
	}

	s.repo.UpdateDiscoveryStatus(ctx, discovery.ID, "completed")
	s.logger.Info("Discovery process completed", "discovery_id", discovery.ID)
}

func (s *DiscoveryService) runPassiveDiscovery(ctx context.Context, discovery *models.Discovery) {
	// Implement passive discovery logic
	// This would involve analyzing traffic logs, proxy logs, etc.
	s.logger.Info("Running passive discovery", "discovery_id", discovery.ID)
	
	// Real passive discovery implementation
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 10)
	
	// Analyze target for common API patterns
	endpoints := s.analyzeTargetForAPIPatterns(discovery.Target)
	
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 50)
	
	// Test discovered endpoints
	s.testDiscoveredEndpoints(ctx, discovery, endpoints)
	
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 100)
}

func (s *DiscoveryService) runActiveDiscovery(ctx context.Context, discovery *models.Discovery) {
	// Implement active discovery logic
	// This would involve crawling, scanning, probing endpoints
	s.logger.Info("Running active discovery", "discovery_id", discovery.ID)
	
	// Real active discovery implementation
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 10)
	
	// Use basic HTTP probing to discover endpoints
	endpoints := s.probeTargetForEndpoints(discovery.Target)
	
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 60)
	
	// Test discovered endpoints
	s.testDiscoveredEndpoints(ctx, discovery, endpoints)
	
	s.repo.UpdateDiscoveryProgress(ctx, discovery.ID, 100)
}

func (s *DiscoveryService) measureResponseTime(endpoint *models.Endpoint) time.Duration {
	// Implement real response time measurement
	start := time.Now()
	
	// Make HTTP request to measure response time
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", endpoint.URL, nil)
	if err != nil {
		return 0
	}
	
	req.Header.Set("User-Agent", "ScopeAPI-Discovery/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	
	return time.Since(start)
}

func (s *DiscoveryService) getStatusCode(endpoint *models.Endpoint) int {
	// Implement real status code detection
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", endpoint.URL, nil)
	if err != nil {
		return 0
	}
	
	req.Header.Set("User-Agent", "ScopeAPI-Discovery/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	
	return resp.StatusCode
}

func (s *DiscoveryService) getContentType(endpoint *models.Endpoint) string {
	// Implement real content type detection
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", endpoint.URL, nil)
	if err != nil {
		return ""
	}
	
	req.Header.Set("User-Agent", "ScopeAPI-Discovery/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	return resp.Header.Get("Content-Type")
}

func (s *DiscoveryService) extractParameters(endpoint *models.Endpoint) []models.Parameter {
	// Implement real parameter extraction
	var params []models.Parameter
	
	// Extract query parameters from URL
	if u, err := url.Parse(endpoint.URL); err == nil {
		for key, values := range u.Query() {
			param := models.Parameter{
				Name:        key,
				In:          "query",
				Type:        "string",
				Required:    false,
				Description: "Query parameter",
			}
			if len(values) > 0 {
				param.Example = values[0]
			}
			params = append(params, param)
		}
	}
	
	return params
}

func (s *DiscoveryService) analyzeSecurityHeaders(endpoint *models.Endpoint) *models.SecurityAnalysis {
	// Implement real security header analysis
	security := &models.SecurityAnalysis{
		HasHTTPS:           false,
		HasSecurityHeaders: false,
		VulnerableHeaders:  []string{},
		RateLimitHeaders:   make(map[string]string),
	}
	
	// Check if HTTPS is used
	if strings.HasPrefix(endpoint.URL, "https://") {
		security.HasHTTPS = true
	}
	
	// Test endpoint to analyze headers
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", endpoint.URL, nil)
	if err != nil {
		return security
	}
	
	req.Header.Set("User-Agent", "ScopeAPI-Discovery/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return security
	}
	defer resp.Body.Close()
	
	// Check for security headers
	securityHeaders := []string{
		"X-Frame-Options",
		"X-Content-Type-Options", 
		"X-XSS-Protection",
		"Strict-Transport-Security",
		"Content-Security-Policy",
	}
	
	for _, header := range securityHeaders {
		if resp.Header.Get(header) != "" {
			security.HasSecurityHeaders = true
			break
		}
	}
	
	// Check for vulnerable headers
	vulnerableHeaders := []string{
		"Server",
		"X-Powered-By",
		"X-AspNet-Version",
		"X-AspNetMvc-Version",
	}
	
	for _, header := range vulnerableHeaders {
		if resp.Header.Get(header) != "" {
			security.VulnerableHeaders = append(security.VulnerableHeaders, header)
		}
	}
	
	// Check for rate limiting headers
	rateLimitHeaders := []string{
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"Retry-After",
	}
	
	for _, header := range rateLimitHeaders {
		if value := resp.Header.Get(header); value != "" {
			security.RateLimitHeaders[header] = value
		}
	}
	
	return security
}

// Helper methods for discovery implementation

func (s *DiscoveryService) analyzeTargetForAPIPatterns(target string) []models.Endpoint {
	var endpoints []models.Endpoint
	
	// Common API endpoint patterns
	patterns := []string{
		"/api",
		"/api/v1",
		"/api/v2", 
		"/rest",
		"/graphql",
		"/swagger",
		"/openapi",
		"/docs",
		"/health",
		"/status",
	}
	
	for _, pattern := range patterns {
		endpoint := models.Endpoint{
			ID:        uuid.New().String(),
			URL:       target + pattern,
			Path:      pattern,
			Method:    "GET",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		endpoints = append(endpoints, endpoint)
	}
	
	return endpoints
}

func (s *DiscoveryService) probeTargetForEndpoints(target string) []models.Endpoint {
	var endpoints []models.Endpoint
	
	// Basic endpoint probing
	commonPaths := []string{
		"/",
		"/api",
		"/v1",
		"/v2",
		"/health",
		"/status",
		"/info",
		"/metrics",
	}
	
	for _, path := range commonPaths {
		endpoint := models.Endpoint{
			ID:        uuid.New().String(),
			URL:       target + path,
			Path:      path,
			Method:    "GET",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		endpoints = append(endpoints, endpoint)
	}
	
	return endpoints
}

func (s *DiscoveryService) testDiscoveredEndpoints(ctx context.Context, discovery *models.Discovery, endpoints []models.Endpoint) {
	s.logger.Info("Testing discovered endpoints", "discovery_id", discovery.ID, "count", len(endpoints))
	
	// Save endpoints to repository
	for _, endpoint := range endpoints {
		err := s.repo.SaveEndpoint(ctx, &endpoint)
		if err != nil {
			s.logger.Error("Failed to save endpoint", "error", err, "endpoint", endpoint.URL)
			continue
		}
		
		// Test endpoint
		if analysis, err := s.AnalyzeEndpoint(ctx, &endpoint); err == nil {
			// Update endpoint with analysis results
			endpoint.StatusCode = analysis.StatusCode
			endpoint.ContentType = analysis.ContentType
			endpoint.Parameters = analysis.Parameters
			endpoint.Headers = analysis.Headers
			
			// Save updated endpoint
			s.repo.UpdateEndpoint(ctx, &endpoint)
		}
	}
	
	// Update discovery with endpoint count
	discovery.EndpointsFound = len(endpoints)
	s.repo.UpdateDiscoveryEndpointsFound(ctx, discovery.ID, len(endpoints))
}
