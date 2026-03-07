# ⚡ GoFlow — Document Processor (Project 3 of 3)

> A concurrent document processing pipeline that consumes events from DocVault's SQS queue, extracts text from PDFs, detects duplicate documents, generates thumbnail info and statistics — completing the Document Management System.

## 🎯 System Overview

This is **Project 3** — the final piece that ties everything together:

```
┌──────────────────────────────────────────────────────────────────┐
│                  Document Management System                      │
│                                                                  │
│  Project 1: DocVault ✅      ──► User uploads PDF                │
│                               ──► Stored in MinIO + SQLite       │
│                               ──► Publishes "file.uploaded" SQS  │
│                                                                  │
│  Project 2: GoAuth ✅        ──► Protects DocVault routes        │
│                               ──► Users own their documents      │
│                                                                  │
│  Project 3: GoFlow           ──► Consumes "file.uploaded" from   │
│  (You are here)                  SQS                             │
│                               ──► Downloads file from MinIO      │
│                               ──► Extracts text from PDF         │
│                               ──► Chunks file into segments      │
│                               ──► Detects duplicate documents    │
│                               ──► Generates thumbnail info       │
│                               ──► Writes results + chunks to     │
│                                   shared SQLite                  │
│                               ──► Caches results for performance │
│                                                                  │
│  Flow:  Upload → SQS → GoFlow → Process → Results in SQLite     │
└──────────────────────────────────────────────────────────────────┘
```

### The Complete Flow

```
User uploads report.pdf via DocVault API (with JWT auth from GoAuth)
    │
    ├──► MinIO: file stored
    ├──► SQLite: metadata saved (documents table)
    └──► SQS: "file.uploaded" event published
              │
              ▼
    GoFlow consumes the event
    GoFlow downloads report.pdf from MinIO
              │
              ├──► Worker 1: Extract text from PDF
              ├──► Worker 2: Chunk text into segments (for search/RAG/indexing)
              ├──► Worker 3: Compute file hash → check for duplicates
              ├──► Worker 4: Extract page count, file info → thumbnail metadata
              │
              ▼
    Results + chunks written to SQLite (processing_results + document_chunks tables)
    Cache results for quick access
    Log processing stats
```

## 🧰 Tech Stack

| Tool | Purpose | Why This? |
|------|---------|-----------|
| **Go (standard library)** | Core | No HTTP framework — CLI/worker tool |
| **MinIO SDK** | Download files | Same MinIO as DocVault |
| **AWS SDK (SQS)** | Consume events | Same SQS as DocVault |
| **SQLite** | Store results | Same DB as DocVault + GoAuth |
| **pdfcpu** | PDF text extraction | Pure Go PDF library |
| **crypto/sha256** | Duplicate detection | Hash file contents |
| **testify** | Testing | Mocks + assertions |

---

## 📚 Go Concepts You Will Learn

### Advanced Concurrency Patterns (NEW — main focus)
- [ ] **Fan-out pattern** — one SQS event triggers multiple processing goroutines
- [ ] **Fan-in pattern** — multiple processing results merged into one output
- [ ] **Pipeline pattern** — stages: consume → download → process → store
- [ ] **Worker pool pattern** — N goroutines processing from shared job channel
- [ ] **Task queue** — producer-consumer with buffered channels
- [ ] **Goroutine lifecycle management** — start, stop, track
- [ ] **Race condition detection** — `go run -race`
- [ ] **Race condition prevention** — mutexes for shared state

### Channels Deep Dive (NEW depth)
- [ ] Unbuffered vs buffered channels
- [ ] Channel direction (`chan<-`, `<-chan`)
- [ ] Closing channels and why it matters
- [ ] Ranging over channels
- [ ] `select` with multiple channels, `default`, `time.After`
- [ ] Channel as semaphore
- [ ] Nil channels in select

### Sync Primitives (NEW)
- [ ] `sync.Mutex` — protect shared state
- [ ] `sync.RWMutex` — multiple readers OR one writer
- [ ] `sync.WaitGroup` — wait for goroutines
- [ ] `sync.Once` — initialize once
- [ ] When to use Mutex vs Channel

### Retry & Resilience (NEW)
- [ ] Exponential backoff with bit shifting (`1 << attempt`)
- [ ] Jitter (prevent thundering herd)
- [ ] `context.WithTimeout` per operation
- [ ] Panic recovery in goroutines

### Caching (NEW)
- [ ] In-memory cache from scratch
- [ ] TTL (Time-To-Live)
- [ ] Cache-aside pattern
- [ ] Thread-safe with `sync.RWMutex`
- [ ] Cache hit/miss statistics

