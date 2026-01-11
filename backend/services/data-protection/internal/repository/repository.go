package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"data-protection/internal/models"

	"scopeapi.local/backend/shared/database/postgresql"
)

// =============================================================================
// REPOSITORY INTERFACES
// =============================================================================

type ClassificationRepositoryInterface interface {
	CreateClassificationRule(ctx context.Context, rule *models.ClassificationRule) error
	GetClassificationRule(ctx context.Context, id string) (*models.ClassificationRule, error)
	GetClassificationRules(ctx context.Context, filter *models.ClassificationRuleFilter) ([]models.ClassificationRule, error)
	UpdateClassificationRule(ctx context.Context, id string, rule *models.ClassificationRule) error
	DeleteClassificationRule(ctx context.Context, id string) error
	EnableClassificationRule(ctx context.Context, id string) error
	DisableClassificationRule(ctx context.Context, id string) error
	GetRulesByCategory(ctx context.Context, category models.DataCategory) ([]models.ClassificationRule, error)
	GetActiveRules(ctx context.Context) ([]models.ClassificationRule, error)
	GetRulesByMethod(ctx context.Context, method models.ClassificationMethod) ([]models.ClassificationRule, error)
	CreateClassificationData(ctx context.Context, data *models.ClassificationData) error
}

type PIIRepositoryInterface interface {
	CreatePIIPattern(ctx context.Context, pattern *models.PIIPattern) error
	GetPIIPattern(ctx context.Context, id string) (*models.PIIPattern, error)
	GetPIIPatterns(ctx context.Context, filter *models.PIIPatternFilter) ([]models.PIIPattern, error)
	UpdatePIIPattern(ctx context.Context, pattern *models.PIIPattern) error
	DeletePIIPattern(ctx context.Context, id string) error
	GetPatternsByType(ctx context.Context, piiType string) ([]models.PIIPattern, error)
	GetActivePatterns(ctx context.Context) ([]models.PIIPattern, error)
	StorePIIDetectionResult(ctx context.Context, result *models.PIIDetectionResult) error
	GetPIIDetectionHistory(ctx context.Context, filter *models.PIIHistoryFilter) ([]models.PIIDetectionResult, error)
	GetPIIStatistics(ctx context.Context, timeRange *models.TimeRange) (*models.PIIStatistics, error)
	CreateRiskAssessment(ctx context.Context, assessment *models.RiskAssessment) error
	CreateRiskProfile(ctx context.Context, profile *models.RiskProfile) error
	UpdateRiskProfile(ctx context.Context, profileID string, profile *models.RiskProfile) error
	CreatePIIData(ctx context.Context, data *models.PIIData) error
	CreateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error
	GetMitigationPlan(ctx context.Context, id string) (*models.MitigationPlan, error)
	UpdateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error
	DeleteMitigationPlan(ctx context.Context, id string) error
	GetRiskAssessments(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error)
}

type ComplianceRepositoryInterface interface {
	CreateComplianceFramework(ctx context.Context, framework *models.ComplianceFrameworkData) error
	GetComplianceFramework(ctx context.Context, id string) (*models.ComplianceFrameworkData, error)
	GetComplianceFrameworks(ctx context.Context, filter *models.ComplianceFrameworkFilter) ([]models.ComplianceFrameworkData, error)
	UpdateComplianceFramework(ctx context.Context, framework *models.ComplianceFrameworkData) error
	DeleteComplianceFramework(ctx context.Context, id string) error
	CreateComplianceReport(ctx context.Context, report *models.ComplianceReport) error
	GetComplianceReport(ctx context.Context, id string) (*models.ComplianceReport, error)
	GetComplianceReports(ctx context.Context, filter *models.ComplianceReportFilter) ([]models.ComplianceReport, error)
	UpdateComplianceReport(ctx context.Context, report *models.ComplianceReport) error
	DeleteComplianceReport(ctx context.Context, id string) error
	GetAuditLog(ctx context.Context, filter *models.AuditLogFilter) ([]models.AuditLogEntry, error)
	StoreAuditLog(ctx context.Context, entry *models.AuditLogEntry) error
	GetComplianceStatistics(ctx context.Context, timeRange *models.TimeRange) (*models.ComplianceStatistics, error)
	CreateComplianceValidation(ctx context.Context, validation *models.ComplianceValidation) error
	GetComplianceValidations(ctx context.Context, filter *models.ComplianceValidationFilter) ([]models.ComplianceValidation, error)
	CreateComplianceRule(ctx context.Context, rule *models.ComplianceRule) error
	UpdateComplianceRule(ctx context.Context, ruleID string, rule *models.ComplianceRule) error
	GetRecentViolations(ctx context.Context, since time.Time) ([]models.ComplianceViolation, error)
	GetActiveViolations(ctx context.Context) ([]models.ComplianceViolation, error)
	CreateComplianceViolation(ctx context.Context, violation *models.ComplianceViolation) error
}

