package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/router"
	"go.uber.org/zap"
)

// Config represents server configuration
type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	AllowedOrigins  []string
	AllowedMethods  []string
	AllowedHeaders  []string
}

// Server represents the HTTP server
type Server struct {
	config         Config
	userService    services.UserService
	tokenService   services.TokenService
	metricsService services.MetricsService
	emailService   services.EmailService
	logger         *zap.Logger
	httpServer     *http.Server
}

// NewServer creates a new server instance
func NewServer(
	config Config,
	userService services.UserService,
	tokenService services.TokenService,
	metricsService services.MetricsService,
	emailService services.EmailService,
	logger *zap.Logger,
) *Server {
	return &Server{
		config:         config,
		userService:    userService,
		tokenService:   tokenService,
		metricsService: metricsService,
		emailService:   emailService,
		logger:         logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	r := router.NewRouter(s.userService, s.tokenService, s.metricsService, s.emailService, s.logger)

	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:        r.Setup(),
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	s.logger.Info("starting HTTP server",
		zap.String("address", s.httpServer.Addr),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping HTTP server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
