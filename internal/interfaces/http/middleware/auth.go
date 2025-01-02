package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"go.uber.org/zap"
)

// AuthMiddleware handles authentication for protected routes
type AuthMiddleware struct {
	tokenService   services.TokenService
	metricsService services.MetricsService
	logger         *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(tokenService services.TokenService, metricsService services.MetricsService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService:   tokenService,
		metricsService: metricsService,
		logger:         logger,
	}
}

// Custom type for context keys
type contextKey string

const (
	userIDKey contextKey = "user_id"
)

// Authenticate verifies the JWT token and adds user information to the context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := m.tokenService.ValidateToken(r.Context(), token, services.TokenTypeAccess)
		if err != nil {
			m.logger.Error("invalid token", zap.Error(err))
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
