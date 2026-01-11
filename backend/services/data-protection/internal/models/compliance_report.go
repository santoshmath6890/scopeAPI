package models

import (
	"time"
)

type ComplianceFramework string

const (
	ComplianceFrameworkGDPR     ComplianceFramework = "gdpr"
	ComplianceFrameworkCCPA     ComplianceFramework = "ccpa"
	ComplianceFrameworkHIPAA    ComplianceFramework = "hipaa"
	ComplianceFrameworkSOX      ComplianceFramework = "sox"
	ComplianceFrameworkPCIDSS   ComplianceFramework = "pci_dss"
	ComplianceFrameworkISO27001 ComplianceFramework = "iso_27001"
	ComplianceFrameworkNIST     ComplianceFramework = "nist"
	ComplianceFrameworkCustom   ComplianceFramework = "custom"
)

type ComplianceStatus string

const (
	ComplianceStatusCompliant          ComplianceStatus = "compliant"
	ComplianceStatusPartiallyCompliant ComplianceStatus = "partially_compliant"
	ComplianceStatusNonCompliant       ComplianceStatus = "non_compliant"
	ComplianceStatusUnknown            ComplianceStatus = "unknown"
	ComplianceStatusInProgress         ComplianceStatus = "in_progress"
	ComplianceStatusViolation          ComplianceStatus = "violation"
	ComplianceStatusWarning            ComplianceStatus = "warning"
)

type ComplianceCategory string

const (
	ComplianceCategoryPrivacy            ComplianceCategory = "privacy"
	ComplianceCategorySecurity           ComplianceCategory = "security"
	ComplianceCategoryPHI                ComplianceCategory = "phi"
	ComplianceCategoryCardholderData     ComplianceCategory = "cardholder_data"
	ComplianceCategoryFinancialReporting ComplianceCategory = "financial_reporting"
	ComplianceCategoryLegal              ComplianceCategory = "legal"
	ComplianceCategoryOther              ComplianceCategory = "other"
	ComplianceCategoryConsent            ComplianceCategory = "consent"
	ComplianceCategoryDataMinimization   ComplianceCategory = "data_minimization"
	ComplianceCategoryPrivacyNotice      ComplianceCategory = "privacy_notice"
	ComplianceCategoryAuditTrail         ComplianceCategory = "audit_trail"
)

type ComplianceSeverity string

const (
	ComplianceSeverityLow      ComplianceSeverity = "low"
	ComplianceSeverityMedium   ComplianceSeverity = "medium"
	ComplianceSeverityHigh     ComplianceSeverity = "high"
	ComplianceSeverityCritical ComplianceSeverity = "critical"
)

type ViolationStatus string

const (
	ViolationStatusOpen          ViolationStatus = "open"
	ViolationStatusInProgress    ViolationStatus = "in_progress"
	ViolationStatusResolved      ViolationStatus = "resolved"
	ViolationStatusIgnored       ViolationStatus = "ignored"
	ViolationStatusFalsePositive ViolationStatus = "false_positive"
)

type ComplianceReport struct {
	ID               string                     `json:"id" db:"id"`
	Name             string                     `json:"name" db:"name"`
	Description      string                     `json:"description" db:"description"`
	Type             string                     `json:"type" db:"type"`
	FrameworkID      string                     `json:"framework_id" db:"framework_id"`
	Framework        ComplianceFramework        `json:"framework" db:"framework"`
	Status           ComplianceStatus           `json:"status" db:"status"`
	Summary          ComplianceReportSummary    `json:"summary" db:"summary"`
	FrameworkReports map[string]FrameworkReport `json:"framework_reports" db:"framework_reports"`
	Violations       []ComplianceViolation      `json:"violations" db:"violations"`
	Recommendations  []ComplianceRecommendation `json:"recommendations" db:"recommendations"`
	Trends           []ComplianceTrend          `json:"trends" db:"trends"`
	Filter           *ComplianceReportFilter    `json:"filter,omitempty" db:"filter"`
	Metadata         map[string]interface{}     `json:"metadata" db:"metadata"`
	GeneratedAt      time.Time                  `json:"generated_at" db:"generated_at"`
	GeneratedBy      string                     `json:"generated_by" db:"generated_by"`
	ValidFrom        time.Time                  `json:"valid_from" db:"valid_from"`
	ValidTo          time.Time                  `json:"valid_to" db:"valid_to"`
	CreatedAt        time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at" db:"updated_at"`
}

