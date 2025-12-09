# Agent Management Service

## Overview

The Agent Management Service is a core component of the ScopeAPI attack-blocking service responsible for managing the lifecycle, configuration, and monitoring of security agents deployed across the infrastructure. This service provides comprehensive agent management capabilities including registration, authentication, health monitoring, policy assignment, and command execution.

## Responsibilities

### Core Agent Management
- **Agent Registration**: Handle new agent registration and re-registration
- **Agent Authentication**: Manage agent sessions and authentication tokens
- **Agent Lifecycle**: Control agent status transitions and lifecycle events
- **Agent Configuration**: Manage and distribute agent configurations
- **Agent Updates**: Deploy and manage agent software updates

### Session Management
- **Session Creation**: Create and manage agent sessions
- **Session Authentication**: Validate agent session tokens
- **Session Refresh**: Handle session token renewal
- **Session Termination**: Gracefully terminate agent sessions
- **Session Monitoring**: Monitor session health and timeout handling

### Health Monitoring
- **Heartbeat Processing**: Process and validate agent heartbeats
- **Health Status Tracking**: Monitor agent health metrics and status
- **Performance Metrics**: Collect and analyze agent performance data
- **Alerting**: Generate alerts for unhealthy or inactive agents
- **Status Reporting**: Provide comprehensive agent health reports

### Policy Management
- **Policy Assignment**: Assign security policies to agents
- **Policy Distribution**: Distribute policy updates to agents
- **Policy Enforcement**: Ensure agents are running with correct policies
- **Policy Compliance**: Monitor policy compliance across agents

### Command and Control
- **Command Execution**: Send commands to agents for remote management
- **Command Tracking**: Track command execution status and results
- **Bulk Operations**: Execute operations across multiple agents
- **Emergency Controls**: Handle emergency shutdown and restart commands

### Group Management
- **Agent Grouping**: Organize agents into logical groups
- **Group Operations**: Perform bulk operations on agent groups
- **Group Policies**: Apply policies at the group level
- **Group Monitoring**: Monitor group-level metrics and health

### Backup and Recovery
- **Configuration Backup**: Create backups of agent configurations
- **Configuration Restore**: Restore agent configurations from backups
- **Disaster Recovery**: Handle agent recovery scenarios
- **State Persistence**: Maintain agent state across service restarts

## Key Features

### Real-time Monitoring
- Continuous agent health monitoring
- Real-time heartbeat processing
- Performance metrics collection
- Automated status updates

### Scalable Architecture
- In-memory caching for performance
- Asynchronous event processing
- Bulk operation support
- Horizontal scaling capabilities

### Security Features
- Secure agent authentication
- Session token management
- Encrypted communication
- Access control and authorization

### High Availability
- Graceful degradation
- Automatic failover
- Service health checks
- Maintenance mode support

## Service Dependencies

### Required Services
- **Repository Layer**: Agent data persistence
- **Kafka Producer**: Event publishing and messaging
- **Logger**: Structured logging and audit trails

### External Dependencies
- **PostgreSQL**: Agent metadata storage
- **Redis**: Session and metrics caching
- **Apache Kafka**: Event streaming and messaging

## Configuration

