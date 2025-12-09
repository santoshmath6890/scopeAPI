package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// PolicyType represents the type of policy
type PolicyType string

const (
	PolicyTypeBlocking    PolicyType = "blocking"
	PolicyTypeRateLimit   PolicyType = "rate_limit"
	PolicyTypeValidation  PolicyType = "validation"
	PolicyTypeTransform   PolicyType = "transform"
	PolicyTypeMonitoring  PolicyType = "monitoring"
	PolicyTypeCompliance  PolicyType = "compliance"
)

// PolicyStatus represents the status of a policy
type PolicyStatus string

const (
	PolicyStatusActive   PolicyStatus = "active"
	PolicyStatusInactive PolicyStatus = "inactive"
	PolicyStatusDraft    PolicyStatus = "draft"
	PolicyStatusArchived PolicyStatus = "archived"
)

// ConditionOperator represents condition operators
type ConditionOperator string

const (
	OperatorEquals         ConditionOperator = "equals"
	OperatorNotEquals      ConditionOperator = "not_equals"
	OperatorContains       ConditionOperator = "contains"
	OperatorNotContains    ConditionOperator = "not_contains"
	OperatorStartsWith     ConditionOperator = "starts_with"
	OperatorEndsWith       ConditionOperator = "ends_with"
	OperatorGreaterThan    ConditionOperator = "greater_than"
	OperatorLessThan       ConditionOperator = "less_than"
	OperatorGreaterOrEqual ConditionOperator = "greater_or_equal"
	OperatorLessOrEqual    ConditionOperator = "less_or_equal"
	OperatorRegex          ConditionOperator = "regex"
	OperatorIn             ConditionOperator = "in"
	OperatorNotIn          ConditionOperator = "not_in"
)

// PolicyActionType represents action types
type PolicyActionType string

const (
	ActionBlock       PolicyActionType = "block"
	ActionAllow       PolicyActionType = "allow"
	ActionRateLimit   PolicyActionType = "rate_limit"
	ActionLog         PolicyActionType = "log"
	ActionAlert       PolicyActionType = "alert"
	ActionRedirect    PolicyActionType = "redirect"
	ActionTransform   PolicyActionType = "transform"
	ActionChallenge   PolicyActionType = "challenge"
)

// PolicyTargetType represents target types
type PolicyTargetType string

const (
	TargetAPI        PolicyTargetType = "api"
	TargetEndpoint   PolicyTargetType = "endpoint"
	TargetUser       PolicyTargetType = "user"
	TargetIP         PolicyTargetType = "ip"
	TargetUserAgent  PolicyTargetType = "user_agent"
	TargetHeader     PolicyTargetType = "header"
	TargetParameter  PolicyTargetType = "parameter"
	TargetBody       PolicyTargetType = "body"
)

// Policy represents a security policy
type Policy struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        PolicyType        `json:"type"`
	Priority    int               `json:"priority"`
	Status      PolicyStatus      `json:"status"`
	Conditions  []*PolicyCondition `json:"conditions"`
	Actions     []*PolicyAction   `json:"actions"`
	Targets     []*PolicyTarget   `json:"targets"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedBy   string            `json:"created_by"`
	Version     int               `json:"version"`
	IsEnabled   bool              `json:"is_enabled"`
}

// PolicyCondition represents a condition in a policy
type PolicyCondition struct {
	Field    string            `json:"field"`
	Operator ConditionOperator `json:"operator"`
	Value    interface{}       `json:"value"`
	LogicOp  string           `json:"logic_op,omitempty"` // AND, OR
}

// PolicyAction represents an action to take when policy matches
type PolicyAction struct {
	Type   PolicyActionType       `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// PolicyTarget represents what the policy applies to
type PolicyTarget struct {
	Type  PolicyTargetType `json:"type"`
	Value string          `json:"value"`
}

// PolicyEngine manages policies and their execution
type PolicyEngine struct {
	policies map[string]*Policy
	mutex    sync.RWMutex
	metrics  *PolicyMetrics
}

