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
	Predictor predictor.Predictor
}

type TemplateData struct {
	ApiBaseURL string
	Teams      []team.Team
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString, cfg.AppBaseDir)
	scraper := scraper.NewScrapeData(db, cfg.ScraperURL)

	predictor := predictor.NewCompositePredictor(
		predictor.NewHomeAdvantagePredictor(),
		predictor.NewAverageGoalsPredictor(db),
		predictor.NewLastYearMatchPredictor(db),
		predictor.NewFlippedLastYearMatchPredictor(db),
		predictor.NewLeaderboardDifferencePredictor(db),
	)

	return &App{
		Config:    cfg,
		DB:        db,
		Scraper:   scraper,
		Predictor: predictor,
	}
}

func (a *App) Run() {
	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL, a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/login", checkAuth(loginHandler(), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/health", healthCheckHandler(a.DB))

	staticDir := filepath.Join(a.Config.AppBaseDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func checkAuth(h http.HandlerFunc, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		log.Printf("Login attempt by %s: %v", username, ok)
		log.Printf("Requested url: %s", r.URL.String())

		if !ok || username != validUsername || password != validPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Authenticated.")
	}
}

func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		isAuthenticated := ok && username == validUsername && password == validPassword
		log.Printf("Index: login attempt by %s: %v", username, isAuthenticated)
		log.Printf("Index: requested url: %s", r.URL.String())

		if !isAuthenticated {
			templateFile := filepath.Join(appBaseDir, "templates/login.html")

			tmpl, err := template.ParseFiles(templateFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			teams, err := db.FetchTeamsFromDB()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
				return
			}

			data := TemplateData{
				ApiBaseURL: apiBaseURL,
				Teams:      teams,
			}

			templateFile := filepath.Join(appBaseDir, "templates/main.html")
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
}

func handlePrediction(s *scraper.ScrapeData, p predictor.Predictor) http.HandlerFunc {
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
			log.Printf("Error: %s", err)
			http.Error(w, "Error while scraping data.", http.StatusInternalServerError)
			return
		}

		prediction, err := p.Predict(homeTeamID, awayTeamID)
		if err != nil {
			log.Printf("Error: %s", err)
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
			log.Printf("Error: %s", err)
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

func healthCheckHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the database connection
		if err := db.Conn.Ping(); err != nil {
			// If there is an error, return a 500 Internal Server Error status code
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}

		// If everything is fine, return a 200 OK status code
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
