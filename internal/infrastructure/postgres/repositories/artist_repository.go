package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type postgresArtistRepository struct {
}

func NewPostgresArtistRepository() *postgresArtistRepository {
	return &postgresArtistRepository{}
}

func (repo *postgresArtistRepository) GetArtistByID(ctx context.Context, querier models.Querier, artistID uuid.UUID) (*entities.ArtistEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetArtistByID(ctx, artistID)
	if err != nil {
		return nil, err
	}

	return entities.NewArtistEntity(row.Artist), nil
}

func (repo *postgresArtistRepository) GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	rows, err := querier.GetArtistsByTitle(ctx, models.GetArtistsByTitleParams{
		Title:         title,
		MinSimilarity: 0.2,
	})
	if err != nil {
		return nil, err
	}

	var artistEntities []*entities.ArtistEntity
	for _, row := range rows {
		artistEntities = append(artistEntities, entities.NewArtistEntity(row.Artist))
	}

	return artistEntities, nil
}

func (repo *postgresArtistRepository) GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	rows, err := querier.GetAllArtists(ctx)
	if err != nil {
		return nil, err
	}

	var artistEntities []*entities.ArtistEntity
	for _, row := range rows {
		artistEntities = append(artistEntities, entities.NewArtistEntity(row.Artist))
	}

	return artistEntities, nil
}

func (repo *postgresArtistRepository) CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateArtist(ctx, models.CreateArtistParams{
		ID:             artist.ID,
		ArtistTitle:    artist.Title,
		ArtistSubtitle: artist.SubTitle,
		Bio:            artist.Bio,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewArtistEntity(row), nil
}

func (repo *postgresArtistRepository) UpdateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.UpdateArtist(ctx, models.UpdateArtistParams{
		ID:             artist.ID,
		ArtistTitle:    artist.Title,
		ArtistSubtitle: artist.SubTitle,
		Bio:            artist.Bio,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewArtistEntity(row), nil
}

func (repo *postgresArtistRepository) DeleteArtist(ctx context.Context, querier models.Querier, artistID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.DeleteArtist(ctx, artistID)
	if err != nil {
		return err
	}

	return nil
}
