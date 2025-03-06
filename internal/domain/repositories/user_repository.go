package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, querier models.Querier, userId uuid.UUID) (*entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error)
	GetUserByHandle(ctx context.Context, querier models.Querier, userHandle string) (*entities.UserEntity, error)
	GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error)
	CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity, token string, expiresAt time.Time) (*entities.UserContextEntity, error)
	SetAvatarImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity, user *entities.UserEntity) (*entities.UserEntity, error)
}
