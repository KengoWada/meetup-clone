package organizations

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/services/organizations/members"
	"github.com/KengoWada/meetup-clone/internal/services/organizations/roles"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/go-chi/chi/v5"
)

var cfg = config.Get()

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
	mux.Get("/", h.getUsersOrganizations)

	mux.Route("/{orgID}", func(orgMux chi.Router) {
		orgMux.Use(getOrganization(h.store, h.cacheStore))

		orgMux.Group(func(r chi.Router) {
			r.Use(middleware.IsStaffOrAdmin)
			r.Patch("/", h.deactivateOrganization)
		})

		orgMux.Get("/", h.getOrganization)
		orgMux.Put(
			"/",
			middleware.HasOrgPermission(
				[]string{internal.OrgUpdate},
				h.store, h.cacheStore,
				h.updateOrganization,
			),
		)
		orgMux.Delete(
			"/",
			middleware.HasOrgPermission(
				[]string{internal.OrgDelete},
				h.store,
				h.cacheStore,
				h.deleteOrganization,
			),
		)

		rolesHandler := roles.NewHandler(h.store, h.cacheStore)
		rolesMux := rolesHandler.RegisterRoutes()
		orgMux.Mount("/roles", rolesMux)

		membersHandler := members.NewHandler(h.store, h.cacheStore)
		membersMux := membersHandler.RegisterRoutes()
		orgMux.Mount("/members", membersMux)
	})

	return mux
}
