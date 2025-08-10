package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/threat-detection/internal/models"
	"scopeapi.local/backend/services/threat-detection/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type ThreatDetectionServiceInterface interface {
	AnalyzeTraffic(ctx context.Context, trafficData []byte) (*models.ThreatAnalysisResult, error)
	AnalyzeThreat(ctx context.Context, request *models.ThreatAnalysisRequest) (*models.ThreatAnalysisResult, error)
	GetThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error)
	GetThreat(ctx context.Context, threatID string) (*models.Threat, error)
	UpdateThreatStatus(ctx context.Context, threatID string, update *models.ThreatUpdateRequest) error
	DeleteThreat(ctx context.Context, threatID string) error
	GetThreatStatistics(ctx context.Context, timeRange time.Duration) (*models.ThreatStatistics, error)
	CreateThreat(ctx context.Context, threat *models.Threat) error
	ProcessSecurityEvent(ctx context.Context, event map[string]interface{}) error
}

type ThreatDetectionService struct {
	threatRepo    repository.ThreatRepositoryInterface
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
}

func NewThreatDetectionService(
	threatRepo repository.ThreatRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *ThreatDetectionService {
	return &ThreatDetectionService{
		threatRepo:    threatRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (s *ThreatDetectionService) AnalyzeTraffic(ctx context.Context, trafficData []byte) (*models.ThreatAnalysisResult, error) {
	startTime := time.Now()

	// Parse traffic data
	var traffic map[string]interface{}
	if err := json.Unmarshal(trafficData, &traffic); err != nil {
		return nil, fmt.Errorf("failed to parse traffic data: %w", err)
	}

	requestID := uuid.New().String()
	result := &models.ThreatAnalysisResult{
		RequestID:      requestID,
		ThreatDetected: false,
		Confidence:     0.0,
		RiskScore:      0.0,
		Indicators:     []models.ThreatIndicator{},
		Recommendations: []string{},
		Metadata:       make(map[string]interface{}),
		AnalyzedAt:     time.Now(),
	}

	// Perform various threat detection analyses
	threats := []models.Threat{}

	// 1. SQL Injection Detection
	sqlInjectionThreats, err := s.detectSQLInjection(ctx, traffic)
	if err != nil {
		s.logger.Error("SQL injection detection failed", "error", err)
	} else {
		threats = append(threats, sqlInjectionThreats...)
	}

	// 2. XSS Detection
	xssThreats, err := s.detectXSS(ctx, traffic)
	if err != nil {
		s.logger.Error("XSS detection failed", "error", err)
	} else {
		threats = append(threats, xssThreats...)
	}

	// 3. DDoS Detection
	ddosThreats, err := s.detectDDoS(ctx, traffic)
	if err != nil {
		s.logger.Error("DDoS detection failed", "error", err)
	} else {
		threats = append(threats, ddosThreats...)
	}

	// 4. Brute Force Detection
	bruteForceThreats, err := s.detectBruteForce(ctx, traffic)
	if err != nil {
		s.logger.Error("Brute force detection failed", "error", err)
	} else {
		threats = append(threats, bruteForceThreats...)
	}

	// 5. Data Exfiltration Detection
	dataExfilThreats, err := s.detectDataExfiltration(ctx, traffic)
	if err != nil {
		s.logger.Error("Data exfiltration detection failed", "error", err)
	} else {
		threats = append(threats, dataExfilThreats...)
	}

	// Process detected threats
	if len(threats) > 0 {
		result.ThreatDetected = true
		
		// Calculate overall risk score and confidence
		totalRiskScore := 0.0
		totalConfidence := 0.0
		highestSeverity := ""
		
		for _, threat := range threats {
			totalRiskScore += threat.RiskScore
			totalConfidence += threat.Confidence
			
			// Determine highest severity
			if s.getSeverityWeight(threat.Severity) > s.getSeverityWeight(highestSeverity) {
				highestSeverity = threat.Severity
				result.ThreatType = threat.Type
			}
			
			// Add indicators
			result.Indicators = append(result.Indicators, threat.Indicators...)
		}
		
		result.RiskScore = totalRiskScore / float64(len(threats))
		result.Confidence = totalConfidence / float64(len(threats))
		result.Severity = highestSeverity
		
		// Store threats in database
		for _, threat := range threats {
			if err := s.CreateThreat(ctx, &threat); err != nil {
				s.logger.Error("Failed to store threat", "threat_id", threat.ID, "error", err)
			}
		}
		
		// Generate recommendations
		result.Recommendations = s.generateRecommendations(threats)
		
		// Publish threat event to Kafka
		if err := s.publishThreatEvent(ctx, threats); err != nil {
			s.logger.Error("Failed to publish threat event", "error", err)
		}
	}

	result.ProcessingTime = time.Since(startTime)
	result.Metadata["threats_analyzed"] = len(threats)
	result.Metadata["analysis_methods"] = []string{"sql_injection", "xss", "ddos", "brute_force", "data_exfiltration"}

	return result, nil
}

func (s *ThreatDetectionService) detectSQLInjection(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat
	
	// Extract request data
	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}
	
	// Check URL parameters
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsSQLInjectionPattern(fmt.Sprintf("%v", value)) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            models.ThreatTypeInjection,
					Severity:        models.ThreatSeverityHigh,
					Status:          models.ThreatStatusNew,
					Title:           "SQL Injection Attempt Detected",
					Description:     fmt.Sprintf("Potential SQL injection detected in parameter '%s'", key),
					DetectionMethod: models.DetectionMethodSignature,
					Confidence:      0.85,
					RiskScore:       8.5,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "sql_injection_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: "Suspicious SQL pattern detected",
							Severity:    models.ThreatSeverityHigh,
							Confidence:  0.85,
						},
					},
					RequestData:  requestData,
					FirstSeen:    time.Now(),
					LastSeen:     time.Now(),
					Count:        1,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				
				// Extract additional context
				if ipAddr, ok := traffic["ip_address"].(string); ok {
					threat.IPAddress = ipAddr
				}
				if userAgent, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = userAgent
				}
				if apiID, ok := traffic["api_id"].(string); ok {
					threat.APIID = apiID
				}
				
				threats = append(threats, threat)
			}
		}
	}
	
	// Check request body
	if body, ok := requestData["body"].(string); ok {
		if s.containsSQLInjectionPattern(body) {
			threat := models.Threat{
				ID:              uuid.New().String(),
				Type:            models.ThreatTypeInjection,
				Severity:        models.ThreatSeverityHigh,
				Status:          models.ThreatStatusNew,
				Title:           "SQL Injection in Request Body",
				Description:     "Potential SQL injection detected in request body",
				DetectionMethod: models.DetectionMethodSignature,
				Confidence:      0.80,
				RiskScore:       8.0,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "sql_injection_body",
						Value:       body,
						Description: "SQL injection pattern in request body",
						Severity:    models.ThreatSeverityHigh,
						Confidence:  0.80,
					},
				},
				RequestData: requestData,
				FirstSeen:   time.Now(),
				LastSeen:    time.Now(),
				Count:       1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			
			// Extract additional context
			if ipAddr, ok := traffic["ip_address"].(string); ok {
				threat.IPAddress = ipAddr
			}
			if userAgent, ok := requestData["user_agent"].(string); ok {
				threat.UserAgent = userAgent
			}
			if apiID, ok := traffic["api_id"].(string); ok {
				threat.APIID = apiID
			}
			
			threats = append(threats, threat)
		}
	}
	
	return threats, nil
}

