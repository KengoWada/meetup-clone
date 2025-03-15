package roles

import (
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/middleware"
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

	mux.Get(
		"/",
		middleware.HasOrgPermission(
			[]string{internal.RoleUpdate, internal.RoleDelete, internal.MemberAdd, internal.MemberRoleUpdate},
			h.store,
			h.cacheStore,
			h.getOrganizationRoles,
		),
	)
	mux.Post(
		"/",
		middleware.HasOrgPermission(
			[]string{internal.RoleCreate},
			h.store,
			h.cacheStore,
			h.createRole,
		),
	)
	mux.Get(
		"/permissions",
		middleware.HasOrgPermission(
			[]string{internal.RoleCreate, internal.RoleUpdate},
			h.store,
			h.cacheStore,
			h.getPermissions,
		),
	)

	mux.Route("/{roleID}", func(roleMux chi.Router) {
		roleMux.Use(getRole(h.store, h.cacheStore))

		roleMux.Get(
			"/",
			middleware.HasOrgPermission(
				[]string{internal.RoleCreate, internal.RoleUpdate, internal.RoleDelete},
				h.store,
				h.cacheStore,
				h.getRole,
			),
		)
	})

	return mux
}
