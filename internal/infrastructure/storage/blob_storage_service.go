package storage

import (
	"context"
	"io"

	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var useSSL = false

type blobStorageService struct {
	config *common.Config
}

func NewBlobStorageService(cfg *common.Config) *blobStorageService {
	return &blobStorageService{
		config: cfg,
	}
}

func (s *blobStorageService) GetObject(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	minioClient, err := minio.New(s.config.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.config.Storage.AccessKey, s.config.Storage.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	objResponse, err := minioClient.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer objResponse.Close()

	respByte, err := io.ReadAll(objResponse)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

func (s *blobStorageService) UploadObject(ctx context.Context, bucketName, objectKey string, object io.Reader, size int64) error {
	minioClient, err := minio.New(s.config.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.config.Storage.AccessKey, s.config.Storage.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	_, err = minioClient.PutObject(ctx, bucketName, objectKey, object, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
