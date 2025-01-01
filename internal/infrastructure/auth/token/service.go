package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"github.com/google/uuid"
)

// Service implements the domain.TokenService interface
type Service struct {
	config     services.TokenConfig
	cache      services.CacheService
	keyManager KeyManager
}

// NewService creates a new token service
func NewService(config services.TokenConfig, cache services.CacheService, keyManager KeyManager) *Service {
	return &Service{
		config:     config,
		cache:      cache,
		keyManager: keyManager,
	}
}

// generateToken creates a new JWT token
func (s *Service) generateToken(ctx context.Context, claims services.TokenClaims, duration time.Duration) (string, error) {
	now := time.Now()
	jwtClaims := jwt.MapClaims{
		"user_id":    claims.UserID.String(),
		"email":      claims.Email,
		"username":   claims.Username,
		"token_type": string(claims.TokenType),
		"iat":        now.Unix(),
		"exp":        now.Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	key, err := s.keyManager.GetSigningKey(ctx, claims.TokenType)
	if err != nil {
		return "", fmt.Errorf("failed to get signing key: %w", err)
	}

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateAccessToken generates a new access token
func (s *Service) GenerateAccessToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.AccessTokenDuration)
}

// GenerateRefreshToken generates a new refresh token
func (s *Service) GenerateRefreshToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.RefreshTokenDuration)
}

// GenerateResetToken generates a password reset token
func (s *Service) GenerateResetToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.ResetTokenDuration)
}

// GenerateVerificationToken generates an email verification token
func (s *Service) GenerateVerificationToken(ctx context.Context, claims services.TokenClaims) (string, error) {
	return s.generateToken(ctx, claims, s.config.VerificationTokenDuration)
}

// ValidateToken validates a token and returns its claims
func (s *Service) ValidateToken(ctx context.Context, tokenString string, tokenType services.TokenType) (*services.TokenClaims, error) {
	// Check if token is revoked
	isRevoked, err := s.IsTokenRevoked(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token revocation: %w", err)
	}
	if isRevoked {
		return nil, fmt.Errorf("token is revoked")
	}

	key, err := s.keyManager.GetSigningKey(ctx, tokenType)
	if err != nil {
		return nil, fmt.Errorf("failed to get signing key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate token type
	if tokenTypeStr, ok := claims["token_type"].(string); !ok || services.TokenType(tokenTypeStr) != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	// Parse user_id into UUID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	return &services.TokenClaims{
		UserID:    userID,
		Email:     claims["email"].(string),
		Username:  claims["username"].(string),
		TokenType: tokenType,
	}, nil
}

// RevokeToken revokes a token
func (s *Service) RevokeToken(ctx context.Context, token string) error {
	// Store the token in the blacklist with an expiration
	err := s.cache.Set(ctx, fmt.Sprintf("revoked_token:%s", token), true, s.config.AccessTokenDuration)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (s *Service) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	var isRevoked bool
	err := s.cache.Get(ctx, fmt.Sprintf("revoked_token:%s", token), &isRevoked)
	if err != nil {
		return false, fmt.Errorf("failed to check token revocation: %w", err)
	}
	return isRevoked, nil
}
