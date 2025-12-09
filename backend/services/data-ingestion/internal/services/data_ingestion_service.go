package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"scopeapi.local/backend/services/data-ingestion/internal/config"
	"scopeapi.local/backend/services/data-ingestion/internal/models"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type DataIngestionServiceInterface interface {
	IngestTraffic(ctx context.Context, request *models.IngestionRequest) (*models.IngestionResponse, error)
	IngestBatch(ctx context.Context, batch *models.BatchTrafficData) (*models.IngestionResponse, error)
	GetIngestionStatus(ctx context.Context, id string) (*models.IngestionStatus, error)
	GetIngestionStats(ctx context.Context, timeRange *models.TimeRange) (*models.IngestionStats, error)
	UpdateConfiguration(ctx context.Context, configType string, config interface{}) error
}

type DataIngestionService struct {
	kafkaProducer kafka.ProducerInterface
	logger        logging.Logger
	config        *config.Config
	statusMap     map[string]*models.IngestionStatus
	stats         *models.IngestionStats
	mutex         sync.RWMutex
	parserService DataParserServiceInterface
	normalizerService DataNormalizerServiceInterface
}

func NewDataIngestionService(
	kafkaProducer kafka.ProducerInterface,
	logger logging.Logger,
	cfg *config.Config,
) DataIngestionServiceInterface {
	service := &DataIngestionService{
		kafkaProducer: kafkaProducer,
		logger:        logger,
		config:        cfg,
		statusMap:     make(map[string]*models.IngestionStatus),
		stats:         &models.IngestionStats{},
	}

	// Initialize parser and normalizer services
	service.parserService = NewDataParserService(logger, cfg)
	service.normalizerService = NewDataNormalizerService(logger, cfg)

	// Start background tasks
	go service.startStatsCollector()
	go service.startStatusCleanup()

	return service
}

func (s *DataIngestionService) IngestTraffic(ctx context.Context, request *models.IngestionRequest) (*models.IngestionResponse, error) {
	startTime := time.Now()
	response := &models.IngestionResponse{
		ID:        uuid.New().String(),
		Status:    "processing",
		CreatedAt: time.Now(),
	}

	s.logger.Info("Starting traffic ingestion", "request_id", request.ID, "type", request.Type)

	// Create ingestion status
	status := &models.IngestionStatus{
		ID:         response.ID,
		Status:     "processing",
		Progress:   0.0,
		TotalItems: 1,
		StartTime:  time.Now(),
	}
	s.updateStatus(status)

	// Process the traffic data
	trafficData, err := s.processTrafficData(request)
	if err != nil {
		s.logger.Error("Failed to process traffic data", "error", err, "request_id", request.ID)
		response.Status = "failed"
		response.Errors = append(response.Errors, err.Error())
		status.Status = "failed"
		status.FailedItems = 1
		now := time.Now()
		status.EndTime = &now
		s.updateStatus(status)
		return response, err
	}

	// Publish to Kafka
	if err := s.publishToKafka(ctx, trafficData); err != nil {
		s.logger.Error("Failed to publish to Kafka", "error", err, "request_id", request.ID)
		response.Status = "failed"
		response.Errors = append(response.Errors, err.Error())
		status.Status = "failed"
		status.FailedItems = 1
		now := time.Now()
		status.EndTime = &now
		s.updateStatus(status)
		return response, err
	}

	// Update response and status
	response.Status = "success"
	response.TrafficIDs = []string{trafficData.ID}
	response.ProcessingTime = time.Since(startTime)
	status.Status = "completed"
	status.ProcessedItems = 1
	status.Progress = 100.0
	now := time.Now()
	status.EndTime = &now
	s.updateStatus(status)

	// Update statistics
	s.updateStats(trafficData, time.Since(startTime), nil)

	s.logger.Info("Traffic ingestion completed", "request_id", request.ID, "traffic_id", trafficData.ID)
	return response, nil
}

