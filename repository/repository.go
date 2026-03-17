package repository

import (
	"context"
	"goflow/entity"
)

type ResultRepository interface {
	Insert(ctx context.Context, result *entity.ProcessingResult) error
	FindByDocID(ctx context.Context, docID string) (*entity.ProcessingResult, error)
	FindByHash(ctx context.Context, hash string) (*entity.ProcessingResult, error)
}

type ChunkRepository interface {
	InsertBatch(ctx context.Context, chunks []entity.DocumentChunk) error
	FindByDocID(ctx context.Context, docID string) ([]entity.DocumentChunk, error)
	Search(ctx context.Context, query string) ([]entity.DocumentChunk, error)
}
