package handlers

import (
	"context"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/application"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/rs/zerolog"
)

type EventHandler struct {
	logger          *zerolog.Logger
	eventAppService application.EventApplicationService
}

func NewEventHandler(logger *zerolog.Logger, eventAppService application.EventApplicationService) *EventHandler {
	return &EventHandler{
		logger:          logger,
		eventAppService: eventAppService,
	}
}

func (h *EventHandler) GetEventByID(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*dto.GetEventByIDResponse, error) {

	query := queries.EventByIDQuery{
		ID: input.ID,
	}

	event, err := h.eventAppService.GetEventByID(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get event by ID", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.GetEventByIDResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) GetCurrentEvent(ctx context.Context, input *struct{}) (*dto.GetCurrentEventResponse, error) {
	query := queries.CurrentEventQuery{}

	event, err := h.eventAppService.GetCurrentEvent(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get event by ID", err)
	}

	if event == nil {
		return &dto.GetCurrentEventResponse{
			Body: nil,
		}, nil
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.GetCurrentEventResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) GetUpcomingEvents(ctx context.Context, input *struct {
}) (*dto.GetEventsResponse, error) {

	query := queries.EventsQuery{}

	events, err := h.eventAppService.GetEvents(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get event by ID", err)
	}

	eventDtos := make([]*dto.EventDto, 0)
	for _, event := range events {
		eventDto := dto.NewEventDtoFromEntity(event)
		eventDtos = append(eventDtos, eventDto)
	}

	return &dto.GetEventsResponse{
		Body: eventDtos,
	}, nil
}

func (h *EventHandler) CreateEvent(ctx context.Context, input *dto.CreateEventRequest) (*dto.CreateEventResponse, error) {

	cmd := commands.CreateNewEventCommand{
		StartTime: input.Body.StartTime,
		EndTime:   input.Body.EndTime,
		EventType: input.Body.EventType,
	}

	event, err := h.eventAppService.CreateEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create event", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.CreateEventResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) UpdateEvent(ctx context.Context, input *dto.UpdateEventRequest) (*dto.UpdateEventResponse, error) {

	cmd := commands.UpdateEventCommand{
		ID:        input.ID,
		StartTime: input.Body.StartTime,
		EndTime:   input.Body.EndTime,
		EventType: input.Body.EventType,
	}

	event, err := h.eventAppService.UpdateEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update event", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.UpdateEventResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) DeleteEvent(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*dto.DeleteEventResponse, error) {

	cmd := commands.DeleteEventCommand{
		ID: input.ID,
	}

	err := h.eventAppService.DeleteEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete event", err)
	}

	msg := dto.DeleteEventResponse{
		Body: "Event deleted",
	}

	return &msg, nil
}

func (h *EventHandler) AddArtistToEvent(ctx context.Context, input *dto.AddArtistToEventEventRequest) (*dto.AddArtistToEventEventResponst, error) {

	cmd := commands.AddArtistToEventCommand{
		EventID:  input.EventID,
		ArtistID: input.Body.ArtistID,
	}

	event, err := h.eventAppService.AddArtistToEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to add artist to event", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	h.eventAppService.MessageBus().Publish(eventDto)

	return &dto.AddArtistToEventEventResponst{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) RemoveArtistFromEvent(ctx context.Context, input *dto.RemoveArtistFromEventEventRequest) (*dto.RemoveArtistFromEventEventResponse, error) {

	cmd := commands.RemoveArtistFromEventCommand{
		EventID:  input.EventID,
		ArtistID: input.Body.ArtistID,
	}

	event, err := h.eventAppService.RemoveArtistFromEvent(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to remove artist from event", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.RemoveArtistFromEventEventResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) SetTimeslotMarker(ctx context.Context, input *dto.SetTimeslotMarkerRequest) (*dto.SetTimeslotMarkerResponse, error) {

	cmd := commands.SetTimeslotMarkerCommand{
		EventID:     input.EventID,
		TimeDisplay: input.Body.TimeDisplay,
		SlotIndex:   input.Body.SlotIndex,
	}

	event, err := h.eventAppService.SetTimeslotMarker(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to set timeslot", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.SetTimeslotMarkerResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) DeleteTimeslotMarker(ctx context.Context, input *dto.DeleteTimeslotMarkerRequest) (*dto.SetTimeslotMarkerResponse, error) {

	cmd := commands.DeleteTimeslotMarkerCommand{
		EventID:      input.EventID,
		SlotMarkerID: input.Body.TimeslotMarkerID,
	}

	event, err := h.eventAppService.DeleteTimeslotMarker(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to set timeslot", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.SetTimeslotMarkerResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) SetSortOrderRequest(ctx context.Context, input *dto.SetSortOrderRequest) (*dto.SetSortOrderResponse, error) {

	cmd := commands.SetSortOrderCommand{
		EventID:       input.EventID,
		BeforeSlotID:  input.Body.BeforeSlotID,
		CurrentSlotID: input.Body.CurrentSlotID,
		AfterSlotID:   input.Body.AfterSlotID,
	}

	event, err := h.eventAppService.SetSortOrder(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to set sort order", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.SetSortOrderResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) UpdateTimeSlot(ctx context.Context, input *dto.UpdateTimeSlotRequest) (*dto.UpdateTimeSlotResponse, error) {

	cmd := commands.UpdateTimeSlotCommand{
		EventID:    input.EventID,
		TimeSlotID: input.TimeSlotID,
		SongCount:  input.Body.SongCount,
	}

	event, err := h.eventAppService.UpdateTimeSlot(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update timeslot", err)
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	return &dto.UpdateTimeSlotResponse{
		Body: eventDto,
	}, nil
}

func (h *EventHandler) ListenForEventChange(ctx context.Context, input *struct {
	ID uuid.UUID `path:"event_id"`
}, send sse.Sender) {

	event, err := h.eventAppService.GetEventByID(ctx, queries.EventByIDQuery{
		ID: input.ID,
	})

	if err != nil {
		h.logger.Err(err).Msg("Failed to get event by ID")
		return
	}

	eventDto := dto.NewEventDtoFromEntity(event)

	send(sse.Message{
		Retry: 5000,
		Data: dto.ListenForChangeEventResponse{
			Body: eventDto,
		},
	})

	var wg sync.WaitGroup

	c, clientID, err := h.eventAppService.MessageBus().Subscribe()
	if err != nil {
		h.logger.Err(err).Msg("Failed to subscribe to event changes")
		return
	}
	defer h.eventAppService.MessageBus().Unsubscribe(clientID)

	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case msg, ok := <-c:
			if !ok {
				return
			}
			send(sse.Message{
				Data: dto.ListenForChangeEventResponse{
					Body: msg,
				},
			})
		case <-ctx.Done():
			return

		}
	}()

	wg.Wait()

	h.eventAppService.MessageBus().Unsubscribe(clientID)
}