### Rate Limiting (NEW)
- [ ] Semaphore pattern (buffered channels)
- [ ] `time.Ticker` based rate limiting

### File Chunking (NEW)
- [ ] Splitting large files into fixed-size chunks with overlap
- [ ] Chunk configuration (chunk size, overlap size)
- [ ] Page-aware chunking (track start/end page per chunk)
- [ ] Batch INSERT for chunks (efficient DB writes)
- [ ] Chunking strategies (fixed-size, by page, by paragraph)
- [ ] Parallel chunk processing across goroutines

### Functional Options Pattern (NEW)
- [ ] `type Option func(*Struct)` pattern
- [ ] Variadic `...Option` constructors
- [ ] `WithXxx()` naming convention
- [ ] Sensible defaults + selective overrides

### Unit Testing (Reinforced + Concurrency)
- [ ] Testing concurrent code with `-race`
- [ ] Testing channels and goroutines
- [ ] Testing with timeouts
- [ ] Mock interfaces
- [ ] Benchmark tests (`func BenchmarkXxx(b *testing.B)`)

### Clean Architecture (Reinforced — CLI/worker context)
- [ ] Entity, repository, service, usecase layers
- [ ] Factory pattern with functional options
- [ ] Dependency inversion in a non-HTTP context

### Observability (NEW)
- [ ] Correlation IDs through pipeline
- [ ] Structured logging
- [ ] Processing metrics (items/sec, errors, duration)
- [ ] Cache statistics

---

## 📁 Project Structure

```
goflow/
├── main.go                     # Entry point: parse args, create factory, start processor
├── go.mod / go.sum
├── .env                        # DB_PATH (same as DocVault!), MinIO creds, SQS endpoint
│
├── config/
│   └── config.go               # CLI flags + env vars
│
├── factory/
│   └── factory.go              # Creates all deps with functional options
│
├── entity/
│   ├── document.go             # Mirror DocVault's entity
│   ├── event.go                # SQS event struct
│   ├── processing_result.go    # ExtractedText, PageCount, FileHash, IsDuplicate, ThumbnailInfo
│   ├── chunk.go                # DocumentChunk: ID, DocumentID, ChunkIndex, Content, StartPage, EndPage
│   └── processing_stats.go     # TotalProcessed, Duplicates, Errors, AvgDuration
│
├── repository/
│   ├── repository.go           # Interface: ResultRepository (Insert, FindByDocID, FindByHash, GetStats)
│   ├── chunk_repository.go     # Interface: ChunkRepository (InsertBatch, FindByDocID, Search)
│   ├── sqlite_result.go        # SQLite impl — writes to shared DB
│   └── sqlite_chunk.go         # SQLite impl — writes to document_chunks table
│
├── service/
│   ├── consumer.go             # Interface: EventConsumer (Consume, Acknowledge)
│   ├── consumer_sqs.go         # SQS impl — consumes from DocVault's queue
│   ├── downloader.go           # Interface: FileDownloader (Download → io.ReadCloser)
│   ├── downloader_minio.go     # MinIO impl — downloads from DocVault's bucket
│   ├── cache.go                # Interface: CacheService (Get, Set, Delete, Stats)
│   ├── cache_memory.go         # In-memory with RWMutex + TTL
│   ├── limiter.go              # Interface: RateLimiter (Acquire, Release)
│   └── limiter_semaphore.go    # Semaphore implementation
│
├── usecase/
│   └── processor.go            # ProcessorUsecase with functional options
│                                #   NewProcessorUsecase(consumer, downloader, repo, ...Option)
│
├── pipeline/
│   ├── stage.go                # Stage interface
│   ├── consumer.go             # Stage 1: Consume SQS events
│   ├── downloader.go           # Stage 2: Download from MinIO (fan-out)
│   ├── extractor.go            # Stage 3: Extract text, hash, page count
│   ├── chunker.go              # Stage 4: Split extracted text into chunks
│   ├── deduplicator.go         # Stage 5: Check hash → detect duplicates
│   ├── aggregator.go           # Stage 6: Fan-in results
│   └── writer.go               # Stage 7: Write results + chunks to SQLite
│
├── worker/
│   ├── pool.go                 # Worker pool
│   └── task.go                 # ProcessingTask definition
│
├── processor/
│   ├── text_extractor.go       # PDF text extraction
│   ├── chunker.go              # Split text into chunks (fixed-size with overlap, or by page/paragraph)
│   ├── hasher.go               # SHA256 hash
│   └── thumbnail.go            # Page count, dimensions
│
├── retry/
│   └── retry.go                # Exponential backoff + jitter
│
├── safesync/
│   └── safe_map.go             # Generic thread-safe map
│
├── logger/
│   └── logger.go               # Structured logger with correlation IDs
│
├── cmd/
│   ├── process.go              # Start pipeline
│   ├── stats.go                # Show processing + cache stats
│   └── health.go               # Check SQS, MinIO, SQLite
│
└── tests/
    ├── mock/                   # MockConsumer, MockDownloader, MockResultRepo, MockCache, MockLimiter
    ├── usecase/                # Processor orchestration tests
    ├── pipeline/               # Stage tests with channels
    │   └── chunker_test.go     # Chunking logic + edge cases
    ├── worker/                 # Pool concurrency tests
    ├── service/                # Cache thread safety, limiter tests
    ├── processor/              # Text extractor, hasher tests
    ├── retry/                  # Backoff + cancellation tests
    ├── safesync/               # Concurrent map with -race
    └── benchmark/              # Cache + pool benchmarks
```

