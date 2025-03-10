package auth

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/go-chi/chi/v5"
)

var cfg = config.Get()

type Handler struct {
	store         store.Store
	cacheStore    cache.Store
	authenticator auth.Authenticator
}

func NewHandler(store store.Store, cacheStore cache.Store, authenticator auth.Authenticator) *Handler {
	return &Handler{store, cacheStore, authenticator}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticatedRoute)
		r.Use(middleware.IsStaffOrAdmin)

		r.Patch("/users/{userID}/deactivate", h.deactivateUser)
	})

	mux.Post("/register", h.registerUser)
	mux.Post("/login", h.loginUser)
	mux.Patch("/activate", h.activateUser)
	mux.Post("/resend-verification-email", h.resendVerificationEmail)
	mux.Post("/password-reset-request", h.passwordResetRequest)
	mux.Post("/reset-password", h.resetUserPassword)

	return mux
}
