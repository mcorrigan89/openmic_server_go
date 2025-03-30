package dto

import (
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type ArtistDto struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	SubTitle *string   `json:"sub_title"`
	Bio      *string   `json:"bio"`
}

func NewArtistDtoFromEntity(entity *entities.ArtistEntity) *ArtistDto {
	return &ArtistDto{
		ID:       entity.ID,
		Title:    entity.Title,
		SubTitle: entity.SubTitle,
		Bio:      entity.Bio,
	}
}

type GetArtistByIDResponse struct {
	Body *ArtistDto `json:"body"`
}

type GetArtistsByTitleResponse struct {
	Body []*ArtistDto `json:"body"`
}

type GetAllArtistsResponse struct {
	Body []*ArtistDto `json:"body"`
}

type CreateArtistRequest struct {
	Body struct {
		Title    string  `json:"title"`
		SubTitle *string `json:"sub_title"`
		Bio      *string `json:"bio"`
	}
}

type CreateArtistResponse struct {
	Body *ArtistDto `json:"body"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Title    string  `json:"title"`
		SubTitle *string `json:"sub_title"`
		Bio      *string `json:"bio"`
	}
}

type UpdateArtistResponse struct {
	Body *ArtistDto `json:"body"`
}

type DeleteArtistRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteArtistResponse struct {
	Body *ArtistDto `json:"body"`
}
