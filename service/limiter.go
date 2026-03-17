package service

import "context"

type RateLimiter interface {
	Acquire(ctx context.Context) error
	Release()
}
