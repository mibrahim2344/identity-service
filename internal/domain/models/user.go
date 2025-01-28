package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User represents the user entity in our domain
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email          string         `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Username       string         `gorm:"type:varchar(255);uniqueIndex" json:"username"`
	PasswordHash   string         `gorm:"type:varchar(255)" json:"-"`
	Status         UserStatus     `gorm:"type:user_status;default:'pending'" json:"status"`
	FirstName      string         `gorm:"type:varchar(255)" json:"first_name"`
	LastName       string         `gorm:"type:varchar(255)" json:"last_name"`
	Role           Role          `gorm:"type:user_role;default:'user'" json:"role"`
	EmailVerified  bool          `gorm:"default:false" json:"email_verified"`
	CreatedAt      time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"not null" json:"updated_at"`
	LastLoginAt    *time.Time    `json:"last_login_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate will update the UpdatedAt timestamp
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with default values
func NewUser(email, username string, role Role) *User {
	return &User{
		Email:         email,
		Username:      username,
		Status:        UserStatusPending,
		Role:          role,
		EmailVerified: false,
	}
}

// UpdatePassword updates the user's password hash
func (u *User) UpdatePassword(passwordHash string) {
	u.PasswordHash = passwordHash
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.Status = UserStatusActive
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}
