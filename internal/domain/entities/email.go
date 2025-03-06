package entities

import (
	"github.com/google/uuid"
)

type EmailEntity struct {
	ID        uuid.UUID
	ToEmail   string
	FromEmail string
	Subject   string
	PlainBody string
	HtmlBody  string
}

type EmailEntityArgs struct {
	ID        uuid.UUID
	ToEmail   string
	FromEmail string
	Subject   string
	PlainBody string
	HtmlBody  string
}

func NewEmailEntity(args EmailEntityArgs) *EmailEntity {
	return &EmailEntity{
		ID:        args.ID,
		ToEmail:   args.ToEmail,
		FromEmail: args.FromEmail,
		Subject:   args.Subject,
		PlainBody: args.PlainBody,
		HtmlBody:  args.HtmlBody,
	}
}
