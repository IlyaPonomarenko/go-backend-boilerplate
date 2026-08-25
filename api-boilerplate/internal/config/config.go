package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config contains runtime settings supplied by the environment.
type Config struct {
	Address  string
	Database string
	APIToken string
}

// Load reads configuration and requires an API token outside explicit local
// development mode, preventing an accidentally unauthenticated deployment.
func Load() (Config, error) {
	// Load local development settings when present. godotenv.Load does not
	// overwrite variables already supplied by the shell or deployment system.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	settings := Config{
		Address:  envOr("API_ADDRESS", ":8080"),
		Database: envOr("DATABASE_PATH", "./data/api.db"),
		APIToken: os.Getenv("API_TOKEN"),
	}
	if settings.APIToken == "" && os.Getenv("APP_ENV") != "development" {
		return Config{}, errors.New("API_TOKEN must be set unless APP_ENV=development")
	}
	return settings, nil
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
