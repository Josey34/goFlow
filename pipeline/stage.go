package pipeline

import (
	"context"
	"goflow/entity"
	"goflow/worker"
)

type Stage interface {
	Process(ctx context.Context, input interface{}) (interface{}, error)
	Name() string
}

type ExtractorOutput struct {
	Task          *worker.ProcessingTask
	ExtractedText string
	PageCount     int
	FileHash      string
}

type ChunkerOutput struct {
	Task          *worker.ProcessingTask
	ExtractedText string
	PageCount     int
	FileHash      string
	Chunks        []entity.DocumentChunk
}

type DeduplicatorOutput struct {
	Task          *worker.ProcessingTask
	ExtractedText string
	PageCount     int
	FileHash      string
	IsDuplicate   bool
	Chunks        []entity.DocumentChunk
}

type AggregatorOutput struct {
	Result *entity.ProcessingResult
	Chunks []entity.DocumentChunk
}

type StageConfig struct {
	ChunkSize    int
	ChunkOverlap int
}
