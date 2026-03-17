package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Retry(ctx context.Context, maxAttempts int, operation func() error) error {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		if attempt == maxAttempts-1 {
			return err
		}

		baseDelay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		jitter := time.Duration(rand.Intn(100)) * time.Millisecond
		delay := baseDelay + jitter

		fmt.Printf("Retrying in %v...\n", delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
