package models

import (
	"time"
)

type DataClassificationLevel string

const (
	DataClassificationPublic       DataClassificationLevel = "public"
	DataClassificationInternal     DataClassificationLevel = "internal"
	DataClassificationConfidential DataClassificationLevel = "confidential"
	DataClassificationRestricted   DataClassificationLevel = "restricted"
	DataClassificationTopSecret    DataClassificationLevel = "top_secret"
)

type DataCategory string

const (
	DataCategoryPersonal   DataCategory = "personal"
	DataCategoryFinancial  DataCategory = "financial"
	DataCategoryHealth     DataCategory = "health"
	DataCategoryBiometric  DataCategory = "biometric"
	DataCategoryLocation   DataCategory = "location"
	DataCategoryBehavioral DataCategory = "behavioral"
	DataCategoryTechnical  DataCategory = "technical"
	DataCategoryBusiness   DataCategory = "business"
	DataCategoryLegal      DataCategory = "legal"
	DataCategoryCustom     DataCategory = "custom"
)

type ClassificationMethod string

const (
	ClassificationMethodRule    ClassificationMethod = "rule_based"
	ClassificationMethodML      ClassificationMethod = "machine_learning"
	ClassificationMethodPattern ClassificationMethod = "pattern_matching"
	ClassificationMethodContext ClassificationMethod = "context_analysis"
	ClassificationMethodHybrid  ClassificationMethod = "hybrid"
	ClassificationMethodManual  ClassificationMethod = "manual"
)

type ClassificationCategory string

const (
	ClassificationCategoryPublic       ClassificationCategory = "public"
	ClassificationCategoryInternal     ClassificationCategory = "internal"
	ClassificationCategoryConfidential ClassificationCategory = "confidential"
	ClassificationCategoryRestricted   ClassificationCategory = "restricted"
)

type ClassificationRule struct {
	ID               string                    `json:"id" db:"id"`
	Name             string                    `json:"name" db:"name"`
	Description      string                    `json:"description" db:"description"`
	Category         DataCategory              `json:"category" db:"category"`
	Classification   DataClassificationLevel   `json:"classification" db:"classification"`
	Method           ClassificationMethod      `json:"method" db:"method"`
	Priority         int                       `json:"priority" db:"priority"`
	Enabled          bool                      `json:"enabled" db:"enabled"`
	Conditions       []ClassificationCondition `json:"conditions" db:"conditions"`
	Actions          []ClassificationAction    `json:"actions" db:"actions"`
	Patterns         []ClassificationPattern   `json:"patterns" db:"patterns"`
	MLModelConfig    *MLModelConfig            `json:"ml_model_config,omitempty" db:"ml_model_config"`
	ComplianceRules  []string                  `json:"compliance_rules" db:"compliance_rules"`
	RetentionPolicy  *RetentionPolicy          `json:"retention_policy,omitempty" db:"retention_policy"`
	AccessControls   []AccessControl           `json:"access_controls" db:"access_controls"`
	EncryptionConfig *EncryptionConfig         `json:"encryption_config,omitempty" db:"encryption_config"`
	Metadata         map[string]interface{}    `json:"metadata" db:"metadata"`
	Labels           []DataLabel               `json:"labels" db:"labels"`
	CreatedAt        time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at" db:"updated_at"`
	CreatedBy        string                    `json:"created_by" db:"created_by"`
	UpdatedBy        string                    `json:"updated_by" db:"updated_by"`
}

type ClassificationCondition struct {
	Field         string                 `json:"field"`
	Operator      string                 `json:"operator"`
	Value         interface{}            `json:"value"`
	CaseSensitive bool                   `json:"case_sensitive"`
	Weight        float64                `json:"weight"`
	Context       map[string]interface{} `json:"context"`
}

type ClassificationAction struct {
	Type       string                 `json:"type"`
	Config     map[string]interface{} `json:"config"`
	Enabled    bool                   `json:"enabled"`
	Priority   int                    `json:"priority"`
	Conditions []string               `json:"conditions"`
}

type ClassificationPattern struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Pattern     string  `json:"pattern"`
	Type        string  `json:"type"` // regex, keyword, ml_signature
	Weight      float64 `json:"weight"`
	Enabled     bool    `json:"enabled"`
	Description string  `json:"description"`
}

type MLModelConfig struct {
	ModelID      string                 `json:"model_id"`
	ModelType    string                 `json:"model_type"`
	Version      string                 `json:"version"`
	Threshold    float64                `json:"threshold"`
	Features     []string               `json:"features"`
	Parameters   map[string]interface{} `json:"parameters"`
	TrainingData string                 `json:"training_data"`
	LastTrained  time.Time              `json:"last_trained"`
}

type RetentionPolicy struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	RetentionPeriod time.Duration `json:"retention_period"`
	DeleteAfter     bool          `json:"delete_after"`
	ArchiveAfter    time.Duration `json:"archive_after"`
	NotifyBefore    time.Duration `json:"notify_before"`
	LegalHold       bool          `json:"legal_hold"`
	Exceptions      []string      `json:"exceptions"`
}

