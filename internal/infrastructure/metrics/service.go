package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsService struct {
	requestDuration *prometheus.HistogramVec
	counters       map[string]*prometheus.CounterVec
	observations   map[string]*prometheus.GaugeVec
}

// NewMetricsService creates a new metrics service using Prometheus
func NewMetricsService() *metricsService {
	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"path", "method", "status"},
	)

	return &metricsService{
		requestDuration: requestDuration,
		counters:       make(map[string]*prometheus.CounterVec),
		observations:   make(map[string]*prometheus.GaugeVec),
	}
}

// RecordRequest records an incoming request with its duration and status
func (m *metricsService) RecordRequest(path string, method string, statusCode int, duration float64) {
	m.requestDuration.WithLabelValues(
		path,
		method,
		string(rune(statusCode)),
	).Observe(duration)
}

// IncrementCounter increments a named counter
func (m *metricsService) IncrementCounter(name string, labels map[string]string) {
	counter, exists := m.counters[name]
	if !exists {
		counter = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: "Custom counter " + name,
			},
			getLabelsKeys(labels),
		)
		m.counters[name] = counter
	}
	counter.With(labels).Inc()
}

// ObserveValue records a value observation for a metric
func (m *metricsService) ObserveValue(name string, value float64, labels map[string]string) {
	gauge, exists := m.observations[name]
	if !exists {
		gauge = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: "Custom gauge " + name,
			},
			getLabelsKeys(labels),
		)
		m.observations[name] = gauge
	}
	gauge.With(labels).Set(value)
}

func getLabelsKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}
