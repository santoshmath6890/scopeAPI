package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"data-protection/internal/models"
	"data-protection/internal/repository"
	"shared/logging"
	"shared/messaging/kafka"
)

type ComplianceServiceInterface interface {
	ValidateCompliance(ctx context.Context, request *models.ComplianceValidationRequest) (*models.ComplianceValidationResult, error)
	CreateComplianceRule(ctx context.Context, rule *models.ComplianceRule) error
	UpdateComplianceRule(ctx context.Context, ruleID string, rule *models.ComplianceRule) error
	GetComplianceRules(ctx context.Context, filter *models.ComplianceRuleFilter) ([]models.ComplianceRule, error)
	GenerateComplianceReport(ctx context.Context, filter *models.ComplianceReportFilter) (*models.ComplianceReport, error)
	GetComplianceStatus(ctx context.Context, filter *models.ComplianceStatusFilter) (*models.ComplianceStatus, error)
	TrackComplianceViolation(ctx context.Context, violation *models.ComplianceViolation) error
}

type ComplianceService struct {
	complianceRepo repository.ComplianceRepositoryInterface
	kafkaProducer  kafka.ProducerInterface
	logger         logging.Logger
	rules          map[string]*models.ComplianceRule
	frameworks     map[string]*models.ComplianceFrameworkData
}


	func NewComplianceService(
	complianceRepo repository.ComplianceRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *ComplianceService {
	service := &ComplianceService{
		complianceRepo: complianceRepo,
		kafkaProducer:  kafkaProducer,
		logger:         logger,
		rules:          make(map[string]*models.ComplianceRule),
		frameworks:     make(map[string]*models.ComplianceFrameworkData),
	}

	// Load default compliance frameworks and rules
	service.loadDefaultFrameworks()
	service.loadDefaultRules()

	return service
}

func (s *ComplianceService) loadDefaultFrameworks() {
	frameworks := []models.ComplianceFrameworkData{
		{
			ID:          "gdpr",
			Name:        "General Data Protection Regulation",
			Description: "EU regulation on data protection and privacy",
			Version:     "2018",
			Region:      "EU",
			Categories: []string{
				"data_protection", "privacy", "consent", "data_subject_rights",
				"data_processing", "data_retention", "data_breach",
			},
			Enabled: true,
		},
		{
			ID:          "ccpa",
			Name:        "California Consumer Privacy Act",
			Description: "California state statute intended to enhance privacy rights",
			Version:     "2020",
			Region:      "California, US",
			Categories: []string{
				"consumer_rights", "data_disclosure", "data_deletion",
				"data_portability", "opt_out_rights",
			},
			Enabled: true,
		},
		{
			ID:          "hipaa",
			Name:        "Health Insurance Portability and Accountability Act",
			Description: "US legislation that provides data privacy and security provisions",
			Version:     "1996",
			Region:      "US",
			Categories: []string{
				"phi_protection", "access_controls", "audit_controls",
				"integrity", "transmission_security",
			},
			Enabled: true,
		},
		{
			ID:          "pci_dss",
			Name:        "Payment Card Industry Data Security Standard",
			Description: "Information security standard for organizations that handle credit cards",
			Version:     "4.0",
			Region:      "Global",
			Categories: []string{
				"network_security", "cardholder_data_protection",
				"vulnerability_management", "access_control",
				"monitoring", "security_policies",
			},
			Enabled: true,
		},
	}

	for _, framework := range frameworks {
		framework.CreatedAt = time.Now()
		framework.UpdatedAt = time.Now()
		s.frameworks[framework.ID] = &framework
	}
}

func (s *ComplianceService) loadDefaultRules() {
	rules := []models.ComplianceRule{
		{
			ID:          "gdpr_pii_encryption",
			Name:        "GDPR PII Encryption Rule",
			Description: "Ensures PII data is encrypted according to GDPR requirements",
			Framework:   "gdpr",
			Category:    "data_protection",
			Severity:    models.ComplianceSeverityHigh,
			Conditions: []models.RuleCondition{
				{
					Field:    "contains_pii",
					Operator: "equals",
					Value:    "true",
				},
				{
					Field:    "encrypted",
					Operator: "equals",
					Value:    "false",
				},
			},
			Actions: []models.RuleAction{
				{
					Type: "violation",
					Config: map[string]interface{}{
						"message":     "PII data must be encrypted",
						"remediation": "Enable encryption for PII data",
					},
				},
				{
					Type: "alert",
					Config: map[string]interface{}{
						"level":      "high",
						"recipients": []string{"dpo@company.com", "security@company.com"},
					},
				},
			},
			Enabled: true,
		},
		{
			ID:          "ccpa_data_retention",
			Name:        "CCPA Data Retention Rule",
			Description: "Validates data retention periods according to CCPA",
			Framework:   "ccpa",
			Category:    "data_retention",
			Severity:    models.ComplianceSeverityMedium,
			Conditions: []models.RuleCondition{
				{
					Field:    "data_age_days",
					Operator: "greater_than",
					Value:    "365",
				},
				{
					Field:    "retention_policy_applied",
					Operator: "equals",
					Value:    "false",
				},
			},
			Actions: []models.RuleAction{
				{
					Type: "violation",
					Config: map[string]interface{}{
						"message":     "Data retention period exceeded",
						"remediation": "Apply data retention policy or delete data",
					},
				},
			},
			Enabled: true,
		},
		{
			ID:          "hipaa_phi_access_control",
			Name:        "HIPAA PHI Access Control Rule",
			Description: "Ensures proper access controls for PHI data",
			Framework:   "hipaa",
			Category:    "access_control",
			Severity:    models.ComplianceSeverityCritical,
			Conditions: []models.RuleCondition{
				{
					Field:    "contains_phi",
					Operator: "equals",
					Value:    "true",
				},
				{
					Field:    "access_control_level",
					Operator: "in",
					Value:    "none,basic",
				},
			},
			Actions: []models.RuleAction{
				{
					Type: "violation",
					Config: map[string]interface{}{
						"message":     "PHI requires strong access controls",
						"remediation": "Implement role-based access control for PHI",
					},
				},
				{
					Type: "block",
					Config: map[string]interface{}{
						"reason": "Insufficient access controls for PHI data",
					},
				},
			},
			Enabled: true,
		},
		{
			ID:          "pci_cardholder_data_encryption",
			Name:        "PCI DSS Cardholder Data Encryption",
			Description: "Ensures cardholder data is properly encrypted",
			Framework:   "pci_dss",
			Category:    "cardholder_data_protection",
			Severity:    models.ComplianceSeverityCritical,
			Conditions: []models.RuleCondition{
				{
					Field:    "contains_cardholder_data",
					Operator: "equals",
					Value:    "true",
				},
				{
					Field:    "encryption_algorithm",
					Operator: "not_in",
					Value:    "AES-256,AES-192",
				},
			},
			Actions: []models.RuleAction{
				{
					Type: "violation",
					Config: map[string]interface{}{
						"message":     "Cardholder data must use approved encryption",
						"remediation": "Use AES-256 or AES-192 encryption for cardholder data",
					},
				},
			},
			Enabled: true,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		s.rules[rule.ID] = &rule
	}
}

func (s *ComplianceService) ValidateCompliance(ctx context.Context, request *models.ComplianceValidationRequest) (*models.ComplianceValidationResult, error) {
	startTime := time.Now()

	result := &models.ComplianceValidationResult{
		RequestID:       request.RequestID,
		OverallStatus:   models.ComplianceStatusCompliant,
		Violations:      []models.ComplianceViolation{},
		Warnings:        []models.ComplianceWarning{},
		FrameworkResults: make(map[string]models.FrameworkComplianceResult),
		ProcessingTime:  0,
		ValidatedAt:     time.Now(),
	}

	// Validate against specified frameworks or all if none specified
	frameworks := request.Frameworks
	if len(frameworks) == 0 {
		for frameworkID := range s.frameworks {
			frameworks = append(frameworks, frameworkID)
		}
	}

	// Process each framework
	for _, frameworkID := range frameworks {
		frameworkResult := s.validateFrameworkCompliance(ctx, frameworkID, request)
		result.FrameworkResults[frameworkID] = frameworkResult

		// Aggregate violations and warnings
		result.Violations = append(result.Violations, frameworkResult.Violations...)
		result.Warnings = append(result.Warnings, frameworkResult.Warnings...)

		// Update overall status
		if frameworkResult.Status == models.ComplianceStatusNonCompliant {
			result.OverallStatus = models.ComplianceStatusNonCompliant
		} else if frameworkResult.Status == models.ComplianceStatusPartiallyCompliant && 
			result.OverallStatus == models.ComplianceStatusCompliant {
			result.OverallStatus = models.ComplianceStatusPartiallyCompliant
		}
	}

	// Calculate processing time
	result.ProcessingTime = time.Since(startTime)

	// Store compliance validation
	validation := &models.ComplianceValidation{
		ID:               uuid.New().String(),
		RequestID:        request.RequestID,
		APIID:            request.APIID,
		EndpointID:       request.EndpointID,
		Frameworks:       frameworks,
		OverallStatus:    result.OverallStatus,
		ViolationCount:   len(result.Violations),
		WarningCount:     len(result.Warnings),
		FrameworkResults: result.FrameworkResults,
		IPAddress:        request.IPAddress,
		UserAgent:        request.UserAgent,
		ValidatedAt:      time.Now(),
		Metadata: map[string]interface{}{
			"processing_time": result.ProcessingTime.Milliseconds(),
			"data_factors":    request.DataFactors,
		},
	}

	if err := s.complianceRepo.CreateComplianceValidation(ctx, validation); err != nil {
		s.logger.Error("Failed to store compliance validation", "error", err)
	}

	// Publish compliance events
	if err := s.publishComplianceEvents(ctx, result, request); err != nil {
		s.logger.Error("Failed to publish compliance events", "error", err)
	}

	return result, nil
}

func (s *ComplianceService) validateFrameworkCompliance(ctx context.Context, frameworkID string, request *models.ComplianceValidationRequest) models.FrameworkComplianceResult {
	framework, exists := s.frameworks[frameworkID]
	if !exists {
		return models.FrameworkComplianceResult{
			FrameworkID:   frameworkID,
			FrameworkName: "Unknown",
			Status:        models.ComplianceStatusNonCompliant,
			Violations: []models.ComplianceViolation{
				{
					ID:          uuid.New().String(),
					RuleID:      "",
					Framework:   frameworkID,
					Severity:    models.ComplianceSeverityHigh,
					Message:     "Unknown compliance framework",
					Description: fmt.Sprintf("Framework %s is not recognized", frameworkID),
					DetectedAt:  time.Now(),
				},
			},
			Warnings:       []models.ComplianceWarning{},
			Score:          0.0,
			RequirementsMet: 0,
			TotalRequirements: 0,
		}
	}

	result := models.FrameworkComplianceResult{
		FrameworkID:       frameworkID,
		FrameworkName:     framework.Name,
		Status:            models.ComplianceStatusCompliant,
		Violations:        []models.ComplianceViolation{},
		Warnings:          []models.ComplianceWarning{},
		Score:             100.0,
		RequirementsMet:   0,
		TotalRequirements: len(framework.Requirements),
	}

	// Evaluate framework-specific rules
	for _, rule := range s.getFrameworkRules(frameworkID) {
		if !rule.Enabled {
			continue
		}

		if s.evaluateComplianceRule(rule, request) {
			// Rule conditions met - this is a violation
			violation := models.ComplianceViolation{
				ID:          uuid.New().String(),
				RuleID:      rule.ID,
				RuleName:    rule.Name,
				Framework:   frameworkID,
				Category:    rule.Category,
				Severity:    rule.Severity,
				Message:     s.getViolationMessage(rule),
				Description: rule.Description,
				Remediation: s.getRemediationAdvice(rule),
				APIID:       request.APIID,
				EndpointID:  request.EndpointID,
				DetectedAt:  time.Now(),
				Status:      models.ViolationStatusOpen,
				Metadata: map[string]interface{}{
					"rule_conditions": rule.Conditions,
					"data_factors":    request.DataFactors,
				},
			}

			result.Violations = append(result.Violations, violation)

			// Adjust score based on severity
			switch rule.Severity {
			case models.ComplianceSeverityCritical:
				result.Score -= 25.0
			case models.ComplianceSeverityHigh:
				result.Score -= 15.0
			case models.ComplianceSeverityMedium:
				result.Score -= 10.0
			case models.ComplianceSeverityLow:
				result.Score -= 5.0
			}

			// Execute rule actions
			s.executeComplianceActions(ctx, rule, violation, request)
		} else {
			result.RequirementsMet++
		}
	}

	// Determine overall framework status
	if len(result.Violations) == 0 {
		result.Status = models.ComplianceStatusCompliant
	} else {
		criticalViolations := s.countViolationsBySeverity(result.Violations, models.ComplianceSeverityCritical)
		highViolations := s.countViolationsBySeverity(result.Violations, models.ComplianceSeverityHigh)

		if criticalViolations > 0 || highViolations > 2 {
			result.Status = models.ComplianceStatusNonCompliant
		} else {
			result.Status = models.ComplianceStatusPartiallyCompliant
		}
	}

	// Ensure score doesn't go below 0
	if result.Score < 0 {
		result.Score = 0
	}

	return result
}

func (s *ComplianceService) getFrameworkRules(frameworkID string) []*models.ComplianceRule {
	var rules []*models.ComplianceRule
	for _, rule := range s.rules {
		if rule.Framework == frameworkID {
			rules = append(rules, rule)
		}
	}
	return rules
}

func (s *ComplianceService) evaluateComplianceRule(rule *models.ComplianceRule, request *models.ComplianceValidationRequest) bool {
	for _, condition := range rule.Conditions {
		if !s.evaluateComplianceCondition(condition, request) {
			return false
		}
	}
	return true
}

func (s *ComplianceService) evaluateComplianceCondition(condition models.ComplianceCondition, request *models.ComplianceValidationRequest) bool {
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

func (s *ComplianceService) evaluateConditionValue(value interface{}, operator, expectedValue string) bool {
	valueStr := fmt.Sprintf("%v", value)

	switch operator {
	case "equals":
		return valueStr == expectedValue
	case "not_equals":
		return valueStr != expectedValue
	case "contains":
		return strings.Contains(strings.ToLower(valueStr), strings.ToLower(expectedValue))
	case "not_contains":
		return !strings.Contains(strings.ToLower(valueStr), strings.ToLower(expectedValue))
	case "in":
		values := strings.Split(expectedValue, ",")
		for _, v := range values {
			if strings.TrimSpace(v) == valueStr {
				return true
			}
		}
		return false
	case "not_in":
		values := strings.Split(expectedValue, ",")
		for _, v := range values {
			if strings.TrimSpace(v) == valueStr {
				return false
			}
		}
		return true
	case "greater_than":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValue); err == nil {
				return numValue > expectedNum
			}
		}
	case "less_than":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValue); err == nil {
				return numValue < expectedNum
			}
		}
	case "greater_equal":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValue); err == nil {
				return numValue >= expectedNum
			}
		}
	case "less_equal":
		if numValue, err := s.parseFloat(valueStr); err == nil {
			if expectedNum, err := s.parseFloat(expectedValue); err == nil {
				return numValue <= expectedNum
			}
		}
	case "regex":
		matched, err := regexp.MatchString(expectedValue, valueStr)
		return err == nil && matched
	}

	return false
}

