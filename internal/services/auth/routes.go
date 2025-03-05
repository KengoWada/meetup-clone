package auth

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-chi/chi/v5"
)

var cfg = config.Get()

type Handler struct {
	store         store.Store
	authenticator auth.Authenticator
}

func NewHandler(store store.Store, authenticator auth.Authenticator) *Handler {
	return &Handler{store, authenticator}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/register", h.registerUser)
	mux.Post("/login", h.loginUser)
	mux.Patch("/activate", h.activateUser)
	mux.Post("/resend-verification-email", h.resendVerificationEmail)
	mux.Post("/password-reset-request", h.passwordResetRequest)
	mux.Post("/reset-password", h.resetUserPassword)

	return mux
}
