package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"data-protection/internal/models"
	"data-protection/internal/repository"

	"github.com/google/uuid"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type RiskScoringServiceInterface interface {
	AssessRisk(ctx context.Context, request *models.RiskAssessmentRequest) (*models.RiskScoringResult, error)
	GetRiskScores(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error) // Returns models.RiskAssessment as "scores"
	GetRiskScore(ctx context.Context, id string) (*models.RiskAssessment, error)
	CreateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error
	GetMitigationPlan(ctx context.Context, id string) (*models.MitigationPlan, error)
	UpdateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error
	DeleteMitigationPlan(ctx context.Context, id string) error
	GetRiskAssessments(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error)

	CalculateRiskScore(ctx context.Context, request *models.RiskAssessmentRequest) (*models.RiskScoringResult, error)
	CreateRiskProfile(ctx context.Context, profile *models.RiskProfile) error
	UpdateRiskProfile(ctx context.Context, profileID string, profile *models.RiskProfile) error
	GetRiskProfiles(ctx context.Context, filter *models.RiskProfileFilter) ([]models.RiskProfile, error)
	GetRiskTrends(ctx context.Context, filter *models.RiskTrendFilter) (*models.RiskTrendAnalysis, error)
	GenerateRiskReport(ctx context.Context, filter *models.RiskReportFilter) (*models.RiskReport, error)
}

type RiskScoringService struct {
	piiRepo       repository.PIIRepositoryInterface
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
	riskProfiles  map[string]*models.RiskProfile
	scoringRules  map[string]*models.RiskScoringRule
	weights       *models.RiskWeights
}