type ComplianceReportSummary struct {
	TotalValidations      int                         `json:"total_validations"`
	ComplianceRate        float64                     `json:"compliance_rate"`
	OverallScore          float64                     `json:"overall_score"`
	TotalViolations       int                         `json:"total_violations"`
	OpenViolations        int                         `json:"open_violations"`
	ResolvedViolations    int                         `json:"resolved_violations"`
	ViolationsByFramework map[string]int              `json:"violations_by_framework"`
	ViolationsBySeverity  map[ComplianceSeverity]int  `json:"violations_by_severity"`
	ViolationsByStatus    map[ViolationStatus]int     `json:"violations_by_status"`
	TopViolations         []TopViolation              `json:"top_violations"`
	ComplianceByAPI       map[string]ComplianceMetric `json:"compliance_by_api"`
	ComplianceByEndpoint  map[string]ComplianceMetric `json:"compliance_by_endpoint"`
	RiskScore             float64                     `json:"risk_score"`
	TrendDirection        string                      `json:"trend_direction"`
}

type FrameworkReport struct {
	FrameworkID       string                       `json:"framework_id"`
	FrameworkName     string                       `json:"framework_name"`
	Status            ComplianceStatus             `json:"status"`
	Score             float64                      `json:"score"`
	TotalValidations  int                          `json:"total_validations"`
	CompliantCount    int                          `json:"compliant_count"`
	ViolationCount    int                          `json:"violation_count"`
	ComplianceRate    float64                      `json:"compliance_rate"`
	AverageScore      float64                      `json:"average_score"`
	RequirementStatus map[string]RequirementStatus `json:"requirement_status"`
	TopViolations     []TopViolation               `json:"top_violations"`
	Recommendations   []string                     `json:"recommendations"`
	LastAssessment    time.Time                    `json:"last_assessment"`
}

type RequirementStatus struct {
	RequirementID   string           `json:"requirement_id"`
	RequirementName string           `json:"requirement_name"`
	Status          ComplianceStatus `json:"status"`
	Score           float64          `json:"score"`
	ViolationCount  int              `json:"violation_count"`
	LastChecked     time.Time        `json:"last_checked"`
}

type ComplianceViolation struct {
	ID          string                 `json:"id" db:"id"`
	RuleID      string                 `json:"rule_id" db:"rule_id"`
	RuleName    string                 `json:"rule_name" db:"rule_name"`
	Framework   string                 `json:"framework" db:"framework"`
	Category    string                 `json:"category" db:"category"`
	Severity    ComplianceSeverity     `json:"severity" db:"severity"`
	Status      ViolationStatus        `json:"status" db:"status"`
	Message     string                 `json:"message" db:"message"`
	Description string                 `json:"description" db:"description"`
	Remediation string                 `json:"remediation" db:"remediation"`
	APIID       string                 `json:"api_id" db:"api_id"`
	EndpointID  string                 `json:"endpoint_id" db:"endpoint_id"`
	RequestID   string                 `json:"request_id" db:"request_id"`
	Evidence    ViolationEvidence      `json:"evidence" db:"evidence"`
	Impact      ViolationImpact        `json:"impact" db:"impact"`
	Resolution  *ViolationResolution   `json:"resolution,omitempty" db:"resolution"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	DetectedAt  time.Time              `json:"detected_at" db:"detected_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	AssignedTo  string                 `json:"assigned_to" db:"assigned_to"`
	DueDate     *time.Time             `json:"due_date,omitempty" db:"due_date"`
}

type ViolationEvidence struct {
	RequestData   string                 `json:"request_data"`
	ResponseData  string                 `json:"response_data"`
	Headers       map[string]string      `json:"headers"`
	QueryParams   map[string]string      `json:"query_params"`
	BodyContent   string                 `json:"body_content"`
	Metadata      map[string]interface{} `json:"metadata"`
	Screenshots   []string               `json:"screenshots"`
	LogEntries    []string               `json:"log_entries"`
	NetworkTraces []string               `json:"network_traces"`
}

type ViolationImpact struct {
	RiskScore       float64  `json:"risk_score"`
	BusinessImpact  string   `json:"business_impact"`
	TechnicalImpact string   `json:"technical_impact"`
	LegalImpact     string   `json:"legal_impact"`
	FinancialImpact string   `json:"financial_impact"`
	AffectedSystems []string `json:"affected_systems"`
	AffectedUsers   int      `json:"affected_users"`
	DataVolume      int64    `json:"data_volume"`
	Urgency         string   `json:"urgency"`
}

