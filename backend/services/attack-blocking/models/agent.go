package models

import (
	"time"
)

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	AgentStatusOnline  AgentStatus = "online"
	AgentStatusOffline AgentStatus = "offline"
	AgentStatusError   AgentStatus = "error"
)

// AgentType represents the type of agent
type AgentType string

const (
	AgentTypeBlocking   AgentType = "blocking"
	AgentTypeDiscovery  AgentType = "discovery"
	AgentTypeMonitoring AgentType = "monitoring"
)

// Agent represents a blocking agent in the system
type Agent struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Type        AgentType              `json:"type" db:"type"`
	Status      AgentStatus            `json:"status" db:"status"`
	IPAddress   string                 `json:"ip_address" db:"ip_address"`
	Port        int                    `json:"port" db:"port"`
	Version     string                 `json:"version" db:"version"`
	LastSeen    time.Time              `json:"last_seen" db:"last_seen"`
	Config      map[string]interface{} `json:"config" db:"config"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// AgentConfig represents agent configuration
type AgentConfig struct {
	ServerURL         string        `json:"server_url"`
	APIKey            string        `json:"api_key"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	ReportingInterval time.Duration `json:"reporting_interval"`
	LogLevel          string        `json:"log_level"`
	EnableMetrics     bool          `json:"enable_metrics"`
	EnableLogging     bool          `json:"enable_logging"`
	MaxConnections    int           `json:"max_connections"`
	BufferSize        int           `json:"buffer_size"`
	RetryAttempts     int           `json:"retry_attempts"`
	TimeoutDuration   time.Duration `json:"timeout_duration"`
}

// IsOnline returns true if the agent is online
func (a *Agent) IsOnline() bool {
	return a.Status == AgentStatusOnline
}

// IsHealthy returns true if the agent is healthy (online and recently seen)
func (a *Agent) IsHealthy() bool {
	if !a.IsOnline() {
		return false
	}
	
	// Consider agent unhealthy if not seen in the last 5 minutes
	return time.Since(a.LastSeen) < 5*time.Minute
}

// UpdateLastSeen updates the last seen timestamp
func (a *Agent) UpdateLastSeen() {
	a.LastSeen = time.Now()
	a.UpdatedAt = time.Now()
}

// SetStatus updates the agent status
func (a *Agent) SetStatus(status AgentStatus) {
	a.Status = status
	a.UpdatedAt = time.Now()
}

// DefaultAgentConfig returns default configuration for agents
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		HeartbeatInterval: 30 * time.Second,
		ReportingInterval: 5 * time.Minute,
		LogLevel:          "info",
		EnableMetrics:     true,
		EnableLogging:     true,
		MaxConnections:    1000,
		BufferSize:        8192,
		RetryAttempts:     3,
		TimeoutDuration:   30 * time.Second,
	}
}
