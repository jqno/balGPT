package config

import (
	"os"
)

type Config struct {
	DBConnectionString string
}

func LoadConfig() *Config {
	return &Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}
