package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	mocks "github.com/mcorrigan89/openmic/internal/mocks/repositories"
	"github.com/stretchr/testify/assert"
)

// Test cases
func TestGetUserByID(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRefLinkRepo := new(mocks.MockReferenceLinkRepository)
	mockImageRepo := new(mocks.MockImageRepository)
	userService := NewUserService(mockUserRepo, mockRefLinkRepo, mockImageRepo)

	ctx := context.Background()
	querier := &models.Queries{}
	userID := uuid.New()
	notUserID := uuid.New()

	t.Run("successful get user by id", func(t *testing.T) {
		expectedUser := &entities.UserEntity{
			ID:    userID,
			Email: "test@example.com",
		}

		mockUserRepo.On("GetUserByID", ctx, querier, userID).Return(expectedUser, nil)

		user, err := userService.GetUserByID(ctx, querier, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)

	})

	t.Run("user by id not found", func(t *testing.T) {

		mockUserRepo.On("GetUserByID", ctx, querier, notUserID).Return(nil, entities.ErrUserNotFound)

		user, err := userService.GetUserByID(ctx, querier, notUserID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, entities.ErrUserNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRefLinkRepo := new(mocks.MockReferenceLinkRepository)
	mockImageRepo := new(mocks.MockImageRepository)
	userService := NewUserService(mockUserRepo, mockRefLinkRepo, mockImageRepo)

	ctx := context.Background()
	querier := &models.Queries{}
	userID := uuid.New()
	email := "test@example.com"
	notEmail := "nope@example.com"
	t.Run("get user by email", func(t *testing.T) {

		expectedUser := &entities.UserEntity{
			ID:    userID,
			Email: email,
		}

		mockUserRepo.On("GetUserByEmail", ctx, querier, email).Return(expectedUser, nil)

		user, err := userService.GetUserByEmail(ctx, querier, email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user by id not found", func(t *testing.T) {

		mockUserRepo.On("GetUserByEmail", ctx, querier, notEmail).Return(nil, entities.ErrUserNotFound)

		user, err := userService.GetUserByEmail(ctx, querier, notEmail)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, entities.ErrUserNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})

}

func TestGetUserByHandle(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRefLinkRepo := new(mocks.MockReferenceLinkRepository)
	mockImageRepo := new(mocks.MockImageRepository)
	userService := NewUserService(mockUserRepo, mockRefLinkRepo, mockImageRepo)

	ctx := context.Background()
	querier := &models.Queries{}
	userID := uuid.New()

	t.Run("get user by handle", func(t *testing.T) {

		handle := "abc123"
		expectedUser := &entities.UserEntity{
			ID:     userID,
			Handle: handle,
		}

		mockUserRepo.On("GetUserByHandle", ctx, querier, handle).Return(expectedUser, nil)

		user, err := userService.GetUserByHandle(ctx, querier, handle)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user by handle not found", func(t *testing.T) {

		handle := "abc456"
		mockUserRepo.On("GetUserByHandle", ctx, querier, handle).Return(nil, entities.ErrUserNotFound)

		user, err := userService.GetUserByHandle(ctx, querier, handle)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, entities.ErrUserNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})

}

func TestGetUserByToken(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRefLinkRepo := new(mocks.MockReferenceLinkRepository)
	mockImageRepo := new(mocks.MockImageRepository)
	userService := NewUserService(mockUserRepo, mockRefLinkRepo, mockImageRepo)

	ctx := context.Background()
	querier := &models.Queries{}
	userID := uuid.New()

	t.Run("get user by token", func(t *testing.T) {

		token := "zzzzzzz"
		expectedUserContext := &entities.UserContextEntity{
			UserID: userID,
			User: &entities.UserEntity{
				ID: userID,
			},
			SessionToken: token,
		}

		mockUserRepo.On("GetUserContextBySessionToken", ctx, querier, token).Return(expectedUserContext, nil)

		userContext, err := userService.GetUserContextBySessionToken(ctx, querier, token)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserContext, userContext)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user by handle not found", func(t *testing.T) {

		token := "yyyyyyy"
		mockUserRepo.On("GetUserContextBySessionToken", ctx, querier, token).Return(nil, entities.ErrUserNotFound)

		userContext, err := userService.GetUserContextBySessionToken(ctx, querier, token)

		assert.Error(t, err)
		assert.Nil(t, userContext)
		assert.Equal(t, entities.ErrUserNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})

}
