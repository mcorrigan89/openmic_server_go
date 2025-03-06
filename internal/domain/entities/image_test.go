package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/assert"
)

func TestNewImageEntity(t *testing.T) {

	t.Run("create entity from model", func(t *testing.T) {
		imageId := uuid.New()
		bucketName := "test-bucket"
		fileName := "test-file"
		height := int32(100)
		width := int32(100)
		size := int32(1000)

		imageModel := models.Image{
			ID:         imageId,
			BucketName: bucketName,
			ObjectID:   fileName,
			Height:     height,
			Width:      width,
			FileSize:   size,
		}

		imageEntity := NewImageEntity(imageModel)

		assert.Equal(t, imageEntity.ID, imageId)
		assert.Equal(t, imageEntity.BucketName, bucketName)
		assert.Equal(t, imageEntity.ObjectID, fileName)
		assert.Equal(t, imageEntity.Height, height)
		assert.Equal(t, imageEntity.Width, width)
		assert.Equal(t, imageEntity.UrlSlug(), "/image/"+imageId.String())
	})
}
