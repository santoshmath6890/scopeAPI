package models

import (
	"time"
)

type Threat struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Severity      string    `json:"severity"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	DetectionMethod string  `json:"detection_method"`
	Confidence    float64   `json:"confidence"`
	RiskScore     float64   `json:"risk_score"`
	Indicators    []ThreatIndicator `json:"indicators"`
	RequestData   map[string]interface{} `json:"request_data"`
	FirstSeen     time.Time `json:"first_seen"`
	LastSeen      time.Time `json:"last_seen"`
	Count         int       `json:"count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	APIID         string    `json:"api_id"`
	EndpointID    string    `json:"endpoint_id"`
	SourceIP      string    `json:"source_ip"`
	AttackType    string    `json:"attack_type"`
	RequestDetail string    `json:"request_detail"`
	ResponseDetail string   `json:"response_detail"`
	ResponseData  map[string]interface{} `json:"response_data"`
	Timestamp     time.Time `json:"timestamp"`
}

type ThreatIndicator struct {
	Type        string      `json:"type"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
	Severity    string      `json:"severity"`
	Confidence  float64     `json:"confidence"`
	Context     interface{} `json:"context,omitempty"`
}

type ThreatAnalysisRequest struct {
	TrafficData   map[string]interface{} `json:"traffic_data"`
	RequestID     string                 `json:"request_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Source        string                 `json:"source"`
	AnalysisType  string                 `json:"analysis_type"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

type ThreatAnalysisResult struct {
	RequestID       string                 `json:"request_id"`
	ThreatDetected  bool                   `json:"threat_detected"`
	ThreatType      string                 `json:"threat_type,omitempty"`
	Severity        string                 `json:"severity,omitempty"`
	Confidence      float64                `json:"confidence"`
	RiskScore       float64                `json:"risk_score"`
	Indicators      []ThreatIndicator      `json:"indicators"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	AnalyzedAt      time.Time              `json:"analyzed_at"`
}

type ThreatFilter struct {
	AttackType string    `json:"attack_type,omitempty"`
	Status     string    `json:"status,omitempty"`
	Since      time.Time `json:"since,omitempty"`
}

type ThreatStatistics struct {
	TotalThreats      int64              `json:"total_threats"`
	ActiveThreats     int64              `json:"active_threats"`
	ResolvedThreats   int64              `json:"resolved_threats"`
	CriticalThreats   int64              `json:"critical_threats"`
	HighThreats       int64              `json:"high_threats"`
	MediumThreats     int64              `json:"medium_threats"`
	LowThreats        int64              `json:"low_threats"`
	ThreatsByType     map[string]int64   `json:"threats_by_type"`
	ThreatsBySource   map[string]int64   `json:"threats_by_source"`
	RecentThreats     int64              `json:"recent_threats"`
	TrendData         []ThreatTrendPoint `json:"trend_data"`
	TopTargetedAPIs   []APIThreatSummary `json:"top_targeted_apis"`
	TopAttackerIPs    []IPThreatSummary  `json:"top_attacker_ips"`
}

type ThreatTrendPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	ThreatCount int64     `json:"threat_count"`
	Severity    string    `json:"severity,omitempty"`
	Type        string    `json:"type,omitempty"`
}

type APIThreatSummary struct {
	APIID       string `json:"api_id"`
	APIName     string `json:"api_name"`
	ThreatCount int64  `json:"threat_count"`
	LastThreat  time.Time `json:"last_threat"`
}

type IPThreatSummary struct {
	IPAddress   string    `json:"ip_address"`
	ThreatCount int64     `json:"threat_count"`
	LastThreat  time.Time `json:"last_threat"`
	Country     string    `json:"country,omitempty"`
	ISP         string    `json:"isp,omitempty"`
}

type ThreatUpdateRequest struct {
	Status     string                 `json:"status,omitempty"`
	Resolved   *bool                  `json:"resolved,omitempty"`
	ResolvedBy string                 `json:"resolved_by,omitempty"`
	Notes      string                 `json:"notes,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Threat severity levels
const (
	ThreatSeverityCritical = "critical"
	ThreatSeverityHigh     = "high"
	ThreatSeverityMedium   = "medium"
	ThreatSeverityLow      = "low"
	ThreatSeverityInfo     = "info"
)

// Threat types
const (
	ThreatTypeInjection        = "injection"
	ThreatTypeXSS              = "xss"
	ThreatTypeCSRF             = "csrf"
	ThreatTypeBruteForce       = "brute_force"
	ThreatTypeDDoS             = "ddos"
	ThreatTypeDataExfiltration = "data_exfiltration"
	ThreatTypePrivilegeEsc     = "privilege_escalation"
	ThreatTypeAnomaly          = "anomaly"
	ThreatTypeMalware          = "malware"
	ThreatTypePhishing         = "phishing"
	ThreatTypeRateLimitAbuse   = "rate_limit_abuse"
	ThreatTypeUnauthorized     = "unauthorized_access"
)

// Threat status
const (
	ThreatStatusNew        = "new"
	ThreatStatusInProgress = "in_progress"
	ThreatStatusResolved   = "resolved"
	ThreatStatusFalsePos   = "false_positive"
	ThreatStatusIgnored    = "ignored"
)

// Detection methods
const (
	DetectionMethodSignature  = "signature"
	DetectionMethodAnomaly    = "anomaly"
	DetectionMethodBehavioral = "behavioral"
	DetectionMethodML         = "machine_learning"
	DetectionMethodRule       = "rule_based"
	DetectionMethodHeuristic  = "heuristic"
)
