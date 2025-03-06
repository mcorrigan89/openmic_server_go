package queries

import (
	"github.com/google/uuid"
)

type ImageByIDQuery struct {
	ID uuid.UUID
}

type ImageDataByIDQuery struct {
	ID        uuid.UUID
	Rendition string
}

type CollectionByIDQuery struct {
	ID uuid.UUID
}

type CollectionByOwnerIDQuery struct {
	OwnerID uuid.UUID
}

type CollectionByOwnerTokenQuery struct {
	Token string
}
