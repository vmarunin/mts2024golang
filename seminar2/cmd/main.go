package main

import (
	"github.com/vmarunin/mts2024golang/seminar2/internal/app"
	"github.com/vmarunin/mts2024golang/seminar2/internal/pkg/config"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.InitFromEnv()
	if err != nil {
		panic(err)
	}

	mainApp := app.NewApp(cfg)
	err = mainApp.Run()
	if err != nil {
		panic(err)
	}
}
