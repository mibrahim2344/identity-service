package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService    services.UserService
	metricsService services.MetricsService
	emailService   services.EmailService
	logger         *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userService services.UserService,
	metricsService services.MetricsService,
	emailService services.EmailService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		metricsService: metricsService,
		emailService:   emailService,
		logger:         logger,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password       string `json:"password"`
}

// RequestPasswordResetRequest represents the request body for password reset request
type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents the request body for password reset
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.RegisterUser(r.Context(), services.RegisterUserInput{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		h.handleError(w, r, err, http.StatusInternalServerError, "failed to register user")
		return
	}

	h.respondJSON(w, http.StatusCreated, user)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.userService.AuthenticateUser(r.Context(), req.EmailOrUsername, req.Password)

	if err != nil {
		h.handleError(w, r, err, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	token := r.URL.Query().Get("token")
	if token == "" {
		h.handleError(w, r, nil, http.StatusBadRequest, "missing verification token")
		return
	}

	if err := h.userService.VerifyEmail(r.Context(), token); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid verification token")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "email verified successfully"})
}

// GetUser handles retrieving user information
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	userID := r.Context().Value("user_id").(string)
	id, err := uuid.Parse(userID)
	if err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err, http.StatusNotFound, "user not found")
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

// RequestPasswordReset handles password reset request
func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	var req RequestPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.userService.RequestPasswordReset(r.Context(), req.Email); err != nil {
		h.handleError(w, r, err, http.StatusInternalServerError, "failed to request password reset")
		return
	}

	// Send success response even if user doesn't exist (security best practice)
	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "if the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.userService.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "failed to reset password")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "password has been reset successfully",
	})
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.metricsService.RecordRequest(r.Method, r.URL.Path, http.StatusOK, time.Since(start).Seconds())
	}()

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.userService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(w, r, err, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	h.respondJSON(w, http.StatusOK, tokens)
}

func (h *UserHandler) handleError(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	h.logger.Error(message,
		zap.Error(err),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)

	h.metricsService.IncrementCounter("http_errors", map[string]string{
		"path":    r.URL.Path,
		"method":  r.Method,
		"message": message,
	})
	h.respondJSON(w, status, map[string]string{"error": message})
}

func (h *UserHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.logger.Error("failed to encode response",
				zap.Error(err),
			)
		}
	}
}
