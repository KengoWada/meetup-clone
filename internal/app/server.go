package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/KengoWada/meetup-clone/docs"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/logger"
	appMiddleware "github.com/KengoWada/meetup-clone/internal/middleware"
	"github.com/KengoWada/meetup-clone/internal/services/auth"
	"github.com/KengoWada/meetup-clone/internal/services/organizations"
	"github.com/KengoWada/meetup-clone/internal/services/profiles"
	"github.com/KengoWada/meetup-clone/internal/services/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const version = "0.0.1"

var l = logger.Get()

func (app *Application) Mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(appMiddleware.LoggerMiddleware)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(appMiddleware.JWTMiddleware(app.Authenticator, app.Store, app.CacheStore))

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("invalid route path")
		response.ErrorResponseRouteNotFound(w, r, err)
	})

	mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("invalid route method")
		response.ErrorResponseRouteMethodNotAllowed(w, r, err)
	})

	mux.Route("/v1", func(r chi.Router) {
		if app.Config.Environment == config.AppEnvDev {
			docsURL := fmt.Sprintf("%s/swagger/doc.json", app.Config.Addr)
			r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		}

		authHandler := auth.NewHandler(app.Store, app.CacheStore, app.Authenticator)
		authMux := authHandler.RegisterRoutes()
		r.Mount("/auth", authMux)

		profileHandler := profiles.NewHandler(app.Store, app.Authenticator)
		profileMux := profileHandler.RegisterRoutes()
		r.Mount("/profiles", profileMux)

		organizationHandler := organizations.NewHandler(app.Store, app.CacheStore)
		organizationMux := organizationHandler.RegisterRoutes()
		r.Mount("/organizations", organizationMux)
	})

	return mux
}

func (app *Application) Run(mux http.Handler) error {
	if app.Config.Environment == config.AppEnvDev {
		docs.SwaggerInfo.Version = version
		docs.SwaggerInfo.Host = app.Config.ApiURL
		docs.SwaggerInfo.BasePath = "/v1"
	}

	svr := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	l.Info().Msgf("%s_env:server is starting on port %s", app.Config.Environment, app.Config.Addr)
	return svr.ListenAndServe()
}