func (s *DataIngestionService) IngestBatch(ctx context.Context, batch *models.BatchTrafficData) (*models.IngestionResponse, error) {
	startTime := time.Now()
	response := &models.IngestionResponse{
		ID:        uuid.New().String(),
		Status:    "processing",
		CreatedAt: time.Now(),
	}

	s.logger.Info("Starting batch ingestion", "batch_id", batch.ID, "count", batch.Count)

	// Create ingestion status
	status := &models.IngestionStatus{
		ID:         response.ID,
		Status:     "processing",
		Progress:   0.0,
		TotalItems: batch.Count,
		StartTime:  time.Now(),
	}
	s.updateStatus(status)

	// Process batch items
	var trafficIDs []string
	var errors []string
	processedCount := 0
	failedCount := 0

	for i, trafficData := range batch.Data {
		// Process individual traffic data
		processed, err := s.processTrafficData(&models.IngestionRequest{
			ID:         trafficData.ID,
			Type:       "single",
			Data:       trafficData,
			Priority:   trafficData.Priority,
			Tags:       trafficData.Tags,
			Validation: true,
			Timestamp:  trafficData.Timestamp,
		})

		if err != nil {
			s.logger.Error("Failed to process batch item", "error", err, "index", i, "traffic_id", trafficData.ID)
			errors = append(errors, fmt.Sprintf("Item %d: %s", i, err.Error()))
			failedCount++
		} else {
			trafficIDs = append(trafficIDs, processed.ID)
			processedCount++
		}

		// Update progress
		progress := float64(i+1) / float64(batch.Count) * 100.0
		status.Progress = progress
		status.ProcessedItems = processedCount
		status.FailedItems = failedCount
		s.updateStatus(status)
	}

	// Publish batch to Kafka
	if len(trafficIDs) > 0 {
		if err := s.publishBatchToKafka(ctx, batch); err != nil {
			s.logger.Error("Failed to publish batch to Kafka", "error", err, "batch_id", batch.ID)
			errors = append(errors, fmt.Sprintf("Kafka publish: %s", err.Error()))
		}
	}

	// Update response
	if len(errors) == 0 {
		response.Status = "success"
	} else if len(errors) == batch.Count {
		response.Status = "failed"
	} else {
		response.Status = "partial_success"
	}

	response.TrafficIDs = trafficIDs
	response.Errors = errors
	response.ProcessingTime = time.Since(startTime)

	// Update status
	status.Status = response.Status
	now := time.Now()
	status.EndTime = &now
	s.updateStatus(status)

	// Update statistics
	s.updateStats(nil, time.Since(startTime), errors)

	s.logger.Info("Batch ingestion completed", "batch_id", batch.ID, "processed", processedCount, "failed", failedCount)
	return response, nil
}

func (s *DataIngestionService) GetIngestionStatus(ctx context.Context, id string) (*models.IngestionStatus, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status, exists := s.statusMap[id]
	if !exists {
		return nil, fmt.Errorf("ingestion status not found: %s", id)
	}

	return status, nil
}

func (s *DataIngestionService) GetIngestionStats(ctx context.Context, timeRange *models.TimeRange) (*models.IngestionStats, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return current stats (in a real implementation, this would query the database)
	return s.stats, nil
}

