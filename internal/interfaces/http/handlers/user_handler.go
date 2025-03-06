package handlers

import (
	"context"
	"net/mail"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/application"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	logger         *zerolog.Logger
	userAppService application.UserApplicationService
}

func NewUserHandler(logger *zerolog.Logger, userAppService application.UserApplicationService) *UserHandler {
	return &UserHandler{
		logger:         logger,
		userAppService: userAppService,
	}
}

func (h *UserHandler) GetUserByID(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*dto.GetUserByIDResponse, error) {

	query := queries.UserByIDQuery{
		ID: input.ID,
	}

	user, err := h.userAppService.GetUserByID(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by ID", err)
	}

	userDto := dto.NewUserDtoFromEntity(user)

	return &dto.GetUserByIDResponse{
		Body: userDto,
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, input *struct {
	Email string `path:"email"`
}) (*dto.GetUserByEmailResponse, error) {

	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to parse email", err)
	}

	query := queries.UserByEmailQuery{
		Email: input.Email,
	}

	user, err := h.userAppService.GetUserByEmail(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by email", err)
	}

	userDto := dto.NewUserDtoFromEntity(user)

	return &dto.GetUserByEmailResponse{
		Body: userDto,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, input *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	cmd := commands.CreateNewUserCommand{
		Email:      input.Body.Email,
		GivenName:  input.Body.GivenName,
		FamilyName: input.Body.FamilyName,
	}

	userSessionEntity, err := h.userAppService.CreateUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create user", err)
	}

	userDto := dto.NewUserDtoFromEntity(userSessionEntity.User)
	sessionDto := dto.SessionDto{
		Token:     userSessionEntity.SessionToken,
		ExpiresAt: userSessionEntity.ExpiresAt(),
	}

	resp := dto.CreateUserResponse{}

	resp.Body.User = userDto
	resp.Body.SessionDto = &sessionDto

	return &resp, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	cmd := commands.UpdateUserCommand{
		ID:         input.ID,
		Email:      input.Body.Email,
		GivenName:  input.Body.GivenName,
		FamilyName: input.Body.FamilyName,
		Handle:     input.Body.Handle,
	}

	userEntity, err := h.userAppService.UpdateUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update user", err)
	}

	userDto := dto.NewUserDtoFromEntity(userEntity)

	resp := dto.UpdateUserResponse{}

	resp.Body.User = userDto

	return &resp, nil
}
