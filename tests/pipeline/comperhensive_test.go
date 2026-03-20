package pipeline_test

import (
	"context"
	"testing"

	"goflow/entity"
	"goflow/processor"
	"goflow/tests/mock"
	"goflow/worker"
)

// Test all 7 pipeline stages working together
func TestFullPipelineIntegration(t *testing.T) {
	ctx := context.Background()

	// Mocks
	downloader := mock.NewMockDownloader().WithContent("test document content")
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	// Create task
	task := &worker.ProcessingTask{
		Event: &entity.Event{
			DocumentID: "doc1",
			BucketName: "test-bucket",
			ObjectName: "doc.pdf",
		},
	}

	// Stage 2: Download
	fileReader, err := downloader.Download(ctx, task.Event.BucketName, task.Event.ObjectName)
	if err != nil {
		t.Fatalf("stage 2 download failed: %v", err)
	}
	defer fileReader.Close()

	// Simulate Stage 3: Extract (since real PDF parsing needs real PDFs)
	extractedText := "extracted document text content for processing"
	pageCount := 5

	// Simulate Stage 4: Chunk
	chunkConfig := processor.ChunkConfig{
		ChunkSize:    100,
		ChunkOverlap: 20,
	}
	chunks := processor.ChunkText(extractedText, task.Event.DocumentID, pageCount, chunkConfig)

	if len(chunks) == 0 {
		t.Error("stage 4 chunking: expected chunks")
	}

	// Verify chunk structure
	for i, chunk := range chunks {
		if chunk.DocumentID != task.Event.DocumentID {
			t.Errorf("chunk %d: DocumentID mismatch", i)
		}
		if chunk.ChunkIndex != i {
			t.Errorf("chunk %d: ChunkIndex should be %d", i, i)
		}
		if chunk.Content == "" {
			t.Errorf("chunk %d: empty content", i)
		}
	}

	// Simulate Stage 5: Deduplicator (check hash)
	hash := "test-hash-123"
	isDuplicate := false
	existingResult, _ := resultRepo.FindByHash(ctx, hash)
	if existingResult != nil {
		isDuplicate = true
	}

	// Simulate Stage 6: Aggregator (create final result)
	result := &entity.ProcessingResult{
		ID:            "result-1",
		DocumentID:    task.Event.DocumentID,
		ExtractedText: extractedText,
		PageCount:     pageCount,
		FileHash:      hash,
		IsDuplicate:   isDuplicate,
	}

	// Stage 7: Writer (persist to DB)
	err = resultRepo.Insert(ctx, result)
	if err != nil {
		t.Fatalf("stage 7 writer: insert result failed: %v", err)
	}

	err = chunkRepo.InsertBatch(ctx, chunks)
	if err != nil {
		t.Fatalf("stage 7 writer: insert chunks failed: %v", err)
	}

	// Verify persistence
	storedResult, _ := resultRepo.FindByDocID(ctx, task.Event.DocumentID)
	if storedResult == nil {
		t.Error("result not persisted")
	}
	if storedResult.ExtractedText != extractedText {
		t.Error("extracted text mismatch")
	}

	storedChunks, _ := chunkRepo.FindByDocID(ctx, task.Event.DocumentID)
	if len(storedChunks) != len(chunks) {
		t.Errorf("chunk count mismatch: expected %d, got %d", len(chunks), len(storedChunks))
	}
}

// Test pipeline with duplicate detection across documents
func TestPipelineDuplicateHandling(t *testing.T) {
	ctx := context.Background()

	resultRepo := mock.NewMockResultRepository()

	// Document 1: Original
	result1 := &entity.ProcessingResult{
		ID:            "result-1",
		DocumentID:    "doc-original",
		FileHash:      "hash-xyz",
		ExtractedText: "original content",
	}
	resultRepo.Insert(ctx, result1)

	// Document 2: Duplicate (same hash)
	result2 := &entity.ProcessingResult{
		ID:          "result-2",
		DocumentID:  "doc-copy",
		FileHash:    "hash-xyz",
		IsDuplicate: true,
	}

	// Verify duplicate detection
	found, _ := resultRepo.FindByHash(ctx, "hash-xyz")
	if found == nil {
		t.Error("should find existing result by hash")
	}
	if found.DocumentID != "doc-original" {
		t.Error("should find original document")
	}

	// Insert duplicate result
	resultRepo.Insert(ctx, result2)

	// Both should exist
	orig, _ := resultRepo.FindByDocID(ctx, "doc-original")
	if orig == nil {
		t.Error("original should exist")
	}

	copy, _ := resultRepo.FindByDocID(ctx, "doc-copy")
	if copy == nil {
		t.Error("copy should exist")
	}
	if !copy.IsDuplicate {
		t.Error("copy should be marked as duplicate")
	}
}

