package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"goflow/retry"
)

func TestRetrySuccess(t *testing.T) {
	attempts := 0
	err := retry.Retry(context.Background(), 3, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("fail")
		}
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetryExhausted(t *testing.T) {
	attempts := 0
	err := retry.Retry(context.Background(), 3, func() error {
		attempts++
		return errors.New("persistent failure")
	})

	if err == nil {
		t.Error("expected error after max attempts")
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := retry.Retry(ctx, 5, func() error {
		return errors.New("fail")
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context deadline exceeded, got %v", err)
	}
}
