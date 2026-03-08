// Middleware package provides checks and methods to be applied
// for all endpoints like auth and metricspackage middleware
package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	httpLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpLatency)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Collects simple metrics like request latency and number of metrics
// The behavior can be extended per endpoint
func Metrics(next http.Handler) http.Handler {
	slog.Info("Metrics middleware initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         200,
		}
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()
		httpRequests.WithLabelValues(
			r.URL.Path,
			r.Method,
			strconv.Itoa(rec.status),
		).Inc()
		httpLatency.WithLabelValues(
			r.URL.Path,
			r.Method,
		).Observe(duration)
	})
}
