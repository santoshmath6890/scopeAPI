# Blocking Rule Model

## Overview

The `blocking_rule.go` file defines the core data structures and models for the attack-blocking service's rule engine. This file contains comprehensive models for creating, managing, and executing blocking rules that protect APIs from various security threats.

## Responsibilities

### Core Data Models
- **BlockingRule**: Main structure representing a security blocking rule
- **RuleCondition**: Defines conditions that trigger rule execution
- **RuleAction**: Specifies actions to take when rules are triggered
- **RuleMatch**: Records when a rule matches against a request
- **RulePerformance**: Tracks rule effectiveness and performance metrics

### Rule Management
- **Rule Validation**: Ensures rules are properly configured and valid
- **Rule Templates**: Predefined templates for common security scenarios
- **Rule Cloning**: Deep copying of rules for modification and testing
- **Batch Operations**: Support for bulk rule operations

### Template System
- **RuleTemplate**: Template structure for creating reusable rule patterns
- **TemplateVariable**: Variables that can be customized in templates
- **Template Application**: Logic for applying variables to create rules

## Key Features

### Flexible Rule Engine
- **Multiple Rule Types**: Support for signature, anomaly, behavioral, reputation, rate limiting, geo-location, ML-based, and adaptive rules
- **Complex Conditions**: Rich set of operators for matching various data types
- **Weighted Conditions**: Confidence scoring based on condition weights
- **Multiple Actions**: Support for blocking, alerting, logging, rate limiting, and custom actions

### Advanced Operators
- **String Operations**: equals, contains, starts_with, ends_with, regex
- **Numeric Operations**: greater_than, less_than, ranges
- **List Operations**: in_list, not_in_list
- **Network Operations**: ip_in_range, geo_location
- **Existence Operations**: exists, not_exists

### Action Types
- **Blocking Actions**: Block requests with configurable duration
- **Rate Limiting**: Apply rate limits with custom windows
- **Alerting**: Generate security alerts with severity levels
- **Logging**: Create detailed security logs
- **Redirection**: Redirect malicious requests
- **Notification**: Send notifications to administrators
- **Rule Updates**: Dynamically update rules based on patterns

### Performance Tracking
- **Effectiveness Metrics**: Precision, recall, F1-score, accuracy
- **Performance Analytics**: Processing time, match rates
- **True/False Positives**: Track rule accuracy over time
- **Adaptive Scoring**: Dynamic effectiveness scoring

## Data Structures

### BlockingRule
```go
type BlockingRule struct {
    ID                 string           // Unique rule identifier
    Name               string           // Human-readable rule name
    Description        string           // Rule description
    Type               BlockingRuleType // Rule type (signature, anomaly, etc.)
    Status             RuleStatus       // Rule status (active, inactive, etc.)
    Priority           int              // Rule execution priority
    Conditions         []*RuleCondition // Rule conditions
    Actions            []*RuleAction    // Actions to execute
    Metadata           map[string]interface{} // Additional metadata
    CreatedBy          string           // Rule creator
    UpdatedBy          string           // Last modifier
    CreatedAt          time.Time        // Creation timestamp
    UpdatedAt          time.Time        // Last update timestamp
    LastTriggered      *time.Time       // Last trigger time
    TriggerCount       int64            // Total trigger count
    EffectivenessScore float64          // Rule effectiveness score
    Tags               []string         // Rule tags
    ExpiresAt          *time.Time       // Rule expiration time
    IsTemporary        bool             // Temporary rule flag
    Version            int              // Rule version
}


type RuleCondition struct {
    ID            string                 // Condition identifier
    Field         string                 // Field to evaluate
    Operator      ConditionOperator      // Comparison operator
    Value         interface{}            // Expected value
    ValueType     ValueType              // Value data type
    CaseSensitive bool                   // Case sensitivity flag
    Negate        bool                   // Negation flag
    Weight        float64                // Condition weight
    Description   string                 // Condition description
    Metadata      map[string]interface{} // Additional metadata
}

### RuleAction
```go
type RuleAction struct {
    ID          string                 // Action identifier
    Type        ActionType             // Action type (block, alert, etc.)
    Parameters  map[string]interface{} // Action parameters
    Priority    int                    // Action execution priority
    Description string                 // Action description
    Enabled     bool                   // Action enabled flag
    Metadata    map[string]interface{} // Additional metadata
}

type RuleMatch struct {
    ID                string            // Match identifier
    RuleID            string            // Matched rule ID
    RuleName          string            // Matched rule name
    RequestID         string            // Request identifier
    MatchedConditions []*ConditionMatch // Matched conditions
    Confidence        float64           // Match confidence score
    Severity          string            // Match severity
    Actions           []*ActionResult   // Executed actions
    Timestamp         time.Time         // Match timestamp
    ProcessingTime    time.Duration     // Processing duration
    Metadata          map[string]interface{} // Additional metadata
}

