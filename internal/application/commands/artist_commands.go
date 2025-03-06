package commands

import (
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type CreateNewArtistCommand struct {
	Title    string
	SubTitle *string
	Bio      *string
}

func (c *CreateNewArtistCommand) ToDomain() *entities.ArtistEntity {
	return &entities.ArtistEntity{
		ID:       uuid.New(),
		Title:    c.Title,
		SubTitle: c.SubTitle,
		Bio:      c.Bio,
	}
}
