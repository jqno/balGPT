package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jqno/balGPT/internal/config"
	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/predictor"
	"github.com/jqno/balGPT/internal/scraper"
	"golang.org/x/oauth2"
	oauth2api "google.golang.org/api/oauth2/v2"
)

type App struct {
	Config    *config.Config
	DB        *database.DB
	Scraper   *scraper.ScrapeData
	Predictor *predictor.Predictor
}

func NewApp(cfg *config.Config) *App {
	db := database.New(cfg.DBConnectionString)
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
	allowedEmail := a.Config.AllowedEmail
	http.HandleFunc("/predict", googleAuthMiddleware(handlePrediction(a.Scraper, a.Predictor), allowedEmail))
	http.HandleFunc("/team_id", googleAuthMiddleware(handleTeamID(a.DB), allowedEmail))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func googleAuthMiddleware(next http.Handler, allowedEmail string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "access token cookie not found", http.StatusUnauthorized)
			return
		}

		accessToken := accessTokenCookie.Value
		client, err := getAuthenticatedClient(accessToken, allowedEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "client", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuthenticatedClient(accessToken string, allowedEmail string) (*http.Client, error) {
	ctx := context.Background()

	// Use the access token to create an authenticated client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	client := oauth2.NewClient(ctx, ts)

	// Check if the token is valid
	oauth2Service, err := oauth2api.New(client)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := oauth2Service.Tokeninfo().AccessToken(accessToken).Do()
	if err != nil {
		return nil, err
	}

	// Check if the email address is allowed
	if allowedEmail != tokenInfo.Email {
		return nil, fmt.Errorf("unauthorized email address")
	}

	return client, nil
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
