// Package middleware provides HTTP middleware functions for authentication,
// authorization, and metrics collection.
//
// Middleware functions wrap HTTP handlers to add cross-cutting concerns like:
// - Authentication and authorization checks
// - Request/response logging
// - Metrics collection
// - Rate limiting (future enhancement)
//
// All middleware in this package follows the standard Go HTTP middleware pattern.
package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

// LocalHostAuth creates a middleware that restricts access to localhost only.
// This is a basic security measure to prevent external access to the API.
//
// Security considerations:
//   - Only allows requests from 127.0.0.1, ::1, or [::1] (IPv4 and IPv6 localhost)
//   - Rejects all other requests with HTTP 401 Unauthorized
//   - In production, additional authentication should be added (API keys, JWT, OAuth)
//
// Returns:
//   - func(http.Handler) http.Handler: Middleware function that wraps the next handler
//
// Example usage:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/data", dataHandler)
//	authMiddleware := middleware.LocalHostAuth()
//	http.ListenAndServe(":8080", authMiddleware(mux))
//
// Future enhancements:
//   - API Key validation
//   - JWT token verification
//   - OAuth 2.0 integration
//   - Role-based access control (RBAC)
func LocalHostAuth() func(http.Handler) http.Handler {
	slog.Info("Local Auth middleware initialized")
	slog.Info("LocalHost authentication middleware initialized")
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP address from RemoteAddr
			// RemoteAddr format: "IP:port" (e.g., "127.0.0.1:12345")
			lastIndex := strings.LastIndex(r.RemoteAddr, ":")
			ip, _ := r.RemoteAddr[:lastIndex], r.RemoteAddr[lastIndex+1:]

			// Check if request is from localhost (IPv4 or IPv6)
			if ip == "127.0.0.1" || ip == "::1" || ip == "[::1]" {
				next.ServeHTTP(w, r) // Allow request to proceed
				return
			}

			// Reject non-localhost requests
			http.Error(w, "Unauthorized request from "+r.RemoteAddr+ip, http.StatusUnauthorized)
		})
	}
}
