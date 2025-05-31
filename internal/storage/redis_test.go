package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/stretchr/testify/assert"
)

var (
	cfg         = config.LoadConfig()
	testRedisDB = 0
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newTestRedisStore(t *testing.T) *RedisStore {
	cfg := config.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       testRedisDB,
	}
	store := NewRedisStore(cfg)
	// Ping to ensure connection
	err := store.client.Ping(context.Background()).Err()
	if err != nil {
		t.Fatalf("Could not connect to Redis: %v", err)
	}
	return store
}

func TestRedisStore_SetGet(t *testing.T) {
	store := newTestRedisStore(t)
	ctx := context.Background()
	key := "test_key_setget"
	value := "12345"

	// Clean up before and after
	store.client.Del(ctx, key)

	err := store.Set(ctx, key, value, 2*time.Second)
	assert.NoError(t, err)

	got, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)

	// Wait for expiration
	time.Sleep(3 * time.Second)
	got, err = store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestRedisStore_SetNoExpiration(t *testing.T) {
	store := newTestRedisStore(t)
	ctx := context.Background()
	key := "test_key_noexp"
	value := "abc"

	store.client.Del(ctx, key)

	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	got, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)

	store.client.Del(ctx, key)
}

func TestRedisStore_Incr(t *testing.T) {
	store := newTestRedisStore(t)
	ctx := context.Background()
	key := "test_key_incr"

	store.client.Del(ctx, key)

	err := store.Incr(ctx, key)
	assert.NoError(t, err)

	got, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "1", got)

	err = store.Incr(ctx, key)
	assert.NoError(t, err)
	got, err = store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "2", got)

	store.client.Del(ctx, key)
}

func TestRedisStore_Expire(t *testing.T) {
	store := newTestRedisStore(t)
	ctx := context.Background()
	key := "test_key_expire"
	value := "expireme"

	store.client.Del(ctx, key)
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	err = store.Expire(ctx, key, 1*time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	got, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestRedisStore_GetNonExistent(t *testing.T) {
	store := newTestRedisStore(t)
	ctx := context.Background()
	key := "nonexistent_key"

	store.client.Del(ctx, key)

	got, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "", got)
}
