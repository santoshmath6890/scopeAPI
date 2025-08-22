package models

import (
	"time"
)

type BehaviorPattern struct {
	ID              string                 `json:"id" db:"id"`
	Name            string                 `json:"name" db:"name"`
	Type            string                 `json:"type" db:"type"`
	Category        string                 `json:"category" db:"category"`
	Description     string                 `json:"description" db:"description"`
	UserID          string                 `json:"user_id,omitempty" db:"user_id"`
	IPAddress       string                 `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent       string                 `json:"user_agent,omitempty" db:"user_agent"`
	SessionID       string                 `json:"session_id,omitempty" db:"session_id"`
	APIID           string                 `json:"api_id,omitempty" db:"api_id"`
	EndpointPattern string                 `json:"endpoint_pattern,omitempty" db:"endpoint_pattern"`
	TimePattern     TimePattern            `json:"time_pattern" db:"time_pattern"`
	RequestPattern  RequestPattern         `json:"request_pattern" db:"request_pattern"`
	ResponsePattern ResponsePattern        `json:"response_pattern" db:"response_pattern"`
	Frequency       FrequencyPattern       `json:"frequency" db:"frequency"`
	Sequence        []SequenceStep         `json:"sequence,omitempty" db:"sequence"`
	Baseline        BaselineMetrics        `json:"baseline" db:"baseline"`
	Thresholds      PatternThresholds      `json:"thresholds" db:"thresholds"`
	RiskScore       float64                `json:"risk_score" db:"risk_score"`
	Confidence      float64                `json:"confidence" db:"confidence"`
	IsNormal        bool                   `json:"is_normal" db:"is_normal"`
	IsSuspicious    bool                   `json:"is_suspicious" db:"is_suspicious"`
	IsMalicious     bool                   `json:"is_malicious" db:"is_malicious"`
	Tags            []string               `json:"tags" db:"tags"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	FirstSeen       time.Time              `json:"first_seen" db:"first_seen"`
	LastSeen        time.Time              `json:"last_seen" db:"last_seen"`
	Occurrences     int64                  `json:"occurrences" db:"occurrences"`
	IsActive        bool                   `json:"is_active" db:"is_active"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

type TimePattern struct {
	HourlyDistribution  map[int]float64 `json:"hourly_distribution"`
	DailyDistribution   map[int]float64 `json:"daily_distribution"`
	WeeklyDistribution  map[int]float64 `json:"weekly_distribution"`
	MonthlyDistribution map[int]float64 `json:"monthly_distribution"`
	PeakHours          []int           `json:"peak_hours"`
	OffHours           []int           `json:"off_hours"`
	Timezone           string          `json:"timezone"`
	Seasonality        string          `json:"seasonality,omitempty"`
}

type RequestPattern struct {
	Methods            map[string]int64       `json:"methods"`
	ContentTypes       map[string]int64       `json:"content_types"`
	HeaderPatterns     map[string]interface{} `json:"header_patterns"`
	ParameterPatterns  map[string]interface{} `json:"parameter_patterns"`
	PayloadPatterns    map[string]interface{} `json:"payload_patterns"`
	SizeDistribution   SizeDistribution       `json:"size_distribution"`
	EncodingPatterns   map[string]int64       `json:"encoding_patterns"`
	AuthenticationUsed bool                   `json:"authentication_used"`
	EncryptionUsed     bool                   `json:"encryption_used"`
}

type ResponsePattern struct {
	StatusCodes        map[int]int64          `json:"status_codes"`
	ContentTypes       map[string]int64       `json:"content_types"`
	SizeDistribution   SizeDistribution       `json:"size_distribution"`
	ResponseTimes      ResponseTimePattern    `json:"response_times"`
	HeaderPatterns     map[string]interface{} `json:"header_patterns"`
	ErrorPatterns      map[string]int64       `json:"error_patterns"`
	CachePatterns      map[string]int64       `json:"cache_patterns"`
}

type FrequencyPattern struct {
	RequestsPerSecond  float64           `json:"requests_per_second"`
	RequestsPerMinute  float64           `json:"requests_per_minute"`
	RequestsPerHour    float64           `json:"requests_per_hour"`
	RequestsPerDay     float64           `json:"requests_per_day"`
	BurstPattern       BurstPattern      `json:"burst_pattern"`
	IntervalPattern    IntervalPattern   `json:"interval_pattern"`
	VelocityPattern    VelocityPattern   `json:"velocity_pattern"`
}

type SequenceStep struct {
	Order       int                    `json:"order"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	ExpectedDelay time.Duration        `json:"expected_delay,omitempty"`
	Optional    bool                   `json:"optional"`
	Weight      float64                `json:"weight"`
}

