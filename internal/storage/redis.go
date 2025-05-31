package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
)

// RedisStore implements a simple Redis-backed store for a rate limiter.
type RedisStore struct {
	Store
	client *redis.Client
	prefix string
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
		prefix: cfg.Prefix,
	}
}

// Get returns the current value and TTL for a given key.
func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	result, err := s.client.Get(ctx, s.prefix+key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return result, nil
}

// Set sets the value and expiration for a given key.
func (s *RedisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if expiration <= 0 {
		return s.client.Set(ctx, s.prefix+key, value, 0).Err()
	}
	return nil
}

// Incr increments the value for a given key and sets expiration if new.
func (s *RedisStore) Incr(ctx context.Context, key string) error {
	return s.client.Incr(ctx, key).Err()
}

func (s *RedisStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := s.client.Expire(ctx, s.prefix+key, expiration).Err()
	return err
}
