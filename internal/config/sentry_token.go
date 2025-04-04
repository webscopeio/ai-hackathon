package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadSentryAuthToken loads the Sentry auth token from environment variables
// It attempts to load from .env file first, then falls back to OS environment variables
func LoadSentryAuthToken() string {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get auth token from environment
	return os.Getenv("SENTRY_AUTH_TOKEN")
}
