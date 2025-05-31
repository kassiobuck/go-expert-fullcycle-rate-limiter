package limiter

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/kassiobuck/go-expert-fullcycle-rate-limiter/config"
	"github.com/stretchr/testify/assert"
)

// mockStore implements storage.Store for testing
type mockStore struct {
	data      map[string]string
	expiries  map[string]time.Time
	incrCalls map[string]int
}

func newMockStore() *mockStore {
	return &mockStore{
		data:      make(map[string]string),
		expiries:  make(map[string]time.Time),
		incrCalls: make(map[string]int),
	}
}

func (m *mockStore) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	m.data[key] = value
	if expiration > 0 {
		m.expiries[key] = time.Now().Add(expiration)
	}
	return nil
}

func (m *mockStore) Get(ctx context.Context, key string) (string, error) {
	// Simulate expiration
	if exp, ok := m.expiries[key]; ok && time.Now().After(exp) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", nil
	}
	val, ok := m.data[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (m *mockStore) Incr(ctx context.Context, key string) error {
	val, _ := m.Get(ctx, key)
	count := 0
	if val != "" {
		count, _ = strconv.Atoi(val)
	}
	count++
	m.data[key] = strconv.Itoa(count)
	m.incrCalls[key]++
	return nil
}

func (m *mockStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if _, ok := m.data[key]; ok {
		m.expiries[key] = time.Now().Add(expiration)
	}
	return nil
}

func TestLimiter_AllowRequest_ByIP(t *testing.T) {
	store := newMockStore()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:       2,
		IpBlockDuration:    1,
		TokenMaxRequest:    5,
		TokenBlockDuration: 1,
	}
	limiter := NewLimiter(context.Background(), cfg, store)
	ip := "1.2.3.4"

	// First request should be allowed
	allowed := limiter.AllowRequest(ip, "")
	assert.True(t, allowed)

	// Second request should be allowed
	allowed = limiter.AllowRequest(ip, "")
	assert.True(t, allowed)

	// Third request should be blocked
	allowed = limiter.AllowRequest(ip, "")
	assert.False(t, allowed)
}

func TestLimiter_AllowRequest_ByToken(t *testing.T) {
	store := newMockStore()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:       2,
		IpBlockDuration:    1,
		TokenMaxRequest:    2,
		TokenBlockDuration: 1,
	}
	limiter := NewLimiter(context.Background(), cfg, store)
	token := "token123"

	// First request should be allowed
	allowed := limiter.AllowRequest("", token)
	assert.True(t, allowed)

	// Second request should be allowed
	allowed = limiter.AllowRequest("", token)
	assert.True(t, allowed)

	// Third request should be blocked
	allowed = limiter.AllowRequest("", token)
	assert.False(t, allowed)
}

func TestLimiter_AllowRequest_IPAndTokenSeparateLimits(t *testing.T) {
	store := newMockStore()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:       1,
		IpBlockDuration:    1,
		TokenMaxRequest:    1,
		TokenBlockDuration: 1,
	}
	limiter := NewLimiter(context.Background(), cfg, store)
	ip := "5.6.7.8"
	token := "tok456"

	// IP limit
	assert.True(t, limiter.AllowRequest(ip, ""))
	assert.False(t, limiter.AllowRequest(ip, ""))

	// Token limit
	assert.True(t, limiter.AllowRequest("", token))
	assert.False(t, limiter.AllowRequest("", token))
}

func TestLimiter_AllowRequest_ResetsAfterBlockDuration(t *testing.T) {
	store := newMockStore()
	cfg := config.RateLimiterConfig{
		IpMaxRequest:       1,
		IpBlockDuration:    1, // 1 second
		TokenMaxRequest:    1,
		TokenBlockDuration: 1,
	}
	limiter := NewLimiter(context.Background(), cfg, store)
	ip := "9.9.9.9"

	assert.True(t, limiter.AllowRequest(ip, ""))
	assert.False(t, limiter.AllowRequest(ip, ""))

	// Wait for block duration to expire
	time.Sleep(1100 * time.Millisecond)
	assert.True(t, limiter.AllowRequest(ip, ""))
}

func TestLimiter_AllowRequest_ErrorOnGet(t *testing.T) {
	// Store that always returns error on Get
	store := &errorStore{}
	cfg := config.RateLimiterConfig{
		IpMaxRequest:       1,
		IpBlockDuration:    1,
		TokenMaxRequest:    1,
		TokenBlockDuration: 1,
	}
	limiter := NewLimiter(context.Background(), cfg, store)
	assert.False(t, limiter.AllowRequest("errip", ""))
}

type errorStore struct{}

func (e *errorStore) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return nil
}
func (e *errorStore) Get(ctx context.Context, key string) (string, error) { return "", assert.AnError }
func (e *errorStore) Incr(ctx context.Context, key string) error          { return nil }
func (e *errorStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}
