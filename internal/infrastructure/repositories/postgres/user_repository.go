package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	domainErrors "github.com/mibrahim2344/identity-service/internal/domain/errors"
	"github.com/mibrahim2344/identity-service/internal/domain/models"
)

// UserRepository implements the domain.UserRepository interface using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domainErrors.ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainErrors.ErrUserNotFound
	}
	return nil
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainErrors.ErrUserNotFound
	}
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users)
	
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
