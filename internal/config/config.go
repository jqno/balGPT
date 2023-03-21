package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBConnectionString string
	ScraperURL         string
}

func LoadConfig() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	scraperURL := os.Getenv("SCRAPER_URL")

	return &Config{
		DBConnectionString: connectionString,
		ScraperURL:         scraperURL,
	}
}