func NewRiskScoringService(
	piiRepo repository.PIIRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *RiskScoringService {
	service := &RiskScoringService{
		piiRepo:       piiRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		riskProfiles:  make(map[string]*models.RiskProfile),
		scoringRules:  make(map[string]*models.RiskScoringRule),
		weights:       &models.RiskWeights{},
	}

	// Load default risk profiles and rules
	service.loadDefaultRiskProfiles()
	service.loadDefaultScoringRules()
	service.setupDefaultWeights()

	return service
}

func (s *RiskScoringService) loadDefaultRiskProfiles() {
	defaultProfiles := []models.RiskProfile{
		{
			ID:          "low_risk_profile",
			Name:        "Low Risk Profile",
			Description: "Profile for low-risk data and operations",
			Category:    models.RiskCategoryLow,
			BaseScore:   10.0,
			Multipliers: map[string]float64{
				"pii_detected":     1.2,
				"public_endpoint":  1.1,
				"encrypted":        0.8,
				"internal_network": 0.9,
			},
			Thresholds: models.RiskThresholds{
				Low:      0.0,
				Medium:   25.0,
				High:     50.0,
				Critical: 75.0,
			},
			Enabled: true,
		},
		{
			ID:          "medium_risk_profile",
			Name:        "Medium Risk Profile",
			Description: "Profile for medium-risk data and operations",
			Category:    models.RiskCategoryMedium,
			BaseScore:   30.0,
			Multipliers: map[string]float64{
				"pii_detected":     1.5,
				"public_endpoint":  1.3,
				"encrypted":        0.7,
				"internal_network": 0.8,
				"sensitive_data":   1.4,
			},
			Thresholds: models.RiskThresholds{
				Low:      0.0,
				Medium:   35.0,
				High:     60.0,
				Critical: 80.0,
			},
			Enabled: true,
		},
		{
			ID:          "high_risk_profile",
			Name:        "High Risk Profile",
			Description: "Profile for high-risk data and operations",
			Category:    models.RiskCategoryHigh,
			BaseScore:   60.0,
			Multipliers: map[string]float64{
				"pii_detected":     2.0,
				"public_endpoint":  1.8,
				"encrypted":        0.6,
				"internal_network": 0.7,
				"sensitive_data":   1.8,
				"financial_data":   2.2,
				"health_data":      2.5,
			},
			Thresholds: models.RiskThresholds{
				Low:      0.0,
				Medium:   45.0,
				High:     70.0,
				Critical: 85.0,
			},
			Enabled: true,
		},
		{
			ID:          "critical_risk_profile",
			Name:        "Critical Risk Profile",
			Description: "Profile for critical and highly sensitive data",
			Category:    models.RiskCategoryCritical,
			BaseScore:   80.0,
			Multipliers: map[string]float64{
				"pii_detected":     2.5,
				"public_endpoint":  2.2,
				"encrypted":        0.5,
				"internal_network": 0.6,
				"sensitive_data":   2.2,
				"financial_data":   2.8,
				"health_data":      3.0,
				"government_data":  3.2,
			},
			Thresholds: models.RiskThresholds{
				Low:      0.0,
				Medium:   55.0,
				High:     75.0,
				Critical: 90.0,
			},
			Enabled: true,
		},
	}

	for _, profile := range defaultProfiles {
		profile.CreatedAt = time.Now()
		profile.UpdatedAt = time.Now()
		s.riskProfiles[profile.ID] = &profile
	}
}

func (s *RiskScoringService) loadDefaultScoringRules() {
	defaultRules := []models.RiskScoringRule{
		{
			ID:          "pii_detection_rule",
			Name:        "PII Detection Risk",
			Description: "Increases risk score when PII is detected",
			Category:    "data_sensitivity",
			Priority:    1,
			Conditions: []models.RiskCondition{
				{
					Field:    "pii_detected",
					Operator: "equals",
					Value:    "true",
				},
			},
			ScoreAdjustment: models.ScoreAdjustment{
				Type:   "multiply",
				Value:  1.5,
				Reason: "PII data detected",
			},
			Enabled: true,
		},
		{
			ID:          "public_endpoint_rule",
			Name:        "Public Endpoint Risk",
			Description: "Increases risk for public-facing endpoints",
			Category:    "exposure",
			Priority:    2,
			Conditions: []models.RiskCondition{
				{
					Field:    "endpoint_visibility",
					Operator: "equals",
					Value:    "public",
				},
			},
			ScoreAdjustment: models.ScoreAdjustment{
				Type:   "add",
				Value:  15.0,
				Reason: "Public endpoint exposure",
			},
			Enabled: true,
		},
		{
			ID:          "encryption_rule",
			Name:        "Encryption Protection",
			Description: "Reduces risk when data is encrypted",
			Category:    "protection",
			Priority:    3,
			Conditions: []models.RiskCondition{
				{
					Field:    "encrypted",
					Operator: "equals",
					Value:    "true",
				},
			},
			ScoreAdjustment: models.ScoreAdjustment{
				Type:   "multiply",
				Value:  0.7,
				Reason: "Data encryption in place",
			},
			Enabled: true,
		},
		{
			ID:          "authentication_rule",
			Name:        "Authentication Protection",
			Description: "Reduces risk when strong authentication is present",
			Category:    "access_control",
			Priority:    4,
			Conditions: []models.RiskCondition{
				{
					Field:    "authentication_strength",
					Operator: "in",
					Value:    "strong,multi_factor",
				},
			},
			ScoreAdjustment: models.ScoreAdjustment{
				Type:   "multiply",
				Value:  0.8,
				Reason: "Strong authentication required",
			},
			Enabled: true,
		},
		{
			ID:          "vulnerability_rule",
			Name:        "Known Vulnerabilities",
			Description: "Increases risk when vulnerabilities are present",
			Category:    "vulnerabilities",
			Priority:    5,
			Conditions: []models.RiskCondition{
				{
					Field:    "vulnerability_count",
					Operator: "greater_than",
					Value:    "0",
				},
			},
			ScoreAdjustment: models.ScoreAdjustment{
				Type:   "add",
				Value:  20.0,
				Reason: "Known vulnerabilities present",
			},
			Enabled: true,
		},
	}

	for _, rule := range defaultRules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		s.scoringRules[rule.ID] = &rule
	}
}

func (s *RiskScoringService) setupDefaultWeights() {
	s.weights = &models.RiskWeights{
		DataSensitivity:  0.30,
		ExposureLevel:    0.25,
		AccessControls:   0.20,
		Vulnerabilities:  0.15,
		ComplianceStatus: 0.10,
	}
}

func (s *RiskScoringService) AssessRisk(ctx context.Context, request *models.RiskAssessmentRequest) (*models.RiskScoringResult, error) {
	return s.CalculateRiskScore(ctx, request)
}

