package usecase

import (
	"context"

	"goflow/entity"
	"goflow/processor"
	"goflow/repository"
	"goflow/service"

	"github.com/google/uuid"
)

type ProcessorUsecase struct {
	consumer   service.EventConsumer
	downloader service.FileDownloader
	resultRepo repository.ResultRepository
	chunkRepo  repository.ChunkRepository

	workers      int
	maxRetries   int
	chunkSize    int
	chunkOverlap int
}

type Option func(*ProcessorUsecase)

func WithWorkers(n int) Option {
	return func(uc *ProcessorUsecase) { uc.workers = n }
}

func WithRetry(max int) Option {
	return func(uc *ProcessorUsecase) { uc.maxRetries = max }
}

func WithChunkSize(size int) Option {
	return func(uc *ProcessorUsecase) { uc.chunkSize = size }
}

func WithChunkOverlap(overlap int) Option {
	return func(uc *ProcessorUsecase) { uc.chunkOverlap = overlap }
}

func NewProcessorUsecase(
	consumer service.EventConsumer,
	downloader service.FileDownloader,
	resultRepo repository.ResultRepository,
	chunkRepo repository.ChunkRepository,
	opts ...Option,
) *ProcessorUsecase {
	uc := &ProcessorUsecase{
		consumer:     consumer,
		downloader:   downloader,
		resultRepo:   resultRepo,
		chunkRepo:    chunkRepo,
		workers:      3,
		maxRetries:   0,
		chunkSize:    1000,
		chunkOverlap: 200,
	}
	for _, opt := range opts {
		opt(uc)
	}
	return uc
}

func (uc *ProcessorUsecase) Process(ctx context.Context, event *entity.Event) error {
	objectName := event.ObjectName
	if objectName == "" {
		objectName = event.Filename
	}

	fileReader, err := uc.downloader.Download(ctx, event.BucketName, objectName)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	text, pageCount, err := processor.ExtractText(fileReader)
	if err != nil {
		return err
	}

	fileReader2, err := uc.downloader.Download(ctx, event.BucketName, objectName)
	if err != nil {
		return err
	}
	defer fileReader2.Close()

	hash, err := processor.ComputeHash(fileReader2)
	if err != nil {
		return err
	}

	existing, err := uc.resultRepo.FindByHash(ctx, hash)
	if err != nil {
		return err
	}

	isDuplicate := existing != nil

	chunkCfg := processor.ChunkConfig{
		ChunkSize:    uc.chunkSize,
		ChunkOverlap: uc.chunkOverlap,
	}
	chunks := processor.ChunkText(text, event.DocumentID, pageCount, chunkCfg)

	thumbnail := processor.ExtractThumbnailInfo(pageCount)

	result := &entity.ProcessingResult{
		ID:            uuid.New().String(),
		DocumentID:    event.DocumentID,
		ExtractedText: text,
		PageCount:     pageCount,
		FileHash:      hash,
		IsDuplicate:   isDuplicate,
		ThumbnailInfo: thumbnail.Title,
		ErrorMessage:  "",
	}

	if err := uc.resultRepo.Insert(ctx, result); err != nil {
		return err
	}

	if err := uc.chunkRepo.InsertBatch(ctx, chunks); err != nil {
		return err
	}

	return nil
}
