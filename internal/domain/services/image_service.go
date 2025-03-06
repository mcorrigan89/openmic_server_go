package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ImageService interface {
	GetImageByID(ctx context.Context, querier models.Querier, imageID uuid.UUID) (*entities.ImageEntity, error)
	CreateImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity) (*entities.ImageEntity, error)
}

type imageService struct {
	imageRepo repositories.ImageRepository
}

func NewImageService(imageRepo repositories.ImageRepository) *imageService {
	return &imageService{imageRepo: imageRepo}
}

func (s *imageService) GetImageByID(ctx context.Context, querier models.Querier, imageID uuid.UUID) (*entities.ImageEntity, error) {
	imageEntity, err := s.imageRepo.GetImageByID(ctx, querier, imageID)
	if err != nil {
		return nil, err
	}

	return imageEntity, nil
}

func (s *imageService) CreateImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity) (*entities.ImageEntity, error) {
	imageEntity, err := s.imageRepo.CreateImage(ctx, querier, image)
	if err != nil {
		return nil, err
	}

	return imageEntity, nil
}
