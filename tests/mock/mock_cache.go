package mock

import (
	"context"
	"goflow/entity"
	"goflow/service"
)

type MockCache struct {
	data   map[string]*entity.ProcessingResult
	hits   int64
	misses int64
}

func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]*entity.ProcessingResult),
	}
}

func (mc *MockCache) Get(ctx context.Context, documentID string) (*entity.ProcessingResult, bool) {
	result, ok := mc.data[documentID]
	if ok {
		mc.hits++
	} else {
		mc.misses++
	}
	return result, ok
}

func (mc *MockCache) Set(ctx context.Context, documentID string, result *entity.ProcessingResult) error {
	mc.data[documentID] = result
	return nil
}

func (mc *MockCache) Delete(ctx context.Context, documentID string) error {
	delete(mc.data, documentID)
	return nil
}

func (mc *MockCache) Stats() service.CacheStats {
	return service.CacheStats{
		Hits:   mc.hits,
		Misses: mc.misses,
		Items:  len(mc.data),
	}
}

func (mc *MockCache) Clear() error {
	mc.data = make(map[string]*entity.ProcessingResult)
	return nil
}
