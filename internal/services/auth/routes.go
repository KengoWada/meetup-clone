package auth

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/register", h.registerUser)

	return mux
}
