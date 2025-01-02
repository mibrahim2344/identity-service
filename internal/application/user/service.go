package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mibrahim2344/identity-service/internal/domain/events"
	"github.com/mibrahim2344/identity-service/internal/domain/models"
	"github.com/mibrahim2344/identity-service/internal/domain/repositories"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"go.uber.org/zap"
)

// Service implements the domain.UserService interface
type Service struct {
	userRepo        repositories.UserRepository
	passwordService services.PasswordService
	tokenService    services.TokenService
	cacheService    services.CacheService
	eventPublisher  services.EventPublisher
	logger          *zap.Logger
	config          services.CacheConfig
	webAppURL       string
}

// NewService creates a new user service
func NewService(
	userRepo repositories.UserRepository,
	passwordService services.PasswordService,
	tokenService services.TokenService,
	cacheService services.CacheService,
	eventPublisher services.EventPublisher,
	logger *zap.Logger,
	config services.CacheConfig,
	webAppURL string,
) *Service {
	return &Service{
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
		cacheService:    cacheService,
		eventPublisher:  eventPublisher,
		logger:          logger,
		config:          config,
		webAppURL:       webAppURL,
	}
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(ctx context.Context, input services.RegisterUserInput) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, services.ErrUserAlreadyExists
	}

	// Validate password
	if err := s.passwordService.ValidatePassword(ctx, input.Password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// Hash password
	hashedPassword, err := s.passwordService.HashPassword(ctx, input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.NewUser(input.Email, input.Username, models.RoleUser)
	user.PasswordHash = hashedPassword

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send verification email
	if err := s.eventPublisher.PublishUserEvent(ctx, string(events.UserRegistered), events.NewUserRegisteredEvent(
		user.ID,
		user.Email,
		user.Username,
		input.FirstName,
		input.LastName,
	)); err != nil {
		s.logger.Error("failed to publish user registered event", zap.Error(err))
	}

	return user, nil
}

// Login authenticates a user and returns access and refresh tokens
func (s *Service) Login(ctx context.Context, input services.LoginUserInput) (*services.LoginResponse, error) {
	// Find user
	var user *models.User
	var err error

	if input.Email != "" {
		user, err = s.userRepo.GetByEmail(ctx, input.Email)
	} else if input.Username != "" {
		user, err = s.userRepo.GetByUsername(ctx, input.Username)
	} else {
		return nil, services.ErrInvalidCredentials
	}

	if err != nil || user == nil {
		return nil, services.ErrInvalidCredentials
	}

	// Verify password
	if err := s.passwordService.VerifyPassword(ctx, input.Password, user.PasswordHash); err != nil {
		return nil, services.ErrInvalidCredentials
	}

	// Generate tokens
	claims := services.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		TokenType: services.TokenTypeAccess,
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login
	user.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update last login time", zap.Error(err))
	}

	return &services.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// AuthenticateUser authenticates a user with email/username and password
func (s *Service) AuthenticateUser(ctx context.Context, emailOrUsername, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Try to find user by email first
	user, err = s.userRepo.GetByEmail(ctx, emailOrUsername)
	if err != nil {
		// If not found by email, try username
		user, err = s.userRepo.GetByUsername(ctx, emailOrUsername)
		if err != nil {
			return nil, services.ErrInvalidCredentials
		}
	}

	// Verify password
	if err := s.passwordService.VerifyPassword(ctx, password, user.PasswordHash); err != nil {
		return nil, services.ErrInvalidCredentials
	}

	return user, nil
}

// VerifyEmail verifies a user's email address
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	claims, err := s.tokenService.ValidateToken(ctx, token, services.TokenTypeVerification)
	if err != nil {
		return fmt.Errorf("invalid verification token: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.VerifyEmail()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Publish email verified event
	if err := s.eventPublisher.PublishUserEvent(ctx, string(events.UserVerified), events.NewUserVerifiedEvent(
		user.ID,
		user.Email,
	)); err != nil {
		s.logger.Error("failed to publish user verified event", zap.Error(err))
	}

	return nil
}

// RequestPasswordReset initiates the password reset process
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return services.ErrNotFound
	}

	claims := services.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: services.TokenTypeReset,
	}

	token, err := s.tokenService.GenerateResetToken(ctx, claims)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Publish password reset requested event
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.webAppURL, token)
	if err := s.eventPublisher.PublishUserEvent(ctx, string(events.UserPasswordReset), events.NewUserPasswordResetEvent(
		user.ID,
		user.Email,
		resetLink,
	)); err != nil {
		s.logger.Error("failed to publish password reset event", zap.Error(err))
	}

	return nil
}

// ResetPassword resets a user's password using a reset token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	claims, err := s.tokenService.ValidateToken(ctx, token, services.TokenTypeReset)
	if err != nil {
		return fmt.Errorf("invalid reset token: %w", err)
	}

	if err := s.passwordService.ValidatePassword(ctx, newPassword); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	hashedPassword, err := s.passwordService.HashPassword(ctx, newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.UpdatePassword(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Publish password changed event
	if err := s.eventPublisher.PublishUserEvent(ctx, string(events.UserPasswordChange), events.NewUserPasswordChangedEvent(
		user.ID,
		user.Email,
	)); err != nil {
		s.logger.Error("failed to publish password changed event", zap.Error(err))
	}

	// Revoke all existing tokens
	if err := s.tokenService.RevokeToken(ctx, token); err != nil {
		s.logger.Error("failed to revoke reset token", zap.Error(err))
	}

	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*services.TokenResponse, error) {
	claims, err := s.tokenService.ValidateToken(ctx, refreshToken, services.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	isRevoked, err := s.tokenService.IsTokenRevoked(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to check token revocation: %w", err)
	}
	if isRevoked {
		return nil, services.ErrTokenRevoked
	}

	newClaims := services.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		TokenType: services.TokenTypeAccess,
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, newClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(ctx, newClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := s.tokenService.RevokeToken(ctx, refreshToken); err != nil {
		s.logger.Error("failed to revoke old refresh token", zap.Error(err))
	}

	return &services.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout invalidates a user's tokens
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	if err := s.tokenService.RevokeToken(ctx, accessToken); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// GetUser retrieves a user by their ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user's profile
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, input services.UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if input.Email != "" && input.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, services.ErrEmailAlreadyExists
		}
		user.Email = input.Email
		user.Status = models.UserStatusPending // Require email verification again
	}

	if input.Username != "" && input.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, input.Username)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, services.ErrUsernameAlreadyExists
		}
		user.Username = input.Username
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser soft deletes a user account
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Publish user deleted event
	event := events.NewUserDeletedEvent(user.ID, user.Email)
	if err := s.eventPublisher.PublishUserEvent(ctx, "user.deleted", event); err != nil {
		s.logger.Error("failed to publish user deleted event", zap.Error(err))
	}

	return nil
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := s.passwordService.VerifyPassword(ctx, currentPassword, user.PasswordHash); err != nil {
		return services.ErrInvalidCredentials
	}

	// Validate new password
	if err := s.passwordService.ValidatePassword(ctx, newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.passwordService.HashPassword(ctx, newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Publish password changed event
	if err := s.eventPublisher.PublishUserEvent(ctx, string(events.UserPasswordChange), events.NewUserPasswordChangedEvent(
		user.ID,
		user.Email,
	)); err != nil {
		s.logger.Error("failed to publish password changed event", zap.Error(err))
	}

	return nil
}