func (s *ThreatDetectionService) containsSQLInjectionPattern(input string) bool {
	// Common SQL injection patterns
	patterns := []string{
		"'",
		"\"",
		";",
		"--",
		"/*",
		"*/",
		"xp_",
		"sp_",
		"union",
		"select",
		"insert",
		"delete",
		"update",
		"drop",
		"create",
		"alter",
		"exec",
		"execute",
		"script",
		"javascript",
		"vbscript",
		"onload",
		"onerror",
		"<script",
		"</script>",
	}
	
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}
	
	return false
}

func (s *ThreatDetectionService) detectXSS(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat
	
	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}
	
	// Check parameters for XSS patterns
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsXSSPattern(fmt.Sprintf("%v", value)) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            models.ThreatTypeXSS,
					Severity:        models.ThreatSeverityMedium,
					Status:          models.ThreatStatusNew,
					Title:           "Cross-Site Scripting (XSS) Attempt",
					Description:     fmt.Sprintf("Potential XSS attack detected in parameter '%s'", key),
					DetectionMethod: models.DetectionMethodSignature,
					Confidence:      0.75,
					RiskScore:       6.5,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "xss_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: "Suspicious script pattern detected",
							Severity:    models.ThreatSeverityMedium,
							Confidence:  0.75,
						},
					},
					RequestData: requestData,
					FirstSeen:   time.Now(),
					LastSeen:    time.Now(),
					Count:       1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				
				// Extract additional context
				if ipAddr, ok := traffic["ip_address"].(string); ok {
					threat.IPAddress = ipAddr
				}
				if userAgent, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = userAgent
				}
				if apiID, ok := traffic["api_id"].(string); ok {
					threat.APIID = apiID
				}
				
				threats = append(threats, threat)
			}
		}
	}
	
	return threats, nil
}