// =============================================================================
// IN-MEMORY REPOSITORY IMPLEMENTATIONS
// =============================================================================

type MemoryClassificationRepository struct {
	rules              map[string]*models.ClassificationRule
	classificationData map[string]*models.ClassificationData
	mutex              sync.RWMutex
}

func NewClassificationRepository(db *postgresql.Connection) ClassificationRepositoryInterface {
	repo := &MemoryClassificationRepository{
		rules:              make(map[string]*models.ClassificationRule),
		classificationData: make(map[string]*models.ClassificationData),
	}

	// Load default rules
	repo.loadDefaultRules()
	return repo
}

func (r *MemoryClassificationRepository) loadDefaultRules() {
	defaultRules := []*models.ClassificationRule{
		{
			ID:             "default_public",
			Name:           "Public Data Classification",
			Description:    "Default rule for public data",
			Category:       models.DataCategoryBusiness,
			Classification: models.DataClassificationPublic,
			Method:         models.ClassificationMethodRule,
			Priority:       1,
			Enabled:        true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "default_internal",
			Name:           "Internal Data Classification",
			Description:    "Default rule for internal data",
			Category:       models.DataCategoryBusiness,
			Classification: models.DataClassificationInternal,
			Method:         models.ClassificationMethodRule,
			Priority:       2,
			Enabled:        true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, rule := range defaultRules {
		r.rules[rule.ID] = rule
	}
}

func (r *MemoryClassificationRepository) CreateClassificationRule(ctx context.Context, rule *models.ClassificationRule) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if rule.ID == "" {
		rule.ID = fmt.Sprintf("rule_%d", time.Now().UnixNano())
	}

	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	r.rules[rule.ID] = rule
	return nil
}

func (r *MemoryClassificationRepository) GetClassificationRule(ctx context.Context, id string) (*models.ClassificationRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rule, exists := r.rules[id]
	if !exists {
		return nil, fmt.Errorf("classification rule not found: %s", id)
	}

	return rule, nil
}

func (r *MemoryClassificationRepository) GetClassificationRules(ctx context.Context, filter *models.ClassificationRuleFilter) ([]models.ClassificationRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var rules []models.ClassificationRule
	for _, rule := range r.rules {
		if filter != nil {
			if filter.Category != "" && string(rule.Category) != filter.Category {
				continue
			}
			if filter.Method != "" && string(rule.Method) != filter.Method {
				continue
			}
			if filter.Enabled != nil && rule.Enabled != *filter.Enabled {
				continue
			}
		}
		rules = append(rules, *rule)
	}

	return rules, nil
}

func (r *MemoryClassificationRepository) UpdateClassificationRule(ctx context.Context, id string, rule *models.ClassificationRule) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rules[id]; !exists {
		return fmt.Errorf("classification rule not found: %s", id)
	}

	rule.ID = id
	rule.UpdatedAt = time.Now()
	r.rules[id] = rule
	return nil
}

func (r *MemoryClassificationRepository) DeleteClassificationRule(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rules[id]; !exists {
		return fmt.Errorf("classification rule not found: %s", id)
	}

	delete(r.rules, id)
	return nil
}

