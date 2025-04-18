package main

import (
	"log"

	"github.com/kirillidk/pvz-service/internal/app"
	"github.com/kirillidk/pvz-service/internal/config"
)

func main() {
	cfg := config.NewConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
