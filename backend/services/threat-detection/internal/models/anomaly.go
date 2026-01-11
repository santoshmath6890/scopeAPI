package models

import (
	"time"
)

type Anomaly struct {
	ID              string                 `json:"id" db:"id"`
	Type            string                 `json:"type" db:"type"`
	Severity        string                 `json:"severity" db:"severity"`
	Score           float64                `json:"score" db:"score"`
	Threshold       float64                `json:"threshold" db:"threshold"`
	Title           string                 `json:"title" db:"title"`
	Description     string                 `json:"description" db:"description"`
	APIID           string                 `json:"api_id" db:"api_id"`
	EndpointID      string                 `json:"endpoint_id" db:"endpoint_id"`
	IPAddress       string                 `json:"ip_address" db:"ip_address"`
	UserID          string                 `json:"user_id,omitempty" db:"user_id"`
	SessionID       string                 `json:"session_id,omitempty" db:"session_id"`
	RequestPattern  map[string]interface{} `json:"request_pattern" db:"request_pattern"`
	BaselineData    map[string]interface{} `json:"baseline_data" db:"baseline_data"`
	DeviationData   map[string]interface{} `json:"deviation_data" db:"deviation_data"`
	Features        []AnomalyFeature       `json:"features" db:"features"`
	ModelVersion    string                 `json:"model_version" db:"model_version"`
	DetectionEngine string                 `json:"detection_engine" db:"detection_engine"`
	Confidence      float64                `json:"confidence" db:"confidence"`
	FalsePositive   bool                   `json:"false_positive" db:"false_positive"`
	Feedback        string                 `json:"feedback,omitempty" db:"feedback"`
	Status          string                 `json:"status" db:"status"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	FirstDetected   time.Time              `json:"first_detected" db:"first_detected"`
	LastDetected    time.Time              `json:"last_detected" db:"last_detected"`
	Count           int64                  `json:"count" db:"count"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

type AnomalyFeature struct {
	Name           string      `json:"name"`
	Value          interface{} `json:"value"`
	BaselineValue  interface{} `json:"baseline_value"`
	DeviationScore float64     `json:"deviation_score"`
	Weight         float64     `json:"weight"`
	Description    string      `json:"description,omitempty"`
}

type AnomalyDetectionRequest struct {
	TrafficData     map[string]interface{} `json:"traffic_data"`
	RequestID       string                 `json:"request_id"`
	Timestamp       time.Time              `json:"timestamp"`
	ModelType       string                 `json:"model_type,omitempty"`
	Sensitivity     float64                `json:"sensitivity,omitempty"`
	IncludeFeatures bool                   `json:"include_features,omitempty"`
	Configuration   map[string]interface{} `json:"configuration,omitempty"`
}

type AnomalyDetectionResult struct {
	RequestID       string                 `json:"request_id"`
	AnomaliesFound  bool                   `json:"anomalies_found"`
	Anomalies       []Anomaly              `json:"anomalies"`
	OverallScore    float64                `json:"overall_score"`
	Threshold       float64                `json:"threshold"`
	ModelVersion    string                 `json:"model_version"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	FeatureScores   map[string]float64     `json:"feature_scores,omitempty"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
	AnalyzedAt      time.Time              `json:"analyzed_at"`
}

type AnomalyFilter struct {
	Type            []string  `json:"type,omitempty"`
	Severity        []string  `json:"severity,omitempty"`
	Status          []string  `json:"status,omitempty"`
	APIID           string    `json:"api_id,omitempty"`
	EndpointID      string    `json:"endpoint_id,omitempty"`
	IPAddress       string    `json:"ip_address,omitempty"`
	UserID          string    `json:"user_id,omitempty"`
	MinScore        float64   `json:"min_score,omitempty"`
	MaxScore        float64   `json:"max_score,omitempty"`
	FalsePositive   *bool     `json:"false_positive,omitempty"`
	DateFrom        time.Time `json:"date_from,omitempty"`
	DateTo          time.Time `json:"date_to,omitempty"`
	ModelVersion    string    `json:"model_version,omitempty"`
	DetectionEngine string    `json:"detection_engine,omitempty"`
	Page            int       `json:"page"`
	Limit           int       `json:"limit"`
	SortBy          string    `json:"sort_by"`
	SortOrder       string    `json:"sort_order"`
}

type AnomalyFeedback struct {
	AnomalyID     string                 `json:"anomaly_id"`
	FalsePositive bool                   `json:"false_positive"`
	Feedback      string                 `json:"feedback"`
	UserID        string                 `json:"user_id"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

type AnomalyStatistics struct {
	TotalAnomalies      int64                    `json:"total_anomalies"`
	ActiveAnomalies     int64                    `json:"active_anomalies"`
	ResolvedAnomalies   int64                    `json:"resolved_anomalies"`
	FalsePositives      int64                    `json:"false_positives"`
	CriticalAnomalies   int64                    `json:"critical_anomalies"`
	HighAnomalies       int64                    `json:"high_anomalies"`
	MediumAnomalies     int64                    `json:"medium_anomalies"`
	LowAnomalies        int64                    `json:"low_anomalies"`
	AnomaliesByType     map[string]int64         `json:"anomalies_by_type"`
	AnomaliesBySeverity map[string]int64         `json:"anomalies_by_severity"`
	AnomaliesByStatus   map[string]int64         `json:"anomalies_by_status"`
	AnomaliesByEngine   map[string]int64         `json:"anomalies_by_engine"`
	RecentAnomalies     int64                    `json:"recent_anomalies"`
	TrendData           []AnomalyTrendPoint      `json:"trend_data"`
	TopAnomalousAPIs    []APIAnomalySummary      `json:"top_anomalous_apis"`
	TopAnomalousIPs     []IPAnomalySummary       `json:"top_anomalous_ips"`
	ModelPerformance    []ModelPerformanceMetric `json:"model_performance"`
	AverageScore        float64                  `json:"average_score"`
	ScoreDistribution   map[string]int64         `json:"score_distribution"`
	GeneratedAt         time.Time                `json:"generated_at"`
}

type AnomalyTrendPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	AnomalyCount int64     `json:"anomaly_count"`
	AverageScore float64   `json:"average_score"`
	Type         string    `json:"type,omitempty"`
	Severity     string    `json:"severity,omitempty"`
}

type APIAnomalySummary struct {
	APIID        string    `json:"api_id"`
	APIName      string    `json:"api_name"`
	AnomalyCount int64     `json:"anomaly_count"`
	AverageScore float64   `json:"average_score"`
	LastAnomaly  time.Time `json:"last_anomaly"`
}

type IPAnomalySummary struct {
	IPAddress    string    `json:"ip_address"`
	AnomalyCount int64     `json:"anomaly_count"`
	AverageScore float64   `json:"average_score"`
	LastAnomaly  time.Time `json:"last_anomaly"`
	Country      string    `json:"country,omitempty"`
	ISP          string    `json:"isp,omitempty"`
}

type ModelPerformanceMetric struct {
	ModelVersion      string    `json:"model_version"`
	Engine            string    `json:"engine"`
	Accuracy          float64   `json:"accuracy"`
	Precision         float64   `json:"precision"`
	Recall            float64   `json:"recall"`
	F1Score           float64   `json:"f1_score"`
	FalsePositiveRate float64   `json:"false_positive_rate"`
	FalseNegativeRate float64   `json:"false_negative_rate"`
	TruePositiveRate  float64   `json:"true_positive_rate"`
	LastUpdated       time.Time `json:"last_updated"`
	SampleSize        int64     `json:"sample_size"`
}

// Anomaly types
const (
	AnomalyTypeTrafficVolume   = "traffic_volume"
	AnomalyTypeRequestPattern  = "request_pattern"
	AnomalyTypeResponseTime    = "response_time"
	AnomalyTypeErrorRate       = "error_rate"
	AnomalyTypeUserBehavior    = "user_behavior"
	AnomalyTypeDataAccess      = "data_access"
	AnomalyTypeGeolocation     = "geolocation"
	AnomalyTypeTimePattern     = "time_pattern"
	AnomalyTypeParameterValues = "parameter_values"
	AnomalyTypeHeaderPattern   = "header_pattern"
	AnomalyTypePayloadSize     = "payload_size"
	AnomalyTypeSessionPattern  = "session_pattern"
)

// Anomaly severity levels
const (
	AnomalySeverityCritical = "critical"
	AnomalySeverityHigh     = "high"
	AnomalySeverityMedium   = "medium"
	AnomalySeverityLow      = "low"
	AnomalySeverityInfo     = "info"
)

// Anomaly status
const (
	AnomalyStatusNew        = "new"
	AnomalyStatusInProgress = "in_progress"
	AnomalyStatusResolved   = "resolved"
	AnomalyStatusFalsePos   = "false_positive"
	AnomalyStatusIgnored    = "ignored"
)

// Detection engines
const (
	DetectionEngineIsolationForest = "isolation_forest"
	DetectionEngineAutoencoder     = "autoencoder"
	DetectionEngineLSTM            = "lstm"
	DetectionEngineEnsemble        = "ensemble"
	DetectionEngineStatistical     = "statistical"
	DetectionEngineKMeans          = "kmeans"
	DetectionEngineSVM             = "svm"
)
