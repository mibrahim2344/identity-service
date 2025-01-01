package events

import "context"

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	// PublishUserRegistered publishes a UserRegisteredEvent
	PublishUserRegistered(ctx context.Context, event UserRegisteredEvent) error

	// PublishUserEmailVerified publishes a UserEmailVerifiedEvent
	PublishUserEmailVerified(ctx context.Context, event UserEmailVerifiedEvent) error

	// PublishPasswordResetRequested publishes a UserPasswordResetRequestedEvent
	PublishPasswordResetRequested(ctx context.Context, event UserPasswordResetRequestedEvent) error

	// PublishPasswordChanged publishes a UserPasswordChangedEvent
	PublishPasswordChanged(ctx context.Context, event UserPasswordChangedEvent) error
}

// EventPublisherFunc is a function type that implements EventPublisher
type EventPublisherFunc func(ctx context.Context, eventType string, event interface{}) error

// NoOpEventPublisher is a no-op implementation of EventPublisher
type NoOpEventPublisher struct{}

func (p *NoOpEventPublisher) PublishUserRegistered(ctx context.Context, event UserRegisteredEvent) error {
	return nil
}

func (p *NoOpEventPublisher) PublishUserEmailVerified(ctx context.Context, event UserEmailVerifiedEvent) error {
	return nil
}

func (p *NoOpEventPublisher) PublishPasswordResetRequested(ctx context.Context, event UserPasswordResetRequestedEvent) error {
	return nil
}

func (p *NoOpEventPublisher) PublishPasswordChanged(ctx context.Context, event UserPasswordChangedEvent) error {
	return nil
}
