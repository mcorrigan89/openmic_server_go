package queries

import (
	"github.com/google/uuid"
)

type ArtistByIDQuery struct {
	ID uuid.UUID
}

type ArtistsByTitleQuery struct {
	Title string
}
