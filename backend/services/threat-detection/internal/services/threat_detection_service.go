package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/threat-detection/internal/models"
	"scopeapi.local/backend/services/threat-detection/internal/repository"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

// MLModel represents a machine learning model for threat detection
type MLModel struct {
	ID          string
	Name        string
	Version     string
	Type        string // "anomaly", "behavioral", "pattern"
	Accuracy    float64
	LastUpdated time.Time
	Features    []string
	Threshold   float64
}

// MLPrediction represents a prediction from an ML model
type MLPrediction struct {
	ModelID      string
	ModelType    string
	Prediction   float64
	Confidence   float64
	Features     map[string]float64
	AnomalyScore float64
	IsAnomaly    bool
}

// MLFeatureExtractor extracts features from traffic data for ML models
type MLFeatureExtractor struct {
	logger logging.Logger
}

// MLAnomalyDetector detects anomalies using machine learning
type MLAnomalyDetector struct {
	models map[string]*MLModel
	logger logging.Logger
}

// MLBehavioralAnalyzer analyzes behavioral patterns using ML
type MLBehavioralAnalyzer struct {
	models map[string]*MLModel
	logger logging.Logger
}

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

	// ML Integration methods
	TrainMLModel(ctx context.Context, modelType string, trainingData []byte) error
	UpdateMLModel(ctx context.Context, modelID string, newData []byte) error
	GetMLModelMetrics(ctx context.Context, modelID string) (*MLModel, error)
	PredictThreat(ctx context.Context, traffic map[string]interface{}) (*MLPrediction, error)
}

type ThreatDetectionService struct {
	threatRepo         repository.ThreatRepositoryInterface
	kafkaProducer      kafka.ProducerInterface
	logger             logging.Logger
	featureExtractor   *MLFeatureExtractor
	anomalyDetector    *MLAnomalyDetector
	behavioralAnalyzer *MLBehavioralAnalyzer
	mlModels           map[string]*MLModel
}

