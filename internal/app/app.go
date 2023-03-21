package app

import (
	"github.com/jqno/balGPT/internal/config"
	"github.com/jqno/balGPT/internal/database"
)

type App struct {
	Config *config.Config
	DB     *database.DB
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString)
	return &App{
		Config: cfg,
		DB:     db,
	}
}

func (a *App) Run() {
	// TODO: Initialize and run your web server here
}
