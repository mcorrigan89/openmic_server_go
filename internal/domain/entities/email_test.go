package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewEmailEntity(t *testing.T) {

	t.Run("create entity from model", func(t *testing.T) {
		emailID := uuid.New()
		toEmail := "to@gmail.com"
		fromEmail := "from@gmail.com"
		subject := "test subject"
		body := "test body"

		emailEntity := NewEmailEntity(EmailEntityArgs{
			ID:        emailID,
			ToEmail:   toEmail,
			FromEmail: fromEmail,
			Subject:   subject,
			PlainBody: body,
		})

		assert.Equal(t, emailEntity.ID, emailID)
		assert.Equal(t, emailEntity.ToEmail, toEmail)
		assert.Equal(t, emailEntity.FromEmail, fromEmail)
		assert.Equal(t, emailEntity.Subject, subject)
		assert.Equal(t, emailEntity.PlainBody, body)
	})
}
