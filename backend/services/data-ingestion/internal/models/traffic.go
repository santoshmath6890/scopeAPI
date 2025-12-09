package models

import (
	"encoding/json"
	"time"
)

// TrafficData represents raw API traffic data
type TrafficData struct {
	ID            string                 `json:"id" db:"id"`
	Timestamp     time.Time              `json:"timestamp" db:"timestamp"`
	SourceIP      string                 `json:"source_ip" db:"source_ip"`
	DestinationIP string                 `json:"destination_ip" db:"destination_ip"`
	Protocol      string                 `json:"protocol" db:"protocol"`
	Method        string                 `json:"method" db:"method"`
	URL           string                 `json:"url" db:"url"`
	Path          string                 `json:"path" db:"path"`
	QueryParams   map[string]string      `json:"query_params" db:"query_params"`
	Headers       map[string]string      `json:"headers" db:"headers"`
	Body          []byte                 `json:"body" db:"body"`
	BodyText      string                 `json:"body_text" db:"body_text"`
	StatusCode    int                    `json:"status_code" db:"status_code"`
	ResponseBody  []byte                 `json:"response_body" db:"response_body"`
	ResponseText  string                 `json:"response_text" db:"response_text"`
	ResponseHeaders map[string]string    `json:"response_headers" db:"response_headers"`
	Duration      time.Duration          `json:"duration" db:"duration"`
	Size          int64                  `json:"size" db:"size"`
	UserAgent     string                 `json:"user_agent" db:"user_agent"`
	ContentType   string                 `json:"content_type" db:"content_type"`
	Encoding      string                 `json:"encoding" db:"encoding"`
	Compressed    bool                   `json:"compressed" db:"compressed"`
	Encrypted     bool                   `json:"encrypted" db:"encrypted"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	Tags          []string               `json:"tags" db:"tags"`
	Priority      int                    `json:"priority" db:"priority"`
	Status        string                 `json:"status" db:"status"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}

// BatchTrafficData represents a batch of traffic data
type BatchTrafficData struct {
	ID        string         `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	Count     int            `json:"count"`
	Data      []TrafficData  `json:"data"`
	Status    string         `json:"status"`
	Errors    []string       `json:"errors,omitempty"`
}

// ParsedData represents parsed traffic data
type ParsedData struct {
	ID           string                 `json:"id"`
	TrafficID    string                 `json:"traffic_id"`
	Format       string                 `json:"format"`
	ContentType  string                 `json:"content_type"`
	Parsed       interface{}            `json:"parsed"`
	Schema       map[string]interface{} `json:"schema"`
	Validation   ValidationResult       `json:"validation"`
	Errors       []ParseError           `json:"errors,omitempty"`
	Warnings     []ParseWarning         `json:"warnings,omitempty"`
	ProcessingTime time.Duration        `json:"processing_time"`
	CreatedAt    time.Time              `json:"created_at"`
}

// ValidationResult represents validation results
type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Score    float64  `json:"score"`
}

// ParseError represents parsing errors
type ParseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Path    string `json:"path,omitempty"`
}

// ParseWarning represents parsing warnings
type ParseWarning struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Path    string `json:"path,omitempty"`
}

// NormalizedData represents normalized traffic data
type NormalizedData struct {
	ID           string                 `json:"id"`
	TrafficID    string                 `json:"traffic_id"`
	ParsedID     string                 `json:"parsed_id"`
	Schema       string                 `json:"schema"`
	SchemaVersion string                `json:"schema_version"`
	Data         map[string]interface{} `json:"data"`
	Transformations []Transformation    `json:"transformations"`
	Validation   ValidationResult       `json:"validation"`
	Errors       []string               `json:"errors,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// Transformation represents data transformations
type Transformation struct {
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Config    map[string]interface{} `json:"config"`
	Input     interface{}            `json:"input"`
	Output    interface{}            `json:"output"`
	Duration  time.Duration          `json:"duration"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
}

// IngestionRequest represents an ingestion request
type IngestionRequest struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "single", "batch", "stream"
	Data        interface{}            `json:"data"`
	Format      string                 `json:"format,omitempty"`
	Schema      string                 `json:"schema,omitempty"`
	Priority    int                    `json:"priority"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Validation  bool                   `json:"validation"`
	Compression bool                   `json:"compression"`
	Encryption  bool                   `json:"encryption"`
	Timestamp   time.Time              `json:"timestamp"`
}

// IngestionResponse represents an ingestion response
type IngestionResponse struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	Message      string                 `json:"message"`
	TrafficIDs   []string               `json:"traffic_ids,omitempty"`
	ParsedIDs    []string               `json:"parsed_ids,omitempty"`
	NormalizedIDs []string              `json:"normalized_ids,omitempty"`
	Errors       []string               `json:"errors,omitempty"`
	Warnings     []string               `json:"warnings,omitempty"`
	ProcessingTime time.Duration        `json:"processing_time"`
	CreatedAt    time.Time              `json:"created_at"`
}

// IngestionStatus represents ingestion status
type IngestionStatus struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	Progress     float64                `json:"progress"`
	TotalItems   int                    `json:"total_items"`
	ProcessedItems int                  `json:"processed_items"`
	FailedItems  int                    `json:"failed_items"`
	Errors       []string               `json:"errors,omitempty"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      *time.Time             `json:"end_time,omitempty"`
	Duration     *time.Duration         `json:"duration,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// IngestionStats represents ingestion statistics
