package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/repositories"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type EventService interface {
	GetEventByID(ctx context.Context, querier models.Querier, eventID uuid.UUID) (*entities.EventEntity, error)
	GetEvents(ctx context.Context, querier models.Querier) ([]*entities.EventEntity, error)
	CreateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error)
	UpdateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error)
	DeleteEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID) error
	UpdateTimeSlot(ctx context.Context, querier models.Querier, timeslot *entities.TimeSlotEntity) error
	AddArtistToEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error
	RemoveArtistFromEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error
	SetTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, index int, timeslotDisplay string) error
	DeleteTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, timeslotMarkerID uuid.UUID) error
}

type eventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) *eventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) GetEventByID(ctx context.Context, querier models.Querier, eventID uuid.UUID) (*entities.EventEntity, error) {
	event, err := s.eventRepo.GetEventByID(ctx, querier, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventService) GetEvents(ctx context.Context, querier models.Querier) ([]*entities.EventEntity, error) {
	events, err := s.eventRepo.GetEvents(ctx, querier)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventService) CreateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error) {
	eventEntity, err := s.eventRepo.CreateEvent(ctx, querier, event)
	if err != nil {
		return nil, err
	}

	return eventEntity, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, querier models.Querier, event *entities.EventEntity) (*entities.EventEntity, error) {
	eventEntity, err := s.eventRepo.UpdateEvent(ctx, querier, event)
	if err != nil {
		return nil, err
	}

	return eventEntity, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID) error {
	err := s.eventRepo.DeleteEvent(ctx, querier, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) UpdateTimeSlot(ctx context.Context, querier models.Querier, timeslot *entities.TimeSlotEntity) error {
	err := s.eventRepo.UpdateTimeSlot(ctx, querier, timeslot)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) AddArtistToEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error {

	event, err := s.eventRepo.GetEventByID(ctx, querier, eventID)
	if err != nil {
		return err
	}
	timeSlots := event.TimeSlots()

	var sortKey string

	if len(timeSlots) > 0 {
		lastTimeslot := timeSlots[len(timeSlots)-1]
		sortKey, err = common.KeyBetween(lastTimeslot.SortKey, "")
		if err != nil {
			return err
		}
	} else {
		sortKey, err = common.KeyBetween("", "")
		if err != nil {
			return err
		}
	}

	err = s.eventRepo.AddArtistToEvent(ctx, querier, eventID, artistID, sortKey, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) RemoveArtistFromEvent(ctx context.Context, querier models.Querier, eventID uuid.UUID, artistID uuid.UUID) error {
	err := s.eventRepo.RemoveArtistFromEvent(ctx, querier, eventID, artistID)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) SetTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, index int, timeslotDisplay string) error {
	event, err := s.eventRepo.GetEventByID(ctx, querier, eventID)
	if err != nil {
		return err
	}

	markerEntity := event.TimeSlotMarkerByDisplay(timeslotDisplay)
	if markerEntity != nil {
		markerEntity.Index = index
		err = s.eventRepo.UpdateTimeslotMarker(ctx, querier, markerEntity)
		if err != nil {
			return err
		}
	} else {
		newMarker := entities.TimeMarkerEntity{
			ID:    uuid.New(),
			Time:  timeslotDisplay,
			Index: index,
			Type:  "TIME",
		}
		err = s.eventRepo.CreateTimeslotMarker(ctx, querier, eventID, &newMarker)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *eventService) DeleteTimeslotMarker(ctx context.Context, querier models.Querier, eventID uuid.UUID, timeslotMarkerID uuid.UUID) error {
	event, err := s.eventRepo.GetEventByID(ctx, querier, eventID)
	if err != nil {
		return err
	}

	markerEntity := event.TimeSlotMarkerByID(timeslotMarkerID)
	if markerEntity != nil {
		err = s.eventRepo.DeleteTimeslotMarker(ctx, querier, markerEntity.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