func (s *DataIngestionService) UpdateConfiguration(ctx context.Context, configType string, config interface{}) error {
	s.logger.Info("Updating configuration", "type", configType)

	switch configType {
	case "parser":
		return s.parserService.UpdateConfiguration(ctx, config)
	case "normalizer":
		return s.normalizerService.UpdateConfiguration(ctx, config)
	case "ingestion":
		// Update ingestion configuration
		if cfg, ok := config.(map[string]interface{}); ok {
			if batchSize, exists := cfg["batch_size"]; exists {
				if bs, ok := batchSize.(int); ok {
					s.config.Ingestion.BatchSize = bs
				}
			}
			if batchTimeout, exists := cfg["batch_timeout"]; exists {
				if bt, ok := batchTimeout.(time.Duration); ok {
					s.config.Ingestion.BatchTimeout = bt
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown configuration type: %s", configType)
	}
}

// Helper methods

func (s *DataIngestionService) processTrafficData(request *models.IngestionRequest) (*models.TrafficData, error) {
	// Convert request data to TrafficData
	trafficData := &models.TrafficData{
		ID:        uuid.New().String(),
		Timestamp: request.Timestamp,
		Priority:  request.Priority,
		Tags:      request.Tags,
		Status:    "processing",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Extract traffic data from request
	if data, ok := request.Data.(map[string]interface{}); ok {
		// Extract basic fields
		if sourceIP, exists := data["source_ip"]; exists {
			trafficData.SourceIP = fmt.Sprintf("%v", sourceIP)
		}
		if destIP, exists := data["destination_ip"]; exists {
			trafficData.DestinationIP = fmt.Sprintf("%v", destIP)
		}
		if method, exists := data["method"]; exists {
			trafficData.Method = fmt.Sprintf("%v", method)
		}
		if url, exists := data["url"]; exists {
			trafficData.URL = fmt.Sprintf("%v", url)
		}
		if path, exists := data["path"]; exists {
			trafficData.Path = fmt.Sprintf("%v", path)
		}
		if statusCode, exists := data["status_code"]; exists {
			if sc, ok := statusCode.(int); ok {
				trafficData.StatusCode = sc
			}
		}
		if userAgent, exists := data["user_agent"]; exists {
			trafficData.UserAgent = fmt.Sprintf("%v", userAgent)
		}
		if contentType, exists := data["content_type"]; exists {
			trafficData.ContentType = fmt.Sprintf("%v", contentType)
		}

		// Extract headers
		if headers, exists := data["headers"]; exists {
			if h, ok := headers.(map[string]interface{}); ok {
				trafficData.Headers = make(map[string]string)
				for k, v := range h {
					trafficData.Headers[k] = fmt.Sprintf("%v", v)
				}
			}
		}

		// Extract query parameters
		if queryParams, exists := data["query_params"]; exists {
			if qp, ok := queryParams.(map[string]interface{}); ok {
				trafficData.QueryParams = make(map[string]string)
				for k, v := range qp {
					trafficData.QueryParams[k] = fmt.Sprintf("%v", v)
				}
			}
		}

		// Extract body
		if body, exists := data["body"]; exists {
			if bodyStr, ok := body.(string); ok {
				trafficData.BodyText = bodyStr
				trafficData.Body = []byte(bodyStr)
			}
		}

		// Extract metadata
		if metadata, exists := data["metadata"]; exists {
			if m, ok := metadata.(map[string]interface{}); ok {
				trafficData.Metadata = m
			}
		}
	}

	// Validate traffic data
	if err := s.validateTrafficData(trafficData); err != nil {
		return nil, fmt.Errorf("traffic data validation failed: %w", err)
	}

	return trafficData, nil
}

func (s *DataIngestionService) validateTrafficData(trafficData *models.TrafficData) error {
	if trafficData.Timestamp.IsZero() {
		trafficData.Timestamp = time.Now()
	}

	if trafficData.Method == "" {
		return fmt.Errorf("method is required")
	}

	if trafficData.URL == "" && trafficData.Path == "" {
		return fmt.Errorf("either URL or path is required")
	}

	return nil
}

func (s *DataIngestionService) publishToKafka(ctx context.Context, trafficData *models.TrafficData) error {
	// Determine topic based on traffic type
	topic := s.config.Ingestion.Topics.APITraffic

	// Serialize traffic data
	data, err := json.Marshal(trafficData)
	if err != nil {
		return fmt.Errorf("failed to marshal traffic data: %w", err)
	}

	// Publish to Kafka
	message := kafka.Message{
		Topic: topic,
		Key:   []byte(trafficData.ID),
		Value: data,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce message to Kafka: %w", err)
	}

	s.logger.Info("Published traffic data to Kafka", "traffic_id", trafficData.ID, "topic", topic)
	return nil
}

func (s *DataIngestionService) publishBatchToKafka(ctx context.Context, batch *models.BatchTrafficData) error {
	// Serialize batch data
	data, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal batch data: %w", err)
	}

	// Publish to Kafka
	message := kafka.Message{
		Topic: s.config.Ingestion.Topics.APITraffic,
		Key:   []byte(batch.ID),
		Value: data,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		return fmt.Errorf("failed to produce batch message to Kafka: %w", err)
	}

	s.logger.Info("Published batch data to Kafka", "batch_id", batch.ID, "count", batch.Count)
	return nil
}

func (s *DataIngestionService) updateStatus(status *models.IngestionStatus) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.statusMap[status.ID] = status
}

func (s *DataIngestionService) updateStats(trafficData *models.TrafficData, processingTime time.Duration, errors []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats.TotalIngested++
	if len(errors) == 0 {
		s.stats.TotalParsed++
		s.stats.TotalNormalized++
		s.stats.SuccessRate = float64(s.stats.TotalParsed) / float64(s.stats.TotalIngested)
	} else {
		s.stats.TotalErrors++
		s.stats.ErrorRate = float64(s.stats.TotalErrors) / float64(s.stats.TotalIngested)
	}

	s.stats.AverageProcessingTime = (s.stats.AverageProcessingTime + processingTime) / 2
	s.stats.LastIngestion = time.Now()

	if len(errors) > 0 {
		s.stats.LastError = time.Now()
	}
}

func (s *DataIngestionService) startStatsCollector() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		// In a real implementation, this would persist stats to database
		s.mutex.Unlock()
	}
}

func (s *DataIngestionService) startStatusCleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		cutoff := time.Now().Add(-24 * time.Hour) // Keep status for 24 hours
		for id, status := range s.statusMap {
			if status.EndTime != nil && status.EndTime.Before(cutoff) {
				delete(s.statusMap, id)
			}
		}
		s.mutex.Unlock()
	}
} 