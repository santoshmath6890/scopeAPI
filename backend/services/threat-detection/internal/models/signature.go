package models

import "time"

type SignatureDetectionRequest struct {
	RequestID   string                 `json:"request_id"`
	EntityID    string                 `json:"entity_id"`
	EntityType  string                 `json:"entity_type"`
	Payload     map[string]interface{} `json:"payload"`
	Timestamp   time.Time              `json:"timestamp"`
	RequestData map[string]interface{} `json:"request_data"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	APIID       string                 `json:"api_id"`
	EndpointID  string                 `json:"endpoint_id"`
}

type SignatureTestRequest struct {
	SignatureID string                   `json:"signature_id"`
	TestData    []map[string]interface{} `json:"test_data"`
	Description string                   `json:"description,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
}

type SignatureDetectionResult struct {
	ResultID    string    `json:"result_id"`
	SignatureID string    `json:"signature_id"`
	Matched     bool      `json:"matched"`
	Details     string    `json:"details"`
	DetectedAt  time.Time `json:"detected_at"`
}

type ThreatSignature struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Pattern     string    `json:"pattern"`
	Severity    string    `json:"severity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	RiskScore   float64   `json:"risk_score"`
	Confidence  float64   `json:"confidence"`
	Tags        []string  `json:"tags"`
	SignatureSet string   `json:"signature_set"`
	Enabled     bool      `json:"enabled"`
	Rules       []SignatureRule `json:"rules"`
}

type SignatureFilter struct {
	SignatureID  string    `json:"signature_id,omitempty"`
	Type         string    `json:"type,omitempty"`
	Category     string    `json:"category,omitempty"`
	Severity     string    `json:"severity,omitempty"`
	Pattern      string    `json:"pattern,omitempty"`
	SignatureSet string    `json:"signature_set,omitempty"`
	Enabled      bool      `json:"enabled,omitempty"`
}

type SignatureTestResult struct {
	TestID        string                `json:"test_id"`
	SignatureID   string                `json:"signature_id"`
	SignatureName string                `json:"signature_name"`
	Passed        bool                  `json:"passed"`
	Details       string                `json:"details"`
	TestCases     []SignatureTestCase   `json:"test_cases"`
	TotalTests    int                   `json:"total_tests"`
	PassedTests   int                   `json:"passed_tests"`
	FailedTests   int                   `json:"failed_tests"`
	TestedAt      time.Time             `json:"tested_at"`
}

type SignatureTestCase struct {
	TestID       string                 `json:"test_id"`
	TestName     string                 `json:"test_name"`
	TestData     map[string]interface{} `json:"test_data"`
	Expected     bool                   `json:"expected"`
	Actual       bool                   `json:"actual"`
	Passed       bool                   `json:"passed"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	MatchedRule  string                 `json:"matched_rule,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

type SignatureMatch struct {
	SignatureID   string    `json:"signature_id"`
	SignatureName string    `json:"signature_name"`
	SignatureType string    `json:"signature_type"`
	Category      string    `json:"category"`
	Severity      string    `json:"severity"`
	RiskScore     float64   `json:"risk_score"`
	Confidence    float64   `json:"confidence"`
	Description   string    `json:"description"`
	MatchedField  string    `json:"matched_field"`
	MatchedValue  string    `json:"matched_value"`
	RuleMatched   string    `json:"rule_matched"`
	RuleOperator  string    `json:"rule_operator"`
	RuleValue     string    `json:"rule_value"`
	MatchedAt     time.Time `json:"matched_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	Matched       bool      `json:"matched"`
	Details       string    `json:"details"`
}

type SignatureRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Field       string `json:"field"`
	Operator    string `json:"operator"`
	Value       string `json:"value"`
	IntValue    int    `json:"int_value"`
	Weight      float64 `json:"weight"`
}

type SignatureStatistics struct {
	TotalSignatures       int                          `json:"total_signatures"`
	Matched               int                          `json:"matched"`
	Unmatched             int                          `json:"unmatched"`
	EnabledSignatures     int                          `json:"enabled_signatures"`
	DisabledSignatures    int                          `json:"disabled_signatures"`
	SignaturesByType      map[string]int               `json:"signatures_by_type"`
	SignaturesByCategory  map[string]int               `json:"signatures_by_category"`
	MatchStatistics       *SignatureMatchStats         `json:"match_statistics"`
	GeneratedAt           time.Time                    `json:"generated_at"`
}

type SignatureMatchStats struct {
	TotalMatches     int            `json:"total_matches"`
	MatchesByType    map[string]int `json:"matches_by_type"`
	MatchesByCategory map[string]int `json:"matches_by_category"`
	MatchesBySeverity map[string]int `json:"matches_by_severity"`
}

type SignatureOptimizationResult struct {
	OptimizedSignatures int       `json:"optimized_signatures"`
	Details             string    `json:"details"`
	RemovedSignatures   []string  `json:"removed_signatures"`
	UpdatedSignatures   []string  `json:"updated_signatures"`
	Recommendations     []string  `json:"recommendations"`
	OptimizedAt         time.Time `json:"optimized_at"`
}

// Signature type constants
const (
	SignatureTypeCustom    = "custom"
	SignatureTypeBuiltIn   = "builtin"
	SignatureTypeCommunity = "community"
)