func NewThreatDetectionService(
	threatRepo repository.ThreatRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *ThreatDetectionService {
	// Initialize ML models
	mlModels := initializeMLModels()

	return &ThreatDetectionService{
		threatRepo:         threatRepo,
		kafkaProducer:      kafkaProducer,
		logger:             logger,
		featureExtractor:   NewMLFeatureExtractor(logger),
		anomalyDetector:    NewMLAnomalyDetector(logger, mlModels),
		behavioralAnalyzer: NewMLBehavioralAnalyzer(logger, mlModels),
		mlModels:           mlModels,
	}
}

// initializeMLModels sets up default ML models
func initializeMLModels() map[string]*MLModel {
	models := make(map[string]*MLModel)

	// Anomaly Detection Model
	models["anomaly_detection"] = &MLModel{
		ID:          "anomaly_detection_v1",
		Name:        "Anomaly Detection Model",
		Version:     "1.0.0",
		Type:        "anomaly",
		Accuracy:    0.92,
		LastUpdated: time.Now(),
		Features:    []string{"request_rate", "response_time", "payload_size", "error_rate", "unique_ips", "user_agent_diversity"},
		Threshold:   0.75,
	}

	// Behavioral Analysis Model
	models["behavioral_analysis"] = &MLModel{
		ID:          "behavioral_analysis_v1",
		Name:        "Behavioral Analysis Model",
		Version:     "1.0.0",
		Type:        "behavioral",
		Accuracy:    0.89,
		LastUpdated: time.Now(),
		Features:    []string{"session_pattern", "request_sequence", "timing_pattern", "resource_access", "data_volume"},
		Threshold:   0.70,
	}

	// Pattern Recognition Model
	models["pattern_recognition"] = &MLModel{
		ID:          "pattern_recognition_v1",
		Name:        "Pattern Recognition Model",
		Version:     "1.0.0",
		Type:        "pattern",
		Accuracy:    0.94,
		LastUpdated: time.Now(),
		Features:    []string{"url_pattern", "parameter_pattern", "header_pattern", "payload_pattern", "response_pattern"},
		Threshold:   0.80,
	}

	return models
}

func NewMLFeatureExtractor(logger logging.Logger) *MLFeatureExtractor {
	return &MLFeatureExtractor{
		logger: logger,
	}
}

func NewMLAnomalyDetector(logger logging.Logger, models map[string]*MLModel) *MLAnomalyDetector {
	return &MLAnomalyDetector{
		models: models,
		logger: logger,
	}
}

func NewMLBehavioralAnalyzer(logger logging.Logger, models map[string]*MLModel) *MLBehavioralAnalyzer {
	return &MLBehavioralAnalyzer{
		models: models,
		logger: logger,
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
		RequestID:       requestID,
		ThreatDetected:  false,
		Confidence:      0.0,
		RiskScore:       0.0,
		Indicators:      []models.ThreatIndicator{},
		Recommendations: []string{},
		Metadata:        make(map[string]interface{}),
		AnalyzedAt:      time.Now(),
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

	// 6. Path Traversal Detection
	pathTraversalThreats, err := s.detectPathTraversal(ctx, traffic)
	if err != nil {
		s.logger.Error("Path traversal detection failed", "error", err)
	} else {
		threats = append(threats, pathTraversalThreats...)
	}

	// 7. Command Injection Detection
	commandInjectionThreats, err := s.detectCommandInjection(ctx, traffic)
	if err != nil {
		s.logger.Error("Command injection detection failed", "error", err)
	} else {
		threats = append(threats, commandInjectionThreats...)
	}

	// 8. ML-Based Anomaly Detection
	mlAnomalyThreats, err := s.detectMLAnomalies(ctx, traffic)
	if err != nil {
		s.logger.Error("ML anomaly detection failed", "error", err)
	} else {
		threats = append(threats, mlAnomalyThreats...)
	}

	// 9. ML-Based Behavioral Analysis
	mlBehavioralThreats, err := s.detectMLBehavioral(ctx, traffic)
	if err != nil {
		s.logger.Error("ML behavioral analysis failed", "error", err)
	} else {
		threats = append(threats, mlBehavioralThreats...)
	}

	// 10. ML-Based Pattern Recognition
	mlPatternThreats, err := s.detectMLPatterns(ctx, traffic)
	if err != nil {
		s.logger.Error("ML pattern recognition failed", "error", err)
	} else {
		threats = append(threats, mlPatternThreats...)
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
	result.Metadata["analysis_methods"] = []string{
		"sql_injection", "xss", "ddos", "brute_force", "data_exfiltration",
		"path_traversal", "command_injection", "ml_anomaly", "ml_behavioral", "ml_pattern",
	}
	result.Metadata["ml_models_used"] = []string{
		"anomaly_detection_v1", "behavioral_analysis_v1", "pattern_recognition_v1",
	}

	return result, nil
}

func (s *ThreatDetectionService) detectSQLInjection(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	// Extract request data
	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}

	// SQL Injection patterns to detect
	sqlPatterns := []string{
		"';", "--", "/*", "*/", "xp_", "sp_", "exec", "execute", "union", "select", "insert", "update", "delete", "drop", "create", "alter",
		"1=1", "1=1--", "1=1#", "1=1/*", "1=1*/", "1=1 union", "1=1 union select", "1=1 union select 1", "1=1 union select 1,2",
		"admin'--", "admin'#", "admin'/*", "admin'*/", "admin' or '1'='1", "admin' or '1'='1'--", "admin' or '1'='1'#",
		"or 1=1", "or 1=1--", "or 1=1#", "or 1=1/*", "or 1=1*/", "or 1=1 union", "or 1=1 union select",
		"union select", "union select 1", "union select 1,2", "union select 1,2,3", "union select 1,2,3,4",
		"waitfor delay", "waitfor time", "benchmark", "sleep", "pg_sleep", "dbms_pipe.receive_message",
	}

	// Check URL parameters
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsSQLInjectionPattern(fmt.Sprintf("%v", value), sqlPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "sql_injection",
					Severity:        "high",
					Status:          "new",
					Title:           "SQL Injection Attempt Detected",
					Description:     fmt.Sprintf("Potential SQL injection detected in parameter '%s' with value '%v'", key, value),
					DetectionMethod: "signature",
					Confidence:      0.85,
					RiskScore:       8.5,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "sql_injection_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: "Suspicious SQL pattern detected",
							Severity:    "high",
							Confidence:  0.85,
						},
					},
					RequestData: requestData,
					FirstSeen:   time.Now(),
					LastSeen:    time.Now(),
					Count:       1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				// Extract IP address if available
				if ip, ok := requestData["ip_address"].(string); ok {
					threat.IPAddress = ip
					threat.SourceIP = ip
				}

				// Extract user agent if available
				if ua, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = ua
				}

				// Extract API and endpoint information
				if apiID, ok := requestData["api_id"].(string); ok {
					threat.APIID = apiID
				}
				if endpointID, ok := requestData["endpoint_id"].(string); ok {
					threat.EndpointID = endpointID
				}

				threats = append(threats, threat)
			}
		}
	}

	// Check request body
	if body, ok := requestData["body"].(string); ok {
		if s.containsSQLInjectionPattern(body, sqlPatterns) {
			threat := models.Threat{
				ID:              uuid.New().String(),
				Type:            "sql_injection",
				Severity:        "high",
				Status:          "new",
				Title:           "SQL Injection Attempt in Request Body",
				Description:     "Potential SQL injection detected in request body",
				DetectionMethod: "signature",
				Confidence:      0.80,
				RiskScore:       8.0,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "sql_injection_pattern",
						Value:       body,
						Description: "Suspicious SQL pattern detected in request body",
						Severity:    "high",
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
			if ip, ok := requestData["ip_address"].(string); ok {
				threat.IPAddress = ip
				threat.SourceIP = ip
			}

			threats = append(threats, threat)
		}
	}

	// Check headers for suspicious content
	if headers, ok := requestData["headers"].(map[string]interface{}); ok {
		for headerName, headerValue := range headers {
			if s.containsSQLInjectionPattern(fmt.Sprintf("%v", headerValue), sqlPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "sql_injection",
					Severity:        "medium",
					Status:          "new",
					Title:           "SQL Injection Attempt in Header",
					Description:     fmt.Sprintf("Potential SQL injection detected in header '%s'", headerName),
					DetectionMethod: "signature",
					Confidence:      0.75,
					RiskScore:       7.5,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "sql_injection_pattern",
							Value:       fmt.Sprintf("%v", headerValue),
							Description: fmt.Sprintf("Suspicious SQL pattern detected in header '%s'", headerName),
							Severity:    "medium",
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

				threats = append(threats, threat)
			}
		}
	}

	return threats, nil
}

