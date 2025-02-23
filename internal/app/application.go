package app

import (
	"database/sql"
	"slices"

	"github.com/KengoWada/meetup-clone/internal/db"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

var (
	environments = []string{"dev", "prod"}
)

type Application struct {
	Config Config
	Store  store.Store
}

type DeferItems struct {
	DB *sql.DB
}

func NewApplication() (*Application, DeferItems, error) {
	environment := utils.GetString("SERVER_ENVIRONMENT", "prod")
	if !slices.Contains(environments, environment) {
		environment = "prod"
	}

	cfg := Config{
		Addr:        utils.GetString("SERVER_ADDR", ""),
		Debug:       utils.GetBool("DEBUG", false),
		Environment: environment,
		FrontendURL: utils.GetString("FRONTEND_URL", ""),
		DBConfig: DBConfig{
			Addr:         utils.GetString("DB_ADDR", ""),
			MaxOpenConns: utils.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: utils.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  utils.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.New(
		cfg.DBConfig.Addr,
		cfg.DBConfig.MaxOpenConns,
		cfg.DBConfig.MaxIdleConns,
		cfg.DBConfig.MaxIdleTime,
	)
	if err != nil {
		return nil, DeferItems{}, err
	}

	store := store.NewStore(db)

	deferItems := DeferItems{
		DB: db,
	}

	app := &Application{
		Config: cfg,
		Store:  store,
	}
	return app, deferItems, nil
}

func NewTestApplication(store store.Store) *Application {
	cfg := Config{
		Addr:        utils.GetString("SERVER_ADDR", ""),
		Debug:       utils.GetBool("DEBUG", false),
		Environment: "test",
		FrontendURL: utils.GetString("FRONTEND_URL", ""),
	}

	app := &Application{
		Config: cfg,
		Store:  store,
	}
	return app
}