type ViolationResolution struct {
	ResolutionType   string             `json:"resolution_type"`
	Description      string             `json:"description"`
	Actions          []ResolutionAction `json:"actions"`
	ResolvedBy       string             `json:"resolved_by"`
	ResolvedAt       time.Time          `json:"resolved_at"`
	VerifiedBy       string             `json:"verified_by"`
	VerifiedAt       *time.Time         `json:"verified_at,omitempty"`
	Notes            string             `json:"notes"`
	Attachments      []string           `json:"attachments"`
	FollowUpRequired bool               `json:"follow_up_required"`
	FollowUpDate     *time.Time         `json:"follow_up_date,omitempty"`
}

type ResolutionAction struct {
	Action      string                 `json:"action"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	AssignedTo  string                 `json:"assigned_to"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Notes       string                 `json:"notes"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type TopViolation struct {
	RuleID      string             `json:"rule_id"`
	RuleName    string             `json:"rule_name"`
	Count       int                `json:"count"`
	Severity    ComplianceSeverity `json:"severity"`
	Category    string             `json:"category"`
	Framework   string             `json:"framework"`
	Percentage  float64            `json:"percentage"`
	TrendChange float64            `json:"trend_change"`
}

type ComplianceRecommendation struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Priority     string                 `json:"priority"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Framework    string                 `json:"framework"`
	Category     string                 `json:"category"`
	Actions      []RecommendationAction `json:"actions"`
	Impact       string                 `json:"impact"`
	Effort       string                 `json:"effort"`
	Timeline     string                 `json:"timeline"`
	Resources    []string               `json:"resources"`
	Dependencies []string               `json:"dependencies"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
}

type RecommendationAction struct {
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Effort      string    `json:"effort"`
	Timeline    string    `json:"timeline"`
	Owner       string    `json:"owner"`
	DueDate     time.Time `json:"due_date"`
}

type ComplianceTrend struct {
	Date             time.Time        `json:"date"`
	Framework        string           `json:"framework"`
	Status           ComplianceStatus `json:"status"`
	Score            float64          `json:"score"`
	TotalValidations int              `json:"total_validations"`
	CompliantCount   int              `json:"compliant_count"`
	ViolationCount   int              `json:"violation_count"`
	ComplianceRate   float64          `json:"compliance_rate"`
	Change           float64          `json:"change"`
	ChangePercent    float64          `json:"change_percent"`
}

type ComplianceMetric struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Status         ComplianceStatus `json:"status"`
	Score          float64          `json:"score"`
	ViolationCount int              `json:"violation_count"`
	ComplianceRate float64          `json:"compliance_rate"`
	LastAssessment time.Time        `json:"last_assessment"`
	TrendDirection string           `json:"trend_direction"`
	RiskLevel      string           `json:"risk_level"`
}

type ComplianceReportFilter struct {
	Frameworks    []string   `json:"frameworks,omitempty"`
	APIIDs        []string   `json:"api_ids,omitempty"`
	EndpointIDs   []string   `json:"endpoint_ids,omitempty"`
	Severities    []string   `json:"severities,omitempty"`
	Statuses      []string   `json:"statuses,omitempty"`
	FrameworkID   string     `json:"framework_id,omitempty"`
	Status        string     `json:"status,omitempty"`
	Since         *time.Time `json:"since,omitempty"`
	Until         *time.Time `json:"until,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Categories    []string   `json:"categories,omitempty"`
	RuleIDs       []string   `json:"rule_ids,omitempty"`
	IncludeTrends bool       `json:"include_trends"`
	TrendPeriod   string     `json:"trend_period"`
	Limit         int        `json:"limit,omitempty"`
	Offset        int        `json:"offset,omitempty"`
}