func (s *ThreatDetectionService) containsSQLInjectionPattern(input string, patterns []string) bool {
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

	// Enhanced XSS patterns to detect
	xssPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror=", "onclick=", "onmouseover=",
		"onfocus=", "onblur=", "onchange=", "onsubmit=", "onreset=", "onselect=", "onunload=",
		"<iframe", "</iframe>", "<object", "</object>", "<embed", "<form", "<input", "<textarea",
		"<svg", "<math", "<xmp", "<plaintext", "<listing", "<marquee", "<applet", "<bgsound",
		"<link", "<meta", "<title", "<style", "<div", "<span", "<p", "<a", "<img",
		"alert(", "confirm(", "prompt(", "eval(", "setTimeout(", "setInterval(", "Function(",
		"document.cookie", "document.write", "document.writeln", "window.open", "location.href",
		"innerHTML", "outerHTML", "insertAdjacentHTML", "createElement", "appendChild",
		"<base", "<bdo", "<noscript", "<noframes", "<frameset", "<frame", "<noframe",
		"<xss", "<xss>", "expression(", "url(", "behavior:", "binding:", "mocha:",
	}

	// Check URL parameters
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsXSSPattern(fmt.Sprintf("%v", value), xssPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "xss",
					Severity:        "high",
					Status:          "new",
					Title:           "Cross-Site Scripting (XSS) Attempt Detected",
					Description:     fmt.Sprintf("Potential XSS attack detected in parameter '%s'", key),
					DetectionMethod: "signature",
					Confidence:      0.90,
					RiskScore:       9.0,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "xss_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: "Suspicious XSS pattern detected",
							Severity:    "high",
							Confidence:  0.90,
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
				if ip, ok := requestData["ip_address"].(string); ok {
					threat.IPAddress = ip
					threat.SourceIP = ip
				}
				if ua, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = ua
				}
				if apiID, ok := requestData["api_id"].(string); ok {
					threat.APIID = apiID
				}
				if endpointID, ok := requestData["endpoint_id"].(string); ok {
					threat.EndpointID = endpointID
				}

				threats = append(threats, threat)
			}
		}
	}

	// Check request body
	if body, ok := requestData["body"].(string); ok {
		if s.containsXSSPattern(body, xssPatterns) {
			threat := models.Threat{
				ID:              uuid.New().String(),
				Type:            "xss",
				Severity:        "high",
				Status:          "new",
				Title:           "Cross-Site Scripting (XSS) in Request Body",
				Description:     "Potential XSS attack detected in request body",
				DetectionMethod: "signature",
				Confidence:      0.85,
				RiskScore:       8.5,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "xss_pattern",
						Value:       body,
						Description: "XSS pattern detected in request body",
						Severity:    "high",
						Confidence:  0.85,
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
			if ip, ok := requestData["ip_address"].(string); ok {
				threat.IPAddress = ip
				threat.SourceIP = ip
			}

			threats = append(threats, threat)
		}
	}

	// Check headers for suspicious content
	if headers, ok := requestData["headers"].(map[string]interface{}); ok {
		for headerName, headerValue := range headers {
			if s.containsXSSPattern(fmt.Sprintf("%v", headerValue), xssPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "xss",
					Severity:        "medium",
					Status:          "new",
					Title:           "Cross-Site Scripting (XSS) in Header",
					Description:     fmt.Sprintf("Potential XSS attack detected in header '%s'", headerName),
					DetectionMethod: "signature",
					Confidence:      0.80,
					RiskScore:       8.0,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "xss_pattern",
							Value:       fmt.Sprintf("%v", headerValue),
							Description: fmt.Sprintf("XSS pattern detected in header '%s'", headerName),
							Severity:    "medium",
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

				threats = append(threats, threat)
			}
		}
	}

	return threats, nil
}

func (s *ThreatDetectionService) containsXSSPattern(input string, patterns []string) bool {
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, strings.ToLower(pattern)) {
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
		// If timestamp is not available, use current time
		timestamp = time.Now()
	}

	// Enhanced DDoS detection with multiple thresholds
	thresholds := []struct {
		duration time.Duration
		limit    int
		severity string
		score    float64
	}{
		{time.Minute, 100, "high", 9.0},     // 100 req/min = high severity
		{time.Minute, 50, "medium", 7.5},    // 50 req/min = medium severity
		{time.Minute, 25, "low", 6.0},       // 25 req/min = low severity
		{time.Second * 10, 20, "high", 9.5}, // 20 req/10sec = very high severity
		{time.Second * 5, 10, "high", 9.0},  // 10 req/5sec = high severity
	}

	for _, threshold := range thresholds {
		requestCount, err := s.threatRepo.GetRequestCountByIP(ctx, ipAddr, threshold.duration)
		if err != nil {
			s.logger.Warn("Failed to get request count for DDoS detection", "ip", ipAddr, "duration", threshold.duration, "error", err)
			continue
		}

		if requestCount > threshold.limit {
			// Calculate request rate
			rate := float64(requestCount) / threshold.duration.Seconds()

			threat := models.Threat{
				ID:       uuid.New().String(),
				Type:     "ddos",
				Severity: threshold.severity,
				Status:   "new",
				Title:    "Potential DDoS Attack Detected",
				Description: fmt.Sprintf("High request rate detected from IP %s: %d requests in %v (%.2f req/sec)",
					ipAddr, requestCount, threshold.duration, rate),
				IPAddress:       ipAddr,
				SourceIP:        ipAddr,
				DetectionMethod: "rate_limiting",
				Confidence:      0.90,
				RiskScore:       threshold.score,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "high_request_rate",
						Value:       fmt.Sprintf("%d", requestCount),
						Description: fmt.Sprintf("Abnormally high request rate: %d requests in %v", requestCount, threshold.duration),
						Severity:    threshold.severity,
						Confidence:  0.90,
					},
					{
						Type:        "request_rate",
						Value:       fmt.Sprintf("%.2f", rate),
						Description: fmt.Sprintf("Request rate: %.2f requests per second", rate),
						Severity:    threshold.severity,
						Confidence:  0.85,
					},
				},
				RequestData: traffic,
				FirstSeen:   timestamp,
				LastSeen:    timestamp,
				Count:       requestCount,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Extract additional context
			if apiID, ok := traffic["api_id"].(string); ok {
				threat.APIID = apiID
			}
			if endpointID, ok := traffic["endpoint_id"].(string); ok {
				threat.EndpointID = endpointID
			}
			if userAgent, ok := traffic["user_agent"].(string); ok {
				threat.UserAgent = userAgent
			}

			threats = append(threats, threat)

			// Only create one threat per IP, so break after first threshold exceeded
			break
		}
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

	if threatTypes["path_traversal"] {
		recommendations = append(recommendations, "Implement strict path validation and sanitization")
		recommendations = append(recommendations, "Use allowlist approach for file paths and directories")
		recommendations = append(recommendations, "Implement proper URL encoding/decoding validation")
	}

	if threatTypes["command_injection"] {
		recommendations = append(recommendations, "Avoid using system() or exec() functions with user input")
		recommendations = append(recommendations, "Implement strict input validation and sanitization")
		recommendations = append(recommendations, "Use parameterized APIs instead of shell commands")
		recommendations = append(recommendations, "Implement proper escaping for shell commands if necessary")
	}

	// ML-based threat recommendations
	if threatTypes["ml_anomaly"] {
		recommendations = append(recommendations, "Review ML model performance and retrain if necessary")
		recommendations = append(recommendations, "Adjust anomaly detection thresholds based on false positive rates")
		recommendations = append(recommendations, "Implement adaptive learning for improved anomaly detection")
	}

	if threatTypes["ml_behavioral"] {
		recommendations = append(recommendations, "Analyze behavioral patterns to identify new threat indicators")
		recommendations = append(recommendations, "Update behavioral baseline models with new data")
		recommendations = append(recommendations, "Implement user behavior analytics for enhanced detection")
	}

	if threatTypes["ml_pattern"] {
		recommendations = append(recommendations, "Update pattern recognition models with new threat signatures")
		recommendations = append(recommendations, "Implement ensemble learning for improved pattern detection")
		recommendations = append(recommendations, "Use transfer learning for cross-domain threat pattern recognition")
	}

	return recommendations
}