// PolicyMetrics tracks policy execution metrics
type PolicyMetrics struct {
	TotalPolicies    int64                    `json:"total_policies"`
	ActivePolicies   int64                    `json:"active_policies"`
	ExecutionCount   map[string]int64         `json:"execution_count"`
	MatchCount       map[string]int64         `json:"match_count"`
	LastExecution    map[string]time.Time     `json:"last_execution"`
	AverageExecTime  map[string]time.Duration `json:"average_exec_time"`
	mutex            sync.RWMutex
}

// NewPolicyEngine creates a new policy engine
func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		policies: make(map[string]*Policy),
		metrics:  NewPolicyMetrics(),
	}
}

// NewPolicyMetrics creates new policy metrics
func NewPolicyMetrics() *PolicyMetrics {
	return &PolicyMetrics{
		ExecutionCount:  make(map[string]int64),
		MatchCount:      make(map[string]int64),
		LastExecution:   make(map[string]time.Time),
		AverageExecTime: make(map[string]time.Duration),
	}
}

// AddPolicy adds a new policy
func (pe *PolicyEngine) AddPolicy(policy *Policy) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	if policy.ID == "" {
		policy.ID = pe.generatePolicyID()
	}
	
	if err := pe.validatePolicy(policy); err != nil {
		return fmt.Errorf("policy validation failed: %w", err)
	}
	
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()
	policy.Version = 1
	policy.IsEnabled = policy.Status == PolicyStatusActive
	
	pe.policies[policy.ID] = policy
	pe.updateMetrics()
	
	return nil
}

// UpdatePolicy updates an existing policy
func (pe *PolicyEngine) UpdatePolicy(policy *Policy) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	existing, exists := pe.policies[policy.ID]
	if !exists {
		return fmt.Errorf("policy not found: %s", policy.ID)
	}
	
	if err := pe.validatePolicy(policy); err != nil {
		return fmt.Errorf("policy validation failed: %w", err)
	}
	
	policy.CreatedAt = existing.CreatedAt
	policy.UpdatedAt = time.Now()
	policy.Version = existing.Version + 1
	policy.IsEnabled = policy.Status == PolicyStatusActive
	
	pe.policies[policy.ID] = policy
	pe.updateMetrics()
	
	return nil
}

// DeletePolicy removes a policy
func (pe *PolicyEngine) DeletePolicy(policyID string) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	if _, exists := pe.policies[policyID]; !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}
	
	delete(pe.policies, policyID)
	pe.updateMetrics()
	
	return nil
}

// GetPolicy retrieves a policy by ID
func (pe *PolicyEngine) GetPolicy(policyID string) (*Policy, error) {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policy, exists := pe.policies[policyID]
	if !exists {
		return nil, fmt.Errorf("policy not found: %s", policyID)
	}
	
	return policy, nil
}

// ListPolicies returns all policies
func (pe *PolicyEngine) ListPolicies() []*Policy {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policies := make([]*Policy, 0, len(pe.policies))
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}
	
	return policies
}

// EvaluateRequest evaluates a request against all active policies
func (pe *PolicyEngine) EvaluateRequest(request *RequestContext) (*PolicyDecision, error) {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	decision := &PolicyDecision{
		Allow:          true,
		MatchedPolicies: make([]*Policy, 0),
		Actions:        make([]*PolicyAction, 0),
		Timestamp:      time.Now(),
	}
	
	// Sort policies by priority (higher priority first)
	sortedPolicies := pe.getSortedPolicies()
	
	for _, policy := range sortedPolicies {
		if !policy.IsEnabled {
			continue
		}
		
		start := time.Now()
		matches, err := pe.evaluatePolicy(policy, request)
		execTime := time.Since(start)
		
		// Update metrics
		pe.updatePolicyMetrics(policy.ID, execTime, matches)
		
		if err != nil {
			continue // Log error but continue evaluation
		}
		
		if matches {
			decision.MatchedPolicies = append(decision.MatchedPolicies, policy)
			decision.Actions = append(decision.Actions, policy.Actions...)
			
			// Check if any action blocks the request
			for _, action := range policy.Actions {
				if action.Type == ActionBlock {
					decision.Allow = false
					decision.BlockReason = fmt.Sprintf("Blocked by policy: %s", policy.Name)
					return decision, nil
				}
			}
		}
	}
	
	return decision, nil
}