type ComplianceValidation struct {
	ID               string                               `json:"id" db:"id"`
	RequestID        string                               `json:"request_id" db:"request_id"`
	APIID            string                               `json:"api_id" db:"api_id"`
	EndpointID       string                               `json:"endpoint_id" db:"endpoint_id"`
	Frameworks       []string                             `json:"frameworks" db:"frameworks"`
	OverallStatus    ComplianceStatus                     `json:"overall_status" db:"overall_status"`
	OverallScore     float64                              `json:"overall_score" db:"overall_score"`
	ViolationCount   int                                  `json:"violation_count" db:"violation_count"`
	WarningCount     int                                  `json:"warning_count" db:"warning_count"`
	FrameworkResults map[string]FrameworkComplianceResult `json:"framework_results" db:"framework_results"`
	Violations       []ComplianceViolation                `json:"violations" db:"violations"`
	Warnings         []ComplianceWarning                  `json:"warnings" db:"warnings"`
	ProcessingTime   time.Duration                        `json:"processing_time" db:"processing_time"`
	Metadata         map[string]interface{}               `json:"metadata" db:"metadata"`
	IPAddress        string                               `json:"ip_address" db:"ip_address"`
	UserAgent        string                               `json:"user_agent" db:"user_agent"`
	ValidatedAt      time.Time                            `json:"validated_at" db:"validated_at"`
	CreatedAt        time.Time                            `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                            `json:"updated_at" db:"updated_at"`
}

type FrameworkValidationResult struct {
	Framework      string                `json:"framework"`
	Status         ComplianceStatus      `json:"status"`
	Score          float64               `json:"score"`
	ViolationCount int                   `json:"violation_count"`
	WarningCount   int                   `json:"warning_count"`
	Violations     []ComplianceViolation `json:"violations"`
	Warnings       []ComplianceWarning   `json:"warnings"`
	Requirements   []RequirementResult   `json:"requirements"`
	ProcessingTime time.Duration         `json:"processing_time"`
}

type RequirementResult struct {
	RequirementID   string           `json:"requirement_id"`
	RequirementName string           `json:"requirement_name"`
	Status          ComplianceStatus `json:"status"`
	Score           float64          `json:"score"`
	Message         string           `json:"message"`
	Evidence        string           `json:"evidence"`
}

type ComplianceWarning struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"rule_id"`
	RuleName    string                 `json:"rule_name"`
	Framework   string                 `json:"framework"`
	Category    string                 `json:"category"`
	Message     string                 `json:"message"`
	Description string                 `json:"description"`
	Suggestion  string                 `json:"suggestion"`
	APIID       string                 `json:"api_id"`
	EndpointID  string                 `json:"endpoint_id"`
	RequestID   string                 `json:"request_id"`
	Metadata    map[string]interface{} `json:"metadata"`
	DetectedAt  time.Time              `json:"detected_at"`
}

type ComplianceValidationRequest struct {
	RequestID       string                 `json:"request_id"`
	APIID           string                 `json:"api_id"`
	EndpointID      string                 `json:"endpoint_id"`
	Frameworks      []string               `json:"frameworks"`
	Content         string                 `json:"content"`
	ContentType     string                 `json:"content_type"`
	Source          string                 `json:"source"`
	Rules           []string               `json:"rules"`
	Options         ValidationOptions      `json:"options"`
	Context         map[string]interface{} `json:"context"`
	IPAddress       string                 `json:"ip_address"`
	UserAgent       string                 `json:"user_agent"`
	DataFactors     map[string]interface{} `json:"data_factors"`
	SecurityFactors map[string]interface{} `json:"security_factors"`
	ContextFactors  map[string]interface{} `json:"context_factors"`
}

type ValidationOptions struct {
	EnableRealTimeValidation bool     `json:"enable_real_time_validation"`
	EnableDeepValidation     bool     `json:"enable_deep_validation"`
	IncludeWarnings          bool     `json:"include_warnings"`
	IncludeRecommendations   bool     `json:"include_recommendations"`
	CustomRules              []string `json:"custom_rules"`
	ValidationDepth          int      `json:"validation_depth"`
	TimeoutSeconds           int      `json:"timeout_seconds"`
}

type ComplianceValidationResult struct {
	RequestID             string                               `json:"request_id"`
	Regulations           []string                             `json:"regulations"`
	OverallStatus         ComplianceStatus                     `json:"overall_status"`
	OverallScore          float64                              `json:"overall_score"`
	ViolationCount        int                                  `json:"violation_count"`
	WarningCount          int                                  `json:"warning_count"`
	FrameworkResults      map[string]FrameworkComplianceResult `json:"framework_results"`
	Violations            []ComplianceViolation                `json:"violations"`
	Warnings              []ComplianceWarning                  `json:"warnings"`
	Recommendations       []ComplianceRecommendation           `json:"recommendations"`
	RecommendationStrings []string                             `json:"recommendation_strings"`
	Issues                []ComplianceIssue                    `json:"issues"`
	ProcessingTime        time.Duration                        `json:"processing_time"`
	Metadata              map[string]interface{}               `json:"metadata"`
	ValidatedAt           time.Time                            `json:"validated_at"`
}