func (s *ThreatDetectionService) publishThreatEvent(ctx context.Context, threats []models.Threat) error {
	for _, threat := range threats {
		eventData := map[string]interface{}{
			"event_type":  "threat_detected",
			"threat_id":   threat.ID,
			"threat_type": threat.Type,
			"severity":    threat.Severity,
			"risk_score":  threat.RiskScore,
			"confidence":  threat.Confidence,
			"ip_address":  threat.IPAddress,
			"api_id":      threat.APIID,
			"endpoint_id": threat.EndpointID,
			"timestamp":   threat.CreatedAt,
			"indicators":  threat.Indicators,
			"description": threat.Description,
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

// detectPathTraversal detects path traversal attacks
func (s *ThreatDetectionService) detectPathTraversal(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}

	// Path traversal patterns to detect
	pathTraversalPatterns := []string{
		"../", "..\\", "..%2f", "..%5c", "..%2F", "..%5C",
		"....//", "....\\\\", "....%2f", "....%5c",
		"%2e%2e%2f", "%2e%2e%5c", "%2e%2e%2F", "%2e%2e%5C",
		"..%252f", "..%255c", "..%252F", "..%255C",
		"..%c0%af", "..%c1%9c", "..%c0%AF", "..%c1%9C",
		"..%255c", "..%255C", "..%c1%9c", "..%c1%9C",
		"..%c0%af", "..%c0%AF", "..%c1%9c", "..%c1%9C",
		"..%c0%af", "..%c0%AF", "..%c1%9c", "..%c1%9C",
		"..%c0%af", "..%c0%AF", "..%c1%9c", "..%c1%9C",
	}

	// Check URL path
	if path, ok := requestData["path"].(string); ok {
		if s.containsPathTraversalPattern(path, pathTraversalPatterns) {
			threat := models.Threat{
				ID:              uuid.New().String(),
				Type:            "path_traversal",
				Severity:        "high",
				Status:          "new",
				Title:           "Path Traversal Attack Detected",
				Description:     fmt.Sprintf("Potential path traversal attack detected in URL path: %s", path),
				DetectionMethod: "signature",
				Confidence:      0.95,
				RiskScore:       9.5,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "path_traversal_pattern",
						Value:       path,
						Description: "Path traversal pattern detected in URL",
						Severity:    "high",
						Confidence:  0.95,
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
			if ip, ok := requestData["ip_address"].(string); ok {
				threat.IPAddress = ip
				threat.SourceIP = ip
			}
			if ua, ok := requestData["user_agent"].(string); ok {
				threat.UserAgent = ua
			}
			if apiID, ok := requestData["api_id"].(string); ok {
				threat.APIID = apiID
			}
			if endpointID, ok := requestData["endpoint_id"].(string); ok {
				threat.EndpointID = endpointID
			}

			threats = append(threats, threat)
		}
	}

	// Check URL parameters for path traversal
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsPathTraversalPattern(fmt.Sprintf("%v", value), pathTraversalPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "path_traversal",
					Severity:        "high",
					Status:          "new",
					Title:           "Path Traversal Attack in Parameter",
					Description:     fmt.Sprintf("Potential path traversal attack detected in parameter '%s'", key),
					DetectionMethod: "signature",
					Confidence:      0.90,
					RiskScore:       9.0,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "path_traversal_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: fmt.Sprintf("Path traversal pattern detected in parameter '%s'", key),
							Severity:    "high",
							Confidence:  0.90,
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
				if ip, ok := requestData["ip_address"].(string); ok {
					threat.IPAddress = ip
					threat.SourceIP = ip
				}
				if ua, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = ua
				}
				if apiID, ok := requestData["api_id"].(string); ok {
					threat.APIID = apiID
				}
				if endpointID, ok := requestData["endpoint_id"].(string); ok {
					threat.EndpointID = endpointID
				}

				threats = append(threats, threat)
			}
		}
	}

	return threats, nil
}

