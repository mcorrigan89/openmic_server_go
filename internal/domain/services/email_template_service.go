package services

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/mcorrigan89/openmic/internal/common"
	"github.com/mcorrigan89/openmic/internal/domain/entities"
)

//go:embed "templates"
var templateFS embed.FS

type EmailTemplateService interface {
	LoginEmail(templateFile string, refLink *entities.ReferenceLinkEntity) (string, string, error)
}

type emailTemplateService struct {
	config *common.Config
}

func NewEmailTemplateService(cfg *common.Config) *emailTemplateService {
	return &emailTemplateService{
		config: cfg,
	}
}

func (s *emailTemplateService) LoginEmail(templateFile string, refLink *entities.ReferenceLinkEntity) (string, string, error) {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return "", "", err
	}

	plainBodyBytes := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBodyBytes, "plainBody", refLink)
	if err != nil {
		return "", "", err
	}

	htmlBodyBytes := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBodyBytes, "htmlBody", refLink)
	if err != nil {
		return "", "", err
	}

	plainBody := plainBodyBytes.String()
	htmlBody := htmlBodyBytes.String()

	return plainBody, htmlBody, nil
}
