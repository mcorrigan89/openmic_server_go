package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type postgresUserRepository struct {
}

func NewPostgresUserRepository() *postgresUserRepository {
	return &postgresUserRepository{}
}

func (repo *postgresUserRepository) GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row.User)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row.User, avatarImageEntity), nil
}

func (repo *postgresUserRepository) GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row.User)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row.User, avatarImageEntity), nil
}

func (repo *postgresUserRepository) GetUserByHandle(ctx context.Context, querier models.Querier, userHandle string) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserByHandle(ctx, userHandle)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row.User)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row.User, avatarImageEntity), nil
}

func (repo *postgresUserRepository) GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row.User)
	if err != nil {
		return nil, err
	}

	userEntity := entities.NewUserEntity(row.User, avatarImageEntity)

	return entities.NewUserContextEntity(userEntity, row.UserSession), nil
}

func (repo *postgresUserRepository) CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateUser(ctx, models.CreateUserParams{
		ID:            user.ID,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Email:         user.Email,
		EmailVerified: false,
		Claimed:       user.Claimed,
		UserHandle:    user.Handle,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, entities.ErrEmailInUse
			case "users_user_handle_key":
				return nil, entities.ErrUserHandleInUse
			default:
				return nil, fmt.Errorf("unique constraint violation: %s", pgErr.ConstraintName)
			}
		}
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row, avatarImageEntity), nil
}

func (repo *postgresUserRepository) UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.UpdateUser(ctx, models.UpdateUserParams{
		ID:            user.ID,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Claimed:       user.Claimed,
		UserHandle:    user.Handle,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, entities.ErrEmailInUse
			case "users_user_handle_key":
				return nil, entities.ErrUserHandleInUse
			default:
				return nil, fmt.Errorf("unique constraint violation: %s", pgErr.ConstraintName)
			}
		}
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row, avatarImageEntity), nil
}

func (repo *postgresUserRepository) CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity, token string, expiresAt time.Time) (*entities.UserContextEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateUserSession(ctx, models.CreateUserSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewUserContextEntity(user, row), nil
}

func (repo *postgresUserRepository) SetAvatarImage(ctx context.Context, querier models.Querier, image *entities.ImageEntity, user *entities.UserEntity) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.SetAvatarImage(ctx, models.SetAvatarImageParams{
		UserID:  user.ID,
		ImageID: &image.ID,
	})
	if err != nil {
		return nil, err
	}

	avatarImageEntity, err := repo.getAvatarImageEntity(ctx, querier, &row)
	if err != nil {
		return nil, err
	}

	return entities.NewUserEntity(row, avatarImageEntity), nil
}

func (repo *postgresUserRepository) getAvatarImageEntity(ctx context.Context, querier models.Querier, user *models.User) (*entities.ImageEntity, error) {
	if user.AvatarID != nil {
		avatarRow, err := querier.GetImageByID(ctx, *user.AvatarID)
		if err != nil {
			return nil, err
		}

		return entities.NewImageEntity(avatarRow.Image), nil
	}

	return nil, nil
}