// RequestContext represents the context of a request being evaluated
type RequestContext struct {
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Headers     map[string]string      `json:"headers"`
	Parameters  map[string]string      `json:"parameters"`
	Body        string                 `json:"body"`
	ClientIP    string                 `json:"client_ip"`
	UserAgent   string                 `json:"user_agent"`
	UserID      string                 `json:"user_id,omitempty"`
	APIKey      string                 `json:"api_key,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PolicyDecision represents the result of policy evaluation
type PolicyDecision struct {
	Allow           bool            `json:"allow"`
	BlockReason     string          `json:"block_reason,omitempty"`
	MatchedPolicies []*Policy       `json:"matched_policies"`
	Actions         []*PolicyAction `json:"actions"`
	Timestamp       time.Time       `json:"timestamp"`
}

// Helper methods

func (pe *PolicyEngine) generatePolicyID() string {
	timestamp := time.Now().Unix()
	hash := sha256.Sum256([]byte(fmt.Sprintf("policy_%d", timestamp)))
	return fmt.Sprintf("policy_%x", hash[:8])
}

func (pe *PolicyEngine) validatePolicy(policy *Policy) error {
	if policy.Name == "" {
		return fmt.Errorf("policy name is required")
	}
	
	if policy.Type == "" {
		return fmt.Errorf("policy type is required")
	}
	
	if policy.Priority < 0 || policy.Priority > 100 {
		return fmt.Errorf("policy priority must be between 0 and 100")
	}
	
	if len(policy.Conditions) == 0 {
		return fmt.Errorf("policy must have at least one condition")
	}
	
	if len(policy.Actions) == 0 {
		return fmt.Errorf("policy must have at least one action")
	}
	
	// Validate conditions
	for i, condition := range policy.Conditions {
		if condition.Field == "" {
			return fmt.Errorf("condition %d: field is required", i)
		}
		if condition.Operator == "" {
			return fmt.Errorf("condition %d: operator is required", i)
		}
		if condition.Value == nil {
			return fmt.Errorf("condition %d: value is required", i)
		}
	}
	
	// Validate actions
	for i, action := range policy.Actions {
		if action.Type == "" {
			return fmt.Errorf("action %d: type is required", i)
		}
	}
	
	return nil
}

func (pe *PolicyEngine) getSortedPolicies() []*Policy {
	policies := make([]*Policy, 0, len(pe.policies))
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}
	
	// Sort by priority (descending)
	for i := 0; i < len(policies)-1; i++ {
		for j := i + 1; j < len(policies); j++ {
			if policies[i].Priority < policies[j].Priority {
				policies[i], policies[j] = policies[j], policies[i]
			}
		}
	}
	
	return policies
}

func (pe *PolicyEngine) evaluatePolicy(policy *Policy, request *RequestContext) (bool, error) {
	// Check if policy targets match the request
	if !pe.evaluateTargets(policy.Targets, request) {
		return false, nil
	}
	
	// Evaluate all conditions
	return pe.evaluateConditions(policy.Conditions, request), nil
}

func (pe *PolicyEngine) evaluateTargets(targets []*PolicyTarget, request *RequestContext) bool {
	if len(targets) == 0 {
		return true // No targets means apply to all
	}
	
	for _, target := range targets {
		if pe.evaluateTarget(target, request) {
			return true
		}
	}
	
	return false
}

func (pe *PolicyEngine) evaluateTarget(target *PolicyTarget, request *RequestContext) bool {
	switch target.Type {
	case TargetAPI:
		return request.Path == target.Value
	case TargetEndpoint:
		return request.Path == target.Value
	case TargetIP:
		return request.ClientIP == target.Value
	case TargetUserAgent:
		return request.UserAgent == target.Value
	case TargetUser:
		return request.UserID == target.Value
	default:
		return false
	}
}

func (pe *PolicyEngine) evaluateConditions(conditions []*PolicyCondition, request *RequestContext) bool {
	if len(conditions) == 0 {
		return true
	}
	
	result := true
	for i, condition := range conditions {
		conditionResult := pe.evaluateCondition(condition, request)
		
		if i == 0 {
			result = conditionResult
		} else {
			switch condition.LogicOp {
			case "OR":
				result = result || conditionResult
			case "AND", "":
				result = result && conditionResult
			}
		}
	}
	
	return result
}

func (pe *PolicyEngine) evaluateCondition(condition *PolicyCondition, request *RequestContext) bool {
	fieldValue := pe.getFieldValue(condition.Field, request)
	if fieldValue == nil {
		return false
	}
	
	switch condition.Operator {
	case OperatorEquals:
		return pe.compareValues(fieldValue, condition.Value, "equals")
	case OperatorNotEquals:
		return !pe.compareValues(fieldValue, condition.Value, "equals")
	case OperatorContains:
		return pe.compareValues(fieldValue, condition.Value, "contains")
	case OperatorNotContains:
		return !pe.compareValues(fieldValue, condition.Value, "contains")
	case OperatorStartsWith:
		return pe.compareValues(fieldValue, condition.Value, "starts_with")
	case OperatorEndsWith:
		return pe.compareValues(fieldValue, condition.Value, "ends_with")
	case OperatorGreaterThan:
		return pe.compareValues(fieldValue, condition.Value, "greater_than")
	case OperatorLessThan:
		return pe.compareValues(fieldValue, condition.Value, "less_than")
	case OperatorGreaterOrEqual:
		return pe.compareValues(fieldValue, condition.Value, "greater_or_equal")
	case OperatorLessOrEqual:
		return pe.compareValues(fieldValue, condition.Value, "less_or_equal")
	case OperatorIn:
		return pe.compareValues(fieldValue, condition.Value, "in")
	case OperatorNotIn:
		return !pe.compareValues(fieldValue, condition.Value, "in")
	default:
		return false
	}
}

func (pe *PolicyEngine) getFieldValue(field string, request *RequestContext) interface{} {
	switch field {
	case "method":
		return request.Method
	case "path":
		return request.Path
	case "client_ip":
		return request.ClientIP
	case "user_agent":
		return request.UserAgent
	case "user_id":
		return request.UserID
	case "api_key":
		return request.APIKey
	default:
		// Check headers
		if headerValue, exists := request.Headers[field]; exists {
			return headerValue
		}
		// Check parameters
		if paramValue, exists := request.Parameters[field]; exists {
			return paramValue
		}
		// Check metadata
		if metaValue, exists := request.Metadata[field]; exists {
			return metaValue
		}
		return nil
	}
}

func (pe *PolicyEngine) compareValues(fieldValue, conditionValue interface{}, operator string) bool {
	switch operator {
	case "equals":
		return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", conditionValue)
	case "contains":
		fieldStr := fmt.Sprintf("%v", fieldValue)
		conditionStr := fmt.Sprintf("%v", conditionValue)
		return len(fieldStr) > 0 && len(conditionStr) > 0 && 
			   fieldStr != conditionStr && 
			   (fieldStr == conditionStr || len(fieldStr) > len(conditionStr))
	case "starts_with":
		fieldStr := fmt.Sprintf("%v", fieldValue)
		conditionStr := fmt.Sprintf("%v", conditionValue)
		return len(fieldStr) >= len(conditionStr) && fieldStr[:len(conditionStr)] == conditionStr
	case "ends_with":
		fieldStr := fmt.Sprintf("%v", fieldValue)
		conditionStr := fmt.Sprintf("%v", conditionValue)
		return len(fieldStr) >= len(conditionStr) && fieldStr[len(fieldStr)-len(conditionStr):] == conditionStr
	case "in":
		if values, ok := conditionValue.([]interface{}); ok {
			for _, value := range values {
				if fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", value) {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}

func (pe *PolicyEngine) updateMetrics() {
	pe.metrics.mutex.Lock()
	defer pe.metrics.mutex.Unlock()
	
	pe.metrics.TotalPolicies = int64(len(pe.policies))
	
	activeCount := int64(0)
	for _, policy := range pe.policies {
		if policy.IsEnabled {
			activeCount++
		}
	}
	pe.metrics.ActivePolicies = activeCount
}

func (pe *PolicyEngine) updatePolicyMetrics(policyID string, execTime time.Duration, matched bool) {
	pe.metrics.mutex.Lock()
	defer pe.metrics.mutex.Unlock()
	
	pe.metrics.ExecutionCount[policyID]++
	pe.metrics.LastExecution[policyID] = time.Now()
	
	if matched {
		pe.metrics.MatchCount[policyID]++
	}
	
	// Update average execution time
	currentAvg := pe.metrics.AverageExecTime[policyID]
	count := pe.metrics.ExecutionCount[policyID]
	newAvg := time.Duration((int64(currentAvg)*(count-1) + int64(execTime)) / count)
	pe.metrics.AverageExecTime[policyID] = newAvg
}

// GetMetrics returns current policy metrics
func (pe *PolicyEngine) GetMetrics() *PolicyMetrics {
	pe.metrics.mutex.RLock()
	defer pe.metrics.mutex.RUnlock()
	
	// Create a copy to avoid race conditions
	metrics := &PolicyMetrics{
		TotalPolicies:   pe.metrics.TotalPolicies,
		ActivePolicies:  pe.metrics.ActivePolicies,
		ExecutionCount:  make(map[string]int64),
		MatchCount:      make(map[string]int64),
		LastExecution:   make(map[string]time.Time),
		AverageExecTime: make(map[string]time.Duration),
	}
	
	for k, v := range pe.metrics.ExecutionCount {
		metrics.ExecutionCount[k] = v
	}
	for k, v := range pe.metrics.MatchCount {
		metrics.MatchCount[k] = v
	}
	for k, v := range pe.metrics.LastExecution {
		metrics.LastExecution[k] = v
	}
	for k, v := range pe.metrics.AverageExecTime {
		metrics.AverageExecTime[k] = v
	}
	
	return metrics
}

// EnablePolicy enables a policy
func (pe *PolicyEngine) EnablePolicy(policyID string) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	policy, exists := pe.policies[policyID]
	if !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}
	
	policy.Status = PolicyStatusActive
	policy.IsEnabled = true
	policy.UpdatedAt = time.Now()
	
	pe.updateMetrics()
	return nil
}

// DisablePolicy disables a policy
func (pe *PolicyEngine) DisablePolicy(policyID string) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	policy, exists := pe.policies[policyID]
	if !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}
	
	policy.Status = PolicyStatusInactive
	policy.IsEnabled = false
	policy.UpdatedAt = time.Now()
	
	pe.updateMetrics()
	return nil
}

// GetPoliciesByType returns policies filtered by type
func (pe *PolicyEngine) GetPoliciesByType(policyType PolicyType) []*Policy {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policies := make([]*Policy, 0)
	for _, policy := range pe.policies {
		if policy.Type == policyType {
			policies = append(policies, policy)
		}
	}
	
	return policies
}

// GetPoliciesByStatus returns policies filtered by status
func (pe *PolicyEngine) GetPoliciesByStatus(status PolicyStatus) []*Policy {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policies := make([]*Policy, 0)
	for _, policy := range pe.policies {
		if policy.Status == status {
			policies = append(policies, policy)
		}
	}
	
	return policies
}

// SearchPolicies searches policies by name or description
func (pe *PolicyEngine) SearchPolicies(query string) []*Policy {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policies := make([]*Policy, 0)
	for _, policy := range pe.policies {
		if pe.containsIgnoreCase(policy.Name, query) || 
		   pe.containsIgnoreCase(policy.Description, query) {
			policies = append(policies, policy)
		}
	}
	
	return policies
}

func (pe *PolicyEngine) containsIgnoreCase(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 && 
		   fmt.Sprintf("%s", str) != fmt.Sprintf("%s", substr)
}

// ExportPolicies exports policies to JSON
func (pe *PolicyEngine) ExportPolicies() ([]byte, error) {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	policies := make([]*Policy, 0, len(pe.policies))
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}
	
	return json.MarshalIndent(policies, "", "  ")
}

// ImportPolicies imports policies from JSON
func (pe *PolicyEngine) ImportPolicies(data []byte, overwrite bool) error {
	var policies []*Policy
	if err := json.Unmarshal(data, &policies); err != nil {
		return fmt.Errorf("failed to unmarshal policies: %w", err)
	}
	
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	for _, policy := range policies {
		if err := pe.validatePolicy(policy); err != nil {
			return fmt.Errorf("invalid policy %s: %w", policy.ID, err)
		}
		
		if _, exists := pe.policies[policy.ID]; exists && !overwrite {
			continue // Skip existing policies if overwrite is false
		}
		
		policy.UpdatedAt = time.Now()
		if _, exists := pe.policies[policy.ID]; !exists {
			policy.CreatedAt = time.Now()
			policy.Version = 1
		} else {
			policy.Version++
		}
		
		pe.policies[policy.ID] = policy
	}
	
	pe.updateMetrics()
	return nil
}

// PolicyStats represents policy statistics
type PolicyStats struct {
	TotalPolicies     int64                    `json:"total_policies"`
	ActivePolicies    int64                    `json:"active_policies"`
	InactivePolicies  int64                    `json:"inactive_policies"`
	DraftPolicies     int64                    `json:"draft_policies"`
	ArchivedPolicies  int64                    `json:"archived_policies"`
	PoliciesByType    map[PolicyType]int64     `json:"policies_by_type"`
	TopExecuted       []PolicyExecutionStat    `json:"top_executed"`
	TopMatched        []PolicyMatchStat        `json:"top_matched"`
}

type PolicyExecutionStat struct {
	PolicyID    string `json:"policy_id"`
	PolicyName  string `json:"policy_name"`
	Executions  int64  `json:"executions"`
}

type PolicyMatchStat struct {
	PolicyID    string `json:"policy_id"`
	PolicyName  string `json:"policy_name"`
	Matches     int64  `json:"matches"`
	MatchRate   float64 `json:"match_rate"`
}

// GetStats returns comprehensive policy statistics
func (pe *PolicyEngine) GetStats() *PolicyStats {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	stats := &PolicyStats{
		PoliciesByType: make(map[PolicyType]int64),
		TopExecuted:    make([]PolicyExecutionStat, 0),
		TopMatched:     make([]PolicyMatchStat, 0),
	}
	
	// Count policies by status and type
	for _, policy := range pe.policies {
		stats.TotalPolicies++
		
		switch policy.Status {
		case PolicyStatusActive:
			stats.ActivePolicies++
		case PolicyStatusInactive:
			stats.InactivePolicies++
		case PolicyStatusDraft:
			stats.DraftPolicies++
		case PolicyStatusArchived:
			stats.ArchivedPolicies++
		}
		
		stats.PoliciesByType[policy.Type]++
	}
	
	// Get top executed and matched policies
	pe.metrics.mutex.RLock()
	defer pe.metrics.mutex.RUnlock()
	
	for policyID, execCount := range pe.metrics.ExecutionCount {
		if policy, exists := pe.policies[policyID]; exists {
			stats.TopExecuted = append(stats.TopExecuted, PolicyExecutionStat{
				PolicyID:   policyID,
				PolicyName: policy.Name,
				Executions: execCount,
			})
		}
	}
	
	for policyID, matchCount := range pe.metrics.MatchCount {
		if policy, exists := pe.policies[policyID]; exists {
			execCount := pe.metrics.ExecutionCount[policyID]
			matchRate := float64(0)
			if execCount > 0 {
				matchRate = float64(matchCount) / float64(execCount) * 100
			}
			
			stats.TopMatched = append(stats.TopMatched, PolicyMatchStat{
				PolicyID:   policyID,
				PolicyName: policy.Name,
				Matches:    matchCount,
				MatchRate:  matchRate,
			})
		}
	}
	
	return stats
}

// ResetMetrics resets all policy metrics
func (pe *PolicyEngine) ResetMetrics() {
	pe.metrics.mutex.Lock()
	defer pe.metrics.mutex.Unlock()
	
	pe.metrics.ExecutionCount = make(map[string]int64)
	pe.metrics.MatchCount = make(map[string]int64)
	pe.metrics.LastExecution = make(map[string]time.Time)
	pe.metrics.AverageExecTime = make(map[string]time.Duration)
}

// ValidateRequest validates a request context
func ValidateRequestContext(request *RequestContext) error {
	if request == nil {
		return fmt.Errorf("request context cannot be nil")
	}
	
	if request.Method == "" {
		return fmt.Errorf("request method is required")
	}
	
	if request.Path == "" {
		return fmt.Errorf("request path is required")
	}
	
	if request.ClientIP == "" {
		return fmt.Errorf("client IP is required")
	}
	
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}
	
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	
	if request.Parameters == nil {
		request.Parameters = make(map[string]string)
	}
	
	if request.Metadata == nil {
		request.Metadata = make(map[string]interface{})
	}
	
	return nil
}

// Clone creates a deep copy of a policy
func (p *Policy) Clone() *Policy {
	clone := &Policy{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Type:        p.Type,
		Priority:    p.Priority,
		Status:      p.Status,
		Tags:        make([]string, len(p.Tags)),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		CreatedBy:   p.CreatedBy,
		Version:     p.Version,
		IsEnabled:   p.IsEnabled,
	}
	
	// Copy tags
	copy(clone.Tags, p.Tags)
	
	// Copy metadata
	for k, v := range p.Metadata {
		clone.Metadata[k] = v
	}
	
	// Copy conditions
	clone.Conditions = make([]*PolicyCondition, len(p.Conditions))
	for i, condition := range p.Conditions {
		clone.Conditions[i] = &PolicyCondition{
			Field:    condition.Field,
			Operator: condition.Operator,
			Value:    condition.Value,
			LogicOp:  condition.LogicOp,
		}
	}
	
	// Copy actions
	clone.Actions = make([]*PolicyAction, len(p.Actions))
	for i, action := range p.Actions {
		clone.Actions[i] = &PolicyAction{
			Type:   action.Type,
			Config: make(map[string]interface{}),
		}
		for k, v := range action.Config {
			clone.Actions[i].Config[k] = v
		}
	}
	
	// Copy targets
	clone.Targets = make([]*PolicyTarget, len(p.Targets))
	for i, target := range p.Targets {
		clone.Targets[i] = &PolicyTarget{
			Type:  target.Type,
			Value: target.Value,
		}
	}
	
	return clone
}

// String returns a string representation of the policy
func (p *Policy) String() string {
	return fmt.Sprintf("Policy{ID: %s, Name: %s, Type: %s, Priority: %d, Status: %s}", 
		p.ID, p.Name, p.Type, p.Priority, p.Status)
}

// IsActive returns true if the policy is active
func (p *Policy) IsActive() bool {
	return p.Status == PolicyStatusActive && p.IsEnabled
}

// HasTag returns true if the policy has the specified tag
func (p *Policy) HasTag(tag string) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the policy
func (p *Policy) AddTag(tag string) {
	if !p.HasTag(tag) {
		p.Tags = append(p.Tags, tag)
	}
}

// RemoveTag removes a tag from the policy
func (p *Policy) RemoveTag(tag string) {
	for i, t := range p.Tags {
		if t == tag {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			break
		}
	}
}

// GetMetadataValue gets a metadata value by key
func (p *Policy) GetMetadataValue(key string) (interface{}, bool) {
	value, exists := p.Metadata[key]
	return value, exists
}

// SetMetadataValue sets a metadata value
func (p *Policy) SetMetadataValue(key string, value interface{}) {
	if p.Metadata == nil {
		p.Metadata = make(map[string]interface{})
	}
	p.Metadata[key] = value
}

// PolicyBuilder helps build policies fluently
type PolicyBuilder struct {
	policy *Policy
}

// NewPolicyBuilder creates a new policy builder
func NewPolicyBuilder() *PolicyBuilder {
	return &PolicyBuilder{
		policy: &Policy{
			Conditions: make([]*PolicyCondition, 0),
			Actions:    make([]*PolicyAction, 0),
			Targets:    make([]*PolicyTarget, 0),
			Tags:       make([]string, 0),
			Metadata:   make(map[string]interface{}),
			Status:     PolicyStatusDraft,
			Priority:   50,
		},
	}
}

// WithName sets the policy name
func (pb *PolicyBuilder) WithName(name string) *PolicyBuilder {
	pb.policy.Name = name
	return pb
}

// WithDescription sets the policy description
func (pb *PolicyBuilder) WithDescription(description string) *PolicyBuilder {
	pb.policy.Description = description
	return pb
}

// WithType sets the policy type
func (pb *PolicyBuilder) WithType(policyType PolicyType) *PolicyBuilder {
	pb.policy.Type = policyType
	return pb
}

// WithPriority sets the policy priority
func (pb *PolicyBuilder) WithPriority(priority int) *PolicyBuilder {
	pb.policy.Priority = priority
	return pb
}

// WithStatus sets the policy status
func (pb *PolicyBuilder) WithStatus(status PolicyStatus) *PolicyBuilder {
	pb.policy.Status = status
	return pb
}

// AddCondition adds a condition to the policy
func (pb *PolicyBuilder) AddCondition(field string, operator ConditionOperator, value interface{}) *PolicyBuilder {
	condition := &PolicyCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
	pb.policy.Conditions = append(pb.policy.Conditions, condition)
	return pb
}

// AddAction adds an action to the policy
func (pb *PolicyBuilder) AddAction(actionType PolicyActionType, config map[string]interface{}) *PolicyBuilder {
	action := &PolicyAction{
		Type:   actionType,
		Config: config,
	}
	pb.policy.Actions = append(pb.policy.Actions, action)
	return pb
}

// AddTarget adds a target to the policy
func (pb *PolicyBuilder) AddTarget(targetType PolicyTargetType, value string) *PolicyBuilder {
	target := &PolicyTarget{
		Type:  targetType,
		Value: value,
	}
	pb.policy.Targets = append(pb.policy.Targets, target)
	return pb
}

// AddTag adds a tag to the policy
func (pb *PolicyBuilder) AddTag(tag string) *PolicyBuilder {
	pb.policy.Tags = append(pb.policy.Tags, tag)
	return pb
}

// WithMetadata sets metadata for the policy
func (pb *PolicyBuilder) WithMetadata(key string, value interface{}) *PolicyBuilder {
	pb.policy.Metadata[key] = value
	return pb
}

// Build builds and returns the policy
func (pb *PolicyBuilder) Build() *Policy {
	pb.policy.CreatedAt = time.Now()
	pb.policy.UpdatedAt = time.Now()
	pb.policy.Version = 1
	pb.policy.IsEnabled = pb.policy.Status == PolicyStatusActive
	
	if pb.policy.ID == "" {
		pb.policy.ID = pb.generateID()
	}
	
	return pb.policy
}

func (pb *PolicyBuilder) generateID() string {
	timestamp := time.Now().Unix()
	hash := sha256.Sum256([]byte(fmt.Sprintf("policy_%s_%d", pb.policy.Name, timestamp)))
	return fmt.Sprintf("policy_%x", hash[:8])
}

// Common policy templates
var (
	// BlockSQLInjectionPolicy creates a policy to block SQL injection attempts
	BlockSQLInjectionPolicy = func() *Policy {
		return NewPolicyBuilder().
			WithName("Block SQL Injection").
			WithDescription("Blocks requests containing SQL injection patterns").
			WithType(PolicyTypeBlocking).
			WithPriority(90).
			WithStatus(PolicyStatusActive).
			AddCondition("body", OperatorRegex, `(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`).
			AddCondition("parameters", OperatorRegex, `(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`).
			AddAction(ActionBlock, map[string]interface{}{
				"message": "SQL injection attempt detected",
				"code":    403,
			}).
			AddAction(ActionLog, map[string]interface{}{
				"level": "high",
				"category": "security",
			}).
			AddTag("security").
			AddTag("sql-injection").
			Build()
	}
	
	// RateLimitPolicy creates a basic rate limiting policy
	RateLimitPolicy = func(limit int, window string) *Policy {
		return NewPolicyBuilder().
			WithName("Rate Limit").
			WithDescription(fmt.Sprintf("Limits requests to %d per %s", limit, window)).
			WithType(PolicyTypeRateLimit).
			WithPriority(70).
			WithStatus(PolicyStatusActive).
			AddAction(ActionRateLimit, map[string]interface{}{
				"limit":  limit,
				"window": window,
			}).
			AddTag("rate-limit").
			Build()
	}
	
	// BlockMaliciousIPPolicy creates a policy to block known malicious IPs
	BlockMaliciousIPPolicy = func(ips []string) *Policy {
		builder := NewPolicyBuilder().
			WithName("Block Malicious IPs").
			WithDescription("Blocks requests from known malicious IP addresses").
			WithType(PolicyTypeBlocking).
			WithPriority(95).
			WithStatus(PolicyStatusActive).
			AddCondition("client_ip", OperatorIn, ips).
			AddAction(ActionBlock, map[string]interface{}{
				"message": "Request from blocked IP address",
				"code":    403,
			}).
			AddAction(ActionLog, map[string]interface{}{
				"level": "high",
				"category": "security",
			}).
			AddTag("security").
			AddTag("ip-blocking")
		
		return builder.Build()
	}
)
