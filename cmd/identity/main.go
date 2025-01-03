package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mibrahim2344/identity-service/docs" // Import swagger docs
	"github.com/mibrahim2344/identity-service/internal/application/config"
	"github.com/mibrahim2344/identity-service/internal/application/user"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/events/kafka"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/metrics"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/persistence/postgres"
	"github.com/mibrahim2344/identity-service/internal/infrastructure/persistence/redis"
	infraservices "github.com/mibrahim2344/identity-service/internal/infrastructure/services"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/server"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Force unbuffered output
	os.Stdout.Sync()

	fmt.Println("Starting identity service...")

	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	fmt.Println("Initializing logger...")
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	fmt.Println("Logger initialized successfully")

	// Load configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.LoadConfig("config/default.json")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	fmt.Println("Configuration loaded successfully")

	// Initialize database connection
	fmt.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	db, err := gorm.Open(pgdriver.New(pgdriver.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	fmt.Println("Database connection established successfully")

	// Get underlying SQL DB
	fmt.Println("Getting underlying SQL DB...")
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB", zap.Error(err))
	}
	fmt.Println("Underlying SQL DB retrieved successfully")

	// Configure connection pool
	fmt.Println("Configuring connection pool...")
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)
	fmt.Println("Connection pool configured successfully")

	// Initialize Redis client
	fmt.Println("Initializing Redis client...")
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
	fmt.Println("Redis client initialized successfully")

	// Initialize cache service with config
	fmt.Println("Initializing cache service...")
	cacheConfig := redis.NewCacheConfig(
		cfg.Cache.DefaultTTL,
		cfg.Cache.MaxEntries,
		cfg.Cache.Prefix,
		cfg.Cache.Namespace,
	)
	cacheService := redis.NewCacheService(redisClient, cacheConfig)
	fmt.Println("Cache service initialized successfully")

	// Initialize Kafka producer
	fmt.Println("Initializing Kafka producer...")
	kafkaProducer := kafka.NewPublisher(cfg.Kafka.Brokers)
	defer kafkaProducer.Close()
	fmt.Println("Kafka producer initialized successfully")

	// Initialize metrics collector
	fmt.Println("Initializing metrics collector...")
	metricsCollector := metrics.NewMetricsService()
	fmt.Println("Metrics collector initialized successfully")

	// Initialize user repository
	fmt.Println("Initializing user repository...")
	userRepo := postgres.NewRepository(db)
	fmt.Println("User repository initialized successfully")

	// Initialize infrastructure services
	fmt.Println("Initializing infrastructure services...")
	services := infraservices.NewServices(
		db,                  // *gorm.DB
		cacheService,        // services.CacheService
		kafkaProducer,       // services.EventPublisher
		metricsCollector,    // MetricsCollector
		userRepo,            // repositories.UserRepository
		cfg.Auth.SigningKey, // tokenSecret string
		time.Duration(cfg.Auth.AccessTokenDuration)*time.Second,  // accessTokenExpiry time.Duration
		time.Duration(cfg.Auth.RefreshTokenDuration)*time.Second, // refreshTokenExpiry time.Duration
	)
	fmt.Println("Infrastructure services initialized successfully")

	// Initialize user application service
	fmt.Println("Initializing user application service...")
	userApp := user.NewService(
		services.UserRepository,
		services.Password,
		services.Token,
		services.Cache,
		services.EventPublisher,
		logger,
		redis.NewCacheConfig(
			cfg.Cache.DefaultTTL,
			cfg.Cache.MaxEntries,
			cfg.Cache.Prefix,
			cfg.Cache.Namespace,
		),
		cfg.WebApp.URL,
	)
	fmt.Println("User application service initialized successfully")

	// Initialize HTTP server
	fmt.Println("Initializing HTTP server...")
	httpServer := server.NewServer(
		server.Config{
			Host:           cfg.Server.Host,
			Port:           cfg.Server.Port,
			ReadTimeout:    10 * time.Second, // default timeout
			WriteTimeout:   10 * time.Second, // default timeout
			MaxHeaderBytes: 1 << 20,          // default 1MB
			AllowedOrigins: []string{"*"},    // allow all origins
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		userApp,
		services.Token,
		services.MetricsCollector,
		logger,
	)
	fmt.Println("HTTP server initialized successfully")

	// Start HTTP server
	fmt.Println("Starting HTTP server...")
	errChan := make(chan error, 1)
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", zap.Error(err))
			errChan <- err
		}
	}()

	// Wait for interrupt signal or error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Server error", zap.Error(err))
		os.Exit(1)
	case sig := <-sigChan:
		logger.Info("Received signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	fmt.Println("Server is running. Press Ctrl+C to stop.")
	<-sigChan
}
