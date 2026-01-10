package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"data-protection/internal/models"
	"data-protection/internal/repository"

	"github.com/google/uuid"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type DataClassificationServiceInterface interface {
	ClassifyData(ctx context.Context, request *models.DataClassificationRequest) (*models.DataClassificationResult, error)
	CreateClassificationRule(ctx context.Context, rule *models.ClassificationRule) error
	UpdateClassificationRule(ctx context.Context, ruleID string, rule *models.ClassificationRule) error
	GetClassificationRules(ctx context.Context, filter *models.ClassificationRuleFilter) ([]models.ClassificationRule, error)
	ApplyDataLabels(ctx context.Context, data map[string]interface{}, labels []models.DataLabel) error
	GetDataClassificationReport(ctx context.Context, filter *models.ClassificationReportFilter) (*models.ClassificationReport, error)
}

type DataClassificationService struct {
	classificationRepo repository.ClassificationRepositoryInterface
	kafkaProducer      kafka.ProducerInterface
	logger             logging.Logger
	rules              map[string]*models.ClassificationRule
	labelHierarchy     map[string][]string
}

func NewDataClassificationService(
	classificationRepo repository.ClassificationRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *DataClassificationService {
	service := &DataClassificationService{
		classificationRepo: classificationRepo,
		kafkaProducer:      kafkaProducer,
		logger:             logger,
		rules:              make(map[string]*models.ClassificationRule),
		labelHierarchy:     make(map[string][]string),
	}

	// Load default classification rules
	service.loadDefaultRules()
	service.setupLabelHierarchy()

	return service
}

func (s *DataClassificationService) loadDefaultRules() {
	defaultRules := []models.ClassificationRule{
		{
			ID:          "public_data",
			Name:        "Public Data",
			Description: "Data that can be freely shared",
			Category:    models.ClassificationCategoryPublic,
			Priority:    1,
			Conditions: []models.ClassificationCondition{
				{
					Field:    "data_type",
					Operator: "equals",
					Value:    "public",
				},
			},
			Labels: []models.DataLabel{
				{
					Key:   "classification",
					Value: "public",
				},
				{
					Key:   "sensitivity",
					Value: "low",
				},
			},
			Actions: []models.ClassificationAction{
				{
					Type:   "label",
					Config: map[string]interface{}{"auto_apply": true},
				},
			},
			Enabled: true,
		},
		{
			ID:          "internal_data",
			Name:        "Internal Data",
			Description: "Data for internal use only",
			Category:    models.ClassificationCategoryInternal,
			Priority:    2,
			Conditions: []models.ClassificationCondition{
				{
					Field:    "field_name",
					Operator: "contains",
					Value:    "internal",
				},
			},
			Labels: []models.DataLabel{
				{
					Key:   "classification",
					Value: "internal",
				},
				{
					Key:   "sensitivity",
					Value: "medium",
				},
			},
			Actions: []models.ClassificationAction{
				{
					Type:   "label",
					Config: map[string]interface{}{"auto_apply": true},
				},
				{
					Type:   "encrypt",
					Config: map[string]interface{}{"algorithm": "AES-256"},
				},
			},
			Enabled: true,
		},
		{
			ID:          "confidential_data",
			Name:        "Confidential Data",
			Description: "Sensitive data requiring protection",
			Category:    models.ClassificationCategoryConfidential,
			Priority:    3,
			Conditions: []models.ClassificationCondition{
				{
					Field:    "pii_detected",
					Operator: "equals",
					Value:    "true",
				},
			},
			Labels: []models.DataLabel{
				{
					Key:   "classification",
					Value: "confidential",
				},
				{
					Key:   "sensitivity",
					Value: "high",
				},
			},
			Actions: []models.ClassificationAction{
				{
					Type:   "label",
					Config: map[string]interface{}{"auto_apply": true},
				},
				{
					Type:   "encrypt",
					Config: map[string]interface{}{"algorithm": "AES-256"},
				},
				{
					Type:   "audit",
					Config: map[string]interface{}{"log_access": true},
				},
			},
			Enabled: true,
		},
		{
			ID:          "restricted_data",
			Name:        "Restricted Data",
			Description: "Highly sensitive data with strict access controls",
			Category:    models.ClassificationCategoryRestricted,
			Priority:    4,
			Conditions: []models.ClassificationCondition{
				{
					Field:    "pii_type",
					Operator: "in",
					Value:    "ssn,credit_card,passport",
				},
			},
			Labels: []models.DataLabel{
				{
					Key:   "classification",
					Value: "restricted",
				},
				{
					Key:   "sensitivity",
					Value: "critical",
				},
			},
			Actions: []models.ClassificationAction{
				{
					Type:   "label",
					Config: map[string]interface{}{"auto_apply": true},
				},
				{
					Type:   "encrypt",
					Config: map[string]interface{}{"algorithm": "AES-256-GCM"},
				},
				{
					Type:   "audit",
					Config: map[string]interface{}{"log_access": true, "alert_access": true},
				},
				{
					Type:   "access_control",
					Config: map[string]interface{}{"require_approval": true},
				},
			},
			Enabled: true,
		},
	}

	for _, rule := range defaultRules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		s.rules[rule.ID] = &rule
	}
}