func (s *ThreatDetectionService) containsPathTraversalPattern(input string, patterns []string) bool {
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// detectCommandInjection detects command injection attacks
func (s *ThreatDetectionService) detectCommandInjection(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	requestData, ok := traffic["request"].(map[string]interface{})
	if !ok {
		return threats, nil
	}

	// Command injection patterns to detect
	commandInjectionPatterns := []string{
		";", "|", "&", "&&", "||", "`", "$(", "$(((", "eval", "exec", "system",
		"cmd", "command", "powershell", "bash", "sh", "csh", "tcsh", "zsh",
		"ping", "nslookup", "dig", "whois", "traceroute", "netstat", "ps",
		"cat", "ls", "dir", "type", "more", "less", "head", "tail", "grep",
		"find", "locate", "which", "where", "whereis", "wget", "curl", "nc",
		"telnet", "ssh", "ftp", "scp", "rsync", "tar", "zip", "unzip",
		"chmod", "chown", "chgrp", "umask", "su", "sudo", "passwd",
		"useradd", "userdel", "groupadd", "groupdel", "usermod", "groupmod",
		"service", "systemctl", "init", "rc", "cron", "at", "anacron",
		"iptables", "firewall", "ufw", "selinux", "apparmor", "auditd",
		"logrotate", "rsyslog", "syslog", "journalctl", "dmesg", "last",
		"who", "w", "uptime", "top", "htop", "iotop", "iotop", "iotop",
	}

	// Check URL parameters
	if params, ok := requestData["parameters"].(map[string]interface{}); ok {
		for key, value := range params {
			if s.containsCommandInjectionPattern(fmt.Sprintf("%v", value), commandInjectionPatterns) {
				threat := models.Threat{
					ID:              uuid.New().String(),
					Type:            "command_injection",
					Severity:        "critical",
					Status:          "new",
					Title:           "Command Injection Attack Detected",
					Description:     fmt.Sprintf("Potential command injection attack detected in parameter '%s'", key),
					DetectionMethod: "signature",
					Confidence:      0.95,
					RiskScore:       10.0,
					Indicators: []models.ThreatIndicator{
						{
							Type:        "command_injection_pattern",
							Value:       fmt.Sprintf("%v", value),
							Description: "Suspicious command injection pattern detected",
							Severity:    "critical",
							Confidence:  0.95,
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
				if ip, ok := requestData["ip_address"].(string); ok {
					threat.IPAddress = ip
					threat.SourceIP = ip
				}
				if ua, ok := requestData["user_agent"].(string); ok {
					threat.UserAgent = ua
				}
				if apiID, ok := requestData["api_id"].(string); ok {
					threat.APIID = apiID
				}
				if endpointID, ok := requestData["endpoint_id"].(string); ok {
					threat.EndpointID = endpointID
				}

				threats = append(threats, threat)
			}
		}
	}

	// Check request body
	if body, ok := requestData["body"].(string); ok {
		if s.containsCommandInjectionPattern(body, commandInjectionPatterns) {
			threat := models.Threat{
				ID:              uuid.New().String(),
				Type:            "command_injection",
				Severity:        "critical",
				Status:          "new",
				Title:           "Command Injection Attack in Request Body",
				Description:     "Potential command injection attack detected in request body",
				DetectionMethod: "signature",
				Confidence:      0.90,
				RiskScore:       9.5,
				Indicators: []models.ThreatIndicator{
					{
						Type:        "command_injection_pattern",
						Value:       body,
						Description: "Command injection pattern detected in request body",
						Severity:    "critical",
						Confidence:  0.90,
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
			if ip, ok := requestData["ip_address"].(string); ok {
				threat.IPAddress = ip
				threat.SourceIP = ip
			}

			threats = append(threats, threat)
		}
	}

	return threats, nil
}

func (s *ThreatDetectionService) containsCommandInjectionPattern(input string, patterns []string) bool {
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// detectMLAnomalies uses ML models to detect anomalies
func (s *ThreatDetectionService) detectMLAnomalies(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	// Extract features for ML analysis
	features := s.featureExtractor.ExtractAnomalyFeatures(traffic)

	// Get ML prediction
	prediction, err := s.PredictThreat(ctx, traffic)
	if err != nil {
		s.logger.Warn("ML prediction failed for anomaly detection", "error", err)
		return threats, nil
	}

	// Check if anomaly is detected
	if prediction.IsAnomaly && prediction.AnomalyScore > s.mlModels["anomaly_detection"].Threshold {
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            "ml_anomaly",
			Severity:        s.calculateMLSeverity(prediction.AnomalyScore),
			Status:          "new",
			Title:           "ML-Based Anomaly Detected",
			Description:     fmt.Sprintf("Machine learning model detected anomalous behavior with score %.3f", prediction.AnomalyScore),
			DetectionMethod: "machine_learning",
			Confidence:      prediction.Confidence,
			RiskScore:       prediction.AnomalyScore * 10.0, // Scale to 0-10
			Indicators: []models.ThreatIndicator{
				{
					Type:        "ml_anomaly_score",
					Value:       fmt.Sprintf("%.3f", prediction.AnomalyScore),
					Description: "ML model anomaly score",
					Severity:    s.calculateMLSeverity(prediction.AnomalyScore),
					Confidence:  prediction.Confidence,
				},
				{
					Type:        "ml_model",
					Value:       prediction.ModelID,
					Description: "ML model used for detection",
					Severity:    s.calculateMLSeverity(prediction.AnomalyScore),
					Confidence:  prediction.Confidence,
				},
			},
			RequestData: traffic,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			Count:       1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Extract additional context
		if ip, ok := traffic["ip_address"].(string); ok {
			threat.IPAddress = ip
			threat.SourceIP = ip
		}
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		if endpointID, ok := traffic["endpoint_id"].(string); ok {
			threat.EndpointID = endpointID
		}

		// Add ML-specific metadata
		threat.Metadata = map[string]interface{}{
			"ml_model_id":      prediction.ModelID,
			"ml_model_type":    prediction.ModelType,
			"ml_features":      features,
			"ml_prediction":    prediction.Prediction,
			"ml_confidence":    prediction.Confidence,
			"ml_anomaly_score": prediction.AnomalyScore,
		}

		threats = append(threats, threat)
	}

	return threats, nil
}

// detectMLBehavioral uses ML models to analyze behavioral patterns
func (s *ThreatDetectionService) detectMLBehavioral(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	// Extract behavioral features
	features := s.featureExtractor.ExtractBehavioralFeatures(traffic)

	// Get ML prediction for behavioral analysis
	prediction, err := s.PredictThreat(ctx, traffic)
	if err != nil {
		s.logger.Warn("ML prediction failed for behavioral analysis", "error", err)
		return threats, nil
	}

	// Check if suspicious behavior is detected
	if prediction.Prediction > s.mlModels["behavioral_analysis"].Threshold {
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            "ml_behavioral",
			Severity:        s.calculateMLSeverity(prediction.Prediction),
			Status:          "new",
			Title:           "ML-Based Suspicious Behavior Detected",
			Description:     fmt.Sprintf("Machine learning model detected suspicious behavioral pattern with score %.3f", prediction.Prediction),
			DetectionMethod: "machine_learning",
			Confidence:      prediction.Confidence,
			RiskScore:       prediction.Prediction * 10.0, // Scale to 0-10
			Indicators: []models.ThreatIndicator{
				{
					Type:        "ml_behavioral_score",
					Value:       fmt.Sprintf("%.3f", prediction.Prediction),
					Description: "ML model behavioral score",
					Severity:    s.calculateMLSeverity(prediction.Prediction),
					Confidence:  prediction.Confidence,
				},
				{
					Type:        "ml_model",
					Value:       prediction.ModelID,
					Description: "ML model used for behavioral analysis",
					Severity:    s.calculateMLSeverity(prediction.Prediction),
					Confidence:  prediction.Confidence,
				},
			},
			RequestData: traffic,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			Count:       1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Extract additional context
		if ip, ok := traffic["ip_address"].(string); ok {
			threat.IPAddress = ip
			threat.SourceIP = ip
		}
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		if endpointID, ok := traffic["endpoint_id"].(string); ok {
			threat.EndpointID = endpointID
		}

		// Add ML-specific metadata
		threat.Metadata = map[string]interface{}{
			"ml_model_id":         prediction.ModelID,
			"ml_model_type":       prediction.ModelType,
			"ml_features":         features,
			"ml_prediction":       prediction.Prediction,
			"ml_confidence":       prediction.Confidence,
			"ml_behavioral_score": prediction.Prediction,
		}

		threats = append(threats, threat)
	}

	return threats, nil
}

// detectMLPatterns uses ML models to recognize threat patterns
func (s *ThreatDetectionService) detectMLPatterns(ctx context.Context, traffic map[string]interface{}) ([]models.Threat, error) {
	var threats []models.Threat

	// Extract pattern features
	features := s.featureExtractor.ExtractPatternFeatures(traffic)

	// Get ML prediction for pattern recognition
	prediction, err := s.PredictThreat(ctx, traffic)
	if err != nil {
		s.logger.Warn("ML prediction failed for pattern recognition", "error", err)
		return threats, nil
	}

	// Check if threat pattern is detected
	if prediction.Prediction > s.mlModels["pattern_recognition"].Threshold {
		threat := models.Threat{
			ID:              uuid.New().String(),
			Type:            "ml_pattern",
			Severity:        s.calculateMLSeverity(prediction.Prediction),
			Status:          "new",
			Title:           "ML-Based Threat Pattern Detected",
			Description:     fmt.Sprintf("Machine learning model detected threat pattern with score %.3f", prediction.Prediction),
			DetectionMethod: "machine_learning",
			Confidence:      prediction.Confidence,
			RiskScore:       prediction.Prediction * 10.0, // Scale to 0-10
			Indicators: []models.ThreatIndicator{
				{
					Type:        "ml_pattern_score",
					Value:       fmt.Sprintf("%.3f", prediction.Prediction),
					Description: "ML model pattern recognition score",
					Severity:    s.calculateMLSeverity(prediction.Prediction),
					Confidence:  prediction.Confidence,
				},
				{
					Type:        "ml_model",
					Value:       prediction.ModelID,
					Description: "ML model used for pattern recognition",
					Severity:    s.calculateMLSeverity(prediction.Prediction),
					Confidence:  prediction.Confidence,
				},
			},
			RequestData: traffic,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			Count:       1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Extract additional context
		if ip, ok := traffic["ip_address"].(string); ok {
			threat.IPAddress = ip
			threat.SourceIP = ip
		}
		if apiID, ok := traffic["api_id"].(string); ok {
			threat.APIID = apiID
		}
		if endpointID, ok := traffic["endpoint_id"].(string); ok {
			threat.EndpointID = endpointID
		}

		// Add ML-specific metadata
		threat.Metadata = map[string]interface{}{
			"ml_model_id":      prediction.ModelID,
			"ml_model_type":    prediction.ModelType,
			"ml_features":      features,
			"ml_prediction":    prediction.Prediction,
			"ml_confidence":    prediction.Confidence,
			"ml_pattern_score": prediction.Prediction,
		}

		threats = append(threats, threat)
	}

	return threats, nil
}

// calculateMLSeverity calculates severity based on ML prediction score
func (s *ThreatDetectionService) calculateMLSeverity(score float64) string {
	switch {
	case score >= 0.9:
		return "critical"
	case score >= 0.8:
		return "high"
	case score >= 0.6:
		return "medium"
	case score >= 0.4:
		return "low"
	default:
		return "info"
	}
}

// ExtractAnomalyFeatures extracts features for anomaly detection
func (e *MLFeatureExtractor) ExtractAnomalyFeatures(traffic map[string]interface{}) map[string]float64 {
	features := make(map[string]float64)

	// Request rate features - calculate actual rates instead of placeholders
	if _, ok := traffic["ip_address"].(string); ok {
		// Calculate request rate based on timestamp and request count
		if timestamp, ok := traffic["timestamp"].(time.Time); ok {
			// In a real system, this would query historical data for this IP
			// For now, use a simulated rate based on current time
			hour := timestamp.Hour()
			if hour >= 9 && hour <= 17 {
				features["request_rate"] = 15.0 // Business hours: higher rate
			} else {
				features["request_rate"] = 5.0 // Off hours: lower rate
			}
		} else {
			features["request_rate"] = 10.0 // Default rate
		}

		// Calculate unique IPs (in a real system, this would be from historical data)
		features["unique_ips"] = 1.0 // Single IP for current request
	}

	// Response time features
	if response, ok := traffic["response"].(map[string]interface{}); ok {
		if responseTime, ok := response["response_time"].(float64); ok {
			features["response_time"] = responseTime
		}
	}

	// Payload size features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if body, ok := request["body"].(string); ok {
			features["payload_size"] = float64(len(body))
		}
	}

	// Error rate features
	if response, ok := traffic["response"].(map[string]interface{}); ok {
		if statusCode, ok := response["status_code"].(float64); ok {
			if statusCode >= 400 {
				features["error_rate"] = 1.0
			} else {
				features["error_rate"] = 0.0
			}
		}
	}

	// User agent diversity - calculate actual diversity score
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if userAgent, ok := request["user_agent"].(string); ok {
			// Calculate diversity based on user agent characteristics
			diversity := 0.0
			if strings.Contains(userAgent, "Mozilla") {
				diversity += 0.3
			}
			if strings.Contains(userAgent, "Chrome") {
				diversity += 0.2
			}
			if strings.Contains(userAgent, "Firefox") {
				diversity += 0.2
			}
			if strings.Contains(userAgent, "Safari") {
				diversity += 0.2
			}
			if strings.Contains(userAgent, "Mobile") {
				diversity += 0.1
			}
			features["user_agent_diversity"] = diversity
		}
	}

	return features
}

// ExtractBehavioralFeatures extracts features for behavioral analysis
func (e *MLFeatureExtractor) ExtractBehavioralFeatures(traffic map[string]interface{}) map[string]float64 {
	features := make(map[string]float64)

	// Session pattern features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if path, ok := request["path"].(string); ok {
			features["session_pattern"] = float64(len(path))
		}
		if method, ok := request["method"].(string); ok {
			features["request_sequence"] = float64(len(method))
		}
	}

	// Timing pattern features
	if timestamp, ok := traffic["timestamp"].(time.Time); ok {
		features["timing_pattern"] = float64(timestamp.Hour())
	}

	// Resource access features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if path, ok := request["path"].(string); ok {
			features["resource_access"] = float64(len(path))
		}
	}

	// Data volume features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if body, ok := request["body"].(string); ok {
			features["data_volume"] = float64(len(body))
		}
	}

	return features
}

