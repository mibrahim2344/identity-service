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
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Server represents the HTTP server
type Server struct {
	config         Config
	userService    services.UserService
	tokenService   services.TokenService
	metricsService services.MetricsService
	logger         *zap.Logger
	httpServer     *http.Server
	router         *router.Router
}

// NewServer creates a new server instance
func NewServer(
	config Config,
	userService services.UserService,
	tokenService services.TokenService,
	metricsService services.MetricsService,
	logger *zap.Logger,
) *Server {
	return &Server{
		config:         config,
		userService:    userService,
		tokenService:   tokenService,
		metricsService: metricsService,
		logger:         logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Setting up routes...")
	s.router = router.NewRouter(s.userService, s.tokenService, s.metricsService, s.logger)
	handler := s.router.Setup()
	
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.logger.Info("Starting HTTP server", 
		zap.String("address", addr),
		zap.Int("port", s.config.Port),
	)
	
	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.WriteTimeout) * time.Second,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	s.logger.Info("Server is listening", zap.String("address", addr))
	return s.httpServer.ListenAndServe()
}

// setupRoutes configures all the routes for our server
func (s *Server) setupRoutes() http.Handler {
	s.router = router.NewRouter(
		s.userService,
		s.tokenService,
		s.metricsService,
		s.logger,
	)
	return s.router.Setup()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping HTTP server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