type BaselineMetrics struct {
	AverageRequestsPerHour   float64                `json:"average_requests_per_hour"`
	AverageResponseTime      time.Duration          `json:"average_response_time"`
	AveragePayloadSize       int64                  `json:"average_payload_size"`
	CommonEndpoints          map[string]float64     `json:"common_endpoints"`
	CommonParameters         map[string]float64     `json:"common_parameters"`
	CommonHeaders            map[string]float64     `json:"common_headers"`
	TypicalErrorRate         float64                `json:"typical_error_rate"`
	TypicalGeolocations      map[string]float64     `json:"typical_geolocations"`
	EstablishedAt            time.Time              `json:"established_at"`
	LastUpdated              time.Time              `json:"last_updated"`
	SampleSize               int64                  `json:"sample_size"`
	ConfidenceLevel          float64                `json:"confidence_level"`
}

type PatternThresholds struct {
	FrequencyThreshold     float64 `json:"frequency_threshold"`
	VelocityThreshold      float64 `json:"velocity_threshold"`
	SizeThreshold          int64   `json:"size_threshold"`
	ResponseTimeThreshold  time.Duration `json:"response_time_threshold"`
	ErrorRateThreshold     float64 `json:"error_rate_threshold"`
	DeviationThreshold     float64 `json:"deviation_threshold"`
	AnomalyScoreThreshold  float64 `json:"anomaly_score_threshold"`
}

type SizeDistribution struct {
	Min        int64   `json:"min"`
	Max        int64   `json:"max"`
	Mean       float64 `json:"mean"`
	Median     float64 `json:"median"`
	StdDev     float64 `json:"std_dev"`
	Percentiles map[string]int64 `json:"percentiles"`
}

type ResponseTimePattern struct {
	Min         time.Duration        `json:"min"`
	Max         time.Duration        `json:"max"`
	Mean        time.Duration        `json:"mean"`
	Median      time.Duration        `json:"median"`
	StdDev      time.Duration        `json:"std_dev"`
	Percentiles map[string]time.Duration `json:"percentiles"`
	Distribution map[string]int64    `json:"distribution"`
}

type BurstPattern struct {
	MaxBurstSize     int64         `json:"max_burst_size"`
	AverageBurstSize float64       `json:"average_burst_size"`
	BurstDuration    time.Duration `json:"burst_duration"`
	BurstFrequency   float64       `json:"burst_frequency"`
	QuietPeriods     time.Duration `json:"quiet_periods"`
}

type IntervalPattern struct {
	MinInterval     time.Duration `json:"min_interval"`
	MaxInterval     time.Duration `json:"max_interval"`
	AverageInterval time.Duration `json:"average_interval"`
	Regularity      float64       `json:"regularity"`
	Jitter          time.Duration `json:"jitter"`
}

type VelocityPattern struct {
	Acceleration    float64 `json:"acceleration"`
	Deceleration    float64 `json:"deceleration"`
	PeakVelocity    float64 `json:"peak_velocity"`
	AverageVelocity float64 `json:"average_velocity"`
	VelocityVariance float64 `json:"velocity_variance"`
}

type BehaviorAnalysisRequest struct {
	TrafficData    map[string]interface{} `json:"traffic_data"`
	UserID         string                 `json:"user_id,omitempty"`
	SessionID      string                 `json:"session_id,omitempty"`
	IPAddress      string                 `json:"ip_address,omitempty"`
	TimeWindow     time.Duration          `json:"time_window,omitempty"`
	AnalysisType   string                 `json:"analysis_type"`
	IncludeBaseline bool                  `json:"include_baseline"`
	Configuration  map[string]interface{} `json:"configuration,omitempty"`
	RequestID      string                 `json:"request_id"`
	Timestamp      time.Time              `json:"timestamp"`
}

type BaselineCreationRequest struct {
	EntityID      string                   `json:"entity_id"`
	EntityType    string                   `json:"entity_type"`
	TrainingData  []map[string]interface{} `json:"training_data,omitempty"`
	Configuration map[string]interface{}   `json:"configuration,omitempty"`
	Description   string                   `json:"description,omitempty"`
	Tags          []string                 `json:"tags,omitempty"`
}

