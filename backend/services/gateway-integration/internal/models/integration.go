package models

import (
	"encoding/json"
	"time"
)

// Integration represents a gateway integration configuration
type Integration struct {
	ID          string                 `json:"id" db:"integration_id"`
	Name        string                 `json:"name" db:"name"`
	Type        GatewayType            `json:"type" db:"type"`
	Status      IntegrationStatus      `json:"status" db:"status"`
	Config      map[string]interface{} `json:"config" db:"config"`
	Credentials *Credentials           `json:"credentials,omitempty" db:"credentials"`
	Endpoints   []Endpoint             `json:"endpoints" db:"endpoints"`
	Health      *HealthStatus          `json:"health" db:"health"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	LastSync    *time.Time             `json:"last_sync" db:"last_sync"`
}

// GatewayType represents supported gateway types
type GatewayType string

const (
	GatewayTypeKong    GatewayType = "kong"
	GatewayTypeNginx   GatewayType = "nginx"
	GatewayTypeTraefik GatewayType = "traefik"
	GatewayTypeEnvoy   GatewayType = "envoy"
	GatewayTypeHAProxy GatewayType = "haproxy"
)

// IntegrationStatus represents the status of an integration
type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusInactive IntegrationStatus = "inactive"
	IntegrationStatusError    IntegrationStatus = "error"
	IntegrationStatusPending  IntegrationStatus = "pending"
)

// Credentials represents authentication credentials for gateway integration
type Credentials struct {
	Type     CredentialType `json:"type"`
	Username string         `json:"username,omitempty"`
	Password string         `json:"password,omitempty"`
	Token    string         `json:"token,omitempty"`
	APIKey   string         `json:"api_key,omitempty"`
	CertFile string         `json:"cert_file,omitempty"`
	KeyFile  string         `json:"key_file,omitempty"`
}

// CredentialType represents the type of credentials
type CredentialType string

const (
	CredentialTypeBasic  CredentialType = "basic"
	CredentialTypeToken  CredentialType = "token"
	CredentialTypeAPIKey CredentialType = "api_key"
	CredentialTypeTLS    CredentialType = "tls"
)

// Endpoint represents a gateway endpoint configuration
type Endpoint struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	URL      string            `json:"url"`
	Protocol string            `json:"protocol"`
	Port     int               `json:"port"`
	Headers  map[string]string `json:"headers,omitempty"`
	Timeout  time.Duration     `json:"timeout"`
}

// HealthStatus represents the health status of an integration
type HealthStatus struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	LastCheck time.Time `json:"last_check"`
	Latency   int64     `json:"latency_ms"`
}

// IntegrationConfig represents configuration for a specific gateway type
type IntegrationConfig struct {
	ID          string                 `json:"id" db:"config_id"`
	IntegrationID string               `json:"integration_id" db:"integration_id"`
	Type        GatewayType            `json:"type" db:"type"`
	Config      map[string]interface{} `json:"config" db:"config"`
	Version     int                    `json:"version" db:"version"`
	Active      bool                   `json:"active" db:"active"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// GatewayConfig represents a gateway configuration with versioning
type GatewayConfig struct {
	ID             int64                  `json:"id" db:"id"`
	IntegrationID  int64                  `json:"integration_id" db:"integration_id"`
	ConfigType     string                 `json:"config_type" db:"config_type"`
	ConfigData     map[string]interface{} `json:"config_data" db:"config_data"`
	Version        int                    `json:"version" db:"version"`
	Status         string                 `json:"status" db:"status"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// IntegrationStats represents statistics about integrations
type IntegrationStats struct {
	Total        int64 `json:"total"`
	Active       int64 `json:"active"`
	Error        int64 `json:"error"`
	Pending      int64 `json:"pending"`
	KongCount   int64 `json:"kong_count"`
	NginxCount  int64 `json:"nginx_count"`
	TraefikCount int64 `json:"traefik_count"`
	EnvoyCount  int64 `json:"envoy_count"`
	HAProxyCount int64 `json:"haproxy_count"`
}

// HAProxyStats represents HAProxy statistics
type HAProxyStats struct {
	TotalConnections int64 `json:"total_connections"`
	ActiveConnections int64 `json:"active_connections"`
	RequestsPerSecond int64 `json:"requests_per_second"`
	BytesIn int64 `json:"bytes_in"`
	BytesOut int64 `json:"bytes_out"`
	SessionRate int64 `json:"session_rate"`
	MaxSessionRate int64 `json:"max_session_rate"`
}

// KongConfig represents Kong-specific configuration
type KongConfig struct {
	AdminURL    string            `json:"admin_url"`
	ProxyURL    string            `json:"proxy_url"`
	Plugins     []KongPlugin      `json:"plugins"`
	Services    []KongService     `json:"services"`
	Routes      []KongRoute       `json:"routes"`
	Consumers   []KongConsumer    `json:"consumers"`
	Upstreams   []KongUpstream    `json:"upstreams"`
	Certificates []KongCertificate `json:"certificates"`
}

// KongPlugin represents a Kong plugin configuration
type KongPlugin struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	ServiceID  string                 `json:"service_id,omitempty"`
	RouteID    string                 `json:"route_id,omitempty"`
	ConsumerID string                 `json:"consumer_id,omitempty"`
	Config     map[string]interface{} `json:"config"`
	Enabled    bool                   `json:"enabled"`
	Protocols  []string               `json:"protocols"`
}

// KongService represents a Kong service
type KongService struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Path         string `json:"path"`
	Retries      int    `json:"retries"`
	ConnectTimeout int  `json:"connect_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
}

// KongRoute represents a Kong route
type KongRoute struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ServiceID    string   `json:"service_id"`
	Protocols    []string `json:"protocols"`
	Methods      []string `json:"methods"`
	Hosts        []string `json:"hosts"`
	Paths        []string `json:"paths"`
	StripPath    bool     `json:"strip_path"`
	PreserveHost bool     `json:"preserve_host"`
}

