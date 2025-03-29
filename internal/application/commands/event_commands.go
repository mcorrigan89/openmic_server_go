package commands

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type CreateNewEventCommand struct {
	StartTime time.Time
	EndTime   time.Time
	EventType string
}

func (cmd *CreateNewEventCommand) ToDomain() *entities.EventEntity {
	return &entities.EventEntity{
		ID:        uuid.New(),
		StartTime: cmd.StartTime,
		EndTime:   cmd.EndTime,
		EventType: cmd.EventType,
	}
}

type UpdateEventCommand struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	EventType string
}

func (cmd *UpdateEventCommand) ToDomain() *entities.EventEntity {
	return &entities.EventEntity{
		ID:        cmd.ID,
		StartTime: cmd.StartTime,
		EndTime:   cmd.EndTime,
		EventType: cmd.EventType,
	}
}

type DeleteEventCommand struct {
	ID uuid.UUID
}

type AddArtistToEventCommand struct {
	EventID  uuid.UUID
	ArtistID uuid.UUID
}

type RemoveArtistFromEventCommand struct {
	EventID  uuid.UUID
	ArtistID uuid.UUID
}

type SetTimeslotMarkerCommand struct {
	EventID     uuid.UUID
	TimeDisplay string
	SlotIndex   int
}

type DeleteTimeslotMarkerCommand struct {
	EventID      uuid.UUID
	SlotMarkerID uuid.UUID
}

type SetSortOrderCommand struct {
	EventID       uuid.UUID
	BeforeSlotID  *uuid.UUID
	CurrentSlotID uuid.UUID
}
