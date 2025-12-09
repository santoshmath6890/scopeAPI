package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"scopeapi.local/backend/services/data-ingestion/internal/config"
	"scopeapi.local/backend/services/data-ingestion/internal/models"
	"scopeapi.local/backend/shared/logging"
)

type DataParserServiceInterface interface {
	ParseData(ctx context.Context, data []byte, format string, contentType string) (*models.ParsedData, error)
	GetSupportedFormats(ctx context.Context) ([]models.FormatInfo, error)
	ValidateFormat(ctx context.Context, data []byte, format string) (*models.ValidationResult, error)
	UpdateConfiguration(ctx context.Context, config interface{}) error
}

type DataParserService struct {
	logger  logging.Logger
	config  *config.Config
	formats map[string]*config.FormatConfig
	mutex   sync.RWMutex
}

func NewDataParserService(logger logging.Logger, cfg *config.Config) DataParserServiceInterface {
	service := &DataParserService{
		logger:  logger,
		config:  cfg,
		formats: make(map[string]*config.FormatConfig),
	}

	// Initialize supported formats
	service.initializeFormats()

	return service
}

func (s *DataParserService) ParseData(ctx context.Context, data []byte, format string, contentType string) (*models.ParsedData, error) {
	startTime := time.Now()
	parsedData := &models.ParsedData{
		ID:            uuid.New().String(),
		Format:        format,
		ContentType:   contentType,
		ProcessingTime: time.Duration(0),
		CreatedAt:     time.Now(),
	}

	s.logger.Info("Starting data parsing", "format", format, "content_type", contentType, "size", len(data))

	// Validate input data
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data provided")
	}

	if len(data) > s.config.Parser.MaxPayloadSize {
		return nil, fmt.Errorf("data size exceeds maximum allowed size: %d > %d", len(data), s.config.Parser.MaxPayloadSize)
	}

	// Detect format if not specified
	if format == "" {
		format = s.detectFormat(data, contentType)
	}

	// Parse data based on format
	var parsed interface{}
	var err error

	switch strings.ToLower(format) {
	case "json":
		parsed, err = s.parseJSON(data)
	case "xml":
		parsed, err = s.parseXML(data)
	case "yaml", "yml":
		parsed, err = s.parseYAML(data)
	case "protobuf":
		parsed, err = s.parseProtobuf(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		s.logger.Error("Failed to parse data", "error", err, "format", format)
		parsedData.AddError("parse_error", err.Error())
		parsedData.Validation.Valid = false
		parsedData.Validation.Errors = append(parsedData.Validation.Errors, err.Error())
		parsedData.ProcessingTime = time.Since(startTime)
		return parsedData, err
	}

	// Validate parsed data
	validationResult := s.validateParsedData(parsed, format)
	parsedData.Validation = validationResult

	// Extract schema
	schema := s.extractSchema(parsed, format)
	parsedData.Schema = schema

	// Set parsed data
	parsedData.Parsed = parsed
	parsedData.ProcessingTime = time.Since(startTime)

	if validationResult.Valid {
		s.logger.Info("Data parsing completed successfully", "format", format, "processing_time", parsedData.ProcessingTime)
	} else {
		s.logger.Warn("Data parsing completed with validation errors", "format", format, "errors", len(validationResult.Errors))
	}

	return parsedData, nil
}

func (s *DataParserService) GetSupportedFormats(ctx context.Context) ([]models.FormatInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var formats []models.FormatInfo
	for name, config := range s.formats {
		if config.Enabled {
			formats = append(formats, models.FormatInfo{
				Name:        name,
				Extensions:  config.Extensions,
				MimeTypes:   config.MimeTypes,
				Enabled:     config.Enabled,
				Priority:    config.Priority,
				Description: s.getFormatDescription(name),
				Config:      config.Config,
			})
		}
	}

	return formats, nil
}

func (s *DataParserService) ValidateFormat(ctx context.Context, data []byte, format string) (*models.ValidationResult, error) {
	// Try to parse the data with the specified format
	parsedData, err := s.ParseData(ctx, data, format, "")
	if err != nil {
		return &models.ValidationResult{
			Valid:  false,
			Errors: []string{err.Error()},
			Score:  0.0,
		}, nil
	}

	return &parsedData.Validation, nil
}

