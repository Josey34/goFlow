package service

import (
	"context"
	"goflow/entity"
)

type CacheStats struct {
	Hits   int64
	Misses int64
	Items  int
}

type CacheService interface {
	Get(ctx context.Context, documentID string) (*entity.ProcessingResult, bool)
	Set(ctx context.Context, documentID string, result *entity.ProcessingResult) error
	Delete(ctx context.Context, documentID string) error
	Stats() CacheStats
	Clear() error
}
