package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher defines the interface for password hashing strategies
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

// HashingAlgorithm represents the type of hashing algorithm
type HashingAlgorithm string

const (
	// BCrypt represents the bcrypt hashing algorithm
	BCrypt HashingAlgorithm = "bcrypt"
	// Add more algorithms as needed
)

// NewPasswordHasher creates a new password hasher based on the algorithm
func NewPasswordHasher(algorithm HashingAlgorithm, options map[string]interface{}) (PasswordHasher, error) {
	switch algorithm {
	case BCrypt:
		cost := bcrypt.DefaultCost
		if costVal, ok := options["cost"]; ok {
			if costInt, ok := costVal.(int); ok {
				cost = costInt
			}
		}
		return NewBCryptHasher(cost), nil
	default:
		return nil, fmt.Errorf("unsupported hashing algorithm: %s", algorithm)
	}
}

// BCryptHasher implements PasswordHasher using bcrypt
type BCryptHasher struct {
	cost int
}

// NewBCryptHasher creates a new BCryptHasher
func NewBCryptHasher(cost int) *BCryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &BCryptHasher{
		cost: cost,
	}
}

// Hash generates a bcrypt hash of the password
func (h *BCryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to generate bcrypt hash: %w", err)
	}
	return string(hash), nil
}

// Verify checks if the password matches the hash
func (h *BCryptHasher) Verify(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}
