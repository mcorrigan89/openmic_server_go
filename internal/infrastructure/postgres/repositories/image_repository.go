package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type postgresImageRepository struct {
}

func NewPostgresImageRepository() *postgresImageRepository {
	return &postgresImageRepository{}
}

func (repo *postgresImageRepository) GetImageByID(ctx context.Context, querier models.Querier, imageID uuid.UUID) (*entities.ImageEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetImageByID(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return entities.NewImageEntity(row.Image), nil
}

func (repo *postgresImageRepository) CreateImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity) (*entities.ImageEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateImage(ctx, models.CreateImageParams{
		ID:         image.ID,
		BucketName: image.BucketName,
		ObjectID:   image.ObjectID,
		Width:      image.Width,
		Height:     image.Height,
		FileSize:   image.Size,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewImageEntity(row), nil
}