func (s *ComplianceService) parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func (s *ComplianceService) getViolationMessage(rule *models.ComplianceRule) string {
	for _, action := range rule.Actions {
		if action.Type == "violation" {
			if message, exists := action.Config["message"]; exists {
				if msgStr, ok := message.(string); ok {
					return msgStr
				}
			}
		}
	}
	return fmt.Sprintf("Compliance rule violation: %s", rule.Name)
}

func (s *ComplianceService) getRemediationAdvice(rule *models.ComplianceRule) string {
	for _, action := range rule.Actions {
		if action.Type == "violation" {
			if remediation, exists := action.Config["remediation"]; exists {
				if remStr, ok := remediation.(string); ok {
					return remStr
				}
			}
		}
	}
	return "Review compliance requirements and implement necessary controls"
}

func (s *ComplianceService) executeComplianceActions(ctx context.Context, rule *models.ComplianceRule, violation models.ComplianceViolation, request *models.ComplianceValidationRequest) {
	for _, action := range rule.Actions {
		switch action.Type {
		case "alert":
			s.sendComplianceAlert(ctx, action, violation, request)
		case "block":
			s.blockRequest(ctx, action, violation, request)
		case "log":
			s.logComplianceEvent(ctx, action, violation, request)
		case "webhook":
			s.callWebhook(ctx, action, violation, request)
		}
	}
}

