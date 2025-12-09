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

type AnomalyDetectionServiceInterface interface {
	DetectAnomalies(ctx context.Context, request *models.AnomalyDetectionRequest) (*models.AnomalyDetectionResult, error)
	GetAnomalies(ctx context.Context, filter *models.AnomalyFilter) ([]models.Anomaly, error)
	GetAnomaly(ctx context.Context, anomalyID string) (*models.Anomaly, error)
	UpdateAnomalyFeedback(ctx context.Context, feedback *models.AnomalyFeedback) error
	GetAnomalyStatistics(ctx context.Context, timeRange time.Duration) (*models.AnomalyStatistics, error)
	TrainBaselineModel(ctx context.Context, trainingData []map[string]interface{}) error
	GetModelPerformance(ctx context.Context, modelVersion string) (*models.ModelPerformanceMetric, error)
}

type AnomalyDetectionService struct {
	anomalyRepo   repository.AnomalyRepositoryInterface
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
	modelThresholds map[string]float64
}

func NewAnomalyDetectionService(
	anomalyRepo repository.AnomalyRepositoryInterface,
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
) *AnomalyDetectionService {
	return &AnomalyDetectionService{
		anomalyRepo:   anomalyRepo,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		modelThresholds: map[string]float64{
			models.DetectionEngineIsolationForest: 0.7,
			models.DetectionEngineAutoencoder:     0.8,
			models.DetectionEngineLSTM:            0.75,
			models.DetectionEngineEnsemble:        0.85,
			models.DetectionEngineStatistical:     0.6,
		},
	}
}

func (s *AnomalyDetectionService) DetectAnomalies(ctx context.Context, request *models.AnomalyDetectionRequest) (*models.AnomalyDetectionResult, error) {
	startTime := time.Now()

	result := &models.AnomalyDetectionResult{
		RequestID:       request.RequestID,
		AnomaliesFound:  false,
		Anomalies:       []models.Anomaly{},
		OverallScore:    0.0,
		Threshold:       0.7, // Default threshold
		ModelVersion:    "v1.0.0",
		FeatureScores:   make(map[string]float64),
		Recommendations: []string{},
		Metadata:        make(map[string]interface{}),
		AnalyzedAt:      time.Now(),
	}

	// Extract features from traffic data
	features, err := s.extractFeatures(request.TrafficData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract features: %w", err)
	}

	// Run different anomaly detection models
	anomalies := []models.Anomaly{}

	// 1. Statistical Anomaly Detection
	statAnomalies, err := s.detectStatisticalAnomalies(ctx, features, request)
	if err != nil {
		s.logger.Error("Statistical anomaly detection failed", "error", err)
	} else {
		anomalies = append(anomalies, statAnomalies...)
	}

	// 2. Traffic Volume Anomaly Detection
	volumeAnomalies, err := s.detectVolumeAnomalies(ctx, features, request)
	if err != nil {
		s.logger.Error("Volume anomaly detection failed", "error", err)
	} else {
		anomalies = append(anomalies, volumeAnomalies...)
	}

	// 3. Request Pattern Anomaly Detection
	patternAnomalies, err := s.detectPatternAnomalies(ctx, features, request)
	if err != nil {
		s.logger.Error("Pattern anomaly detection failed", "error", err)
	} else {
		anomalies = append(anomalies, patternAnomalies...)
	}

	// 4. Response Time Anomaly Detection
	responseTimeAnomalies, err := s.detectResponseTimeAnomalies(ctx, features, request)
	if err != nil {
		s.logger.Error("Response time anomaly detection failed", "error", err)
	} else {
		anomalies = append(anomalies, responseTimeAnomalies...)
	}

	// 5. Geolocation Anomaly Detection
	geoAnomalies, err := s.detectGeolocationAnomalies(ctx, features, request)
	if err != nil {
		s.logger.Error("Geolocation anomaly detection failed", "error", err)
	} else {
		anomalies = append(anomalies, geoAnomalies...)
	}

	// Process detected anomalies
	if len(anomalies) > 0 {
		result.AnomaliesFound = true
		result.Anomalies = anomalies

		// Calculate overall score
		totalScore := 0.0
		for _, anomaly := range anomalies {
			totalScore += anomaly.Score
			result.FeatureScores[anomaly.Type] = anomaly.Score
		}
		result.OverallScore = totalScore / float64(len(anomalies))

		// Store anomalies in database
		for _, anomaly := range anomalies {
			if err := s.anomalyRepo.CreateAnomaly(ctx, &anomaly); err != nil {
				s.logger.Error("Failed to store anomaly", "anomaly_id", anomaly.ID, "error", err)
			}
		}

		// Generate recommendations
		result.Recommendations = s.generateAnomalyRecommendations(anomalies)

		// Publish anomaly events
		if err := s.publishAnomalyEvents(ctx, anomalies); err != nil {
			s.logger.Error("Failed to publish anomaly events", "error", err)
		}
	}

	result.ProcessingTime = time.Since(startTime)
	result.Metadata["anomalies_detected"] = len(anomalies)
	result.Metadata["detection_methods"] = []string{"statistical", "volume", "pattern", "response_time", "geolocation"}

	return result, nil
}