type BehaviorAnalysisResult struct {
	RequestID          string                 `json:"request_id"`
	PatternsDetected   []BehaviorPattern      `json:"patterns_detected"`
	AnomaliesDetected  []BehaviorAnomaly      `json:"anomalies_detected"`
	RiskAssessment     RiskAssessment         `json:"risk_assessment"`
	BaselineComparison BaselineComparison     `json:"baseline_comparison"`
	Recommendations    []string               `json:"recommendations"`
	ProcessingTime     time.Duration          `json:"processing_time"`
	Metadata           map[string]interface{} `json:"metadata"`
	AnalyzedAt         time.Time              `json:"analyzed_at"`
}

type BehaviorAnomaly struct {
	Type           string                 `json:"type"`
	Description    string                 `json:"description"`
	Severity       string                 `json:"severity"`
	Score          float64                `json:"score"`
	Pattern        BehaviorPattern        `json:"pattern"`
	Deviation      map[string]interface{} `json:"deviation"`
	Context        map[string]interface{} `json:"context"`
	Recommendations []string              `json:"recommendations"`
}

type RiskAssessment struct {
	OverallRiskScore   float64            `json:"overall_risk_score"`
	RiskLevel          string             `json:"risk_level"`
	RiskFactors        []RiskFactor       `json:"risk_factors"`
	MitigationActions  []string           `json:"mitigation_actions"`
	ConfidenceLevel    float64            `json:"confidence_level"`
	AssessmentBasis    string             `json:"assessment_basis"`
}

type RiskFactor struct {
	Factor      string  `json:"factor"`
	Score       float64 `json:"score"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
	Evidence    []string `json:"evidence"`
}

type BaselineComparison struct {
	HasBaseline        bool                   `json:"has_baseline"`
	BaselineAge        time.Duration          `json:"baseline_age,omitempty"`
	DeviationScore     float64                `json:"deviation_score"`
	SignificantChanges []BaselineDeviation    `json:"significant_changes"`
	Stability          float64                `json:"stability"`
	Confidence         float64                `json:"confidence"`
}

type BaselineDeviation struct {
	Metric        string      `json:"metric"`
	BaselineValue interface{} `json:"baseline_value"`
	CurrentValue  interface{} `json:"current_value"`
	DeviationPct  float64     `json:"deviation_percentage"`
	Significance  string      `json:"significance"`
	Description   string      `json:"description"`
}

// Behavior pattern types
const (
	BehaviorTypeUser        = "user"
	BehaviorTypeSession     = "session"
	BehaviorTypeIP          = "ip_address"
	BehaviorTypeAPI         = "api"
	BehaviorTypeEndpoint    = "endpoint"
	BehaviorTypeApplication = "application"
	BehaviorTypeSystem      = "system"
)

// Behavior categories
const (
	BehaviorCategoryAccess    = "access"
	BehaviorCategoryUsage     = "usage"
	BehaviorCategoryTiming    = "timing"
	BehaviorCategoryVolume    = "volume"
	BehaviorCategorySequence  = "sequence"
	BehaviorCategoryLocation  = "location"
	BehaviorCategoryDevice    = "device"
	BehaviorCategoryContent   = "content"
)

// Risk levels
const (
	RiskLevelCritical = "critical"
	RiskLevelHigh     = "high"
	RiskLevelMedium   = "medium"
	RiskLevelLow      = "low"
	RiskLevelMinimal  = "minimal"
)

type BaselineProfile struct {
	AccessPatterns   *AccessPatterns   `json:"access_patterns,omitempty"`
	UsagePatterns    *UsagePatterns    `json:"usage_patterns,omitempty"`
	TimingPatterns   *TimingPatterns   `json:"timing_patterns,omitempty"`
	LocationPatterns *LocationPatterns `json:"location_patterns,omitempty"`
	// Add more fields as needed
}

type AccessPatterns struct {
	NormalAccessHours      []int   `json:"normal_access_hours"`
	AverageHourlyAccess    float64 `json:"average_hourly_access"`
}

type UsagePatterns struct {
	CommonEndpoints map[string]float64 `json:"common_endpoints"`
	MethodFrequency map[string]float64 `json:"method_frequency"`
}

type TimingPatterns struct {}
type LocationPatterns struct {}

type BehaviorPatternFilter struct {
	EntityID    string `json:"entity_id,omitempty"`
	EntityType  string `json:"entity_type,omitempty"`
	PatternType string `json:"pattern_type,omitempty"`
	Severity    string `json:"severity,omitempty"`
	Status      string `json:"status,omitempty"`
}

type BehaviorPatternUpdate struct {
	Status  string `json:"status,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type BehaviorChange struct {
	ChangeType string    `json:"change_type"`
	EntityID   string    `json:"entity_id"`
	Timestamp  time.Time `json:"timestamp"`
	Details    string    `json:"details"`
}