### Service Configuration
```go
type AgentManagementServiceConfig struct {
    SessionTimeout    time.Duration // Agent session timeout
    HeartbeatInterval time.Duration // Heartbeat check interval
    MaxInactiveTime   time.Duration // Maximum inactive time before marking unhealthy
}

Agent Configuration
•	Log levels and reporting intervals
•	Connection limits and buffer sizes
•	Metrics and tracing settings
•	Custom metadata support
API Operations
Agent Lifecycle
•	RegisterAgent() - Register new agent
•	UnregisterAgent() - Remove agent from system
•	GetAgent() - Retrieve agent information
•	GetAgents() - List agents with filtering
•	UpdateAgent() - Update agent metadata
Session Management
•	AuthenticateAgent() - Validate agent credentials
•	RefreshSession() - Renew session token
•	TerminateSession() - End agent session
Health and Monitoring
•	ProcessHeartbeat() - Handle agent heartbeats
•	GetAgentHealth() - Retrieve agent health status
•	GetAgentMetrics() - Get performance metrics
•	GetAgentStatistics() - System-wide statistics
Configuration Management
•	UpdateAgentConfiguration() - Update agent config
•	GetAgentConfiguration() - Retrieve current config
Policy Management
•	AssignPolicyToAgent() - Assign security policy
•	UnassignPolicyFromAgent() - Remove policy assignment
•	GetAgentPolicies() - List assigned policies
Command and Control
•	SendCommandToAgent() - Execute remote command
•	GetAgentCommands() - List command history
•	UpdateCommandStatus() - Update command execution status
Bulk Operations
•	BulkUpdateAgentConfiguration() - Update multiple agents
•	BulkAssignPolicy() - Assign policy to multiple agents
•	BulkSendCommand() - Send command to multiple agents
Group Management
•	CreateAgentGroup() - Create agent group
•	AddAgentToGroup() - Add agent to group
•	RemoveAgentFromGroup() - Remove agent from group
•	GetAgentGroups() - List agent groups
Backup and Recovery
•	BackupAgentConfiguration() - Create configuration backup
•	RestoreAgentConfiguration() - Restore from backup
Event Publishing
Agent Events
•	agent_registered - New agent registration
•	agent_unregistered - Agent removal
•	agent_updated - Agent information updated
•	agent_inactive - Agent became inactive
Heartbeat Events
•	agent_heartbeat - Regular heartbeat received
•	Health status changes
•	Performance metric updates
Error Handling
Validation Errors
•	Invalid registration requests
•	Malformed heartbeats
•	Invalid configurations
•	Command validation failures
Runtime Errors
•	Repository access failures
•	Kafka publishing errors
•	Session timeout handling
•	Agent communication failures
Performance Considerations
Optimization Strategies
•	In-memory caching for frequently accessed data
•	Asynchronous event processing
•	Bulk operation batching
•	Connection pooling
Scalability Features
•	Horizontal scaling support
•	Load balancing capabilities
•	Resource usage monitoring
•	Performance metrics collection
Monitoring and Observability
Health Checks
•	Service health validation
•	Repository connectivity
•	Kafka producer status
•	In-memory state verification
Metrics
•	Agent count and status distribution
•	Session activity metrics
•	Command execution statistics
•	Performance indicators
Logging
•	Structured logging with context
•	Audit trail for all operations
•	Error tracking and alerting
•	Debug information for troubleshooting
Security Considerations
Authentication
•	Secure session token generation
•	Token validation and expiration
•	Multi-factor authentication support
•	Certificate-based authentication
Authorization
•	Role-based access control
•	Operation-level permissions
•	Agent-specific access rights
•	Group-based authorization
Data Protection
•	Encrypted data transmission
•	Secure configuration storage
•	Sensitive data masking
•	Audit logging
Maintenance Operations
Routine Maintenance
•	Expired session cleanup
•	Old metrics purging
•	Command history cleanup
•	Log rotation
Backup Operations
•	Configuration backups
•	State snapshots
•	Recovery procedures
•	Data migration
Integration Points
Internal Services
•	Attack Blocking Service
•	Policy Enforcement Service
•	Cloud Intelligence Service
•	Analytics and Reporting
External Systems
•	Security Information and Event Management (SIEM)
•	Network monitoring tools
•	Configuration management systems
•	Alerting and notification systems
Usage Examples

*Basic Agent Registration
request := &models.AgentRegistrationRequest{
    AgentID:      "agent-001",
    Name:         "Web Server Agent",
    Hostname:     "web-server-01",
    IPAddress:    "192.168.1.100",
    Version:      "1.0.0",
    Type:         models.AgentTypeBlocking,
    Capabilities: []string{"http_blocking", "rate_limiting"},
}

response, err := service.RegisterAgent(ctx, request)

*Health Monitoring
health, err := service.GetAgentHealth(ctx, "agent-001")
if err != nil {
    log.Error("Failed to get agent health", err)
}

if health.OverallHealth != "healthy" {
    // Handle unhealthy agent
}

*Bulk Operations
agentIDs := []string{"agent-001", "agent-002", "agent-003"}
config := &models.AgentConfiguration{
    LogLevel:       "debug",
    ReportInterval: 30 * time.Second,
}

result, err := service.BulkUpdateAgentConfiguration(ctx, agentIDs, config)


Best Practices
Implementation Guidelines
•	Always validate input parameters
•	Use structured logging with context
•	Handle errors gracefully
•	Implement proper timeout handling
•	Use bulk operations for efficiency
Performance Tips
•	Cache frequently accessed data
•	Use asynchronous processing where possible
•	Implement proper connection pooling
•	Monitor resource usage
•	Optimize database queries
Security Best Practices
•	Validate all agent communications
•	Use secure session management
•	Implement proper access controls
•	Audit all operations
•	Encrypt sensitive data