func (s *AnomalyDetectionService) extractFeatures(trafficData map[string]interface{}) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Extract request features
	if request, ok := trafficData["request"].(map[string]interface{}); ok {
		// Request size
		if body, ok := request["body"].(string); ok {
			features["request_size"] = len(body)
		}

		// Request method
		if method, ok := request["method"].(string); ok {
			features["request_method"] = method
		}

		// Parameter count
		if params, ok := request["parameters"].(map[string]interface{}); ok {
			features["parameter_count"] = len(params)
		}

		// Header count
		if headers, ok := request["headers"].(map[string]interface{}); ok {
			features["header_count"] = len(headers)
		}

		// Path depth
		if path, ok := request["path"].(string); ok {
			features["path_depth"] = len(strings.Split(path, "/")) - 1
		}

		// User agent entropy
		if userAgent, ok := request["user_agent"].(string); ok {
			features["user_agent_entropy"] = s.calculateEntropy(userAgent)
		}
	}

	// Extract response features
	if response, ok := trafficData["response"].(map[string]interface{}); ok {
		// Response size
		if size, ok := response["size"].(float64); ok {
			features["response_size"] = size
		}

		// Status code
		if statusCode, ok := response["status_code"].(float64); ok {
			features["status_code"] = statusCode
		}

		// Response time
		if responseTime, ok := response["response_time"].(float64); ok {
			features["response_time"] = responseTime
		}
	}

	// Extract temporal features
	if timestamp, ok := trafficData["timestamp"].(time.Time); ok {
		features["hour_of_day"] = timestamp.Hour()
		features["day_of_week"] = int(timestamp.Weekday())
		features["day_of_month"] = timestamp.Day()
	}

	// Extract IP-based features
	if ipAddr, ok := trafficData["ip_address"].(string); ok {
		features["ip_address"] = ipAddr
		// Add geolocation features if available
		if geo, ok := trafficData["geolocation"].(map[string]interface{}); ok {
			if country, ok := geo["country"].(string); ok {
				features["country"] = country
			}
			if city, ok := geo["city"].(string); ok {
				features["city"] = city
			}
		}
	}

	return features, nil
}

