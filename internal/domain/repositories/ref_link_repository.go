package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ReferenceLinkRepository interface {
	CreateReferenceLink(ctx context.Context, querier models.Querier, refLink *entities.ReferenceLinkEntity) (*entities.ReferenceLinkEntity, error)
	GetReferenceLinkByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.ReferenceLinkEntity, error)
	GetReferenceLinkByToken(ctx context.Context, querier models.Querier, token string) (*entities.ReferenceLinkEntity, error)
	DeleteReferenceLink(ctx context.Context, querier models.Querier, refLink *entities.ReferenceLinkEntity) error
}
