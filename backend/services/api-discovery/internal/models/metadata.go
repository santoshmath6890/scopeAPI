package models

import "time"

type Metadata struct {
	ID                  string                 `json:"id" db:"id"`
	EndpointID          string                 `json:"endpoint_id" db:"endpoint_id"`
	APIID               string                 `json:"api_id" db:"api_id"`
	URL                 string                 `json:"url" db:"url"`
	Method              string                 `json:"method" db:"method"`
	Title               string                 `json:"title" db:"title"`
	Description         string                 `json:"description" db:"description"`
	Tags                []string               `json:"tags" db:"tags"`
	Category            string                 `json:"category" db:"category"`
	BusinessOwner       string                 `json:"business_owner" db:"business_owner"`
	TechnicalOwner      string                 `json:"technical_owner" db:"technical_owner"`
	DataSensitivity     string                 `json:"data_sensitivity" db:"data_sensitivity"`
	ComplianceReqs      []string               `json:"compliance_requirements" db:"compliance_requirements"`
	Parameters          []Parameter            `json:"parameters" db:"parameters"`
	ResponseSchema      map[string]interface{} `json:"response_schema" db:"response_schema"`
	RequestSchema       map[string]interface{} `json:"request_schema" db:"request_schema"`
	Examples            []MetadataExample      `json:"examples" db:"examples"`
	Documentation       *Documentation         `json:"documentation" db:"documentation"`
	Versioning          map[string]interface{} `json:"versioning" db:"versioning"`
	Performance         *PerformanceMetrics    `json:"performance" db:"performance"`
	Security            *SecurityMetadata      `json:"security" db:"security"`
	Quality             *QualityMetrics        `json:"quality" db:"quality"`
	Usage               *UsageMetrics          `json:"usage" db:"usage"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

type SecurityMetadata struct {
	HasHTTPS              bool              `json:"has_https"`
	HasSecurityHeaders    bool              `json:"has_security_headers"`
	SecurityHeaders       map[string]string `json:"security_headers"`
	AuthenticationRequired bool             `json:"authentication_required"`
	AuthenticationMethods []string          `json:"authentication_methods"`
	DataClassification    string            `json:"data_classification"`
	EncryptionInTransit   bool              `json:"encryption_in_transit"`
	EncryptionAtRest      bool              `json:"encryption_at_rest"`
	VulnerabilityScans    []VulnerabilityScan `json:"vulnerability_scans"`
	ComplianceStatus      map[string]string `json:"compliance_status"`
	LastSecurityScan      *time.Time        `json:"last_security_scan"`
}

type QualityMetrics struct {
	DocumentationScore float64   `json:"documentation_score"`
	APIDesignScore    float64   `json:"api_design_score"`
	ConsistencyScore  float64   `json:"consistency_score"`
	QualityIssues     []string  `json:"quality_issues"`
	LastQualityCheck  time.Time `json:"last_quality_check"`
}

type PerformanceMetrics struct {
	LastMeasured     time.Time `json:"last_measured"`
	SecurityOverhead float64   `json:"security_overhead"`
}

type UsageMetrics struct {
	TotalRequests    int64     `json:"total_requests"`
	RequestsLast24h  int64     `json:"requests_last_24h"`
	RequestsLastWeek int64     `json:"requests_last_week"`
	RequestsLastMonth int64    `json:"requests_last_month"`
	LastUsed         time.Time `json:"last_used"`
}

type VulnerabilityScan struct {
	Name      string    `json:"name"`
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

type MetadataExample struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Example     interface{}            `json:"example"`
	Request     map[string]interface{} `json:"request"`
	Response    map[string]interface{} `json:"response"`
	StatusCode  int                    `json:"status_code"`
}

type Documentation struct {
	Summary      string                 `json:"summary"`
	ExternalDocs map[string]interface{} `json:"external_docs"`
} 