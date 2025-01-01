package services

import "context"

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	// PublishUserEvent publishes user-related events
	PublishUserEvent(ctx context.Context, eventType string, payload interface{}) error
}
