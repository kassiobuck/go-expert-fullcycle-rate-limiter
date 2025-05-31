package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/limiter"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/storage"
)

func main() {

	cfg := config.LoadConfig()
	ctx := context.Background()
	tokenStore := storage.NewRedisStore(cfg.Redis)
	ipStore := storage.NewRedisStore(cfg.Redis)

	rateLimiter := limiter.NewLimiter(ctx, cfg.RateLimiter, ipStore, tokenStore)

	mux := http.NewServeMux()

	middleware := internal.RateLimitMiddleware(rateLimiter)

	mux.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
