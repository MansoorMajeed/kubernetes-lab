package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPMetrics holds all HTTP-related Prometheus metrics
type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

// NewHTTPMetrics creates and registers HTTP metrics
func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "catalog_http_requests_total",
				Help: "Total number of HTTP requests processed by the catalog service",
			},
			[]string{"method", "path", "status_code"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "catalog_http_request_duration_seconds",
				Help:    "Duration of HTTP requests processed by the catalog service",
				Buckets: prometheus.DefBuckets, // Default buckets: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
			},
			[]string{"method", "path"},
		),
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "catalog_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed by the catalog service",
			},
		),
	}
}

// RecordRequest records metrics for an HTTP request
func (m *HTTPMetrics) RecordRequest(method, path, statusCode string, duration float64) {
	// Record request count
	m.RequestsTotal.WithLabelValues(method, path, statusCode).Inc()

	// Record request duration
	m.RequestDuration.WithLabelValues(method, path).Observe(duration)
}

// IncInFlight increments the in-flight requests counter
func (m *HTTPMetrics) IncInFlight() {
	m.RequestsInFlight.Inc()
}

// DecInFlight decrements the in-flight requests counter
func (m *HTTPMetrics) DecInFlight() {
	m.RequestsInFlight.Dec()
}