### Factory with Functional Options

```go
// usecase/processor.go
type Option func(*ProcessorUsecase)

func WithWorkers(n int) Option     { return func(uc *ProcessorUsecase) { uc.workers = n } }
func WithRetry(max int) Option     { return func(uc *ProcessorUsecase) { uc.maxRetries = max } }
func WithChunkSize(size int) Option { return func(uc *ProcessorUsecase) { uc.chunkSize = size } }
func WithChunkOverlap(overlap int) Option { return func(uc *ProcessorUsecase) { uc.chunkOverlap = overlap } }
func WithCache(c service.CacheService) Option { return func(uc *ProcessorUsecase) { uc.cache = c } }
func WithRateLimiter(l service.RateLimiter) Option { return func(uc *ProcessorUsecase) { uc.limiter = l } }
func WithVerbose(v bool) Option    { return func(uc *ProcessorUsecase) { uc.verbose = v } }

func NewProcessorUsecase(
    consumer service.EventConsumer,
    downloader service.FileDownloader,
    resultRepo repository.ResultRepository,
    opts ...Option,
) *ProcessorUsecase {
    uc := &ProcessorUsecase{
        consumer: consumer, downloader: downloader, resultRepo: resultRepo,
        workers: 3, maxRetries: 0, chunkSize: 1000, chunkOverlap: 200, verbose: false,
    }
    for _, opt := range opts { opt(uc) }
    return uc
}
```

```go
// factory/factory.go
func New(cfg *config.Config) *Factory {
    db, _ := database.OpenSQLite(cfg.DBPath)  // SAME DB!
    consumer := service.NewSQSConsumer(sqsClient, cfg.QueueURL)
    downloader := service.NewMinIODownloader(minioClient, cfg.BucketName)
    resultRepo := repository.NewSQLiteResultRepo(db)

    processorUC := usecase.NewProcessorUsecase(
        consumer, downloader, resultRepo,
        usecase.WithWorkers(cfg.Workers),
        usecase.WithRetry(cfg.MaxRetries),
        usecase.WithChunkSize(cfg.ChunkSize),
        usecase.WithChunkOverlap(cfg.ChunkOverlap),
        usecase.WithCache(service.NewMemoryCache(cfg.CacheTTL)),
        usecase.WithRateLimiter(service.NewSemaphoreLimiter(cfg.RateLimit)),
        usecase.WithVerbose(cfg.Verbose),
    )
    return &Factory{ProcessorUsecase: processorUC}
}
```

---

## 🗺️ Phase-by-Phase Roadmap

### Phase 1: CLI, Entity & Shared DB (Day 1)
**Goal:** CLI structure, entities, connect to DocVault's shared SQLite.

- `entity/processing_result.go`, `entity/event.go`, `entity/processing_stats.go`
- `config/config.go` with CLI flags
- `database/sqlite.go` → opens SAME DB as DocVault
- `cmd/health.go` → check SQS + MinIO + SQLite

**Test:** `goflow health` shows all services connected, same DB file as DocVault.

### Phase 2: SQS Consumer & MinIO Downloader (Day 2–3)
**Goal:** Consume DocVault's SQS events, download files from MinIO.

- `service/consumer.go` interface + `consumer_sqs.go`
- `service/downloader.go` interface + `downloader_minio.go`
- `worker/pool.go` — worker pool consumes events, downloads files

**Test:** Upload via DocVault → GoFlow logs "received file.uploaded for report.pdf" → file downloaded.

