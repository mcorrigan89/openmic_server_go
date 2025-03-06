package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mcorrigan89/openmic/internal/application"
	"github.com/mcorrigan89/openmic/internal/application/commands"
	"github.com/mcorrigan89/openmic/internal/application/queries"
	"github.com/mcorrigan89/openmic/internal/infrastructure/media"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/middleware"
	"github.com/rs/zerolog"
)

type ImageHandler struct {
	logger          *zerolog.Logger
	imageAppService application.ImageApplicationService
}

func NewImageHandler(logger *zerolog.Logger, imageAppService application.ImageApplicationService) *ImageHandler {
	return &ImageHandler{
		logger:          logger,
		imageAppService: imageAppService,
	}
}

func (h *ImageHandler) GetImageByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	imageID := r.PathValue("id")

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		http.Error(w, "Failed to parse UUID", http.StatusInternalServerError)
		return
	}

	query := queries.ImageByIDQuery{
		ID: imageUUID,
	}

	image, err := h.imageAppService.GetImageByID(ctx, query)
	if err != nil {
		http.Error(w, "Failed to get user by ID", http.StatusInternalServerError)
		return
	}

	imageDto := dto.NewImageDtoFromEntity(image)

	imageJson, err := imageDto.ToJson()
	if err != nil {
		http.Error(w, "Failed to marshal image to JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(imageJson)
}

func (h *ImageHandler) GetImageDataByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	imageID := r.PathValue("id")

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		http.Error(w, "Failed to parse UUID", http.StatusInternalServerError)
		return
	}

	query := queries.ImageDataByIDQuery{
		ID:        imageUUID,
		Rendition: media.RenditionMedium,
	}

	image, contentType, err := h.imageAppService.GetImageDataByID(ctx, query)
	if err != nil {
		h.logger.Err(err).Ctx(ctx).Msg("Failed to get image")
		http.Error(w, "Failed to get image", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", contentType)
	w.Write(image)
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		h.logger.Err(err).Ctx(ctx).Msg("Failed to get file from form")
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userContextEntity := middleware.GetUserFromContext(ctx)
	if userContextEntity == nil {
		h.logger.Error().Ctx(ctx).Msg("User is not authenticated")
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}

	fileName := fmt.Sprintf("%s-%s", handler.Filename, uuid.New().String())

	cmd := commands.CreateNewAvatarImageCommand{
		UserID:   userContextEntity.UserID,
		BucketID: "image",
		ObjectID: fileName,
		File:     file,
		Size:     handler.Size,
	}

	image, err := h.imageAppService.UploadAvatarImage(ctx, cmd)
	if err != nil {
		http.Error(w, "Failed to get user by ID", http.StatusInternalServerError)
		return
	}

	imageDto := dto.NewImageDtoFromEntity(image)

	imageJson, err := imageDto.ToJson()
	if err != nil {
		http.Error(w, "Failed to marshal image to JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(imageJson)
}
