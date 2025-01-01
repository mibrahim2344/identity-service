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
