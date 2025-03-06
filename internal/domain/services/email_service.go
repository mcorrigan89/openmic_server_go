package services

import (
	"context"

	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/mcorrigan89/openmic/internal/domain/external"
	"github.com/mcorrigan89/openmic/internal/infrastructure/postgres/models"
)

type EmailService interface {
	SendEmail(ctx context.Context, querier models.Querier, email *entities.EmailEntity) (*entities.EmailEntity, error)
}

type emailService struct {
	smptService external.SmtpService
}

func NewEmailService(smptService external.SmtpService) *emailService {
	return &emailService{smptService: smptService}
}

func (s *emailService) SendEmail(ctx context.Context, querier models.Querier, email *entities.EmailEntity) (*entities.EmailEntity, error) {
	err := s.smptService.SendEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return email, nil
}
