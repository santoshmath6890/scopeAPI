package models

import (
	"time"
)

type APISpec struct {
	ID          string                 `json:"id" db:"id"`
	APIID       string                 `json:"api_id" db:"api_id"`
	Version     string                 `json:"version" db:"version"`
	Title       string                 `json:"title" db:"title"`
	Description string                 `json:"description" db:"description"`
	OpenAPIVersion string              `json:"openapi_version" db:"openapi_version"`
	Info        *SpecInfo              `json:"info" db:"info"`
	Servers     []SpecServer           `json:"servers" db:"servers"`
	Paths       map[string]*PathItem   `json:"paths" db:"paths"`
	Components  *SpecComponents        `json:"components,omitempty" db:"components"`
	Security    []SecurityRequirement  `json:"security,omitempty" db:"security"`
	Tags        []SpecTag              `json:"tags,omitempty" db:"tags"`
	ExternalDocs *ExternalDocumentation `json:"external_docs,omitempty" db:"external_docs"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

type SpecInfo struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"terms_of_service,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type SpecServer struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

type PathItem struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	OperationID  string                `json:"operation_id,omitempty"`
	Parameters   []Parameter           `json:"parameters,omitempty"`
	RequestBody  *RequestBody          `json:"request_body,omitempty"`
	Responses    map[string]Response   `json:"responses"`
	Callbacks    map[string]Callback   `json:"callbacks,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Servers      []SpecServer          `json:"servers,omitempty"`
}

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Content     map[string]MediaType  `json:"content"`
	Required    bool                  `json:"required,omitempty"`
}

type MediaType struct {
	Schema   *Schema                `json:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty"`
	Examples map[string]Example     `json:"examples,omitempty"`
	Encoding map[string]Encoding    `json:"encoding,omitempty"`
}

type Schema struct {
	Type                 string             `json:"type,omitempty"`
	Format               string             `json:"format,omitempty"`
	Title                string             `json:"title,omitempty"`
	Description          string             `json:"description,omitempty"`
	Default              interface{}        `json:"default,omitempty"`
	Example              interface{}        `json:"example,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	AdditionalProperties interface{}        `json:"additional_properties,omitempty"`
	Enum                 []interface{}      `json:"enum,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
	Maximum              *float64           `json:"maximum,omitempty"`
	MinLength            *int               `json:"min_length,omitempty"`
	MaxLength            *int               `json:"max_length,omitempty"`
	Pattern              string             `json:"pattern,omitempty"`
	MinItems             *int               `json:"min_items,omitempty"`
	MaxItems             *int               `json:"max_items,omitempty"`
	UniqueItems          bool               `json:"unique_items,omitempty"`
}

type Example struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"external_value,omitempty"`
}

type Encoding struct {
	ContentType   string             `json:"content_type,omitempty"`
	Headers       map[string]Header  `json:"headers,omitempty"`
	Style         string             `json:"style,omitempty"`
	Explode       bool               `json:"explode,omitempty"`
	AllowReserved bool               `json:"allow_reserved,omitempty"`
}

type Header struct {
	Description     string      `json:"description,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"`
	AllowEmptyValue bool        `json:"allow_empty_value,omitempty"`
	Style           string      `json:"style,omitempty"`
	Explode         bool        `json:"explode,omitempty"`
	AllowReserved   bool        `json:"allow_reserved,omitempty"`
	Schema          *Schema     `json:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Examples        map[string]Example `json:"examples,omitempty"`
}

type Callback map[string]PathItem

type SecurityRequirement map[string][]string

type SpecComponents struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"`
	Responses       map[string]Response        `json:"responses,omitempty"`
	Parameters      map[string]Parameter       `json:"parameters,omitempty"`
	Examples        map[string]Example         `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBody     `json:"request_bodies,omitempty"`
	Headers         map[string]Header          `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme  `json:"security_schemes,omitempty"`
	Links           map[string]Link            `json:"links,omitempty"`
	Callbacks       map[string]Callback        `json:"callbacks,omitempty"`
}

type SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Scheme           string            `json:"scheme,omitempty"`
	BearerFormat     string            `json:"bearer_format,omitempty"`
	Flows            *OAuthFlows       `json:"flows,omitempty"`
	OpenIDConnectURL string            `json:"open_id_connect_url,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"client_credentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorization_code,omitempty"`
}

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorization_url,omitempty"`
	TokenURL         string            `json:"token_url,omitempty"`
	RefreshURL       string            `json:"refresh_url,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

type Link struct {
	OperationRef string                 `json:"operation_ref,omitempty"`
	OperationID  string                 `json:"operation_id,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	RequestBody  interface{}            `json:"request_body,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Server       *SpecServer            `json:"server,omitempty"`
}

type SpecTag struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"external_docs,omitempty"`
}

type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}
