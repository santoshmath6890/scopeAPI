package repository

import (
	"context"
	"fmt"
	"time"
	"scopeapi.local/backend/services/threat-detection/internal/models"
)

type AnomalyRepositoryInterface interface {
	// GetRecentAnomalies returns recent anomalies for an entity within a time window
	GetRecentAnomalies(ctx context.Context, entityID string, entityType string, since time.Time) ([]models.Anomaly, error)
	// GetAnomalies returns anomalies based on filter criteria
	GetAnomalies(ctx context.Context, filter *models.AnomalyFilter) ([]models.Anomaly, error)
	// SaveAnomaly saves a detected anomaly
	SaveAnomaly(ctx context.Context, anomaly *models.Anomaly) error
	// CreateAnomaly creates a new anomaly record
	CreateAnomaly(ctx context.Context, anomaly *models.Anomaly) error
	// GetBaselineStatistics returns baseline statistics for an entity
	GetBaselineStatistics(ctx context.Context, entityID string, entityType string) (map[string]interface{}, error)
	// GetRecentRequestCount returns the recent request count for an entity
	GetRecentRequestCount(ctx context.Context, entityID string, entityType string, duration time.Duration) (int, error)
	// GetBaselineRequestCount returns the baseline request count for an entity
	GetBaselineRequestCount(ctx context.Context, entityID string, entityType string) (int, error)
	// GetBaselineResponseTime returns the baseline response time for an entity
	GetBaselineResponseTime(ctx context.Context, entityID string, entityType string) (float64, error)
	// GetHistoricalCountries returns a list of historical countries for an entity
	GetHistoricalCountries(ctx context.Context, entityID string, entityType string) ([]string, error)
	// GetAnomaly returns a single anomaly by ID
	GetAnomaly(ctx context.Context, anomalyID string) (*models.Anomaly, error)
	// UpdateAnomalyFeedback updates feedback for an anomaly
	UpdateAnomalyFeedback(ctx context.Context, feedback *models.AnomalyFeedback) error
	// GetAnomalyStatistics returns statistics for anomalies
	GetAnomalyStatistics(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyStatistics, error)
	// StoreBaselineStatistics stores baseline statistics for an entity
	StoreBaselineStatistics(ctx context.Context, entityID string, baseline map[string]interface{}) error
	// GetModelPerformance returns model performance metrics
	GetModelPerformance(ctx context.Context, modelVersion string) (*models.ModelPerformanceMetric, error)
}

type ThreatRepositoryInterface interface {
	// GetThreatByID fetches a threat by its ID
	GetThreatByID(ctx context.Context, threatID string) (*models.Threat, error)
	// SaveThreat saves a threat event
	SaveThreat(ctx context.Context, threat *models.Threat) error
	// ListThreats returns a list of threats matching a filter
	ListThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error)

	// Signature management methods
	GetThreatSignatures(ctx context.Context, filter *models.SignatureFilter) ([]models.ThreatSignature, error)
	UpdateThreatSignature(ctx context.Context, id string, signature *models.ThreatSignature) error
	CreateThreatSignature(ctx context.Context, signature *models.ThreatSignature) error
	DeleteThreatSignature(ctx context.Context, id string) error
	GetSignatureMatchStatistics(ctx context.Context) (*models.SignatureMatchStats, error)
	
	// Additional threat management methods
	GetThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error)
	GetThreat(ctx context.Context, threatID string) (*models.Threat, error)
	CreateThreat(ctx context.Context, threat *models.Threat) error
	UpdateThreat(ctx context.Context, threatID string, threat *models.Threat) error
	DeleteThreat(ctx context.Context, threatID string) error
	GetThreatStatistics(ctx context.Context, timeRange time.Duration) (*models.ThreatStatistics, error)
	GetRequestCountByIP(ctx context.Context, ipAddress string, timeWindow time.Duration) (int, error)
	GetFailedAuthAttempts(ctx context.Context, ipAddress string, timeWindow time.Duration) (int, error)
}

type PatternRepositoryInterface interface {
	GetBehaviorPattern(ctx context.Context, patternID string) (*models.BehaviorPattern, error)
	SaveBehaviorPattern(ctx context.Context, pattern *models.BehaviorPattern) error
	ListBehaviorPatterns(ctx context.Context, filter *models.BehaviorPatternFilter) ([]models.BehaviorPattern, error)
}

// In-memory implementation of ThreatRepositoryInterface

type MemoryThreatRepository struct {
	threats     map[string]*models.Threat
	signatures  map[string]*models.ThreatSignature
}

type MemoryPatternRepository struct {
	patterns map[string]*models.BehaviorPattern
}

type MemoryAnomalyRepository struct {
	anomalies map[string]*models.Anomaly
}

func NewMemoryThreatRepository() *MemoryThreatRepository {
	return &MemoryThreatRepository{
		threats:    make(map[string]*models.Threat),
		signatures: make(map[string]*models.ThreatSignature),
	}
}