func (s *RiskScoringService) CalculateRiskScore(ctx context.Context, request *models.RiskAssessmentRequest) (*models.RiskScoringResult, error) {
	startTime := time.Now()

	result := &models.RiskScoringResult{
		RequestID:       request.RequestID,
		RiskScore:       0.0,
		RiskLevel:       models.RiskLevelLow,
		ProfileUsed:     "",
		ScoreBreakdown:  models.RiskScoreBreakdown{},
		AppliedRules:    []models.AppliedRiskRule{},
		Recommendations: []string{},
		ProcessingTime:  0,
		CalculatedAt:    time.Now(),
	}

	// Determine appropriate risk profile
	profile := s.selectRiskProfile(request)
	if profile == nil {
		return nil, fmt.Errorf("no suitable risk profile found")
	}

	result.ProfileUsed = profile.ID

	// Start with base score
	currentScore := profile.BaseScore
	result.ScoreBreakdown.BaseScore = currentScore

	// Apply scoring rules
	ruleScores := make(map[string]float64)
	for _, rule := range s.getSortedRules() {
		if !rule.Enabled {
			continue
		}

		if s.evaluateRiskRule(rule, request) {
			adjustment := s.applyScoreAdjustment(currentScore, rule.ScoreAdjustment)
			ruleScores[rule.ID] = adjustment
			currentScore = adjustment

			appliedRule := models.AppliedRiskRule{
				RuleID:      rule.ID,
				RuleName:    rule.Name,
				Category:    rule.Category,
				Adjustment:  rule.ScoreAdjustment,
				ScoreChange: adjustment - result.ScoreBreakdown.BaseScore,
				Reason:      rule.ScoreAdjustment.Reason,
			}
			result.AppliedRules = append(result.AppliedRules, appliedRule)
		}
	}

	// Apply profile multipliers
	multiplierScore := s.applyProfileMultipliers(currentScore, profile, request)
	result.ScoreBreakdown.MultiplierAdjustment = multiplierScore - currentScore
	currentScore = multiplierScore

	// Apply weighted scoring
	weightedScore := s.applyWeightedScoring(currentScore, request)
	result.ScoreBreakdown.WeightedAdjustment = weightedScore - currentScore
	currentScore = weightedScore

	// Normalize score (0-100)
	result.RiskScore = s.normalizeScore(currentScore)
	result.RiskLevel = s.determineRiskLevel(result.RiskScore, profile)

	// Generate recommendations
	result.Recommendations = s.generateRiskRecommendations(result, profile, request)

	// Calculate processing time
	result.ProcessingTime = time.Since(startTime)

	// Store risk assessment
	riskAssessment := &models.RiskAssessment{
		ID:              uuid.New().String(),
		RequestID:       request.RequestID,
		APIID:           request.APIID,
		EndpointID:      request.EndpointID,
		RiskScore:       result.RiskScore,
		RiskLevel:       result.RiskLevel,
		ProfileID:       profile.ID,
		ScoreBreakdown:  result.ScoreBreakdown,
		AppliedRules:    result.AppliedRules,
		Recommendations: result.Recommendations,
		IPAddress:       request.IPAddress,
		UserAgent:       request.UserAgent,
		AssessedAt:      time.Now(),
		Metadata: map[string]interface{}{
			"processing_time": result.ProcessingTime.Milliseconds(),
			"rule_count":      len(result.AppliedRules),
			"data_factors":    request.DataFactors,
		},
	}

	if err := s.piiRepo.CreateRiskAssessment(ctx, riskAssessment); err != nil {
		s.logger.Error("Failed to store risk assessment", "error", err)
	}

	// Publish risk events
	if err := s.publishRiskEvents(ctx, result, request); err != nil {
		s.logger.Error("Failed to publish risk events", "error", err)
	}

	return result, nil
}

func (s *RiskScoringService) selectRiskProfile(request *models.RiskAssessmentRequest) *models.RiskProfile {
	// Default to medium risk profile
	defaultProfile := s.riskProfiles["medium_risk_profile"]

	// Check data factors to determine appropriate profile
	if request.DataFactors != nil {
		// High-risk data types
		if s.hasHighRiskData(request.DataFactors) {
			if profile, exists := s.riskProfiles["critical_risk_profile"]; exists {
				return profile
			}
		}

		// Medium-risk data types
		if s.hasMediumRiskData(request.DataFactors) {
			if profile, exists := s.riskProfiles["high_risk_profile"]; exists {
				return profile
			}
		}

		// Low-risk indicators
		if s.hasLowRiskData(request.DataFactors) {
			if profile, exists := s.riskProfiles["low_risk_profile"]; exists {
				return profile
			}
		}
	}

	return defaultProfile
}

func (s *RiskScoringService) hasHighRiskData(factors map[string]interface{}) bool {
	highRiskIndicators := []string{
		"health_data", "medical_records", "phi",
		"financial_data", "credit_card", "bank_account",
		"government_data", "classified", "national_security",
		"biometric_data", "genetic_data",
	}

	for _, indicator := range highRiskIndicators {
		if value, exists := factors[indicator]; exists {
			if boolValue, ok := value.(bool); ok && boolValue {
				return true
			}
			if stringValue, ok := value.(string); ok && stringValue == "true" {
				return true
			}
		}
	}

	return false
}

