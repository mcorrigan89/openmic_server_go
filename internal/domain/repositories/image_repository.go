package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ImageRepository interface {
	GetImageByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.ImageEntity, error)
	CreateImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity) (*entities.ImageEntity, error)
}
