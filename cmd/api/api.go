package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type application struct {
	config config
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
