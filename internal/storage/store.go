package storage

import (
	"context"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value int64, expiration time.Duration) error
	Incr(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