func (s *ThreatDetectionService) containsXSSPattern(input string) bool {
	patterns := []string{
		"<script",
		"</script>",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
		"onmouseover=",
		"onfocus=",
		"onblur=",
		"alert(",
		"confirm(",
		"prompt(",
		"document.cookie",
		"document.write",
		"window.location",
		"eval(",
		"setTimeout(",
		"setInterval(",
		"<iframe",
		"<object",
		"<embed",
		"<form",
		"<img",
		"<svg",
	}
	
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}
	
	return false
}

func (s *ThreatDetectionService) detectDDoS(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat
	
	// Extract IP address and timestamp
	ipAddr, ok := traffic["ip_address"].(string)
	if !ok {
		return threats, nil
	}
	
	timestamp, ok := traffic["timestamp"].(time.Time)
	if !ok {
		return threats, nil
	}
	
	// Check request rate from this IP in the last minute
	requestCount, err := s.threatRepo.GetRequestCountByIP(ctx, ipAddr, time.Minute)
	if err != nil {
		return threats, fmt.Errorf("failed to get request count: %w", err)
	}
	
	// DDoS threshold: more than 100 requests per minute from single IP
	if requestCount > 100 {
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            models.ThreatTypeDDoS,
			Severity:        models.ThreatSeverityHigh,
			Status:          models.ThreatStatusNew,
			Title:           "Potential DDoS Attack",
			Description:     fmt.Sprintf("High request rate detected from IP %s: %d requests/minute", ipAddr, requestCount),
			IPAddress:       ipAddr,
			DetectionMethod: models.DetectionMethodRule,
			Confidence:      0.90,
			RiskScore:       9.0,
			Indicators: []models.ThreatIndicator{
				{
					Type:        "high_request_rate",
					Value:       fmt.Sprintf("%d", requestCount),
					Description: "Abnormally high request rate",
					Severity:    models.ThreatSeverityHigh,
					Confidence:  0.90,
				},
			},
			RequestData: traffic,
			FirstSeen:   timestamp,
			LastSeen:    timestamp,
			Count:       requestCount,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		
		threats = append(threats, threat)
	}
	
	return threats, nil
}