func (s *RiskScoringService) hasMediumRiskData(factors map[string]interface{}) bool {
	mediumRiskIndicators := []string{
		"pii_detected", "personal_data", "contact_info",
		"internal_data", "proprietary", "confidential",
		"employee_data", "customer_data",
	}

	for _, indicator := range mediumRiskIndicators {
		if value, exists := factors[indicator]; exists {
			if boolValue, ok := value.(bool); ok && boolValue {
				return true
			}
			if stringValue, ok := value.(string); ok && stringValue == "true" {
				return true
			}
		}
	}

	return false
}

func (s *RiskScoringService) hasLowRiskData(factors map[string]interface{}) bool {
	lowRiskIndicators := []string{
		"public_data", "marketing_data", "anonymous_data",
		"aggregated_data", "statistical_data",
	}

	for _, indicator := range lowRiskIndicators {
		if value, exists := factors[indicator]; exists {
			if boolValue, ok := value.(bool); ok && boolValue {
				return true
			}
			if stringValue, ok := value.(string); ok && stringValue == "true" {
				return true
			}
		}
	}

	return false
}

func (s *RiskScoringService) getSortedRules() []*models.RiskScoringRule {
	var rules []*models.RiskScoringRule
	for _, rule := range s.scoringRules {
		rules = append(rules, rule)
	}

	// Sort by priority (higher priority first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return rules
}

func (s *RiskScoringService) evaluateRiskRule(rule *models.RiskScoringRule, request *models.RiskAssessmentRequest) bool {
	for _, condition := range rule.Conditions {
		if !s.evaluateRiskCondition(condition, request) {
			return false
		}
	}
	return true
}

func (s *RiskScoringService) evaluateRiskCondition(condition models.RiskCondition, request *models.RiskAssessmentRequest) bool {
	var value interface{}

	// Get value from appropriate source
	if request.DataFactors != nil {
		if v, exists := request.DataFactors[condition.Field]; exists {
			value = v
		}
	}

	if value == nil && request.SecurityFactors != nil {
		if v, exists := request.SecurityFactors[condition.Field]; exists {
			value = v
		}
	}

	if value == nil && request.ContextFactors != nil {
		if v, exists := request.ContextFactors[condition.Field]; exists {
			value = v
		}
	}

	if value == nil {
		return false
	}

	return s.evaluateConditionValue(value, condition.Operator, condition.Value)
}

func (s *RiskScoringService) evaluateConditionValue(value interface{}, operator string, expectedValue interface{}) bool {
	valueStr := fmt.Sprintf("%v", value)
	expectedValueStr := fmt.Sprintf("%v", expectedValue)

	switch operator {
	case "equals":
		return valueStr == expectedValueStr
	case "not_equals":
		return valueStr != expectedValueStr
	case "contains":
		return strings.Contains(strings.ToLower(valueStr), strings.ToLower(expectedValueStr))
	case "in":
		values := strings.Split(expectedValueStr, ",")
		for _, v := range values {
			if strings.TrimSpace(v) == valueStr {
				return true
			}
		}
		return false
	case "greater_than":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValueStr); err == nil {
				return numValue > expectedNum
			}
		}
	case "less_than":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValueStr); err == nil {
				return numValue < expectedNum
			}
		}
	case "greater_equal":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValueStr); err == nil {
				return numValue >= expectedNum
			}
		}
	case "less_equal":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValueStr); err == nil {
				return numValue <= expectedNum
			}
		}
	}

	return false
}

func (s *RiskScoringService) parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func (s *RiskScoringService) applyScoreAdjustment(currentScore float64, adjustment models.ScoreAdjustment) float64 {
	switch adjustment.Type {
	case "add":
		return currentScore + adjustment.Value
	case "subtract":
		return currentScore - adjustment.Value
	case "multiply":
		return currentScore * adjustment.Value
	case "divide":
		if adjustment.Value != 0 {
			return currentScore / adjustment.Value
		}
		return currentScore
	case "set":
		return adjustment.Value
	default:
		return currentScore
	}
}

func (s *RiskScoringService) applyProfileMultipliers(currentScore float64, profile *models.RiskProfile, request *models.RiskAssessmentRequest) float64 {
	adjustedScore := currentScore

	// Apply multipliers based on data factors
	if request.DataFactors != nil {
		for factor, value := range request.DataFactors {
			if multiplier, exists := profile.Multipliers[factor]; exists {
				if s.isFactorActive(value) {
					adjustedScore *= multiplier
				}
			}
		}
	}

	// Apply multipliers based on security factors
	if request.SecurityFactors != nil {
		for factor, value := range request.SecurityFactors {
			if multiplier, exists := profile.Multipliers[factor]; exists {
				if s.isFactorActive(value) {
					adjustedScore *= multiplier
				}
			}
		}
	}

	return adjustedScore
}

