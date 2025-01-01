package services

import "context"

// PasswordService defines the interface for password-related operations
type PasswordService interface {
	// HashPassword generates a hash for the given password
	HashPassword(ctx context.Context, password string) (string, error)

	// VerifyPassword verifies if a password matches its hash
	VerifyPassword(ctx context.Context, password, hash string) error

	// GenerateRandomPassword generates a random password
	GenerateRandomPassword(ctx context.Context) (string, error)

	// ValidatePassword validates password strength
	ValidatePassword(ctx context.Context, password string) error
}

// PasswordConfig represents the configuration for password operations
type PasswordConfig struct {
	MinLength           int
	RequireUppercase    bool
	RequireLowercase    bool
	RequireNumbers      bool
	RequireSpecialChars bool
	MaxLength           int
}
