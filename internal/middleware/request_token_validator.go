package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/golang-jwt/jwt/v5"
)

var cfg = config.Get()

func JWTMiddleware(jwtAuthenticator auth.Authenticator, appStore store.Store, cacheStore cache.Store) func(next http.Handler) http.Handler {
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
			user, err := getUser(ctx, userID, appStore, cacheStore)
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

func getUser(ctx context.Context, ID int64, appStore store.Store, cacheStore cache.Store) (*models.User, error) {
	if !cfg.CacheConfig.Enabled {
		fields, values := []string{"id"}, []any{ID}
		return appStore.Users.GetWithProfile(ctx, false, fields, values)
	}

	user, err := cacheStore.Users.Get(ID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	fields, values := []string{"id"}, []any{ID}
	user, err = appStore.Users.GetWithProfile(ctx, false, fields, values)
	if err != nil {
		return nil, err
	}

	if err := cacheStore.Users.Set(user); err != nil {
		return nil, err
	}

	return user, nil
}
