package main

import (
	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/logger"
)

func main() {
	log := logger.Get()
	appItems, err := app.NewApplication()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create application")
	}

	defer appItems.DB.Close()

	mux := appItems.App.Mount()
	log.Fatal().Err(appItems.App.Run(mux)).Msg("Server has stopped")
}
