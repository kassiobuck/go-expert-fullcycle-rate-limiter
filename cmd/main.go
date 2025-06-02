package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/auth"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/limiter"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/storage"
)

func main() {

	cfg := config.LoadConfig()
	ctx := context.Background()
	store := storage.NewRedisStore(cfg.Redis)
	auth := auth.NewAuth([]byte(cfg.Auth.JwtSecret))

	rateLimiter := limiter.NewLimiter(ctx, cfg.RateLimiter, store, auth)

	mux := http.NewServeMux()

	middleware := internal.RateLimitMiddleware(rateLimiter)

	mux.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	mux.Handle("/genToken", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		intervalAccess, err := strconv.ParseInt(r.URL.Query().Get("interval"), 10, 64)
		if err != nil || intervalAccess <= 0 {
			http.Error(w, "Send interval access valid number", http.StatusBadRequest)
			return
		}
		maxAccess, err := strconv.ParseInt(r.URL.Query().Get("max"), 10, 64)
		if err != nil || maxAccess <= 0 {
			http.Error(w, "Send max access valid number", http.StatusBadRequest)
			return
		}

		timeDuration := time.Duration(9999) * time.Minute
		token, err := auth.GenerateToken("_generic", maxAccess, intervalAccess, timeDuration)
		if err != nil {
			http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(token))
	}))

	mux.Handle("/decodeToken", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")

		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		w.Write([]byte(fmt.Sprintf("Token is valid. Max Access: %d, interval: %d", claims.MaxAccess, claims.IntervalAccess)))
	}))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
