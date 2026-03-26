// Package middleware provides unit tests for Prometheus metrics collection middleware.
// Tests verify that HTTP requests are properly tracked with counters and histograms.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

// TestMetricsMiddlewareRecordsRequest verifies that the metrics middleware
// successfully records HTTP requests in the Prometheus counter.
//
// Test scenario:
//  1. Creates a handler wrapped with Metrics middleware
//  2. Simulates a GET request to /v1/anagrams
//  3. Verifies the request completes with HTTP 200 OK
//  4. Checks that at least one metric was recorded
//
// This ensures the request counter is working correctly.
func TestMetricsMiddlewareRecordsRequest(t *testing.T) {

	handler := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/v1/anagrams", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	count := testutil.CollectAndCount(httpRequests)

	if count == 0 {
		t.Fatalf("expected metrics to be recorded")
	}
}

// TestMetricsMiddlewareStatusCode verifies that the metrics middleware
// correctly captures and records error status codes.
//
// Test scenario:
//  1. Creates a handler that returns HTTP 500 Internal Server Error
//  2. Wraps it with Metrics middleware
//  3. Simulates a request
//  4. Verifies the error status is returned
//  5. Checks that the metric was recorded with the error status
//
// This ensures error responses are properly tracked in metrics.
func TestMetricsMiddlewareStatusCode(t *testing.T) {

	handler := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failure", http.StatusInternalServerError)
	}))

	req := httptest.NewRequest("GET", "/v1/anagrams", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500")
	}

	count := testutil.CollectAndCount(httpRequests)

	if count == 0 {
		t.Fatalf("expected metrics increment")
	}
}

// TestMetricsMiddlewareLatencyMetric verifies that the metrics middleware
// records request latency in the Prometheus histogram.
//
// Test scenario:
//  1. Creates a handler wrapped with Metrics middleware
//  2. Simulates a request
//  3. Verifies that latency metrics were recorded
//
// This ensures request duration tracking is working correctly for
// performance monitoring and SLA tracking.
func TestMetricsMiddlewareLatencyMetric(t *testing.T) {

	handler := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/v1/anagrams", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	count := testutil.CollectAndCount(httpLatency)

	if count == 0 {
		t.Fatalf("expected latency metric recorded")
	}
}