func (s *ThreatDetectionService) detectBruteForce(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat
	
	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}
	
	responseData, ok := traffic["response"].(map[string]interface{})
	if !ok {
		return threats, nil
	}
	
	// Check if this is an authentication endpoint
	path, ok := requestData["path"].(string)
	if !ok {
		return threats, nil
	}
	
	isAuthEndpoint := strings.Contains(strings.ToLower(path), "login") ||
		strings.Contains(strings.ToLower(path), "auth") ||
		strings.Contains(strings.ToLower(path), "signin")
	
	if !isAuthEndpoint {
		return threats, nil
	}
	
	// Check for failed authentication (401, 403 status codes)
	statusCode, ok := responseData["status_code"].(float64)
	if !ok || (statusCode != 401 && statusCode != 403) {
		return threats, nil
	}
	
	ipAddr, ok := traffic["ip_address"].(string)
	if !ok {
		return threats, nil
	}
	
	// Count failed authentication attempts from this IP in the last 5 minutes
	failedAttempts, err := s.threatRepo.GetFailedAuthAttempts(ctx, ipAddr, 5*time.Minute)
	if err != nil {
		return threats, fmt.Errorf("failed to get failed auth attempts: %w", err)
	}
	
	// Brute force threshold: more than 10 failed attempts in 5 minutes
	if failedAttempts > 10 {
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            models.ThreatTypeBruteForce,
			Severity:        models.ThreatSeverityHigh,
			Status:          models.ThreatStatusNew,
			Title:           "Brute Force Attack Detected",
			Description:     fmt.Sprintf("Multiple failed authentication attempts from IP %s: %d attempts in 5 minutes", ipAddr, failedAttempts),
			IPAddress:       ipAddr,
			DetectionMethod: models.DetectionMethodRule,
			Confidence:      0.85,
			RiskScore:       8.5,
			Indicators: []models.ThreatIndicator{
				{
					Type:        "failed_auth_attempts",
					Value:       fmt.Sprintf("%d", failedAttempts),
					Description: "Multiple failed authentication attempts",
					Severity:    models.ThreatSeverityHigh,
					Confidence:  0.85,
				},
			},
			RequestData:  requestData,
			ResponseData: responseData,
			FirstSeen:    time.Now(),
			LastSeen:     time.Now(),
			Count:        failedAttempts,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		if userAgent, ok := requestData["user_agent"].(string); ok {
			threat.UserAgent = userAgent
		}
		
		threats = append(threats, threat)
	}
	
	return threats, nil
}

func (s *ThreatDetectionService) detectDataExfiltration(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat
	
	responseData, ok := traffic["response"].(map[string]interface{})
	if !ok {
		return threats, nil
	}
	
	// Check response size
	responseSize, ok := responseData["size"].(float64)
	if !ok {
		return threats, nil
	}
	
	// Large response threshold: more than 10MB
	if responseSize > 10*1024*1024 {
		requestData, _ := traffic["request"].(map[string]interface{})
		
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            models.ThreatTypeDataExfiltration,
			Severity:        models.ThreatSeverityMedium,
			Status:          models.ThreatStatusNew,
			Title:           "Potential Data Exfiltration",
			Description:     fmt.Sprintf("Large response size detected: %.2f MB", responseSize/(1024*1024)),
			DetectionMethod: models.DetectionMethodRule,
			Confidence:      0.70,
			RiskScore:       7.0,
			Indicators: []models.ThreatIndicator{
				{
					Type:        "large_response",
					Value:       fmt.Sprintf("%.2f MB", responseSize/(1024*1024)),
					Description: "Unusually large response size",
					Severity:    models.ThreatSeverityMedium,
					Confidence:  0.70,
				},
			},
			RequestData:  requestData,
			ResponseData: responseData,
			FirstSeen:    time.Now(),
			LastSeen:     time.Now(),
			Count:        1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		
		if ipAddr, ok := traffic["ip_address"].(string); ok {
			threat.IPAddress = ipAddr
		}
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		if userAgent, ok := requestData["user_agent"].(string); ok {
			threat.UserAgent = userAgent
		}
		
		threats = append(threats, threat)
	}
	
	return threats, nil
}

func (s *ThreatDetectionService) getSeverityWeight(severity string) int {
	switch severity {
	case models.ThreatSeverityCritical:
		return 5
	case models.ThreatSeverityHigh:
		return 4
	case models.ThreatSeverityMedium:
		return 3
	case models.ThreatSeverityLow:
		return 2
	case models.ThreatSeverityInfo:
		return 1
	default:
		return 0
	}
}

