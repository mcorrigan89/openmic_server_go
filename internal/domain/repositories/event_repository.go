package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type EventRepository interface {
	GetEventByID(ctx context.Context, querier models.Querier, id uuid.UUID) (*entities.EventEntity, error)
	GetEvents(ctx context.Context, querier models.Querier, afterDate time.Time) ([]*entities.EventEntity, error)
	CreateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error)
	UpdateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error)
	DeleteEvent(ctx context.Context, querier models.Querier, id uuid.UUID) error
	UpdateTimeSlot(ctx context.Context, querier models.Querier, timeslot *entities.TimeSlotEntity) error
	AddArtistToEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID, sortKet string, artistNameOverride *string) error
	RemoveArtistFromEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error
	CreateTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, markerEntity *entities.TimeMarkerEntity) error
	UpdateTimeslotMarker(ctx context.Context, querier models.Querier, markerEntity *entities.TimeMarkerEntity) error
	DeleteTimeslotMarker(ctx context.Context, querier models.Querier, timeslotMarkerID uuid.UUID) error
}