func (s *DataClassificationService) setupLabelHierarchy() {
	s.labelHierarchy = map[string][]string{
		"classification": {"public", "internal", "confidential", "restricted"},
		"sensitivity":    {"low", "medium", "high", "critical"},
		"retention":      {"short", "medium", "long", "permanent"},
		"geography":      {"global", "regional", "country", "local"},
	}
}

func (s *DataClassificationService) ClassifyData(ctx context.Context, request *models.DataClassificationRequest) (*models.DataClassificationResult, error) {
	startTime := time.Now()

	result := &models.DataClassificationResult{
		RequestID:       request.RequestID,
		Classifications: []models.DataClassification{},
		AppliedLabels:   []models.DataLabel{},
		ExecutedActions: []models.ClassificationAction{},
		RulesMatched:    []string{},
		ProcessingTime:  0,
		ClassifiedAt:    time.Now(),
	}

	// Analyze data structure
	dataAnalysis := s.analyzeDataStructure(request.Data)

	// Apply classification rules
	for _, rule := range s.getSortedRules() {
		if !rule.Enabled {
			continue
		}

		matched, matchContext := s.evaluateRule(rule, request.Data, dataAnalysis)
		if matched {
			classification := models.DataClassification{
				ID:           uuid.New().String(),
				RuleID:       rule.ID,
				RuleName:     rule.Name,
				Category:     rule.Category,
				Labels:       rule.Labels,
				Confidence:   s.calculateConfidence(rule, matchContext),
				MatchContext: matchContext,
				ClassifiedAt: time.Now(),
			}

			result.Classifications = append(result.Classifications, classification)
			result.RulesMatched = append(result.RulesMatched, rule.ID)

			// Apply labels
			for _, label := range rule.Labels {
				if !s.labelExists(result.AppliedLabels, label) {
					result.AppliedLabels = append(result.AppliedLabels, label)
				}
			}

			// Execute actions
			for _, action := range rule.Actions {
				if err := s.executeAction(ctx, action, request.Data, classification); err != nil {
					s.logger.Warn("Failed to execute classification action",
						"rule_id", rule.ID, "action_type", action.Type, "error", err)
				} else {
					result.ExecutedActions = append(result.ExecutedActions, action)
				}
			}
		}
	}

	result.ProcessingTime = time.Since(startTime)

	// Store classification results
	if len(result.Classifications) > 0 {
		classificationData := &models.ClassificationData{
			ID:              uuid.New().String(),
			RequestID:       request.RequestID,
			APIID:           request.APIID,
			EndpointID:      request.EndpointID,
			Classifications: result.Classifications,
			Labels:          result.AppliedLabels,
			DataHash:        s.calculateDataHash(request.Data),
			IPAddress:       request.IPAddress,
			UserAgent:       request.UserAgent,
			ClassifiedAt:    time.Now(),
			Metadata: map[string]interface{}{
				"data_source":     request.DataSource,
				"rules_matched":   result.RulesMatched,
				"processing_time": result.ProcessingTime.Milliseconds(),
				"data_analysis":   dataAnalysis,
			},
		}

		if err := s.classificationRepo.CreateClassificationData(ctx, classificationData); err != nil {
			s.logger.Error("Failed to store classification data", "error", err)
		}

		// Publish classification events
		if err := s.publishClassificationEvents(ctx, result, request); err != nil {
			s.logger.Error("Failed to publish classification events", "error", err)
		}
	}

	return result, nil
}