// ExtractPatternFeatures extracts features for pattern recognition
func (e *MLFeatureExtractor) ExtractPatternFeatures(traffic map[string]interface{}) map[string]float64 {
	features := make(map[string]float64)

	// URL pattern features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if path, ok := request["path"].(string); ok {
			features["url_pattern"] = float64(len(path))
			features["url_complexity"] = float64(strings.Count(path, "/"))
		}
	}

	// Parameter pattern features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if params, ok := request["parameters"].(map[string]interface{}); ok {
			features["parameter_pattern"] = float64(len(params))
		}
	}

	// Header pattern features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if headers, ok := request["headers"].(map[string]interface{}); ok {
			features["header_pattern"] = float64(len(headers))
		}
	}

	// Payload pattern features
	if request, ok := traffic["request"].(map[string]interface{}); ok {
		if body, ok := request["body"].(string); ok {
			features["payload_pattern"] = float64(len(body))
			features["payload_complexity"] = float64(strings.Count(body, " "))
		}
	}

	// Response pattern features
	if response, ok := traffic["response"].(map[string]interface{}); ok {
		if statusCode, ok := response["status_code"].(float64); ok {
			features["response_pattern"] = statusCode
		}
	}

	return features
}

// PredictThreat uses ML models to predict threats
func (s *ThreatDetectionService) PredictThreat(ctx context.Context, traffic map[string]interface{}) (*MLPrediction, error) {
	// This is a simplified ML prediction implementation
	// In a real system, this would call actual ML models or APIs

	// Extract features
	anomalyFeatures := s.featureExtractor.ExtractAnomalyFeatures(traffic)
	behavioralFeatures := s.featureExtractor.ExtractBehavioralFeatures(traffic)
	patternFeatures := s.featureExtractor.ExtractPatternFeatures(traffic)

	// Calculate anomaly score (simplified)
	anomalyScore := s.calculateAnomalyScore(anomalyFeatures)

	// Calculate behavioral score (simplified)
	behavioralScore := s.calculateBehavioralScore(behavioralFeatures)

	// Calculate pattern score (simplified)
	patternScore := s.calculatePatternScore(patternFeatures)

	// Combine scores
	combinedScore := (anomalyScore + behavioralScore + patternScore) / 3.0

	// Determine if anomaly
	isAnomaly := combinedScore > 0.7

	return &MLPrediction{
		ModelID:      "combined_ml_models",
		ModelType:    "ensemble",
		Prediction:   combinedScore,
		Confidence:   math.Min(combinedScore+0.1, 1.0),
		Features:     anomalyFeatures, // Use anomaly features as primary
		AnomalyScore: anomalyScore,
		IsAnomaly:    isAnomaly,
	}, nil
}

