package application

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/services"
	"github.com/mcorrigan89/openmic/internal/infrastructure/bus"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"

	"github.com/rs/zerolog"
)

type EventApplicationService interface {
	GetEventByID(ctx context.Context, query queries.EventByIDQuery) (*entities.EventEntity, error)
	GetCurrentEvent(ctx context.Context, query queries.CurrentEventQuery) (*entities.EventEntity, error)
	GetEvents(ctx context.Context, query queries.EventsQuery) ([]*entities.EventEntity, error)
	CreateEvent(ctx context.Context, cmd commands.CreateNewEventCommand) (*entities.EventEntity, error)
	UpdateEvent(ctx context.Context, cmd commands.UpdateEventCommand) (*entities.EventEntity, error)
	DeleteEvent(ctx context.Context, query commands.DeleteEventCommand) error
	AddArtistToEvent(ctx context.Context, cmd commands.AddArtistToEventCommand) (*entities.EventEntity, error)
	RemoveArtistFromEvent(ctx context.Context, cmd commands.RemoveArtistFromEventCommand) (*entities.EventEntity, error)
	SetTimeslotMarker(ctx context.Context, cmd commands.SetTimeslotMarkerCommand) (*entities.EventEntity, error)
	DeleteTimeslotMarker(ctx context.Context, cmd commands.DeleteTimeslotMarkerCommand) (*entities.EventEntity, error)
	SetSortOrder(ctx context.Context, cmd commands.SetSortOrderCommand) (*entities.EventEntity, error)
	UpdateTimeSlot(ctx context.Context, cmd commands.UpdateTimeSlotCommand) (*entities.EventEntity, error)
	SetNowPlaying(ctx context.Context, cmd commands.SetNowPlayingCommand) (*entities.EventEntity, error)
	MessageBus() *bus.MessageBus[*dto.EventDto]
}

type eventApplicationService struct {
	config       *common.Config
	wg           *sync.WaitGroup
	logger       *zerolog.Logger
	db           *pgxpool.Pool
	queries      models.Querier
	bus          *bus.MessageBus[*dto.EventDto]
	eventService services.EventService
}

func NewEventApplicationService(db *pgxpool.Pool, wg *sync.WaitGroup, cfg *common.Config, logger *zerolog.Logger, bus *bus.MessageBus[*dto.EventDto], eventService services.EventService) *eventApplicationService {
	dbQueries := models.New(db)
	return &eventApplicationService{
		db:           db,
		config:       cfg,
		wg:           wg,
		logger:       logger,
		queries:      dbQueries,
		bus:          bus,
		eventService: eventService,
	}
}

func (app *eventApplicationService) MessageBus() *bus.MessageBus[*dto.EventDto] {
	return app.bus
}

func (app *eventApplicationService) GetEventByID(ctx context.Context, query queries.EventByIDQuery) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting event by ID")

	event, err := app.eventService.GetEventByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	app.logger.Info().Ctx(ctx).Msg("Got event by ID")

	return event, nil
}

func (app *eventApplicationService) GetCurrentEvent(ctx context.Context, query queries.CurrentEventQuery) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting current event")

	events, err := app.eventService.GetEvents(ctx, app.queries)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get all events")
		return nil, err
	}

	// Show up 24 hours before start date

	var currentEventID uuid.UUID
	currentTime := time.Now()

	for _, event := range events {
		eventTimeMinus := event.StartTime.Add(-30 * time.Hour)
		eventTimePlus := event.StartTime.Add(8 * time.Hour)
		if eventTimeMinus.After(currentTime) && eventTimePlus.After(currentTime) && event.EventType == "OPEN_MIC" {
			currentEventID = event.ID
			break
		}
	}

	if currentEventID == uuid.Nil {
		return nil, nil
	}

	event, err := app.eventService.GetEventByID(ctx, app.queries, events[0].ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) GetEvents(ctx context.Context, query queries.EventsQuery) ([]*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting all events")

	events, err := app.eventService.GetEvents(ctx, app.queries)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get all events")
		return nil, err
	}

	return events, nil
}

