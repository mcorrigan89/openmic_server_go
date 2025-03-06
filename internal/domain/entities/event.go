package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type EventEntity struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	EventType string
	timeSlots []*TimeSlotEntity
}

func NewEventEntity(eventModel models.Event, timeSlotEntities []*TimeSlotEntity) *EventEntity {
	return &EventEntity{
		ID:        eventModel.ID,
		StartTime: eventModel.StartTime,
		EndTime:   eventModel.EndTime,
		EventType: eventModel.EventType,
		timeSlots: timeSlotEntities,
	}
}

func (e *EventEntity) TimeSlots() []*TimeSlotEntity {
	return e.timeSlots
}

type TimeSlotEntity struct {
	ID      uuid.UUID
	SortKey string
	Artist  *ArtistEntity
}

func NewTimeSlotEntity(timeSlotModel models.Timeslot, artistModel models.Artist) *TimeSlotEntity {
	return &TimeSlotEntity{
		ID:      timeSlotModel.ID,
		SortKey: timeSlotModel.SortOrder,
		Artist:  NewArtistEntity(artistModel),
	}
}