// KongConsumer represents a Kong consumer
type KongConsumer struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	CustomID string            `json:"custom_id"`
	Tags     []string          `json:"tags"`
}

// KongUpstream represents a Kong upstream
type KongUpstream struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Algorithm   string `json:"algorithm"`
	Slots       int    `json:"slots"`
	OrderList   []int  `json:"order_list"`
	HashOn      string `json:"hash_on"`
	HashFallback string `json:"hash_fallback"`
	HashOnHeader string `json:"hash_on_header"`
	HashFallbackHeader string `json:"hash_fallback_header"`
	HashOnCookie string `json:"hash_on_cookie"`
	HashOnCookiePath string `json:"hash_on_cookie_path"`
}

// KongCertificate represents a Kong certificate
type KongCertificate struct {
	ID   string `json:"id"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
	SNIs []string `json:"snis"`
}

// NginxConfig represents NGINX-specific configuration
type NginxConfig struct {
	ConfigPath    string            `json:"config_path"`
	ReloadCommand string            `json:"reload_command"`
	TestCommand   string            `json:"test_command"`
	Upstreams     []NginxUpstream   `json:"upstreams"`
	Locations     []NginxLocation   `json:"locations"`
	Servers       []NginxServer     `json:"servers"`
	SSLConfigs    []NginxSSLConfig  `json:"ssl_configs"`
}

// NginxUpstream represents an NGINX upstream configuration
type NginxUpstream struct {
	Name    string            `json:"name"`
	Servers []NginxServer     `json:"servers"`
	Options map[string]string `json:"options"`
}

// NginxLocation represents an NGINX location configuration
type NginxLocation struct {
	Path        string            `json:"path"`
	ProxyPass   string            `json:"proxy_pass"`
	Headers     map[string]string `json:"headers"`
	Options     map[string]string `json:"options"`
}

// NginxServer represents an NGINX server configuration
type NginxServer struct {
	Listen    []string          `json:"listen"`
	ServerName []string         `json:"server_name"`
	Locations []NginxLocation   `json:"locations"`
	SSL       *NginxSSLConfig   `json:"ssl"`
}

// NginxSSLConfig represents NGINX SSL configuration
type NginxSSLConfig struct {
	Certificate     string `json:"certificate"`
	CertificateKey string `json:"certificate_key"`
	Protocols       string `json:"protocols"`
	Ciphers         string `json:"ciphers"`
}

// TraefikConfig represents Traefik-specific configuration
type TraefikConfig struct {
	APIEndpoint    string                `json:"api_endpoint"`
	Dashboard      bool                  `json:"dashboard"`
	Providers      []TraefikProvider     `json:"providers"`
	Middlewares    []TraefikMiddleware   `json:"middlewares"`
	Routers        []TraefikRouter       `json:"routers"`
	Services       []TraefikService      `json:"services"`
	TLSConfigs     []TraefikTLSConfig    `json:"tls_configs"`
}

// TraefikProvider represents a Traefik provider
type TraefikProvider struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// TraefikMiddleware represents a Traefik middleware
type TraefikMiddleware struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// TraefikRouter represents a Traefik router
type TraefikRouter struct {
	Name       string   `json:"name"`
	EntryPoints []string `json:"entry_points"`
	Rule       string   `json:"rule"`
	Service    string   `json:"service"`
	Middlewares []string `json:"middlewares"`
	TLS        *TraefikTLSConfig `json:"tls"`
}

// TraefikService represents a Traefik service
type TraefikService struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
}

// TraefikTLSConfig represents Traefik TLS configuration
type TraefikTLSConfig struct {
	CertResolver string   `json:"cert_resolver"`
	Domains      []string `json:"domains"`
}

// EnvoyConfig represents Envoy-specific configuration
type EnvoyConfig struct {
	AdminAddress string           `json:"admin_address"`
	Clusters     []EnvoyCluster   `json:"clusters"`
	Listeners    []EnvoyListener  `json:"listeners"`
	Routes       []EnvoyRoute     `json:"routes"`
	Filters      []EnvoyFilter    `json:"filters"`
}

// EnvoyCluster represents an Envoy cluster
type EnvoyCluster struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Endpoints []EnvoyEndpoint   `json:"endpoints"`
	Options   map[string]string `json:"options"`
}

// EnvoyEndpoint represents an Envoy endpoint
type EnvoyEndpoint struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// EnvoyListener represents an Envoy listener
type EnvoyListener struct {
	Name    string        `json:"name"`
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Filters []EnvoyFilter `json:"filters"`
}

// EnvoyRoute represents an Envoy route
type EnvoyRoute struct {
	Name    string            `json:"name"`
	Match   map[string]string `json:"match"`
	Action  string            `json:"action"`
	Cluster string            `json:"cluster"`
}

// EnvoyFilter represents an Envoy filter
type EnvoyFilter struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// HAProxyConfig represents HAProxy-specific configuration
type HAProxyConfig struct {
	ConfigPath    string           `json:"config_path"`
	ReloadCommand string           `json:"reload_command"`
	TestCommand   string           `json:"test_command"`
	Frontends     []HAProxyFrontend `json:"frontends"`
	Backends      []HAProxyBackend  `json:"backends"`
	Servers       []HAProxyServer   `json:"servers"`
	Defaults      map[string]string `json:"defaults"`
}

// HAProxyFrontend represents an HAProxy frontend
type HAProxyFrontend struct {
	Name    string            `json:"name"`
	Bind    string            `json:"bind"`
	Mode    string            `json:"mode"`
	Options map[string]string `json:"options"`
	Rules   []HAProxyRule     `json:"rules"`
}

// HAProxyBackend represents an HAProxy backend
type HAProxyBackend struct {
	Name    string            `json:"name"`
	Mode    string            `json:"mode"`
	Balance string            `json:"balance"`
	Servers []HAProxyServer   `json:"servers"`
	Options map[string]string `json:"options"`
}

// HAProxyServer represents an HAProxy server
type HAProxyServer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Check   bool   `json:"check"`
}

// HAProxyRule represents an HAProxy rule
type HAProxyRule struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Backend string `json:"backend"`
}

// IntegrationEvent represents an integration event
type IntegrationEvent struct {
	ID          string                 `json:"id"`
	IntegrationID string               `json:"integration_id"`
	Type        string                 `json:"type"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
}

// SyncResult represents the result of a configuration sync
type SyncResult struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Changes     []Change               `json:"changes"`
	Errors      []string               `json:"errors"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`
}

// Change represents a configuration change
type Change struct {
	Type      string `json:"type"`
	Resource  string `json:"resource"`
	Action    string `json:"action"`
	Details   string `json:"details"`
}

// MarshalJSON custom marshaling for Integration to handle credentials securely
func (i *Integration) MarshalJSON() ([]byte, error) {
	type Alias Integration
	return json.Marshal(&struct {
		*Alias
		Credentials *Credentials `json:"credentials,omitempty"`
	}{
		Alias:       (*Alias)(i),
		Credentials: i.Credentials, // Will be nil in JSON response for security
	})
} 