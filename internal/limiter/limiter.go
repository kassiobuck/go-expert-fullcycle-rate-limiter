package limiter

import (
	"context"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/auth"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/internal/storage"
)

type Limiter struct {
	IPMaxRequestsPerSecond int64
	IPBlockDurationSeconds int64
	store                  storage.Store
	ctx                    context.Context
	auth                   auth.AuthInterface
}

func NewLimiter(ctx context.Context, cfg config.RateLimiterConfig, store storage.Store, auth auth.AuthInterface) *Limiter {
	return &Limiter{
		IPMaxRequestsPerSecond: cfg.IpMaxRequest,
		IPBlockDurationSeconds: cfg.IpBlockDuration,
		store:                  store,
		ctx:                    ctx,
		auth:                   auth,
	}
}

func (l *Limiter) AllowRequest(ip string, token string) bool {
	if token != "" {
		return l.isAllowedByToken(token)
	}
	return l.isAllowedByIP(ip)
}

func (l *Limiter) isAllowedByIP(ip string) bool {
	key := "[ip]" + ip
	return l.isAllowed(key, l.IPMaxRequestsPerSecond, l.IPBlockDurationSeconds, l.store)
}

func (l *Limiter) isAllowedByToken(token string) bool {
	key := "[token]" + token
	c, error := l.auth.ValidateToken(token)
	if error != nil {
		println("Invalid token: ", error.Error())
		return false
	}

	return l.isAllowed(key, c.MaxAccess, c.IntervalAccess, l.store)
}

func (l *Limiter) isAllowed(key string, maxRequestsPerSecond int64, blockDurationSeconds int64, store storage.Store) bool {
	count, err := store.Get(l.ctx, key)
	if err != nil {
		println(err.Error())
		return false
	}

	if count >= maxRequestsPerSecond {
		println("Rate limit exceeded for key: ", key)
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
