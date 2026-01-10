package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"data-protection/internal/models"
	"data-protection/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type PIIDetectionServiceInterface interface {
	DetectPII(ctx context.Context, request *models.PIIDetectionRequest) (*models.PIIDetectionResult, error)
	ScanData(ctx context.Context, data map[string]interface{}) ([]models.PIIFinding, error)
	GetPIIPatterns(ctx context.Context) ([]models.PIIPattern, error)
	CreateCustomPattern(ctx context.Context, pattern *models.PIIPattern) error
	UpdatePattern(ctx context.Context, patternID string, pattern *models.PIIPattern) error
	ValidateDataCompliance(ctx context.Context, data map[string]interface{}, regulations []string) (*models.ComplianceValidationResult, error)
}

type PIIDetectionService struct {
	piiRepo       repository.PIIRepositoryInterface
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
	patterns      map[string]*models.PIIPattern
	compiledRegex map[string]*regexp.Regexp
}

func NewPIIDetectionService(
	piiRepo repository.PIIRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *PIIDetectionService {
	service := &PIIDetectionService{
		piiRepo:       piiRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		patterns:      make(map[string]*models.PIIPattern),
		compiledRegex: make(map[string]*regexp.Regexp),
	}

	// Load default PII patterns
	service.loadDefaultPatterns()

	return service
}

func (s *PIIDetectionService) loadDefaultPatterns() {
	defaultPatterns := []models.PIIPattern{
		{
			ID:          "ssn",
			Name:        "Social Security Number",
			Type:        models.PIITypeSSN,
			Category:    models.PIICategoryIdentifier,
			Pattern:     `\b\d{3}-?\d{2}-?\d{4}\b`,
			Sensitivity: models.SensitivityHigh,
			Confidence:  0.9,
			Enabled:     true,
			Description: "US Social Security Number pattern",
		},
		{
			ID:          "email",
			Name:        "Email Address",
			Type:        models.PIITypeEmail,
			Category:    models.PIICategoryContact,
			Pattern:     `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
			Sensitivity: models.SensitivityMedium,
			Confidence:  0.8,
			Enabled:     true,
			Description: "Email address pattern",
		},
		{
			ID:          "phone",
			Name:        "Phone Number",
			Type:        models.PIITypePhone,
			Category:    models.PIICategoryContact,
			Pattern:     `\b(?:\+?1[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})\b`,
			Sensitivity: models.SensitivityMedium,
			Confidence:  0.7,
			Enabled:     true,
			Description: "US phone number pattern",
		},
		{
			ID:          "credit_card",
			Name:        "Credit Card Number",
			Type:        models.PIITypeCreditCard,
			Category:    models.PIICategoryFinancial,
			Pattern:     `\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3[0-9]{13}|6(?:011|5[0-9]{2})[0-9]{12})\b`,
			Sensitivity: models.SensitivityHigh,
			Confidence:  0.9,
			Enabled:     true,
			Description: "Credit card number pattern (Visa, MasterCard, Amex, Discover)",
		},
		{
			ID:          "ip_address",
			Name:        "IP Address",
			Type:        models.PIITypeIPAddress,
			Category:    models.PIICategoryTechnical,
			Pattern:     `\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`,
			Sensitivity: models.SensitivityLow,
			Confidence:  0.8,
			Enabled:     true,
			Description: "IPv4 address pattern",
		},
		{
			ID:          "date_of_birth",
			Name:        "Date of Birth",
			Type:        models.PIITypeDateOfBirth,
			Category:    models.PIICategoryPersonal,
			Pattern:     `\b(?:0[1-9]|1[0-2])[-/](?:0[1-9]|[12][0-9]|3[01])[-/](?:19|20)\d{2}\b`,
			Sensitivity: models.SensitivityHigh,
			Confidence:  0.7,
			Enabled:     true,
			Description: "Date of birth pattern (MM/DD/YYYY or MM-DD-YYYY)",
		},
		{
			ID:          "passport",
			Name:        "Passport Number",
			Type:        models.PIITypePassport,
			Category:    models.PIICategoryIdentifier,
			Pattern:     `\b[A-Z]{1,2}[0-9]{6,9}\b`,
			Sensitivity: models.SensitivityHigh,
			Confidence:  0.6,
			Enabled:     true,
			Description: "Passport number pattern",
		},
		{
			ID:          "driver_license",
			Name:        "Driver License",
			Type:        models.PIITypeDriverLicense,
			Category:    models.PIICategoryIdentifier,
			Pattern:     `\b[A-Z]{1,2}[0-9]{6,8}\b`,
			Sensitivity: models.SensitivityHigh,
			Confidence:  0.5,
			Enabled:     true,
			Description: "Driver license number pattern",
		},
	}

	for _, pattern := range defaultPatterns {
		pattern.CreatedAt = time.Now()
		pattern.UpdatedAt = time.Now()
		s.patterns[pattern.ID] = &pattern

		// Compile regex
		if regex, err := regexp.Compile(pattern.Pattern); err == nil {
			s.compiledRegex[pattern.ID] = regex
		} else {
			s.logger.Warn("Failed to compile PII pattern regex", "pattern_id", pattern.ID, "error", err)
		}
	}
}

func (s *PIIDetectionService) DetectPII(ctx context.Context, request *models.PIIDetectionRequest) (*models.PIIDetectionResult, error) {
	startTime := time.Now()

	result := &models.PIIDetectionResult{
		RequestID:      request.RequestID,
		PIIFindings:    []models.PIIFinding{},
		TotalPatterns:  len(s.patterns),
		MatchedCount:   0,
		ProcessingTime: 0,
		RiskScore:      0.0,
		ComplianceIssues: []models.ComplianceIssue{},
		Recommendations:  []string{},
		ScannedAt:       time.Now(),
	}

	// Scan the data for PII
	findings, err := s.ScanData(ctx, request.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to scan data for PII: %w", err)
	}

	result.PIIFindings = findings
	result.MatchedCount = len(findings)
	result.ProcessingTime = time.Since(startTime)

	// Calculate risk score based on findings
	result.RiskScore = s.calculateRiskScore(findings)

	// Check compliance if regulations specified
	if len(request.Regulations) > 0 {
		complianceResult, err := s.ValidateDataCompliance(ctx, request.Data, request.Regulations)
		if err != nil {
			s.logger.Warn("Failed to validate compliance", "error", err)
		} else {
			result.ComplianceIssues = complianceResult.Issues
			result.Recommendations = append(result.Recommendations, complianceResult.Recommendations...)
		}
	}

	// Generate general recommendations
	result.Recommendations = append(result.Recommendations, s.generateRecommendations(findings)...)

	// Store findings if significant PII detected
	if result.RiskScore > 5.0 {
		for _, finding := range findings {
			piiData := &models.PIIData{
				ID:           uuid.New().String(),
				RequestID:    request.RequestID,
				DataType:     finding.Type,
				Category:     finding.Category,
				Sensitivity:  finding.Sensitivity,
				Location:     finding.Location,
				Value:        finding.MaskedValue,
				RiskScore:    finding.RiskScore,
				Confidence:   finding.Confidence,
				APIID:        request.APIID,
				EndpointID:   request.EndpointID,
				IPAddress:    request.IPAddress,
				UserAgent:    request.UserAgent,
				DetectedAt:   time.Now(),
				Metadata: map[string]interface{}{
					"pattern_id":     finding.PatternID,
					"pattern_name":   finding.PatternName,
					"field_name":     finding.FieldName,
					"data_source":    request.DataSource,
					"regulations":    request.Regulations,
				},
			}

			if err := s.piiRepo.CreatePIIData(ctx, piiData); err != nil {
				s.logger.Error("Failed to store PII data", "pii_id", piiData.ID, "error", err)
			}
		}

		// Publish PII detection events
		if err := s.publishPIIEvents(ctx, findings, request); err != nil {
			s.logger.Error("Failed to publish PII events", "error", err)
		}
	}

	return result, nil
}

func (s *PIIDetectionService) ScanData(ctx context.Context, data map[string]interface{}) ([]models.PIIFinding, error) {
	var findings []models.PIIFinding

	// Recursively scan all data fields
	findings = append(findings, s.scanDataRecursive(data, "")...)

	return findings, nil
}

func (s *PIIDetectionService) scanDataRecursive(data interface{}, path string) []models.PIIFinding {
	var findings []models.PIIFinding

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			currentPath := key
			if path != "" {
				currentPath = path + "." + key
			}
			findings = append(findings, s.scanDataRecursive(value, currentPath)...)
		}
	case []interface{}:
		for i, item := range v {
			currentPath := fmt.Sprintf("%s[%d]", path, i)
			findings = append(findings, s.scanDataRecursive(item, currentPath)...)
		}
	case string:
		// Scan string values for PII patterns
		findings = append(findings, s.scanStringValue(v, path)...)
	}

	return findings
}

func (s *PIIDetectionService) scanStringValue(value, location string) []models.PIIFinding {
	var findings []models.PIIFinding

	for patternID, pattern := range s.patterns {
		if !pattern.Enabled {
			continue
		}

		regex, exists := s.compiledRegex[patternID]
		if !exists {
			continue
		}

		matches := regex.FindAllString(value, -1)
		for _, match := range matches {
			// Additional validation for certain PII types
			if s.validatePIIMatch(pattern.Type, match) {
				finding := models.PIIFinding{
					ID:           uuid.New().String(),
					PatternID:    patternID,
					PatternName:  pattern.Name,
					Type:         pattern.Type,
					Category:     pattern.Category,
					Sensitivity:  pattern.Sensitivity,
					Location:     location,
					FieldName:    s.extractFieldName(location),
					Value:        match,
					MaskedValue:  s.maskValue(match, pattern.Type),
					RiskScore:    s.calculateFindingRiskScore(pattern, match),
					Confidence:   pattern.Confidence,
					DetectedAt:   time.Now(),
					Metadata: map[string]interface{}{
						"pattern_description": pattern.Description,
						"match_length":        len(match),
						"field_context":       s.getFieldContext(location),
					},
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

func (s *PIIDetectionService) validatePIIMatch(piiType, value string) bool {
	switch piiType {
	case models.PIITypeCreditCard:
		return s.validateCreditCard(value)
	case models.PIITypeSSN:
		return s.validateSSN(value)
	case models.PIITypeEmail:
		return s.validateEmail(value)
	case models.PIITypePhone:
		return s.validatePhone(value)
	default:
		return true // No additional validation needed
	}
}

func (s *PIIDetectionService) validateCreditCard(number string) bool {
	// Remove non-digits
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(number, "")
	
	// Check length
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return false
	}

	// Luhn algorithm validation
	return s.luhnCheck(cleaned)
}

func (s *PIIDetectionService) luhnCheck(number string) bool {
	var sum int
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		
		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

func (s *PIIDetectionService) validateSSN(ssn string) bool {
	// Remove non-digits
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(ssn, "")
	
	// Check length
	if len(cleaned) != 9 {
		return false
	}

	// Check for invalid patterns
	invalidPatterns := []string{
		"000000000", "111111111", "222222222", "333333333",
		"444444444", "555555555", "666666666", "777777777",
		"888888888", "999999999",
	}

	for _, invalid := range invalidPatterns {
		if cleaned == invalid {
			return false
		}
	}

	// Check area number (first 3 digits)
	areaNumber := cleaned[:3]
	if areaNumber == "000" || areaNumber == "666" || areaNumber[0] == '9' {
		return false
	}

	// Check group number (middle 2 digits)
	groupNumber := cleaned[3:5]
	if groupNumber == "00" {
		return false
	}

	// Check serial number (last 4 digits)
	serialNumber := cleaned[5:]
	if serialNumber == "0000" {
		return false
	}

	return true
}

func (s *PIIDetectionService) validateEmail(email string) bool {
	// Basic email validation beyond regex
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]
	
	// Check local part
	if len(local) == 0 || len(local) > 64 {
		return false
	}

	// Check domain part
	if len(domain) == 0 || len(domain) > 255 {
		return false
	}

	// Domain must contain at least one dot
	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

func (s *PIIDetectionService) validatePhone(phone string) bool {
	// Remove non-digits
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// US phone numbers should be 10 or 11 digits (with country code)
	if len(cleaned) == 10 {
		return true
	}
	if len(cleaned) == 11 && cleaned[0] == '1' {
		return true
	}

	return false
}

func (s *PIIDetectionService) maskValue(value, piiType string) string {
	switch piiType {
	case models.PIITypeSSN:
		if len(value) >= 4 {
			return "***-**-" + value[len(value)-4:]
		}
	case models.PIITypeCreditCard:
		cleaned := regexp.MustCompile(`\D`).ReplaceAllString(value, "")
		if len(cleaned) >= 4 {
			return "**** **** **** " + cleaned[len(cleaned)-4:]
		}
	case models.PIITypeEmail:
		parts := strings.Split(value, "@")
		if len(parts) == 2 && len(parts[0]) > 2 {
			return parts[0][:2] + "***@" + parts[1]
		}
	case models.PIITypePhone:
		cleaned := regexp.MustCompile(`\D`).ReplaceAllString(value, "")
		if len(cleaned) >= 4 {
			return "***-***-" + cleaned[len(cleaned)-4:]
		}
	default:
		if len(value) > 4 {
			return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
		}
	}

	return strings.Repeat("*", len(value))
}

func (s *PIIDetectionService) calculateFindingRiskScore(pattern *models.PIIPattern, value string) float64 {
	baseScore := 0.0

	// Base score from sensitivity level
	switch pattern.Sensitivity {
	case models.SensitivityHigh:
		baseScore = 8.0
	case models.SensitivityMedium:
		baseScore = 5.0
	case models.SensitivityLow:
		baseScore = 2.0
	}

	// Adjust based on confidence
	baseScore *= pattern.Confidence

	// Adjust based on PII type criticality
	switch pattern.Type {
	case models.PIITypeSSN, models.PIITypeCreditCard, models.PIITypePassport:
		baseScore *= 1.2
	case models.PIITypeDateOfBirth, models.PIITypeDriverLicense:
		baseScore *= 1.1
	case models.PIITypeEmail, models.PIITypePhone:
		baseScore *= 0.9
	case models.PIITypeIPAddress:
		baseScore *= 0.7
	}

	// Cap at 10.0
	if baseScore > 10.0 {
		baseScore = 10.0
	}

	return baseScore
}

func (s *PIIDetectionService) calculateRiskScore(findings []models.PIIFinding) float64 {
	if len(findings) == 0 {
		return 0.0
	}

	totalScore := 0.0
	highSensitivityCount := 0
	uniqueTypes := make(map[string]bool)

	for _, finding := range findings {
		totalScore += finding.RiskScore
		uniqueTypes[finding.Type] = true
		
		if finding.Sensitivity == models.SensitivityHigh {
			highSensitivityCount++
		}
	}

	// Average risk score
	avgScore := totalScore / float64(len(findings))

	// Boost score for multiple high-sensitivity findings
	if highSensitivityCount > 1 {
		avgScore *= 1.2
	}

	// Boost score for diverse PII types
	if len(uniqueTypes) > 3 {
		avgScore *= 1.1
	}

	// Cap at 10.0
	if avgScore > 10.0 {
		avgScore = 10.0
	}

	return avgScore
}

func (s *PIIDetectionService) extractFieldName(location string) string {
	parts := strings.Split(location, ".")
	if len(parts) > 0 {
		// Remove array indices
		fieldName := parts[len(parts)-1]
		return regexp.MustCompile(`\[\d+\]`).ReplaceAllString(fieldName, "")
	}
	return location
}

func (s *PIIDetectionService) getFieldContext(location string) string {
	fieldName := strings.ToLower(s.extractFieldName(location))
	
	// Determine context based on field name
	contexts := map[string]string{
		"email":     "contact_information",
		"phone":     "contact_information",
		"address":   "address_information",
		"name":      "personal_information",
		"ssn":       "identification",
		"id":        "identification",
		"card":      "financial_information",
		"account":   "financial_information",
		"password":  "authentication",
		"token":     "authentication",
	}

	for keyword, context := range contexts {
		if strings.Contains(fieldName, keyword) {
			return context
		}
	}

	return "general"
}

func (s *PIIDetectionService) generateRecommendations(findings []models.PIIFinding) []string {
	var recommendations []string
	
	if len(findings) == 0 {
		return []string{"No PII detected - continue monitoring"}
	}

	// General recommendations
	recommendations = append(recommendations, "Implement data encryption for sensitive fields")
	recommendations = append(recommendations, "Consider data masking for non-production environments")
	recommendations = append(recommendations, "Review data retention policies")

	// Type-specific recommendations
	typeCount := make(map[string]int)
	for _, finding := range findings {
		typeCount[finding.Type]++
	}

	if typeCount[models.PIITypeCreditCard] > 0 {
		recommendations = append(recommendations, "Ensure PCI DSS compliance for credit card data")
		recommendations = append(recommendations, "Implement tokenization for credit card numbers")
	}

	if typeCount[models.PIITypeSSN] > 0 {
		recommendations = append(recommendations, "Restrict access to SSN data")
		recommendations = append(recommendations, "Implement additional authentication for SSN access")
	}

	if typeCount[models.PIITypeEmail] > 0 {
		recommendations = append(recommendations, "Consider email hashing for analytics")
		recommendations = append(recommendations, "Implement opt-out mechanisms for email data")
	}

	// Sensitivity-based recommendations
	highSensitivityCount := 0
	for _, finding := range findings {
		if finding.Sensitivity == models.SensitivityHigh {
			highSensitivityCount++
		}
	}

	if highSensitivityCount > 2 {
		recommendations = append(recommendations, "High-sensitivity data detected - implement enhanced security controls")
		recommendations = append(recommendations, "Consider data loss prevention (DLP) solutions")
	}

	return recommendations
}

func (s *PIIDetectionService) GetPIIPatterns(ctx context.Context) ([]models.PIIPattern, error) {
	var patterns []models.PIIPattern
	for _, pattern := range s.patterns {
		patterns = append(patterns, *pattern)
	}
	return patterns, nil
}

func (s *PIIDetectionService) CreateCustomPattern(ctx context.Context, pattern *models.PIIPattern) error {
	// Validate pattern
	if err := s.validatePattern(pattern); err != nil {
		return fmt.Errorf("pattern validation failed: %w", err)
	}

	// Set metadata
	pattern.ID = uuid.New().String()
	pattern.CreatedAt = time.Now()
	pattern.UpdatedAt = time.Now()

	// Test regex compilation
	if _, err := regexp.Compile(pattern.Pattern); err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Store in repository
	if err := s.piiRepo.CreatePIIPattern(ctx, pattern); err != nil {
		return fmt.Errorf("failed to create pattern: %w", err)
	}

	// Add to memory
	s.patterns[pattern.ID] = pattern
	if regex, err := regexp.Compile(pattern.Pattern); err == nil {
		s.compiledRegex[pattern.ID] = regex
	}

	s.logger.Info("Created custom PII pattern", "pattern_id", pattern.ID, "name", pattern.Name)
	return nil
}

func (s *PIIDetectionService) UpdatePattern(ctx context.Context, patternID string, pattern *models.PIIPattern) error {
	// Validate pattern
	if err := s.validatePattern(pattern); err != nil {
		return fmt.Errorf("pattern validation failed: %w", err)
	}

	// Test regex compilation
	if _, err := regexp.Compile(pattern.Pattern); err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Update in repository
	pattern.UpdatedAt = time.Now()
	if err := s.piiRepo.UpdatePIIPattern(ctx, patternID, pattern); err != nil {
		return fmt.Errorf("failed to update pattern: %w", err)
	}

	// Update in memory
	pattern.ID = patternID
	s.patterns[patternID] = pattern
	if regex, err := regexp.Compile(pattern.Pattern); err == nil {
		s.compiledRegex[patternID] = regex
	} else {
		delete(s.compiledRegex, patternID)
	}

	s.logger.Info("Updated PII pattern", "pattern_id", patternID, "name", pattern.Name)
	return nil
}

func (s *PIIDetectionService) validatePattern(pattern *models.PIIPattern) error {
	if pattern.Name == "" {
		return fmt.Errorf("pattern name is required")
	}

	if pattern.Pattern == "" {
		return fmt.Errorf("pattern regex is required")
	}

	if pattern.Type == "" {
		return fmt.Errorf("pattern type is required")
	}

	if pattern.Category == "" {
		return fmt.Errorf("pattern category is required")
	}

	if pattern.Sensitivity == "" {
		return fmt.Errorf("pattern sensitivity is required")
	}

	validSensitivities := map[string]bool{
		models.SensitivityHigh:   true,
		models.SensitivityMedium: true,
		models.SensitivityLow:    true,
	}

	if !validSensitivities[pattern.Sensitivity] {
		return fmt.Errorf("invalid sensitivity level: %s", pattern.Sensitivity)
	}

	if pattern.Confidence < 0.0 || pattern.Confidence > 1.0 {
		return fmt.Errorf("confidence must be between 0.0 and 1.0")
	}

	return nil
}

func (s *PIIDetectionService) ValidateDataCompliance(ctx context.Context, data map[string]interface{}, regulations []string) (*models.ComplianceValidationResult, error) {
	result := &models.ComplianceValidationResult{
		Regulations:     regulations,
		Issues:          []models.ComplianceIssue{},
		Recommendations: []string{},
		OverallStatus:   models.ComplianceStatusCompliant,
		ValidatedAt:     time.Now(),
	}

	// Scan for PII first
	findings, err := s.ScanData(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to scan data for compliance validation: %w", err)
	}

	// Check each regulation
	for _, regulation := range regulations {
		issues := s.checkRegulationCompliance(regulation, findings, data)
		result.Issues = append(result.Issues, issues...)
	}

	// Determine overall status
	if len(result.Issues) > 0 {
		hasViolations := false
		for _, issue := range result.Issues {
			if issue.Severity == models.ComplianceSeverityHigh || issue.Severity == models.ComplianceSeverityCritical {
				hasViolations = true
				break
			}
		}
		
		if hasViolations {
			result.OverallStatus = models.ComplianceStatusViolation
		} else {
			result.OverallStatus = models.ComplianceStatusWarning
		}
	}

	// Generate recommendations
	result.Recommendations = s.generateComplianceRecommendations(result.Issues, regulations)

	return result, nil
}

func (s *PIIDetectionService) checkRegulationCompliance(regulation string, findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	switch strings.ToUpper(regulation) {
	case "GDPR":
		issues = append(issues, s.checkGDPRCompliance(findings, data)...)
	case "CCPA":
		issues = append(issues, s.checkCCPACompliance(findings, data)...)
	case "HIPAA":
		issues = append(issues, s.checkHIPAACompliance(findings, data)...)
	case "PCI_DSS":
		issues = append(issues, s.checkPCIDSSCompliance(findings, data)...)
	case "SOX":
		issues = append(issues, s.checkSOXCompliance(findings, data)...)
	}

	return issues
}

func (s *PIIDetectionService) checkGDPRCompliance(findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	// Check for personal data without consent indicators
	personalDataTypes := []string{
		models.PIITypeEmail,
		models.PIITypePhone,
		models.PIITypeDateOfBirth,
		models.PIITypePassport,
		models.PIITypeDriverLicense,
	}

	hasPersonalData := false
	for _, finding := range findings {
		for _, pdType := range personalDataTypes {
			if finding.Type == pdType {
				hasPersonalData = true
				break
			}
		}
	}

	if hasPersonalData {
		// Check for consent indicators
		consentFields := []string{"consent", "agreement", "opt_in", "permission"}
		hasConsent := false
		
		for field := range data {
			fieldLower := strings.ToLower(field)
			for _, consentField := range consentFields {
				if strings.Contains(fieldLower, consentField) {
					hasConsent = true
					break
				}
			}
		}

		if !hasConsent {
			issues = append(issues, models.ComplianceIssue{
				ID:          uuid.New().String(),
				Regulation:  "GDPR",
				Article:     "Article 6",
				Title:       "Missing Consent for Personal Data Processing",
				Description: "Personal data detected without clear consent indicators",
				Severity:    models.ComplianceSeverityHigh,
				Category:    models.ComplianceCategoryConsent,
				DetectedAt:  time.Now(),
			})
		}

		// Check for data minimization
		if len(findings) > 5 {
			issues = append(issues, models.ComplianceIssue{
				ID:          uuid.New().String(),
				Regulation:  "GDPR",
				Article:     "Article 5(1)(c)",
				Title:       "Data Minimization Concern",
				Description: "Large amount of personal data detected - review data minimization principles",
				Severity:    models.ComplianceSeverityMedium,
				Category:    models.ComplianceCategoryDataMinimization,
				DetectedAt:  time.Now(),
			})
		}
	}

	return issues
}

func (s *PIIDetectionService) checkCCPACompliance(findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	// Check for personal information categories under CCPA
	ccpaCategories := map[string]bool{
		models.PIITypeEmail:         true,
		models.PIITypePhone:         true,
		models.PIITypeSSN:           true,
		models.PIITypeDriverLicense: true,
		models.PIITypeIPAddress:     true,
	}

	hasCCPAData := false
	for _, finding := range findings {
		if ccpaCategories[finding.Type] {
			hasCCPAData = true
			break
		}
	}

	if hasCCPAData {
		// Check for privacy notice indicators
		privacyFields := []string{"privacy_notice", "privacy_policy", "data_usage", "opt_out"}
		hasPrivacyNotice := false
		
		for field := range data {
			fieldLower := strings.ToLower(field)
			for _, privacyField := range privacyFields {
				if strings.Contains(fieldLower, privacyField) {
					hasPrivacyNotice = true
					break
				}
			}
		}

		if !hasPrivacyNotice {
			issues = append(issues, models.ComplianceIssue{
				ID:          uuid.New().String(),
				Regulation:  "CCPA",
				Article:     "Section 1798.100",
				Title:       "Missing Privacy Notice",
				Description: "Personal information detected without privacy notice indicators",
				Severity:    models.ComplianceSeverityHigh,
				Category:    models.ComplianceCategoryPrivacyNotice,
				DetectedAt:  time.Now(),
			})
		}
	}

	return issues
}

func (s *PIIDetectionService) checkHIPAACompliance(findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	// Check for PHI (Protected Health Information)
	phiIndicators := []string{"medical", "health", "patient", "diagnosis", "treatment", "medication"}
	hasPHI := false

	// Check field names for health-related terms
	for field := range data {
		fieldLower := strings.ToLower(field)
		for _, indicator := range phiIndicators {
			if strings.Contains(fieldLower, indicator) {
				hasPHI = true
				break
			}
		}
	}

	// Check for identifiers that could be PHI when combined with health data
	if hasPHI {
		for _, finding := range findings {
			if finding.Type == models.PIITypeSSN || finding.Type == models.PIITypeEmail || finding.Type == models.PIITypePhone {
				issues = append(issues, models.ComplianceIssue{
					ID:          uuid.New().String(),
					Regulation:  "HIPAA",
					Article:     "45 CFR 164.514",
					Title:       "PHI Identifier Detected",
					Description: fmt.Sprintf("Health data combined with %s may constitute PHI", finding.Type),
					Severity:    models.ComplianceSeverityHigh,
					Category:    models.ComplianceCategoryPHI,
					DetectedAt:  time.Now(),
				})
			}
		}
	}

	return issues
}

func (s *PIIDetectionService) checkPCIDSSCompliance(findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	// Check for credit card data
	for _, finding := range findings {
		if finding.Type == models.PIITypeCreditCard {
			issues = append(issues, models.ComplianceIssue{
				ID:          uuid.New().String(),
				Regulation:  "PCI DSS",
				Article:     "Requirement 3",
				Title:       "Cardholder Data Detected",
				Description: "Credit card number detected - ensure PCI DSS compliance",
				Severity:    models.ComplianceSeverityCritical,
				Category:    models.ComplianceCategoryCardholderData,
				DetectedAt:  time.Now(),
			})
		}
	}

	return issues
}

func (s *PIIDetectionService) checkSOXCompliance(findings []models.PIIFinding, data map[string]interface{}) []models.ComplianceIssue {
	var issues []models.ComplianceIssue

	// Check for financial data indicators
	financialFields := []string{"financial", "accounting", "revenue", "expense", "audit", "transaction"}
	hasFinancialData := false

	for field := range data {
		fieldLower := strings.ToLower(field)
		for _, finField := range financialFields {
			if strings.Contains(fieldLower, finField) {
				hasFinancialData = true
				break
			}
		}
	}

	if hasFinancialData {
		// Check for audit trail indicators
		auditFields := []string{"audit_trail", "timestamp", "user_id", "action", "change_log"}
		hasAuditTrail := false

		for field := range data {
			fieldLower := strings.ToLower(field)
			for _, auditField := range auditFields {
				if strings.Contains(fieldLower, auditField) {
					hasAuditTrail = true
					break
				}
			}
		}

		if !hasAuditTrail {
			issues = append(issues, models.ComplianceIssue{
				ID:          uuid.New().String(),
				Regulation:  "SOX",
				Article:     "Section 404",
				Title:       "Missing Audit Trail for Financial Data",
				Description: "Financial data detected without proper audit trail indicators",
				Severity:    models.ComplianceSeverityHigh,
				Category:    models.ComplianceCategoryAuditTrail,
				DetectedAt:  time.Now(),
			})
		}
	}

	return issues
}

func (s *PIIDetectionService) generateComplianceRecommendations(issues []models.ComplianceIssue, regulations []string) []string {
	var recommendations []string

	if len(issues) == 0 {
		return []string{"Data appears compliant with specified regulations"}
	}

	// General recommendations
	recommendations = append(recommendations, "Implement data governance policies")
	recommendations = append(recommendations, "Regular compliance audits recommended")

	// Regulation-specific recommendations
	regulationMap := make(map[string]bool)
	for _, regulation := range regulations {
		regulationMap[regulation] = true
	}

	if regulationMap["GDPR"] {
		recommendations = append(recommendations, "Implement consent management system")
		recommendations = append(recommendations, "Establish data subject rights procedures")
		recommendations = append(recommendations, "Conduct Data Protection Impact Assessment (DPIA)")
	}

	if regulationMap["CCPA"] {
		recommendations = append(recommendations, "Implement consumer rights request handling")
		recommendations = append(recommendations, "Update privacy policy with CCPA disclosures")
	}

	if regulationMap["HIPAA"] {
		recommendations = append(recommendations, "Implement HIPAA security safeguards")
		recommendations = append(recommendations, "Conduct risk assessment for PHI")
		recommendations = append(recommendations, "Establish business associate agreements")
	}

	if regulationMap["PCI_DSS"] {
		recommendations = append(recommendations, "Implement PCI DSS security controls")
		recommendations = append(recommendations, "Use tokenization for card data")
		recommendations = append(recommendations, "Regular PCI compliance validation")
	}

	// Issue-specific recommendations
	categoryCount := make(map[string]int)
	for _, issue := range issues {
		categoryCount[issue.Category]++
	}

	if categoryCount[models.ComplianceCategoryConsent] > 0 {
		recommendations = append(recommendations, "Implement explicit consent mechanisms")
	}

	if categoryCount[models.ComplianceCategoryCardholderData] > 0 {
		recommendations = append(recommendations, "Encrypt cardholder data at rest and in transit")
	}

	return recommendations
}

func (s *PIIDetectionService) publishPIIEvents(ctx context.Context, findings []models.PIIFinding, request *models.PIIDetectionRequest) error {
	for _, finding := range findings {
		eventData := map[string]interface{}{
			"event_type":     "pii_detected",
			"finding_id":     finding.ID,
			"pattern_id":     finding.PatternID,
			"pattern_name":   finding.PatternName,
			"pii_type":       finding.Type,
			"category":       finding.Category,
			"sensitivity":    finding.Sensitivity,
			"location":       finding.Location,
			"field_name":     finding.FieldName,
			"masked_value":   finding.MaskedValue,
			"risk_score":     finding.RiskScore,
			"confidence":     finding.Confidence,
			"request_id":     request.RequestID,
			"api_id":         request.APIID,
			"endpoint_id":    request.EndpointID,
			"ip_address":     request.IPAddress,
			"user_agent":     request.UserAgent,
			"data_source":    request.DataSource,
			"regulations":    request.Regulations,
			"timestamp":      finding.DetectedAt,
			"metadata":       finding.Metadata,
		}

		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal PII event: %w", err)
		}

		message := kafka.Message{
			Topic: "pii_events",
			Key:   finding.ID,
			Value: eventJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce PII event: %w", err)
		}
	}

	return nil
}


