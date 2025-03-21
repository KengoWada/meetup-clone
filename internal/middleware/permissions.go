package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/logger"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/go-chi/chi/v5"
)

func AuthenticatedRoute(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(internal.UserCtx).(*models.User)
		if !ok || user == nil {
			err := errors.New("no user in context")
			response.ErrorResponseUnauthorized(w, r, err)
			return
		}

		if user.IsDeactivated() || !user.IsActive {
			err := errors.New("user is deactivated or email not verified")
			response.ErrorResponseUnauthorized(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func IsStaffOrAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(internal.UserCtx).(*models.User)

		if user.Role == models.UserClientRole {
			err := errors.New("client role user tried to access staff or admin route")
			response.ErrorResponseForbidden(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func IsAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(internal.UserCtx).(*models.User)

		if user.Role != models.UserAdminRole {
			err := fmt.Errorf("%s role tried to access admin route", user.Role)
			response.ErrorResponseForbidden(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func HasOrgPermission(permissions []string, appStore store.Store, cacheStore cache.Store, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		orgID, err := strconv.ParseInt(chi.URLParam(r, "orgID"), 10, 64)
		if err != nil {
			errorMessage := response.ErrorResponse{Message: "Invalid organization ID"}
			response.ErrorResponseBadRequest(w, r, err, errorMessage)
			return
		}
		ctx := r.Context()
		user, _ := ctx.Value(internal.UserCtx).(*models.User)

		member, err := getOrganizationMember(ctx, r, appStore, cacheStore, user.UserProfile.ID, orgID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				err := errors.New("user is not an active member of the organization")
				response.ErrorResponseForbidden(w, r, err)
			default:
				response.ErrorResponseInternalServerErr(w, r, err)
			}
			return
		}

		role, err := getRole(ctx, r, appStore, cacheStore, member.RoleID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				err := errors.New("user has a deleted role")
				response.ErrorResponseForbidden(w, r, err)
			default:
				response.ErrorResponseInternalServerErr(w, r, err)
			}
			return
		}

		var hasPermission bool
		for _, permission := range permissions {
			if slices.Contains(role.Permissions, permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			err := errors.New("user does not have permissions to perform action")
			response.ErrorResponseForbidden(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func getOrganizationMember(ctx context.Context, r *http.Request, appStore store.Store, cacheStore cache.Store, userID, orgID int64) (*models.OrganizationMember, error) {
	var fields = []string{"user_id", "org_id"}
	var values = []any{userID, orgID}

	if !cfg.CacheConfig.Enabled {
		return appStore.OrganizationMembers.Get(ctx, false, fields, values)
	}

	member, err := cacheStore.OrganizationMembers.Get(userID, orgID)
	if err != nil {
		logger.ErrLoggerCache(r, err)
	}

	if member == nil {
		member, err = appStore.OrganizationMembers.Get(ctx, false, fields, values)
		if err != nil {
			return nil, err
		}

		if err := cacheStore.OrganizationMembers.Set(member); err != nil {
			logger.ErrLoggerCache(r, err)
		}
	}

	// This is to ensure that cached members that were deleted
	// don't have access
	if member.DeletedAt != nil {
		return nil, store.ErrNotFound
	}

	return member, nil
}

func getRole(ctx context.Context, r *http.Request, appStore store.Store, cacheStore cache.Store, roleID int64) (*models.Role, error) {
	var fields = []string{"id"}
	var values = []any{roleID}

	if !cfg.CacheConfig.Enabled {
		return appStore.Roles.Get(ctx, false, fields, values)
	}

	role, err := cacheStore.Roles.Get(roleID)
	if err != nil {
		logger.ErrLoggerCache(r, err)
	}

	if role == nil {
		role, err = appStore.Roles.Get(ctx, false, fields, values)
		if err != nil {
			return nil, err
		}

		if err := cacheStore.Roles.Set(role); err != nil {
			logger.ErrLoggerCache(r, err)
		}
	}

	// This is to ensure that cached roles that were deleted
	// don't have access
	if role.DeletedAt != nil {
		return nil, store.ErrNotFound
	}

	return role, nil
}
