package internal

import (
	"net"
	"net/http"
	"strings"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/limiter" // Replace 'yourmodule' with your actual module name
)

// RateLimitMiddleware returns a middleware that applies the given limiter.
func RateLimitMiddleware(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getIPFromRequest(r)
			token := getTokenFromRequest(r)
			if !l.AllowRequest(key, token) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenFromRequest(r *http.Request) string {
	// Check Authorization header for Bearer token
	authHeader := r.Header.Get("API_KEY")
	if authHeader != "" {
		return authHeader
	}
	return ""
}

func getIPFromRequest(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// The first IP in the list is the client IP
		for _, ip := range splitAndTrim(xff, ",") {
			if ip != "" {
				return ip
			}
		}
	}
	// Check X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	// Fallback to RemoteAddr (may include port)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
