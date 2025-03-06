package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(jwtAuthenticator auth.Authenticator, appStore store.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				err := fmt.Errorf("malformed token: %s", authHeader)
				response.ErrorResponseUnauthorized(w, r, err)
				return
			}

			token := headerParts[1]
			jwtToken, err := jwtAuthenticator.ValidateToken(token)
			if err != nil {
				err := fmt.Errorf("%s, token failed validation: %s", err.Error(), authHeader)
				response.ErrorResponseUnauthorized(w, r, err)
				return
			}

			claims, _ := jwtToken.Claims.(jwt.MapClaims)

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				response.ErrorResponseUnauthorized(w, r, err)
				return
			}

			ctx := r.Context()
			user, err := appStore.Users.GetByID(ctx, int(userID))
			if err != nil {
				switch err {
				case store.ErrNotFound:
					response.ErrorResponseUnauthorized(w, r, err)
				default:
					response.ErrorResponseInternalServerErr(w, r, err)
				}
				return
			}

			if user.IsDeactivated() || !user.IsActivated() {
				response.ErrorResponseUnauthorized(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, internal.UserCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

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
