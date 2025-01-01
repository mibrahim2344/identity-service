package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/mibrahim2344/identity-service/internal/domain/events"
)

const (
	topicUserRegistered          = "user.registered"
	topicUserEmailVerified      = "user.email.verified"
	topicPasswordResetRequested = "user.password.reset.requested"
	topicPasswordChanged        = "user.password.changed"
)

// Publisher implements the domain.EventPublisher interface using Kafka
type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher creates a new Kafka event publisher
func NewPublisher(brokers []string) *Publisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Publisher{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *Publisher) Close() error {
	return p.writer.Close()
}

// PublishUserRegistered publishes a UserRegisteredEvent
func (p *Publisher) PublishUserRegistered(ctx context.Context, event events.UserRegisteredEvent) error {
	return p.publishEvent(ctx, topicUserRegistered, event)
}

// PublishUserEmailVerified publishes a UserEmailVerifiedEvent
func (p *Publisher) PublishUserEmailVerified(ctx context.Context, event events.UserEmailVerifiedEvent) error {
	return p.publishEvent(ctx, topicUserEmailVerified, event)
}

// PublishPasswordResetRequested publishes a UserPasswordResetRequestedEvent
func (p *Publisher) PublishPasswordResetRequested(ctx context.Context, event events.UserPasswordResetRequestedEvent) error {
	return p.publishEvent(ctx, topicPasswordResetRequested, event)
}

// PublishPasswordChanged publishes a UserPasswordChangedEvent
func (p *Publisher) PublishPasswordChanged(ctx context.Context, event events.UserPasswordChangedEvent) error {
	return p.publishEvent(ctx, topicPasswordChanged, event)
}

// PublishUserEvent implements the services.EventPublisher interface
func (p *Publisher) PublishUserEvent(ctx context.Context, eventType string, payload interface{}) error {
	return p.publishEvent(ctx, eventType, payload)
}

// publishEvent is a helper function to publish events to Kafka
func (p *Publisher) publishEvent(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic: topic,
		Value: data,
	}

	return p.writer.WriteMessages(ctx, message)
}
