package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/handlers"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Router handles all routing logic
type Router struct {
	userService    services.UserService
	tokenService   services.TokenService
	metricsService services.MetricsService
	emailService   services.EmailService
	logger         *zap.Logger
}

// NewRouter creates a new router instance
func NewRouter(
	userService services.UserService,
	tokenService services.TokenService,
	metricsService services.MetricsService,
	emailService services.EmailService,
	logger *zap.Logger,
) *Router {
	return &Router{
		userService:    userService,
		tokenService:   tokenService,
		metricsService: metricsService,
		emailService:   emailService,
		logger:         logger,
	}
}

// Setup sets up all routes and middleware
func (r *Router) Setup() http.Handler {
	router := mux.NewRouter()

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(r.tokenService, r.metricsService, r.logger)
	loggingMiddleware := middleware.NewLoggingMiddleware(r.logger, r.metricsService)

	// Create handlers
	userHandler := handlers.NewUserHandler(r.userService, r.metricsService, r.emailService, r.logger)

	// Apply global middleware
	router.Use(loggingMiddleware.LogRequest)

	// Public routes
	router.HandleFunc("/api/v1/users/register", userHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/verify-email", userHandler.VerifyEmail).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/request-password-reset", userHandler.RequestPasswordReset).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/reset-password", userHandler.ResetPassword).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/refresh-token", userHandler.RefreshToken).Methods(http.MethodPost)

	// Protected routes
	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(authMiddleware.Authenticate)
	protected.HandleFunc("/users/me", userHandler.GetUser).Methods(http.MethodGet)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Not found handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	return router
}
