package app

import (
	"database/sql"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/db"
	"github.com/KengoWada/meetup-clone/internal/store"
)

type Application struct {
	Config        config.Config
	Store         store.Store
	Authenticator auth.Authenticator
}

type AppItems struct {
	App *Application
	DB  *sql.DB
}

func NewApplication() (*AppItems, error) {
	cfg := config.NewConfig()

	appItems := &AppItems{}

	// Make database connection
	db, err := db.New(
		cfg.DBConfig.Addr,
		cfg.DBConfig.MaxOpenConns,
		cfg.DBConfig.MaxIdleConns,
		cfg.DBConfig.MaxIdleTime,
		cfg.Environment,
	)
	if err != nil {
		return appItems, err
	}
	l.Info().Msg("successfully connected to postgres")
	appItems.DB = db

	// Create JWT Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.AuthConfig.Secret, cfg.AuthConfig.Audience, cfg.AuthConfig.Issuer)

	// Create Global App Store
	store := store.NewStore(db)

	app := &Application{
		Config:        cfg,
		Store:         store,
		Authenticator: jwtAuthenticator,
	}
	appItems.App = app

	return appItems, nil
}
