package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/storage"
)

type Limiter struct {
	IPMaxRequestsPerSecond    int
	IPBlockDurationSeconds    int
	TokenMaxRequestsPerSecond int
	TokenBlockDurationSeconds int
	store                     storage.Store
	ctx                       context.Context
}

func NewLimiter(ctx context.Context, cfg config.RateLimiterConfig, store storage.Store) *Limiter {
	return &Limiter{
		IPMaxRequestsPerSecond:    cfg.IpMaxRequest,
		IPBlockDurationSeconds:    cfg.IpBlockDuration,
		TokenMaxRequestsPerSecond: cfg.TokenMaxRequest,
		TokenBlockDurationSeconds: cfg.TokenBlockDuration,
		store:                     store,
		ctx:                       ctx,
	}
}

func (l *Limiter) AllowRequest(ip string, token string) bool {
	if token != "" {
		return l.isAllowedByToken(token)
	}
	return l.isAllowedByIP(ip)
}

func (l *Limiter) isAllowedByIP(ip string) bool {
	key := "[ip] " + ip
	return l.isAllowed(key, l.IPMaxRequestsPerSecond, l.IPBlockDurationSeconds, l.store)
}

func (l *Limiter) isAllowedByToken(token string) bool {
	key := "[token] " + token
	return l.isAllowed(key, l.TokenMaxRequestsPerSecond, l.TokenBlockDurationSeconds, l.store)
}

func (l *Limiter) isAllowed(key string, maxRequestsPerSecond, blockDurationSeconds int, store storage.Store) bool {
	countStr, err := store.Get(l.ctx, key)
	if err != nil {
		println(err.Error())
		return false
	}

	count, _ := strconv.Atoi(countStr)
	println(count, maxRequestsPerSecond)
	if count >= maxRequestsPerSecond {
		println("Rate limit exceeded for key: %s", key)
		return false
	}

	expiration := time.Duration(blockDurationSeconds) * time.Second
	if err := store.Incr(l.ctx, key); err != nil {
		return false
	}
	if err := store.Expire(l.ctx, key, expiration); err != nil {
		return false
	}

	return true
}
