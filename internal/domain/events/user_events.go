package events

import (
	"github.com/google/uuid"
)

// UserEmailVerifiedEvent is emitted when a user verifies their email
type UserEmailVerifiedEvent struct {
	BaseEvent
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}

// UserPasswordResetRequestedEvent is emitted when a user requests a password reset
type UserPasswordResetRequestedEvent struct {
	BaseEvent
	UserID     uuid.UUID `json:"userId"`
	Email      string    `json:"email"`
	ResetToken string    `json:"resetToken"`
}

// NewUserEmailVerifiedEvent creates a new email verified event
func NewUserEmailVerifiedEvent(userID uuid.UUID, email string) UserEmailVerifiedEvent {
	event := UserEmailVerifiedEvent{
		BaseEvent: NewBaseEvent("UserVerified"),
		UserID:    userID,
		Email:     email,
	}
	return event
}

// NewUserPasswordResetRequestedEvent creates a new password reset requested event
func NewUserPasswordResetRequestedEvent(userID uuid.UUID, email, resetToken string) UserPasswordResetRequestedEvent {
	event := UserPasswordResetRequestedEvent{
		BaseEvent:  NewBaseEvent("UserPasswordReset"),
		UserID:     userID,
		Email:      email,
		ResetToken: resetToken,
	}
	return event
}
