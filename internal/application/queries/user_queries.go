package queries

import (
	"github.com/google/uuid"
)

type UserByIDQuery struct {
	ID uuid.UUID
}

type UserByEmailQuery struct {
	Email string `validate:"required,email"`
}

type UserByHandleQuery struct {
	Handle string
}

type UserBySessionTokenQuery struct {
	SessionToken string
}
