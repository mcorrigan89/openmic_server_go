package commands

import (
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/rs/xid"
)

type CreateNewUserCommand struct {
	Email      string  `json:"email" validate:"required,email"`
	GivenName  *string `json:"firstName" validate:"-"`
	FamilyName *string `json:"lastName" validate:"-"`
}

func (cmd *CreateNewUserCommand) ToDomain() *entities.UserEntity {
	return &entities.UserEntity{
		ID:         uuid.New(),
		Email:      cmd.Email,
		GivenName:  cmd.GivenName,
		FamilyName: cmd.FamilyName,
		Claimed:    true,
		Handle:     xid.New().String(),
	}
}

type UpdateUserCommand struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	GivenName  *string   `json:"firstName" validate:"-"`
	FamilyName *string   `json:"lastName" validate:"-"`
	Handle     string    `json:"handle" validate:"required"`
}

func (cmd *UpdateUserCommand) ToDomain() *entities.UserEntity {
	return &entities.UserEntity{
		ID:         cmd.ID,
		Email:      cmd.Email,
		GivenName:  cmd.GivenName,
		FamilyName: cmd.FamilyName,
		Claimed:    true,
		Handle:     cmd.Handle,
	}
}

type RequestEmailLoginCommand struct {
	Email string `json:"email" validate:"required,email"`
}

func (cmd *RequestEmailLoginCommand) ToDomain() string {
	return cmd.Email
}

type LoginWithReferenceLinkCommand struct {
	ReferenceLinkToken string `json:"token" validate:"required"`
}

type InviteUserCommand struct {
	Email      string  `json:"email" validate:"required,email"`
	GivenName  *string `json:"firstName" validate:"-"`
	FamilyName *string `json:"lastName" validate:"-"`
}

func (cmd *InviteUserCommand) ToDomain() *entities.UserEntity {
	return &entities.UserEntity{
		ID:         uuid.New(),
		Email:      cmd.Email,
		GivenName:  cmd.GivenName,
		FamilyName: cmd.FamilyName,
		Claimed:    false,
		Handle:     xid.New().String(),
	}
}

type AcceptInviteReferenceLinkCommand struct {
	ReferenceLinkToken string `json:"token" validate:"required"`
}