type ComplianceRule struct {
	ID           string                 `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name"`
	Description  string                 `json:"description" db:"description"`
	Framework    string                 `json:"framework" db:"framework"`
	Category     string                 `json:"category" db:"category"`
	Severity     ComplianceSeverity     `json:"severity" db:"severity"`
	Type         string                 `json:"type" db:"type"`
	Enabled      bool                   `json:"enabled" db:"enabled"`
	Priority     int                    `json:"priority" db:"priority"`
	Conditions   []ComplianceCondition  `json:"conditions" db:"conditions"`
	Actions      []ComplianceAction     `json:"actions" db:"actions"`
	Requirements []string               `json:"requirements" db:"requirements"`
	Tags         []string               `json:"tags" db:"tags"`
	References   []string               `json:"references" db:"references"`
	Remediation  string                 `json:"remediation" db:"remediation"`
	Examples     []RuleExample          `json:"examples" db:"examples"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy    string                 `json:"created_by" db:"created_by"`
	UpdatedBy    string                 `json:"updated_by" db:"updated_by"`
	Version      string                 `json:"version" db:"version"`
	LastTested   *time.Time             `json:"last_tested,omitempty" db:"last_tested"`
}

type RuleCondition struct {
	Field         string                 `json:"field"`
	Operator      string                 `json:"operator"`
	Value         string                 `json:"value"`
	CaseSensitive bool                   `json:"case_sensitive"`
	Weight        float64                `json:"weight"`
	Context       map[string]interface{} `json:"context"`
	Negated       bool                   `json:"negated"`
}

type RuleAction struct {
	Type       string                 `json:"type"`
	Config     map[string]interface{} `json:"config"`
	Enabled    bool                   `json:"enabled"`
	Priority   int                    `json:"priority"`
	Conditions []string               `json:"conditions"`
	Parameters map[string]string      `json:"parameters"`
}

type RuleExample struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Expected    string `json:"expected"`
	Violation   bool   `json:"violation"`
}

type ComplianceRuleFilter struct {
	Framework string   `json:"framework,omitempty"`
	Category  string   `json:"category,omitempty"`
	Severity  string   `json:"severity,omitempty"`
	Enabled   *bool    `json:"enabled,omitempty"`
	Name      string   `json:"name,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Offset    int      `json:"offset,omitempty"`
}

type ComplianceStatusFilter struct {
	Frameworks []string `json:"frameworks,omitempty"`
	APIIDs     []string `json:"api_ids,omitempty"`
	Statuses   []string `json:"statuses,omitempty"`
}

type ComplianceStatusInfo struct {
	OverallStatus     ComplianceStatus           `json:"overall_status"`
	FrameworkStatuses map[string]FrameworkStatus `json:"framework_statuses"`
	Summary           ComplianceStatusSummary    `json:"summary"`
	LastUpdated       time.Time                  `json:"last_updated"`
}

type FrameworkStatus struct {
	FrameworkID      string           `json:"framework_id"`
	FrameworkName    string           `json:"framework_name"`
	Status           ComplianceStatus `json:"status"`
	Score            float64          `json:"score"`
	ActiveViolations int              `json:"active_violations"`
	LastAssessment   time.Time        `json:"last_assessment"`
}

type ComplianceStatusSummary struct {
	TotalFrameworks     int `json:"total_frameworks"`
	CompliantFrameworks int `json:"compliant_frameworks"`
	ActiveViolations    int `json:"active_violations"`
	RecentViolations    int `json:"recent_violations"`
}

type ComplianceValidationFilter struct {
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Frameworks  []string   `json:"frameworks,omitempty"`
	APIIDs      []string   `json:"api_ids,omitempty"`
	EndpointIDs []string   `json:"endpoint_ids,omitempty"`
	Statuses    []string   `json:"statuses,omitempty"`
	Limit       int        `json:"limit,omitempty"`
	Offset      int        `json:"offset,omitempty"`
}

// ComplianceFrameworkData - Compliance framework definition
type ComplianceFrameworkData struct {
	ID           string                  `json:"id" db:"id"`
	Name         string                  `json:"name" db:"name"`
	Description  string                  `json:"description" db:"description"`
	Type         string                  `json:"type" db:"type"`
	Version      string                  `json:"version" db:"version"`
	Region       string                  `json:"region" db:"region"`
	Categories   []string                `json:"categories" db:"categories"`
	Requirements []ComplianceRequirement `json:"requirements" db:"requirements"`
	Enabled      bool                    `json:"enabled" db:"enabled"`
	CreatedAt    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at" db:"updated_at"`
}

