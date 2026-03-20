package mock

import (
	"context"
	"goflow/entity"
)

type MockResultRepository struct {
	results map[string]*entity.ProcessingResult
}

func NewMockResultRepository() *MockResultRepository {
	return &MockResultRepository{
		results: make(map[string]*entity.ProcessingResult),
	}
}

func (mr *MockResultRepository) Insert(ctx context.Context, result *entity.ProcessingResult) error {
	mr.results[result.ID] = result
	return nil
}

func (mr *MockResultRepository) FindByDocID(ctx context.Context, docID string) (*entity.ProcessingResult, error) {
	for _, r := range mr.results {
		if r.DocumentID == docID {
			return r, nil
		}
	}
	return nil, nil
}

func (mr *MockResultRepository) FindByHash(ctx context.Context, hash string) (*entity.ProcessingResult, error) {
	for _, r := range mr.results {
		if r.FileHash == hash {
			return r, nil
		}
	}
	return nil, nil
}

func (mr *MockResultRepository) GetStats(ctx context.Context) (*entity.ProcessingStats, error) {
	return &entity.ProcessingStats{
		TotalProcessed:    len(mr.results),
		DuplicatesFound:   0,
		ErrorsEncountered: 0,
		AvgProcessingTime: 0,
	}, nil
}
