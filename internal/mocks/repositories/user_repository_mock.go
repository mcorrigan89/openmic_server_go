package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

// Implement complete UserRepository interface
func (m *MockUserRepository) GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetUserByHandle(ctx context.Context, querier models.Querier, userHandle string) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, userHandle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error) {
	args := m.Called(ctx, querier, sessionToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserContextEntity), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func (m *MockUserRepository) CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity, token string, expiresAt time.Time) (*entities.UserContextEntity, error) {
	args := m.Called(ctx, querier, user, token, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserContextEntity), args.Error(1)
}

func (m *MockUserRepository) SetAvatarImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity, user *entities.UserEntity) (*entities.UserEntity, error) {
	args := m.Called(ctx, querier, image, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}
