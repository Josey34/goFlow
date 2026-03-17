package pipeline

import (
	"context"
	"goflow/processor"
)

type Chunker struct {
	name   string
	config StageConfig
}

func NewChunker(config StageConfig) *Chunker {
	return &Chunker{
		name:   "Chunker",
		config: config,
	}
}

func (c *Chunker) Name() string {
	return c.name
}

func (c *Chunker) Process(ctx context.Context, input interface{}) (interface{}, error) {
	extractorOutput := input.(*ExtractorOutput)

	chunkCfg := processor.ChunkConfig{
		ChunkSize:    c.config.ChunkSize,
		ChunkOverlap: c.config.ChunkOverlap,
	}

	chunks := processor.ChunkText(
		extractorOutput.ExtractedText,
		extractorOutput.Task.Event.DocumentID,
		extractorOutput.PageCount,
		chunkCfg,
	)

	return &ChunkerOutput{
		Task:          extractorOutput.Task,
		ExtractedText: extractorOutput.ExtractedText,
		PageCount:     extractorOutput.PageCount,
		FileHash:      extractorOutput.FileHash,
		Chunks:        chunks,
	}, nil
}