func (s *ComplianceService) sendComplianceAlert(ctx context.Context, action models.ComplianceAction, violation models.ComplianceViolation, request *models.ComplianceValidationRequest) {
	alert := map[string]interface{}{
		"event_type":    "compliance_violation_alert",
		"violation_id":  violation.ID,
		"rule_id":       violation.RuleID,
		"framework":     violation.Framework,
		"severity":      violation.Severity,
		"message":       violation.Message,
		"api_id":        request.APIID,
		"endpoint_id":   request.EndpointID,
		"timestamp":     time.Now(),
	}

	if level, exists := action.Config["level"]; exists {
		alert["alert_level"] = level
	}

	if recipients, exists := action.Config["recipients"]; exists {
		alert["recipients"] = recipients
	}

	alertJSON, err := json.Marshal(alert)
	if err != nil {
		s.logger.Error("Failed to marshal compliance alert", "error", err)
		return
	}

	message := kafka.Message{
		Topic: "compliance_alerts",
		Key:   violation.ID,
		Value: alertJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to produce compliance alert", "error", err)
	}
}

func (s *ComplianceService) blockRequest(ctx context.Context, action models.ComplianceAction, violation models.ComplianceViolation, request *models.ComplianceValidationRequest) {
	blockEvent := map[string]interface{}{
		"event_type":    "compliance_block_request",
		"violation_id":  violation.ID,
		"rule_id":       violation.RuleID,
		"framework":     violation.Framework,
		"api_id":        request.APIID,
		"endpoint_id":   request.EndpointID,
		"ip_address":    request.IPAddress,
		"user_agent":    request.UserAgent,
		"timestamp":     time.Now(),
	}

	if reason, exists := action.Config["reason"]; exists {
		blockEvent["reason"] = reason
	}

	blockJSON, err := json.Marshal(blockEvent)
	if err != nil {
		s.logger.Error("Failed to marshal compliance block event", "error", err)
		return
	}

	message := kafka.Message{
		Topic: "attack_blocking_events",
		Key:   violation.ID,
		Value: blockJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to produce compliance block event", "error", err)
	}
}

