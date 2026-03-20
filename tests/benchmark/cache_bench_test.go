package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"

	"goflow/entity"
	"goflow/service"
)

func BenchmarkCacheGet(b *testing.B) {
	cache := service.NewMemoryCache(5 * time.Minute)
	ctx := context.Background()

	result := &entity.ProcessingResult{ID: "test"}
	cache.Set(ctx, "test", result)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, "test")
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache := service.NewMemoryCache(5 * time.Minute)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := &entity.ProcessingResult{ID: fmt.Sprintf("doc%d", i)}
		cache.Set(ctx, fmt.Sprintf("doc%d", i), result)
	}
}
