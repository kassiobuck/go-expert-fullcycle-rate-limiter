package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/auth"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/limiter"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/storage"
)

func TestLimiter_AllowRequest_IP(t *testing.T) {
	ctx := context.Background()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:    3,
		IpBlockDuration: 2,
	}
	auth := auth.NewAuth([]byte("test-secret"))
	store := storage.NewMockStore()
	lim := limiter.NewLimiter(ctx, cfg, store, auth)

	ip := "127.0.0.1"
	token := ""

	lim.AllowRequest(ip, token)

	// First request should be allowed
	if !lim.AllowRequest(ip, token) {
		t.Error("expected first request to be allowed")
	}
	// Second request should be allowed
	if !lim.AllowRequest(ip, token) {
		t.Error("expected second request to be allowed")
	}
	// Third request should be blocked
	if lim.AllowRequest(ip, token) {
		t.Error("expected third request to be blocked")
	}
}

func TestLimiter_AllowRequest_Token(t *testing.T) {
	ctx := context.Background()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:    2,
		IpBlockDuration: 1,
	}
	auth := auth.NewAuth([]byte("test-secret"))
	store := storage.NewMockStore()
	lim := limiter.NewLimiter(ctx, cfg, store, auth)

	ip := ""
	token, error := auth.GenerateToken("_generic", 1, 2, 2*time.Minute)
	if error != nil {
		t.Fatalf("error generating token: %v", error)
	}
	// First request should be allowed
	if !lim.AllowRequest(ip, token) {
		t.Error("expected first token request to be allowed")
	}
	// Second request should be blocked
	if lim.AllowRequest(ip, token) {
		t.Error("expected second token request to be blocked")
	}
}

func TestServer_WithLimiter(t *testing.T) {
	ctx := context.Background()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:    1,
		IpBlockDuration: 1,
	}
	store := storage.NewMockStore()
	auth := auth.NewAuth([]byte("test-secret"))
	lim := limiter.NewLimiter(ctx, cfg, store, auth)

	middleware := internal.RateLimitMiddleware(lim)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("X-Token")
		if !lim.AllowRequest(ip, token) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}

	server := httptest.NewServer(middleware(http.HandlerFunc(handler)))
	defer server.Close()

	client := &http.Client{}

	// First request should succeed
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.RemoteAddr = "1.2.3.4:5678"
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %v", resp.Status)
	}

	// Second request should be rate limited
	req2, _ := http.NewRequest("GET", server.URL, nil)
	req2.RemoteAddr = "1.2.3.4:5678"
	resp2, err := client.Do(req2)
	if err != nil || resp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %v", resp2.Status)
	}

	token, error := auth.GenerateToken("_generic", 1, 5, 9999*time.Minute)
	if error != nil {
		t.Fatalf("error generating token: %v", error)
	}
	// First request should succeed
	req3, _ := http.NewRequest("GET", server.URL, nil)
	req3.Header.Add("API_KEY", token)
	resp3, err := client.Do(req3)
	if err != nil || resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %v", resp.Status)
	}

	// Second request should be rate limited
	req4, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Add("API_KEY", token)
	resp4, err := client.Do(req4)
	if err != nil || resp4.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %v", resp2.Status)
	}
}