func (s *ComplianceService) logComplianceEvent(ctx context.Context, action models.ComplianceAction, violation models.ComplianceViolation, request *models.ComplianceValidationRequest) {
	logLevel := "info"
	if level, exists := action.Config["level"]; exists {
		if levelStr, ok := level.(string); ok {
			logLevel = levelStr
		}
	}

	logData := map[string]interface{}{
		"violation_id": violation.ID,
		"rule_id":      violation.RuleID,
		"framework":    violation.Framework,
		"severity":     violation.Severity,
		"message":      violation.Message,
		"api_id":       request.APIID,
		"endpoint_id":  request.EndpointID,
	}

	switch logLevel {
	case "error":
		s.logger.Error("Compliance violation", logData)
	case "warn":
		s.logger.Warn("Compliance violation", logData)
	case "debug":
		s.logger.Debug("Compliance violation", logData)
	default:
		s.logger.Info("Compliance violation", logData)
	}
}

func (s *ComplianceService) callWebhook(ctx context.Context, action models.ComplianceAction, violation models.ComplianceViolation, request *models.ComplianceValidationRequest) {
	webhookEvent := map[string]interface{}{
		"event_type":    "compliance_webhook",
		"violation_id":  violation.ID,
		"rule_id":       violation.RuleID,
		"framework":     violation.Framework,
		"severity":      violation.Severity,
		"message":       violation.Message,
		"api_id":        request.APIID,
		"endpoint_id":   request.EndpointID,
		"timestamp":     time.Now(),
	}

	webhookJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		s.logger.Error("Failed to marshal webhook event", "error", err)
		return
	}

	message := kafka.Message{
		Topic: "webhook_events",
		Key:   violation.ID,
		Value: webhookJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to produce webhook event", "error", err)
	}
}

