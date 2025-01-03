package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/handlers"
	"github.com/mibrahim2344/identity-service/internal/interfaces/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Router handles all routing logic
type Router struct {
	userService    services.UserService
	tokenService   services.TokenService
	metricsService services.MetricsService
	logger         *zap.Logger
}

// NewRouter creates a new router instance
func NewRouter(
	userService services.UserService,
	tokenService services.TokenService,
	metricsService services.MetricsService,
	logger *zap.Logger,
) *Router {
	return &Router{
		userService:    userService,
		tokenService:   tokenService,
		metricsService: metricsService,
		logger:         logger,
	}
}

// Setup sets up all routes and middleware
func (r *Router) Setup() http.Handler {
	r.logger.Info("Setting up router...")
	router := mux.NewRouter()

	// Apply CORS middleware
	r.logger.Debug("Applying CORS middleware...")
	router.Use(middleware.CORSMiddleware([]string{"*"}))

	// Health check
	r.logger.Debug("Setting up health check endpoint...")
	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			r.logger.Error("failed to write response", zap.Error(err))
		}
	}).Methods(http.MethodGet)

	// API v1 routes
	r.logger.Debug("Setting up API v1 routes...")
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	r.logger.Debug("Setting up auth routes...")
	auth := v1.PathPrefix("/auth").Subrouter()
	userHandler := handlers.NewUserHandler(r.userService, r.metricsService, r.logger)
	auth.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	auth.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", userHandler.RefreshToken).Methods(http.MethodPost)
	auth.HandleFunc("/forgot-password", userHandler.RequestPasswordReset).Methods(http.MethodPost)
	auth.HandleFunc("/reset-password", userHandler.ResetPassword).Methods(http.MethodPost)
	auth.HandleFunc("/verify-email", userHandler.VerifyEmail).Methods(http.MethodGet)

	// Protected routes
	r.logger.Debug("Setting up protected routes...")
	protected := v1.PathPrefix("/").Subrouter()
	authMiddleware := middleware.NewAuthMiddleware(r.tokenService, r.metricsService, r.logger)
	protected.Use(authMiddleware.Authenticate)

	// User routes
	r.logger.Debug("Setting up user routes...")
	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("/me", userHandler.GetUser).Methods(http.MethodGet)
	users.HandleFunc("/me/password", userHandler.ChangePassword).Methods(http.MethodPut)

	// Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Not found handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("not found"))
		if err != nil {
			r.logger.Error("failed to write response", zap.Error(err))
		}
	})

	r.logger.Info("Router setup completed successfully")
	return router
}
