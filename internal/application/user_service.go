package application

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/services"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"

	"github.com/rs/zerolog"
)

type UserApplicationService interface {
	GetUserByID(ctx context.Context, query queries.UserByIDQuery) (*entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, query queries.UserByEmailQuery) (*entities.UserEntity, error)
	GetUserByHandle(ctx context.Context, query queries.UserByHandleQuery) (*entities.UserEntity, error)
	GetUserBySessionToken(ctx context.Context, query queries.UserBySessionTokenQuery) (*entities.UserEntity, error)
	CreateUser(ctx context.Context, cmd commands.CreateNewUserCommand) (*entities.UserContextEntity, error)
	UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entities.UserEntity, error)
	RequestEmailLogin(ctx context.Context, cmd commands.RequestEmailLoginCommand) error
	LoginWithReferenceLink(ctx context.Context, cmd commands.LoginWithReferenceLinkCommand) (*entities.UserContextEntity, error)
	InviteUser(ctx context.Context, cmd commands.InviteUserCommand) error
	AcceptInviteReferenceLink(ctx context.Context, cmd commands.AcceptInviteReferenceLinkCommand) (*entities.UserContextEntity, error)
}

type userApplicationService struct {
	config               *common.Config
	wg                   *sync.WaitGroup
	logger               *zerolog.Logger
	db                   *pgxpool.Pool
	queries              models.Querier
	userService          services.UserService
	emailService         services.EmailService
	emailTemplateService services.EmailTemplateService
}

func NewUserApplicationService(db *pgxpool.Pool, wg *sync.WaitGroup, cfg *common.Config, logger *zerolog.Logger, userService services.UserService, emailService services.EmailService, emailTemplateService services.EmailTemplateService) *userApplicationService {
	dbQueries := models.New(db)
	return &userApplicationService{
		db:                   db,
		config:               cfg,
		wg:                   wg,
		logger:               logger,
		queries:              dbQueries,
		userService:          userService,
		emailService:         emailService,
		emailTemplateService: emailTemplateService,
	}
}

func (app *userApplicationService) GetUserByID(ctx context.Context, query queries.UserByIDQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by ID")

	user, err := app.userService.GetUserByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by ID")
		return nil, err
	}

	return user, nil
}

func (app *userApplicationService) GetUserByEmail(ctx context.Context, query queries.UserByEmailQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by email")

	user, err := app.userService.GetUserByEmail(ctx, app.queries, query.Email)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by email")
		return nil, err
	}

	return user, nil
}

func (app *userApplicationService) GetUserByHandle(ctx context.Context, query queries.UserByHandleQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by handle")

	user, err := app.userService.GetUserByHandle(ctx, app.queries, query.Handle)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by handle")
		return nil, err
	}

	return user, nil
}

func (app *userApplicationService) GetUserBySessionToken(ctx context.Context, query queries.UserBySessionTokenQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by sessionToken")

	userContext, err := app.userService.GetUserContextBySessionToken(ctx, app.queries, query.SessionToken)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by sessionToken")
		return nil, err
	}

	return userContext.User, nil
}

func (app *userApplicationService) CreateUser(ctx context.Context, cmd commands.CreateNewUserCommand) (*entities.UserContextEntity, error) {
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Creating new user")

	userEntity := cmd.ToDomain()

	createdUser, err := app.userService.CreateUser(ctx, qtx, userEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new user")
		return nil, err
	}

	userSession, err := app.userService.CreateSession(ctx, qtx, createdUser)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new session")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return userSession, nil
}

func (app *userApplicationService) UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entities.UserEntity, error) {
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Updating user")

	userEntity := cmd.ToDomain()

	updatedUser, err := app.userService.UpdateUser(ctx, qtx, userEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update user")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return updatedUser, nil
}

func (app *userApplicationService) InviteUser(ctx context.Context, cmd commands.InviteUserCommand) error {
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Creating new user")

	userEntity := cmd.ToDomain()

	var userForInvite *entities.UserEntity

	foundUser, err := app.userService.GetUserByEmail(ctx, qtx, userEntity.Email)
	if err != nil && !errors.Is(err, entities.ErrUserNotFound) {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by email")
		return err
	}

	if foundUser != nil && foundUser.Claimed {
		return entities.ErrUserClaimed
	}

	if foundUser == nil {
		createdUser, err := app.userService.CreateUser(ctx, qtx, userEntity)
		if err != nil {
			app.logger.Err(err).Ctx(ctx).Msg("Failed to create new user")
			return err
		}
		userForInvite = createdUser
	} else {
		userForInvite = foundUser
	}

	inviteLink, err := app.userService.CreateInviteLink(ctx, app.queries, userForInvite)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create login link")
		return err
	}

	plainBody, htmlBody, err := app.emailTemplateService.LoginEmail("invite.go.tmpl", inviteLink)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create email template")
		return err
	}

	emailEntity := entities.EmailEntity{
		ID:        uuid.New(),
		ToEmail:   userForInvite.Email,
		FromEmail: "mcorrigan89@gmail.com",
		Subject:   "Invite to Big App",
		PlainBody: plainBody,
		HtmlBody:  htmlBody,
	}

	_, err = app.emailService.SendEmail(ctx, app.queries, &emailEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to send email")
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

func (app *userApplicationService) RequestEmailLogin(ctx context.Context, cmd commands.RequestEmailLoginCommand) error {
	app.logger.Info().Ctx(ctx).Msg("Requesting email login")

	email := cmd.ToDomain()

	loginLink, err := app.userService.CreateLoginLink(ctx, app.queries, email)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create login link")
		return err
	}

	plainBody, htmlBody, err := app.emailTemplateService.LoginEmail("login.go.tmpl", loginLink)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create email template")
		return err
	}

	emailEntity := entities.EmailEntity{
		ID:        uuid.New(),
		ToEmail:   email,
		FromEmail: "mcorrigan89@gmail.com",
		Subject:   "Login to Big App",
		PlainBody: plainBody,
		HtmlBody:  htmlBody,
	}

	_, err = app.emailService.SendEmail(ctx, app.queries, &emailEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to send email")
		return err
	}

	return nil
}

func (app *userApplicationService) LoginWithReferenceLink(ctx context.Context, cmd commands.LoginWithReferenceLinkCommand) (*entities.UserContextEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Login with reference link")

	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	userSession, err := app.userService.LoginWithLink(ctx, qtx, cmd.ReferenceLinkToken)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create login link")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return userSession, nil
}

func (app *userApplicationService) AcceptInviteReferenceLink(ctx context.Context, cmd commands.AcceptInviteReferenceLinkCommand) (*entities.UserContextEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Login with reference link")

	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	userSession, err := app.userService.AcceptInviteLink(ctx, qtx, cmd.ReferenceLinkToken)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to accept invite link")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return userSession, nil
}
