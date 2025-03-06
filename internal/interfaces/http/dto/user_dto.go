package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type UserDto struct {
	ID         uuid.UUID `json:"id"`
	GivenName  *string   `json:"given_name"`
	FamilyName *string   `json:"family_name"`
	Email      string    `json:"email"`
}

func NewUserDtoFromEntity(entity *entities.UserEntity) *UserDto {
	return &UserDto{
		ID:         entity.ID,
		GivenName:  entity.GivenName,
		FamilyName: entity.FamilyName,
		Email:      entity.Email,
	}
}

type SessionDto struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetUserByIDResponse struct {
	Body *UserDto `json:"body"`
}

type GetUserByEmailResponse struct {
	Body *UserDto `json:"body"`
}

type CreateUserRequest struct {
	Body struct {
		Email      string  `json:"email"`
		GivenName  *string `json:"given_name"`
		FamilyName *string `json:"family_name"`
	}
}

type CreateUserResponse struct {
	Body struct {
		User       *UserDto    `json:"user"`
		SessionDto *SessionDto `json:"session"`
	} `json:"body"`
}

type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Email      string  `json:"email"`
		GivenName  *string `json:"given_name"`
		FamilyName *string `json:"family_name"`
		Handle     string  `json:"handle"`
	}
}

type UpdateUserResponse struct {
	Body struct {
		User *UserDto `json:"user"`
	} `json:"body"`
}
