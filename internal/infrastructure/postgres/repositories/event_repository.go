package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type postgresEventRepository struct {
}

func NewPostgresEventRepository() *postgresEventRepository {
	return &postgresEventRepository{}
}

func (repo *postgresEventRepository) GetEventByID(ctx context.Context, querier models.Querier, eventID uuid.UUID) (*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	timeslotRows, err := querier.TimeSlotsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	timeslots := make([]*entities.TimeSlotEntity, 0)
	for _, timeslotRow := range timeslotRows {
		timeslots = append(timeslots, entities.NewTimeSlotEntity(timeslotRow.Timeslot, timeslotRow.Artist))
	}

	return entities.NewEventEntity(row.Event, timeslots), nil
}

func (repo *postgresEventRepository) GetEvents(ctx context.Context, querier models.Querier) ([]*entities.EventEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	rows, err := querier.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	var eventEntities []*entities.EventEntity
	for _, row := range rows {
		eventEntities = append(eventEntities, entities.NewEventEntity(row.Event, nil))
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

	return entities.NewEventEntity(row, nil), nil
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

	return entities.NewEventEntity(row, nil), nil
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

func (repo *postgresEventRepository) AddArtistToEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID, sortKey string, artistNameOverride *string) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	err := querier.AddArtistToEvent(ctx, models.AddArtistToEventParams{
		ID:                 uuid.New(),
		EventID:            eventID,
		ArtistID:           artistID,
		ArtistNameOverride: artistNameOverride,
		SortOrder:          sortKey,
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