func (s *DataClassificationService) analyzeDataStructure(data map[string]interface{}) *models.DataAnalysis {
	analysis := &models.DataAnalysis{
		FieldCount:     0,
		FieldTypes:     make(map[string]int),
		FieldNames:     []string{},
		DataPatterns:   []string{},
		NestedLevels:   0,
		ArrayFields:    []string{},
		SensitiveHints: []string{},
	}

	s.analyzeDataRecursive(data, "", analysis, 0)

	// Detect sensitive field name patterns
	analysis.SensitiveHints = s.detectSensitiveFieldNames(analysis.FieldNames)

	return analysis
}

func (s *DataClassificationService) analyzeDataRecursive(data interface{}, path string, analysis *models.DataAnalysis, level int) {
	if level > analysis.NestedLevels {
		analysis.NestedLevels = level
	}

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			currentPath := key
			if path != "" {
				currentPath = path + "." + key
			}

			analysis.FieldCount++
			analysis.FieldNames = append(analysis.FieldNames, key)

			s.analyzeDataRecursive(value, currentPath, analysis, level+1)
		}
	case []interface{}:
		analysis.ArrayFields = append(analysis.ArrayFields, path)
		analysis.FieldTypes["array"]++

		for i, item := range v {
			currentPath := fmt.Sprintf("%s[%d]", path, i)
			s.analyzeDataRecursive(item, currentPath, analysis, level+1)
		}
	case string:
		analysis.FieldTypes["string"]++
		if len(v) > 100 {
			analysis.DataPatterns = append(analysis.DataPatterns, "long_string")
		}
		if s.containsSpecialChars(v) {
			analysis.DataPatterns = append(analysis.DataPatterns, "special_chars")
		}
	case float64:
		analysis.FieldTypes["number"]++
	case bool:
		analysis.FieldTypes["boolean"]++
	case nil:
		analysis.FieldTypes["null"]++
	default:
		analysis.FieldTypes["unknown"]++
	}
}

func (s *DataClassificationService) detectSensitiveFieldNames(fieldNames []string) []string {
	var hints []string

	sensitivePatterns := map[string]string{
		"password":  "authentication",
		"token":     "authentication",
		"key":       "authentication",
		"secret":    "authentication",
		"email":     "contact",
		"phone":     "contact",
		"address":   "location",
		"ssn":       "identifier",
		"id":        "identifier",
		"card":      "financial",
		"account":   "financial",
		"salary":    "financial",
		"medical":   "health",
		"health":    "health",
		"diagnosis": "health",
	}

	for _, fieldName := range fieldNames {
		fieldLower := strings.ToLower(fieldName)
		for pattern, category := range sensitivePatterns {
			if strings.Contains(fieldLower, pattern) {
				hint := fmt.Sprintf("%s_field_detected:%s", category, fieldName)
				if !s.containsString(hints, hint) {
					hints = append(hints, hint)
				}
			}
		}
	}

	return hints
}

func (s *DataClassificationService) containsSpecialChars(str string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, char := range specialChars {
		if strings.ContainsRune(str, char) {
			return true
		}
	}
	return false
}

func (s *DataClassificationService) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *DataClassificationService) getSortedRules() []*models.ClassificationRule {
	var rules []*models.ClassificationRule
	for _, rule := range s.rules {
		rules = append(rules, rule)
	}

	// Sort by priority (higher priority first)
	for i := 0; i < len(rules)-1; i++ {
		for j := i + 1; j < len(rules); j++ {
			if rules[i].Priority < rules[j].Priority {
				rules[i], rules[j] = rules[j], rules[i]
			}
		}
	}

	return rules
}

func (s *DataClassificationService) evaluateRule(rule *models.ClassificationRule, data map[string]interface{}, analysis *models.DataAnalysis) (bool, map[string]interface{}) {
	matchContext := make(map[string]interface{})
	allConditionsMet := true

	for _, condition := range rule.Conditions {
		conditionMet, context := s.evaluateCondition(condition, data, analysis)
		if conditionMet {
			matchContext[condition.Field] = context
		} else {
			allConditionsMet = false
			break
		}
	}

	return allConditionsMet, matchContext
}