// ComplianceFrameworkFilter - Filter for compliance frameworks
type ComplianceFrameworkFilter struct {
	Type       string   `json:"type,omitempty"`
	Region     string   `json:"region,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}

// =============================================================================
// ADDITIONAL COMPLIANCE AND RISK TYPES
// =============================================================================

// TimeRange - Time range for filtering
type TimeRange struct {
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Since     time.Time `json:"since"`
	Until     time.Time `json:"until"`
}

// AuditLogFilter - Filter for audit logs
type AuditLogFilter struct {
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Since        *time.Time `json:"since,omitempty"`
	Until        *time.Time `json:"until,omitempty"`
	Action       string     `json:"action,omitempty"`
	UserID       string     `json:"user_id,omitempty"`
	ResourceType string     `json:"resource_type,omitempty"`
	Resource     string     `json:"resource,omitempty"`
	Limit        int        `json:"limit,omitempty"`
	Offset       int        `json:"offset,omitempty"`
}

// AuditLogEntry - Audit log entry
type AuditLogEntry struct {
	ID        string                 `json:"id" db:"id"`
	Timestamp time.Time              `json:"timestamp" db:"timestamp"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Action    string                 `json:"action" db:"action"`
	Resource  string                 `json:"resource" db:"resource"`
	Details   map[string]interface{} `json:"details" db:"details"`
	IPAddress string                 `json:"ip_address" db:"ip_address"`
	UserAgent string                 `json:"user_agent" db:"user_agent"`
	Success   bool                   `json:"success" db:"success"`
	Error     string                 `json:"error,omitempty" db:"error"`
}

// ComplianceStatistics - Statistics for compliance
type ComplianceStatistics struct {
	TotalFrameworks     int                    `json:"total_frameworks"`
	CompliantFrameworks int                    `json:"compliant_frameworks"`
	ViolationCount      int                    `json:"violation_count"`
	ComplianceRate      float64                `json:"compliance_rate"`
	StatisticsByDate    map[string]interface{} `json:"statistics_by_date"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// RiskScoringRequest - Risk scoring request
type RiskScoringRequest struct {
	RequestID       string                 `json:"request_id"`
	APIID           string                 `json:"api_id"`
	EndpointID      string                 `json:"endpoint_id"`
	Data            map[string]interface{} `json:"data"`
	Context         map[string]interface{} `json:"context"`
	DataFactors     map[string]interface{} `json:"data_factors"`
	SecurityFactors map[string]interface{} `json:"security_factors"`
	ContextFactors  map[string]interface{} `json:"context_factors"`
	IPAddress       string                 `json:"ip_address"`
	UserAgent       string                 `json:"user_agent"`
}

// RiskScoringResult - Risk scoring result
type RiskScoringResult struct {
	RequestID       string                 `json:"request_id"`
	RiskScore       float64                `json:"risk_score"`
	RiskLevel       RiskLevel              `json:"risk_level"`
	ProfileUsed     string                 `json:"profile_used"`
	ScoreBreakdown  RiskScoreBreakdown     `json:"score_breakdown"`
	AppliedRules    []AppliedRiskRule      `json:"applied_rules"`
	Recommendations []string               `json:"recommendations"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	Metadata        map[string]interface{} `json:"metadata"`
	CalculatedAt    time.Time              `json:"calculated_at"`
}

type RiskAssessmentRequest struct {
	RequestID       string                 `json:"request_id"`
	APIID           string                 `json:"api_id"`
	EndpointID      string                 `json:"endpoint_id"`
	DataFactors     map[string]interface{} `json:"data_factors"`
	SecurityFactors map[string]interface{} `json:"security_factors"`
	ContextFactors  map[string]interface{} `json:"context_factors"`
	IPAddress       string                 `json:"ip_address"`
	UserAgent       string                 `json:"user_agent"`
	DataSource      string                 `json:"data_source"`
}

type RiskScoreFilter struct {
	DataSource string     `json:"data_source"`
	RiskLevel  string     `json:"risk_level"`
	Since      *time.Time `json:"since"`
	Until      *time.Time `json:"until"`
}

type MitigationPlan struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	RiskID            string             `json:"risk_id"`
	MitigationActions []MitigationAction `json:"mitigation_actions"`
	Status            string             `json:"status"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type MitigationAction struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  string    `json:"assigned_to"`
	DueDate     time.Time `json:"due_date"`
}

type RiskScoreBreakdown struct {
	BaseScore            float64 `json:"base_score"`
	MultiplierAdjustment float64 `json:"multiplier_adjustment"`
	WeightedAdjustment   float64 `json:"weighted_adjustment"`
	RuleAdjustments      float64 `json:"rule_adjustments"`
}

type AppliedRiskRule struct {
	RuleID      string          `json:"rule_id"`
	RuleName    string          `json:"rule_name"`
	Category    string          `json:"category"`
	Adjustment  ScoreAdjustment `json:"adjustment"`
	ScoreChange float64         `json:"score_change"`
	Reason      string          `json:"reason"`
}

// RiskProfile - Risk assessment profile
type RiskProfile struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    RiskCategory       `json:"category"`
	BaseScore   float64            `json:"base_score"`
	Multipliers map[string]float64 `json:"multipliers"`
	Thresholds  RiskThresholds     `json:"thresholds"`
	Enabled     bool               `json:"enabled"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type RiskCategory string