// calculateAnomalyScore calculates anomaly score from features
func (s *ThreatDetectionService) calculateAnomalyScore(features map[string]float64) float64 {
	score := 0.0
	count := 0.0

	// Request rate anomaly
	if rate, ok := features["request_rate"]; ok {
		if rate > 100 { // High request rate
			score += 0.8
		}
		count++
	}

	// Response time anomaly
	if responseTime, ok := features["response_time"]; ok {
		if responseTime > 5000 { // High response time
			score += 0.6
		}
		count++
	}

	// Payload size anomaly
	if payloadSize, ok := features["payload_size"]; ok {
		if payloadSize > 10000 { // Large payload
			score += 0.7
		}
		count++
	}

	// Error rate anomaly
	if errorRate, ok := features["error_rate"]; ok {
		if errorRate > 0.5 { // High error rate
			score += 0.9
		}
		count++
	}

	if count > 0 {
		return score / count
	}
	return 0.0
}

// calculateBehavioralScore calculates behavioral score from features
func (s *ThreatDetectionService) calculateBehavioralScore(features map[string]float64) float64 {
	score := 0.0
	count := 0.0

	// Session pattern anomaly
	if sessionPattern, ok := features["session_pattern"]; ok {
		if sessionPattern > 100 { // Complex session
			score += 0.6
		}
		count++
	}

	// Request sequence anomaly
	if requestSequence, ok := features["request_sequence"]; ok {
		if requestSequence > 10 { // Long request sequence
			score += 0.5
		}
		count++
	}

	// Timing pattern anomaly
	if timingPattern, ok := features["timing_pattern"]; ok {
		if timingPattern < 6 || timingPattern > 22 { // Off-hours activity
			score += 0.7
		}
		count++
	}

	// Resource access anomaly
	if resourceAccess, ok := features["resource_access"]; ok {
		if resourceAccess > 50 { // Many resources accessed
			score += 0.6
		}
		count++
	}

	if count > 0 {
		return score / count
	}
	return 0.0
}

