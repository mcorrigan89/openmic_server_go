package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
	"github.com/stretchr/testify/assert"
)

func TestNewReferenceLinkEntity(t *testing.T) {

	t.Run("create entity from model", func(t *testing.T) {
		refLinkID := uuid.New()
		linkID := uuid.New()
		expiresAt, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		if err != nil {
			t.Fatal(err)
		}

		referenceLinkModel := models.ReferenceLink{
			ID:        refLinkID,
			LinkID:    linkID,
			LinkType:  RefLinkTypeLogin,
			Token:     "test-token",
			ExpiresAt: expiresAt,
		}

		referenceLinkEntity := NewReferenceLinkEntity(referenceLinkModel)

		assert.Equal(t, referenceLinkEntity.ID, refLinkID)
		assert.Equal(t, referenceLinkEntity.LinkID, linkID)
		assert.Equal(t, referenceLinkEntity.Type, RefLinkTypeLogin)
		assert.Equal(t, referenceLinkEntity.Token, "test-token")
		assert.Equal(t, referenceLinkEntity.ExpiresAt, expiresAt)
		assert.Equal(t, referenceLinkEntity.IsExpired(), true)
	})

	t.Run("create entity from model with invite type", func(t *testing.T) {
		refLinkID := uuid.New()
		linkID := uuid.New()
		expiresAt, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
		if err != nil {
			t.Fatal(err)
		}

		referenceLinkModel := models.ReferenceLink{
			ID:        refLinkID,
			LinkID:    linkID,
			LinkType:  RefLinkTypeInvite,
			Token:     "test-token",
			ExpiresAt: expiresAt,
		}

		referenceLinkEntity := NewReferenceLinkEntity(referenceLinkModel)

		assert.Equal(t, referenceLinkEntity.ID, refLinkID)
		assert.Equal(t, referenceLinkEntity.LinkID, linkID)
		assert.Equal(t, referenceLinkEntity.Type, RefLinkTypeInvite)
		assert.Equal(t, referenceLinkEntity.Token, "test-token")
		assert.Equal(t, referenceLinkEntity.ExpiresAt, expiresAt)
	})

	t.Run("create not expired entity from model", func(t *testing.T) {
		refLinkID := uuid.New()
		linkID := uuid.New()
		expiresAt := time.Now().Add(time.Hour)

		referenceLinkModel := models.ReferenceLink{
			ID:        refLinkID,
			LinkID:    linkID,
			LinkType:  RefLinkTypeLogin,
			Token:     "test-token",
			ExpiresAt: expiresAt,
		}

		referenceLinkEntity := NewReferenceLinkEntity(referenceLinkModel)

		assert.Equal(t, referenceLinkEntity.ID, refLinkID)
		assert.Equal(t, referenceLinkEntity.LinkID, linkID)
		assert.Equal(t, referenceLinkEntity.Type, RefLinkTypeLogin)
		assert.Equal(t, referenceLinkEntity.Token, "test-token")
		assert.Equal(t, referenceLinkEntity.ExpiresAt, expiresAt)
		assert.Equal(t, referenceLinkEntity.IsExpired(), false)
	})
}
