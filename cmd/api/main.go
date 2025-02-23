package main

import (
	"log"

	"github.com/KengoWada/meetup-clone/internal/app"
)

func main() {
	app, deferItems, err := app.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	defer deferItems.DB.Close()

	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