func (s *AnomalyDetectionService) detectStatisticalAnomalies(ctx context.Context, features map[string]interface{}, request *models.AnomalyDetectionRequest) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly

	// Get baseline statistics for comparison
	baseline, err := s.anomalyRepo.GetBaselineStatistics(ctx, "statistical", "entityType")
	if err != nil {
		return anomalies, fmt.Errorf("failed to get baseline statistics: %w", err)
	}

	// Check numerical features for statistical anomalies
	numericalFeatures := []string{"request_size", "response_size", "response_time", "parameter_count", "header_count"}

	for _, feature := range numericalFeatures {
		if value, ok := features[feature]; ok {
			if numValue, ok := value.(float64); ok {
				score := s.calculateZScore(numValue, baseline, feature)
				
				// Threshold for statistical anomaly (Z-score > 3)
				if math.Abs(score) > 3.0 {
					anomaly := models.Anomaly{
						ID:              uuid.New().String(),
						Type:            models.AnomalyTypeTrafficVolume,
						Severity:        s.getSeverityFromScore(math.Abs(score)),
						Score:           math.Abs(score),
						Threshold:       3.0,
						Title:           fmt.Sprintf("Statistical Anomaly in %s", feature),
						Description:     fmt.Sprintf("Unusual %s detected: %.2f (Z-score: %.2f)", feature, numValue, score),
						DetectionEngine: models.DetectionEngineStatistical,
						Confidence:      s.calculateConfidence(math.Abs(score)),
						Features: []models.AnomalyFeature{
							{
								Name:           feature,
								Value:          numValue,
								BaselineValue:  baseline[feature+"_mean"],
								DeviationScore: math.Abs(score),
								Weight:         1.0,
								Description:    fmt.Sprintf("Z-score deviation for %s", feature),
							},
						},
						ModelVersion:  "statistical_v1.0",
						Status:        models.AnomalyStatusNew,
						FirstDetected: time.Now(),
						LastDetected:  time.Now(),
						Count:         1,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					}

					// Add context from request
					if apiID, ok := request.TrafficData["api_id"].(string); ok {
						anomaly.APIID = apiID
					}
					if ipAddr, ok := request.TrafficData["ip_address"].(string); ok {
						anomaly.IPAddress = ipAddr
					}

					anomalies = append(anomalies, anomaly)
				}
			}
		}
	}

	return anomalies, nil
}