func (app *eventApplicationService) CreateEvent(ctx context.Context, cmd commands.CreateNewEventCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Creating new event")

	eventEntity := cmd.ToDomain()

	event, err := app.eventService.CreateEvent(ctx, app.queries, eventEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new event")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) UpdateEvent(ctx context.Context, cmd commands.UpdateEventCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Updating event")

	eventEntity := cmd.ToDomain()

	event, err := app.eventService.UpdateEvent(ctx, app.queries, eventEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update event")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) DeleteEvent(ctx context.Context, query commands.DeleteEventCommand) error {
	app.logger.Info().Ctx(ctx).Msg("Deleting event")

	err := app.eventService.DeleteEvent(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to delete event")
		return err
	}

	return nil
}

func (app *eventApplicationService) AddArtistToEvent(ctx context.Context, cmd commands.AddArtistToEventCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Adding artist to event")

	err := app.eventService.AddArtistToEvent(ctx, app.queries, cmd.EventID, cmd.ArtistID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to add artist to event")
		return nil, err
	}

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) RemoveArtistFromEvent(ctx context.Context, cmd commands.RemoveArtistFromEventCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Removing artist from event")

	err := app.eventService.RemoveArtistFromEvent(ctx, app.queries, cmd.EventID, cmd.ArtistID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to remove artist from event")
		return nil, err
	}

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) SetTimeslotMarker(ctx context.Context, cmd commands.SetTimeslotMarkerCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Setting timeslot marker")

	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	err = app.eventService.SetTimeslotMarker(ctx, qtx, cmd.EventID, cmd.SlotIndex, cmd.TimeDisplay)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to set timeslot")
		return nil, err
	}

	event, err := app.eventService.GetEventByID(ctx, qtx, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) DeleteTimeslotMarker(ctx context.Context, cmd commands.DeleteTimeslotMarkerCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Deleting timeslot marker")

	err := app.eventService.DeleteTimeslotMarker(ctx, app.queries, cmd.EventID, cmd.SlotMarkerID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to set timeslot")
		return nil, err
	}

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) SetSortOrder(ctx context.Context, cmd commands.SetSortOrderCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Setting sort order")

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	currentSlot := event.TimeSlotByID(cmd.CurrentSlotID)
	var beforeSlot, afterSlot *entities.TimeSlotEntity
	var beforeSlotSortKey, afterSlotSortKey string
	if currentSlot == nil {
		return event, nil
	}

	if cmd.BeforeSlotID != nil {
		beforeSlot = event.TimeSlotByID(*cmd.BeforeSlotID)
		afterSlot = event.NextTimeSlotByID(*cmd.BeforeSlotID)
	}

	if cmd.AfterSlotID != nil {
		afterSlot = event.TimeSlotByID(*cmd.AfterSlotID)
		beforeSlot = event.PreviousTimeSlotByID(*cmd.AfterSlotID)
	}

	if beforeSlot != nil {
		beforeSlotSortKey = beforeSlot.SortKey
	} else {
		beforeSlotSortKey = ""
	}

	if afterSlot != nil {
		afterSlotSortKey = afterSlot.SortKey
	} else {
		afterSlotSortKey = ""
	}

	sortKey, err := common.KeyBetween(beforeSlotSortKey, afterSlotSortKey)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to generate sort key")
		return nil, err
	}

	currentSlot.SortKey = sortKey

	err = app.eventService.UpdateTimeSlot(ctx, app.queries, currentSlot)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update timeslot")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) UpdateTimeSlot(ctx context.Context, cmd commands.UpdateTimeSlotCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Updating timeslot")

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	timeslot := event.TimeSlotByID(cmd.TimeSlotID)
	if timeslot == nil {
		return event, nil
	}

	timeslot.SongCount = cmd.SongCount

	err = app.eventService.UpdateTimeSlot(ctx, app.queries, timeslot)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update timeslot")
		return nil, err
	}

	return event, nil
}

func (app *eventApplicationService) SetNowPlaying(ctx context.Context, cmd commands.SetNowPlayingCommand) (*entities.EventEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Setting now playing")

	err := app.eventService.SetNowPlaying(ctx, app.queries, cmd.EventID, cmd.Index)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update timeslot")
		return nil, err
	}

	event, err := app.eventService.GetEventByID(ctx, app.queries, cmd.EventID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	return event, nil
}
