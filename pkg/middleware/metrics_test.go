package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

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
