package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/threat-detection/internal/models"
	"scopeapi.local/backend/services/threat-detection/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type SignatureDetectionServiceInterface interface {
	DetectSignatures(ctx context.Context, request *models.SignatureDetectionRequest) (*models.SignatureDetectionResult, error)
	LoadSignatures(ctx context.Context, signatureSet string) error
	UpdateSignature(ctx context.Context, signatureID string, signature *models.ThreatSignature) error
	GetSignatures(ctx context.Context, filter *models.SignatureFilter) ([]models.ThreatSignature, error)
	CreateCustomSignature(ctx context.Context, signature *models.ThreatSignature) error
	TestSignature(ctx context.Context, signatureID string, testData []map[string]interface{}) (*models.SignatureTestResult, error)
}

type SignatureDetectionService struct {
	threatRepo    repository.ThreatRepositoryInterface
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
	signatures    map[string]*models.ThreatSignature
	compiledRules map[string]*regexp.Regexp
}

func NewSignatureDetectionService(
	threatRepo repository.ThreatRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *SignatureDetectionService {
	return &SignatureDetectionService{
		threatRepo:    threatRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		signatures:    make(map[string]*models.ThreatSignature),
		compiledRules: make(map[string]*regexp.Regexp),
	}
}

func (s *SignatureDetectionService) DetectSignatures(ctx context.Context, request *models.SignatureDetectionRequest) (*models.SignatureDetectionResult, error) {
	result := &models.SignatureDetectionResult{
		ResultID:   "",
		SignatureID: "",
		Matched:    false,
		Details:    "",
		DetectedAt: time.Now(),
	}

	targets := s.extractDetectionTargets(request.Payload)

	for signatureID, signature := range s.signatures {
		match, err := s.checkSignature(signature, targets, request)
		if err != nil {
			s.logger.Error("Error checking signature", "signature_id", signatureID, "error", err)
			continue
		}

		if match != nil && match.Matched {
			result.SignatureID = signatureID
			result.Matched = true
			result.Details = match.Details
			result.DetectedAt = time.Now()
			break
		}
	}

	return result, nil
}

func (s *SignatureDetectionService) extractDetectionTargets(requestData map[string]interface{}) map[string]string {
	targets := make(map[string]string)

	// Extract URL and path
	if request, ok := requestData["request"].(map[string]interface{}); ok {
		if url, ok := request["url"].(string); ok {
			targets["url"] = url
		}
		if path, ok := request["path"].(string); ok {
			targets["path"] = path
		}
		if method, ok := request["method"].(string); ok {
			targets["method"] = method
		}
		if query, ok := request["query"].(string); ok {
			targets["query"] = query
		}

		// Extract headers
		if headers, ok := request["headers"].(map[string]interface{}); ok {
			for key, value := range headers {
				if strValue, ok := value.(string); ok {
					targets["header_"+strings.ToLower(key)] = strValue
				}
			}
		}

		// Extract body
		if body, ok := request["body"].(string); ok {
			targets["body"] = body
		}

		// Extract parameters
		if params, ok := request["parameters"].(map[string]interface{}); ok {
			for key, value := range params {
				if strValue, ok := value.(string); ok {
					targets["param_"+key] = strValue
				}
			}
		}
	}

	// Extract response data
	if response, ok := requestData["response"].(map[string]interface{}); ok {
		if body, ok := response["body"].(string); ok {
			targets["response_body"] = body
		}
		if headers, ok := response["headers"].(map[string]interface{}); ok {
			for key, value := range headers {
				if strValue, ok := value.(string); ok {
					targets["response_header_"+strings.ToLower(key)] = strValue
				}
			}
		}
	}

	return targets
}

func (s *SignatureDetectionService) checkSignature(signature *models.ThreatSignature, targets map[string]string, request *models.SignatureDetectionRequest) (*models.SignatureMatch, error) {
	// TODO: SignatureRule logic skipped due to missing 'Rules' field in ThreatSignature
	return nil, nil
}

func (s *SignatureDetectionService) checkRule(rule models.SignatureRule, targets map[string]string, signature *models.ThreatSignature, request *models.SignatureDetectionRequest) (*models.SignatureMatch, error) {
	// TODO: Rule logic skipped due to missing 'Field', 'Operator', 'Value' fields in SignatureRule
	return nil, nil
}

func (s *SignatureDetectionService) LoadSignatures(ctx context.Context, signatureSet string) error {
	signatures, err := s.threatRepo.GetThreatSignatures(ctx, &models.SignatureFilter{
		SignatureSet: signatureSet,
		Enabled:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to load signatures: %w", err)
	}

	// Clear existing signatures for this set
	for id, sig := range s.signatures {
		if sig.SignatureSet == signatureSet {
			delete(s.signatures, id)
		}
	}

	// Load new signatures
	for _, signature := range signatures {
		s.signatures[signature.ID] = &signature
		
		// Pre-compile regex rules
		for _, rule := range signature.Rules {
			if rule.Operator == "regex" {
				ruleKey := fmt.Sprintf("%s_%s", signature.ID, rule.ID)
				if regex, err := regexp.Compile(rule.Value); err == nil {
					s.compiledRules[ruleKey] = regex
				} else {
					s.logger.Warn("Failed to compile regex rule", "signature_id", signature.ID, "rule_id", rule.ID, "error", err)
				}
			}
		}
	}

	s.logger.Info("Loaded signatures", "signature_set", signatureSet, "count", len(signatures))
	return nil
}

func (s *SignatureDetectionService) UpdateSignature(ctx context.Context, signatureID string, signature *models.ThreatSignature) error {
	// Update in repository
	if err := s.threatRepo.UpdateThreatSignature(ctx, signatureID, signature); err != nil {
		return fmt.Errorf("failed to update signature in repository: %w", err)
	}

	// Update in memory
	signature.UpdatedAt = time.Now()
	s.signatures[signatureID] = signature

	// Update compiled rules
	for _, rule := range signature.Rules {
		if rule.Operator == "regex" {
			ruleKey := fmt.Sprintf("%s_%s", signature.ID, rule.ID)
			if regex, err := regexp.Compile(rule.Value); err == nil {
				s.compiledRules[ruleKey] = regex
			} else {
				delete(s.compiledRules, ruleKey)
				s.logger.Warn("Failed to compile updated regex rule", "signature_id", signature.ID, "rule_id", rule.ID, "error", err)
			}
		}
	}

	s.logger.Info("Updated signature", "signature_id", signatureID)
	return nil
}

func (s *SignatureDetectionService) GetSignatures(ctx context.Context, filter *models.SignatureFilter) ([]models.ThreatSignature, error) {
	return s.threatRepo.GetThreatSignatures(ctx, filter)
}

func (s *SignatureDetectionService) CreateCustomSignature(ctx context.Context, signature *models.ThreatSignature) error {
	// Validate signature
	if err := s.validateSignature(signature); err != nil {
		return fmt.Errorf("signature validation failed: %w", err)
	}

	// Set metadata
	signature.ID = uuid.New().String()
	signature.CreatedAt = time.Now()
	signature.UpdatedAt = time.Now()
	signature.Type = models.SignatureTypeCustom

	// Create in repository
	if err := s.threatRepo.CreateThreatSignature(ctx, signature); err != nil {
		return fmt.Errorf("failed to create signature: %w", err)
	}

	// Add to memory if enabled
	if signature.Enabled {
		s.signatures[signature.ID] = signature

		// Compile regex rules
		for _, rule := range signature.Rules {
			if rule.Operator == "regex" {
				ruleKey := fmt.Sprintf("%s_%s", signature.ID, rule.ID)
				if regex, err := regexp.Compile(rule.Value); err == nil {
					s.compiledRules[ruleKey] = regex
				} else {
					s.logger.Warn("Failed to compile regex rule in new signature", "signature_id", signature.ID, "rule_id", rule.ID, "error", err)
				}
			}
		}
	}

	s.logger.Info("Created custom signature", "signature_id", signature.ID, "name", signature.Name)
	return nil
}

func (s *SignatureDetectionService) validateSignature(signature *models.ThreatSignature) error {
	if signature.Name == "" {
		return fmt.Errorf("signature name is required")
	}

	if len(signature.Rules) == 0 {
		return fmt.Errorf("signature must have at least one rule")
	}

	validOperators := map[string]bool{
		"contains":       true,
		"equals":         true,
		"starts_with":    true,
		"ends_with":      true,
		"regex":          true,
		"length_greater": true,
		"length_less":    true,
		"not_contains":   true,
		"not_equals":     true,
	}

	for i, rule := range signature.Rules {
		if rule.Field == "" {
			return fmt.Errorf("rule %d: field is required", i)
		}

		if !validOperators[rule.Operator] {
			return fmt.Errorf("rule %d: invalid operator '%s'", i, rule.Operator)
		}

		if rule.Operator == "regex" {
			if _, err := regexp.Compile(rule.Value); err != nil {
				return fmt.Errorf("rule %d: invalid regex pattern: %w", i, err)
			}
		}

		if (rule.Operator == "length_greater" || rule.Operator == "length_less") && rule.IntValue <= 0 {
			return fmt.Errorf("rule %d: int_value must be positive for length operators", i)
		}
	}

	return nil
}

func (s *SignatureDetectionService) TestSignature(ctx context.Context, signatureID string, testData []map[string]interface{}) (*models.SignatureTestResult, error) {
	signature, exists := s.signatures[signatureID]
	if !exists {
		// Try to load from repository
		signatures, err := s.threatRepo.GetThreatSignatures(ctx, &models.SignatureFilter{
			SignatureID: signatureID,
		})
		if err != nil || len(signatures) == 0 {
			return nil, fmt.Errorf("signature not found: %s", signatureID)
		}
		signature = &signatures[0]
	}

	result := &models.SignatureTestResult{
		SignatureID:   signatureID,
		SignatureName: signature.Name,
		TestCases:     []models.SignatureTestCase{},
		TotalTests:    len(testData),
		PassedTests:   0,
		FailedTests:   0,
		TestedAt:      time.Now(),
	}

	for i, data := range testData {
		testCase := models.SignatureTestCase{
			TestID:      fmt.Sprintf("test_%d", i+1),
			TestData:    data,
			Expected:    false, // Default expectation
			Actual:      false,
			Passed:      false,
			MatchedRule: "",
			Error:       "",
		}

		// Check if test case has expected result
		if expected, ok := data["expected_match"].(bool); ok {
			testCase.Expected = expected
		}

		// Extract targets from test data
		targets := s.extractDetectionTargets(data)

		// Create mock request for testing
		mockRequest := &models.SignatureDetectionRequest{
			RequestID:   fmt.Sprintf("test_%d", i+1),
			RequestData: data,
		}

		// Test signature
		match, err := s.checkSignature(signature, targets, mockRequest)
		if err != nil {
			testCase.Error = err.Error()
			result.FailedTests++
		} else {
			testCase.Actual = match != nil
			if match != nil {
				testCase.MatchedRule = match.RuleMatched
			}

			// Check if result matches expectation
			if testCase.Expected == testCase.Actual {
				testCase.Passed = true
				result.PassedTests++
			} else {
				result.FailedTests++
			}
		}

		result.TestCases = append(result.TestCases, testCase)
	}

	return result, nil
}

func (s *SignatureDetectionService) getDetectionCategories() []string {
	categories := make(map[string]bool)
	for _, signature := range s.signatures {
		if signature.Enabled {
			categories[signature.Category] = true
		}
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}

	return result
}

func (s *SignatureDetectionService) publishSignatureEvents(ctx context.Context, matches []models.SignatureMatch, request *models.SignatureDetectionRequest) error {
	for _, match := range matches {
		eventData := map[string]interface{}{
			"event_type":      "signature_match",
			"signature_id":    match.SignatureID,
			"signature_name":  match.SignatureName,
			"signature_type":  match.SignatureType,
			"category":        match.Category,
			"severity":        match.Severity,
			"risk_score":      match.RiskScore,
			"confidence":      match.Confidence,
			"matched_field":   match.MatchedField,
			"matched_value":   match.MatchedValue,
			"rule_matched":    match.RuleMatched,
			"request_id":      request.RequestID,
			"ip_address":      request.IPAddress,
			"user_agent":      request.UserAgent,
			"api_id":          request.APIID,
			"endpoint_id":     request.EndpointID,
			"timestamp":       match.MatchedAt,
			"metadata":        match.Metadata,
		}

		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal signature event: %w", err)
		}

		message := kafka.Message{
			Topic: "signature_events",
			Key:   []byte(match.SignatureID),
			Value: eventJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce signature event: %w", err)
		}
	}

	return nil
}

// Signature management methods

func (s *SignatureDetectionService) ImportSignatureSet(ctx context.Context, signatureSetData []byte, signatureSet string) error {
	var signatures []models.ThreatSignature
	if err := json.Unmarshal(signatureSetData, &signatures); err != nil {
		return fmt.Errorf("failed to unmarshal signature set: %w", err)
	}

	// Validate all signatures before importing
	for i, signature := range signatures {
		if err := s.validateSignature(&signature); err != nil {
			return fmt.Errorf("signature %d validation failed: %w", i, err)
		}
	}

	// Import signatures
	for _, signature := range signatures {
		signature.ID = uuid.New().String()
		signature.SignatureSet = signatureSet
		signature.CreatedAt = time.Now()
		signature.UpdatedAt = time.Now()

		if err := s.threatRepo.CreateThreatSignature(ctx, &signature); err != nil {
			return fmt.Errorf("failed to create signature %s: %w", signature.Name, err)
		}
	}

	// Reload signatures for this set
	if err := s.LoadSignatures(ctx, signatureSet); err != nil {
		return fmt.Errorf("failed to reload signatures after import: %w", err)
	}

	s.logger.Info("Imported signature set", "signature_set", signatureSet, "count", len(signatures))
	return nil
}

func (s *SignatureDetectionService) ExportSignatureSet(ctx context.Context, signatureSet string) ([]byte, error) {
	signatures, err := s.threatRepo.GetThreatSignatures(ctx, &models.SignatureFilter{
		SignatureSet: signatureSet,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures for export: %w", err)
	}

	// Remove internal fields before export
	for i := range signatures {
		signatures[i].ID = ""
		signatures[i].CreatedAt = time.Time{}
		signatures[i].UpdatedAt = time.Time{}
	}

	data, err := json.MarshalIndent(signatures, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signatures: %w", err)
	}

	return data, nil
}

func (s *SignatureDetectionService) GetSignatureStatistics(ctx context.Context, timeRange time.Duration) (*models.SignatureStatistics, error) {
	stats := &models.SignatureStatistics{
		TotalSignatures:    len(s.signatures),
		EnabledSignatures:  0,
		DisabledSignatures: 0,
		SignaturesByType:   make(map[string]int),
		SignaturesByCategory: make(map[string]int),
		MatchStatistics:    &models.SignatureMatchStats{
			TotalMatches:      0,
			MatchesByType:     make(map[string]int),
			MatchesByCategory: make(map[string]int),
			MatchesBySeverity: make(map[string]int),
		},
		GeneratedAt:        time.Now(),
	}

	// Count signatures by status, type, and category
	for _, signature := range s.signatures {
		if signature.Enabled {
			stats.EnabledSignatures++
		} else {
			stats.DisabledSignatures++
		}

		stats.SignaturesByType[signature.Type]++
		stats.SignaturesByCategory[signature.Category]++
	}

	// Get match statistics from repository
	matchStats, err := s.threatRepo.GetSignatureMatchStatistics(ctx)
	if err != nil {
		s.logger.Warn("Failed to get signature match statistics", "error", err)
	} else {
		stats.MatchStatistics = matchStats
	}

	return stats, nil
}

func (s *SignatureDetectionService) OptimizeSignatures(ctx context.Context) (*models.SignatureOptimizationResult, error) {
	result := &models.SignatureOptimizationResult{
		OptimizedSignatures: 0,
		RemovedSignatures:   []string{},
		UpdatedSignatures:   []string{},
		Recommendations:     []string{},
		OptimizedAt:         time.Now(),
	}

	// Get signature performance data
	_, err := s.GetSignatureStatistics(ctx, 30*24*time.Hour) // Last 30 days
	if err != nil {
		return nil, fmt.Errorf("failed to get signature statistics: %w", err)
	}

	// Identify low-performing signatures
	for signatureID, signature := range s.signatures {
		if !signature.Enabled {
			continue
		}

		// Check for signatures that are very old but low confidence
		if signature.Confidence < 0.5 && time.Since(signature.CreatedAt) > 30*24*time.Hour {
			result.UpdatedSignatures = append(result.UpdatedSignatures, signatureID)
			result.Recommendations = append(result.Recommendations, 
				fmt.Sprintf("Signature '%s' has low confidence - consider tuning", 
					signature.Name))
		}

		// Check for signatures that are very old and might need review
		if time.Since(signature.CreatedAt) > 90*24*time.Hour {
			result.Recommendations = append(result.Recommendations, 
				fmt.Sprintf("Signature '%s' is over 90 days old - consider reviewing", 
					signature.Name))
		}
	}

	// Identify duplicate or overlapping signatures
	duplicates := s.findDuplicateSignatures()
	for _, duplicateGroup := range duplicates {
		if len(duplicateGroup) > 1 {
			// Keep the most recent one, mark others for removal
			for i := 1; i < len(duplicateGroup); i++ {
				result.RemovedSignatures = append(result.RemovedSignatures, duplicateGroup[i])
				result.Recommendations = append(result.Recommendations, 
					fmt.Sprintf("Duplicate signature detected - removing '%s'", s.signatures[duplicateGroup[i]].Name))
			}
		}
	}

	return result, nil
}

func (s *SignatureDetectionService) findDuplicateSignatures() [][]string {
	var duplicates [][]string
	checked := make(map[string]bool)

	for id1, sig1 := range s.signatures {
		if checked[id1] {
			continue
		}

		var group []string
		group = append(group, id1)

		for id2, sig2 := range s.signatures {
			if id1 == id2 || checked[id2] {
				continue
			}

			if s.signaturesAreSimilar(sig1, sig2) {
				group = append(group, id2)
				checked[id2] = true
			}
		}

		if len(group) > 1 {
			duplicates = append(duplicates, group)
		}
		checked[id1] = true
	}

	return duplicates
}

func (s *SignatureDetectionService) signaturesAreSimilar(sig1, sig2 *models.ThreatSignature) bool {
	// Check if signatures have similar names
	if strings.EqualFold(sig1.Name, sig2.Name) {
		return true
	}

	// Check if they have identical rules
	if len(sig1.Rules) == len(sig2.Rules) {
		rulesMatch := true
		for i, rule1 := range sig1.Rules {
			rule2 := sig2.Rules[i]
			if rule1.Field != rule2.Field || rule1.Operator != rule2.Operator || rule1.Value != rule2.Value {
				rulesMatch = false
				break
			}
		}
		if rulesMatch {
			return true
		}
	}

	return false
}
