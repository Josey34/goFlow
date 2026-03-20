package pipeline

import (
	"context"
	"goflow/repository"
	"goflow/service"
)

type Deduplicator struct {
	name       string
	resultRepo repository.ResultRepository
	cache      service.CacheService
}

func NewDeduplicator(resultRepo repository.ResultRepository, cache service.CacheService) *Deduplicator {
	return &Deduplicator{
		name:       "Deduplicator",
		resultRepo: resultRepo,
		cache:      cache,
	}
}

func (d *Deduplicator) Name() string {
	return d.name
}

func (d *Deduplicator) Process(ctx context.Context, input interface{}) (interface{}, error) {
	chunkerOutput := input.(*ChunkerOutput)
	documentID := chunkerOutput.Task.Event.DocumentID

	if _, found := d.cache.Get(ctx, documentID); found {
		return &DeduplicatorOutput{
			Task:          chunkerOutput.Task,
			ExtractedText: chunkerOutput.ExtractedText,
			PageCount:     chunkerOutput.PageCount,
			FileHash:      chunkerOutput.FileHash,
			IsDuplicate:   true,
			Chunks:        chunkerOutput.Chunks,
		}, nil
	}

	isDuplicate := false
	existing, err := d.resultRepo.FindByHash(ctx, chunkerOutput.FileHash)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		isDuplicate = true
		d.cache.Set(ctx, documentID, existing)
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