func (s *DataParserService) UpdateConfiguration(ctx context.Context, config interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if cfg, ok := config.(map[string]interface{}); ok {
		// Update parser configuration
		if maxPayloadSize, exists := cfg["max_payload_size"]; exists {
			if mps, ok := maxPayloadSize.(int); ok {
				s.config.Parser.MaxPayloadSize = mps
			}
		}
		if timeout, exists := cfg["timeout"]; exists {
			if t, ok := timeout.(time.Duration); ok {
				s.config.Parser.Timeout = t
			}
		}

		// Update format configurations
		if formats, exists := cfg["formats"]; exists {
			if formatConfigs, ok := formats.(map[string]interface{}); ok {
				for formatName, formatConfig := range formatConfigs {
					if fc, ok := formatConfig.(map[string]interface{}); ok {
						s.updateFormatConfig(formatName, fc)
					}
				}
			}
		}
	}

	s.logger.Info("Parser configuration updated")
	return nil
}

// Helper methods

func (s *DataParserService) initializeFormats() {
	// JSON format
	s.formats["json"] = &config.FormatConfig{
		Name:       "json",
		Extensions: []string{".json"},
		MimeTypes:  []string{"application/json", "text/json"},
		Enabled:    true,
		Priority:   1,
		Config: map[string]interface{}{
			"allow_comments": false,
			"strict":         true,
		},
	}

	// XML format
	s.formats["xml"] = &config.FormatConfig{
		Name:       "xml",
		Extensions: []string{".xml"},
		MimeTypes:  []string{"application/xml", "text/xml"},
		Enabled:    true,
		Priority:   2,
		Config: map[string]interface{}{
			"strict": true,
		},
	}

	// YAML format
	s.formats["yaml"] = &config.FormatConfig{
		Name:       "yaml",
		Extensions: []string{".yaml", ".yml"},
		MimeTypes:  []string{"application/yaml", "text/yaml", "application/x-yaml"},
		Enabled:    true,
		Priority:   3,
		Config: map[string]interface{}{
			"strict": true,
		},
	}

	// Protobuf format
	s.formats["protobuf"] = &config.FormatConfig{
		Name:       "protobuf",
		Extensions: []string{".proto", ".pb"},
		MimeTypes:  []string{"application/x-protobuf", "application/protobuf"},
		Enabled:    true,
		Priority:   4,
		Config: map[string]interface{}{
			"strict": true,
		},
	}
}

func (s *DataParserService) detectFormat(data []byte, contentType string) string {
	// Try to detect format from content type
	if contentType != "" {
		switch {
		case strings.Contains(contentType, "json"):
			return "json"
		case strings.Contains(contentType, "xml"):
			return "xml"
		case strings.Contains(contentType, "yaml"):
			return "yaml"
		case strings.Contains(contentType, "protobuf"):
			return "protobuf"
		}
	}

	// Try to detect format from data content
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return "json" // Default to JSON
	}

	switch trimmed[0] {
	case '{', '[':
		return "json"
	case '<':
		return "xml"
	case '#', '-', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		return "yaml"
	default:
		return "json" // Default to JSON
	}
}

