package commands

import (
	"io"

	"github.com/google/uuid"
)

type CreateNewAvatarImageCommand struct {
	UserID   uuid.UUID
	BucketID string
	ObjectID string
	File     io.Reader
	Size     int64
}

type ImageUploadData struct {
	ObjectID string
	File     io.Reader
	Size     int64
}