func (r *MemoryClassificationRepository) EnableClassificationRule(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	rule, exists := r.rules[id]
	if !exists {
		return fmt.Errorf("classification rule not found: %s", id)
	}

	rule.Enabled = true
	rule.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryClassificationRepository) DisableClassificationRule(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	rule, exists := r.rules[id]
	if !exists {
		return fmt.Errorf("classification rule not found: %s", id)
	}

	rule.Enabled = false
	rule.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryClassificationRepository) GetRulesByCategory(ctx context.Context, category models.DataCategory) ([]models.ClassificationRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var rules []models.ClassificationRule
	for _, rule := range r.rules {
		if rule.Category == category && rule.Enabled {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (r *MemoryClassificationRepository) GetActiveRules(ctx context.Context) ([]models.ClassificationRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var rules []models.ClassificationRule
	for _, rule := range r.rules {
		if rule.Enabled {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (r *MemoryClassificationRepository) GetRulesByMethod(ctx context.Context, method models.ClassificationMethod) ([]models.ClassificationRule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var rules []models.ClassificationRule
	for _, rule := range r.rules {
		if rule.Method == method && rule.Enabled {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (r *MemoryClassificationRepository) CreateClassificationData(ctx context.Context, data *models.ClassificationData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.classificationData[data.ID] = data
	return nil
}

// =============================================================================
// PII REPOSITORY IMPLEMENTATION
// =============================================================================

type MemoryPIIRepository struct {
	patterns        map[string]*models.PIIPattern
	results         map[string]*models.PIIDetectionResult
	piiData         map[string]*models.PIIData
	riskAssessments map[string]*models.RiskAssessment
	riskProfiles    map[string]*models.RiskProfile
	mitigationPlans map[string]*models.MitigationPlan
	mutex           sync.RWMutex
}

func NewPIIRepository(db *postgresql.Connection) PIIRepositoryInterface {
	repo := &MemoryPIIRepository{
		patterns:        make(map[string]*models.PIIPattern),
		results:         make(map[string]*models.PIIDetectionResult),
		piiData:         make(map[string]*models.PIIData),
		riskAssessments: make(map[string]*models.RiskAssessment),
		riskProfiles:    make(map[string]*models.RiskProfile),
		mitigationPlans: make(map[string]*models.MitigationPlan),
	}

	// Load default PII patterns
	repo.loadDefaultPatterns()
	return repo
}

func (r *MemoryPIIRepository) loadDefaultPatterns() {
	defaultPatterns := []*models.PIIPattern{
		{
			ID:         "email_pattern",
			Name:       "Email Address Pattern",
			Pattern:    `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
			Type:       "regex",
			PIIType:    "email",
			Confidence: 0.95,
			Enabled:    true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "ssn_pattern",
			Name:       "Social Security Number Pattern",
			Pattern:    `\b\d{3}-\d{2}-\d{4}\b`,
			Type:       "regex",
			PIIType:    "ssn",
			Confidence: 0.98,
			Enabled:    true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "credit_card_pattern",
			Name:       "Credit Card Number Pattern",
			Pattern:    `\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`,
			Type:       "regex",
			PIIType:    "credit_card",
			Confidence: 0.96,
			Enabled:    true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, pattern := range defaultPatterns {
		r.patterns[pattern.ID] = pattern
	}
}

func (r *MemoryPIIRepository) CreatePIIPattern(ctx context.Context, pattern *models.PIIPattern) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if pattern.ID == "" {
		pattern.ID = fmt.Sprintf("pattern_%d", time.Now().UnixNano())
	}

	pattern.CreatedAt = time.Now()
	pattern.UpdatedAt = time.Now()

	r.patterns[pattern.ID] = pattern
	return nil
}

func (r *MemoryPIIRepository) GetPIIPattern(ctx context.Context, id string) (*models.PIIPattern, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pattern, exists := r.patterns[id]
	if !exists {
		return nil, fmt.Errorf("PII pattern not found: %s", id)
	}

	return pattern, nil
}

func (r *MemoryPIIRepository) GetPIIPatterns(ctx context.Context, filter *models.PIIPatternFilter) ([]models.PIIPattern, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var patterns []models.PIIPattern
	for _, pattern := range r.patterns {
		if filter != nil {
			if filter.PIIType != "" && pattern.PIIType != filter.PIIType {
				continue
			}
			if filter.Type != "" && string(pattern.Type) != filter.Type {
				continue
			}
			if filter.Enabled != nil && pattern.Enabled != *filter.Enabled {
				continue
			}
		}
		patterns = append(patterns, *pattern)
	}

	return patterns, nil
}

func (r *MemoryPIIRepository) UpdatePIIPattern(ctx context.Context, pattern *models.PIIPattern) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.patterns[pattern.ID]; !exists {
		return fmt.Errorf("PII pattern not found: %s", pattern.ID)
	}

	pattern.UpdatedAt = time.Now()
	r.patterns[pattern.ID] = pattern
	return nil
}

func (r *MemoryPIIRepository) DeletePIIPattern(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.patterns[id]; !exists {
		return fmt.Errorf("PII pattern not found: %s", id)
	}

	delete(r.patterns, id)
	return nil
}

func (r *MemoryPIIRepository) GetPatternsByType(ctx context.Context, piiType string) ([]models.PIIPattern, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var patterns []models.PIIPattern
	for _, pattern := range r.patterns {
		if pattern.PIIType == piiType && pattern.Enabled {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns, nil
}

func (r *MemoryPIIRepository) GetActivePatterns(ctx context.Context) ([]models.PIIPattern, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var patterns []models.PIIPattern
	for _, pattern := range r.patterns {
		if pattern.Enabled {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns, nil
}

func (r *MemoryPIIRepository) StorePIIDetectionResult(ctx context.Context, result *models.PIIDetectionResult) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if result.ID == "" {
		result.ID = fmt.Sprintf("result_%d", time.Now().UnixNano())
	}

	result.CreatedAt = time.Now()
	r.results[result.ID] = result
	return nil
}

func (r *MemoryPIIRepository) GetPIIDetectionHistory(ctx context.Context, filter *models.PIIHistoryFilter) ([]models.PIIDetectionResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var results []models.PIIDetectionResult
	for _, result := range r.results {
		if filter != nil {
			if filter.PIIType != "" && result.PIIType != filter.PIIType {
				continue
			}
			if filter.Since != nil && result.CreatedAt.Before(*filter.Since) {
				continue
			}
		}
		results = append(results, *result)
	}

	return results, nil
}

func (r *MemoryPIIRepository) GetPIIStatistics(ctx context.Context, timeRange *models.TimeRange) (*models.PIIStatistics, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := &models.PIIStatistics{
		TotalPIIDetected:  0,
		PIIByType:         make(map[models.PIIType]int),
		PIIBySensitivity:  make(map[models.PIISensitivityLevel]int),
		PIIByAPI:          make(map[string]int),
		PIIByEndpoint:     make(map[string]int),
		DetectionTrends:   []models.PIITrend{},
		ComplianceImpact:  make(map[string]int),
		RiskDistribution:  make(map[string]int),
		ProcessingActions: make(map[string]int),
		GeneratedAt:       time.Now(),
	}

	for _, result := range r.results {
		if timeRange != nil {
			if result.CreatedAt.Before(timeRange.Since) {
				continue
			}
			if result.CreatedAt.After(timeRange.Until) {
				continue
			}
		}

		stats.TotalPIIDetected++
		// Note: PIIType conversion would need to be handled properly
		// For now, using a default type
		stats.PIIByType[models.PIIType("unknown")]++
	}

	return stats, nil
}

func (r *MemoryPIIRepository) CreateRiskAssessment(ctx context.Context, assessment *models.RiskAssessment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.riskAssessments[assessment.ID] = assessment
	return nil
}

func (r *MemoryPIIRepository) CreateRiskProfile(ctx context.Context, profile *models.RiskProfile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.riskProfiles[profile.ID] = profile
	return nil
}

func (r *MemoryPIIRepository) UpdateRiskProfile(ctx context.Context, profileID string, profile *models.RiskProfile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.riskProfiles[profileID] = profile
	return nil
}

func (r *MemoryPIIRepository) CreatePIIData(ctx context.Context, data *models.PIIData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.piiData[data.ID] = data
	return nil
}

func (r *MemoryPIIRepository) CreateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if plan.ID == "" {
		plan.ID = fmt.Sprintf("plan_%d", time.Now().UnixNano())
	}
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	r.mitigationPlans[plan.ID] = plan
	return nil
}

func (r *MemoryPIIRepository) GetMitigationPlan(ctx context.Context, id string) (*models.MitigationPlan, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plan, exists := r.mitigationPlans[id]
	if !exists {
		return nil, fmt.Errorf("mitigation plan not found: %s", id)
	}
	return plan, nil
}

func (r *MemoryPIIRepository) UpdateMitigationPlan(ctx context.Context, plan *models.MitigationPlan) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.mitigationPlans[plan.ID]; !exists {
		return fmt.Errorf("mitigation plan not found: %s", plan.ID)
	}

	plan.UpdatedAt = time.Now()
	r.mitigationPlans[plan.ID] = plan
	return nil
}

func (r *MemoryPIIRepository) DeleteMitigationPlan(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.mitigationPlans[id]; !exists {
		return fmt.Errorf("mitigation plan not found: %s", id)
	}

	delete(r.mitigationPlans, id)
	return nil
}

func (r *MemoryPIIRepository) GetRiskAssessments(ctx context.Context, filter *models.RiskScoreFilter) ([]models.RiskAssessment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var assessments []models.RiskAssessment
	for _, assessment := range r.riskAssessments {
		if filter != nil {
			// DataSource isn't in RiskAssessment yet, but we'll filter what we can
			if filter.RiskLevel != "" && string(assessment.RiskLevel) != filter.RiskLevel {
				continue
			}
			if filter.Since != nil && assessment.AssessedAt.Before(*filter.Since) {
				continue
			}
			if filter.Until != nil && assessment.AssessedAt.After(*filter.Until) {
				continue
			}
		}
		assessments = append(assessments, *assessment)
	}
	return assessments, nil
}

func (r *MemoryComplianceRepository) CreateComplianceValidation(ctx context.Context, validation *models.ComplianceValidation) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.validations[validation.ID] = validation
	return nil
}

func (r *MemoryComplianceRepository) CreateComplianceRule(ctx context.Context, rule *models.ComplianceRule) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.rules[rule.ID] = rule
	return nil
}

func (r *MemoryComplianceRepository) UpdateComplianceRule(ctx context.Context, ruleID string, rule *models.ComplianceRule) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.rules[ruleID] = rule
	return nil
}

// =============================================================================
// COMPLIANCE REPOSITORY IMPLEMENTATION
// =============================================================================

type MemoryComplianceRepository struct {
	frameworks  map[string]*models.ComplianceFrameworkData
	reports     map[string]*models.ComplianceReport
	auditLogs   map[string]*models.AuditLogEntry
	validations map[string]*models.ComplianceValidation
	rules       map[string]*models.ComplianceRule
	violations  map[string]*models.ComplianceViolation
	mutex       sync.RWMutex
}

func NewComplianceRepository(db *postgresql.Connection) ComplianceRepositoryInterface {
	repo := &MemoryComplianceRepository{
		frameworks:  make(map[string]*models.ComplianceFrameworkData),
		reports:     make(map[string]*models.ComplianceReport),
		auditLogs:   make(map[string]*models.AuditLogEntry),
		validations: make(map[string]*models.ComplianceValidation),
		rules:       make(map[string]*models.ComplianceRule),
		violations:  make(map[string]*models.ComplianceViolation),
	}

	// Load default compliance frameworks
	repo.loadDefaultFrameworks()
	return repo
}

func (r *MemoryComplianceRepository) loadDefaultFrameworks() {
	defaultFrameworks := []*models.ComplianceFrameworkData{
		{
			ID:          "gdpr",
			Name:        "General Data Protection Regulation",
			Description: "EU data protection and privacy regulation",
			Version:     "2018",
			Region:      "EU",
			Categories:  []string{"data_protection", "privacy", "consent"},
			Enabled:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "hipaa",
			Name:        "Health Insurance Portability and Accountability Act",
			Description: "US healthcare data protection regulation",
			Version:     "1996",
			Region:      "US",
			Categories:  []string{"healthcare", "data_protection", "privacy"},
			Enabled:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "pci_dss",
			Name:        "Payment Card Industry Data Security Standard",
			Description: "Payment card data security standard",
			Version:     "4.0",
			Region:      "Global",
			Categories:  []string{"payment", "security", "data_protection"},
			Enabled:     true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, framework := range defaultFrameworks {
		r.frameworks[framework.ID] = framework
	}
}

func (r *MemoryComplianceRepository) CreateComplianceFramework(ctx context.Context, framework *models.ComplianceFrameworkData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if framework.ID == "" {
		framework.ID = fmt.Sprintf("framework_%d", time.Now().UnixNano())
	}

	framework.CreatedAt = time.Now()
	framework.UpdatedAt = time.Now()

	r.frameworks[framework.ID] = framework
	return nil
}

func (r *MemoryComplianceRepository) GetComplianceFramework(ctx context.Context, id string) (*models.ComplianceFrameworkData, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	framework, exists := r.frameworks[id]
	if !exists {
		return nil, fmt.Errorf("compliance framework not found: %s", id)
	}

	return framework, nil
}

func (r *MemoryComplianceRepository) GetComplianceFrameworks(ctx context.Context, filter *models.ComplianceFrameworkFilter) ([]models.ComplianceFrameworkData, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var frameworks []models.ComplianceFrameworkData
	for _, framework := range r.frameworks {
		if filter != nil {

			if filter.Region != "" && framework.Region != filter.Region {
				continue
			}
			if filter.Enabled != nil && framework.Enabled != *filter.Enabled {
				continue
			}
		}
		frameworks = append(frameworks, *framework)
	}

	return frameworks, nil
}

func (r *MemoryComplianceRepository) UpdateComplianceFramework(ctx context.Context, framework *models.ComplianceFrameworkData) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.frameworks[framework.ID]; !exists {
		return fmt.Errorf("compliance framework not found: %s", framework.ID)
	}

	framework.UpdatedAt = time.Now()
	r.frameworks[framework.ID] = framework
	return nil
}

func (r *MemoryComplianceRepository) DeleteComplianceFramework(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.frameworks[id]; !exists {
		return fmt.Errorf("compliance framework not found: %s", id)
	}

	delete(r.frameworks, id)
	return nil
}

func (r *MemoryComplianceRepository) CreateComplianceReport(ctx context.Context, report *models.ComplianceReport) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if report.ID == "" {
		report.ID = fmt.Sprintf("report_%d", time.Now().UnixNano())
	}

	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	r.reports[report.ID] = report
	return nil
}

func (r *MemoryComplianceRepository) GetComplianceReport(ctx context.Context, id string) (*models.ComplianceReport, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	report, exists := r.reports[id]
	if !exists {
		return nil, fmt.Errorf("compliance report not found: %s", id)
	}

	return report, nil
}

func (r *MemoryComplianceRepository) GetComplianceReports(ctx context.Context, filter *models.ComplianceReportFilter) ([]models.ComplianceReport, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var reports []models.ComplianceReport
	for _, report := range r.reports {
		if filter != nil {
			if len(filter.Frameworks) > 0 {
				found := false
				for _, framework := range filter.Frameworks {
					if report.Framework == models.ComplianceFramework(framework) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if len(filter.Statuses) > 0 {
				found := false
				for _, status := range filter.Statuses {
					if report.Status == models.ComplianceStatus(status) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if filter.StartDate != nil && report.CreatedAt.Before(*filter.StartDate) {
				continue
			}
			if filter.EndDate != nil && report.CreatedAt.After(*filter.EndDate) {
				continue
			}
		}
		reports = append(reports, *report)
	}

	return reports, nil
}

func (r *MemoryComplianceRepository) UpdateComplianceReport(ctx context.Context, report *models.ComplianceReport) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.reports[report.ID]; !exists {
		return fmt.Errorf("compliance report not found: %s", report.ID)
	}

	report.UpdatedAt = time.Now()
	r.reports[report.ID] = report
	return nil
}

func (r *MemoryComplianceRepository) DeleteComplianceReport(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.reports[id]; !exists {
		return fmt.Errorf("compliance report not found: %s", id)
	}

	delete(r.reports, id)
	return nil
}

func (r *MemoryComplianceRepository) GetAuditLog(ctx context.Context, filter *models.AuditLogFilter) ([]models.AuditLogEntry, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var entries []models.AuditLogEntry
	for _, entry := range r.auditLogs {
		if filter != nil {
			if filter.Action != "" && entry.Action != filter.Action {
				continue
			}
			if filter.Resource != "" && entry.Resource != filter.Resource {
				continue
			}
			if filter.StartDate != nil && entry.Timestamp.Before(*filter.StartDate) {
				continue
			}
			if filter.EndDate != nil && entry.Timestamp.After(*filter.EndDate) {
				continue
			}
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *MemoryComplianceRepository) StoreAuditLog(ctx context.Context, entry *models.AuditLogEntry) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if entry.ID == "" {
		entry.ID = fmt.Sprintf("audit_%d", time.Now().UnixNano())
	}

	entry.Timestamp = time.Now()
	r.auditLogs[entry.ID] = entry
	return nil
}

func (r *MemoryComplianceRepository) GetComplianceStatistics(ctx context.Context, timeRange *models.TimeRange) (*models.ComplianceStatistics, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := &models.ComplianceStatistics{
		TotalFrameworks:     0,
		CompliantFrameworks: 0,
		ViolationCount:      0,
		ComplianceRate:      0.0,
		StatisticsByDate:    make(map[string]interface{}),
		Metadata:            make(map[string]interface{}),
	}

	for _, framework := range r.frameworks {
		if timeRange != nil {
			if framework.CreatedAt.Before(timeRange.Since) {
				continue
			}
			if framework.CreatedAt.After(timeRange.Until) {
				continue
			}
		}

		stats.TotalFrameworks++
		if framework.Enabled {
			stats.CompliantFrameworks++
		}
	}

	// Calculate compliance rate
	if stats.TotalFrameworks > 0 {
		stats.ComplianceRate = float64(stats.CompliantFrameworks) / float64(stats.TotalFrameworks) * 100.0
	}

	return stats, nil
}

// GetComplianceValidations retrieves compliance validations based on filter criteria
func (r *MemoryComplianceRepository) GetComplianceValidations(ctx context.Context, filter *models.ComplianceValidationFilter) ([]models.ComplianceValidation, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var validations []models.ComplianceValidation
	for _, validation := range r.validations {
		if filter != nil {
			// Filter by date range
			if filter.StartDate != nil && validation.ValidatedAt.Before(*filter.StartDate) {
				continue
			}
			if filter.EndDate != nil && validation.ValidatedAt.After(*filter.EndDate) {
				continue
			}

			// Filter by frameworks
			if len(filter.Frameworks) > 0 {
				found := false
				for _, framework := range filter.Frameworks {
					for _, vFramework := range validation.Frameworks {
						if framework == vFramework {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by API IDs
			if len(filter.APIIDs) > 0 {
				found := false
				for _, apiID := range filter.APIIDs {
					if apiID == validation.APIID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by endpoint IDs
			if len(filter.EndpointIDs) > 0 {
				found := false
				for _, endpointID := range filter.EndpointIDs {
					if endpointID == validation.EndpointID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by statuses
			if len(filter.Statuses) > 0 {
				found := false
				for _, status := range filter.Statuses {
					if status == string(validation.OverallStatus) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}

		validations = append(validations, *validation)
	}

	return validations, nil
}

// GetRecentViolations retrieves violations that occurred after a specific time
func (r *MemoryComplianceRepository) GetRecentViolations(ctx context.Context, since time.Time) ([]models.ComplianceViolation, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var violations []models.ComplianceViolation
	for _, violation := range r.violations {
		if violation.DetectedAt.After(since) {
			violations = append(violations, *violation)
		}
	}

	return violations, nil
}

// GetActiveViolations retrieves all violations with status "open" or "in_progress"
func (r *MemoryComplianceRepository) GetActiveViolations(ctx context.Context) ([]models.ComplianceViolation, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var violations []models.ComplianceViolation
	for _, violation := range r.violations {
		if violation.Status == models.ViolationStatusOpen || violation.Status == models.ViolationStatusInProgress {
			violations = append(violations, *violation)
		}
	}

	return violations, nil
}

// CreateComplianceViolation stores a new compliance violation
func (r *MemoryComplianceRepository) CreateComplianceViolation(ctx context.Context, violation *models.ComplianceViolation) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if violation.ID == "" {
		violation.ID = fmt.Sprintf("violation_%d", time.Now().UnixNano())
	}

	violation.CreatedAt = time.Now()
	violation.UpdatedAt = time.Now()

	r.violations[violation.ID] = violation
	return nil
}
