package queries

import (
	"github.com/google/uuid"
)

type EventByIDQuery struct {
	ID uuid.UUID
}

type EventsQuery struct {
}
