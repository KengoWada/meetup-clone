package app

import (
	"log"
	"net/http"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route("/v1", func(r chi.Router) {
		authHandler := auth.NewHandler(app.Store)
		authMux := authHandler.RegisterRoutes()
		r.Mount("/auth", authMux)
	})

	return mux
}

func (app *Application) Run(mux http.Handler) error {
	svr := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server is starting on port %s", app.Config.Addr)
	return svr.ListenAndServe()
}
