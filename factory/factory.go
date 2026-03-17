package factory

import (
	"context"
	"database/sql"
	"goflow/config"
	"goflow/database"
	"goflow/repository"
	"goflow/service"
	"goflow/usecase"
	"time"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Factory struct {
	DB               *sql.DB
	Config           *config.Config
	SQSConsumer      service.EventConsumer
	MinIODownloader  service.FileDownloader
	ResultRepository repository.ResultRepository
	ChunkRepository  repository.ChunkRepository
	ProcessorUC      usecase.ProcessorUsecase
	Cache            service.CacheService
}

func New(c *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(c.DBPath)
	if err != nil {
		return nil, err
	}

	if err := database.InitSchema(db); err != nil {
		return nil, err
	}

	resultRepo := repository.NewSQLiteResultRepo(db)
	chunkRepo := repository.NewSQLiteChunkRepo(db)

	awsCfg, err := awscfg.LoadDefaultConfig(context.Background(),
		awscfg.WithRegion(c.SQSRegion),
	)
	if err != nil {
		return nil, err
	}

	sqsClient := sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		o.BaseEndpoint = &c.SQSEndpoint
	})

	minioClient, err := minio.New(c.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIOAccessKey, c.MinIOSecretKey, ""),
		Secure: c.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}

	sqsConsumer := service.NewSQSConsumer(sqsClient, c.SQSQueueUrl)
	minioDownloader := service.NewMinIODownloader(minioClient, c.MinIOBucket)
	cache := service.NewMemoryCache(time.Duration(c.CacheTTL) * time.Second)

	processorUC := usecase.NewProcessorUsecase(
		sqsConsumer, minioDownloader, resultRepo, chunkRepo,
		usecase.WithWorkers(c.Workers),
		usecase.WithChunkSize(c.ChunkSize),
		usecase.WithChunkOverlap(c.ChunkOverlap),
	)

	return &Factory{
		DB:               db,
		Config:           c,
		SQSConsumer:      sqsConsumer,
		MinIODownloader:  minioDownloader,
		ResultRepository: resultRepo,
		ChunkRepository:  chunkRepo,
		ProcessorUC:      *processorUC,
		Cache:            cache,
	}, nil
}
