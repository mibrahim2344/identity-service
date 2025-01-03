package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenType represents the type of token
type TokenType string

const (
	// TokenTypeAccess represents an access token
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh represents a refresh token
	TokenTypeRefresh TokenType = "refresh"
	// TokenTypeReset represents a password reset token
	TokenTypeReset TokenType = "reset"
	// TokenTypeVerification represents an email verification token
	TokenTypeVerification TokenType = "verification"
)

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
}

// TokenService defines the interface for token-related operations
type TokenService interface {
	// GenerateAccessToken generates a new access token
	GenerateAccessToken(ctx context.Context, claims TokenClaims) (string, error)

	// GenerateRefreshToken generates a new refresh token
	GenerateRefreshToken(ctx context.Context, claims TokenClaims) (string, error)

	// GenerateResetToken generates a password reset token
	GenerateResetToken(ctx context.Context, claims TokenClaims) (string, error)

	// GenerateVerificationToken generates an email verification token
	GenerateVerificationToken(ctx context.Context, claims TokenClaims) (string, error)

	// ValidateToken validates a token and returns its claims
	ValidateToken(ctx context.Context, token string, tokenType TokenType) (*TokenClaims, error)

	// RevokeToken revokes a token
	RevokeToken(ctx context.Context, token string) error

	// IsTokenRevoked checks if a token has been revoked
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// TokenConfig represents the configuration for token generation
type TokenConfig struct {
	AccessTokenDuration       time.Duration
	RefreshTokenDuration      time.Duration
	ResetTokenDuration        time.Duration
	VerificationTokenDuration time.Duration
	SigningKey                []byte
}
