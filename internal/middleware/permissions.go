package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
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
			response.ErrorResponseUnauthorized(w, r, err)
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
			response.ErrorResponseUnauthorized(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
