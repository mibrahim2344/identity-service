package services

// MetricsService defines the interface for collecting and managing application metrics
type MetricsService interface {
	// RecordRequest records an incoming request with its duration and status
	RecordRequest(path string, method string, statusCode int, duration float64)
	
	// IncrementCounter increments a named counter
	IncrementCounter(name string, labels map[string]string)
	
	// ObserveValue records a value observation for a metric
	ObserveValue(name string, value float64, labels map[string]string)
}
