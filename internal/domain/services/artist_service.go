package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ArtistService interface {
	GetArtistByID(ctx context.Context, querier models.Querier, artistID uuid.UUID) (*entities.ArtistEntity, error)
	GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error)
	GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error)
	CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error)
}

type artistService struct {
	artistRepo repositories.ArtistRepository
}

func NewArtistService(artistRepo repositories.ArtistRepository) *artistService {
	return &artistService{artistRepo: artistRepo}
}

func (s *artistService) GetArtistByID(ctx context.Context, querier models.Querier, artistID uuid.UUID) (*entities.ArtistEntity, error) {
	artist, err := s.artistRepo.GetArtistByID(ctx, querier, artistID)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *artistService) GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error) {
	artists, err := s.artistRepo.GetArtistsByTitle(ctx, querier, title)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (s *artistService) GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error) {
	artists, err := s.artistRepo.GetAllArtists(ctx, querier)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (s *artistService) CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error) {
	createdArtist, err := s.artistRepo.CreateArtist(ctx, querier, artist)
	if err != nil {
		return nil, err
	}

	return createdArtist, nil
}