func (s *AnomalyDetectionService) detectVolumeAnomalies(ctx context.Context, features map[string]interface{}, request *models.AnomalyDetectionRequest) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly

	// Check for unusual request volume patterns
	if ipAddr, ok := features["ip_address"].(string); ok {
		// Get recent request count from this IP
		requestCount, err := s.anomalyRepo.GetRecentRequestCount(ctx, ipAddr, "ip", time.Hour)
		if err != nil {
			return anomalies, fmt.Errorf("failed to get request count: %w", err)
		}

		// Get baseline request count for this IP
		baselineCount, err := s.anomalyRepo.GetBaselineRequestCount(ctx, ipAddr, "ip")
		if err != nil {
			return anomalies, fmt.Errorf("failed to get baseline request count: %w", err)
		}

		// Calculate volume anomaly score
		if baselineCount > 0 {
			volumeRatio := float64(requestCount) / float64(baselineCount)
			
			// Anomaly if volume is 5x higher or 10x lower than baseline
			if volumeRatio > 5.0 || volumeRatio < 0.1 {
				severity := models.AnomalySeverityMedium
				if volumeRatio > 10.0 || volumeRatio < 0.05 {
					severity = models.AnomalySeverityHigh
				}

				anomaly := models.Anomaly{
					ID:              uuid.New().String(),
					Type:            models.AnomalyTypeTrafficVolume,
					Severity:        severity,
					Score:           math.Max(volumeRatio, 1.0/volumeRatio),
					Threshold:       5.0,
					Title:           "Traffic Volume Anomaly",
					Description:     fmt.Sprintf("Unusual traffic volume from IP %s: %d requests (%.2fx baseline)", ipAddr, requestCount, volumeRatio),
					IPAddress:       ipAddr,
					DetectionEngine: models.DetectionEngineStatistical,
					Confidence:      s.calculateVolumeConfidence(volumeRatio),
					Features: []models.AnomalyFeature{
						{
							Name:           "request_volume",
							Value:          requestCount,
							BaselineValue:  baselineCount,
							DeviationScore: math.Max(volumeRatio, 1.0/volumeRatio),
							Weight:         1.0,
							Description:    "Request volume deviation from baseline",
						},
					},
					ModelVersion:  "volume_v1.0",
					Status:        models.AnomalyStatusNew,
					FirstDetected: time.Now(),
					LastDetected:  time.Now(),
					Count:         1,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				if apiID, ok := request.TrafficData["api_id"].(string); ok {
					anomaly.APIID = apiID
				}

				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies, nil
}

func (s *AnomalyDetectionService) detectPatternAnomalies(ctx context.Context, features map[string]interface{}, request *models.AnomalyDetectionRequest) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly

	// Check for unusual request patterns
	if requestData, ok := request.TrafficData["request"].(map[string]interface{}); ok {
		// Check for unusual parameter patterns
		if params, ok := requestData["parameters"].(map[string]interface{}); ok {
			for key, value := range params {
				valueStr := fmt.Sprintf("%v", value)
				
				// Check for unusual parameter values
				if s.isUnusualParameterValue(valueStr) {
					anomaly := models.Anomaly{
						ID:              uuid.New().String(),
						Type:            models.AnomalyTypeParameterValues,
						Severity:        models.AnomalySeverityMedium,
						Score:           0.8,
						Threshold:       0.7,
						Title:           "Unusual Parameter Pattern",
						Description:     fmt.Sprintf("Unusual parameter value detected in '%s'", key),
						DetectionEngine: models.DetectionEngineStatistical,
						Confidence:      0.75,
						Features: []models.AnomalyFeature{
							{
								Name:           fmt.Sprintf("parameter_%s", key),
								Value:          valueStr,
								DeviationScore: 0.8,
								Weight:         1.0,
								Description:    "Unusual parameter value pattern",
							},
						},
						ModelVersion:  "pattern_v1.0",
						Status:        models.AnomalyStatusNew,
						FirstDetected: time.Now(),
						LastDetected:  time.Now(),
						Count:         1,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					}

					if apiID, ok := request.TrafficData["api_id"].(string); ok {
						anomaly.APIID = apiID
					}
					if ipAddr, ok := request.TrafficData["ip_address"].(string); ok {
						anomaly.IPAddress = ipAddr
					}

					anomalies = append(anomalies, anomaly)
				}
			}
		}

		// Check for unusual header patterns
		if headers, ok := requestData["headers"].(map[string]interface{}); ok {
			userAgent, hasUA := headers["user-agent"].(string)
			if hasUA && s.isUnusualUserAgent(userAgent) {
				anomaly := models.Anomaly{
					ID:              uuid.New().String(),
					Type:            models.AnomalyTypeHeaderPattern,
					Severity:        models.AnomalySeverityLow,
					Score:           0.6,
					Threshold:       0.5,
					Title:           "Unusual User Agent",
					Description:     "Suspicious or unusual user agent detected",
					DetectionEngine: models.DetectionEngineStatistical,
					Confidence:      0.65,
					Features: []models.AnomalyFeature{
						{
							Name:           "user_agent",
							Value:          userAgent,
							DeviationScore: 0.6,
							Weight:         0.8,
							Description:    "Unusual user agent pattern",
						},
					},
					ModelVersion:  "pattern_v1.0",
					Status:        models.AnomalyStatusNew,
					FirstDetected: time.Now(),
					LastDetected:  time.Now(),
					Count:         1,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				if apiID, ok := request.TrafficData["api_id"].(string); ok {
					anomaly.APIID = apiID
				}
				if ipAddr, ok := request.TrafficData["ip_address"].(string); ok {
					anomaly.IPAddress = ipAddr
				}

				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies, nil
}

func (s *AnomalyDetectionService) detectResponseTimeAnomalies(ctx context.Context, features map[string]interface{}, request *models.AnomalyDetectionRequest) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly

	if responseTime, ok := features["response_time"].(float64); ok {
		// Get baseline response time for this endpoint
		var endpointID string
		if apiID, ok := request.TrafficData["api_id"].(string); ok {
			if requestData, ok := request.TrafficData["request"].(map[string]interface{}); ok {
				if path, ok := requestData["path"].(string); ok {
					endpointID = fmt.Sprintf("%s:%s", apiID, path)
				}
			}
		}

		if endpointID != "" {
			baselineRT, err := s.anomalyRepo.GetBaselineResponseTime(ctx, endpointID, "endpoint")
			if err != nil {
				return anomalies, fmt.Errorf("failed to get baseline response time: %w", err)
			}

			if baselineRT > 0 {
				rtRatio := responseTime / baselineRT
				
				// Anomaly if response time is 3x higher than baseline
				if rtRatio > 3.0 {
					severity := models.AnomalySeverityMedium
					if rtRatio > 10.0 {
						severity = models.AnomalySeverityHigh
					}

					anomaly := models.Anomaly{
						ID:              uuid.New().String(),
						Type:            models.AnomalyTypeResponseTime,
						Severity:        severity,
						Score:           rtRatio,
						Threshold:       3.0,
						Title:           "Response Time Anomaly",
						Description:     fmt.Sprintf("Unusually slow response time: %.2fms (%.2fx baseline)", responseTime, rtRatio),
						DetectionEngine: models.DetectionEngineStatistical,
						Confidence:      s.calculateResponseTimeConfidence(rtRatio),
						Features: []models.AnomalyFeature{
							{
								Name:           "response_time",
								Value:          responseTime,
								BaselineValue:  baselineRT,
								DeviationScore: rtRatio,
								Weight:         1.0,
								Description:    "Response time deviation from baseline",
							},
						},
						ModelVersion:  "response_time_v1.0",
						Status:        models.AnomalyStatusNew,
						FirstDetected: time.Now(),
						LastDetected:  time.Now(),
						Count:         1,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					}

					if apiID, ok := request.TrafficData["api_id"].(string); ok {
						anomaly.APIID = apiID
					}
					if ipAddr, ok := request.TrafficData["ip_address"].(string); ok {
						anomaly.IPAddress = ipAddr
					}

					anomalies = append(anomalies, anomaly)
				}
			}
		}
	}

	return anomalies, nil
}

func (s *AnomalyDetectionService) detectGeolocationAnomalies(ctx context.Context, features map[string]interface{}, request *models.AnomalyDetectionRequest) ([]models.Anomaly, error) {
	var anomalies []models.Anomaly

	if country, ok := features["country"].(string); ok {
		if ipAddr, ok := features["ip_address"].(string); ok {
			// Get historical countries for this IP
			historicalCountries, err := s.anomalyRepo.GetHistoricalCountries(ctx, ipAddr, "ip")
			if err != nil {
				return anomalies, fmt.Errorf("failed to get historical countries: %w", err)
			}

			// Check if current country is unusual for this IP
			isUnusual := true
			for _, histCountry := range historicalCountries {
				if histCountry == country {
					isUnusual = false
					break
				}
			}

			if isUnusual && len(historicalCountries) > 0 {
				anomaly := models.Anomaly{
					ID:              uuid.New().String(),
					Type:            models.AnomalyTypeGeolocation,
					Severity:        models.AnomalySeverityMedium,
					Score:           0.8,
					Threshold:       0.7,
					Title:           "Geolocation Anomaly",
					Description:     fmt.Sprintf("Unusual location detected for IP %s: %s", ipAddr, country),
					IPAddress:       ipAddr,
					DetectionEngine: models.DetectionEngineStatistical,
					Confidence:      0.75,
					Features: []models.AnomalyFeature{
						{
							Name:           "country",
							Value:          country,
							DeviationScore: 0.8,
							Weight:         1.0,
							Description:    "Unusual country for this IP address",
						},
					},
					ModelVersion:  "geolocation_v1.0",
					Status:        models.AnomalyStatusNew,
					FirstDetected: time.Now(),
					LastDetected:  time.Now(),
					Count:         1,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				if apiID, ok := request.TrafficData["api_id"].(string); ok {
					anomaly.APIID = apiID
				}

				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies, nil
}

// Helper functions

func (s *AnomalyDetectionService) calculateEntropy(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	frequency := make(map[rune]int)
	for _, char := range text {
		frequency[char]++
	}

	entropy := 0.0
	length := float64(len(text))
	for _, count := range frequency {
		p := float64(count) / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

func (s *AnomalyDetectionService) calculateZScore(value float64, baseline map[string]interface{}, feature string) float64 {
	meanKey := feature + "_mean"
	stdKey := feature + "_std"

	mean, ok1 := baseline[meanKey].(float64)
	std, ok2 := baseline[stdKey].(float64)

	if !ok1 || !ok2 || std == 0 {
		return 0
	}

	return (value - mean) / std
}

func (s *AnomalyDetectionService) getSeverityFromScore(score float64) string {
	if score >= 5.0 {
		return models.AnomalySeverityCritical
	} else if score >= 4.0 {
		return models.AnomalySeverityHigh
	} else if score >= 3.0 {
		return models.AnomalySeverityMedium
	} else if score >= 2.0 {
		return models.AnomalySeverityLow
	}
	return models.AnomalySeverityInfo
}

func (s *AnomalyDetectionService) calculateConfidence(score float64) float64 {
	// Sigmoid function to convert score to confidence (0-1)
	return 1.0 / (1.0 + math.Exp(-score+3.0))
}

func (s *AnomalyDetectionService) calculateVolumeConfidence(ratio float64) float64 {
	deviation := math.Max(ratio, 1.0/ratio)
	return math.Min(0.95, deviation/10.0)
}

func (s *AnomalyDetectionService) calculateResponseTimeConfidence(ratio float64) float64 {
	return math.Min(0.9, ratio/10.0)
}

func (s *AnomalyDetectionService) isUnusualParameterValue(value string) bool {
	// Check for suspicious patterns in parameter values
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"eval(",
		"union select",
		"drop table",
		"../",
		"..\\",
		"cmd.exe",
		"/bin/sh",
		"base64",
	}

	valueLower := strings.ToLower(value)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(valueLower, pattern) {
			return true
		}
	}

	// Check for unusual length
	if len(value) > 1000 {
		return true
	}

	// Check for high entropy (random-looking strings)
	entropy := s.calculateEntropy(value)
	if entropy > 4.5 && len(value) > 20 {
		return true
	}

	return false
}

func (s *AnomalyDetectionService) isUnusualUserAgent(userAgent string) bool {
	// Check for suspicious user agent patterns
	suspiciousPatterns := []string{
		"sqlmap",
		"nikto",
		"nmap",
		"burp",
		"scanner",
		"bot",
		"crawler",
		"spider",
		"wget",
		"curl",
		"python",
		"perl",
		"php",
	}

	userAgentLower := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgentLower, pattern) {
			return true
		}
	}

	// Check for very short or very long user agents
	if len(userAgent) < 10 || len(userAgent) > 500 {
		return true
	}

	return false
}

func (s *AnomalyDetectionService) generateAnomalyRecommendations(anomalies []models.Anomaly) []string {
	recommendations := []string{}
	anomalyTypes := make(map[string]bool)

	for _, anomaly := range anomalies {
		anomalyTypes[anomaly.Type] = true
	}

	if anomalyTypes[models.AnomalyTypeTrafficVolume] {
		recommendations = append(recommendations, "Implement rate limiting to control traffic volume")
		recommendations = append(recommendations, "Monitor and alert on unusual traffic patterns")
	}

	if anomalyTypes[models.AnomalyTypeResponseTime] {
		recommendations = append(recommendations, "Investigate performance bottlenecks")
		recommendations = append(recommendations, "Consider implementing caching or load balancing")
	}

	if anomalyTypes[models.AnomalyTypeParameterValues] {
		recommendations = append(recommendations, "Implement input validation and sanitization")
		recommendations = append(recommendations, "Enable Web Application Firewall (WAF) protection")
	}

	if anomalyTypes[models.AnomalyTypeGeolocation] {
		recommendations = append(recommendations, "Review access patterns from unusual locations")
		recommendations = append(recommendations, "Consider implementing geo-blocking for high-risk regions")
	}

	if anomalyTypes[models.AnomalyTypeHeaderPattern] {
		recommendations = append(recommendations, "Monitor for automated tools and bots")
		recommendations = append(recommendations, "Implement user agent filtering and validation")
	}

	return recommendations
}

func (s *AnomalyDetectionService) publishAnomalyEvents(ctx context.Context, anomalies []models.Anomaly) error {
	for _, anomaly := range anomalies {
		eventData := map[string]interface{}{
			"event_type":       "anomaly_detected",
			"anomaly_id":       anomaly.ID,
			"anomaly_type":     anomaly.Type,
			"severity":         anomaly.Severity,
			"score":            anomaly.Score,
			"confidence":       anomaly.Confidence,
			"ip_address":       anomaly.IPAddress,
			"api_id":           anomaly.APIID,
			"endpoint_id":      anomaly.EndpointID,
			"detection_engine": anomaly.DetectionEngine,
			"timestamp":        anomaly.CreatedAt,
			"features":         anomaly.Features,
			"description":      anomaly.Description,
		}

		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			return fmt.Errorf("failed to marshal anomaly event: %w", err)
		}

		message := kafka.Message{
			Topic: "anomaly_events",
			Key:   []byte(anomaly.ID),
			Value: eventJSON,
		}

		if err := s.kafkaProducer.Produce(ctx, message); err != nil {
			return fmt.Errorf("failed to produce anomaly event: %w", err)
		}
	}

	return nil
}

func (s *AnomalyDetectionService) GetAnomalies(ctx context.Context, filter *models.AnomalyFilter) ([]models.Anomaly, error) {
	return s.anomalyRepo.GetAnomalies(ctx, filter)
}

func (s *AnomalyDetectionService) GetAnomaly(ctx context.Context, anomalyID string) (*models.Anomaly, error) {
	return s.anomalyRepo.GetAnomaly(ctx, anomalyID)
}

func (s *AnomalyDetectionService) UpdateAnomalyFeedback(ctx context.Context, feedback *models.AnomalyFeedback) error {
	return s.anomalyRepo.UpdateAnomalyFeedback(ctx, feedback)
}

func (s *AnomalyDetectionService) GetAnomalyStatistics(ctx context.Context, timeRange time.Duration) (*models.AnomalyStatistics, error) {
	filter := &models.AnomalyFilter{
		DateFrom: time.Now().Add(-timeRange),
		DateTo:   time.Now(),
	}
	return s.anomalyRepo.GetAnomalyStatistics(ctx, filter)
}

func (s *AnomalyDetectionService) TrainBaselineModel(ctx context.Context, trainingData []map[string]interface{}) error {
	// Extract features from training data
	var features []map[string]interface{}
	for _, data := range trainingData {
		extractedFeatures, err := s.extractFeatures(data)
		if err != nil {
			s.logger.Error("Failed to extract features from training data", "error", err)
			continue
		}
		features = append(features, extractedFeatures)
	}

	// Calculate baseline statistics
	baseline := s.calculateBaselineStatistics(features)

	// Store baseline in repository
	return s.anomalyRepo.StoreBaselineStatistics(ctx, "statistical", baseline)
}

func (s *AnomalyDetectionService) calculateBaselineStatistics(features []map[string]interface{}) map[string]interface{} {
	baseline := make(map[string]interface{})
	numericalFeatures := []string{"request_size", "response_size", "response_time", "parameter_count", "header_count"}

	for _, feature := range numericalFeatures {
		var values []float64
		for _, featureSet := range features {
			if value, ok := featureSet[feature].(float64); ok {
				values = append(values, value)
			}
		}

		if len(values) > 0 {
			mean := s.calculateMean(values)
			std := s.calculateStandardDeviation(values, mean)
			baseline[feature+"_mean"] = mean
			baseline[feature+"_std"] = std
		}
	}

	return baseline
}

func (s *AnomalyDetectionService) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func (s *AnomalyDetectionService) calculateStandardDeviation(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	sumSquaredDiff := 0.0
	for _, value := range values {
		diff := value - mean
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(values)-1)
	return math.Sqrt(variance)
}

func (s *AnomalyDetectionService) GetModelPerformance(ctx context.Context, modelVersion string) (*models.ModelPerformanceMetric, error) {
	return s.anomalyRepo.GetModelPerformance(ctx, modelVersion)
}
