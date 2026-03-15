package pipeline

import (
	"context"
	"time"

	"goflow/entity"

	"github.com/google/uuid"
)

type Aggregator struct {
	name string
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		name: "Aggregator",
	}
}

func (a *Aggregator) Name() string {
	return a.name
}

func (a *Aggregator) Process(ctx context.Context, input interface{}) (interface{}, error) {
	dedupOutput := input.(*DeduplicatorOutput)

	result := &entity.ProcessingResult{
		ID:            uuid.New().String(),
		DocumentID:    dedupOutput.Task.Event.DocumentID,
		ExtractedText: dedupOutput.ExtractedText,
		PageCount:     dedupOutput.PageCount,
		FileHash:      dedupOutput.FileHash,
		IsDuplicate:   dedupOutput.IsDuplicate,
		ThumbnailInfo: "",
		ProcessedAt:   time.Now(),
		ErrorMessage:  "",
	}

	return &AggregatorOutput{
		Result: result,
		Chunks: dedupOutput.Chunks,
	}, nil
}
