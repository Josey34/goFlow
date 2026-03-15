package pipeline

import (
	"context"

	"goflow/repository"
)

type Writer struct {
	name       string
	resultRepo repository.ResultRepository
	chunkRepo  repository.ChunkRepository
}

func NewWriter(
	resultRepo repository.ResultRepository,
	chunkRepo repository.ChunkRepository,
) *Writer {
	return &Writer{
		name:       "Writer",
		resultRepo: resultRepo,
		chunkRepo:  chunkRepo,
	}
}

func (w *Writer) Name() string {
	return w.name
}

func (w *Writer) Process(ctx context.Context, input interface{}) (interface{}, error) {
	aggOutput := input.(*AggregatorOutput)

	if err := w.resultRepo.Insert(ctx, aggOutput.Result); err != nil {
		return nil, err
	}

	if len(aggOutput.Chunks) > 0 {
		if err := w.chunkRepo.InsertBatch(ctx, aggOutput.Chunks); err != nil {
			return nil, err
		}
	}

	return aggOutput, nil
}
