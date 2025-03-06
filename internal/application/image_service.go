package application

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/external"
	"github.com/mcorrigan89/openmic/internal/domain/services"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"

	"github.com/rs/zerolog"
)

type ImageApplicationService interface {
	GetImageByID(ctx context.Context, query queries.ImageByIDQuery) (*entities.ImageEntity, error)
	GetImageDataByID(ctx context.Context, query queries.ImageDataByIDQuery) ([]byte, string, error)
	UploadAvatarImage(ctx context.Context, cmd commands.CreateNewAvatarImageCommand) (*entities.ImageEntity, error)
}

type imageApplicationService struct {
	config            *common.Config
	wg                *sync.WaitGroup
	logger            *zerolog.Logger
	db                *pgxpool.Pool
	queries           models.Querier
	imageService      services.ImageService
	userService       services.UserService
	imageMediaService external.ImageMediaService
}

func NewImageApplicationService(db *pgxpool.Pool, wg *sync.WaitGroup, cfg *common.Config, logger *zerolog.Logger, imageService services.ImageService, userService services.UserService, imageMediaService external.ImageMediaService) *imageApplicationService {
	dbQueries := models.New(db)
	return &imageApplicationService{
		db:                db,
		config:            cfg,
		wg:                wg,
		logger:            logger,
		queries:           dbQueries,
		imageService:      imageService,
		userService:       userService,
		imageMediaService: imageMediaService,
	}
}

func (app *imageApplicationService) GetImageByID(ctx context.Context, query queries.ImageByIDQuery) (*entities.ImageEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting image by ID")

	image, err := app.imageService.GetImageByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get image by ID")
		return nil, err
	}

	return image, nil
}

func (app *imageApplicationService) GetImageDataByID(ctx context.Context, query queries.ImageDataByIDQuery) ([]byte, string, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting image data by ID")

	imageEntity, err := app.imageService.GetImageByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get image by ID")
		return nil, "", err
	}

	imageData, contentType, err := app.imageMediaService.GetImageDataByID(ctx, imageEntity, query.Rendition)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get image data by ID")
		return nil, "", err
	}

	return imageData, contentType, nil
}

func (app *imageApplicationService) UploadAvatarImage(ctx context.Context, cmd commands.CreateNewAvatarImageCommand) (*entities.ImageEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Uploading image")
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Uploading image data")
	imageData, err := app.imageMediaService.UploadImage(ctx, cmd.BucketID, cmd.ObjectID, cmd.File, cmd.Size)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to upload image")
		return nil, err
	}

	app.logger.Info().Ctx(ctx).Msg("Creating image in database")
	image, err := app.imageService.CreateImage(ctx, qtx, imageData)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create image")
		return nil, err
	}

	app.logger.Info().Ctx(ctx).Msg("Getting user by ID")
	user, err := app.userService.GetUserByID(ctx, app.queries, cmd.UserID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by ID")
		return nil, err
	}

	app.logger.Info().Ctx(ctx).Msg("Setting image as avatar on user")
	_, err = app.userService.SetAvatarImage(ctx, qtx, image, user)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to set avatar image")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return image, nil
}
