package mock

import (
	"context"
	"errors"
	"io"
	"strings"
)

type MockDownloader struct {
	fileContent string
	shouldError bool
}

func NewMockDownloader() *MockDownloader {
	return &MockDownloader{
		fileContent: "sample pdf content for testing",
	}
}

func (md *MockDownloader) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if md.shouldError {
		return nil, errors.New("consume failed")
	}

	return io.NopCloser(strings.NewReader(md.fileContent)), nil
}

func (md *MockDownloader) WithContent(content string) *MockDownloader {
	md.fileContent = content
	return md
}

func (md *MockDownloader) WithError() *MockDownloader {
	md.shouldError = true
	return md
}
