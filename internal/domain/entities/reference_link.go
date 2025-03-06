package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

var (
	ErrLinkNotFound = errors.New("link not found")
	ErrLinkExpired  = errors.New("link expired")
	ErrLinkInvalid  = errors.New("link type invalid")
)

var (
	RefLinkTypeLogin  = "login"
	RefLinkTypeInvite = "invite"
)

type ReferenceLinkEntity struct {
	ID        uuid.UUID
	LinkID    uuid.UUID
	Token     string
	Type      string
	ExpiresAt time.Time
}

func NewReferenceLinkEntity(refLinkModel models.ReferenceLink) *ReferenceLinkEntity {
	return &ReferenceLinkEntity{
		ID:        refLinkModel.ID,
		LinkID:    refLinkModel.LinkID,
		Token:     refLinkModel.Token,
		Type:      refLinkModel.LinkType,
		ExpiresAt: refLinkModel.ExpiresAt,
	}
}

func (el *ReferenceLinkEntity) IsExpired() bool {
	return el.ExpiresAt.Before(time.Now())
}
