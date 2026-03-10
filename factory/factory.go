package factory

import (
	"context"
	"database/sql"
	"goflow/config"
	"goflow/database"
	"goflow/service"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Factory struct {
	DB              *sql.DB
	Config          *config.Config
	SQSConsumer     service.EventConsumer
	MinIODownloader service.FileDownloader
}

func New(c *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(c.DBPath)
	if err != nil {
		return nil, err
	}

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

	return &Factory{
		DB:              db,
		Config:          c,
		SQSConsumer:     service.NewSQSConsumer(sqsClient, c.SQSQueueUrl),
		MinIODownloader: service.NewMinIODownloader(minioClient, c.MinIOBucket),
	}, nil
}
