package main

import (
	"github.com/jqno/balGPT/internal/app"
	"github.com/jqno/balGPT/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := app.NewApp(cfg)
	app.Run()
}
