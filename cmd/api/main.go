package main

import (
	"log"

	"github.com/kirillidk/pvz-service/internal/app"
	"github.com/kirillidk/pvz-service/internal/config"
)

func main() {
	cfg := config.NewConfig()
	app := app.NewApp(cfg)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