### Phase 3: Text Extraction, Chunking & Hashing + Functional Options (Day 4–5)
**Goal:** Process files — extract text, chunk it, compute hash, get page info. Implement usecase with functional options.

- `processor/text_extractor.go` — extract text from PDF
- `processor/chunker.go` — split extracted text into chunks:
  ```go
  type ChunkConfig struct {
      ChunkSize    int  // characters per chunk (e.g., 1000)
      ChunkOverlap int  // overlap between chunks (e.g., 200)
  }

  // Chunker splits text into overlapping segments
  func Chunk(text string, cfg ChunkConfig) []entity.DocumentChunk
  ```
  Strategies: fixed-size with overlap, by page, by paragraph
- `processor/hasher.go` — SHA256 hash of file contents
- `processor/thumbnail.go` — page count, dimensions
- `entity/chunk.go`:
  ```go
  type DocumentChunk struct {
      ID          string
      DocumentID  string
      ChunkIndex  int       // 0, 1, 2, ...
      Content     string    // the chunk text
      StartPage   int       // which page this chunk starts on
      EndPage     int       // which page this chunk ends on
      CharCount   int
  }
  ```
- `pipeline/extractor.go` + `pipeline/chunker.go` — pipeline stages
- `usecase/processor.go` with `WithWorkers()`, `WithRetry()`, `WithChunkSize()`, etc.

**Test:** PDF text extracted + split into chunks. `NewProcessorUsecase(consumer, downloader, repo)` works with defaults.

### Phase 4: Duplicate Detection & Chunk Storage (Day 6)
**Goal:** Check file hash for duplicates, store chunks in DB.

- `repository/repository.go` — `ResultRepository` interface + `ChunkRepository` interface
- `repository/sqlite_result.go` + `repository/sqlite_chunk.go` — SQLite impls
- `pipeline/deduplicator.go`
- `document_chunks` table in shared SQLite:
  ```sql
  CREATE TABLE IF NOT EXISTS document_chunks (
      id TEXT PRIMARY KEY,
      document_id TEXT NOT NULL,
      chunk_index INTEGER NOT NULL,
      content TEXT NOT NULL,
      start_page INTEGER,
      end_page INTEGER,
      char_count INTEGER,
      FOREIGN KEY (document_id) REFERENCES documents(id)
  );
  ```

**Test:** Same file twice → duplicate. Chunks stored with correct indices. Query chunks by document ID.

### Phase 5: Full Pipeline (Day 7)
**Goal:** All stages connected: consume → download → extract → chunk → deduplicate → aggregate → write.

- `pipeline/aggregator.go` (fan-in), `pipeline/writer.go`
- Writer stores BOTH `processing_results` AND `document_chunks`
- Full pipeline in usecase

**Test:** Upload PDF via DocVault → GoFlow processes → results in `processing_results` + chunks in `document_chunks`.

### Phase 6: Cache Service (Day 8)
**Goal:** Skip already-processed files.

- `service/cache.go` interface + `service/cache_memory.go` (RWMutex + TTL)
- Cache-aside in worker pool
- `goflow stats` shows hit/miss

**Test:** Re-upload same file → cache hit, no re-processing. `-race` clean.

### Phase 7: Retry & Rate Limiting (Day 9)
**Goal:** Retry failures, limit concurrency.

- `retry/retry.go` — exponential backoff + jitter
- `service/limiter.go` interface + `limiter_semaphore.go`

**Test:** Failed download retried 3x. Rate limit of 2 → only 2 concurrent downloads.

### Phase 8: Thread-Safe Map & Race Detection (Day 9 cont.)
**Goal:** SafeMap for pipeline stats, race detection practice.

- `safesync/safe_map.go` — generic `SafeMap[K, V]`
- Intentional race → detect → fix

**Test:** `-race` passes. Can explain WHY the race happened.

### Phase 9: Unit Tests & Benchmarks (Day 10–11)
**Goal:** Comprehensive tests for all layers.

- Mocks for all interfaces
- Usecase: process success, cache hit, duplicate, retry, cancellation
- Pipeline: each stage with channels
- Cache: thread safety with 100 goroutines
- Retry: Nth attempt success, context cancellation
- Benchmarks: cache read/write, pool throughput

**Test:** `go test ./...` + `go test -race ./...` + `go test -bench=.`

### Phase 10: Graceful Shutdown & Observability (Day 12)
**Goal:** Clean shutdown, correlation IDs, stats.

- Signal handling → cancel context → WaitGroup
- Panic recovery per worker
- Correlation ID via `context.WithValue`
- `goflow stats` → total processed, duplicates, errors, avg duration

