package application

import (
	"fmt"
	"time"

	"github.com/mibrahim2344/identity-service/internal/application/user"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/auth/password"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/auth/token"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/events/kafka"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/metrics"
	pgdb "github.com/mibrahim2344/identity-service/internal/infrastructure/persistence/postgres"
	pgrepo "github.com/mibrahim2344/identity-service/internal/infrastructure/persistence/postgres/repositories"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/persistence/redis"
	"go.uber.org/zap"
)

// Config holds all the configuration needed for the application services
type Config struct {
	Database struct {
		Host                   string
		Port                   int
		User                   string
		Password               string
		DBName                 string
		SSLMode                string
		MaxIdleConns           int
		MaxOpenConns           int
		ConnMaxLifetimeMinutes int
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	Kafka struct {
		Brokers []string
		Topic   string
	}
	Auth struct {
		AccessTokenDuration  int // in minutes
		RefreshTokenDuration int // in minutes
		SigningKey           string
		HashingCost          int
	}
	Cache struct {
		DefaultTTL time.Duration
		MaxEntries int
		Prefix     string
		Namespace  string
	}
	WebApp struct {
		URL string
	}
	Server struct {
		Host           string
		Port           int
		ReadTimeout    int // in seconds
		WriteTimeout   int // in seconds
		MaxHeaderBytes int
	}
}

// Factory is responsible for creating and wiring application services
type Factory struct {
	config Config
	logger *zap.Logger
}

// NewFactory creates a new application service factory
func NewFactory(config Config, logger *zap.Logger) *Factory {
	return &Factory{
		config: config,
		logger: logger,
	}
}

// CreateUserService creates and configures the user service with all its dependencies
func (f *Factory) CreateUserService() (services.UserService, error) {
	// Create database connection
	db, err := pgdb.NewConnection(pgdb.Config{
		Host:     f.config.Database.Host,
		Port:     f.config.Database.Port,
		User:     f.config.Database.User,
		Password: f.config.Database.Password,
		DBName:   f.config.Database.DBName,
		SSLMode:  f.config.Database.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create Redis client
	redisClient, err := redis.NewClient(redis.Config{
		Host:     f.config.Redis.Host,
		Port:     f.config.Redis.Port,
		Password: f.config.Redis.Password,
		DB:       f.config.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	// Create repositories
	userRepo := pgrepo.NewUserRepository(db)

	// Create cache service
	defaultCacheConfig := &defaultCacheConfig{
		defaultTTL: f.config.Cache.DefaultTTL,
		maxEntries: f.config.Cache.MaxEntries,
		prefix:     f.config.Cache.Prefix,
		namespace:  f.config.Cache.Namespace,
	}
	cacheService := redis.NewCacheService(redisClient, defaultCacheConfig)

	// Create event publisher
	eventPublisher := kafka.NewPublisher(f.config.Kafka.Brokers)

	// Create password service
	passwordHasher, err := password.NewPasswordHasher(password.BCrypt, map[string]interface{}{
		"cost": f.config.Auth.HashingCost,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create password hasher: %w", err)
	}

	passwordService := password.NewService(passwordHasher, services.PasswordConfig{
		MinLength:           8,
		MaxLength:           72, // bcrypt max length
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireNumbers:      true,
		RequireSpecialChars: true,
	}, nil)

	// Create token service
	keyManager := token.NewRedisKeyManager(cacheService)
	tokenService := token.NewService(services.TokenConfig{
		AccessTokenDuration:  time.Duration(f.config.Auth.AccessTokenDuration) * time.Minute,
		RefreshTokenDuration: time.Duration(f.config.Auth.RefreshTokenDuration) * time.Minute,
	}, cacheService, keyManager)

	// Create user service
	userService := user.NewService(
		userRepo,
		passwordService,
		tokenService,
		cacheService,
		eventPublisher,
		f.logger,
		defaultCacheConfig,
		f.config.WebApp.URL,
	)

	return userService, nil
}

// CreateMetricsService creates and configures the metrics service
func (f *Factory) CreateMetricsService() (services.MetricsService, error) {
	metricsService := metrics.NewMetricsService()
	return metricsService, nil
}

// CreateTokenService creates and configures the token service
func (f *Factory) CreateTokenService() (services.TokenService, error) {
	// Create Redis client for token revocation storage
	redisClient, err := redis.NewClient(redis.Config{
		Host:     f.config.Redis.Host,
		Port:     f.config.Redis.Port,
		Password: f.config.Redis.Password,
		DB:       f.config.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	// Configure token service
	tokenConfig := services.TokenConfig{
		AccessTokenDuration:       time.Duration(f.config.Auth.AccessTokenDuration) * time.Minute,
		RefreshTokenDuration:      time.Duration(f.config.Auth.RefreshTokenDuration) * time.Minute,
		ResetTokenDuration:        24 * time.Hour, // Default 24 hours for reset tokens
		VerificationTokenDuration: 48 * time.Hour, // Default 48 hours for verification tokens
		SigningKey:                []byte(f.config.Auth.SigningKey),
	}

	// Create key manager for JWT signing
	keyManager := token.NewLocalKeyManager()

	// Create Redis cache service wrapper
	cacheService := redis.NewCacheService(redisClient, &defaultCacheConfig{})

	// Create token service with Redis-based revocation storage
	tokenService := token.NewService(tokenConfig, cacheService, keyManager)
	return tokenService, nil
}

// Close closes all connections and resources
func (f *Factory) Close() error {
	// TODO: Implement cleanup of resources
	return nil
}

// defaultCacheConfig implements services.CacheConfig
type defaultCacheConfig struct {
	defaultTTL time.Duration
	maxEntries int
	prefix     string
	namespace  string
}

func (c *defaultCacheConfig) GetDefaultTTL() time.Duration {
	return c.defaultTTL
}

func (c *defaultCacheConfig) GetMaxEntries() int {
	return c.maxEntries
}

func (c *defaultCacheConfig) GetPrefix() string {
	return c.prefix
}

func (c *defaultCacheConfig) GetNamespace() string {
	return c.namespace
}
