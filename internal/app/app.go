package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jqno/balGPT/internal/config"
	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/predictor"
	"github.com/jqno/balGPT/internal/scraper"
)

type App struct {
	Config    *config.Config
	DB        *database.DB
	Scraper   *scraper.ScrapeData
	Predictor *predictor.Predictor
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString)
	scraper := scraper.NewScrapeData(db)
	predictor := predictor.NewPredictor(db)
	return &App{
		Config:    cfg,
		DB:        db,
		Scraper:   scraper,
		Predictor: predictor,
	}
}

func (a *App) Run() {
	http.HandleFunc("/predict", handlePrediction(a.Scraper, a.Predictor))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func handlePrediction(s *scraper.ScrapeData, p *predictor.Predictor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeTeamIDStr := r.URL.Query().Get("home_team_id")
		awayTeamIDStr := r.URL.Query().Get("away_team_id")

		if homeTeamIDStr == "" || awayTeamIDStr == "" {
			http.Error(w, "Both home_team_id and away_team_id are required.", http.StatusBadRequest)
			return
		}

		homeTeamID, err := strconv.Atoi(homeTeamIDStr)
		if err != nil {
			http.Error(w, "Invalid home_team_id.", http.StatusBadRequest)
			return
		}

		awayTeamID, err := strconv.Atoi(awayTeamIDStr)
		if err != nil {
			http.Error(w, "Invalid away_team_id.", http.StatusBadRequest)
			return
		}

		err = s.Scrape()
		if err != nil {
			http.Error(w, "Error while scraping data.", http.StatusInternalServerError)
			return
		}

		prediction, err := p.Predict(homeTeamID, awayTeamID)
		if err != nil {
			http.Error(w, "Error while generating prediction.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prediction)
	}
}
