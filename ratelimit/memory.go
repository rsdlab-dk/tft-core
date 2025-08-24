package ratelimit

import (
	"context"
	"sync"
	"time"
)

type MemoryLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*bucket
}

type bucket struct {
	count     int
	resetTime time.Time
}

func NewMemoryLimiter() *MemoryLimiter {
	limiter := &MemoryLimiter{
		buckets: make(map[string]*bucket),
	}

	go limiter.cleanup()

	return limiter
}

func (m *MemoryLimiter) Allow(ctx context.Context, key string, rate int, window time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	b, exists := m.buckets[key]
	if !exists || now.After(b.resetTime) {
		m.buckets[key] = &bucket{
			count:     1,
			resetTime: now.Add(window),
		}
		return true, nil
	}

	if b.count >= rate {
		return false, nil
	}

	b.count++
	return true, nil
}

func (m *MemoryLimiter) Reset(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.buckets, key)
	return nil
}

func (m *MemoryLimiter) GetCount(ctx context.Context, key string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, exists := m.buckets[key]
	if !exists || time.Now().After(b.resetTime) {
		return 0, nil
	}

	return b.count, nil
}

func (m *MemoryLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, bucket := range m.buckets {
			if now.After(bucket.resetTime) {
				delete(m.buckets, key)
			}
		}
		m.mu.Unlock()
	}
}