func (s *RiskScoringService) isFactorActive(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "yes" || v == "1"
	case int:
		return v > 0
	case float64:
		return v > 0
	default:
		return false
	}
}

func (s *RiskScoringService) applyWeightedScoring(currentScore float64, request *models.RiskAssessmentRequest) float64 {
	// Calculate weighted components
	dataSensitivityScore := s.calculateDataSensitivityScore(request)
	exposureLevelScore := s.calculateExposureLevelScore(request)
	accessControlScore := s.calculateAccessControlScore(request)
	vulnerabilityScore := s.calculateVulnerabilityScore(request)
	complianceScore := s.calculateComplianceScore(request)

	// Apply weights
	weightedScore := (dataSensitivityScore * s.weights.DataSensitivity) +
		(exposureLevelScore * s.weights.ExposureLevel) +
		(accessControlScore * s.weights.AccessControls) +
		(vulnerabilityScore * s.weights.Vulnerabilities) +
		(complianceScore * s.weights.ComplianceStatus)

	// Blend with current score (70% weighted, 30% current)
	return (weightedScore * 0.7) + (currentScore * 0.3)
}

func (s *RiskScoringService) calculateDataSensitivityScore(request *models.RiskAssessmentRequest) float64 {
	score := 0.0

	if request.DataFactors != nil {
		// High sensitivity data
		if s.hasHighRiskData(request.DataFactors) {
			score += 80.0
		} else if s.hasMediumRiskData(request.DataFactors) {
			score += 50.0
		} else if s.hasLowRiskData(request.DataFactors) {
			score += 20.0
		} else {
			score += 30.0 // Default medium sensitivity
		}

		// Additional factors
		if value, exists := request.DataFactors["data_volume"]; exists {
			if volume, ok := value.(string); ok {
				switch volume {
				case "high":
					score += 15.0
				case "medium":
					score += 10.0
				case "low":
					score += 5.0
				}
			}
		}
	}

	return math.Min(score, 100.0)
}

func (s *RiskScoringService) calculateExposureLevelScore(request *models.RiskAssessmentRequest) float64 {
	score := 0.0

	if request.ContextFactors != nil {
		// Endpoint visibility
		if value, exists := request.ContextFactors["endpoint_visibility"]; exists {
			if visibility, ok := value.(string); ok {
				switch visibility {
				case "public":
					score += 70.0
				case "partner":
					score += 50.0
				case "internal":
					score += 30.0
				case "private":
					score += 10.0
				}
			}
		}

		// Network exposure
		if value, exists := request.ContextFactors["network_exposure"]; exists {
			if exposure, ok := value.(string); ok {
				switch exposure {
				case "internet":
					score += 60.0
				case "extranet":
					score += 40.0
				case "intranet":
					score += 20.0
				case "isolated":
					score += 5.0
				}
			}
		}
	}

	return math.Min(score, 100.0)
}

func (s *RiskScoringService) calculateAccessControlScore(request *models.RiskAssessmentRequest) float64 {
	score := 100.0 // Start high, reduce based on controls

	if request.SecurityFactors != nil {
		// Authentication strength
		if value, exists := request.SecurityFactors["authentication_strength"]; exists {
			if auth, ok := value.(string); ok {
				switch auth {
				case "none":
					score -= 0.0 // No reduction, keep high risk
				case "basic":
					score -= 20.0
				case "strong":
					score -= 40.0
				case "multi_factor":
					score -= 60.0
				}
			}
		}

		// Authorization controls
		if value, exists := request.SecurityFactors["authorization_controls"]; exists {
			if controls, ok := value.(string); ok {
				switch controls {
				case "none":
					score -= 0.0
				case "basic":
					score -= 15.0
				case "rbac":
					score -= 30.0
				case "abac":
					score -= 40.0
				}
			}
		}

		// Rate limiting
		if value, exists := request.SecurityFactors["rate_limiting"]; exists {
			if s.isFactorActive(value) {
				score -= 10.0
			}
		}
	}

	return math.Max(score, 0.0)
}