const (
	RiskCategoryLow      RiskCategory = "low"
	RiskCategoryMedium   RiskCategory = "medium"
	RiskCategoryHigh     RiskCategory = "high"
	RiskCategoryCritical RiskCategory = "critical"
)

type RiskThresholds struct {
	Low      float64 `json:"low"`
	Medium   float64 `json:"medium"`
	High     float64 `json:"high"`
	Critical float64 `json:"critical"`
}

// RiskWeights - Risk calculation weights
type RiskWeights struct {
	DataSensitivity  float64 `json:"data_sensitivity"`
	ExposureLevel    float64 `json:"exposure_level"`
	AccessControls   float64 `json:"access_controls"`
	Vulnerabilities  float64 `json:"vulnerabilities"`
	ComplianceStatus float64 `json:"compliance_status"`
}

// RiskScoringRule - Risk scoring rule
type RiskScoringRule struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Category        string          `json:"category"`
	Priority        int             `json:"priority"`
	Weight          float64         `json:"weight"`
	Conditions      []RiskCondition `json:"conditions"`
	ScoreAdjustment ScoreAdjustment `json:"score_adjustment"`
	Enabled         bool            `json:"enabled"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// RiskCondition - Risk scoring condition
type RiskCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Weight   float64     `json:"weight"`
}

// ScoreAdjustment - Risk score adjustment
type ScoreAdjustment struct {
	Type   string  `json:"type"`
	Value  float64 `json:"value"`
	Reason string  `json:"reason"`
}

// RiskAssessment - Risk assessment result storage
type RiskAssessment struct {
	ID              string                 `json:"id" db:"id"`
	RequestID       string                 `json:"request_id" db:"request_id"`
	APIID           string                 `json:"api_id" db:"api_id"`
	EndpointID      string                 `json:"endpoint_id" db:"endpoint_id"`
	RiskScore       float64                `json:"risk_score" db:"risk_score"`
	RiskLevel       RiskLevel              `json:"risk_level" db:"risk_level"`
	ProfileID       string                 `json:"profile_id" db:"profile_id"`
	ScoreBreakdown  RiskScoreBreakdown     `json:"score_breakdown" db:"score_breakdown"`
	AppliedRules    []AppliedRiskRule      `json:"applied_rules" db:"applied_rules"`
	Recommendations []string               `json:"recommendations" db:"recommendations"`
	IPAddress       string                 `json:"ip_address" db:"ip_address"`
	UserAgent       string                 `json:"user_agent" db:"user_agent"`
	AssessedAt      time.Time              `json:"assessed_at" db:"assessed_at"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
}

// FrameworkComplianceResult - Compliance result for a framework
type FrameworkComplianceResult struct {
	FrameworkID       string                `json:"framework_id"`
	FrameworkName     string                `json:"framework_name"`
	Framework         string                `json:"framework"`
	Status            ComplianceStatus      `json:"status"`
	Score             float64               `json:"score"`
	ViolationCount    int                   `json:"violation_count"`
	WarningCount      int                   `json:"warning_count"`
	Violations        []ComplianceViolation `json:"violations"`
	Warnings          []ComplianceWarning   `json:"warnings"`
	Requirements      []RequirementResult   `json:"requirements"`
	RequirementsMet   int                   `json:"requirements_met"`
	TotalRequirements int                   `json:"total_requirements"`
	ProcessingTime    time.Duration         `json:"processing_time"`
}

// ComplianceCondition - Compliance rule condition
type ComplianceCondition struct {
	Field       string      `json:"field"`
	Operator    string      `json:"operator"`
	Value       interface{} `json:"value"`
	Weight      float64     `json:"weight"`
	Description string      `json:"description"`
}

