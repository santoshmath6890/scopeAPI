package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Config holds Kafka configuration
type Config struct {
	Brokers        []string
	TopicPrefix    string
	ProducerConfig ProducerConfig
	ConsumerConfig ConsumerConfig
}

// ProducerConfig holds Kafka producer configuration
type ProducerConfig struct {
	Acks           string
	Retries        int
	BatchSize      int
	BatchTimeout   time.Duration
	Compression    string
	MaxMessageSize int
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	GroupID           string
	AutoOffsetReset   string
	SessionTimeout    time.Duration
	HeartbeatInterval time.Duration
	MaxPollRecords    int
	MaxPollInterval   time.Duration
}

// Message represents a Kafka message
type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
	Time      time.Time
}

// Producer represents a Kafka producer
type Producer struct {
	writer *kafka.Writer
	config Config
}

// Consumer represents a Kafka consumer
type Consumer struct {
	reader *kafka.Reader
	config Config
}

// ProducerInterface defines the interface for a Kafka producer
// Used by services for dependency injection and mocking
type ProducerInterface interface {
	Produce(ctx context.Context, message Message) error
	Close() error
}

// NewProducer creates a new Kafka producer
func NewProducer(config Config) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		BatchSize:    config.ProducerConfig.BatchSize,
		BatchTimeout: config.ProducerConfig.BatchTimeout,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &Producer{
		writer: writer,
		config: config,
	}, nil
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config Config, topics []string) (*Consumer, error) {
	if len(topics) == 0 {
		return nil, fmt.Errorf("at least one topic must be specified")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         config.Brokers,
		GroupID:         config.ConsumerConfig.GroupID,
		Topic:           topics[0], // For simplicity, use first topic
		MinBytes:        10e3,      // 10KB
		MaxBytes:        10e6,      // 10MB
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	})

	return &Consumer{
		reader: reader,
		config: config,
	}, nil
}

// SendMessage sends a message to a Kafka topic
func (p *Producer) SendMessage(topic string, key []byte, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	return p.writer.WriteMessages(context.Background(), msg)
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consume consumes messages from Kafka
func (c *Consumer) Consume(ctx context.Context, batchSize int) ([]Message, error) {
	var messages []Message

	for i := 0; i < batchSize; i++ {
		select {
		case <-ctx.Done():
			return messages, ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return messages, err
			}

			messages = append(messages, Message{
				Topic:     msg.Topic,
				Partition: msg.Partition,
				Offset:    msg.Offset,
				Key:       msg.Key,
				Value:     msg.Value,
				Time:      msg.Time,
			})
		}
	}

	return messages, nil
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// SendJSONMessage sends a JSON message to a Kafka topic
func (p *Producer) SendJSONMessage(topic string, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return p.SendMessage(topic, []byte(key), jsonData)
}

// Implement Produce for Producer to satisfy ProducerInterface
func (p *Producer) Produce(ctx context.Context, message Message) error {
	msg := kafka.Message{
		Topic: message.Topic,
		Key:   message.Key,
		Value: message.Value,
	}
	return p.writer.WriteMessages(ctx, msg)
}