func (s *RiskScoringService) calculateVulnerabilityScore(request *models.RiskAssessmentRequest) float64 {
	score := 0.0

	if request.SecurityFactors != nil {
		// Known vulnerabilities
		if value, exists := request.SecurityFactors["vulnerability_count"]; exists {
			if count, err := s.parseFloat(fmt.Sprintf("%v", value)); err == nil {
				score += math.Min(count*10.0, 80.0) // Cap at 80
			}
		}

		// Vulnerability severity
		if value, exists := request.SecurityFactors["max_vulnerability_severity"]; exists {
			if severity, ok := value.(string); ok {
				switch severity {
				case "critical":
					score += 40.0
				case "high":
					score += 30.0
				case "medium":
					score += 20.0
				case "low":
					score += 10.0
				}
			}
		}

		// Security patches
		if value, exists := request.SecurityFactors["security_patches_current"]; exists {
			if !s.isFactorActive(value) {
				score += 25.0
			}
		}
	}

	return math.Min(score, 100.0)
}

func (s *RiskScoringService) calculateComplianceScore(request *models.RiskAssessmentRequest) float64 {
	score := 0.0

	if request.ContextFactors != nil {
		// Compliance violations
		if value, exists := request.ContextFactors["compliance_violations"]; exists {
			if violations, err := s.parseFloat(fmt.Sprintf("%v", value)); err == nil {
				score += math.Min(violations*15.0, 60.0) // Cap at 60
			}
		}

		// Regulatory requirements
		if value, exists := request.ContextFactors["regulatory_requirements"]; exists {
			if reqs, ok := value.([]string); ok {
				// More regulations = higher compliance risk
				score += math.Min(float64(len(reqs))*5.0, 30.0)
			}
		}

		// Compliance status
		if value, exists := request.ContextFactors["compliance_status"]; exists {
			if status, ok := value.(string); ok {
				switch status {
				case "non_compliant":
					score += 50.0
				case "partially_compliant":
					score += 30.0
				case "compliant":
					score += 10.0
				case "fully_compliant":
					score += 0.0
				}
			}
		}
	}

	return math.Min(score, 100.0)
}

func (s *RiskScoringService) normalizeScore(score float64) float64 {
	// Ensure score is between 0 and 100
	if score < 0 {
		return 0.0
	}
	if score > 100 {
		return 100.0
	}
	return score
}

func (s *RiskScoringService) determineRiskLevel(score float64, profile *models.RiskProfile) models.RiskLevel {
	if score >= profile.Thresholds.Critical {
		return models.RiskLevelCritical
	}
	if score >= profile.Thresholds.High {
		return models.RiskLevelHigh
	}
	if score >= profile.Thresholds.Medium {
		return models.RiskLevelMedium
	}
	return models.RiskLevelLow
}

func (s *RiskScoringService) generateRiskRecommendations(result *models.RiskScoringResult, profile *models.RiskProfile, request *models.RiskAssessmentRequest) []string {
	var recommendations []string

	// Risk level based recommendations
	switch result.RiskLevel {
	case models.RiskLevelCritical:
		recommendations = append(recommendations, "Immediate action required - Critical risk level detected")
		recommendations = append(recommendations, "Consider implementing additional access controls")
		recommendations = append(recommendations, "Enable real-time monitoring and alerting")
		recommendations = append(recommendations, "Conduct security audit")
	case models.RiskLevelHigh:
		recommendations = append(recommendations, "High risk detected - Review security controls")
		recommendations = append(recommendations, "Implement enhanced monitoring")
		recommendations = append(recommendations, "Consider data encryption")
	case models.RiskLevelMedium:
		recommendations = append(recommendations, "Medium risk - Monitor closely")
		recommendations = append(recommendations, "Review access permissions")
	case models.RiskLevelLow:
		recommendations = append(recommendations, "Low risk - Continue standard monitoring")
	}

	// Factor-specific recommendations
	if request.DataFactors != nil {
		if s.hasHighRiskData(request.DataFactors) {
			recommendations = append(recommendations, "Sensitive data detected - Implement data loss prevention")
			recommendations = append(recommendations, "Consider data masking or tokenization")
		}
	}

	if request.SecurityFactors != nil {
		if value, exists := request.SecurityFactors["authentication_strength"]; exists {
			if auth, ok := value.(string); ok && (auth == "none" || auth == "basic") {
				recommendations = append(recommendations, "Strengthen authentication mechanisms")
				recommendations = append(recommendations, "Implement multi-factor authentication")
			}
		}

		if value, exists := request.SecurityFactors["vulnerability_count"]; exists {
			if count, err := s.parseFloat(fmt.Sprintf("%v", value)); err == nil && count > 0 {
				recommendations = append(recommendations, "Address identified vulnerabilities")
				recommendations = append(recommendations, "Implement regular security scanning")
			}
		}
	}

	if request.ContextFactors != nil {
		if value, exists := request.ContextFactors["endpoint_visibility"]; exists {
			if visibility, ok := value.(string); ok && visibility == "public" {
				recommendations = append(recommendations, "Public endpoint detected - Implement rate limiting")
				recommendations = append(recommendations, "Consider API gateway protection")
			}
		}
	}

	return recommendations
}