func (r *MemoryThreatRepository) GetThreatByID(ctx context.Context, threatID string) (*models.Threat, error) {
	if threat, ok := r.threats[threatID]; ok {
		return threat, nil
	}
	return nil, nil
}

func (r *MemoryThreatRepository) SaveThreat(ctx context.Context, threat *models.Threat) error {
	r.threats[threat.ID] = threat
	return nil
}

func (r *MemoryThreatRepository) ListThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error) {
	var result []models.Threat
	for _, threat := range r.threats {
		result = append(result, *threat)
	}
	return result, nil
}

// Signature management methods
func (r *MemoryThreatRepository) GetThreatSignatures(ctx context.Context, filter *models.SignatureFilter) ([]models.ThreatSignature, error) {
	var result []models.ThreatSignature
	for _, sig := range r.signatures {
		// Basic filtering by severity, pattern, signature set, enabled
		if filter != nil {
			if filter.Severity != "" && sig.Severity != filter.Severity {
				continue
			}
			if filter.Pattern != "" && sig.Pattern != filter.Pattern {
				continue
			}
			if filter.SignatureSet != "" && sig.SignatureSet != filter.SignatureSet {
				continue
			}
			if filter.Enabled && !sig.Enabled {
				continue
			}
		}
		result = append(result, *sig)
	}
	return result, nil
}

func (r *MemoryThreatRepository) UpdateThreatSignature(ctx context.Context, id string, signature *models.ThreatSignature) error {
	if _, ok := r.signatures[id]; !ok {
		return nil // Not found
	}
	r.signatures[id] = signature
	return nil
}

func (r *MemoryThreatRepository) CreateThreatSignature(ctx context.Context, signature *models.ThreatSignature) error {
	r.signatures[signature.ID] = signature
	return nil
}

func (r *MemoryThreatRepository) DeleteThreatSignature(ctx context.Context, id string) error {
	delete(r.signatures, id)
	return nil
}

func (r *MemoryThreatRepository) GetSignatureMatchStatistics(ctx context.Context) (*models.SignatureMatchStats, error) {
	return &models.SignatureMatchStats{
		TotalMatches:      0,
		MatchesByType:     make(map[string]int),
		MatchesByCategory: make(map[string]int),
		MatchesBySeverity: make(map[string]int),
	}, nil
}

func (r *MemoryThreatRepository) GetThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error) {
	return r.ListThreats(ctx, filter)
}

func (r *MemoryThreatRepository) GetThreat(ctx context.Context, threatID string) (*models.Threat, error) {
	return r.GetThreatByID(ctx, threatID)
}

func (r *MemoryThreatRepository) CreateThreat(ctx context.Context, threat *models.Threat) error {
	return r.SaveThreat(ctx, threat)
}

func (r *MemoryThreatRepository) UpdateThreat(ctx context.Context, threatID string, threat *models.Threat) error {
	r.threats[threatID] = threat
	return nil
}

func (r *MemoryThreatRepository) DeleteThreat(ctx context.Context, threatID string) error {
	delete(r.threats, threatID)
	return nil
}

func (r *MemoryThreatRepository) GetThreatStatistics(ctx context.Context, timeRange time.Duration) (*models.ThreatStatistics, error) {
	return &models.ThreatStatistics{
		TotalThreats:      int64(len(r.threats)),
		ActiveThreats:     0,
		ResolvedThreats:   0,
		CriticalThreats:   0,
		HighThreats:       0,
		MediumThreats:     0,
		LowThreats:        0,
		ThreatsByType:     make(map[string]int64),
		ThreatsBySource:   make(map[string]int64),
		RecentThreats:     0,
		TrendData:         []models.ThreatTrendPoint{},
		TopTargetedAPIs:   []models.APIThreatSummary{},
		TopAttackerIPs:    []models.IPThreatSummary{},
	}, nil
}

func (r *MemoryThreatRepository) GetRequestCountByIP(ctx context.Context, ipAddress string, timeWindow time.Duration) (int, error) {
	// Simple implementation - count threats from this IP in time window
	count := 0
	cutoff := time.Now().Add(-timeWindow)
	for _, threat := range r.threats {
		if threat.IPAddress == ipAddress && threat.CreatedAt.After(cutoff) {
			count++
		}
	}
	return count, nil
}

func (r *MemoryThreatRepository) GetFailedAuthAttempts(ctx context.Context, ipAddress string, timeWindow time.Duration) (int, error) {
	// Simple implementation - count auth-related threats from this IP
	count := 0
	cutoff := time.Now().Add(-timeWindow)
	for _, threat := range r.threats {
		if threat.IPAddress == ipAddress && threat.CreatedAt.After(cutoff) && 
		   (threat.Type == "authentication" || threat.AttackType == "brute_force") {
			count++
		}
	}
	return count, nil
}

// Constructor functions
func NewThreatRepository(db interface{}) ThreatRepositoryInterface {
	// For now, return the in-memory implementation
	// In the future, this could be extended to use the database connection
	return &MemoryThreatRepository{
		threats:    make(map[string]*models.Threat),
		signatures: make(map[string]*models.ThreatSignature),
	}
}

