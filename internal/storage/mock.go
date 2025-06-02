package storage

import (
	"context"
	"sync"
	"time"
)

type MockStore struct {
	Store
	Data map[string]mockEntry
	mu   sync.Mutex
}

type mockEntry struct {
	Value  int64
	Expiry time.Duration
}

func NewMockStore() *MockStore {
	return &MockStore{
		Data: make(map[string]mockEntry),
	}
}

func (m *MockStore) Get(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.Data[key]
	if !ok {
		return 0, nil
	}

	return entry.Value, nil
}

func (m *MockStore) Set(ctx context.Context, key string, value int64, expiry time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Data[key] = mockEntry{Value: value, Expiry: expiry}
	return nil
}

func (m *MockStore) Incr(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.Data[key]
	if !ok {
		entry = mockEntry{Value: 0}
	}
	entry.Value = entry.Value + 1
	m.Data[key] = entry
	return nil
}

func (m *MockStore) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Data, key)
	return nil
}

func (m *MockStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.Data[key]
	if !ok {
		return nil
	}
	entry.Expiry = expiration
	m.Data[key] = entry
	return nil
}
