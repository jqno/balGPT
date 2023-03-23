package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBConnectionString string
	AuthUsername       string
	AuthPassword       string
	ScraperURL         string
	ApiBaseURL         string
}

func LoadConfig() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	authUsername := os.Getenv("AUTH_USERNAME")
	authPassword := os.Getenv("AUTH_PASSWORD")

	scraperURL := os.Getenv("SCRAPER_URL")
	apiBaseURL := os.Getenv("API_BASE_URL")

	return &Config{
		DBConnectionString: connectionString,
		AuthUsername:       authUsername,
		AuthPassword:       authPassword,
		ScraperURL:         scraperURL,
		ApiBaseURL:         apiBaseURL,
	}
}