**Test:** Ctrl+C → clean shutdown. `goflow stats` shows real data.

---

## 🏋️ Concurrency Exercises

### Exercise 1: Concurrent Logger (after Phase 2)
### Exercise 2: Race Condition Lab (Phase 8)
### Exercise 3: Task Queue (after Phase 2)
### Exercise 4: Channel Streams (after Phase 5)
Numbers 1-100 → filter even → square → sum. All channels.
### Exercise 5: Task Racer (after Phase 7)
Download from 3 mirrors, return first, cancel others.
### Exercise 6: Fan-In Combiner (after Phase 5)
Merge 5 channels into 1, ordered by timestamp.

---

## 🧪 Testing Cheat Sheet

```bash
# Start GoFlow (long-running, listens to SQS)
goflow process --workers 5 --retry 3 --chunk-size 1000 --chunk-overlap 200 --cache --verbose

# Upload via DocVault (triggers processing)
curl -H "Authorization: Bearer $TOKEN" -F "file=@test.pdf" http://localhost:8080/api/documents/upload

# Check results
sqlite3 docvault.db "SELECT * FROM processing_results ORDER BY processed_at DESC LIMIT 5;"

# Stats & health
goflow stats
goflow health

# Tests
go test ./...
go test -race ./...
go test -cover ./...
go test -bench=. ./tests/benchmark/
```

## 📖 Dependencies
```bash
go get github.com/minio/minio-go/v7
go get github.com/aws/aws-sdk-go-v2 github.com/aws/aws-sdk-go-v2/service/sqs github.com/aws/aws-sdk-go-v2/config
go get github.com/mattn/go-sqlite3
go get github.com/pdfcpu/pdfcpu
go get github.com/google/uuid
go get github.com/stretchr/testify
```

## 💡 Tutor Instructions
1. Draw channel diagrams — which goroutine sends where
2. Don't fix races — help me understand WHY
3. Suggest print statements to visualize concurrency
4. Challenge: "why buffered here vs unbuffered?"
5. If struggling with pipeline, do Exercise 4 first
6. Check: usecase must not import MinIO/SQS directly

### Common mistakes:
- Goroutine leak (forgetting to close channels)
- Deadlock (unbuffered channel, no receiver)
- Closing channels twice / sending on closed (panic)
- `wg.Add()` after `go func()` (race)
- Capturing loop variable in goroutine
- Not checking `ctx.Done()`
- Usecase importing `minio-go`

---

## ✅ Completion Checklist

- [ ] Phase 1: CLI, entity, shared DB, health
- [ ] Phase 2: SQS consumer + MinIO downloader (DocVault integration!)
- [ ] Phase 3: Text extraction + file chunking + hashing + functional options
- [ ] Phase 4: Duplicate detection + chunk storage via shared DB
- [ ] Phase 5: Full pipeline (consume → process → write)
- [ ] Phase 6: Cache service
- [ ] Phase 7: Retry + rate limiting
- [ ] Phase 8: Thread-safe map + race detection
- [ ] Phase 9: Unit tests + benchmarks
- [ ] Phase 10: Graceful shutdown, correlation IDs, stats
- [ ] Exercises 1–6

---

## 🎯 System Complete!

```
User registers (GoAuth) → logs in → gets JWT
    │
    ▼
User uploads PDF (DocVault) → MinIO + SQLite + SQS event
    │
    ▼
GoFlow consumes event → downloads PDF → extracts text → chunks file → detects dupes → stores results + chunks
    │
    ▼
User queries document → sees extracted text, duplicate status, processing stats
```

| Category | DocVault | GoAuth | GoFlow |
|----------|----------|--------|--------|
| Core Language | ✅ | ✅ | ✅ |
| HTTP & Web | ✅ Server | ✅ Server | — |
| Database | ✅ | ✅ | ✅ |
| Auth & Security | — | ✅ | — |
| Basic Concurrency | ✅ | — | ✅ |
| Advanced Concurrency | — | — | ✅ |
| Caching | — | — | ✅ |
| Resilience | ✅ | — | ✅ |
| Clean Architecture | ✅ | ✅ | ✅ |
| Factory DI | ✅ | ✅ | ✅ |
| Functional Options | — | — | ✅ |
| Unit Testing | ✅ | ✅ | ✅ |
| File Chunking | — | — | ✅ |
| Benchmarks | — | — | ✅ |
| Event-Driven | Publisher | — | Consumer |

**Total: ~35 days → Production-ready Go developer.** 🚀