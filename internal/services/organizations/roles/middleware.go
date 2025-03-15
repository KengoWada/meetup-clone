package roles

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/go-chi/chi/v5"
)

func getRole(appStore store.Store, cacheStore cache.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			roleID, err := strconv.ParseInt(chi.URLParam(r, "roleID"), 10, 64)
			invalidIDErrorMessage := response.ErrorResponse{Message: "Invalid role ID"}
			if err != nil {
				response.ErrorResponseBadRequest(w, r, err, invalidIDErrorMessage)
				return
			}

			var role *models.Role
			if cfg.CacheConfig.Enabled {
				role, err = cacheStore.Roles.Get(roleID)
				if err != nil {
					logger.ErrLoggerCache(r, err)
				}
			}

			ctx := r.Context()

			if role == nil {
				role, err = appStore.Roles.Get(ctx, false, []string{"id"}, []any{roleID})
				if err != nil {
					switch err {
					case store.ErrNotFound:
						response.ErrorResponseForbidden(w, r, err)
					default:
						response.ErrorResponseInternalServerErr(w, r, err)
					}
					return
				}

				if cfg.CacheConfig.Enabled {
					if err := cacheStore.Roles.Set(role); err != nil {
						logger.ErrLoggerCache(r, err)
					}
				}
			}

			if role.DeletedAt != nil {
				response.ErrorResponseForbidden(w, r, store.ErrNotFound)
				return
			}

			ctx = context.WithValue(ctx, internal.RoleCtx, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
