package dto

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

type ImageDto struct {
	ID     uuid.UUID `json:"id"`
	Url    string    `json:"url"`
	Width  int32     `json:"width"`
	Height int32     `json:"height"`
	Size   int32     `json:"size"`
}

func NewImageDtoFromEntity(entity *entities.ImageEntity) *ImageDto {
	return &ImageDto{
		ID:     entity.ID,
		Url:    entity.UrlSlug(),
		Width:  entity.Width,
		Height: entity.Height,
		Size:   entity.Size,
	}
}

func (dto *ImageDto) ToJson() ([]byte, error) {
	return json.Marshal(dto)
}
