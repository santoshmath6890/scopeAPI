package services

import (
	"context"
	"fmt"
	"sync"
	"time"
	"encoding/json"

	"scopeapi.local/backend/services/data-ingestion/internal/config"
	"scopeapi.local/backend/services/data-ingestion/internal/models"
	"scopeapi.local/backend/shared/logging"
	"scopeapi.local/backend/shared/messaging/kafka"
)

type QueueServiceInterface interface {
	Enqueue(ctx context.Context, data *models.TrafficData) error
	Dequeue(ctx context.Context) (*models.TrafficData, error)
	Flush(ctx context.Context) error
	GetQueueStatus(ctx context.Context) (*models.QueueStatus, error)
	GetTopics(ctx context.Context) ([]models.TopicInfo, error)
	UpdateConfiguration(ctx context.Context, config interface{}) error
}

type QueueService struct {
	logger        logging.Logger
	config        *config.Config
	kafkaProducer kafka.ProducerInterface
	queue         []*models.TrafficData
	maxSize       int
	mutex         sync.Mutex
	lastFlush     time.Time
	flushInterval time.Duration
}

func NewQueueService(kafkaProducer kafka.ProducerInterface, logger logging.Logger, cfg *config.Config) QueueServiceInterface {
	return &QueueService{
		logger:        logger,
		config:        cfg,
		kafkaProducer: kafkaProducer,
		queue:         make([]*models.TrafficData, 0, cfg.Queue.MaxSize),
		maxSize:       cfg.Queue.MaxSize,
		flushInterval: cfg.Queue.FlushInterval,
		lastFlush:     time.Now(),
	}
}

func (s *QueueService) Enqueue(ctx context.Context, data *models.TrafficData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.queue) >= s.maxSize {
		return fmt.Errorf("queue is full")
	}

	s.queue = append(s.queue, data)
	s.logger.Info("Enqueued traffic data", "traffic_id", data.ID, "queue_size", len(s.queue))
	return nil
}

func (s *QueueService) Dequeue(ctx context.Context) (*models.TrafficData, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.queue) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	item := s.queue[0]
	s.queue = s.queue[1:]
	s.logger.Info("Dequeued traffic data", "traffic_id", item.ID, "queue_size", len(s.queue))
	return item, nil
}

func (s *QueueService) Flush(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.queue) == 0 {
		s.logger.Info("Queue is empty, nothing to flush")
		return nil
	}

	// Batch publish to Kafka
	batch := &models.BatchTrafficData{
		ID:        fmt.Sprintf("batch-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Count:     len(s.queue),
		Data:      make([]models.TrafficData, len(s.queue)),
		Status:    "flushed",
	}
	for i, item := range s.queue {
		batch.Data[i] = *item
	}

	data, err := json.Marshal(batch)
	if err != nil {
		s.logger.Error("Failed to marshal batch for Kafka", "error", err)
		return err
	}

	message := kafka.Message{
		Topic: s.config.Ingestion.Topics.APITraffic,
		Key:   []byte(batch.ID),
		Value: data,
	}

	if err := s.kafkaProducer.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to produce batch to Kafka", "error", err)
		return err
	}

	s.logger.Info("Flushed queue to Kafka", "batch_id", batch.ID, "count", batch.Count)
	s.queue = s.queue[:0]
	s.lastFlush = time.Now()
	return nil
}

func (s *QueueService) GetQueueStatus(ctx context.Context) (*models.QueueStatus, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	status := &models.QueueStatus{
		Name:          "traffic_queue",
		Size:          len(s.queue),
		MaxSize:       s.maxSize,
		Pending:       len(s.queue),
		Processing:    0,
		Completed:     0,
		Failed:        0,
		RetryCount:    0,
		LastFlush:     s.lastFlush,
		NextFlush:     s.lastFlush.Add(s.flushInterval),
		FlushInterval: s.flushInterval,
		Health:        "healthy",
		Errors:        []string{},
		Metrics:       map[string]interface{}{},
	}
	return status, nil
}

func (s *QueueService) GetTopics(ctx context.Context) ([]models.TopicInfo, error) {
	// In a real implementation, this would query Kafka for topic metadata
	topics := []models.TopicInfo{
		{
			Name:       s.config.Ingestion.Topics.APITraffic,
			Partitions: 1,
			Replicas:   1,
			Config:     map[string]string{"retention.ms": "86400000"},
			Messages:   0,
			Size:       0,
			ConsumerGroups: []string{"default"},
			Health:     "unknown",
			LastMessage: time.Now(),
		},
	}
	return topics, nil
}

func (s *QueueService) UpdateConfiguration(ctx context.Context, config interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if cfg, ok := config.(map[string]interface{}); ok {
		if maxSize, exists := cfg["max_size"]; exists {
			if ms, ok := maxSize.(int); ok {
				s.maxSize = ms
			}
		}
		if flushInterval, exists := cfg["flush_interval"]; exists {
			if fi, ok := flushInterval.(time.Duration); ok {
				s.flushInterval = fi
			}
		}
	}

	s.logger.Info("Queue configuration updated")
	return nil
} 