package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/rs/zerolog"
)

type postgresEventRepository struct {
	logger *zerolog.Logger
}

func NewPostgresEventRepository(logger *zerolog.Logger) *postgresEventRepository {
	return &postgresEventRepository{
		logger: logger,
	}
}

func (repo *postgresEventRepository) GetEventByID(ctx context.Context, querier models.Querier, eventID uuid.UUID) (*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetEventByID(ctx, eventID)
	if err != nil {
		repo.logger.Err(err).Ctx(ctx).Msg("Failed to get event by ID")
		return nil, err
	}

	markerModels := make([]*models.TimeslotMarker, 0)

	if row.Markers != nil {
		err = json.Unmarshal(row.Markers, &markerModels)
		if err != nil {
			repo.logger.Err(err).Ctx(ctx).Msg("Failed to unmarshal markers")
			return nil, err
		}
	}

	timeslotRows, err := querier.TimeSlotsByEventID(ctx, eventID)
	if err != nil {
		repo.logger.Err(err).Ctx(ctx).Msg("Failed to get timeslots by event ID")
		return nil, err
	}

	timeslotArgs := make([]*entities.NewEventEntitySlotsArgs, 0)
	for _, timeslotRow := range timeslotRows {
		timeslotArgs = append(timeslotArgs, &entities.NewEventEntitySlotsArgs{
			TimeSlot: timeslotRow.Timeslot,
			Artist:   timeslotRow.Artist,
		})
	}

	return entities.NewEventEntity(row.Event, timeslotArgs, markerModels), nil
}

func (repo *postgresEventRepository) GetEvents(ctx context.Context, querier models.Querier, afterDate time.Time) ([]*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	rows, err := querier.GetAllEvents(ctx, afterDate)
	if err != nil {
		return nil, err
	}

	var eventEntities []*entities.EventEntity
	for _, row := range rows {
		markerModels := make([]*models.TimeslotMarker, 0)

		if row.Markers != nil {
			err = json.Unmarshal(row.Markers, &markerModels)
			if err != nil {
				repo.logger.Err(err).Ctx(ctx).Msg("Failed to unmarshal markers")
				return nil, err
			}
		}

		timeslotRows, err := querier.TimeSlotsByEventID(ctx, row.Event.ID)
		if err != nil {
			repo.logger.Err(err).Ctx(ctx).Msg("Failed to get timeslots by event ID")
			return nil, err
		}

		timeslotArgs := make([]*entities.NewEventEntitySlotsArgs, 0)
		for _, timeslotRow := range timeslotRows {
			timeslotArgs = append(timeslotArgs, &entities.NewEventEntitySlotsArgs{
				TimeSlot: timeslotRow.Timeslot,
				Artist:   timeslotRow.Artist,
			})
		}

		eventEntities = append(eventEntities, entities.NewEventEntity(row.Event, timeslotArgs, markerModels))
	}

	return eventEntities, nil
}

func (repo *postgresEventRepository) CreateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateEvent(ctx, models.CreateEventParams{
		ID:        event.ID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		EventType: event.EventType,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewEventEntity(row, nil, nil), nil
}

func (repo *postgresEventRepository) UpdateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.UpdateEvent(ctx, models.UpdateEventParams{
		ID:        event.ID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		EventType: event.EventType,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewEventEntity(row, nil, nil), nil
}

func (repo *postgresEventRepository) DeleteEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.DeleteEvent(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresEventRepository) UpdateTimeSlot(ctx context.Context, querier models.Querier, timeslot *entities.TimeSlotEntity) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	_, err := querier.UpdateTimeSlot(ctx, models.UpdateTimeSlotParams{
		ID:                 timeslot.ID,
		ArtistNameOverride: timeslot.NameOverride,
		SortKey:            timeslot.SortKey,
		SongCount:          timeslot.SongCount,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresEventRepository) AddArtistToEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID, sortKey string, artistNameOverride *string) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.AddArtistToEvent(ctx, models.AddArtistToEventParams{
		ID:                 uuid.New(),
		EventID:            eventID,
		ArtistID:           artistID,
		ArtistNameOverride: artistNameOverride,
		SortKey:            sortKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresEventRepository) RemoveArtistFromEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.RemoveArtistFromEvent(ctx, models.RemoveArtistFromEventParams{
		EventID:  eventID,
		ArtistID: artistID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresEventRepository) CreateTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, timeSlotMarker *entities.TimeMarkerEntity) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	_, err := querier.CreateTimeslotMarker(ctx, models.CreateTimeslotMarkerParams{
		EventID:       eventID,
		ID:            timeSlotMarker.ID,
		MarkerType:    timeSlotMarker.Type,
		MarkerValue:   timeSlotMarker.Time,
		TimeslotIndex: int32(timeSlotMarker.Index),
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *postgresEventRepository) UpdateTimeslotMarker(ctx context.Context, querier models.Querier, timeSlotMarker *entities.TimeMarkerEntity) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	_, err := querier.UpdateTimeslotMarker(ctx, models.UpdateTimeslotMarkerParams{
		ID:            timeSlotMarker.ID,
		TimeslotIndex: int32(timeSlotMarker.Index),
		MarkerValue:   timeSlotMarker.Time,
		MarkerType:    timeSlotMarker.Type,
	})
	if err != nil {
		return err
	}

	return nil

}

func (repo *postgresEventRepository) DeleteTimeslotMarker(ctx context.Context, querier models.Querier, timeslotMarkerID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.DeleteTimeslotMarker(ctx, timeslotMarkerID)
	if err != nil {
		return err
	}

	return nil
}