func (s *ThreatDetectionService) generateRecommendations(threats []models.Threat) []string {
	recommendations := []string{}
	threatTypes := make(map[string]bool)
	
	for _, threat := range threats {
		threatTypes[threat.Type] = true
	}
	
	if threatTypes[models.ThreatTypeInjection] {
		recommendations = append(recommendations, "Implement input validation and parameterized queries")
		recommendations = append(recommendations, "Enable Web Application Firewall (WAF) rules for SQL injection")
	}
	
	if threatTypes[models.ThreatTypeXSS] {
		recommendations = append(recommendations, "Implement output encoding and Content Security Policy (CSP)")
		recommendations = append(recommendations, "Sanitize user inputs and validate data types")
	}
	
	if threatTypes[models.ThreatTypeDDoS] {
		recommendations = append(recommendations, "Implement rate limiting and IP blocking")
		recommendations = append(recommendations, "Consider using DDoS protection services")
	}
	
	if threatTypes[models.ThreatTypeBruteForce] {
		recommendations = append(recommendations, "Implement account lockout policies")
		recommendations = append(recommendations, "Enable multi-factor authentication (MFA)")
		recommendations = append(recommendations, "Implement CAPTCHA for repeated failed attempts")
	}
	
	if threatTypes[models.ThreatTypeDataExfiltration] {
		recommendations = append(recommendations, "Review data access permissions and implement data loss prevention (DLP)")
		recommendations = append(recommendations, "Monitor and alert on large data transfers")
	}
	
	return recommendations
}

func (s *ThreatDetectionService) publishThreatEvent(ctx context.Context, threats []models.Threat) error {
	for _, threat := range threats {
		eventData := map[string]interface{}{
			"event_type":    "threat_detected",
			"threat_id":     threat.ID,
			"threat_type":   threat.Type,
			"severity":      threat.Severity,
			"risk_score":    threat.RiskScore,
			"confidence":    threat.Confidence,
			"ip_address":    threat.IPAddress,
			"api_id":        threat.APIID,
			"endpoint_id":   threat.EndpointID,
			"timestamp":     threat.CreatedAt,
			"indicators":    threat.Indicators,
			"description":   threat.Description,
		}
		
		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal threat event: %w", err)
		}
		
		message := kafka.Message{
			Topic: "threat_events",
			Key:   []byte(threat.ID),
			Value: eventJSON,
		}
		
		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce threat event: %w", err)
		}
	}
	
	return nil
}

func (s *ThreatDetectionService) AnalyzeThreat(ctx context.Context, request *models.ThreatAnalysisRequest) (*models.ThreatAnalysisResult, error) {
	trafficJSON, err := json.Marshal(request.TrafficData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal traffic data: %w", err)
	}
	
	return s.AnalyzeTraffic(ctx, trafficJSON)
}

func (s *ThreatDetectionService) GetThreats(ctx context.Context, filter *models.ThreatFilter) ([]models.Threat, error) {
	return s.threatRepo.GetThreats(ctx, filter)
}

func (s *ThreatDetectionService) GetThreat(ctx context.Context, threatID string) (*models.Threat, error) {
	return s.threatRepo.GetThreat(ctx, threatID)
}

func (s *ThreatDetectionService) UpdateThreatStatus(ctx context.Context, threatID string, update *models.ThreatUpdateRequest) error {
	// Get existing threat
	threat, err := s.threatRepo.GetThreat(ctx, threatID)
	if err != nil {
		return err
	}
	
	// Update fields from request
	if update.Status != "" {
		threat.Status = update.Status
	}
	// Note: Additional update logic can be added here for other fields
	
	return s.threatRepo.UpdateThreat(ctx, threatID, threat)
}

func (s *ThreatDetectionService) DeleteThreat(ctx context.Context, threatID string) error {
	return s.threatRepo.DeleteThreat(ctx, threatID)
}

func (s *ThreatDetectionService) GetThreatStatistics(ctx context.Context, timeRange time.Duration) (*models.ThreatStatistics, error) {
	return s.threatRepo.GetThreatStatistics(ctx, timeRange)
}

func (s *ThreatDetectionService) CreateThreat(ctx context.Context, threat *models.Threat) error {
	return s.threatRepo.CreateThreat(ctx, threat)
}

func (s *ThreatDetectionService) ProcessSecurityEvent(ctx context.Context, event map[string]interface{}) error {
	// Convert event to traffic format and analyze
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal security event: %w", err)
	}
	
	result, err := s.AnalyzeTraffic(ctx, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to analyze security event: %w", err)
	}
	
	if result.ThreatDetected {
		s.logger.Info("Threat detected from security event", 
			"threat_type", result.ThreatType, 
			"severity", result.Severity,
			"risk_score", result.RiskScore)
	}
	
	return nil
}
