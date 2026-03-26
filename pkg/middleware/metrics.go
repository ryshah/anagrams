// Package middleware provides HTTP middleware functions for authentication,
// authorization, and metrics collection using Prometheus.
package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// httpRequests is a Prometheus counter that tracks the total number of HTTP requests.
	// It includes labels for path, method, and status code to enable detailed analysis.
	//
	// Labels:
	//   - path: The request URL path (e.g., "/v1/anagrams")
	//   - method: The HTTP method (e.g., "GET", "POST")
	//   - status: The HTTP status code (e.g., "200", "404", "500")
	//
	// Example query: http_requests_total{path="/v1/anagrams",method="GET",status="200"}
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	// httpLatency is a Prometheus histogram that tracks request duration in seconds.
	// It uses default buckets to measure latency distribution.
	//
	// Labels:
	//   - path: The request URL path
	//   - method: The HTTP method
	//
	// Buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
	// Example query: histogram_quantile(0.95, http_request_duration_seconds)
	httpLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

// init registers the Prometheus metrics with the default registry.
// This runs automatically when the package is imported.
func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpLatency)
}

// statusRecorder wraps http.ResponseWriter to capture the HTTP status code.
// This is necessary because the standard ResponseWriter doesn't expose the status
// after it's been written.
type statusRecorder struct {
	http.ResponseWriter     // Embedded ResponseWriter for delegation
	status              int // Captured status code
}

// WriteHeader captures the status code before delegating to the underlying ResponseWriter.
// This allows the middleware to record the status code in metrics.
//
// Parameters:
//   - code: The HTTP status code to write
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Metrics creates a middleware that collects Prometheus metrics for HTTP requests.
// It tracks request count, status codes, and latency for all requests.
//
// Collected metrics:
//   - http_requests_total: Counter of total requests by path, method, and status
//   - http_request_duration_seconds: Histogram of request latency by path and method
//
// Returns:
//   - http.Handler: Middleware-wrapped handler that collects metrics
//
// Example usage:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/data", dataHandler)
//	http.ListenAndServe(":8080", middleware.Metrics(mux))
//
// Metrics can be scraped by Prometheus at the /metrics endpoint.
func Metrics(next http.Handler) http.Handler {
	slog.Info("Metrics middleware initialized")
	slog.Info("Metrics collection middleware initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record start time for latency calculation
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         200, // Default status if WriteHeader is not called
		}

		// Process the request
		next.ServeHTTP(rec, r)

		// Calculate request duration in seconds
		duration := time.Since(start).Seconds()

		// Increment request counter with labels
		httpRequests.WithLabelValues(
			r.URL.Path,
			r.Method,
			strconv.Itoa(rec.status),
		).Inc()

		// Record request latency
		httpLatency.WithLabelValues(
			r.URL.Path,
			r.Method,
		).Observe(duration)
	})
}
