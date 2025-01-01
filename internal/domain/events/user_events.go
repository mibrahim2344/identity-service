package events

import (
	"time"

	"github.com/google/uuid"
)

// UserEvent represents a base event for user-related domain events
type UserEvent struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// UserRegisteredEvent is emitted when a new user is registered
type UserRegisteredEvent struct {
	UserEvent
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UserEmailVerifiedEvent is emitted when a user verifies their email
type UserEmailVerifiedEvent struct {
	UserEvent
}

// UserPasswordResetRequestedEvent is emitted when a user requests a password reset
type UserPasswordResetRequestedEvent struct {
	UserEvent
	ResetToken string `json:"reset_token"`
}

// UserPasswordChangedEvent is emitted when a user changes their password
type UserPasswordChangedEvent struct {
	UserEvent
}

// UserDeletedEvent is emitted when a user is deleted
type UserDeletedEvent struct {
	UserEvent
}

// NewUserEvent creates a new base user event
func NewUserEvent(userID uuid.UUID) UserEvent {
	return UserEvent{
		ID:        uuid.New().String(),
		UserID:    userID,
		Timestamp: time.Now(),
	}
}

// NewUserRegisteredEvent creates a new user registered event
func NewUserRegisteredEvent(userID uuid.UUID, email, username string) UserRegisteredEvent {
	return UserRegisteredEvent{
		UserEvent: NewUserEvent(userID),
		Email:     email,
		Username:  username,
	}
}

// NewUserEmailVerifiedEvent creates a new email verified event
func NewUserEmailVerifiedEvent(userID uuid.UUID) UserEmailVerifiedEvent {
	return UserEmailVerifiedEvent{
		UserEvent: NewUserEvent(userID),
	}
}

// NewUserPasswordResetRequestedEvent creates a new password reset requested event
func NewUserPasswordResetRequestedEvent(userID uuid.UUID, resetToken string) UserPasswordResetRequestedEvent {
	return UserPasswordResetRequestedEvent{
		UserEvent:  NewUserEvent(userID),
		ResetToken: resetToken,
	}
}

// NewUserPasswordChangedEvent creates a new password changed event
func NewUserPasswordChangedEvent(userID uuid.UUID) UserPasswordChangedEvent {
	return UserPasswordChangedEvent{
		UserEvent: NewUserEvent(userID),
	}
}

// NewUserDeletedEvent creates a new user deleted event
func NewUserDeletedEvent(userID uuid.UUID) UserDeletedEvent {
	return UserDeletedEvent{
		UserEvent: NewUserEvent(userID),
	}
}