type IngestionStats struct {
	TotalIngested    int64     `json:"total_ingested"`
	TotalParsed      int64     `json:"total_parsed"`
	TotalNormalized  int64     `json:"total_normalized"`
	TotalErrors      int64     `json:"total_errors"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	SuccessRate      float64   `json:"success_rate"`
	ErrorRate        float64   `json:"error_rate"`
	LastIngestion    time.Time `json:"last_ingestion"`
	LastError        time.Time `json:"last_error"`
	Formats          map[string]int64 `json:"formats"`
	Schemas          map[string]int64 `json:"schemas"`
	Sources          map[string]int64 `json:"sources"`
	TimeRange        TimeRange `json:"time_range"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// QueueStatus represents queue status
type QueueStatus struct {
	Name           string                 `json:"name"`
	Size           int                    `json:"size"`
	MaxSize        int                    `json:"max_size"`
	Pending        int                    `json:"pending"`
	Processing     int                    `json:"processing"`
	Completed      int                    `json:"completed"`
	Failed         int                    `json:"failed"`
	RetryCount     int                    `json:"retry_count"`
	LastFlush      time.Time              `json:"last_flush"`
	NextFlush      time.Time              `json:"next_flush"`
	FlushInterval  time.Duration          `json:"flush_interval"`
	Health         string                 `json:"health"`
	Errors         []string               `json:"errors,omitempty"`
	Metrics        map[string]interface{} `json:"metrics,omitempty"`
}

// TopicInfo represents Kafka topic information
type TopicInfo struct {
	Name           string                 `json:"name"`
	Partitions     int                    `json:"partitions"`
	Replicas       int                    `json:"replicas"`
	Config         map[string]string      `json:"config"`
	Messages       int64                  `json:"messages"`
	Size           int64                  `json:"size"`
	ConsumerGroups []string               `json:"consumer_groups"`
	Health         string                 `json:"health"`
	LastMessage    time.Time              `json:"last_message"`
}

// FormatInfo represents supported format information
type FormatInfo struct {
	Name        string   `json:"name"`
	Extensions  []string `json:"extensions"`
	MimeTypes   []string `json:"mime_types"`
	Enabled     bool     `json:"enabled"`
	Priority    int      `json:"priority"`
	Description string   `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// SchemaInfo represents schema information
type SchemaInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Fields      map[string]FieldInfo   `json:"fields"`
	Required    []string               `json:"required"`
	Validators  []ValidatorInfo        `json:"validators"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// FieldInfo represents field information
type FieldInfo struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Format      string      `json:"format,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	Min         interface{} `json:"min,omitempty"`
	Max         interface{} `json:"max,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
}

// ValidatorInfo represents validator information
type ValidatorInfo struct {
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Config  map[string]interface{} `json:"config"`
	Message string                 `json:"message"`
}

// Methods for TrafficData
func (t *TrafficData) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TrafficData) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *TrafficData) GetHeader(key string) string {
	return t.Headers[key]
}

func (t *TrafficData) SetHeader(key, value string) {
	if t.Headers == nil {
		t.Headers = make(map[string]string)
	}
	t.Headers[key] = value
}

func (t *TrafficData) GetQueryParam(key string) string {
	return t.QueryParams[key]
}

func (t *TrafficData) SetQueryParam(key, value string) {
	if t.QueryParams == nil {
		t.QueryParams = make(map[string]string)
	}
	t.QueryParams[key] = value
}

func (t *TrafficData) AddTag(tag string) {
	t.Tags = append(t.Tags, tag)
}

func (t *TrafficData) HasTag(tag string) bool {
	for _, t := range t.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (t *TrafficData) SetMetadata(key string, value interface{}) {
	if t.Metadata == nil {
		t.Metadata = make(map[string]interface{})
	}
	t.Metadata[key] = value
}

func (t *TrafficData) GetMetadata(key string) interface{} {
	if t.Metadata == nil {
		return nil
	}
	return t.Metadata[key]
}

// Methods for BatchTrafficData
func (b *BatchTrafficData) AddData(data TrafficData) {
	b.Data = append(b.Data, data)
	b.Count = len(b.Data)
}

func (b *BatchTrafficData) AddError(error string) {
	b.Errors = append(b.Errors, error)
}

func (b *BatchTrafficData) IsComplete() bool {
	return b.Status == "completed" || b.Status == "failed"
}

// Methods for ParsedData
func (p *ParsedData) IsValid() bool {
	return p.Validation.Valid
}

func (p *ParsedData) AddError(errorType, message string) {
	p.Errors = append(p.Errors, ParseError{
		Type:    errorType,
		Message: message,
	})
}

func (p *ParsedData) AddWarning(warningType, message string) {
	p.Warnings = append(p.Warnings, ParseWarning{
		Type:    warningType,
		Message: message,
	})
}

// Methods for NormalizedData
func (n *NormalizedData) AddTransformation(transformation Transformation) {
	n.Transformations = append(n.Transformations, transformation)
}

func (n *NormalizedData) GetField(path string) interface{} {
	// Simple field access - could be enhanced with JSON path support
	if value, exists := n.Data[path]; exists {
		return value
	}
	return nil
}

func (n *NormalizedData) SetField(path string, value interface{}) {
	n.Data[path] = value
} 