func (s *ComplianceService) countViolationsBySeverity(violations []models.ComplianceViolation, severity models.ComplianceSeverity) int {
	count := 0
	for _, violation := range violations {
		if violation.Severity == severity {
			count++
		}
	}
	return count
}

func (s *ComplianceService) CreateComplianceRule(ctx context.Context, rule *models.ComplianceRule) error {
	// Validate rule
	if err := s.validateComplianceRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	// Set metadata
	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	// Store in repository
	if err := s.complianceRepo.CreateComplianceRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to create compliance rule: %w", err)
	}

	// Add to memory
	s.rules[rule.ID] = rule

	s.logger.Info("Created compliance rule", "rule_id", rule.ID, "name", rule.Name, "framework", rule.Framework)
	return nil
}

func (s *ComplianceService) UpdateComplianceRule(ctx context.Context, ruleID string, rule *models.ComplianceRule) error {
	// Validate rule
	if err := s.validateComplianceRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	// Update in repository
	rule.UpdatedAt = time.Now()
	if err := s.complianceRepo.UpdateComplianceRule(ctx, ruleID, rule); err != nil {
		return fmt.Errorf("failed to update compliance rule: %w", err)
	}

	// Update in memory
	rule.ID = ruleID
	s.rules[ruleID] = rule

	s.logger.Info("Updated compliance rule", "rule_id", ruleID, "name", rule.Name, "framework", rule.Framework)
	return nil
}

