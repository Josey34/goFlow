package pipeline

import (
	"context"
	"io"
	"os"

	"goflow/processor"
	"goflow/worker"
)

type Extractor struct {
	name string
}

func NewExtractor() *Extractor {
	return &Extractor{
		name: "Extractor",
	}
}

func (e *Extractor) Name() string {
	return e.name
}

func (e *Extractor) Process(ctx context.Context, input interface{}) (interface{}, error) {
	result := input.(*worker.ProcessingResult)

	if result.Error != nil {
		return nil, result.Error
	}
	defer result.File.Close()

	tempFile, err := os.CreateTemp("", "pdf-*.bin")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, result.File); err != nil {
		return nil, err
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, err
	}

	text, pageCount, err := processor.ExtractText(tempFile)
	if err != nil {
		return nil, err
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, err
	}

	hash, err := processor.ComputeHash(tempFile)
	if err != nil {
		return nil, err
	}

	return &ExtractorOutput{
		Task:          result.Task,
		ExtractedText: text,
		PageCount:     pageCount,
		FileHash:      hash,
	}, nil
}
