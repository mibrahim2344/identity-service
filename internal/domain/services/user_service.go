package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mibrahim2344/identity-service/internal/domain/models"
)

// RegisterUserInput represents the input for user registration
type RegisterUserInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
	Role      models.Role
}

// UpdateUserInput represents the input for updating user details
type UpdateUserInput struct {
	Email    string
	Username string
	Status   models.UserStatus
	Role     models.Role
}

// LoginUserInput represents the input for user login
type LoginUserInput struct {
	Email    string
	Username string
	Password string
}

// LoginResponse represents the response for a successful login
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

// ResetPasswordInput represents the input for password reset
type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

// UserService defines the interface for user-related business operations
type UserService interface {
	// RegisterUser registers a new user
	RegisterUser(ctx context.Context, input RegisterUserInput) (*models.User, error)

	// AuthenticateUser authenticates a user with email/username and password
	AuthenticateUser(ctx context.Context, emailOrUsername, password string) (*models.User, error)

	// GetUser retrieves a user by their ID
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)

	// UpdateUser updates user details
	UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*models.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error

	// RequestPasswordReset initiates a password reset process
	RequestPasswordReset(ctx context.Context, email string) error

	// ResetPassword resets a user's password using a reset token
	ResetPassword(ctx context.Context, token, newPassword string) error

	// VerifyEmail verifies a user's email address
	VerifyEmail(ctx context.Context, token string) error

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
}
