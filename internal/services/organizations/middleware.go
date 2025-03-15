package organizations

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

func getOrganization(appStore store.Store, cacheStore cache.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			getOrgDetails := func(orgID int64) (org *models.Organization, err error) {
				if cfg.CacheConfig.Enabled {
					org, err = cacheStore.Organizations.Get(orgID)
					if err != nil {
						logger.ErrLoggerCache(r, err)
					}
				}

				if org == nil {
					fields := []string{"id", "is_active"}
					values := []any{orgID, true}
					org, err = appStore.Organizations.Get(ctx, false, fields, values)
					if err != nil {
						return nil, err
					}

					if cfg.CacheConfig.Enabled {
						if err := cacheStore.Organizations.Set(org); err != nil {
							logger.ErrLoggerCache(r, err)
						}
					}
				}

				if org.DeletedAt != nil || !org.IsActive {
					return nil, store.ErrNotFound
				}

				return org, nil
			}

			orgID, err := strconv.ParseInt(chi.URLParam(r, "orgID"), 10, 64)
			if err != nil {
				errorMessage := response.ErrorResponse{Message: "Invalid organization ID"}
				response.ErrorResponseBadRequest(w, r, err, errorMessage)
				return
			}

			organization, err := getOrgDetails(orgID)
			if err != nil {
				switch err {
				case store.ErrNotFound:
					response.ErrorResponseForbidden(w, r, err)
				default:
					response.ErrorResponseInternalServerErr(w, r, err)
				}
				return
			}

			ctx = context.WithValue(ctx, internal.OrgCtx, organization)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
