package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/mock"
)

type MockReferenceLinkRepository struct {
	mock.Mock
}

func (m *MockReferenceLinkRepository) CreateReferenceLink(ctx context.Context, querier models.Querier, link *entities.ReferenceLinkEntity) (*entities.ReferenceLinkEntity, error) {
	args := m.Called(ctx, querier, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ReferenceLinkEntity), args.Error(1)
}

func (m *MockReferenceLinkRepository) GetReferenceLinkByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.ReferenceLinkEntity, error) {
	args := m.Called(ctx, querier, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ReferenceLinkEntity), args.Error(1)
}

func (m *MockReferenceLinkRepository) GetReferenceLinkByToken(ctx context.Context, querier models.Querier, token string) (*entities.ReferenceLinkEntity, error) {
	args := m.Called(ctx, querier, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ReferenceLinkEntity), args.Error(1)
}

func (m *MockReferenceLinkRepository) DeleteReferenceLink(ctx context.Context, querier models.Querier, link *entities.ReferenceLinkEntity) error {
	args := m.Called(ctx, querier, link)
	return args.Error(0)
}
