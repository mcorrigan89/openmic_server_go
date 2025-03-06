package external

import (
	"context"
	"io"

	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type ImageMediaService interface {
	GetImageDataByID(ctx context.Context, imageEntity *entities.ImageEntity, rendition string) ([]byte, string, error)
	UploadImage(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) (*entities.ImageEntity, error)
}
