package models

import (
	"time"
	"encoding/json"
)

// BlockingRule represents a rule for blocking malicious traffic
type BlockingRule struct {
	ID                string                 `json:"id" db:"id"`
	Name              string                 `json:"name" db:"name"`
	Description       string                 `json:"description" db:"description"`
	Type              BlockingRuleType       `json:"type" db:"type"`
	Status            RuleStatus             `json:"status" db:"status"`
	Priority          int                    `json:"priority" db:"priority"`
	Conditions        []*RuleCondition       `json:"conditions" db:"conditions"`
	Actions           []*RuleAction          `json:"actions" db:"actions"`
	Metadata          map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedBy         string                 `json:"created_by" db:"created_by"`
	UpdatedBy         string                 `json:"updated_by" db:"updated_by"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
	LastTriggered     *time.Time             `json:"last_triggered" db:"last_triggered"`
	TriggerCount      int64                  `json:"trigger_count" db:"trigger_count"`
	EffectivenessScore float64               `json:"effectiveness_score" db:"effectiveness_score"`
	Tags              []string               `json:"tags" db:"tags"`
	ExpiresAt         *time.Time             `json:"expires_at" db:"expires_at"`
	IsTemporary       bool                   `json:"is_temporary" db:"is_temporary"`
	SourceRuleID      string                 `json:"source_rule_id" db:"source_rule_id"`
	Version           int                    `json:"version" db:"version"`
}

// BlockingRuleType defines the type of blocking rule
type BlockingRuleType string

const (
	BlockingRuleTypeSignature    BlockingRuleType = "signature"
	BlockingRuleTypeAnomaly      BlockingRuleType = "anomaly"
	BlockingRuleTypeBehavioral   BlockingRuleType = "behavioral"
	BlockingRuleTypeReputation   BlockingRuleType = "reputation"
	BlockingRuleTypeRateLimit    BlockingRuleType = "rate_limit"
	BlockingRuleTypeGeoLocation  BlockingRuleType = "geo_location"
	BlockingRuleTypeCustom       BlockingRuleType = "custom"
	BlockingRuleTypeML           BlockingRuleType = "ml_based"
	BlockingRuleTypeAdaptive     BlockingRuleType = "adaptive"
)

// RuleCondition represents a condition within a blocking rule
type RuleCondition struct {
	ID          string                 `json:"id"`
	Field       string                 `json:"field"`
	Operator    ConditionOperator      `json:"operator"`
	Value       interface{}            `json:"value"`
	ValueType   ValueType              `json:"value_type"`
	CaseSensitive bool                 `json:"case_sensitive"`
	Negate      bool                   `json:"negate"`
	Weight      float64                `json:"weight"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ConditionOperator defines operators for rule conditions
type ConditionOperator string

const (
	OperatorEquals          ConditionOperator = "equals"
	OperatorNotEquals       ConditionOperator = "not_equals"
	OperatorContains        ConditionOperator = "contains"
	OperatorNotContains     ConditionOperator = "not_contains"
	OperatorStartsWith      ConditionOperator = "starts_with"
	OperatorEndsWith        ConditionOperator = "ends_with"
	OperatorRegex           ConditionOperator = "regex"
	OperatorGreaterThan     ConditionOperator = "greater_than"
	OperatorLessThan        ConditionOperator = "less_than"
	OperatorGreaterEqual    ConditionOperator = "greater_equal"
	OperatorLessEqual       ConditionOperator = "less_equal"
	OperatorInList          ConditionOperator = "in_list"
	OperatorNotInList       ConditionOperator = "not_in_list"
	OperatorExists          ConditionOperator = "exists"
	OperatorNotExists       ConditionOperator = "not_exists"
	OperatorIPInRange       ConditionOperator = "ip_in_range"
	OperatorIPNotInRange    ConditionOperator = "ip_not_in_range"
	OperatorGeoLocation     ConditionOperator = "geo_location"
	OperatorTimeRange       ConditionOperator = "time_range"
)

// ValueType defines the type of condition value
type ValueType string

const (
	ValueTypeString    ValueType = "string"
	ValueTypeNumber    ValueType = "number"
	ValueTypeBoolean   ValueType = "boolean"
	ValueTypeList      ValueType = "list"
	ValueTypeRegex     ValueType = "regex"
	ValueTypeIPAddress ValueType = "ip_address"
	ValueTypeIPRange   ValueType = "ip_range"
	ValueTypeTimestamp ValueType = "timestamp"
	ValueTypeDuration  ValueType = "duration"
)

// RuleAction represents an action to take when a rule is triggered
type RuleAction struct {
	ID          string                 `json:"id"`
	Type        ActionType             `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ActionType defines the type of action to take
type ActionType string

const (
	ActionTypeBlock           ActionType = "block"
	ActionTypeAllow           ActionType = "allow"
	ActionTypeLog             ActionType = "log"
	ActionTypeAlert           ActionType = "alert"
	ActionTypeRateLimit       ActionType = "rate_limit"
	ActionTypeRedirect        ActionType = "redirect"
	ActionTypeModifyRequest   ActionType = "modify_request"
	ActionTypeModifyResponse  ActionType = "modify_response"
	ActionTypeQuarantine      ActionType = "quarantine"
	ActionTypeNotify          ActionType = "notify"
	ActionTypeUpdateRule      ActionType = "update_rule"
	ActionTypeCreateRule      ActionType = "create_rule"
	ActionTypeCustom          ActionType = "custom"
)

// RuleStatus defines the status of a blocking rule
type RuleStatus string

const (
	RuleStatusActive    RuleStatus = "active"
	RuleStatusInactive  RuleStatus = "inactive"
	RuleStatusTesting   RuleStatus = "testing"
	RuleStatusArchived  RuleStatus = "archived"
	RuleStatusExpired   RuleStatus = "expired"
	RuleStatusError     RuleStatus = "error"
)

// BlockingRuleFilter represents filters for querying blocking rules
type BlockingRuleFilter struct {
	IDs          []string           `json:"ids"`
	Names        []string           `json:"names"`
	Types        []BlockingRuleType `json:"types"`
	Statuses     []RuleStatus       `json:"statuses"`
	Tags         []string           `json:"tags"`
	CreatedBy    string             `json:"created_by"`
	CreatedAfter *time.Time         `json:"created_after"`
	CreatedBefore *time.Time        `json:"created_before"`
	UpdatedAfter *time.Time         `json:"updated_after"`
	UpdatedBefore *time.Time        `json:"updated_before"`
	Priority     *int               `json:"priority"`
	MinPriority  *int               `json:"min_priority"`
	MaxPriority  *int               `json:"max_priority"`
	Search       string             `json:"search"`
	Limit        int                `json:"limit"`
	Offset       int                `json:"offset"`
	SortBy       string             `json:"sort_by"`
	SortOrder    string             `json:"sort_order"`
}

// RuleMatch represents a match of a blocking rule against a request
type RuleMatch struct {
	ID                string                 `json:"id"`
	RuleID            string                 `json:"rule_id"`
	RuleName          string                 `json:"rule_name"`
	RequestID         string                 `json:"request_id"`
	MatchedConditions []*ConditionMatch      `json:"matched_conditions"`
	Confidence        float64                `json:"confidence"`
	Severity          string                 `json:"severity"`
	Actions           []*ActionResult        `json:"actions"`
	Timestamp         time.Time              `json:"timestamp"`
	ProcessingTime    time.Duration          `json:"processing_time"`
	Metadata          map[string]interface{} `json:"metadata"`
	SourceIP          string                 `json:"source_ip"`
	UserAgent         string                 `json:"user_agent"`
	RequestPath       string                 `json:"request_path"`
	RequestMethod     string                 `json:"request_method"`
}

// ConditionMatch represents a matched condition
type ConditionMatch struct {
	ConditionID   string      `json:"condition_id"`
	Field         string      `json:"field"`
	Operator      string      `json:"operator"`
	ExpectedValue interface{} `json:"expected_value"`
	ActualValue   interface{} `json:"actual_value"`
	Confidence    float64     `json:"confidence"`
	Weight        float64     `json:"weight"`
	Timestamp     time.Time   `json:"timestamp"`
}

// ActionResult represents the result of executing an action
type ActionResult struct {
	ActionID    string                 `json:"action_id"`
	ActionType  ActionType             `json:"action_type"`
	Status      ActionStatus           `json:"status"`
	Message     string                 `json:"message"`
	ExecutedAt  time.Time              `json:"executed_at"`
	Duration    time.Duration          `json:"duration"`
	Parameters  map[string]interface{} `json:"parameters"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ActionStatus defines the status of action execution
type ActionStatus string

const (
	ActionStatusSuccess ActionStatus = "success"
	ActionStatusFailed  ActionStatus = "failed"
	ActionStatusSkipped ActionStatus = "skipped"
	ActionStatusPending ActionStatus = "pending"
)

// RulePerformance represents performance metrics for a blocking rule
type RulePerformance struct {
	RuleID              string    `json:"rule_id"`
	RuleName            string    `json:"rule_name"`
	TotalMatches        int64     `json:"total_matches"`
	TruePositives       int64     `json:"true_positives"`
	FalsePositives      int64     `json:"false_positives"`
	TrueNegatives       int64     `json:"true_negatives"`
	FalseNegatives      int64     `json:"false_negatives"`
	Precision           float64   `json:"precision"`
	Recall              float64   `json:"recall"`
	F1Score             float64   `json:"f1_score"`
	Accuracy            float64   `json:"accuracy"`
	EffectivenessScore  float64   `json:"effectiveness_score"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	LastEvaluated       time.Time `json:"last_evaluated"`
	Period              string    `json:"period"`
}

// RuleTemplate represents a template for creating blocking rules
type RuleTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Type        BlockingRuleType       `json:"type"`
	Template    *BlockingRule          `json:"template"`
	Variables   []*TemplateVariable    `json:"variables"`
	Tags        []string               `json:"tags"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	UsageCount  int64                  `json:"usage_count"`
	Rating      float64                `json:"rating"`
	IsPublic    bool                   `json:"is_public"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TemplateVariable represents a variable in a rule template
type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         ValueType   `json:"type"`
	Description  string      `json:"description"`
	DefaultValue interface{} `json:"default_value"`
	Required     bool        `json:"required"`
	Validation   string      `json:"validation"`
	Options      []string    `json:"options,omitempty"`
}

// RuleBatch represents a batch operation on blocking rules
type RuleBatch struct {
	ID          string                 `json:"id"`
	Operation   BatchOperation         `json:"operation"`
	Rules       []*BlockingRule        `json:"rules"`
	Status      BatchStatus            `json:"status"`
	Progress    int                    `json:"progress"`
	Total       int                    `json:"total"`
	Results     []*BatchResult         `json:"results"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// BatchOperation defines the type of batch operation
type BatchOperation string

const (
	BatchOperationCreate   BatchOperation = "create"
	BatchOperationUpdate   BatchOperation = "update"
	BatchOperationDelete   BatchOperation = "delete"
	BatchOperationActivate BatchOperation = "activate"
	BatchOperationDeactivate BatchOperation = "deactivate"
	BatchOperationTest     BatchOperation = "test"
)

// BatchStatus defines the status of a batch operation
type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "pending"
	BatchStatusRunning    BatchStatus = "running"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
	BatchStatusCancelled  BatchStatus = "cancelled"
)

// BatchResult represents the result of a single operation in a batch
type BatchResult struct {
	RuleID    string    `json:"rule_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Helper methods for BlockingRule

// IsActive returns true if the rule is active
func (r *BlockingRule) IsActive() bool {
	return r.Status == RuleStatusActive
}

// IsExpired returns true if the rule has expired
func (r *BlockingRule) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// HasTag returns true if the rule has the specified tag
func (r *BlockingRule) HasTag(tag string) bool {
	for _, t := range r.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the rule if it doesn't already exist
func (r *BlockingRule) AddTag(tag string) {
	if !r.HasTag(tag) {
		r.Tags = append(r.Tags, tag)
	}
}

// RemoveTag removes a tag from the rule
func (r *BlockingRule) RemoveTag(tag string) {
	for i, t := range r.Tags {
		if t == tag {
			r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
			break
		}
	}
}

// GetConditionByID returns a condition by its ID
func (r *BlockingRule) GetConditionByID(id string) *RuleCondition {
	for _, condition := range r.Conditions {
		if condition.ID == id {
			return condition
		}
	}
	return nil
}

// GetActionByID returns an action by its ID
func (r *BlockingRule) GetActionByID(id string) *RuleAction {
	for _, action := range r.Actions {
		if action.ID == id {
			return action
		}
	}
	return nil
}

// Validate validates the blocking rule
func (r *BlockingRule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if len(r.Conditions) == 0 {
		return fmt.Errorf("at least one condition is required")
	}

	if len(r.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	// Validate conditions
	for i, condition := range r.Conditions {
		if err := condition.Validate(); err != nil {
			return fmt.Errorf("condition %d: %w", i, err)
		}
	}

	// Validate actions
	for i, action := range r.Actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("action %d: %w", i, err)
		}
	}

	return nil
}

// Clone creates a deep copy of the blocking rule
func (r *BlockingRule) Clone() *BlockingRule {
	clone := &BlockingRule{
		ID:                 r.ID,
		Name:               r.Name,
		Description:        r.Description,
		Type:               r.Type,
		Status:             r.Status,
		Priority:           r.Priority,
		CreatedBy:          r.CreatedBy,
		UpdatedBy:          r.UpdatedBy,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
		TriggerCount:       r.TriggerCount,
		EffectivenessScore: r.EffectivenessScore,
		IsTemporary:        r.IsTemporary,
		SourceRuleID:       r.SourceRuleID,
		Version:            r.Version,
	}

	if r.LastTriggered != nil {
		lastTriggered := *r.LastTriggered
		clone.LastTriggered = &lastTriggered
	}

	if r.ExpiresAt != nil {
		expiresAt := *r.ExpiresAt
		clone.ExpiresAt = &expiresAt
	}

	// Clone tags
	clone.Tags = make([]string, len(r.Tags))
	copy(clone.Tags, r.Tags)

	// Clone conditions
	clone.Conditions = make([]*RuleCondition, len(r.Conditions))
	for i, condition := range r.Conditions {
		clone.Conditions[i] = condition.Clone()
	}

	// Clone actions
	clone.Actions = make([]*RuleAction, len(r.Actions))
	for i, action := range r.Actions {
		clone.Actions[i] = action.Clone()
	}

	// Clone metadata
	if r.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range r.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

// Helper methods for RuleCondition

// Validate validates the rule condition
func (c *RuleCondition) Validate() error {
	if c.Field == "" {
		return fmt.Errorf("condition field is required")
	}

	if c.Operator == "" {
		return fmt.Errorf("condition operator is required")
	}

	if c.Value == nil && c.Operator != OperatorExists && c.Operator != OperatorNotExists {
		return fmt.Errorf("condition value is required for operator %s", c.Operator)
	}

	// Validate operator-specific requirements
	switch c.Operator {
	case OperatorRegex:
		if c.ValueType != ValueTypeRegex && c.ValueType != ValueTypeString {
			return fmt.Errorf("regex operator requires regex or string value type")
		}
	case OperatorIPInRange, OperatorIPNotInRange:
		if c.ValueType != ValueTypeIPRange {
			return fmt.Errorf("IP range operator requires ip_range value type")
		}
	case OperatorInList, OperatorNotInList:
		if c.ValueType != ValueTypeList {
			return fmt.Errorf("list operator requires list value type")
		}
	}

	return nil
}

// Clone creates a deep copy of the rule condition
func (c *RuleCondition) Clone() *RuleCondition {
	clone := &RuleCondition{
		ID:            c.ID,
		Field:         c.Field,
		Operator:      c.Operator,
		Value:         c.Value,
		ValueType:     c.ValueType,
		CaseSensitive: c.CaseSensitive,
		Negate:        c.Negate,
		Weight:        c.Weight,
		Description:   c.Description,
	}

	// Clone metadata
	if c.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range c.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

// Helper methods for RuleAction

// Validate validates the rule action
func (a *RuleAction) Validate() error {
	if a.Type == "" {
		return fmt.Errorf("action type is required")
	}

	// Validate action-specific parameters
	switch a.Type {
	case ActionTypeBlock:
		if duration, exists := a.Parameters["duration"]; exists {
			if _, ok := duration.(float64); !ok {
				return fmt.Errorf("block action duration must be a number")
			}
		}
	case ActionTypeRateLimit:
		if limit, exists := a.Parameters["limit"]; exists {
			if _, ok := limit.(float64); !ok {
				return fmt.Errorf("rate limit action limit must be a number")
			}
		}
		if window, exists := a.Parameters["window"]; exists {
			if _, ok := window.(float64); !ok {
				return fmt.Errorf("rate limit action window must be a number")
			}
		}
	case ActionTypeRedirect:
		if url, exists := a.Parameters["url"]; exists {
			if _, ok := url.(string); !ok {
				return fmt.Errorf("redirect action url must be a string")
			}
		} else {
			return fmt.Errorf("redirect action requires url parameter")
		}
	case ActionTypeNotify:
		if recipients, exists := a.Parameters["recipients"]; exists {
			if _, ok := recipients.([]interface{}); !ok {
				return fmt.Errorf("notify action recipients must be a list")
			}
		} else {
			return fmt.Errorf("notify action requires recipients parameter")
		}
	}

	return nil
}

// Clone creates a deep copy of the rule action
func (a *RuleAction) Clone() *RuleAction {
	clone := &RuleAction{
		ID:          a.ID,
		Type:        a.Type,
		Priority:    a.Priority,
		Description: a.Description,
		Enabled:     a.Enabled,
	}

	// Clone parameters
	if a.Parameters != nil {
		clone.Parameters = make(map[string]interface{})
		for k, v := range a.Parameters {
			clone.Parameters[k] = v
		}
	}

	// Clone metadata
	if a.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range a.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

// Helper methods for RuleMatch

// GetTotalConfidence calculates the total confidence score
func (m *RuleMatch) GetTotalConfidence() float64 {
	if len(m.MatchedConditions) == 0 {
		return m.Confidence
	}

	totalWeight := 0.0
	weightedConfidence := 0.0

	for _, condition := range m.MatchedConditions {
		totalWeight += condition.Weight
		weightedConfidence += condition.Confidence * condition.Weight
	}

	if totalWeight == 0 {
		return m.Confidence
	}

	return weightedConfidence / totalWeight
}

// GetSuccessfulActions returns actions that executed successfully
func (m *RuleMatch) GetSuccessfulActions() []*ActionResult {
	successful := make([]*ActionResult, 0)
	for _, action := range m.Actions {
		if action.Status == ActionStatusSuccess {
			successful = append(successful, action)
		}
	}
	return successful
}

// GetFailedActions returns actions that failed to execute
func (m *RuleMatch) GetFailedActions() []*ActionResult {
	failed := make([]*ActionResult, 0)
	for _, action := range m.Actions {
		if action.Status == ActionStatusFailed {
			failed = append(failed, action)
		}
	}
	return failed
}

// Helper methods for RulePerformance

// CalculateMetrics calculates performance metrics
func (p *RulePerformance) CalculateMetrics() {
	if p.TruePositives+p.FalsePositives > 0 {
		p.Precision = float64(p.TruePositives) / float64(p.TruePositives+p.FalsePositives)
	}

	if p.TruePositives+p.FalseNegatives > 0 {
		p.Recall = float64(p.TruePositives) / float64(p.TruePositives+p.FalseNegatives)
	}

	if p.Precision+p.Recall > 0 {
		p.F1Score = 2 * (p.Precision * p.Recall) / (p.Precision + p.Recall)
	}

	total := p.TruePositives + p.TrueNegatives + p.FalsePositives + p.FalseNegatives
	if total > 0 {
		p.Accuracy = float64(p.TruePositives+p.TrueNegatives) / float64(total)
	}

	// Calculate effectiveness score (weighted combination of precision, recall, and accuracy)
	p.EffectivenessScore = (0.4*p.Precision + 0.4*p.Recall + 0.2*p.Accuracy)
}

// IsEffective returns true if the rule is considered effective
func (p *RulePerformance) IsEffective() bool {
	return p.EffectivenessScore >= 0.7 && p.F1Score >= 0.6
}

// Helper methods for RuleTemplate

// CreateRule creates a blocking rule from the template
func (t *RuleTemplate) CreateRule(variables map[string]interface{}) (*BlockingRule, error) {
	if t.Template == nil {
		return nil, fmt.Errorf("template rule is nil")
	}

	// Validate required variables
	for _, variable := range t.Variables {
		if variable.Required {
			if _, exists := variables[variable.Name]; !exists {
				if variable.DefaultValue == nil {
					return nil, fmt.Errorf("required variable %s is missing", variable.Name)
				}
				variables[variable.Name] = variable.DefaultValue
			}
		}
	}

	// Clone the template rule
	rule := t.Template.Clone()

	// Apply variables to the rule
	if err := t.applyVariables(rule, variables); err != nil {
		return nil, fmt.Errorf("failed to apply variables: %w", err)
	}

	// Generate new ID and update timestamps
	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	rule.Version = 1

	return rule, nil
}

// applyVariables applies template variables to the rule
func (t *RuleTemplate) applyVariables(rule *BlockingRule, variables map[string]interface{}) error {
	// Convert rule to JSON for template processing
	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	// Apply template variables
	ruleStr := string(ruleJSON)
	for name, value := range variables {
		placeholder := fmt.Sprintf("{{.%s}}", name)
		valueStr := fmt.Sprintf("%v", value)
		ruleStr = strings.ReplaceAll(ruleStr, placeholder, valueStr)
	}

	// Convert back to rule
	return json.Unmarshal([]byte(ruleStr), rule)
}

// Validate validates the rule template
func (t *RuleTemplate) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if t.Template == nil {
		return fmt.Errorf("template rule is required")
	}

	if err := t.Template.Validate(); err != nil {
		return fmt.Errorf("template rule validation failed: %w", err)
	}

	// Validate variables
	for i, variable := range t.Variables {
		if variable.Name == "" {
			return fmt.Errorf("variable %d: name is required", i)
		}
		if variable.Type == "" {
			return fmt.Errorf("variable %d: type is required", i)
		}
	}

	return nil
}

// JSON marshaling helpers

// MarshalJSON customizes JSON marshaling for BlockingRule
func (r *BlockingRule) MarshalJSON() ([]byte, error) {
	type Alias BlockingRule
	return json.Marshal(&struct {
		*Alias
		Conditions string `json:"conditions"`
		Actions    string `json:"actions"`
		Metadata   string `json:"metadata"`
		Tags       string `json:"tags"`
	}{
		Alias:      (*Alias)(r),
		Conditions: r.conditionsToJSON(),
		Actions:    r.actionsToJSON(),
		Metadata:   r.metadataToJSON(),
		Tags:       r.tagsToJSON(),
	})
}

// UnmarshalJSON customizes JSON unmarshaling for BlockingRule
func (r *BlockingRule) UnmarshalJSON(data []byte) error {
	type Alias BlockingRule
	aux := &struct {
		*Alias
		Conditions string `json:"conditions"`
		Actions    string `json:"actions"`
		Metadata   string `json:"metadata"`
		Tags       string `json:"tags"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse conditions
	if aux.Conditions != "" {
		if err := json.Unmarshal([]byte(aux.Conditions), &r.Conditions); err != nil {
			return fmt.Errorf("failed to unmarshal conditions: %w", err)
		}
	}

	// Parse actions
	if aux.Actions != "" {
		if err := json.Unmarshal([]byte(aux.Actions), &r.Actions); err != nil {
			return fmt.Errorf("failed to unmarshal actions: %w", err)
		}
	}

	// Parse metadata
	if aux.Metadata != "" {
		if err := json.Unmarshal([]byte(aux.Metadata), &r.Metadata); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	// Parse tags
	if aux.Tags != "" {
		if err := json.Unmarshal([]byte(aux.Tags), &r.Tags); err != nil {
			return fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	return nil
}

// Helper methods for JSON serialization
func (r *BlockingRule) conditionsToJSON() string {
	if r.Conditions == nil {
		return "[]"
	}
	data, _ := json.Marshal(r.Conditions)
	return string(data)
}

func (r *BlockingRule) actionsToJSON() string {
	if r.Actions == nil {
		return "[]"
	}
	data, _ := json.Marshal(r.Actions)
	return string(data)
}

func (r *BlockingRule) metadataToJSON() string {
	if r.Metadata == nil {
		return "{}"
	}
	data, _ := json.Marshal(r.Metadata)
	return string(data)
}

func (r *BlockingRule) tagsToJSON() string {
	if r.Tags == nil {
		return "[]"
	}
	data, _ := json.Marshal(r.Tags)
	return string(data)
}

// Database scanning helpers
func (r *BlockingRule) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return fmt.Errorf("cannot scan %T into BlockingRule", value)
	}
}

// Value implements the driver.Valuer interface
func (r *BlockingRule) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Constants for common rule fields
const (
	FieldSourceIP       = "source_ip"
	FieldDestinationIP  = "destination_ip"
	FieldUserAgent      = "user_agent"
	FieldRequestPath    = "request_path"
	FieldRequestMethod  = "request_method"
	FieldRequestHeaders = "request_headers"
	FieldRequestBody    = "request_body"
	FieldResponseCode   = "response_code"
	FieldResponseSize   = "response_size"
	FieldRequestSize    = "request_size"
	FieldTimestamp      = "timestamp"
	FieldDuration       = "duration"
	FieldGeoLocation    = "geo_location"
	FieldASN            = "asn"
	FieldDomain         = "domain"
	FieldReferer        = "referer"
	FieldCookies        = "cookies"
	FieldQueryParams    = "query_params"
	FieldContentType    = "content_type"
	FieldProtocol       = "protocol"
	FieldPort           = "port"
	FieldTLSVersion     = "tls_version"
	FieldCertificate    = "certificate"
)

// Predefined rule templates
var (
	SQLInjectionTemplate = &RuleTemplate{
		ID:          "sql-injection-basic",
		Name:        "SQL Injection Detection",
		Description: "Detects basic SQL injection attempts",
		Category:    "Web Application Security",
		Type:        BlockingRuleTypeSignature,
		Template: &BlockingRule{
			Name:        "SQL Injection - {{.severity}}",
			Description: "Detects SQL injection patterns in {{.target_field}}",
			Type:        BlockingRuleTypeSignature,
			Priority:    90,
			Conditions: []*RuleCondition{
				{
					Field:    "{{.target_field}}",
					Operator: OperatorRegex,
					Value:    "(?i)(union|select|insert|update|delete|drop|create|alter)\\s",
					ValueType: ValueTypeRegex,
					Weight:   1.0,
				},
			},
			Actions: []*RuleAction{
				{
					Type: ActionTypeBlock,
					Parameters: map[string]interface{}{
						"duration": 300,
						"reason":   "SQL injection attempt detected",
					},
				},
				{
					Type: ActionTypeAlert,
					Parameters: map[string]interface{}{
						"severity": "{{.severity}}",
						"message":  "SQL injection detected in {{.target_field}}",
					},
				},
			},
		},
		Variables: []*TemplateVariable{
			{
				Name:         "target_field",
				Type:         ValueTypeString,
				Description:  "Field to monitor for SQL injection",
				DefaultValue: "request_body",
				Required:     true,
				Options:      []string{"request_body", "query_params", "request_headers"},
			},
			{
				Name:         "severity",
				Type:         ValueTypeString,
				Description:  "Alert severity level",
				DefaultValue: "high",
				Required:     false,
				Options:      []string{"low", "medium", "high", "critical"},
			},
		},
	}

	XSSTemplate = &RuleTemplate{
		ID:          "xss-basic",
		Name:        "Cross-Site Scripting Detection",
		Description: "Detects basic XSS attempts",
		Category:    "Web Application Security",
		Type:        BlockingRuleTypeSignature,
		Template: &BlockingRule{
			Name:        "XSS Detection - {{.severity}}",
			Description: "Detects XSS patterns in {{.target_field}}",
			Type:        BlockingRuleTypeSignature,
			Priority:    85,
			Conditions: []*RuleCondition{
				{
					Field:    "{{.target_field}}",
					Operator: OperatorRegex,
					Value:    "(?i)<script|javascript:|on\\w+\\s*=",
					ValueType: ValueTypeRegex,
					Weight:   1.0,
				},
			},
			Actions: []*RuleAction{
				{
					Type: ActionTypeBlock,
					Parameters: map[string]interface{}{
						"duration": 180,
						"reason":   "XSS attempt detected",
					},
				},
			},
		},
		Variables: []*TemplateVariable{
			{
				Name:         "target_field",
				Type:         ValueTypeString,
				Description:  "Field to monitor for XSS",
				DefaultValue: "request_body",
				Required:     true,
			},
			{
				Name:         "severity",
				Type:         ValueTypeString,
				Description:  "Alert severity level",
				DefaultValue: "high",
				Required:     false,
			},
		},
	}

	RateLimitTemplate = &RuleTemplate{
		ID:          "rate-limit-basic",
		Name:        "Rate Limiting",
		Description: "Basic rate limiting by IP address",
		Category:    "Traffic Control",
		Type:        BlockingRuleTypeRateLimit,
		Template: &BlockingRule{
			Name:        "Rate Limit - {{.requests_per_minute}} req/min",
			Description: "Limits requests to {{.requests_per_minute}} per minute per IP",
			Type:        BlockingRuleTypeRateLimit,
			Priority:    50,
			Conditions: []*RuleCondition{
				{
					Field:    "source_ip",
					Operator: OperatorExists,
					Weight:   1.0,
				},
			},
			Actions: []*RuleAction{
				{
					Type: ActionTypeRateLimit,
					Parameters: map[string]interface{}{
						"limit":  "{{.requests_per_minute}}",
						"window": 60,
						"key":    "source_ip",
					},
				},
			},
		},
		Variables: []*TemplateVariable{
			{
				Name:         "requests_per_minute",
				Type:         ValueTypeNumber,
				Description:  "Maximum requests per minute",
				DefaultValue: 100,
				Required:     true,
			},
		},
	}

	GeoBlockTemplate = &RuleTemplate{
		ID:          "geo-block-basic",
		Name:        "Geographic Blocking",
		Description: "Block traffic from specific countries",
		Category:    "Geographic Security",
		Type:        BlockingRuleTypeGeoLocation,
		Template: &BlockingRule{
			Name:        "Geo Block - {{.blocked_countries}}",
			Description: "Blocks traffic from specified countries",
			Type:        BlockingRuleTypeGeoLocation,
			Priority:    70,
			Conditions: []*RuleCondition{
				{
					Field:    "geo_location.country",
					Operator: OperatorInList,
					Value:    "{{.blocked_countries}}",
					ValueType: ValueTypeList,
					Weight:   1.0,
				},
			},
			Actions: []*RuleAction{
				{
					Type: ActionTypeBlock,
					Parameters: map[string]interface{}{
						"duration": 3600,
						"reason":   "Geographic restriction",
					},
				},
			},
		},
		Variables: []*TemplateVariable{
			{
				Name:        "blocked_countries",
				Type:        ValueTypeList,
				Description: "List of country codes to block",
				Required:    true,
			},
		},
	}
)

// GetPredefinedTemplates returns all predefined rule templates
func GetPredefinedTemplates() []*RuleTemplate {
	return []*RuleTemplate{
		SQLInjectionTemplate,
		XSSTemplate,
		RateLimitTemplate,
		GeoBlockTemplate,
	}
}

// GetTemplateByID returns a predefined template by ID
func GetTemplateByID(id string) *RuleTemplate {
	templates := GetPredefinedTemplates()
	for _, template := range templates {
		if template.ID == id {
			return template
		}
	}
	return nil
}