func (s *DataClassificationService) evaluateCondition(condition models.ClassificationCondition, data map[string]interface{}, analysis *models.DataAnalysis) (bool, interface{}) {
	switch condition.Field {
	case "field_name":
		return s.evaluateFieldNameCondition(condition, analysis.FieldNames)
	case "field_count":
		return s.evaluateNumericCondition(condition, float64(analysis.FieldCount))
	case "data_type":
		return s.evaluateDataTypeCondition(condition, analysis.FieldTypes)
	case "nested_levels":
		return s.evaluateNumericCondition(condition, float64(analysis.NestedLevels))
	case "sensitive_hints":
		return s.evaluateSensitiveHintsCondition(condition, analysis.SensitiveHints)
	case "pii_detected":
		return s.evaluatePIICondition(condition, data)
	case "pii_type":
		return s.evaluatePIITypeCondition(condition, data)
	default:
		// Try to evaluate against actual data fields
		return s.evaluateDataFieldCondition(condition, data)
	}
}

func (s *DataClassificationService) evaluateFieldNameCondition(condition models.ClassificationCondition, fieldNames []string) (bool, interface{}) {
	switch condition.Operator {
	case "contains":
		for _, fieldName := range fieldNames {
			if strings.Contains(strings.ToLower(fieldName), strings.ToLower(condition.Value)) {
				return true, fieldName
			}
		}
	case "equals":
		for _, fieldName := range fieldNames {
			if strings.EqualFold(fieldName, condition.Value) {
				return true, fieldName
			}
		}
	case "regex":
		// Implementation for regex matching
		// This would require compiling the regex and matching against field names
	}
	return false, nil
}

func (s *DataClassificationService) evaluateNumericCondition(condition models.ClassificationCondition, value float64) (bool, interface{}) {
	conditionValue := 0.0
	if v, ok := condition.Value.(float64); ok {
		conditionValue = v
	} else if v, ok := condition.Value.(string); ok {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			conditionValue = parsed
		}
	}

	switch condition.Operator {
	case "equals":
		return value == conditionValue, value
	case "greater_than":
		return value > conditionValue, value
	case "less_than":
		return value < conditionValue, value
	case "greater_equal":
		return value >= conditionValue, value
	case "less_equal":
		return value <= conditionValue, value
	}
	return false, nil
}

func (s *DataClassificationService) evaluateDataTypeCondition(condition models.ClassificationCondition, fieldTypes map[string]int) (bool, interface{}) {
	switch condition.Operator {
	case "contains":
		if count, exists := fieldTypes[condition.Value]; exists && count > 0 {
			return true, count
		}
	case "dominant":
		// Check if the specified type is the most common
		maxCount := 0
		dominantType := ""
		for dataType, count := range fieldTypes {
			if count > maxCount {
				maxCount = count
				dominantType = dataType
			}
		}
		return dominantType == condition.Value, dominantType
	}
	return false, nil
}

func (s *DataClassificationService) evaluateSensitiveHintsCondition(condition models.ClassificationCondition, hints []string) (bool, interface{}) {
	switch condition.Operator {
	case "contains":
		for _, hint := range hints {
			if strings.Contains(hint, condition.Value) {
				return true, hint
			}
		}
	case "count_greater":
		conditionValue := 0
		if v, ok := condition.Value.(string); ok {
			if parsed, err := strconv.Atoi(v); err == nil {
				conditionValue = parsed
			}
		}
		return len(hints) > conditionValue, len(hints)
	}
	return false, nil
}

func (s *DataClassificationService) evaluatePIICondition(condition models.ClassificationCondition, data map[string]interface{}) (bool, interface{}) {
	// This would typically integrate with the PII detection service
	// For now, we'll do a simple check for common PII patterns
	hasPII := s.hasBasicPIIPatterns(data)

	switch condition.Operator {
	case "equals":
		expectedValue := condition.Value == "true"
		return hasPII == expectedValue, hasPII
	}
	return false, nil
}

func (s *DataClassificationService) evaluatePIITypeCondition(condition models.ClassificationCondition, data map[string]interface{}) (bool, interface{}) {
	detectedTypes := s.detectBasicPIITypes(data)

	switch condition.Operator {
	case "in":
		targetTypes := strings.Split(condition.Value, ",")
		for _, targetType := range targetTypes {
			targetType = strings.TrimSpace(targetType)
			for _, detectedType := range detectedTypes {
				if detectedType == targetType {
					return true, detectedTypes
				}
			}
		}
	case "contains":
		for _, detectedType := range detectedTypes {
			if detectedType == condition.Value {
				return true, detectedTypes
			}
		}
	}
	return false, nil
}

