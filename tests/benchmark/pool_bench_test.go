package benchmark

import (
	"context"
	"fmt"
	"testing"

	"goflow/retry"
	"goflow/safesync"
	"goflow/service"
)

// Benchmark SafeMap concurrent access
func BenchmarkSafeMapGet(b *testing.B) {
	sm := safesync.New[string, string]()
	sm.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Get("key")
	}
}

func BenchmarkSafeMapSet(b *testing.B) {
	sm := safesync.New[string, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Set(fmt.Sprintf("key%d", i), i)
	}
}

// Benchmark retry logic
func BenchmarkRetrySuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		retry.Retry(context.Background(), 3, func() error {
			return nil
		})
	}
}

// Benchmark rate limiter
func BenchmarkLimiterAcquireRelease(b *testing.B) {
	limiter := service.NewSemaphoreLimiter(5)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Acquire(ctx)
		limiter.Release()
	}
}
