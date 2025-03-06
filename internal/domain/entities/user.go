package entities

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailInUse      = errors.New("email in use")
	ErrUserHandleInUse = errors.New("handle in use")
	ErrUserClaimed     = errors.New("user already claimed")
)

type UserEntity struct {
	ID            uuid.UUID
	GivenName     *string
	FamilyName    *string
	Email         string
	EmailVerified bool
	Claimed       bool
	Handle        string
	Avatar        *ImageEntity
}

func NewUserEntity(userModel models.User, imageEntity *ImageEntity) *UserEntity {
	return &UserEntity{
		ID:            userModel.ID,
		GivenName:     userModel.GivenName,
		FamilyName:    userModel.FamilyName,
		Email:         userModel.Email,
		EmailVerified: userModel.EmailVerified,
		Claimed:       userModel.Claimed,
		Handle:        userModel.UserHandle,
		Avatar:        imageEntity,
	}
}

func (u *UserEntity) FullName() *string {
	if u.GivenName != nil && u.FamilyName != nil {
		fullName := fmt.Sprintf("%s %s", *u.GivenName, *u.FamilyName)
		return &fullName
	}
	if u.GivenName != nil {
		return u.GivenName
	}

	if u.FamilyName != nil {
		return u.FamilyName
	}
	return nil
}
