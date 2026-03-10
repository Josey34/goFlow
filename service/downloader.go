package service

import (
	"context"
	"io"
)

type FileDownloader interface {
	Download(ctx context.Context, documentID, filename string) (io.ReadCloser, error)
}
