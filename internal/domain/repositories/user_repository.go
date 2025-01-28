package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mibrahim2344/identity-service/internal/domain/models"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByIdentifier retrieves a user by email or username
	GetByIdentifier(ctx context.Context, identifier string) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user by their ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}
