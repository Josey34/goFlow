package pipeline

import (
	"context"
	"goflow/repository"
)

type Deduplicator struct {
	name       string
	resultRepo repository.ResultRepository
}

func NewDeduplicator(resultRepo repository.ResultRepository) *Deduplicator {
	return &Deduplicator{
		name:       "Deduplicator",
		resultRepo: resultRepo,
	}
}

func (d *Deduplicator) Name() string {
	return d.name
}

func (d *Deduplicator) Process(ctx context.Context, input interface{}) (interface{}, error) {
	chunkerOutput := input.(*ChunkerOutput)

	isDuplicate := false
	existing, err := d.resultRepo.FindByHash(ctx, chunkerOutput.FileHash)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		isDuplicate = true
	}

	return &DeduplicatorOutput{
		Task:          chunkerOutput.Task,
		ExtractedText: chunkerOutput.ExtractedText,
		PageCount:     chunkerOutput.PageCount,
		FileHash:      chunkerOutput.FileHash,
		IsDuplicate:   isDuplicate,
		Chunks:        chunkerOutput.Chunks,
	}, nil
}