type AccessControl struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // role, user, group, attribute
	Principals  []string `json:"principals"`
	Permissions []string `json:"permissions"`
	Conditions  []string `json:"conditions"`
	Enabled     bool     `json:"enabled"`
}

type EncryptionConfig struct {
	Algorithm     string            `json:"algorithm"`
	KeySize       int               `json:"key_size"`
	Mode          string            `json:"mode"`
	KeyRotation   time.Duration     `json:"key_rotation"`
	KeyManagement string            `json:"key_management"`
	Parameters    map[string]string `json:"parameters"`
}

type DataClassificationRequest struct {
	RequestID   string                 `json:"request_id"`
	APIID       string                 `json:"api_id"`
	EndpointID  string                 `json:"endpoint_id"`
	Content     string                 `json:"content"`
	ContentType string                 `json:"content_type"`
	Source      string                 `json:"source"`
	Rules       []string               `json:"rules"`
	Options     ClassificationOptions  `json:"options"`
	Context     map[string]interface{} `json:"context"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Data        map[string]interface{} `json:"data"`
	DataSource  string                 `json:"data_source"`
}

type ClassificationOptions struct {
	EnableMLClassification        bool     `json:"enable_ml_classification"`
	EnableRuleBasedClassification bool     `json:"enable_rule_based_classification"`
	EnablePatternMatching         bool     `json:"enable_pattern_matching"`
	MinConfidenceScore            float64  `json:"min_confidence_score"`
	MaxClassificationDepth        int      `json:"max_classification_depth"`
	IncludeMetadata               bool     `json:"include_metadata"`
	EnableRealTimeProcessing      bool     `json:"enable_real_time_processing"`
	CustomRules                   []string `json:"custom_rules"`
}

type DataClassificationResult struct {
	RequestID       string                 `json:"request_id"`
	Classifications []DataClassification   `json:"classifications"`
	Summary         ClassificationSummary  `json:"summary"`
	Recommendations []string               `json:"recommendations"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	Metadata        map[string]interface{} `json:"metadata"`
	ClassifiedAt    time.Time              `json:"classified_at"`
	RulesMatched    []string               `json:"rules_matched"`
	AppliedLabels   []DataLabel            `json:"applied_labels"`
	ExecutedActions []ClassificationAction `json:"executed_actions"`
}

