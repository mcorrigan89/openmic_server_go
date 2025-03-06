package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/assert"
)

func TestNewUserContextEntity(t *testing.T) {

	userId := uuid.New()
	givenName := "John"
	familyName := "Doe"
	email := "johndoe@gmail.com"
	handle := "johndoe"

	userModel := models.User{
		ID:            userId,
		GivenName:     &givenName,
		FamilyName:    &familyName,
		Email:         email,
		EmailVerified: true,
		UserHandle:    handle,
		Claimed:       true,
		AvatarID:      nil,
	}

	userEntity := NewUserEntity(userModel, nil)

	t.Run("create entity from model", func(t *testing.T) {

		sessionID := uuid.New()
		expireTime := time.Now().Add(time.Hour)

		userSessionModel := models.UserSession{
			ID:          sessionID,
			UserID:      userId,
			Token:       "test-token",
			ExpiresAt:   expireTime,
			UserExpired: false,
		}

		userContextEntity := NewUserContextEntity(userEntity, userSessionModel)

		assert.Equal(t, userContextEntity.SessionToken, "test-token")
		assert.Equal(t, userContextEntity.UserID, userId)
		assert.Equal(t, userContextEntity.User, userEntity)
		assert.Equal(t, userContextEntity.userExpired, false)
		assert.Equal(t, userContextEntity.expiresAt, expireTime)
		assert.Equal(t, userContextEntity.ExpiresAt(), expireTime)
		assert.Equal(t, userContextEntity.IsExpired(), false)
	})

	t.Run("is expired by time", func(t *testing.T) {

		sessionID := uuid.New()
		expireTime := time.Now().Add(-time.Hour)

		userSessionModel := models.UserSession{
			ID:          sessionID,
			UserID:      userId,
			Token:       "test-token",
			ExpiresAt:   expireTime,
			UserExpired: false,
		}

		userContextEntity := NewUserContextEntity(userEntity, userSessionModel)

		assert.Equal(t, userContextEntity.SessionToken, "test-token")
		assert.Equal(t, userContextEntity.UserID, userId)
		assert.Equal(t, userContextEntity.User, userEntity)
		assert.Equal(t, userContextEntity.userExpired, false)
		assert.Equal(t, userContextEntity.expiresAt, expireTime)
		assert.Equal(t, userContextEntity.ExpiresAt(), expireTime)
		assert.Equal(t, userContextEntity.IsExpired(), true)
	})

	t.Run("is expired by user", func(t *testing.T) {

		sessionID := uuid.New()
		expireTime := time.Now().Add(time.Hour)

		userSessionModel := models.UserSession{
			ID:          sessionID,
			UserID:      userId,
			Token:       "test-token",
			ExpiresAt:   expireTime,
			UserExpired: true,
		}

		userContextEntity := NewUserContextEntity(userEntity, userSessionModel)

		assert.Equal(t, userContextEntity.SessionToken, "test-token")
		assert.Equal(t, userContextEntity.UserID, userId)
		assert.Equal(t, userContextEntity.User, userEntity)
		assert.Equal(t, userContextEntity.userExpired, true)
		assert.Equal(t, userContextEntity.expiresAt, expireTime)
		assert.Equal(t, userContextEntity.ExpiresAt(), expireTime)
		assert.Equal(t, userContextEntity.IsExpired(), true)
	})
}