func (s *DataParserService) parseJSON(data []byte) (interface{}, error) {
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func (s *DataParserService) parseXML(data []byte) (interface{}, error) {
	var result map[string]interface{}
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	return result, nil
}

func (s *DataParserService) parseYAML(data []byte) (interface{}, error) {
	var result interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}

func (s *DataParserService) parseProtobuf(data []byte) (interface{}, error) {
	// This is a simplified implementation
	// In a real implementation, you would need protobuf schema and proper parsing
	return nil, fmt.Errorf("protobuf parsing not implemented")
}

func (s *DataParserService) validateParsedData(data interface{}, format string) models.ValidationResult {
	result := models.ValidationResult{
		Valid:  true,
		Score:  1.0,
		Errors: []string{},
		Warnings: []string{},
	}

	// Basic validation
	if data == nil {
		result.Valid = false
		result.Score = 0.0
		result.Errors = append(result.Errors, "parsed data is nil")
		return result
	}

	// Format-specific validation
	switch format {
	case "json":
		result = s.validateJSON(data)
	case "xml":
		result = s.validateXML(data)
	case "yaml":
		result = s.validateYAML(data)
	case "protobuf":
		result = s.validateProtobuf(data)
	}

	return result
}

func (s *DataParserService) validateJSON(data interface{}) models.ValidationResult {
	result := models.ValidationResult{
		Valid:  true,
		Score:  1.0,
		Errors: []string{},
		Warnings: []string{},
	}

	// Check if data is a valid JSON structure
	switch v := data.(type) {
	case map[string]interface{}:
		// Valid JSON object
		if len(v) == 0 {
			result.Warnings = append(result.Warnings, "empty JSON object")
		}
	case []interface{}:
		// Valid JSON array
		if len(v) == 0 {
			result.Warnings = append(result.Warnings, "empty JSON array")
		}
	case string, float64, bool, nil:
		// Valid JSON primitive
	default:
		result.Valid = false
		result.Score = 0.0
		result.Errors = append(result.Errors, "invalid JSON structure")
	}

	return result
}

func (s *DataParserService) validateXML(data interface{}) models.ValidationResult {
	result := models.ValidationResult{
		Valid:  true,
		Score:  1.0,
		Errors: []string{},
		Warnings: []string{},
	}

	// Check if data is a valid XML structure
	if _, ok := data.(map[string]interface{}); !ok {
		result.Valid = false
		result.Score = 0.0
		result.Errors = append(result.Errors, "invalid XML structure")
	}

	return result
}

func (s *DataParserService) validateYAML(data interface{}) models.ValidationResult {
	result := models.ValidationResult{
		Valid:  true,
		Score:  1.0,
		Errors: []string{},
		Warnings: []string{},
	}

	// YAML can be any valid data structure
	// Basic validation - check if data is not nil
	if data == nil {
		result.Valid = false
		result.Score = 0.0
		result.Errors = append(result.Errors, "parsed YAML data is nil")
	}

	return result
}

func (s *DataParserService) validateProtobuf(data interface{}) models.ValidationResult {
	result := models.ValidationResult{
		Valid:  false,
		Score:  0.0,
		Errors: []string{"protobuf validation not implemented"},
		Warnings: []string{},
	}

	return result
}

func (s *DataParserService) extractSchema(data interface{}, format string) map[string]interface{} {
	schema := make(map[string]interface{})

	switch v := data.(type) {
	case map[string]interface{}:
		schema["type"] = "object"
		schema["properties"] = s.extractObjectSchema(v)
	case []interface{}:
		schema["type"] = "array"
		if len(v) > 0 {
			schema["items"] = s.extractSchema(v[0], format)
		}
	case string:
		schema["type"] = "string"
	case float64:
		schema["type"] = "number"
	case bool:
		schema["type"] = "boolean"
	case nil:
		schema["type"] = "null"
	}

	return schema
}

func (s *DataParserService) extractObjectSchema(obj map[string]interface{}) map[string]interface{} {
	properties := make(map[string]interface{})

	for key, value := range obj {
		properties[key] = s.extractSchema(value, "")
	}

	return properties
}

func (s *DataParserService) getFormatDescription(format string) string {
	descriptions := map[string]string{
		"json":      "JavaScript Object Notation - lightweight data interchange format",
		"xml":       "Extensible Markup Language - markup language for data representation",
		"yaml":      "YAML Ain't Markup Language - human-readable data serialization format",
		"protobuf":  "Protocol Buffers - language-neutral data serialization format",
	}

	if desc, exists := descriptions[format]; exists {
		return desc
	}
	return "Unknown format"
}

func (s *DataParserService) updateFormatConfig(formatName string, config map[string]interface{}) {
	if format, exists := s.formats[formatName]; exists {
		if enabled, exists := config["enabled"]; exists {
			if e, ok := enabled.(bool); ok {
				format.Enabled = e
			}
		}
		if priority, exists := config["priority"]; exists {
			if p, ok := priority.(int); ok {
				format.Priority = p
			}
		}
		if extensions, exists := config["extensions"]; exists {
			if ext, ok := extensions.([]string); ok {
				format.Extensions = ext
			}
		}
		if mimeTypes, exists := config["mime_types"]; exists {
			if mime, ok := mimeTypes.([]string); ok {
				format.MimeTypes = mime
			}
		}
		if cfg, exists := config["config"]; exists {
			if c, ok := cfg.(map[string]interface{}); ok {
				format.Config = c
			}
		}
	}
} 