package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"scopeapi.local/backend/services/data-ingestion/internal/config"
	"scopeapi.local/backend/services/data-ingestion/internal/models"
	"scopeapi.local/backend/shared/logging"
)

type DataNormalizerServiceInterface interface {
	NormalizeData(ctx context.Context, parsed *models.ParsedData, schemaName string) (*models.NormalizedData, error)
	GetSchemas(ctx context.Context) ([]models.SchemaInfo, error)
	CreateSchema(ctx context.Context, schema models.SchemaInfo) error
	UpdateConfiguration(ctx context.Context, config interface{}) error
}

type DataNormalizerService struct {
	logger  logging.Logger
	config  *config.Config
	schemas map[string]models.SchemaInfo
	mutex   sync.RWMutex
}

func NewDataNormalizerService(logger logging.Logger, cfg *config.Config) DataNormalizerServiceInterface {
	service := &DataNormalizerService{
		logger:  logger,
		config:  cfg,
		schemas: make(map[string]models.SchemaInfo),
	}

	// Initialize with default schema if provided
	if cfg.Normalizer.DefaultSchema != "" {
		if schema, ok := cfg.Normalizer.Schemas[cfg.Normalizer.DefaultSchema]; ok {
			service.schemas[cfg.Normalizer.DefaultSchema] = models.SchemaInfo{
				Name:        schema.Name,
				Version:     schema.Version,
				Fields:      make(map[string]models.FieldInfo),
				Required:    schema.Required,
				Validators:  []models.ValidatorInfo{},
				Description: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
		}
	}

	return service
}

func (s *DataNormalizerService) NormalizeData(ctx context.Context, parsed *models.ParsedData, schemaName string) (*models.NormalizedData, error) {
	s.mutex.RLock()
	schema, exists := s.schemas[schemaName]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("schema not found: %s", schemaName)
	}

	// Map parsed data to schema fields
	parsedMap, ok := parsed.Parsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("parsed data is not an object")
	}

	normalized := &models.NormalizedData{
		ID:            parsed.ID,
		TrafficID:     parsed.TrafficID,
		ParsedID:      parsed.ID,
		Schema:        schema.Name,
		SchemaVersion: schema.Version,
		Data:          make(map[string]interface{}),
		Transformations: []models.Transformation{},
		CreatedAt:     time.Now(),
	}

	// Apply field mapping and transformations
	for field, fieldInfo := range schema.Fields {
		if value, exists := parsedMap[field]; exists {
			normalized.Data[field] = value
		} else if fieldInfo.Required {
			normalized.Data[field] = fieldInfo.Default
		}
	}

	// Validate required fields
	validation := s.validateNormalizedData(normalized, schema)
	normalized.Validation = validation

	if !validation.Valid {
		s.logger.Warn("Normalization validation failed", "errors", validation.Errors)
		return normalized, fmt.Errorf("normalization validation failed: %v", validation.Errors)
	}

	s.logger.Info("Data normalization completed", "schema", schema.Name, "id", normalized.ID)
	return normalized, nil
}

func (s *DataNormalizerService) GetSchemas(ctx context.Context) ([]models.SchemaInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var schemas []models.SchemaInfo
	for _, schema := range s.schemas {
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func (s *DataNormalizerService) CreateSchema(ctx context.Context, schema models.SchemaInfo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.schemas[schema.Name]; exists {
		return fmt.Errorf("schema already exists: %s", schema.Name)
	}
	s.schemas[schema.Name] = schema
	s.logger.Info("Schema created", "name", schema.Name)
	return nil
}

func (s *DataNormalizerService) UpdateConfiguration(ctx context.Context, config interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if cfg, ok := config.(map[string]interface{}); ok {
		// Update schemas
		if schemas, exists := cfg["schemas"]; exists {
			if schemaMap, ok := schemas.(map[string]models.SchemaInfo); ok {
				for name, schema := range schemaMap {
					s.schemas[name] = schema
				}
			}
		}
	}

	s.logger.Info("Normalizer configuration updated")
	return nil
}

// Helper methods
func (s *DataNormalizerService) validateNormalizedData(normalized *models.NormalizedData, schema models.SchemaInfo) models.ValidationResult {
	result := models.ValidationResult{
		Valid:    true,
		Score:    1.0,
		Errors:   []string{},
		Warnings: []string{},
	}

	for _, field := range schema.Required {
		if value, exists := normalized.Data[field]; !exists || value == nil {
			result.Valid = false
			result.Score = 0.0
			result.Errors = append(result.Errors, fmt.Sprintf("missing required field: %s", field))
		}
	}

	return result
} 