func NewPatternRepository(db interface{}) PatternRepositoryInterface {
	// For now, return the in-memory implementation
	return &MemoryPatternRepository{
		patterns: make(map[string]*models.BehaviorPattern),
	}
}

func NewAnomalyRepository(db interface{}) AnomalyRepositoryInterface {
	// For now, return the in-memory implementation
	return &MemoryAnomalyRepository{
		anomalies: make(map[string]*models.Anomaly),
	}
}

// MemoryPatternRepository implementations
func (r *MemoryPatternRepository) GetBehaviorPattern(ctx context.Context, patternID string) (*models.BehaviorPattern, error) {
	pattern, exists := r.patterns[patternID]
	if !exists {
		return nil, fmt.Errorf("pattern not found: %s", patternID)
	}
	return pattern, nil
}

func (r *MemoryPatternRepository) SaveBehaviorPattern(ctx context.Context, pattern *models.BehaviorPattern) error {
	r.patterns[pattern.ID] = pattern
	return nil
}

func (r *MemoryPatternRepository) ListBehaviorPatterns(ctx context.Context, filter *models.BehaviorPatternFilter) ([]models.BehaviorPattern, error) {
	var patterns []models.BehaviorPattern
	for _, pattern := range r.patterns {
		// Apply filter if needed
		if filter != nil {
			if filter.EntityID != "" && pattern.UserID != filter.EntityID {
				continue
			}
			if filter.EntityType != "" && pattern.Type != filter.EntityType {
				continue
			}
		}
		patterns = append(patterns, *pattern)
	}
	return patterns, nil
}

// MemoryAnomalyRepository implementations
func (r *MemoryAnomalyRepository) GetRecentAnomalies(ctx context.Context, entityID string, entityType string, since time.Time) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly
	for _, anomaly := range r.anomalies {
		if anomaly.FirstDetected.After(since) {
			anomalies = append(anomalies, *anomaly)
		}
	}
	return anomalies, nil
}

func (r *MemoryAnomalyRepository) GetAnomalies(ctx context.Context, filter *models.AnomalyFilter) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly
	for _, anomaly := range r.anomalies {
		// Apply filter if provided
		if filter != nil {
			if len(filter.Type) > 0 {
				found := false
				for _, t := range filter.Type {
					if anomaly.Type == t {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if len(filter.Severity) > 0 {
				found := false
				for _, s := range filter.Severity {
					if anomaly.Severity == s {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		anomalies = append(anomalies, *anomaly)
	}
	return anomalies, nil
}

func (r *MemoryAnomalyRepository) GetAnomaly(ctx context.Context, anomalyID string) (*models.Anomaly, error) {
	anomaly, exists := r.anomalies[anomalyID]
	if !exists {
		return nil, fmt.Errorf("anomaly not found: %s", anomalyID)
	}
	return anomaly, nil
}

func (r *MemoryAnomalyRepository) SaveAnomaly(ctx context.Context, anomaly *models.Anomaly) error {
	r.anomalies[anomaly.ID] = anomaly
	return nil
}

func (r *MemoryAnomalyRepository) CreateAnomaly(ctx context.Context, anomaly *models.Anomaly) error {
	return r.SaveAnomaly(ctx, anomaly)
}

// Stub implementations for remaining interface methods
func (r *MemoryAnomalyRepository) GetBaselineStatistics(ctx context.Context, entityID string, entityType string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (r *MemoryAnomalyRepository) GetRecentRequestCount(ctx context.Context, entityID string, entityType string, duration time.Duration) (int, error) {
	return 0, nil
}

func (r *MemoryAnomalyRepository) GetBaselineRequestCount(ctx context.Context, entityID string, entityType string) (int, error) {
	return 0, nil
}

func (r *MemoryAnomalyRepository) GetBaselineResponseTime(ctx context.Context, entityID string, entityType string) (float64, error) {
	return 0.0, nil
}

func (r *MemoryAnomalyRepository) GetHistoricalCountries(ctx context.Context, entityID string, entityType string) ([]string, error) {
	return []string{}, nil
}

func (r *MemoryAnomalyRepository) UpdateAnomalyFeedback(ctx context.Context, feedback *models.AnomalyFeedback) error {
	// TODO: Implement anomaly feedback update
	return nil
}

func (r *MemoryAnomalyRepository) GetAnomalyStatistics(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyStatistics, error) {
	// TODO: Implement anomaly statistics
	return &models.AnomalyStatistics{}, nil
}

func (r *MemoryAnomalyRepository) StoreBaselineStatistics(ctx context.Context, entityID string, baseline map[string]interface{}) error {
	// TODO: Implement baseline statistics storage
	return nil
}

func (r *MemoryAnomalyRepository) GetModelPerformance(ctx context.Context, modelVersion string) (*models.ModelPerformanceMetric, error) {
	// TODO: Implement model performance retrieval
	return &models.ModelPerformanceMetric{}, nil
}
