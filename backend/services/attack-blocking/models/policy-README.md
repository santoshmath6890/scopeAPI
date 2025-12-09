# policy.go Documentation

## Overview
This file is part of the ScopeAPI backend services architecture, specifically the attack-blocking service.

## Location
`backend/services/attack-blocking/internal/models/policy.go`

## Responsibility
policy.go is responsible for:

- **Policy Management**: Defining and managing security policies for attack blocking
- **Policy Engine**: Core engine for evaluating requests against security policies
- **Policy Evaluation**: Real-time evaluation of incoming requests against configured policies
- **Policy Types**: Supporting various policy types (blocking, rate limiting, validation, etc.)
- **Condition Evaluation**: Complex condition matching with multiple operators
- **Action Execution**: Defining and managing actions to take when policies match
- **Policy Metrics**: Tracking policy execution statistics and performance metrics
- **Policy Builder**: Fluent API for building policies programmatically
- **Policy Templates**: Pre-built policy templates for common security scenarios

## Key Components

### Core Structures
- `Policy`: Main policy structure with conditions, actions, and targets
- `PolicyEngine`: Central engine for managing and evaluating policies
- `PolicyCondition`: Individual conditions within a policy
- `PolicyAction`: Actions to execute when policy matches
- `PolicyTarget`: Targets that policies apply to
- `RequestContext`: Context of incoming requests being evaluated
- `PolicyDecision`: Result of policy evaluation

### Policy Types
- `PolicyTypeBlocking`: Policies that block malicious requests
- `PolicyTypeRateLimit`: Policies that limit request rates
- `PolicyTypeValidation`: Policies that validate request data
- `PolicyTypeTransform`: Policies that transform requests
- `PolicyTypeMonitoring`: Policies for monitoring and logging
- `PolicyTypeCompliance`: Policies for compliance enforcement

### Condition Operators
- Equality operators (equals, not_equals)
- String operators (contains, starts_with, ends_with)
- Comparison operators (greater_than, less_than, etc.)
- Collection operators (in, not_in)
- Pattern matching (regex)

### Action Types
- `ActionBlock`: Block the request
- `ActionAllow`: Explicitly allow the request
- `ActionRateLimit`: Apply rate limiting
- `ActionLog`: Log the event
- `ActionAlert`: Send alerts
- `ActionRedirect`: Redirect the request
- `ActionTransform`: Transform the request
- `ActionChallenge`: Challenge the user

## Key Features

### Policy Evaluation
- Real-time evaluation of requests against active policies
- Priority-based policy ordering
- Complex condition logic with AND/OR operators
- Target-based policy application
- Performance metrics tracking

### Policy Management
- CRUD operations for policies
- Policy versioning and history
- Policy import/export functionality
- Policy validation and testing
- Policy templates and builders

### Metrics and Monitoring
- Execution count tracking
- Match rate statistics
- Performance metrics (execution time)
- Policy effectiveness analysis
- Comprehensive reporting

### Security Features
- SQL injection detection templates
- Malicious IP blocking
- Rate limiting policies
- Custom security rules
- Compliance policy enforcement

## Usage Examples

### Creating a Basic Blocking Policy
```go
policy := NewPolicyBuilder().
    WithName("Block SQL Injection").
    WithType(PolicyTypeBlocking).
    AddCondition("body", OperatorRegex, sqlInjectionPattern).
    AddAction(ActionBlock, blockConfig).
    Build()


#Evaluating Requests
engine := NewPolicyEngine()
decision, err := engine.EvaluateRequest(requestContext)
if !decision.Allow {
    // Handle blocked request
}

#Policy Templates
// Use pre-built templates
sqlPolicy := BlockSQLInjectionPolicy()
rateLimitPolicy := RateLimitPolicy(100, "1m")
ipBlockPolicy := BlockMaliciousIPPolicy(maliciousIPs)


Dependencies
•	Standard Go libraries (crypto/sha256, encoding/json, fmt, sync, time)
•	No external dependencies for core functionality
Integration Points
•	Integrates with attack-blocking service handlers
•	Used by policy enforcement services
•	Connects to metrics and monitoring systems
•	Interfaces with logging and alerting systems
Performance Considerations
•	Thread-safe operations with mutex protection
•	Efficient policy evaluation with early termination
•	Metrics collection with minimal overhead
•	Memory-efficient policy storage and retrieval
Security Considerations
•	Input validation for all policy components
•	Safe evaluation of user-provided conditions
•	Protection against policy injection attacks
•	Secure handling of sensitive policy data


