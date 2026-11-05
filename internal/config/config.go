package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the values the application needs to start.
type Config struct {
	DatabaseURL string
	Platform    string
	Address     string
}

// Load reads configuration from the environment, loading a local .env file
// first when one is present.
func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL: os.Getenv("DB_URL"),
		Platform:    os.Getenv("PLATFORM"),
		Address:     envOrDefault("ADDRESS", ":8080"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DB_URL must be set")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
