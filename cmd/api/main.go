package main

import (
	"log"

	"github.com/KengoWada/meetup-clone/internal/db"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/utils"
)

func main() {
	cfg := config{
		addr:        utils.GetString("SERVER_ADDR", ""),
		debug:       utils.GetBool("DEBUG", false),
		environment: utils.GetString("SERVER_ENVIRONMENT", "prod"),
		frontendURL: utils.GetString("FRONTEND_URL", ""),
		dbConfig: dbConfig{
			addr:         utils.GetString("DB_ADDR", ""),
			maxOpenConns: utils.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: utils.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  utils.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.dbConfig.addr,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("connected to db")

	store := store.NewStore(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
