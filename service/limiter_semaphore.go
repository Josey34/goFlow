package service

import "context"

type SemaphoreLimiter struct {
	semaphore chan struct{}
}

func NewSemaphoreLimiter(maxConcurrent int) *SemaphoreLimiter {
	return &SemaphoreLimiter{
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

func (sl *SemaphoreLimiter) Acquire(ctx context.Context) error {
	select {
	case sl.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (sl *SemaphoreLimiter) Release() {
	<-sl.semaphore
}
