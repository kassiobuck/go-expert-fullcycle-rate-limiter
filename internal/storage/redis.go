package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
)

// RedisStore implements a simple Redis-backed store for a rate limiter.
type RedisStore struct {
	Store
	client *redis.Client
}

// NewRedisStore creates a new RedisStore.
func NewRedisStore(cfg config.RedisConfig) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisStore{
		client: rdb,
	}
}

// Get returns the current value and TTL for a given key.
func (s *RedisStore) Get(ctx context.Context, key string) (int64, error) {
	result, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	val, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// Set sets the value and expiration for a given key.
func (s *RedisStore) Set(ctx context.Context, key string, value int64, expiration time.Duration) error {
	if expiration <= 0 {
		return s.client.Set(ctx, key, value, 0).Err()
	}
	return nil
}

// Incr increments the value for a given key and sets expiration if new.
func (s *RedisStore) Incr(ctx context.Context, key string) error {
	return s.client.Incr(ctx, key).Err()
}

func (s *RedisStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return s.client.Expire(ctx, key, expiration).Err()
}
