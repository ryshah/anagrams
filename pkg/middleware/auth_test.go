// Package middleware provides unit tests for authentication middleware.
// Tests verify that localhost requests are allowed and external requests are blocked.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthMiddleware verifies that requests from localhost are allowed through
// the authentication middleware.
//
// Test scenario:
//  1. Creates a handler wrapped with LocalHostAuth middleware
//  2. Simulates a request from 127.0.0.1 (localhost)
//  3. Verifies the request receives HTTP 200 OK status
//
// This ensures legitimate localhost requests can access the API.
func TestAuthMiddleware(t *testing.T) {

	handler := LocalHostAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080" // Simulate localhost request
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal("expected authorized request")
	}
}

// TestAuthMiddlewareError verifies that requests from non-localhost addresses
// are rejected by the authentication middleware.
//
// Test scenario:
//  1. Creates a handler wrapped with LocalHostAuth middleware
//  2. Simulates a request without setting RemoteAddr (defaults to empty/external)
//  3. Verifies the request receives HTTP 401 Unauthorized status
//
// This ensures external requests are properly blocked for security.
func TestAuthMiddlewareError(t *testing.T) {

	handler := LocalHostAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	// RemoteAddr not set, simulating external request
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatal("expected unauthorized request")
	}
}
