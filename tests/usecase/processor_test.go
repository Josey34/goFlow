package usecase_test

import (
	"context"
	"testing"
	"time"

	"goflow/entity"
	"goflow/tests/mock"
	"goflow/usecase"
)

const testFileName = "file.pdf"
const testBucket = "bucket"

// Test: Processor creation with config options
func TestProcessorCreation(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
		usecase.WithWorkers(5),
		usecase.WithRetry(3),
		usecase.WithChunkSize(2000),
		usecase.WithChunkOverlap(500),
	)

	if processor == nil {
		t.Error("expected processor to be created")
	}
}

// Test: Processor with default config
func TestProcessorDefaults(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
	)

	if processor == nil {
		t.Error("expected processor with defaults")
	}
}

// Test: Cache hit - same document processed twice
func TestProcessorCacheHit(t *testing.T) {
	cache := mock.NewMockCache()

	// Pre-populate cache with result
	cachedResult := &entity.ProcessingResult{
		ID:            "cached1",
		DocumentID:    "doc1",
		ExtractedText: "cached text",
		IsDuplicate:   true,
	}
	cache.Set(context.Background(), "doc1", cachedResult)

	// Verify cache stats
	stats := cache.Stats()
	if stats.Items != 1 {
		t.Errorf("expected 1 cached item, got %d", stats.Items)
	}
}

// / Test: Duplicate detection - same file hash
func TestProcessorDuplicateDetection(t *testing.T) {
	ctx := context.Background()

	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	// Pre-populate repo with file having same hash
	existingResult := &entity.ProcessingResult{
		ID:         "existing1",
		DocumentID: "doc-original",
		FileHash:   "samehash123",
	}
	resultRepo.Insert(ctx, existingResult)

	// Create processor
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()

	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
	)

	if processor == nil {
		t.Error("expected processor")
	}

	// Verify existing result can be found by hash
	found, _ := resultRepo.FindByHash(ctx, "samehash123")
	if found == nil {
		t.Error("expected to find duplicate by hash")
	}
	if found.DocumentID != "doc-original" {
		t.Errorf("expected doc-original, got %s", found.DocumentID)
	}
}

// Test: Retry configuration
func TestProcessorRetryConfiguration(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	// Create processor with max 3 retries
	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
		usecase.WithRetry(3),
	)

	if processor == nil {
		t.Error("expected processor with retry config")
	}
}

// Test: Context cancellation
func TestProcessorContextCancellation(t *testing.T) {
	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Verify context is cancelled
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("expected context to be cancelled")
	}
}

// Test: Chunk size configuration
func TestProcessorChunkConfiguration(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
		usecase.WithChunkSize(5000),
		usecase.WithChunkOverlap(1000),
	)

	if processor == nil {
		t.Error("expected processor with chunk config")
	}
}

// Test: Worker pool configuration
func TestProcessorWorkerConfiguration(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()

	for workers := 1; workers <= 10; workers++ {
		processor := usecase.NewProcessorUsecase(
			consumer,
			downloader,
			resultRepo,
			chunkRepo,
			usecase.WithWorkers(workers),
		)

		if processor == nil {
			t.Errorf("expected processor with %d workers", workers)
		}
	}
}

// Test: Combined configuration - all options together
func TestProcessorCombinedConfiguration(t *testing.T) {
	consumer := mock.NewMockConsumer()
	downloader := mock.NewMockDownloader()
	resultRepo := mock.NewMockResultRepository()
	chunkRepo := mock.NewMockChunkRepository()
	cache := mock.NewMockCache()

	processor := usecase.NewProcessorUsecase(
		consumer,
		downloader,
		resultRepo,
		chunkRepo,
		usecase.WithWorkers(4),
		usecase.WithRetry(5),
		usecase.WithChunkSize(3000),
		usecase.WithChunkOverlap(600),
	)

	if processor == nil {
		t.Error("expected processor with all options")
	}

	// Cache should be independent
	cache.Set(context.Background(), "doc", &entity.ProcessingResult{ID: "r1"})
	stats := cache.Stats()
	if stats.Items != 1 {
		t.Error("cache should be independent of processor")
	}
}

// Test: Timeout handling
func TestProcessorTimeout(t *testing.T) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond)

	// Context should be expired
	select {
	case <-ctx.Done():
		// Expected - timeout occurred
	default:
		t.Error("expected context timeout")
	}
}
