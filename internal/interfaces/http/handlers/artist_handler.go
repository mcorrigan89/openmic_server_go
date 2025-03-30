package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/application"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/rs/zerolog"
)

type ArtistHandler struct {
	logger           *zerolog.Logger
	artistAppService application.ArtistApplicationService
}

func NewArtistHandler(logger *zerolog.Logger, artistAppService application.ArtistApplicationService) *ArtistHandler {
	return &ArtistHandler{
		logger:           logger,
		artistAppService: artistAppService,
	}
}

func (h *ArtistHandler) GetArtistByID(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*dto.GetArtistByIDResponse, error) {

	query := queries.ArtistByIDQuery{
		ID: input.ID,
	}

	artist, err := h.artistAppService.GetArtistByID(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get artist by ID", err)
	}

	artistDto := dto.NewArtistDtoFromEntity(artist)

	return &dto.GetArtistByIDResponse{
		Body: artistDto,
	}, nil
}

func (h *ArtistHandler) GetArtistsByTitle(ctx context.Context, input *struct {
	Title string `query:"title"`
}) (*dto.GetArtistsByTitleResponse, error) {

	query := queries.ArtistsByTitleQuery{
		Title: input.Title,
	}

	artists, err := h.artistAppService.GetArtistsByTitle(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get artist by title", err)
	}

	artistDtos := make([]*dto.ArtistDto, 0, len(artists))
	for _, artist := range artists {
		artistDtos = append(artistDtos, dto.NewArtistDtoFromEntity(artist))
	}

	return &dto.GetArtistsByTitleResponse{
		Body: artistDtos,
	}, nil
}

func (h *ArtistHandler) GetAllArtists(ctx context.Context, input *struct {
}) (*dto.GetAllArtistsResponse, error) {

	artists, err := h.artistAppService.GetAllArtists(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get all artists", err)
	}

	artistDtos := make([]*dto.ArtistDto, 0, len(artists))
	for _, artist := range artists {
		artistDtos = append(artistDtos, dto.NewArtistDtoFromEntity(artist))
	}

	return &dto.GetAllArtistsResponse{
		Body: artistDtos,
	}, nil
}

func (h *ArtistHandler) CreateArtist(ctx context.Context, input *dto.CreateArtistRequest) (*dto.CreateArtistResponse, error) {

	cmd := commands.CreateNewArtistCommand{
		Title:    input.Body.Title,
		SubTitle: input.Body.SubTitle,
		Bio:      input.Body.Bio,
	}

	artist, err := h.artistAppService.CreateArtist(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create artist", err)
	}

	artistDto := dto.NewArtistDtoFromEntity(artist)

	return &dto.CreateArtistResponse{
		Body: artistDto,
	}, nil
}

func (h *ArtistHandler) UpdateArtist(ctx context.Context, input *dto.UpdateArtistRequest) (*dto.UpdateArtistResponse, error) {

	cmd := commands.UpdateArtistCommand{
		ID:       input.ID,
		Title:    input.Body.Title,
		SubTitle: input.Body.SubTitle,
		Bio:      input.Body.Bio,
	}

	artist, err := h.artistAppService.UpdateArtist(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update artist", err)
	}

	artistDto := dto.NewArtistDtoFromEntity(artist)

	return &dto.UpdateArtistResponse{
		Body: artistDto,
	}, nil
}

func (h *ArtistHandler) DeleteArtist(ctx context.Context, input *dto.DeleteArtistRequest) (*dto.DeleteArtistResponse, error) {

	cmd := commands.DeleteArtistCommand{
		ID: input.ID,
	}

	err := h.artistAppService.DeleteArtist(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete artist", err)
	}

	return &dto.DeleteArtistResponse{}, nil
}
