package services

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong     = errors.New("password must be at most 72 characters long")
	ErrPasswordNoUpper     = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower     = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial   = errors.New("password must contain at least one special character")
	ErrPasswordHashFailed  = errors.New("failed to hash password")
	ErrPasswordInvalid     = errors.New("invalid password")
)

// PasswordService handles password-related operations
type PasswordService struct{}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// ValidatePassword validates a password
func (s *PasswordService) ValidatePassword(ctx context.Context, password string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 72 {
		return ErrPasswordTooLong
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

// HashPassword hashes a password
func (s *PasswordService) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrPasswordHashFailed
	}
	return string(hash), nil
}

// VerifyPassword verifies if a password matches its hash
func (s *PasswordService) VerifyPassword(ctx context.Context, password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrPasswordInvalid
	}
	return nil
}

// GenerateRandomPassword generates a random password
func (s *PasswordService) GenerateRandomPassword(ctx context.Context) (string, error) {
	const (
		upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars   = "abcdefghijklmnopqrstuvwxyz"
		numberChars  = "0123456789"
		specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		length       = 16
	)

	// Create a password with at least one of each required character type
	password := make([]byte, length)
	
	// Add one of each required character type
	password[0] = upperChars[secureRand(len(upperChars))]
	password[1] = lowerChars[secureRand(len(lowerChars))]
	password[2] = numberChars[secureRand(len(numberChars))]
	password[3] = specialChars[secureRand(len(specialChars))]

	// Fill the rest with random characters from all types
	allChars := upperChars + lowerChars + numberChars + specialChars
	for i := 4; i < length; i++ {
		password[i] = allChars[secureRand(len(allChars))]
	}

	// Shuffle the password to make it more random
	for i := len(password) - 1; i > 0; i-- {
		j := secureRand(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// secureRand returns a secure random number in range [0, n)
func secureRand(n int) int {
	if n <= 0 {
		panic("invalid input to secureRand")
	}
	
	const maxUint = ^uint(0)
	const maxInt = int(maxUint >> 1)
	
	// Calculate the threshold for rejection sampling
	threshold := maxInt - (maxInt % n)
	
	for {
		result, err := rand.Int(rand.Reader, big.NewInt(int64(threshold)))
		if err != nil {
			continue
		}
		
		if r := int(result.Int64()); r < threshold {
			return r % n
		}
	}
}
