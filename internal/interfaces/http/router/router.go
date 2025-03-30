package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/dto"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/handlers"
	"github.com/mcorrigan89/openmic/internal/interfaces/http/middleware"
)

func NewRouter(mux *http.ServeMux, middleware middleware.Middleware, userHandler *handlers.UserHandler, imageHandler *handlers.ImageHandler, eventHandler *handlers.EventHandler, artistHandler *handlers.ArtistHandler) http.Handler {

	api := humago.New(mux, huma.DefaultConfig("OpenMic API", "1.0.0"))

	// User routes
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"User"},
	}, userHandler.GetUserByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-email",
		Method:      http.MethodGet,
		Path:        "/user/email/{email}",
		Summary:     "Get a user by email",
		Tags:        []string{"User"},
	}, userHandler.GetUserByEmail)

	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/user",
		Summary:     "Create user",
		Tags:        []string{"User"},
	}, userHandler.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/user/{id}",
		Summary:     "Update user",
		Tags:        []string{"User"},
	}, userHandler.UpdateUser)

	// Image routes
	mux.HandleFunc("POST /image/upload", imageHandler.UploadImage)
	mux.HandleFunc("GET /image/{id}/metadata", middleware.Authorization(imageHandler.GetImageByID))
	mux.HandleFunc("GET /image/{id}", imageHandler.GetImageDataByID)

	// Event routes
	huma.Register(api, huma.Operation{
		OperationID: "get-event",
		Method:      http.MethodGet,
		Path:        "/event/{id}",
		Summary:     "Event by ID",
		Tags:        []string{"Event"},
	}, eventHandler.GetEventByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-current-event",
		Method:      http.MethodGet,
		Path:        "/event/now",
		Summary:     "Current Event",
		Tags:        []string{"Event"},
	}, eventHandler.GetCurrentEvent)

	huma.Register(api, huma.Operation{
		OperationID: "get-events",
		Method:      http.MethodGet,
		Path:        "/events",
		Summary:     "Upcoming Events",
		Tags:        []string{"Event"},
	}, eventHandler.GetUpcomingEvents)

	huma.Register(api, huma.Operation{
		OperationID: "create-event",
		Method:      http.MethodPost,
		Path:        "/event",
		Summary:     "Create Event",
		Tags:        []string{"Event"},
	}, eventHandler.CreateEvent)

	huma.Register(api, huma.Operation{
		OperationID: "update-event",
		Method:      http.MethodPut,
		Path:        "/event/{id}",
		Summary:     "Update Event",
		Tags:        []string{"Event"},
	}, eventHandler.UpdateEvent)

	huma.Register(api, huma.Operation{
		OperationID: "delete-event",
		Method:      http.MethodDelete,
		Path:        "/event/{id}",
		Summary:     "Delete Event",
		Tags:        []string{"Event"},
	}, eventHandler.DeleteEvent)

	huma.Register(api, huma.Operation{
		OperationID: "add-artist-to-event",
		Method:      http.MethodPost,
		Path:        "/event/{event_id}/add",
		Summary:     "Add Artist to Event",
		Tags:        []string{"Event"},
	}, eventHandler.AddArtistToEvent)

	huma.Register(api, huma.Operation{
		OperationID: "remove-artist-from-event",
		Method:      http.MethodPost,
		Path:        "/event/{event_id}/remove",
		Summary:     "Remove Artist from Event",
		Tags:        []string{"Event"},
	}, eventHandler.RemoveArtistFromEvent)

	huma.Register(api, huma.Operation{
		OperationID: "set-timeslot",
		Method:      http.MethodPost,
		Path:        "/event/{event_id}/timeslot",
		Summary:     "Set Timeslot",
		Tags:        []string{"Event"},
	}, eventHandler.SetTimeslotMarker)

	huma.Register(api, huma.Operation{
		OperationID: "delete-timeslot",
		Method:      http.MethodDelete,
		Path:        "/event/{event_id}/timeslot",
		Summary:     "Delete Timeslot",
		Tags:        []string{"Event"},
	}, eventHandler.DeleteTimeslotMarker)

	huma.Register(api, huma.Operation{
		OperationID: "set-sort-order",
		Method:      http.MethodPut,
		Path:        "/event/{event_id}/sort",
		Summary:     "Set Sort Order",
		Tags:        []string{"Event"},
	}, eventHandler.SetSortOrderRequest)

	sse.Register(api, huma.Operation{
		OperationID: "sse",
		Method:      http.MethodGet,
		Path:        "/sse/{event_id}",
		Summary:     "Server sent events example",
		Tags:        []string{"Event"},
	}, map[string]any{
		// Mapping of event type name to Go struct for that event.
		"message": dto.ListenForChangeEventResponse{},
	}, eventHandler.ListenForEventChange)

	// Artist routes
	huma.Register(api, huma.Operation{
		OperationID: "get-artist",
		Method:      http.MethodGet,
		Path:        "/artist/{id}",
		Summary:     "Get Artist by ID",
		Tags:        []string{"Artist"},
	}, artistHandler.GetArtistByID)

	huma.Register(api, huma.Operation{
		OperationID: "get-artists-by-title",
		Method:      http.MethodGet,
		Path:        "/artists/search",
		Summary:     "Get Artists by Title",
		Tags:        []string{"Artist"},
	}, artistHandler.GetArtistsByTitle)

	huma.Register(api, huma.Operation{
		OperationID: "get-all-artists",
		Method:      http.MethodGet,
		Path:        "/artists",
		Summary:     "Get All Artists",
		Tags:        []string{"Artist"},
	}, artistHandler.GetAllArtists)

	huma.Register(api, huma.Operation{
		OperationID: "create-artist",
		Method:      http.MethodPost,
		Path:        "/artist",
		Summary:     "Create Artist",
		Tags:        []string{"Artist"},
	}, artistHandler.CreateArtist)

	huma.Register(api, huma.Operation{
		OperationID: "update-artist",
		Method:      http.MethodPut,
		Path:        "/artist/{id}",
		Summary:     "Update Artist",
		Tags:        []string{"Artist"},
	}, artistHandler.UpdateArtist)

	huma.Register(api, huma.Operation{
		OperationID: "delete-artist",
		Method:      http.MethodDelete,
		Path:        "/artist/{id}",
		Summary:     "Delete Artist",
		Tags:        []string{"Artist"},
	}, artistHandler.DeleteArtist)

	return middleware.RecoverPanic(middleware.EnabledCORS(middleware.ContextBuilder(mux)))
}
