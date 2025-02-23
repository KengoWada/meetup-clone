package main

import (
	"log"
	"net/http"
	"time"

	"github.com/KengoWada/meetup-clone/internal/services/auth"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Store
}

type config struct {
	addr        string
	debug       bool
	environment string
	frontendURL string
	dbConfig    dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route("/v1", func(r chi.Router) {
		authHandler := auth.NewHandler(app.store)
		authMux := authHandler.RegisterRoutes()
		r.Mount("/auth", authMux)
	})

	return mux
}

func (app *application) run(mux http.Handler) error {
	svr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server is starting on port %s", app.config.addr)
	return svr.ListenAndServe()
}
