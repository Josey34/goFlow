package mock

import (
	"context"
	"goflow/entity"
)

type MockChunkRepository struct {
	chunks map[string][]entity.DocumentChunk
}

func NewMockChunkRepository() *MockChunkRepository {
	return &MockChunkRepository{
		chunks: make(map[string][]entity.DocumentChunk),
	}
}

func (mcr *MockChunkRepository) InsertBatch(ctx context.Context, chunks []entity.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}
	docID := chunks[0].DocumentID
	mcr.chunks[docID] = chunks
	return nil
}

func (mcr *MockChunkRepository) FindByDocID(ctx context.Context, docID string) ([]entity.DocumentChunk, error) {
	return mcr.chunks[docID], nil
}

func (mcr *MockChunkRepository) Search(ctx context.Context, query string) ([]entity.DocumentChunk, error) {
	return []entity.DocumentChunk{}, nil
}
