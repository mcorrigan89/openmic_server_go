package media

import (
	"context"
	"io"
	"math"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/external"
)

type imageMediaService struct {
	blobStorageService external.BlobStorageService
}

func NewImageMediaService(blobStorageService external.BlobStorageService) *imageMediaService {
	return &imageMediaService{
		blobStorageService: blobStorageService,
	}
}

func (service *imageMediaService) GetImageDataByID(ctx context.Context, imageEntity *entities.ImageEntity, rendition string) ([]byte, string, error) {
	imageData, err := service.blobStorageService.GetObject(ctx, imageEntity.BucketName, imageEntity.ObjectID)
	if err != nil {
		return nil, "", err
	}

	image := bimg.NewImage(imageData)

	metadata, err := image.Metadata()
	if err != nil {
		return nil, "", err
	}

	newWidth, newHeight := calculateDimensions(metadata.Size.Width, metadata.Size.Height, getRendition(rendition))

	processedImage, err := image.Process(bimg.Options{
		Type:   bimg.WEBP,
		Height: newHeight,
		Width:  newWidth,
	})
	if err != nil {
		return nil, "", err
	}

	contentType := "image/webp"

	return processedImage, contentType, nil
}

type CreateImageArgs struct {
	BucketName string
	FileName   string
	File       io.Reader
	Size       int64
}

func (service *imageMediaService) UploadImage(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) (*entities.ImageEntity, error) {
	err := service.blobStorageService.UploadObject(ctx, bucketName, fileName, file, size)
	if err != nil {
		return nil, err
	}

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	image := bimg.NewImage(imageBytes)

	metadata, err := image.Metadata()
	if err != nil {
		return nil, err
	}

	imageEntity := entities.ImageEntity{
		ID:         uuid.New(),
		BucketName: bucketName,
		ObjectID:   fileName,
		Width:      int32(metadata.Size.Width),
		Height:     int32(metadata.Size.Height),
		Size:       int32(size),
	}

	return &imageEntity, nil
}

func calculateDimensions(width, height, maxSize int) (int, int) {
	aspectRatio := float64(width) / float64(height)

	var newWidth, newHeight int

	if width > height {
		newWidth = maxSize
		newHeight = int(math.Round(float64(newWidth) / aspectRatio))
	} else {
		newHeight = maxSize
		newWidth = int(math.Round(float64(newHeight) * aspectRatio))
	}
	return newWidth, newHeight
}

var (
	RenditionAvatar = "avatar"
	RenditionSmall  = "small"
	RenditionMedium = "medium"
	RenditionLarge  = "large"
	RenditionXLarge = "xlarge"
)

var renditions = map[string]int{
	"avatar": 120,
	"small":  720,
	"medium": 1080,
	"large":  2160,
	"xlarge": 4320,
}

func getRendition(r string) int {
	size := renditions[r]
	if size == 0 {
		return renditions["large"]
	}

	return size
}
