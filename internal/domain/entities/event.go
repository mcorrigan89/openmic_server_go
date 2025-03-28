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
	markers   []*TimeMarkerEntity
}

type TimeSlotEntity struct {
	ID          uuid.UUID
	SortKey     string
	Artist      *ArtistEntity
	SongCount   int32
	TimeDisplay time.Time
}

type TimeMarkerEntity struct {
	ID    uuid.UUID
	Index int
	Type  string
	Time  string
}

type NewEventEntitySlotsArgs struct {
	TimeSlot models.Timeslot
	Artist   models.Artist
}

func NewEventEntity(eventModel models.Event, timeSlotArgs []*NewEventEntitySlotsArgs, timeMarkers []*models.TimeslotMarker) *EventEntity {

	timeSlotAggregator := eventModel.StartTime
	timeSlotEntities := make([]*TimeSlotEntity, 0)
	for _, timeslotArg := range timeSlotArgs {
		timeSlotEntities = append(timeSlotEntities, newTimeSlotEntity(timeslotArg.TimeSlot, timeslotArg.Artist, timeSlotAggregator))

		if timeslotArg.TimeSlot.SongCount == 1 {
			dur := time.Minute * 5
			timeSlotAggregator = timeSlotAggregator.Add(dur)
		} else {
			dur := time.Minute * 8
			timeSlotAggregator = timeSlotAggregator.Add(dur)
		}
	}

	timeMarkerEntities := make([]*TimeMarkerEntity, 0)

	for _, timeMarker := range timeMarkers {
		timeMarkerEntities = append(timeMarkerEntities, &TimeMarkerEntity{
			ID:    timeMarker.ID,
			Index: int(timeMarker.TimeslotIndex),
			Type:  timeMarker.MarkerType,
			Time:  timeMarker.MarkerValue,
		})
	}

	return &EventEntity{
		ID:        eventModel.ID,
		StartTime: eventModel.StartTime,
		EndTime:   eventModel.EndTime,
		EventType: eventModel.EventType,
		timeSlots: timeSlotEntities,
		markers:   timeMarkerEntities,
	}
}

func (e *EventEntity) TimeSlots() []*TimeSlotEntity {
	return e.timeSlots
}

func (e *EventEntity) TimeMarkers() []*TimeMarkerEntity {
	return e.markers
}

func (e *EventEntity) TimeSlotMarkerByDisplay(timeDisplay string) *TimeMarkerEntity {
	var timeSlot *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.Time == timeDisplay {
			timeSlot = slot
			break
		}
	}
	if timeSlot == nil {
		return nil
	}

	return timeSlot
}

func (e *EventEntity) TimeSlotMarkerByID(id uuid.UUID) *TimeMarkerEntity {
	var timeSlot *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.ID == id {
			timeSlot = slot
			break
		}
	}
	if timeSlot == nil {
		return nil
	}

	return timeSlot
}

func newTimeSlotEntity(timeSlotModel models.Timeslot, artistModel models.Artist, slotTime time.Time) *TimeSlotEntity {
	return &TimeSlotEntity{
		ID:          timeSlotModel.ID,
		SortKey:     timeSlotModel.SortKey,
		SongCount:   timeSlotModel.SongCount,
		TimeDisplay: slotTime,
		Artist:      NewArtistEntity(artistModel),
	}
}
