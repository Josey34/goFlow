package processor_test

import (
	"strings"
	"testing"

	"goflow/processor"
)

func TestChunkText(t *testing.T) {
	text := "This is a long text that needs to be chunked into smaller pieces for processing and analysis."
	cfg := processor.ChunkConfig{
		ChunkSize:    20,
		ChunkOverlap: 5,
	}

	chunks := processor.ChunkText(text, "doc1", 5, cfg)

	if len(chunks) == 0 {
		t.Error("expected chunks to be created")
	}

	for i, chunk := range chunks {
		if chunk.DocumentID != "doc1" {
			t.Errorf("expected DocumentID=doc1, got %s", chunk.DocumentID)
		}
		if chunk.ChunkIndex != i {
			t.Errorf("expected ChunkIndex=%d, got %d", i, chunk.ChunkIndex)
		}
		if chunk.Content == "" {
			t.Error("expected non-empty chunk content")
		}
	}
}

func TestChunkEmptyText(t *testing.T) {
	text := ""
	cfg := processor.ChunkConfig{
		ChunkSize:    100,
		ChunkOverlap: 10,
	}

	chunks := processor.ChunkText(text, "doc1", 1, cfg)

	if len(chunks) > 0 {
		t.Error("expected no chunks for empty text")
	}
}

func TestChunkWithOverlap(t *testing.T) {
	text := "abcdefghijklmnopqrstuvwxyz"
	cfg := processor.ChunkConfig{
		ChunkSize:    5,
		ChunkOverlap: 2,
	}

	chunks := processor.ChunkText(text, "doc1", 1, cfg)

	if len(chunks) > 1 {
		chunk1 := chunks[0]
		chunk2 := chunks[1]

		if !strings.Contains(chunk2.Content, chunk1.Content[len(chunk1.Content)-2:]) {
			t.Error("expected overlap between consecutive chunks")
		}
	}
}
