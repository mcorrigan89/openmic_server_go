package external

import (
	"context"
	"io"
)

type BlobStorageService interface {
	GetObject(ctx context.Context, bucketName, objectKey string) ([]byte, error)
	UploadObject(ctx context.Context, bucketName, objectKey string, object io.Reader, size int64) error
}