// ComplianceAction - Compliance action
type ComplianceAction struct {
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
	Enabled     bool                   `json:"enabled"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
}

// ComplianceDataAnalysis - Data analysis result for compliance
type ComplianceDataAnalysis struct {
	DataType        string                 `json:"data_type"`
	Classification  string                 `json:"classification"`
	Sensitivity     string                 `json:"sensitivity"`
	Confidence      float64                `json:"confidence"`
	PIIDetected     bool                   `json:"pii_detected"`
	ComplianceRisk  float64                `json:"compliance_risk"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ComplianceIssue - Compliance issue found during validation
type ComplianceIssue struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    ComplianceSeverity     `json:"severity"`
	Framework   string                 `json:"framework"`
	Regulation  string                 `json:"regulation"`
	Article     string                 `json:"article"`
	Title       string                 `json:"title"`
	Rule        string                 `json:"rule"`
	Description string                 `json:"description"`
	Location    string                 `json:"location"`
	Evidence    string                 `json:"evidence"`
	Suggestion  string                 `json:"suggestion"`
	Category    ComplianceCategory     `json:"category"`
	Metadata    map[string]interface{} `json:"metadata"`
	DetectedAt  time.Time              `json:"detected_at"`
}

// ScoreAdjustment is already moved up

// RiskLevel - Risk level constants
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// ComplianceRequirement - Compliance framework requirement
type ComplianceRequirement struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Mandatory   bool                   `json:"mandatory"`
	Controls    []string               `json:"controls"`
	References  []string               `json:"references"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RiskProfileFilter - Filter for risk profiles
type RiskProfileFilter struct {
	Category string `json:"category,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
	Name     string `json:"name,omitempty"`
}

// RiskTrendFilter - Filter for risk trends
type RiskTrendFilter struct {
	APIID      string    `json:"api_id,omitempty"`
	EndpointID string    `json:"endpoint_id,omitempty"`
	TimeRange  TimeRange `json:"time_range,omitempty"`
}

// RiskTrendAnalysis - Risk trend analysis result
type RiskTrendAnalysis struct {
	ID              string           `json:"id"`
	GeneratedAt     time.Time        `json:"generated_at"`
	Filter          *RiskTrendFilter `json:"filter"`
	TimeRange       TimeRange        `json:"time_range"`
	Trends          []RiskTrend      `json:"trends"`
	Summary         RiskTrendSummary `json:"summary"`
	Recommendations []string         `json:"recommendations"`
}

// RiskTrend - Single point in a risk trend
type RiskTrend struct {
	Date         time.Time         `json:"date"`
	AverageScore float64           `json:"average_score"`
	RiskCounts   map[RiskLevel]int `json:"risk_counts"`
}

// RiskTrendSummary - Summary of risk trends
type RiskTrendSummary struct {
	TotalAssessments   int     `json:"total_assessments"`
	AverageScore       float64 `json:"average_score"`
	ScoreChange        float64 `json:"score_change"`
	TrendDirection     string  `json:"trend_direction"`
	HighestRiskAPI     string  `json:"highest_risk_api"`
	MostCommonRiskType string  `json:"most_common_risk_type"`
}

// RiskReportFilter - Filter for risk reports
type RiskReportFilter struct {
	APIID      string    `json:"api_id,omitempty"`
	EndpointID string    `json:"endpoint_id,omitempty"`
	TimeRange  TimeRange `json:"time_range,omitempty"`
}

// RiskReport - Risk assessment report
type RiskReport struct {
	ID              string                 `json:"id"`
	GeneratedAt     time.Time              `json:"generated_at"`
	Filter          *RiskReportFilter      `json:"filter"`
	Summary         RiskReportSummary      `json:"summary"`
	Details         []RiskAssessmentDetail `json:"details"`
	Trends          []RiskTrend            `json:"trends"`
	Recommendations []string               `json:"recommendations"`
}

// RiskReportSummary - Summary of a risk report
type RiskReportSummary struct {
	TotalAssessments int               `json:"total_assessments"`
	RiskDistribution map[RiskLevel]int `json:"risk_distribution"`
	AverageScore     float64           `json:"average_score"`
	TopRisks         []TopRisk         `json:"top_risks"`
}

// TopRisk - Most significant risk identified
type TopRisk struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    RiskLevel `json:"severity"`
	Score       float64   `json:"score"`
}

// RiskAssessmentDetail - Detailed risk assessment result
type RiskAssessmentDetail struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	APIID      string    `json:"api_id"`
	EndpointID string    `json:"endpoint_id"`
	RiskScore  float64   `json:"risk_score"`
	RiskLevel  RiskLevel `json:"risk_level"`
	AssessedAt time.Time `json:"assessed_at"`
}
