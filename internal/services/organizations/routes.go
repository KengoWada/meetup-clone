package organizations

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store      store.Store
	cacheStore cache.Store
}

func NewHandler(store store.Store, cacheStore cache.Store) *Handler {
	return &Handler{store, cacheStore}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.AuthenticatedRoute)

	mux.Post("/", h.createOrganization)

	return mux
}
