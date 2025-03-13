package organizations

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/organizations/roles"
	"github.com/KengoWada/meetup-clone/internal/services/response"
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
			middleware.HasOrgPermission(internal.OrgUpdate, h.store, h.cacheStore, h.updateOrganization),
		)
		orgMux.Delete(
			"/",
			middleware.HasOrgPermission(internal.OrgDelete, h.store, h.cacheStore, h.deleteOrganization),
		)

		rolesHandler := roles.NewHandler(h.store, h.cacheStore)
		rolesMux := rolesHandler.RegisterRoutes()
		orgMux.Mount("/roles", rolesMux)
	})

	return mux
}

func getOrganization(appStore store.Store, cacheStore cache.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			orgID, err := strconv.ParseInt(chi.URLParam(r, "orgID"), 10, 64)
			if err != nil {
				errorMessage := response.ErrorResponse{Message: "Invalid organization ID"}
				response.ErrorResponseBadRequest(w, r, err, errorMessage)
				return
			}

			var organization *models.Organization
			if cfg.CacheConfig.Enabled {
				organization, err = cacheStore.Organizations.Get(orgID)
				if err != nil {
					logger.ErrLoggerCache(r, err)
				}
			}

			ctx := r.Context()

			if organization == nil {
				organization, err = appStore.Organizations.Get(ctx, false, []string{"id"}, []any{orgID})
				if err != nil {
					switch err {
					case store.ErrNotFound:
						errorMessage := response.ErrorResponse{Message: "Invalid organization ID"}
						response.ErrorResponseBadRequest(w, r, err, errorMessage)
						return
					default:
						response.ErrorResponseInternalServerErr(w, r, err)
					}
					return
				}

				if cfg.CacheConfig.Enabled {
					if err := cacheStore.Organizations.Set(organization); err != nil {
						logger.ErrLoggerCache(r, err)
					}
				}
			}

			if organization.DeletedAt != nil || !organization.IsActive {
				err := errors.New("organization was either deleted or deactivated")
				response.ErrorResponseUnauthorized(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, internal.OrgCtx, organization)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
