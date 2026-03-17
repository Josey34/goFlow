package service

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinIODownloader struct {
	client     *minio.Client
	bucketName string
}

func NewMinIODownloader(client *minio.Client, bucketName string) *MinIODownloader {
	return &MinIODownloader{
		client:     client,
		bucketName: bucketName,
	}
}

func (d *MinIODownloader) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	obj, err := d.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
