package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type postgresReferenceLinkRepository struct {
}

func NewPostgresReferenceLinkRepository() *postgresReferenceLinkRepository {
	return &postgresReferenceLinkRepository{}
}

func (repo *postgresReferenceLinkRepository) CreateReferenceLink(ctx context.Context, querier models.Querier, emailLink *entities.ReferenceLinkEntity) (*entities.ReferenceLinkEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateReferenceLink(ctx, models.CreateReferenceLinkParams{
		ID:        emailLink.ID,
		LinkID:    emailLink.LinkID,
		LinkType:  emailLink.Type,
		Token:     emailLink.Token,
		ExpiresAt: emailLink.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewReferenceLinkEntity(row), nil
}

func (repo *postgresReferenceLinkRepository) GetReferenceLinkByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.ReferenceLinkEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetReferenceLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return entities.NewReferenceLinkEntity(row.ReferenceLink), nil
}

func (repo *postgresReferenceLinkRepository) GetReferenceLinkByToken(ctx context.Context, querier models.Querier, token string) (*entities.ReferenceLinkEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetReferenceLinkByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return entities.NewReferenceLinkEntity(row.ReferenceLink), nil
}

func (repo *postgresReferenceLinkRepository) DeleteReferenceLink(ctx context.Context, querier models.Querier, emailLink *entities.ReferenceLinkEntity) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	_, err := querier.DeleteReferenceLink(ctx, emailLink.ID)
	if err != nil {
		return err
	}
	return nil
}
