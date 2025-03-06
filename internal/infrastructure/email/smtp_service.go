package email

import (
	"context"

	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
	"github.com/wneessen/go-mail"
)

type smtpService struct {
	config *common.Config
}

func NewSmtpService(cfg *common.Config) *smtpService {
	return &smtpService{
		config: cfg,
	}
}

func (service *smtpService) SendEmail(ctx context.Context, email *entities.EmailEntity) error {
	msg := mail.NewMsg()
	msg.Subject(email.Subject)
	msg.SetBodyString(mail.TypeTextPlain, email.PlainBody)
	msg.AddAlternativeString(mail.TypeTextHTML, email.HtmlBody)

	err := service.addEmails(msg, email.ToEmail, email.FromEmail)
	if err != nil {
		return err
	}

	client, err := service.createClient()
	if err != nil {
		return err
	}

	err = client.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

func (service *smtpService) createClient() (*mail.Client, error) {
	smtpServer := service.config.Mail.SMTPServer
	smtpUsername := service.config.Mail.SMTPUsername
	smtpPassword := service.config.Mail.SMTPPassword

	client, err := mail.NewClient(smtpServer, mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory), mail.WithUsername(smtpUsername), mail.WithPassword(smtpPassword))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (service *smtpService) addEmails(msg *mail.Msg, toEmail, fromEmail string) error {
	err := msg.From(fromEmail)
	if err != nil {
		return err
	}

	err = msg.To(toEmail)
	if err != nil {
		return err
	}

	return nil
}
