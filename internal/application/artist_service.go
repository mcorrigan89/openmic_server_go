package application

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/services"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"

	"github.com/rs/zerolog"
)

type ArtistApplicationService interface {
	GetArtistByID(ctx context.Context, query queries.ArtistByIDQuery) (*entities.ArtistEntity, error)
	GetArtistsByTitle(ctx context.Context, query queries.ArtistsByTitleQuery) ([]*entities.ArtistEntity, error)
	GetAllArtists(ctx context.Context) ([]*entities.ArtistEntity, error)
	CreateArtist(ctx context.Context, cmd commands.CreateNewArtistCommand) (*entities.ArtistEntity, error)
}

type artistApplicationService struct {
	config        *common.Config
	wg            *sync.WaitGroup
	logger        *zerolog.Logger
	db            *pgxpool.Pool
	queries       models.Querier
	artistService services.ArtistService
}

func NewArtistApplicationService(db *pgxpool.Pool, wg *sync.WaitGroup, cfg *common.Config, logger *zerolog.Logger, artistService services.ArtistService) *artistApplicationService {
	dbQueries := models.New(db)
	return &artistApplicationService{
		db:            db,
		config:        cfg,
		wg:            wg,
		logger:        logger,
		queries:       dbQueries,
		artistService: artistService,
	}
}

func (app *artistApplicationService) GetArtistByID(ctx context.Context, query queries.ArtistByIDQuery) (*entities.ArtistEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting artist by ID")

	artist, err := app.artistService.GetArtistByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get artist by ID")
		return nil, err
	}

	return artist, nil
}

func (app *artistApplicationService) GetArtistsByTitle(ctx context.Context, query queries.ArtistsByTitleQuery) ([]*entities.ArtistEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting artist by ID")

	artists, err := app.artistService.GetArtistsByTitle(ctx, app.queries, query.Title)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get artist by ID")
		return nil, err
	}

	return artists, nil
}

func (app *artistApplicationService) GetAllArtists(ctx context.Context) ([]*entities.ArtistEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting all artists")

	artists, err := app.artistService.GetAllArtists(ctx, app.queries)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get all artists")
		return nil, err
	}

	return artists, nil
}

func (app *artistApplicationService) CreateArtist(ctx context.Context, cmd commands.CreateNewArtistCommand) (*entities.ArtistEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Creating new artist")

	artistEntity := cmd.ToDomain()

	artist, err := app.artistService.CreateArtist(ctx, app.queries, artistEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new artist")
		return nil, err
	}

	return artist, nil
}