func (s *RiskScoringService) CreateRiskProfile(ctx context.Context, profile *models.RiskProfile) error {
	// Validate profile
	if err := s.validateRiskProfile(profile); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}

	// Set metadata
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	// Store in repository
	if err := s.piiRepo.CreateRiskProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to create risk profile: %w", err)
	}

	// Add to memory
	s.riskProfiles[profile.ID] = profile

	s.logger.Info("Created risk profile", "profile_id", profile.ID, "name", profile.Name)
	return nil
}

func (s *RiskScoringService) UpdateRiskProfile(ctx context.Context, profileID string, profile *models.RiskProfile) error {
	// Validate profile
	if err := s.validateRiskProfile(profile); err != nil {
		return fmt.Errorf("profile validation failed: %w", err)
	}

	// Update in repository
	profile.UpdatedAt = time.Now()
	if err := s.piiRepo.UpdateRiskProfile(ctx, profileID, profile); err != nil {
		return fmt.Errorf("failed to update risk profile: %w", err)
	}

	// Update in memory
	profile.ID = profileID
	s.riskProfiles[profileID] = profile

	s.logger.Info("Updated risk profile", "profile_id", profileID, "name", profile.Name)
	return nil
}

func (s *RiskScoringService) validateRiskProfile(profile *models.RiskProfile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}

	if profile.BaseScore < 0 || profile.BaseScore > 100 {
		return fmt.Errorf("base score must be between 0 and 100")
	}

	if profile.Thresholds.Low < 0 || profile.Thresholds.Low > 100 {
		return fmt.Errorf("low threshold must be between 0 and 100")
	}

	if profile.Thresholds.Medium < profile.Thresholds.Low || profile.Thresholds.Medium > 100 {
		return fmt.Errorf("medium threshold must be greater than low threshold and less than or equal to 100")
	}

	if profile.Thresholds.High < profile.Thresholds.Medium || profile.Thresholds.High > 100 {
		return fmt.Errorf("high threshold must be greater than medium threshold and less than or equal to 100")
	}

	if profile.Thresholds.Critical < profile.Thresholds.High || profile.Thresholds.Critical > 100 {
		return fmt.Errorf("critical threshold must be greater than high threshold and less than or equal to 100")
	}

	return nil
}

func (s *RiskScoringService) GetRiskProfiles(ctx context.Context, filter *models.RiskProfileFilter) ([]models.RiskProfile, error) {
	var profiles []models.RiskProfile

	for _, profile := range s.riskProfiles {
		if s.matchesProfileFilter(profile, filter) {
			profiles = append(profiles, *profile)
		}
	}

	return profiles, nil
}

func (s *RiskScoringService) matchesProfileFilter(profile *models.RiskProfile, filter *models.RiskProfileFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Category != "" && profile.Category != models.RiskCategory(filter.Category) {
		return false
	}

	if filter.Enabled != nil && profile.Enabled != *filter.Enabled {
		return false
	}

	if filter.Name != "" && !strings.Contains(strings.ToLower(profile.Name), strings.ToLower(filter.Name)) {
		return false
	}

	return true
}

func (s *RiskScoringService) GetRiskTrends(ctx context.Context, filter *models.RiskTrendFilter) (*models.RiskTrendAnalysis, error) {
	// This would typically query the repository for historical risk data
	// For now, we'll return a basic trend analysis structure

	analysis := &models.RiskTrendAnalysis{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		Filter:      filter,
		TimeRange: models.TimeRange{
			Start: time.Now().AddDate(0, -1, 0), // Last month
			End:   time.Now(),
		},
		Trends: []models.RiskTrend{
			{
				Date:         time.Now().AddDate(0, 0, -7),
				AverageScore: 45.2,
				RiskCounts: map[models.RiskLevel]int{
					models.RiskLevelLow:      120,
					models.RiskLevelMedium:   85,
					models.RiskLevelHigh:     32,
					models.RiskLevelCritical: 8,
				},
			},
			{
				Date:         time.Now(),
				AverageScore: 42.8,
				RiskCounts: map[models.RiskLevel]int{
					models.RiskLevelLow:      135,
					models.RiskLevelMedium:   78,
					models.RiskLevelHigh:     28,
					models.RiskLevelCritical: 6,
				},
			},
		},
		Summary: models.RiskTrendSummary{
			TotalAssessments:   492,
			AverageScore:       44.0,
			ScoreChange:        -2.4,
			TrendDirection:     "improving",
			HighestRiskAPI:     "payment-api",
			MostCommonRiskType: "data_exposure",
		},
		Recommendations: []string{
			"Risk scores are trending downward - continue current security measures",
			"Focus on reducing high-risk assessments in payment APIs",
			"Consider implementing additional data protection controls",
		},
	}

	return analysis, nil
}

