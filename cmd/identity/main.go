package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mibrahim2344/identity-service/internal/application"
	"github.com/mibrahim2344/identity-service/internal/application/config"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/server"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig("config/default.json")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Create factory
	factory := application.NewFactory(cfg, logger)

	// Create services
	userService, err := factory.CreateUserService()
	if err != nil {
		logger.Fatal("failed to create user service", zap.Error(err))
	}

	emailService, err := factory.CreateEmailService()
	if err != nil {
		logger.Fatal("failed to create email service", zap.Error(err))
	}

	metricsService, err := factory.CreateMetricsService()
	if err != nil {
		logger.Fatal("failed to create metrics service", zap.Error(err))
	}

	tokenService, err := factory.CreateTokenService()
	if err != nil {
		logger.Fatal("failed to create token service", zap.Error(err))
	}

	// Create and start HTTP server
	srv := server.NewServer(
		server.Config{
			Host:           cfg.Server.Host,
			Port:           cfg.Server.Port,
			ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
			MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		userService,
		tokenService,
		metricsService,
		emailService,
		logger,
	)

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Fatal("failed to stop server", zap.Error(err))
	}

	if err := factory.Close(); err != nil {
		logger.Fatal("failed to close factory", zap.Error(err))
	}

	logger.Info("server stopped")
}