func (s *DataClassificationService) evaluateDataFieldCondition(condition models.ClassificationCondition, data map[string]interface{}) (bool, interface{}) {
	value := s.getNestedValue(data, condition.Field)
	if value == nil {
		return false, nil
	}

	switch condition.Operator {
	case "equals":
		return fmt.Sprintf("%v", value) == condition.Value, value
	case "contains":
		if str, ok := value.(string); ok {
			return strings.Contains(strings.ToLower(str), strings.ToLower(condition.Value)), value
		}
	case "exists":
		return true, value
	case "not_empty":
		if str, ok := value.(string); ok {
			return strings.TrimSpace(str) != "", value
		}
		return value != nil, value
	}
	return false, nil
}

func (s *DataClassificationService) getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

func (s *DataClassificationService) hasBasicPIIPatterns(data map[string]interface{}) bool {
	// Simple PII detection patterns
	patterns := []string{
		`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`, // Email
		`\b\d{3}-?\d{2}-?\d{4}\b`,                             // SSN
		`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`,          // Credit Card
	}

	dataStr := fmt.Sprintf("%v", data)
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, dataStr); matched {
			return true
		}
	}
	return false
}

func (s *DataClassificationService) detectBasicPIITypes(data map[string]interface{}) []string {
	var types []string
	dataStr := fmt.Sprintf("%v", data)

	patterns := map[string]string{
		"email":       `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
		"ssn":         `\b\d{3}-?\d{2}-?\d{4}\b`,
		"credit_card": `\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`,
		"phone":       `\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`,
	}

	for piiType, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, dataStr); matched {
			types = append(types, piiType)
		}
	}

	return types
}

func (s *DataClassificationService) calculateConfidence(rule *models.ClassificationRule, matchContext map[string]interface{}) float64 {
	baseConfidence := 0.8 // Default confidence

	// Adjust confidence based on number of conditions matched
	conditionCount := len(rule.Conditions)
	matchCount := len(matchContext)

	if conditionCount > 0 {
		matchRatio := float64(matchCount) / float64(conditionCount)
		baseConfidence *= matchRatio
	}

	// Adjust based on rule priority
	priorityBoost := float64(rule.Priority) * 0.05
	baseConfidence += priorityBoost

	// Cap at 1.0
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}

	return baseConfidence
}

func (s *DataClassificationService) labelExists(labels []models.DataLabel, newLabel models.DataLabel) bool {
	for _, label := range labels {
		if label.Key == newLabel.Key && label.Value == newLabel.Value {
			return true
		}
	}
	return false
}

func (s *DataClassificationService) executeAction(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	switch action.Type {
	case "label":
		return s.executeLabeling(ctx, action, data, classification)
	case "encrypt":
		return s.executeEncryption(ctx, action, data, classification)
	case "audit":
		return s.executeAuditing(ctx, action, data, classification)
	case "access_control":
		return s.executeAccessControl(ctx, action, data, classification)
	case "notify":
		return s.executeNotification(ctx, action, data, classification)
	case "quarantine":
		return s.executeQuarantine(ctx, action, data, classification)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (s *DataClassificationService) executeLabeling(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	autoApply, ok := action.Config["auto_apply"].(bool)
	if !ok || !autoApply {
		return nil
	}

	// Labels are already applied in the main classification logic
	s.logger.Debug("Labels applied automatically", "rule_id", classification.RuleID)
	return nil
}

func (s *DataClassificationService) executeEncryption(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	algorithm, ok := action.Config["algorithm"].(string)
	if !ok {
		algorithm = "AES-256"
	}

	s.logger.Info("Encryption action triggered",
		"rule_id", classification.RuleID,
		"algorithm", algorithm,
		"classification", classification.Category)

	// In a real implementation, this would trigger encryption of the data
	// For now, we'll just log the action
	return nil
}

func (s *DataClassificationService) executeAuditing(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	logAccess, _ := action.Config["log_access"].(bool)
	alertAccess, _ := action.Config["alert_access"].(bool)

	if logAccess {
		auditEvent := map[string]interface{}{
			"event_type":     "data_access_audit",
			"classification": classification.Category,
			"rule_id":        classification.RuleID,
			"rule_name":      classification.RuleName,
			"confidence":     classification.Confidence,
			"timestamp":      time.Now(),
			"alert_required": alertAccess,
		}

		auditJSON, err := json.Marshal(auditEvent)
		if err != nil {
			return fmt.Errorf("failed to marshal audit event: %w", err)
		}

		message := kafka.Message{
			Topic: "audit_events",
			Key:   classification.ID,
			Value: auditJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce audit event: %w", err)
		}
	}

	return nil
}

func (s *DataClassificationService) executeAccessControl(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	requireApproval, _ := action.Config["require_approval"].(bool)

	if requireApproval {
		s.logger.Warn("Access control triggered - approval required",
			"rule_id", classification.RuleID,
			"classification", classification.Category)

		// In a real implementation, this would integrate with an approval workflow system
		accessControlEvent := map[string]interface{}{
			"event_type":        "access_control_required",
			"classification":    classification.Category,
			"rule_id":           classification.RuleID,
			"requires_approval": true,
			"timestamp":         time.Now(),
		}

		eventJSON, err := json.Marshal(accessControlEvent)
		if err != nil {
			return fmt.Errorf("failed to marshal access control event: %w", err)
		}

		message := kafka.Message{
			Topic: "access_control_events",
			Key:   classification.ID,
			Value: eventJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce access control event: %w", err)
		}
	}

	return nil
}

func (s *DataClassificationService) executeNotification(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	recipients, ok := action.Config["recipients"].([]string)
	if !ok {
		return fmt.Errorf("no recipients specified for notification")
	}

	notificationEvent := map[string]interface{}{
		"event_type":     "classification_notification",
		"classification": classification.Category,
		"rule_id":        classification.RuleID,
		"rule_name":      classification.RuleName,
		"confidence":     classification.Confidence,
		"recipients":     recipients,
		"timestamp":      time.Now(),
		"message":        fmt.Sprintf("Data classified as %s with confidence %.2f", classification.Category, classification.Confidence),
	}

	eventJSON, err := json.Marshal(notificationEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal notification event: %w", err)
	}

	message := kafka.Message{
		Topic: "notification_events",
		Key:   classification.ID,
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce notification event: %w", err)
	}

	return nil
}

func (s *DataClassificationService) executeQuarantine(ctx context.Context, action models.ClassificationAction, data map[string]interface{}, classification models.DataClassification) error {
	reason, ok := action.Config["reason"].(string)
	if !ok {
		reason = "Automatic quarantine due to data classification"
	}

	quarantineEvent := map[string]interface{}{
		"event_type":     "data_quarantine",
		"classification": classification.Category,
		"rule_id":        classification.RuleID,
		"reason":         reason,
		"timestamp":      time.Now(),
		"data_hash":      s.calculateDataHash(data),
	}

	eventJSON, err := json.Marshal(quarantineEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal quarantine event: %w", err)
	}

	message := kafka.Message{
		Topic: "quarantine_events",
		Key:   classification.ID,
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce quarantine event: %w", err)
	}

	s.logger.Warn("Data quarantined",
		"rule_id", classification.RuleID,
		"classification", classification.Category,
		"reason", reason)

	return nil
}

func (s *DataClassificationService) calculateDataHash(data map[string]interface{}) string {
	dataJSON, _ := json.Marshal(data)
	hash := sha256.Sum256(dataJSON)
	return hex.EncodeToString(hash[:])
}

func (s *DataClassificationService) CreateClassificationRule(ctx context.Context, rule *models.ClassificationRule) error {
	// Validate rule
	if err := s.validateRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	// Set metadata
	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	// Store in repository
	if err := s.classificationRepo.CreateClassificationRule(ctx, rule); err != nil {
		return fmt.Errorf("failed to create classification rule: %w", err)
	}

	// Add to memory
	s.rules[rule.ID] = rule

	s.logger.Info("Created classification rule", "rule_id", rule.ID, "name", rule.Name)
	return nil
}

func (s *DataClassificationService) UpdateClassificationRule(ctx context.Context, ruleID string, rule *models.ClassificationRule) error {
	// Validate rule
	if err := s.validateRule(rule); err != nil {
		return fmt.Errorf("rule validation failed: %w", err)
	}

	// Update in repository
	rule.UpdatedAt = time.Now()
	if err := s.classificationRepo.UpdateClassificationRule(ctx, ruleID, rule); err != nil {
		return fmt.Errorf("failed to update classification rule: %w", err)
	}

	// Update in memory
	rule.ID = ruleID
	s.rules[ruleID] = rule

	s.logger.Info("Updated classification rule", "rule_id", ruleID, "name", rule.Name)
	return nil
}

func (s *DataClassificationService) validateRule(rule *models.ClassificationRule) error {
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if rule.Category == "" {
		return fmt.Errorf("rule category is required")
	}

	if len(rule.Conditions) == 0 {
		return fmt.Errorf("at least one condition is required")
	}

	if len(rule.Labels) == 0 {
		return fmt.Errorf("at least one label is required")
	}

	// Validate conditions
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

	// Validate labels
	for i, label := range rule.Labels {
		if label.Key == "" {
			return fmt.Errorf("label %d: key is required", i)
		}
		if label.Value == "" {
			return fmt.Errorf("label %d: value is required", i)
		}
	}

	// Validate actions
	for i, action := range rule.Actions {
		if action.Type == "" {
			return fmt.Errorf("action %d: type is required", i)
		}
	}

	return nil
}

func (s *DataClassificationService) GetClassificationRules(ctx context.Context, filter *models.ClassificationRuleFilter) ([]models.ClassificationRule, error) {
	var rules []models.ClassificationRule

	for _, rule := range s.rules {
		if s.matchesFilter(rule, filter) {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (s *DataClassificationService) matchesFilter(rule *models.ClassificationRule, filter *models.ClassificationRuleFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Category != "" && rule.Category != filter.Category {
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

func (s *DataClassificationService) ApplyDataLabels(ctx context.Context, data map[string]interface{}, labels []models.DataLabel) error {
	// In a real implementation, this would apply labels to the actual data storage
	// For now, we'll create a labeling event

	labelingEvent := map[string]interface{}{
		"event_type": "data_labeling",
		"labels":     labels,
		"data_hash":  s.calculateDataHash(data),
		"timestamp":  time.Now(),
	}

	eventJSON, err := json.Marshal(labelingEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal labeling event: %w", err)
	}

	message := kafka.Message{
		Topic: "labeling_events",
		Key:   s.calculateDataHash(data),
		Value: eventJSON,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce labeling event: %w", err)
	}

	s.logger.Info("Applied data labels", "label_count", len(labels))
	return nil
}

func (s *DataClassificationService) GetDataClassificationReport(ctx context.Context, filter *models.ClassificationReportFilter) (*models.ClassificationReport, error) {
	// This would typically query the repository for classification data
	// For now, we'll return a basic report structure

	report := &models.ClassificationReport{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		Filter:      filter,
		Summary: models.ClassificationSummary{
			TotalClassifications: 0,
			CategoryBreakdown:    make(map[string]int),
			LabelBreakdown:       make(map[string]int),
			RuleUsage:            make(map[string]int),
		},
		Classifications: []models.DataClassification{},
		Trends:          []models.ClassificationTrend{},
		Recommendations: []string{},
	}

	// In a real implementation, this would query the database
	// and populate the report with actual data

	return report, nil
}

func (s *DataClassificationService) publishClassificationEvents(ctx context.Context, result *models.DataClassificationResult, request *models.DataClassificationRequest) error {
	for _, classification := range result.Classifications {
		eventData := map[string]interface{}{
			"event_type":        "data_classified",
			"classification_id": classification.ID,
			"rule_id":           classification.RuleID,
			"rule_name":         classification.RuleName,
			"category":          classification.Category,
			"labels":            classification.Labels,
			"confidence":        classification.Confidence,
			"match_context":     classification.MatchContext,
			"request_id":        request.RequestID,
			"api_id":            request.APIID,
			"endpoint_id":       request.EndpointID,
			"ip_address":        request.IPAddress,
			"user_agent":        request.UserAgent,
			"data_source":       request.DataSource,
			"timestamp":         classification.ClassifiedAt,
		}

		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal classification event: %w", err)
		}

		message := kafka.Message{
			Topic: "classification_events",
			Key:   classification.ID,
			Value: eventJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce classification event: %w", err)
		}
	}

	return nil
}
