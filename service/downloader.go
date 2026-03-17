package service

import (
	"context"
	"io"
)

type FileDownloader interface {
	Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
}
