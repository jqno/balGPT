package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBConnectionString string
	ScraperURL         string
	AllowedEmail       string
	GoogleClientID     string
	ApiBaseURL         string
}

func LoadConfig() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	scraperURL := os.Getenv("SCRAPER_URL")
	allowedEmail := os.Getenv("ALLOWED_EMAIL")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	apiBaseURL := os.Getenv("API_BASE_URL")

	return &Config{
		DBConnectionString: connectionString,
		ScraperURL:         scraperURL,
		AllowedEmail:       allowedEmail,
		GoogleClientID:     googleClientId,
		ApiBaseURL:         apiBaseURL,
	}
}
