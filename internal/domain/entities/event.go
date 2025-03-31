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
	ID           uuid.UUID
	NameOverride *string
	SortKey      string
	Artist       *ArtistEntity
	SongCount    int32
	TimeDisplay  time.Time
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

func (e *EventEntity) TimeSlotByID(id uuid.UUID) *TimeSlotEntity {
	var timeSlot *TimeSlotEntity
	for _, slot := range e.timeSlots {
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

func (e *EventEntity) PreviousTimeSlotByID(id uuid.UUID) *TimeSlotEntity {
	var timeSlot *TimeSlotEntity
	for idx, slot := range e.timeSlots {
		if slot.ID == id {
			if idx > 0 && e.timeSlots[idx-1] != nil {
				timeSlot = e.timeSlots[idx-1]
			}
			break
		}
	}
	if timeSlot == nil {
		return nil
	}

	return timeSlot
}

func (e *EventEntity) NextTimeSlotByID(id uuid.UUID) *TimeSlotEntity {
	var timeSlot *TimeSlotEntity
	for idx, slot := range e.timeSlots {
		if slot.ID == id {
			if idx+1 < len(e.timeSlots) {
				timeSlot = e.timeSlots[idx+1]
			}
			break
		}
	}
	if timeSlot == nil {
		return nil
	}

	return timeSlot
}

func (e *EventEntity) TimeSlotMarkerByDisplay(timeDisplay string) *TimeMarkerEntity {
	var marker *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.Time == timeDisplay {
			marker = slot
			break
		}
	}
	if marker == nil {
		return nil
	}

	return marker
}

func (e *EventEntity) TimeSlotMarkerDupeByIndex(index int, display string) *TimeMarkerEntity {
	var marker *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.Index == index && slot.Time != display {
			marker = slot
			break
		}
	}
	if marker == nil {
		return nil
	}

	return marker
}

func (e *EventEntity) NowPlayingTimeSlotMarker() *TimeMarkerEntity {
	var marker *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.Type == "PLAYING" {
			marker = slot
			break
		}
	}
	if marker == nil {
		return nil
	}

	return marker
}

func (e *EventEntity) TimeSlotMarkerByID(id uuid.UUID) *TimeMarkerEntity {
	var marker *TimeMarkerEntity
	for _, slot := range e.markers {
		if slot.ID == id {
			marker = slot
			break
		}
	}
	if marker == nil {
		return nil
	}

	return marker
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
