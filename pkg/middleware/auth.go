// Middleware package provides checks and methods to be applied
// for all endpoints like auth and metrics
package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

// A simple middleware that rejects all the requests not from localhost
// Additional middlewares that should be added for auth/authz are API Key and/or JWT
// verifications
func LocalHostAuth() func(http.Handler) http.Handler {
	slog.Info("Local Auth middleware initialized")
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// A Simple Auth middleware that only allows requests from localhost
			// r.RemoteAddr includes port (e.g., "127.0.0.1:12345")
			lastIndex := strings.LastIndex(r.RemoteAddr, ":")
			ip, _ := r.RemoteAddr[:lastIndex], r.RemoteAddr[lastIndex+1:]
			if ip == "127.0.0.1" || ip == "::1" || ip == "[::1]" {
				next.ServeHTTP(w, r) // Proceed to the actual handler
				return
			}
			http.Error(w, "Unauthorized request from "+r.RemoteAddr+ip, http.StatusUnauthorized)
		})
	}
}
