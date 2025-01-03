package services

import (
	"time"

	"github.com/mibrahim2344/identity-service/internal/domain/repositories"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"gorm.io/gorm"
)

// Services holds all infrastructure service dependencies
type Services struct {
	DB               *gorm.DB
	Cache            services.CacheService
	EventPublisher   services.EventPublisher
	MetricsCollector services.MetricsService
	Password         services.PasswordService
	Token            services.TokenService
	UserRepository   repositories.UserRepository
}

// CacheService is an alias for domain.CacheService to avoid import cycles
type CacheService = services.CacheService

// NewServices creates a new instance of Services with all dependencies
func NewServices(
	db *gorm.DB,
	cache services.CacheService,
	eventPublisher services.EventPublisher,
	metricsCollector services.MetricsService,
	userRepo repositories.UserRepository,
	tokenSecret string,
	accessTokenExpiry,
	refreshTokenExpiry time.Duration,
) *Services {
	return &Services{
		DB:               db,
		Cache:            cache,
		EventPublisher:   eventPublisher,
		MetricsCollector: metricsCollector,
		Password:         NewPasswordService(),
		Token:            NewTokenService(tokenSecret, accessTokenExpiry, refreshTokenExpiry),
		UserRepository:   userRepo,
	}
}
