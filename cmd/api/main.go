package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/cyber_bed/internal/app"
	"github.com/cyber_bed/internal/config"
)

func main() {
	var configPath string
	config.ParseFlag(&configPath)
	flag.Parse()

	cfg := config.New()
	if err := cfg.Open(configPath); err != nil {
		log.Print("Failed to open config file")
	}

	if err := cfg.Open(cfg.EnvFile); err != nil {
		log.Fatalf("Error loading %s file", cfg.EnvFile)
	}
	// if err := godotenv.Load(cfg.EnvFile); err != nil {
	// 	log.Fatal("Error loading %s file", cfg.EnvFile)
	// }

	e := echo.New()
	app := app.New(e, cfg)
	if err := app.Start(); err != nil {
		app.Echo.Logger.Error(err)
	}
}
