package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jqno/balGPT/internal/config"
	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/predictor"
	"github.com/jqno/balGPT/internal/scraper"
	"github.com/jqno/balGPT/internal/team"
)

type App struct {
	Config    *config.Config
	DB        *database.DB
	Scraper   *scraper.ScrapeData
	Predictor *predictor.Predictor
}

type TemplateData struct {
	ApiBaseURL string
	Teams      []team.Team
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString, cfg.AppBaseDir)
	scraper := scraper.NewScrapeData(db, cfg.ScraperURL)
	predictor := predictor.NewPredictor(db)
	return &App{
		Config:    cfg,
		DB:        db,
		Scraper:   scraper,
		Predictor: predictor,
	}
}

func (a *App) Run() {
	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL))
	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))

	staticDir := filepath.Join(a.Config.AppBaseDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func checkAuth(h http.HandlerFunc, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		log.Printf("Login attempt by %s: %v", username, ok)
		if !ok || username != validUsername || password != validPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := db.FetchTeamsFromDB()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			ApiBaseURL: apiBaseURL,
			Teams:      teams,
		}

		templateFile := filepath.Join(appBaseDir, "templates/index.html")
		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
			return
		}
	}
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

		log.Printf("Prediction for home_team_id=%d, away_team_id=%d: %d - %d", homeTeamID, awayTeamID, prediction.HomeGoals, prediction.AwayGoals)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prediction)
	}
}

func handleScrape(s *scraper.ScrapeData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Scrape()
		if err != nil {
			http.Error(w, "Error while scraping data.", http.StatusInternalServerError)
			return
		}

		log.Printf("Scraping complete")

		w.WriteHeader(http.StatusOK)
	}
}

func handleTeamID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		if teamName == "" {
			http.Error(w, "team_name is required.", http.StatusBadRequest)
			return
		}

		teamID, err := db.GetTeamID(teamName)
		if err != nil {
			http.Error(w, "Error while fetching team ID.", http.StatusInternalServerError)
			return
		}

		if teamID == 0 {
			http.Error(w, "Team not found.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"team_id": teamID})
	}
}
