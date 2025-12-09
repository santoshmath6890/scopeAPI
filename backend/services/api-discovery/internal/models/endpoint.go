package models

import (
	"time"
)

type Endpoint struct {
	ID          string            `json:"id" db:"id"`
	APIID       string            `json:"api_id" db:"api_id"`
	URL         string            `json:"url" db:"url"`
	Path        string            `json:"path" db:"path"`
	Method      string            `json:"method" db:"method"`
	Headers     map[string]string `json:"headers" db:"headers"`
	Body        string            `json:"body" db:"body"`
	StatusCode  int               `json:"status_code" db:"status_code"`
	ContentType string            `json:"content_type" db:"content_type"`
	Summary     string            `json:"summary" db:"summary"`
	Description string            `json:"description" db:"description"`
	Parameters  []Parameter       `json:"parameters" db:"parameters"`
	Responses   map[string]Response `json:"responses" db:"responses"`
	Tags        []string          `json:"tags" db:"tags"`
	IsActive    bool              `json:"is_active" db:"is_active"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"` // query, header, path, body
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Example     interface{} `json:"example,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
}

type Response struct {
	StatusCode  int                    `json:"status_code"`
	Description string                 `json:"description"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
}

type EndpointAnalysis struct {
	EndpointID    string            `json:"endpoint_id"`
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	ResponseTime time.Duration     `json:"response_time"`
	StatusCode   int               `json:"status_code"`
	ContentType  string            `json:"content_type"`
	Parameters   []Parameter       `json:"parameters"`
	Headers      map[string]string `json:"headers"`
	Security     *SecurityAnalysis `json:"security"`
	CreatedAt    time.Time         `json:"created_at"`
}

type SecurityAnalysis struct {
	HasHTTPS           bool     `json:"has_https"`
	HasSecurityHeaders bool     `json:"has_security_headers"`
	VulnerableHeaders  []string `json:"vulnerable_headers"`
	AuthenticationRequired bool `json:"authentication_required"`
	RateLimitHeaders   map[string]string `json:"rate_limit_headers"`
}

type API struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	URL         string    `json:"url" db:"url"`
	BaseURL     string    `json:"base_url" db:"base_url"`
	Version     string    `json:"version" db:"version"`
	Protocol    string    `json:"protocol" db:"protocol"`
	Status      string    `json:"status" db:"status"`
	Description string    `json:"description" db:"description"`
	Tags        []string  `json:"tags" db:"tags"`
	Endpoints   []Endpoint `json:"endpoints,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Discovery struct {
	ID             string           `json:"id" db:"id"`
	Target         string           `json:"target" db:"target"`
	Method         string           `json:"method" db:"method"`
	Status         string           `json:"status" db:"status"`
	Progress       int              `json:"progress" db:"progress"`
	StartTime      time.Time        `json:"start_time" db:"start_time"`
	EndTime        *time.Time       `json:"end_time,omitempty" db:"end_time"`
	EndpointsFound int              `json:"endpoints_found" db:"endpoints_found"`
	ErrorMessage   string           `json:"error_message,omitempty" db:"error_message"`
	Config         *DiscoveryConfig `json:"config" db:"config"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

type DiscoveryConfig struct {
	Target      string            `json:"target"`
	Method      string            `json:"method"`
	Options     map[string]string `json:"options"`
	Credentials *Credentials      `json:"credentials,omitempty"`
}

type Credentials struct {
	Type     string `json:"type"` // basic, bearer, api_key
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
}

type DiscoveryStatus struct {
	ID             string     `json:"id"`
	Status         string     `json:"status"`
	Progress       int        `json:"progress"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	EndpointsFound int        `json:"endpoints_found"`
	ErrorMessage   string     `json:"error_message,omitempty"`
}

type DiscoveryResults struct {
	DiscoveryID string     `json:"discovery_id"`
	Total       int        `json:"total"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
	Endpoints   []Endpoint `json:"endpoints"`
}

type APIInventory struct {
	Total int   `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	APIs  []API `json:"apis"`
}

type APIDetails struct {
	API       API        `json:"api"`
	Endpoints []Endpoint `json:"endpoints"`
	Statistics APIStats  `json:"statistics"`
}

type APIStats struct {
	TotalEndpoints    int `json:"total_endpoints"`
	ActiveEndpoints   int `json:"active_endpoints"`
	InactiveEndpoints int `json:"inactive_endpoints"`
	LastScanned       *time.Time `json:"last_scanned,omitempty"`
}

type APIStatistics struct {
	TotalAPIs         int            `json:"total_apis"`
	ActiveAPIs        int            `json:"active_apis"`
	InactiveAPIs      int            `json:"inactive_apis"`
	TotalEndpoints    int            `json:"total_endpoints"`
	RecentDiscoveries int            `json:"recent_discoveries"`
	ProtocolBreakdown map[string]int `json:"protocol_breakdown"`
	StatusBreakdown   map[string]int `json:"status_breakdown"`
}