func (s *ComplianceService) validateComplianceRule(rule *models.ComplianceRule) error {
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if rule.Framework == "" {
		return fmt.Errorf("framework is required")
	}

	if _, exists := s.frameworks[rule.Framework]; !exists {
		return fmt.Errorf("unknown framework: %s", rule.Framework)
	}

	if len(rule.Conditions) == 0 {
		return fmt.Errorf("at least one condition is required")
	}

	for i, condition := range rule.Conditions {
		if condition.Field == "" {
			return fmt.Errorf("condition %d: field is required", i)
		}
		if condition.Operator == "" {
			return fmt.Errorf("condition %d: operator is required", i)
		}
		if condition.Value == "" {
			return fmt.Errorf("condition %d: value is required", i)
		}
	}

	if len(rule.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	return nil
}

func (s *ComplianceService) GetComplianceRules(ctx context.Context, filter *models.ComplianceRuleFilter) ([]models.ComplianceRule, error) {
	var rules []models.ComplianceRule

	for _, rule := range s.rules {
		if s.matchesRuleFilter(rule, filter) {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (s *ComplianceService) matchesRuleFilter(rule *models.ComplianceRule, filter *models.ComplianceRuleFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Framework != "" && rule.Framework != filter.Framework {
		return false
	}

	if filter.Category != "" && rule.Category != filter.Category {
		return false
	}

	if filter.Severity != "" && rule.Severity != models.ComplianceSeverity(filter.Severity) {
		return false
	}

	if filter.Enabled != nil && rule.Enabled != *filter.Enabled {
		return false
	}

	if filter.Name != "" && !strings.Contains(strings.ToLower(rule.Name), strings.ToLower(filter.Name)) {
		return false
	}

	return true
}

func (s *ComplianceService) GenerateComplianceReport(ctx context.Context, filter *models.ComplianceReportFilter) (*models.ComplianceReport, error) {
	report := &models.ComplianceReport{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		Filter:      filter,
		Summary: models.ComplianceReportSummary{
			TotalValidations:     0,
			ComplianceRate:       0.0,
			ViolationsByFramework: make(map[string]int),
			ViolationsBySeverity:  make(map[models.ComplianceSeverity]int),
			TopViolations:        []models.TopViolation{},
		},
		FrameworkReports: make(map[string]models.FrameworkReport),
		Trends:          []models.ComplianceTrend{},
		Recommendations: []string{},
	}

	// Get compliance data from repository
	validations, err := s.complianceRepo.GetComplianceValidations(ctx, &models.ComplianceValidationFilter{
		StartDate:  filter.StartDate,
		EndDate:    filter.EndDate,
		Frameworks: filter.Frameworks,
		APIIDs:     filter.APIIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get compliance validations: %w", err)
	}

	// Process validations
	report.Summary.TotalValidations = len(validations)
	compliantCount := 0

	for _, validation := range validations {
		if validation.OverallStatus == models.ComplianceStatusCompliant {
			compliantCount++
		}

		// Count violations by framework
		for framework, result := range validation.FrameworkResults {
			report.Summary.ViolationsByFramework[framework] += len(result.Violations)
		}

		// Count violations by severity
		for _, result := range validation.FrameworkResults {
			for _, violation := range result.Violations {
				report.Summary.ViolationsBySeverity[violation.Severity]++
			}
		}
	}

	// Calculate compliance rate
	if report.Summary.TotalValidations > 0 {
		report.Summary.ComplianceRate = float64(compliantCount) / float64(report.Summary.TotalValidations) * 100
	}

	// Generate framework-specific reports
	for _, frameworkID := range filter.Frameworks {
		frameworkReport := s.generateFrameworkReport(ctx, frameworkID, validations)
		report.FrameworkReports[frameworkID] = frameworkReport
	}

	// Generate trends
	report.Trends = s.generateComplianceTrends(ctx, validations, filter)

	// Generate recommendations
	report.Recommendations = s.generateComplianceRecommendations(ctx, report)

	return report, nil
}

func (s *ComplianceService) generateFrameworkReport(ctx context.Context, frameworkID string, validations []models.ComplianceValidation) models.FrameworkReport {
	framework := s.frameworks[frameworkID]
	
	report := models.FrameworkReport{
		FrameworkID:       frameworkID,
		FrameworkName:     framework.Name,
		TotalValidations:  0,
		CompliantCount:    0,
		ViolationCount:    0,
		ComplianceRate:    0.0,
		AverageScore:      0.0,
		RequirementStatus: make(map[string]models.RequirementStatus),
		TopViolations:     []models.TopViolation{},
	}

	var totalScore float64
	violationCounts := make(map[string]int)

	for _, validation := range validations {
		if result, exists := validation.FrameworkResults[frameworkID]; exists {
			report.TotalValidations++
			totalScore += result.Score

			if result.Status == models.ComplianceStatusCompliant {
				report.CompliantCount++
			}

			report.ViolationCount += len(result.Violations)

			// Count violations by rule
			for _, violation := range result.Violations {
				violationCounts[violation.RuleID]++
			}
		}
	}

	// Calculate averages
	if report.TotalValidations > 0 {
		report.ComplianceRate = float64(report.CompliantCount) / float64(report.TotalValidations) * 100
		report.AverageScore = totalScore / float64(report.TotalValidations)
	}

	// Generate top violations
	for ruleID, count := range violationCounts {
		if rule, exists := s.rules[ruleID]; exists {
			report.TopViolations = append(report.TopViolations, models.TopViolation{
				RuleID:      ruleID,
				RuleName:    rule.Name,
				Count:       count,
				Severity:    rule.Severity,
				Category:    rule.Category,
				Percentage:  float64(count) / float64(report.TotalValidations) * 100,
			})
		}
	}

	// Sort top violations by count
	sort.Slice(report.TopViolations, func(i, j int) bool {
		return report.TopViolations[i].Count > report.TopViolations[j].Count
	})

	// Limit to top 10
	if len(report.TopViolations) > 10 {
		report.TopViolations = report.TopViolations[:10]
	}

	return report
}

func (s *ComplianceService) generateComplianceTrends(ctx context.Context, validations []models.ComplianceValidation, filter *models.ComplianceReportFilter) []models.ComplianceTrend {
	// Group validations by day
	dailyData := make(map[string]*models.ComplianceTrend)

	for _, validation := range validations {
		dateKey := validation.ValidatedAt.Format("2006-01-02")
		
		if trend, exists := dailyData[dateKey]; exists {
			trend.TotalValidations++
			if validation.OverallStatus == models.ComplianceStatusCompliant {
				trend.CompliantCount++
			}
			trend.ViolationCount += validation.ViolationCount
		} else {
			compliantCount := 0
			if validation.OverallStatus == models.ComplianceStatusCompliant {
				compliantCount = 1
			}

			dailyData[dateKey] = &models.ComplianceTrend{
				Date:             validation.ValidatedAt,
				TotalValidations: 1,
				CompliantCount:   compliantCount,
				ViolationCount:   validation.ViolationCount,
				ComplianceRate:   0.0, // Will be calculated below
			}
		}
	}

	// Convert to slice and calculate compliance rates
	var trends []models.ComplianceTrend
	for _, trend := range dailyData {
		if trend.TotalValidations > 0 {
			trend.ComplianceRate = float64(trend.CompliantCount) / float64(trend.TotalValidations) * 100
		}
		trends = append(trends, *trend)
	}

	// Sort by date
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Date.Before(trends[j].Date)
	})

	return trends
}

func (s *ComplianceService) generateComplianceRecommendations(ctx context.Context, report *models.ComplianceReport) []string {
	var recommendations []string

	// Overall compliance rate recommendations
	if report.Summary.ComplianceRate < 80 {
		recommendations = append(recommendations, "Compliance rate is below 80% - immediate attention required")
		recommendations = append(recommendations, "Review and strengthen compliance controls")
	} else if report.Summary.ComplianceRate < 95 {
		recommendations = append(recommendations, "Good compliance rate - continue monitoring and improvement")
	} else {
		recommendations = append(recommendations, "Excellent compliance rate - maintain current practices")
	}

	// Severity-based recommendations
	if criticalCount, exists := report.Summary.ViolationsBySeverity[models.ComplianceSeverityCritical]; exists && criticalCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d critical compliance violations immediately", criticalCount))
	}

	if highCount, exists := report.Summary.ViolationsBySeverity[models.ComplianceSeverityHigh]; exists && highCount > 5 {
		recommendations = append(recommendations, fmt.Sprintf("High number of high-severity violations (%d) - prioritize remediation", highCount))
	}

	// Framework-specific recommendations
	for frameworkID, violationCount := range report.Summary.ViolationsByFramework {
		if violationCount > 10 {
			if framework, exists := s.frameworks[frameworkID]; exists {
				recommendations = append(recommendations, fmt.Sprintf("High violation count for %s (%d) - review framework implementation", framework.Name, violationCount))
			}
		}
	}

	// General recommendations
	recommendations = append(recommendations, "Implement automated compliance monitoring")
	recommendations = append(recommendations, "Regular compliance training for development teams")
	recommendations = append(recommendations, "Consider compliance-as-code practices")

	return recommendations
}

func (s *ComplianceService) GetComplianceStatus(ctx context.Context, filter *models.ComplianceStatusFilter) (*models.ComplianceStatus, error) {
	status := &models.ComplianceStatus{
		OverallStatus:     models.ComplianceStatusCompliant,
		FrameworkStatuses: make(map[string]models.FrameworkStatus),
		LastUpdated:       time.Now(),
		Summary: models.ComplianceStatusSummary{
			TotalFrameworks:   len(s.frameworks),
			CompliantFrameworks: 0,
			ActiveViolations:  0,
			RecentViolations:  0,
		},
	}

	// Get recent violations
	recentViolations, err := s.complianceRepo.GetRecentViolations(ctx, time.Now().AddDate(0, 0, -7))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent violations: %w", err)
	}

	status.Summary.RecentViolations = len(recentViolations)

	// Get active violations
	activeViolations, err := s.complianceRepo.GetActiveViolations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active violations: %w", err)
	}

	status.Summary.ActiveViolations = len(activeViolations)

	// Process each framework
	for frameworkID, framework := range s.frameworks {
		if filter != nil && len(filter.Frameworks) > 0 {
			found := false
			for _, f := range filter.Frameworks {
				if f == frameworkID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		frameworkStatus := s.getFrameworkStatus(ctx, frameworkID, activeViolations)
		status.FrameworkStatuses[frameworkID] = frameworkStatus

		if frameworkStatus.Status == models.ComplianceStatusCompliant {
			status.Summary.CompliantFrameworks++
		} else if frameworkStatus.Status == models.ComplianceStatusNonCompliant {
			status.OverallStatus = models.ComplianceStatusNonCompliant
		} else if status.OverallStatus == models.ComplianceStatusCompliant {
			status.OverallStatus = models.ComplianceStatusPartiallyCompliant
		}
	}

	return status, nil
}

func (s *ComplianceService) getFrameworkStatus(ctx context.Context, frameworkID string, activeViolations []models.ComplianceViolation) models.FrameworkStatus {
	framework := s.frameworks[frameworkID]
	
	status := models.FrameworkStatus{
		FrameworkID:      frameworkID,
		FrameworkName:    framework.Name,
		Status:           models.ComplianceStatusCompliant,
		ActiveViolations: 0,
		LastAssessment:   time.Now(),
		Score:            100.0,
	}

	// Count violations for this framework
	criticalCount := 0
	highCount := 0
	
	for _, violation := range activeViolations {
		if violation.Framework == frameworkID {
			status.ActiveViolations++
			
			switch violation.Severity {
			case models.ComplianceSeverityCritical:
				criticalCount++
				status.Score -= 25.0
			case models.ComplianceSeverityHigh:
				highCount++
				status.Score -= 15.0
			case models.ComplianceSeverityMedium:
				status.Score -= 10.0
			case models.ComplianceSeverityLow:
				status.Score -= 5.0
			}
		}
	}

	// Determine status based on violations
	if criticalCount > 0 || highCount > 2 {
		status.Status = models.ComplianceStatusNonCompliant
	} else if status.ActiveViolations > 0 {
		status.Status = models.ComplianceStatusPartiallyCompliant
	}

	// Ensure score doesn't go below 0
	if status.Score < 0 {
		status.Score = 0
	}

	return status
}

func (s *ComplianceService) TrackComplianceViolation(ctx context.Context, violation *models.ComplianceViolation) error {
	// Set metadata
	violation.ID = uuid.New().String()
	violation.DetectedAt = time.Now()
	violation.Status = models.ViolationStatusOpen

	// Store in repository
	if err := s.complianceRepo.CreateComplianceViolation(ctx, violation); err != nil {
		return fmt.Errorf("failed to track compliance violation: %w", err)
	}

	// Publish violation event
	violationEvent := map[string]interface{}{
		"event_type":    "compliance_violation_tracked",
		"violation_id":  violation.ID,
		"rule_id":       violation.RuleID,
		"framework":     violation.Framework,
		"severity":      violation.Severity,
		"message":       violation.Message,
		"api_id":        violation.APIID,
		"endpoint_id":   violation.EndpointID,
		"detected_at":   violation.DetectedAt,
		"status":        violation.Status,
	}

	eventJSON, err := json.Marshal(violationEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal violation event: %w", err)
	}

	message := kafka.Message{
		Topic: "compliance_violations",
		Key:   violation.ID,
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce violation event: %w", err)
	}

	s.logger.Info("Tracked compliance violation", 
		"violation_id", violation.ID, 
		"rule_id", violation.RuleID, 
		"framework", violation.Framework,
		"severity", violation.Severity)

	return nil
}

func (s *ComplianceService) publishComplianceEvents(ctx context.Context, result *models.ComplianceValidationResult, request *models.ComplianceValidationRequest) error {
	// Publish validation completed event
	validationEvent := map[string]interface{}{
		"event_type":        "compliance_validation_completed",
		"request_id":        request.RequestID,
		"api_id":            request.APIID,
		"endpoint_id":       request.EndpointID,
		"overall_status":    result.OverallStatus,
		"violation_count":   len(result.Violations),
		"warning_count":     len(result.Warnings),
		"frameworks":        request.Frameworks,
		"processing_time":   result.ProcessingTime.Milliseconds(),
		"validated_at":      result.ValidatedAt,
		"ip_address":        request.IPAddress,
		"user_agent":        request.UserAgent,
	}

	eventJSON, err := json.Marshal(validationEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal validation event: %w", err)
	}

	message := kafka.Message{
		Topic: "compliance_validation_events",
		Key:   request.RequestID,
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce validation event: %w", err)
	}

	// Publish individual violation events
	for _, violation := range result.Violations {
		violationEvent := map[string]interface{}{
			"event_type":    "compliance_violation_detected",
			"request_id":    request.RequestID,
			"violation_id":  violation.ID,
			"rule_id":       violation.RuleID,
			"framework":     violation.Framework,
			"severity":      violation.Severity,
			"message":       violation.Message,
			"api_id":        request.APIID,
			"endpoint_id":   request.EndpointID,
			"detected_at":   violation.DetectedAt,
		}

		violationJSON, err := json.Marshal(violationEvent)
		if err != nil {
			s.logger.Error("Failed to marshal violation event", "error", err)
			continue
		}

		violationMessage := kafka.Message{
			Topic: "compliance_violations",
			Key:   violation.ID,
			Value: violationJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, violationMessage); err != nil {
			s.logger.Error("Failed to produce violation event", "error", err)
		}
	}

	return nil
}
