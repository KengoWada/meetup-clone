package main

import (
	_ "github.com/KengoWada/meetup-clone/docs"
	"github.com/KengoWada/meetup-clone/internal/app"
	"github.com/KengoWada/meetup-clone/internal/logger"
)

//	@title			MeetUp Clone API
//	@description	API for MeetUp Clone
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	MIT License
//	@license.url	https://opensource.org/license/mit

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
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
