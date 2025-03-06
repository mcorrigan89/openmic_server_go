package entities

import (
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ArtistEntity struct {
	ID       uuid.UUID
	Title    string
	SubTitle *string
	Bio      *string
}

func NewArtistEntity(artistModel models.Artist) *ArtistEntity {
	return &ArtistEntity{
		ID:       artistModel.ID,
		Title:    artistModel.ArtistTitle,
		SubTitle: artistModel.ArtistSubtitle,
		Bio:      artistModel.Bio,
	}
}