func (s *RiskScoringService) GenerateRiskReport(ctx context.Context, filter *models.RiskReportFilter) (*models.RiskReport, error) {
	report := &models.RiskReport{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		Filter:      filter,
		Summary: models.RiskReportSummary{
			TotalAssessments: 0,
			RiskDistribution: make(map[models.RiskLevel]int),
			AverageScore:     0.0,
			TopRisks:         []models.TopRisk{},
		},
		Details: []models.RiskAssessmentDetail{},
		Trends:  []models.RiskTrend{},
		Recommendations: []string{
			"Implement continuous risk monitoring",
			"Regular security assessments recommended",
			"Consider automated risk remediation",
		},
	}

	// In a real implementation, this would query the database
	// and populate the report with actual assessment data

	return report, nil
}

func (s *RiskScoringService) publishRiskEvents(ctx context.Context, result *models.RiskScoringResult, request *models.RiskAssessmentRequest) error {
	// Publish risk assessment event
	assessmentEvent := map[string]interface{}{
		"event_type":      "risk_assessment_completed",
		"request_id":      request.RequestID,
		"api_id":          request.APIID,
		"endpoint_id":     request.EndpointID,
		"risk_score":      result.RiskScore,
		"risk_level":      result.RiskLevel,
		"profile_used":    result.ProfileUsed,
		"applied_rules":   result.AppliedRules,
		"recommendations": result.Recommendations,
		"processing_time": result.ProcessingTime.Milliseconds(),
		"ip_address":      request.IPAddress,
		"user_agent":      request.UserAgent,
		"timestamp":       result.CalculatedAt,
	}

	eventJSON, err := json.Marshal(assessmentEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal risk assessment event: %w", err)
	}

	message := kafka.Message{
		Topic: "risk_assessment_events",
		Key:   []byte(request.RequestID),
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce risk assessment event: %w", err)
	}

	// Publish high-risk alert if needed
	if result.RiskLevel == models.RiskLevelHigh || result.RiskLevel == models.RiskLevelCritical {
		alertEvent := map[string]interface{}{
			"event_type":  "high_risk_alert",
			"request_id":  request.RequestID,
			"api_id":      request.APIID,
			"endpoint_id": request.EndpointID,
			"risk_score":  result.RiskScore,
			"risk_level":  result.RiskLevel,
			"alert_level": "high",
			"message":     fmt.Sprintf("High risk detected: %s (Score: %.2f)", result.RiskLevel, result.RiskScore),
			"timestamp":   time.Now(),
		}

		alertJSON, err := json.Marshal(alertEvent)
		if err != nil {
			return fmt.Errorf("failed to marshal high risk alert: %w", err)
		}

		alertMessage := kafka.Message{
			Topic: "security_alerts",
			Key:   []byte(request.RequestID),
			Value: alertJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, alertMessage); err != nil {
			return fmt.Errorf("failed to produce high risk alert: %w", err)
		}
	}

	return nil
}
func (s *RiskScoringService) GetRiskScores(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error) {
	return s.piiRepo.GetRiskAssessments(ctx, filter)
}

func (s *RiskScoringService) GetRiskScore(ctx context.Context, id string) (*models.RiskAssessment, error) {
	// PII repo doesn't have GetRiskAssessment by ID? Let's check.
	// Actually, I can search in GetRiskAssessments.
	assessments, err := s.piiRepo.GetRiskAssessments(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range assessments {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("risk assessment not found: %s", id)
}

func (s *RiskScoringService) CreateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error {
	return s.piiRepo.CreateMitigationPlan(ctx, plan)
}

func (s *RiskScoringService) GetMitigationPlan(ctx context.Context, id string) (*models.MitigationPlan, error) {
	return s.piiRepo.GetMitigationPlan(ctx, id)
}

func (s *RiskScoringService) UpdateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error {
	return s.piiRepo.UpdateMitigationPlan(ctx, plan)
}

func (s *RiskScoringService) DeleteMitigationPlan(ctx context.Context, id string) error {
	return s.piiRepo.DeleteMitigationPlan(ctx, id)
}

func (s *RiskScoringService) GetRiskAssessments(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error) {
	return s.piiRepo.GetRiskAssessments(ctx, filter)
}
