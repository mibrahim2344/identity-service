package events

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	// User-related events
	UserRegistered     EventType = "user.registered"
	UserVerified       EventType = "user.verified"
	UserPasswordReset  EventType = "user.password.reset"
	UserPasswordChange EventType = "user.password.changed"
	UserDeleted        EventType = "user.deleted"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// UserRegisteredEvent is published when a new user registers
type UserRegisteredEvent struct {
	BaseEvent
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Locale    string    `json:"locale"`
}

// UserVerifiedEvent is published when a user verifies their email
type UserVerifiedEvent struct {
	BaseEvent
	UserID  uuid.UUID `json:"userId"`
	Email   string    `json:"email"`
}

// UserPasswordResetEvent is published when a password reset is requested
type UserPasswordResetEvent struct {
	BaseEvent
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	ResetLink string    `json:"resetLink"`
}

// UserPasswordChangedEvent is published when a password is changed
type UserPasswordChangedEvent struct {
	BaseEvent
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}

// UserDeletedEvent is published when a user is deleted
type UserDeletedEvent struct {
	BaseEvent
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType EventType) BaseEvent {
	return BaseEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}
}

// NewUserRegisteredEvent creates a new user registered event
func NewUserRegisteredEvent(userID uuid.UUID, email, username, firstName, lastName string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: NewBaseEvent(UserRegistered),
		UserID:    userID,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Locale:    "en", // Default locale, could be made configurable
	}
}

// NewUserVerifiedEvent creates a new user verified event
func NewUserVerifiedEvent(userID uuid.UUID, email string) *UserVerifiedEvent {
	return &UserVerifiedEvent{
		BaseEvent: NewBaseEvent(UserVerified),
		UserID:    userID,
		Email:     email,
	}
}

// NewUserPasswordResetEvent creates a new password reset event
func NewUserPasswordResetEvent(userID uuid.UUID, email, resetLink string) *UserPasswordResetEvent {
	return &UserPasswordResetEvent{
		BaseEvent: NewBaseEvent(UserPasswordReset),
		UserID:    userID,
		Email:     email,
		ResetLink: resetLink,
	}
}

// NewUserPasswordChangedEvent creates a new password changed event
func NewUserPasswordChangedEvent(userID uuid.UUID, email string) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseEvent: NewBaseEvent(UserPasswordChange),
		UserID:    userID,
		Email:     email,
	}
}

// NewUserDeletedEvent creates a new user deleted event
func NewUserDeletedEvent(userID uuid.UUID, email string) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: NewBaseEvent(UserDeleted),
		UserID:    userID,
		Email:     email,
	}
}
