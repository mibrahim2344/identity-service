package services

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrConflict is returned when there is a conflict with existing resource
	ErrConflict = errors.New("resource conflict")

	// ErrCacheKeyNotFound is returned when a key is not found in the cache
	ErrCacheKeyNotFound = errors.New("cache key not found")

	// ErrInvalidCacheValue is returned when a cached value cannot be unmarshaled
	ErrInvalidCacheValue = errors.New("invalid cache value")

	// ErrCacheOperationFailed is returned when a cache operation fails
	ErrCacheOperationFailed = errors.New("cache operation failed")

	// ErrCacheConnectionFailed is returned when the connection to the cache fails
	ErrCacheConnectionFailed = errors.New("cache connection failed")

	// ErrAuthentication is returned when authentication fails
	ErrAuthentication = errors.New("authentication failed")

	// ErrEmailAlreadyExists is returned when attempting to use an email that is already registered
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrUsernameAlreadyExists is returned when attempting to use a username that is already taken
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrUserAlreadyExists is returned when a user with the same email or username already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrTokenRevoked is returned when attempting to use a revoked token
	ErrTokenRevoked = errors.New("token has been revoked")
)

// IsNotFoundError checks if the given error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// NewConflictError returns a conflict error with a custom message
func NewConflictError(msg string) error {
	return fmt.Errorf("%w: %s", ErrConflict, msg)
}

// NewAuthenticationError returns an authentication error with a custom message
func NewAuthenticationError(msg string) error {
	return fmt.Errorf("%w: %s", ErrAuthentication, msg)
}
