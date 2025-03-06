package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/mock"
)

type MockImageRepository struct {
	mock.Mock
}

func (m *MockImageRepository) GetImageByID(ctx context.Context, querier models.Querier, imageID uuid.UUID) (*entities.ImageEntity, error) {
	args := m.Called(ctx, querier, imageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ImageEntity), args.Error(1)
}

func (m *MockImageRepository) CreateImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity) (*entities.ImageEntity, error) {
	args := m.Called(ctx, querier, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ImageEntity), args.Error(1)
}
