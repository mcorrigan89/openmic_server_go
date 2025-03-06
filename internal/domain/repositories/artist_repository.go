package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type ArtistRepository interface {
	GetArtistByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.ArtistEntity, error)
	GetArtistsByTitle(ctx context.Context, querier models.Querier, title string) ([]*entities.ArtistEntity, error)
	GetAllArtists(ctx context.Context, querier models.Querier) ([]*entities.ArtistEntity, error)
	CreateArtist(ctx context.Context, querier models.Querier, artist *entities.ArtistEntity) (*entities.ArtistEntity, error)
}
