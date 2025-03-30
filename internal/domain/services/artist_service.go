package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/rs/zerolog"
)

type ArtistService interface {
	GetArtistByID(ctx context.Context, querier models.Querier, artistID uuid.UUID) (*entities.ArtistEntity, error)
	GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error)
	GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error)
	CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error)
	UpdateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error)
	DeleteArtist(ctx context.Context, querier models.Querier, artistID uuid.UUID) error
}

type artistService struct {
	logger     *zerolog.Logger
	artistRepo repositories.ArtistRepository
}

func NewArtistService(logger *zerolog.Logger, artistRepo repositories.ArtistRepository) *artistService {
	return &artistService{logger: logger, artistRepo: artistRepo}
}

func (s *artistService) GetArtistByID(ctx context.Context, querier models.Querier, artistID uuid.UUID) (*entities.ArtistEntity, error) {
	artist, err := s.artistRepo.GetArtistByID(ctx, querier, artistID)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to get artist by ID")
		return nil, err
	}

	return artist, nil
}

func (s *artistService) GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error) {
	artists, err := s.artistRepo.GetArtistsByTitle(ctx, querier, title)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to get artist by title")
		return nil, err
	}

	return artists, nil
}

func (s *artistService) GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error) {
	artists, err := s.artistRepo.GetAllArtists(ctx, querier)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to get all artists")
		return nil, err
	}

	return artists, nil
}

func (s *artistService) CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error) {
	createdArtist, err := s.artistRepo.CreateArtist(ctx, querier, artist)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to create artist")
		return nil, err
	}

	return createdArtist, nil
}

func (s *artistService) UpdateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error) {
	updatedArtist, err := s.artistRepo.UpdateArtist(ctx, querier, artist)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to update artist")
		return nil, err
	}

	return updatedArtist, nil
}

func (s *artistService) DeleteArtist(ctx context.Context, querier models.Querier, artistID uuid.UUID) error {
	err := s.artistRepo.DeleteArtist(ctx, querier, artistID)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Failed to delete artist")
		return err
	}

	return nil
}