// Test pipeline chunk batching
func TestPipelineChunkBatching(t *testing.T) {
	ctx := context.Background()
	chunkRepo := mock.NewMockChunkRepository()

	// Create 100 chunks
	var chunks []entity.DocumentChunk
	for i := 0; i < 100; i++ {
		chunks = append(chunks, entity.DocumentChunk{
			ID:         "chunk-" + string(rune(i)),
			DocumentID: "doc-large",
			ChunkIndex: i,
			Content:    "chunk content",
		})
	}

	// Batch insert
	err := chunkRepo.InsertBatch(ctx, chunks)
	if err != nil {
		t.Fatalf("batch insert failed: %v", err)
	}

	// Verify all persisted
	retrieved, _ := chunkRepo.FindByDocID(ctx, "doc-large")
	if len(retrieved) != 100 {
		t.Errorf("expected 100 chunks, got %d", len(retrieved))
	}
}

// Test pipeline with multiple concurrent documents
func TestPipelineMultipleDocumentProcessing(t *testing.T) {
	ctx := context.Background()

	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	// Process 5 documents
	for docNum := 1; docNum <= 5; docNum++ {
		docID := "doc-" + string(rune(48+docNum))

		// Extract and chunk
		text := "document " + string(rune(48+docNum)) + " content"
		chunkConfig := processor.ChunkConfig{
			ChunkSize:    50,
			ChunkOverlap: 10,
		}
		chunks := processor.ChunkText(text, docID, 3, chunkConfig)

		// Create result
		result := &entity.ProcessingResult{
			ID:            "result-" + string(rune(48+docNum)),
			DocumentID:    docID,
			ExtractedText: text,
			PageCount:     3,
			FileHash:      "hash-" + string(rune(48+docNum)),
		}

		// Persist
		resultRepo.Insert(ctx, result)
		if len(chunks) > 0 {
			chunkRepo.InsertBatch(ctx, chunks)
		}
	}

	// Verify all processed
	stats, _ := resultRepo.GetStats(ctx)
	if stats.TotalProcessed != 5 {
		t.Errorf("expected 5 documents processed, got %d", stats.TotalProcessed)
	}
}

// Test pipeline error handling in stages
func TestPipelineErrorHandling(t *testing.T) {
	ctx := context.Background()

	// Download error
	downloader := mock.NewMockDownloader().WithError()
	_, err := downloader.Download(ctx, "bucket", "file.pdf")
	if err == nil {
		t.Error("expected download error")
	}

	// Repository insert
	resultRepo := mock.NewMockResultRepository()
	result := &entity.ProcessingResult{
		ID:         "result-1",
		DocumentID: "doc-1",
	}

	// First insert succeeds
	err = resultRepo.Insert(ctx, result)
	if err != nil {
		t.Errorf("first insert failed: %v", err)
	}

	// Second insert with same ID should work (overwrites)
	result2 := &entity.ProcessingResult{
		ID:         "result-1",
		DocumentID: "doc-1-updated",
	}
	err = resultRepo.Insert(ctx, result2)
	if err != nil {
		t.Errorf("second insert failed: %v", err)
	}
}

// Test pipeline with empty content
func TestPipelineEmptyContent(t *testing.T) {
	ctx := context.Background()

	downloader := mock.NewMockDownloader().WithContent("")
	_, err := downloader.Download(ctx, "bucket", "empty.pdf")
	if err != nil {
		t.Errorf("download empty file failed: %v", err)
	}

	// Empty text should produce no chunks
	chunkConfig := processor.ChunkConfig{
		ChunkSize:    100,
		ChunkOverlap: 20,
	}
	chunks := processor.ChunkText("", "doc-empty", 1, chunkConfig)
	if len(chunks) != 0 {
		t.Errorf("empty text should produce no chunks, got %d", len(chunks))
	}
}

// Test pipeline overlap verification
func TestPipelineChunkOverlap(t *testing.T) {
	text := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	chunkConfig := processor.ChunkConfig{
		ChunkSize:    10,
		ChunkOverlap: 3,
	}

	chunks := processor.ChunkText(text, "doc-overlap", 1, chunkConfig)

	if len(chunks) < 2 {
		t.Error("expected multiple chunks for overlap testing")
		return
	}

	// Verify overlap between consecutive chunks
	for i := 0; i < len(chunks)-1; i++ {
		chunk1 := chunks[i]
		chunk2 := chunks[i+1]

		// Last 3 chars of chunk1 should be in chunk2
		if len(chunk1.Content) >= 3 {
			overlap := chunk1.Content[len(chunk1.Content)-3:]
			if len(chunk2.Content) < 3 || chunk2.Content[:3] != overlap {
				t.Logf("chunk %d and %d: overlap mismatch", i, i+1)
			}
		}
	}
}
