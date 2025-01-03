package password

import (
	"context"
	"fmt"
	"regexp"
	"unicode"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

// Service implements the domain.PasswordService interface
type Service struct {
	hasher  PasswordHasher
	config  services.PasswordConfig
	entropy EntropyProvider
}

// NewService creates a new password service
func NewService(hasher PasswordHasher, config services.PasswordConfig, entropy EntropyProvider) *Service {
	if entropy == nil {
		entropy = &CryptoEntropyProvider{}
	}
	
	return &Service{
		hasher:  hasher,
		config:  config,
		entropy: entropy,
	}
}

// HashPassword generates a hash for the given password
func (s *Service) HashPassword(ctx context.Context, password string) (string, error) {
	if err := s.ValidatePassword(ctx, password); err != nil {
		return "", fmt.Errorf("password validation failed: %w", err)
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

// VerifyPassword verifies if a password matches its hash
func (s *Service) VerifyPassword(ctx context.Context, password, hash string) error {
	if err := s.hasher.Verify(password, hash); err != nil {
		return fmt.Errorf("password verification failed: %w", err)
	}
	return nil
}

// GenerateRandomPassword generates a random password
func (s *Service) GenerateRandomPassword(ctx context.Context) (string, error) {
	const (
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		digits       = "0123456789"
		specials     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	// Calculate the minimum length needed to satisfy all requirements
	minLength := 0
	if s.config.RequireUppercase {
		minLength++
	}
	if s.config.RequireLowercase {
		minLength++
	}
	if s.config.RequireNumbers {
		minLength++
	}
	if s.config.RequireSpecialChars {
		minLength++
	}

	// Use the larger of minLength and config.MinLength
	if s.config.MinLength > minLength {
		minLength = s.config.MinLength
	}

	// Generate random bytes
	length := minLength + 4 // Add some extra length for better entropy
	randomBytes := make([]byte, length)
	if _, err := s.entropy.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	var password []byte
	
	// Ensure we have at least one of each required character type
	if s.config.RequireUppercase {
		idx := randomBytes[0] % byte(len(upperLetters))
		password = append(password, upperLetters[idx])
	}
	if s.config.RequireLowercase {
		idx := randomBytes[1] % byte(len(lowerLetters))
		password = append(password, lowerLetters[idx])
	}
	if s.config.RequireNumbers {
		idx := randomBytes[2] % byte(len(digits))
		password = append(password, digits[idx])
	}
	if s.config.RequireSpecialChars {
		idx := randomBytes[3] % byte(len(specials))
		password = append(password, specials[idx])
	}

	// Fill the rest with random characters from all allowed sets
	allChars := ""
	if s.config.RequireUppercase {
		allChars += upperLetters
	}
	if s.config.RequireLowercase {
		allChars += lowerLetters
	}
	if s.config.RequireNumbers {
		allChars += digits
	}
	if s.config.RequireSpecialChars {
		allChars += specials
	}

	for i := len(password); i < length; i++ {
		idx := randomBytes[i] % byte(len(allChars))
		password = append(password, allChars[idx])
	}

	// Shuffle the password
	for i := len(password) - 1; i > 0; i-- {
		j := randomBytes[i] % byte(i+1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// ValidatePassword validates password strength
func (s *Service) ValidatePassword(ctx context.Context, password string) error {
	if len(password) < s.config.MinLength {
		return fmt.Errorf("password must be at least %d characters long", s.config.MinLength)
	}

	if len(password) > s.config.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", s.config.MaxLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

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

	if s.config.RequireUppercase && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if s.config.RequireLowercase && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if s.config.RequireNumbers && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if s.config.RequireSpecialChars && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for common patterns that might make the password weak
	commonPatterns := []string{
		`\d{4}`,                // Four digits in a row
		`(?i)password`,         // The word "password" (case insensitive)
		`(?i)admin`,            // The word "admin" (case insensitive)
		`(.)\1{2,}`,           // Same character repeated 3 or more times
		`(?i)(qwerty|asdf)`,   // Common keyboard patterns
	}

	for _, pattern := range commonPatterns {
		if matched, _ := regexp.MatchString(pattern, password); matched {
			return fmt.Errorf("password contains common unsafe pattern")
		}
	}

	return nil
}
