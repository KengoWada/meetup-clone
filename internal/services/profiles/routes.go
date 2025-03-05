package profiles

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store         store.Store
	authenticator auth.Authenticator
}

func NewHandler(store store.Store, authenticator auth.Authenticator) *Handler {
	return &Handler{store, authenticator}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticatedRoute)

		r.Get("/", h.getPersonalProfile)
		r.Put("/", h.updateUserProfile)
	})

	return mux
}
