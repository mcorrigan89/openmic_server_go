package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/mcorrigan89/openmic/internal/application"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/services"
	"github.com/mcorrigan89/openmic/internal/infrastructure/bus"
	"github.com/mcorrigan89/openmic/internal/infrastructure/email"
	"github.com/mcorrigan89/openmic/internal/infrastructure/media"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/storage"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/handlers"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/middleware"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/router"

	"github.com/rs/zerolog"
)

type appServer struct {
	config *common.Config
	wg     *sync.WaitGroup
	logger *zerolog.Logger
}

func main() {
	logger := getLogger()

	logger.Info().Msg("Starting server")

	cfg := common.Config{}
	common.LoadConfig(&cfg)

	db, err := postgres.OpenDBPool(&cfg)
	if err != nil {
		logger.Err(err).Msg("Failed to open database connection")
		os.Exit(1)
	}
	defer db.Close()

	wg := sync.WaitGroup{}
	mux := http.NewServeMux()

	messageBus := bus.NewMessageBus[*dto.EventDto]()
	defer messageBus.Close()

	postgresUserRepository := repositories.NewPostgresUserRepository()
	postgresArtistRepositoy := repositories.NewPostgresArtistRepository()
	postgresEventRepositoy := repositories.NewPostgresEventRepository()
	postgresReferenceLinkRepository := repositories.NewPostgresReferenceLinkRepository()
	postgresImageRepository := repositories.NewPostgresImageRepository()
	blobStorageService := storage.NewBlobStorageService(&cfg)
	smtpService := email.NewSmtpService(&cfg)
	imageMediaService := media.NewImageMediaService(blobStorageService)

	userService := services.NewUserService(postgresUserRepository, postgresReferenceLinkRepository, postgresImageRepository)
	artistService := services.NewArtistService(&logger, postgresArtistRepositoy)
	eventService := services.NewEventService(&logger, postgresEventRepositoy)
	emailService := services.NewEmailService(smtpService)
	emailTemplateService := services.NewEmailTemplateService(&cfg)
	imageService := services.NewImageService(postgresImageRepository)

	userApplicationService := application.NewUserApplicationService(db, &wg, &cfg, &logger, userService, emailService, emailTemplateService)
	imageApplicationService := application.NewImageApplicationService(db, &wg, &cfg, &logger, imageService, userService, imageMediaService)
	artistApplicationService := application.NewArtistApplicationService(db, &wg, &cfg, &logger, artistService)
	eventApplicationService := application.NewEventApplicationService(db, &wg, &cfg, &logger, messageBus, eventService)
	userHandler := handlers.NewUserHandler(&logger, userApplicationService)
	imageHandler := handlers.NewImageHandler(&logger, imageApplicationService)
	artistHandler := handlers.NewArtistHandler(&logger, artistApplicationService)
	eventHandler := handlers.NewEventHandler(&logger, eventApplicationService)

	mdlwr := middleware.CreateMiddleware(&cfg, db, &logger, userService)

	// HTTP Routes
	httpRoutes := router.NewRouter(mux, mdlwr, userHandler, imageHandler, eventHandler, artistHandler)

	server := &appServer{
		wg:     &wg,
		config: &cfg,
		logger: &logger,
	}

	err = server.serve(httpRoutes)
	if err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
