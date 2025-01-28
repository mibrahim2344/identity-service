package errors

import "errors"

var (
	// ErrUserNotFound indicates that the requested user was not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists indicates that a user with the given email or username already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials indicates that the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken indicates that the provided token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired token")

	// ErrUnauthorized indicates that the user is not authorized to perform the action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidInput indicates that the provided input is invalid
	ErrInvalidInput = errors.New("invalid input")
)

// DomainError represents a domain-specific error with operation context
type DomainError struct {
	Op  string
	Err error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err == nil {
		return e.Op
	}
	return e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// WrapError wraps an error with operation context
func WrapError(op string, err error) error {
	return &DomainError{Op: op, Err: err}
}

// Is implements error matching for wrapped errors
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Err, target)
}
