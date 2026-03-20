package pipeline_test

import (
	"io"
	"strings"
	"testing"

	"goflow/processor"
)

func TestComputeHash(t *testing.T) {
	content := "test content for hashing"
	reader := io.NopCloser(strings.NewReader(content))
	defer reader.Close()

	hash, err := processor.ComputeHash(reader)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("expected hash value")
	}
}

func TestHashConsistent(t *testing.T) {
	content := "consistent content"

	reader1 := io.NopCloser(strings.NewReader(content))
	hash1, _ := processor.ComputeHash(reader1)
	reader1.Close()

	reader2 := io.NopCloser(strings.NewReader(content))
	hash2, _ := processor.ComputeHash(reader2)
	reader2.Close()

	if hash1 != hash2 {
		t.Error("hashes should be consistent for same content")
	}
}
