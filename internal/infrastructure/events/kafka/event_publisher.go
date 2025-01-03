package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

// Ensure EventPublisher implements services.EventPublisher
var _ services.EventPublisher = (*EventPublisher)(nil)

// EventPublisher implements the domain.EventPublisher interface using Kafka
type EventPublisher struct {
	writer *kafka.Writer
}

// NewEventPublisher creates a new Kafka event publisher
func NewEventPublisher(brokers []string, topic string) *EventPublisher {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &EventPublisher{
		writer: writer,
	}
}

// PublishUserEvent publishes a user-related event
func (p *EventPublisher) PublishUserEvent(ctx context.Context, eventType string, event interface{}) error {
	// Marshal event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(eventType),
		Value: eventJSON,
	}

	// Write message to Kafka
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close closes the Kafka writer
func (p *EventPublisher) Close() error {
	return p.writer.Close()
}
