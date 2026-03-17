package service

import (
	"context"
	"goflow/entity"
	"sync"
	"time"
)

type cacheEntry struct {
	result    *entity.ProcessingResult
	expiresAt time.Time
}

type MemoryCache struct {
	mu     sync.RWMutex
	data   map[string]*cacheEntry
	ttl    time.Duration
	hits   int64
	misses int64
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*cacheEntry),
		ttl:  ttl,
	}
}

func (c *MemoryCache) Get(ctx context.Context, documentID string) (*entity.ProcessingResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[documentID]
	if !exists {
		c.misses++
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		c.misses++
		return nil, false
	}

	c.hits++
	return entry.result, true
}

func (c *MemoryCache) Set(ctx context.Context, documentID string, result *entity.ProcessingResult) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[documentID] = &cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(c.ttl),
	}

	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, documentID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, documentID)
	return nil
}

func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Hits:   c.hits,
		Misses: c.misses,
		Items:  len(c.data),
	}
}

func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*cacheEntry)
	return nil
}
