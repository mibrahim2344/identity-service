package middleware

import (
	"net/http"
	"time"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
	"go.uber.org/zap"
)

// LoggingMiddleware handles request logging
type LoggingMiddleware struct {
	logger         *zap.Logger
	metricsService services.MetricsService
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *zap.Logger, metricsService services.MetricsService) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger:         logger,
		metricsService: metricsService,
	}
}

// LogRequest logs information about incoming requests
func (m *LoggingMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture the status code
		rw := &responseWriter{w, http.StatusOK}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request details
		duration := time.Since(start)
		m.logger.Info("request processed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Record metrics
		m.metricsService.RecordRequest(r.URL.Path, r.Method, rw.status, duration.Seconds())
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
