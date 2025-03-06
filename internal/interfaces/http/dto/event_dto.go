package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type TimeslotDto struct {
	ID     uuid.UUID  `json:"id"`
	Artist *ArtistDto `json:"artist"`
}

type EventDto struct {
	ID        uuid.UUID      `json:"id"`
	StartTime string         `json:"start_time"`
	EndTime   string         `json:"end_time"`
	EventType string         `json:"event_type"`
	TimeSlots []*TimeslotDto `json:"time_slots"`
}

func NewEventDtoFromEntity(entity *entities.EventEntity) *EventDto {

	timeslotDtos := make([]*TimeslotDto, 0)
	for _, timeslot := range entity.TimeSlots() {
		timeslotDtos = append(timeslotDtos, &TimeslotDto{
			ID:     timeslot.ID,
			Artist: NewArtistDtoFromEntity(timeslot.Artist),
		})
	}

	return &EventDto{
		ID:        entity.ID,
		StartTime: entity.StartTime.Format(time.RFC1123Z),
		EndTime:   entity.EndTime.Format(time.RFC1123Z),
		EventType: entity.EventType,
		TimeSlots: timeslotDtos,
	}
}

type GetEventByIDResponse struct {
	Body *EventDto `json:"body"`
}

type GetEventsResponse struct {
	Body []*EventDto `json:"body"`
}

type CreateEventRequest struct {
	Body struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		EventType string    `json:"event_type"`
	}
}

type CreateEventResponse struct {
	Body *EventDto `json:"body"`
}

type UpdateEventRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		EventType string    `json:"event_type"`
	}
}

type UpdateEventResponse struct {
	Body *EventDto `json:"body"`
}

type DeleteEventResponse struct {
	Body string `json:"body"`
}

type AddArtistToEventEventRequest struct {
	EventID uuid.UUID `path:"event_id"`
	Body    struct {
		ArtistID uuid.UUID `json:"artist_id"`
	}
}

type AddArtistToEventEventResponst struct {
	Body *EventDto `json:"body"`
}

type RemoveArtistFromEventEventRequest struct {
	EventID uuid.UUID `path:"event_id"`
	Body    struct {
		ArtistID uuid.UUID `json:"artist_id"`
	}
}

type RemoveArtistFromEventEventResponse struct {
	Body *EventDto `json:"body"`
}

type ListenForChangeEventResponse struct {
	Body *EventDto `json:"body"`
}
