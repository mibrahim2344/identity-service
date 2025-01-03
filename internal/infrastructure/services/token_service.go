package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

// TokenService handles JWT token operations
type TokenService struct {
	config services.TokenConfig
}

// NewTokenService creates a new token service
func NewTokenService(secret string, accessTokenExpiry, refreshTokenExpiry time.Duration) *TokenService {
	return &TokenService{
		config: services.TokenConfig{
			AccessTokenDuration:        accessTokenExpiry,
			RefreshTokenDuration:      refreshTokenExpiry,
			ResetTokenDuration:        24 * time.Hour,    // 24 hours
			VerificationTokenDuration: 72 * time.Hour,    // 72 hours
			SigningKey:               []byte(secret),
		},
	}
}

// GenerateAccessToken generates a new access token
func (s *TokenService) GenerateAccessToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.AccessTokenDuration)
}

// GenerateRefreshToken generates a new refresh token
func (s *TokenService) GenerateRefreshToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.RefreshTokenDuration)
}

// GenerateResetToken generates a password reset token
func (s *TokenService) GenerateResetToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.ResetTokenDuration)
}

// GenerateVerificationToken generates an email verification token
func (s *TokenService) GenerateVerificationToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.VerificationTokenDuration)
}

// ValidateToken validates a token and returns its claims
func (s *TokenService) ValidateToken(ctx context.Context, tokenString string, tokenType services.TokenType) (*services.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.config.SigningKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Validate token type
	if claims["token_type"].(string) != string(tokenType) {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &services.TokenClaims{
		UserID:    userID,
		Email:     claims["email"].(string),
		Username:  claims["username"].(string),
		Role:      claims["role"].(string),
		TokenType: services.TokenType(claims["token_type"].(string)),
	}, nil
}

// RevokeToken revokes a token
func (s *TokenService) RevokeToken(ctx context.Context, token string) error {
	// TODO: Implement token revocation using Redis
	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (s *TokenService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	// TODO: Implement token revocation check using Redis
	return false, nil
}

// generateToken generates a new JWT token
func (s *TokenService) generateToken(ctx context.Context, claims services.TokenClaims, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID.String(),
		"email":      claims.Email,
		"username":   claims.Username,
		"role":       claims.Role,
		"token_type": string(claims.TokenType),
		"exp":        time.Now().Add(duration).Unix(),
	})

	return token.SignedString(s.config.SigningKey)
}
