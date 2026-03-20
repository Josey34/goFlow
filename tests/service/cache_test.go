package service_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"goflow/entity"
	"goflow/service"
)

func TestCacheThreadSafety(t *testing.T) {
	cache := service.NewMemoryCache(5 * time.Second)
	ctx := context.Background()

	var wg sync.WaitGroup

	// 50 goroutines writing
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result := &entity.ProcessingResult{
				ID:            fmt.Sprintf("doc%d", id),
				DocumentID:    fmt.Sprintf("doc%d", id),
				ExtractedText: "test",
			}
			cache.Set(ctx, fmt.Sprintf("doc%d", id), result)
		}(i)
	}

	// 50 goroutines reading
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Get(ctx, fmt.Sprintf("doc%d", id%25))
		}(i)
	}

	wg.Wait()

	stats := cache.Stats()
	if stats.Items == 0 {
		t.Error("cache should have items")
	}
}

func TestCacheTTL(t *testing.T) {
	cache := service.NewMemoryCache(100 * time.Millisecond)
	ctx := context.Background()

	result := &entity.ProcessingResult{ID: "doc1"}
	cache.Set(ctx, "doc1", result)

	// Should exist immediately
	_, found := cache.Get(ctx, "doc1")
	if !found {
		t.Error("cache should have item immediately after set")
	}

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, found = cache.Get(ctx, "doc1")
	if found {
		t.Error("cache item should be expired after TTL")
	}
}