Predefined Templates
SQL Injection Template
•	Purpose: Detects SQL injection attempts in request data
•	Conditions: Regex patterns for SQL keywords and syntax
•	Actions: Block request and generate high-severity alert
•	Variables: target_field, severity level
XSS Template
•	Purpose: Detects cross-site scripting attempts
•	Conditions: Regex patterns for script tags and JavaScript
•	Actions: Block request with shorter duration
•	Variables: target_field, severity level
Rate Limiting Template
•	Purpose: Controls request rate per IP address
•	Conditions: IP address existence check
•	Actions: Apply rate limiting with configurable thresholds
•	Variables: requests_per_minute
Geographic Blocking Template
•	Purpose: Blocks traffic from specific countries
•	Conditions: Country code matching
•	Actions: Block with geographic restriction reason
•	Variables: blocked_countries list
Validation Rules
Rule Validation
•	Rule name must be provided
•	At least one condition is required
•	At least one action is required
•	All conditions must be valid
•	All actions must be valid
Condition Validation
•	Field name is required
•	Operator is required
•	Value is required (except for exists/not_exists operators)
•	Operator-specific validation (regex patterns, IP ranges, etc.)
Action Validation
•	Action type is required
•	Type-specific parameter validation
•	Required parameters for specific action types
Helper Methods
BlockingRule Methods
•	IsActive(): Check if rule is active
•	IsExpired(): Check if rule has expired
•	HasTag(tag): Check if rule has specific tag
•	AddTag(tag): Add tag to rule
•	RemoveTag(tag): Remove tag from rule
•	GetConditionByID(id): Get condition by ID
•	GetActionByID(id): Get action by ID
•	Validate(): Validate rule structure
•	Clone(): Create deep copy of rule
RuleCondition Methods
•	Validate(): Validate condition structure
•	Clone(): Create deep copy of condition
RuleAction Methods
•	Validate(): Validate action structure
•	Clone(): Create deep copy of action
RuleMatch Methods
•	GetTotalConfidence(): Calculate weighted confidence score
•	GetSuccessfulActions(): Get successfully executed actions
•	GetFailedActions(): Get failed action executions
RulePerformance Methods
•	CalculateMetrics(): Calculate precision, recall, F1-score, accuracy
•	IsEffective(): Determine if rule is effective based on metrics
RuleTemplate Methods
•	CreateRule(variables): Create rule from template with variables
•	Validate(): Validate template structure
•	applyVariables(): Apply template variables to rule
JSON Serialization
Custom Marshaling
•	Complex fields (conditions, actions, metadata, tags) are serialized as JSON strings
•	Proper handling of nested structures
•	Database-compatible serialization
Custom Unmarshaling
•	Parse JSON strings back to complex structures
•	Error handling for malformed data
•	Backward compatibility support
Database Integration
Scanning Support
•	Scan() method for database row scanning
•	Value() method for database value conversion
•	JSON-based storage for complex fields
Field Constants
Predefined constants for common request fields:
•	FieldSourceIP: Source IP address
•	FieldUserAgent: User agent string
•	FieldRequestPath: Request path
•	FieldRequestMethod: HTTP method
•	FieldRequestHeaders: Request headers
•	FieldRequestBody: Request body
•	FieldResponseCode: Response status code
•	FieldGeoLocation: Geographic location
•	FieldDomain: Domain name
•	And many more...


Usage Examples
Creating a Basic Blocking Rule

rule := &BlockingRule{
    Name:        "Block Malicious IPs",
    Description: "Blocks known malicious IP addresses",
    Type:        BlockingRuleTypeReputation,
    Status:      RuleStatusActive,
    Priority:    90,
    Conditions: []*RuleCondition{
        {
            Field:     FieldSourceIP,
            Operator:  OperatorInList,
            Value:     []string{"192.168.1.100", "10.0.0.50"},
            ValueType: ValueTypeList,
            Weight:    1.0,
        },
    },
    Actions: []*RuleAction{
        {
            Type: ActionTypeBlock,
            Parameters: map[string]interface{}{
                "duration": 3600,
                "reason":   "Malicious IP detected",
            },
        },
    },
}

Using Templates
template := GetTemplateByID("sql-injection-basic")
variables := map[string]interface{}{
    "target_field": "request_body",
    "severity":     "critical",
}
rule, err := template.CreateRule(variables)


Performance Tracking
performance := &RulePerformance{
    RuleID:         "rule-123",
    TruePositives:  95,
    FalsePositives: 5,
    TrueNegatives:  900,
    FalseNegatives: 10,
}
performance.CalculateMetrics()
// Results in calculated precision, recall, F1-score, accuracy


Integration Points
Repository Layer
Rule persistence and retrieval
Performance metrics storage
Template management
Service Layer
Rule evaluation engine
Action execution
Performance monitoring
Event Processing
Real-time rule matching
Action result tracking
Performance data collection
Security Considerations
Rule Security
Input validation for all rule components
Sanitization of regex patterns
Protection against rule injection attacks
Performance Security
Rate limiting for rule evaluation
Memory usage monitoring
CPU usage optimization
Data Security
Sensitive data masking in logs
Secure storage of rule parameters
Access control for rule management
This model provides a comprehensive foundation for the attack-blocking service's rule engine, supporting complex security scenarios while maintaining performance and flexibility.