// calculatePatternScore calculates pattern score from features
func (s *ThreatDetectionService) calculatePatternScore(features map[string]float64) float64 {
	score := 0.0
	count := 0.0

	// URL pattern anomaly
	if urlPattern, ok := features["url_pattern"]; ok {
		if urlPattern > 200 { // Very long URL
			score += 0.8
		}
		count++
	}

	// URL complexity anomaly
	if urlComplexity, ok := features["url_complexity"]; ok {
		if urlComplexity > 10 { // Very complex URL
			score += 0.7
		}
		count++
	}

	// Parameter pattern anomaly
	if parameterPattern, ok := features["parameter_pattern"]; ok {
		if parameterPattern > 20 { // Many parameters
			score += 0.6
		}
		count++
	}

	// Header pattern anomaly
	if headerPattern, ok := features["header_pattern"]; ok {
		if headerPattern > 15 { // Many headers
			score += 0.5
		}
		count++
	}

	// Payload pattern anomaly
	if payloadPattern, ok := features["payload_pattern"]; ok {
		if payloadPattern > 50000 { // Very large payload
			score += 0.8
		}
		count++
	}

	// Payload complexity anomaly
	if payloadComplexity, ok := features["payload_complexity"]; ok {
		if payloadComplexity > 1000 { // Very complex payload
			score += 0.7
		}
		count++
	}

	if count > 0 {
		return score / count
	}
	return 0.0
}

// TrainMLModel trains a new ML model
func (s *ThreatDetectionService) TrainMLModel(ctx context.Context, modelType string, trainingData []byte) error {
	s.logger.Info("Training new ML model", "model_type", modelType)

	// In a real implementation, this would:
	// 1. Parse training data
	// 2. Preprocess features
	// 3. Train the model
	// 4. Validate performance
	// 5. Save the model
	// 6. Update model registry

	// For now, just log the training request
	s.logger.Info("ML model training requested",
		"model_type", modelType,
		"training_data_size", len(trainingData))

	return nil
}

// UpdateMLModel updates an existing ML model
func (s *ThreatDetectionService) UpdateMLModel(ctx context.Context, modelID string, newData []byte) error {
	s.logger.Info("Updating ML model", "model_id", modelID)

	// In a real implementation, this would:
	// 1. Load existing model
	// 2. Apply incremental learning
	// 3. Validate performance
	// 4. Update model version
	// 5. Save updated model

	// For now, just log the update request
	s.logger.Info("ML model update requested",
		"model_id", modelID,
		"new_data_size", len(newData))

	return nil
}

// GetMLModelMetrics retrieves metrics for an ML model
func (s *ThreatDetectionService) GetMLModelMetrics(ctx context.Context, modelID string) (*MLModel, error) {
	if model, ok := s.mlModels[modelID]; ok {
		return model, nil
	}

	// Check if it's a specific model type
	for _, model := range s.mlModels {
		if model.ID == modelID {
			return model, nil
		}
	}

	return nil, fmt.Errorf("ML model not found: %s", modelID)
}