type DataClassification struct {
	ID                string                  `json:"id"`
	RequestID         string                  `json:"request_id"`
	APIID             string                  `json:"api_id"`
	EndpointID        string                  `json:"endpoint_id"`
	FieldName         string                  `json:"field_name"`
	FieldPath         string                  `json:"field_path"`
	Value             string                  `json:"value"`
	Category          DataCategory            `json:"category"`
	Classification    DataClassificationLevel `json:"classification"`
	Method            ClassificationMethod    `json:"method"`
	ConfidenceScore   float64                 `json:"confidence_score"`
	RuleID            string                  `json:"rule_id"`
	RuleName          string                  `json:"rule_name"`
	Location          DataLocation            `json:"location"`
	Context           ClassificationContext   `json:"context"`
	MatchContext      map[string]interface{}  `json:"match_context"`
	Confidence        float64                 `json:"confidence"`
	Labels            []DataLabel             `json:"labels"`
	ComplianceFlags   []string                `json:"compliance_flags"`
	ProcessingActions []ProcessingAction      `json:"processing_actions"`
	Metadata          map[string]interface{}  `json:"metadata"`
	ClassifiedAt      time.Time               `json:"classified_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

type DataLocation struct {
	Source       string `json:"source"`
	Section      string `json:"section"`
	StartIndex   int    `json:"start_index"`
	EndIndex     int    `json:"end_index"`
	LineNumber   int    `json:"line_number"`
	ColumnNumber int    `json:"column_number"`
}

type ClassificationContext struct {
	SurroundingText   string            `json:"surrounding_text"`
	FieldDescription  string            `json:"field_description"`
	DataFormat        string            `json:"data_format"`
	BusinessContext   string            `json:"business_context"`
	TechnicalContext  string            `json:"technical_context"`
	RelatedFields     []string          `json:"related_fields"`
	DataFlow          string            `json:"data_flow"`
	ProcessingPurpose string            `json:"processing_purpose"`
	LegalBasis        string            `json:"legal_basis"`
	Attributes        map[string]string `json:"attributes"`
}

type ProcessingAction struct {
	Action    string                 `json:"action"`
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
	Result    string                 `json:"result"`
	Error     string                 `json:"error,omitempty"`
}

type ClassificationSummary struct {
	TotalClassifications      int                             `json:"total_classifications"`
	ClassificationsByLevel    map[DataClassificationLevel]int `json:"classifications_by_level"`
	ClassificationsByCategory map[DataCategory]int            `json:"classifications_by_category"`
	ClassificationsByMethod   map[ClassificationMethod]int    `json:"classifications_by_method"`
	HighRiskClassifications   int                             `json:"high_risk_classifications"`
	ComplianceImpact          []string                        `json:"compliance_impact"`
	RecommendedActions        []string                        `json:"recommended_actions"`
	OverallRiskScore          float64                         `json:"overall_risk_score"`
	ProcessingErrors          []string                        `json:"processing_errors"`
	CategoryBreakdown         map[string]int                  `json:"category_breakdown"`
	LabelBreakdown            map[string]int                  `json:"label_breakdown"`
	RuleUsage                 map[string]int                  `json:"rule_usage"`
}

type ClassificationFilter struct {
	APIID           string                  `json:"api_id,omitempty"`
	EndpointID      string                  `json:"endpoint_id,omitempty"`
	Category        DataCategory            `json:"category,omitempty"`
	Classification  DataClassificationLevel `json:"classification,omitempty"`
	Method          ClassificationMethod    `json:"method,omitempty"`
	RuleID          string                  `json:"rule_id,omitempty"`
	StartDate       *time.Time              `json:"start_date,omitempty"`
	EndDate         *time.Time              `json:"end_date,omitempty"`
	MinConfidence   *float64                `json:"min_confidence,omitempty"`
	MaxConfidence   *float64                `json:"max_confidence,omitempty"`
	ComplianceFlags []string                `json:"compliance_flags,omitempty"`
	Limit           int                     `json:"limit,omitempty"`
	Offset          int                     `json:"offset,omitempty"`
}

type ClassificationStatistics struct {
	TotalClassifications      int                             `json:"total_classifications"`
	ClassificationsByLevel    map[DataClassificationLevel]int `json:"classifications_by_level"`
	ClassificationsByCategory map[DataCategory]int            `json:"classifications_by_category"`
	ClassificationsByAPI      map[string]int                  `json:"classifications_by_api"`
	ClassificationsByEndpoint map[string]int                  `json:"classifications_by_endpoint"`
	ClassificationTrends      []ClassificationTrend           `json:"classification_trends"`
	ComplianceImpact          map[string]int                  `json:"compliance_impact"`
	RiskDistribution          map[string]int                  `json:"risk_distribution"`
	ProcessingActions         map[string]int                  `json:"processing_actions"`
	GeneratedAt               time.Time                       `json:"generated_at"`
}

type ClassificationTrend struct {
	Date                    time.Time `json:"date"`
	ClassificationsCreated  int       `json:"classifications_created"`
	HighRiskClassifications int       `json:"high_risk_classifications"`
	ComplianceIssues        int       `json:"compliance_issues"`
}

// =============================================================================
// ADDITIONAL CLASSIFICATION TYPES
// =============================================================================

// ClassificationRuleFilter - Filter for classification rules
type ClassificationRuleFilter struct {
	Category string   `json:"category,omitempty"`
	Name     string   `json:"name,omitempty"`
	Method   string   `json:"method,omitempty"`
	Enabled  *bool    `json:"enabled,omitempty"`
	Priority *int     `json:"priority,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

// DataLabel - Data classification label
type DataLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ClassificationReportFilter - Filter for classification reports
type ClassificationReportFilter struct {
	Category  string     `json:"category,omitempty"`
	Since     *time.Time `json:"since,omitempty"`
	Until     *time.Time `json:"until,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

// ClassificationReport - Classification analysis report
type ClassificationReport struct {
	ID              string                      `json:"id"`
	Details         map[string]interface{}      `json:"details"`
	GeneratedAt     time.Time                   `json:"generated_at"`
	Filter          *ClassificationReportFilter `json:"filter"`
	Summary         ClassificationSummary       `json:"summary"`
	Classifications []DataClassification        `json:"classifications"`
	Trends          []ClassificationTrend       `json:"trends"`
	Recommendations []string                    `json:"recommendations"`
}

type ClassificationData struct {
	ID              string                 `json:"id" db:"id"`
	RequestID       string                 `json:"request_id" db:"request_id"`
	APIID           string                 `json:"api_id" db:"api_id"`
	EndpointID      string                 `json:"endpoint_id" db:"endpoint_id"`
	Classifications []DataClassification   `json:"classifications" db:"classifications"`
	Labels          []DataLabel            `json:"labels" db:"labels"`
	DataHash        string                 `json:"data_hash" db:"data_hash"`
	IPAddress       string                 `json:"ip_address" db:"ip_address"`
	UserAgent       string                 `json:"user_agent" db:"user_agent"`
	ClassifiedAt    time.Time              `json:"classified_at" db:"classified_at"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
}

type DataAnalysis struct {
	FieldCount     int            `json:"field_count"`
	FieldTypes     map[string]int `json:"field_types"`
	FieldNames     []string       `json:"field_names"`
	DataPatterns   []string       `json:"data_patterns"`
	NestedLevels   int            `json:"nested_levels"`
	ArrayFields    []string       `json:"array_fields"`
	SensitiveHints []string       `json:"sensitive_hints"`
}
