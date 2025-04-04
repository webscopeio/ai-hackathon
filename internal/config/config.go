package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Environment     string
	APIKey          string
	SentryAuthToken string
	UmamiURL        string
	UmamiAPIKey     string
	UmamiWebsiteId  string
}

func Load() *Config {
	cfg := &Config{
		Port:            "8080",
		Environment:     "development",
		APIKey:          "",
		SentryAuthToken: "",
		UmamiURL:        "https://api.umami.is/v1",
		UmamiAPIKey:     "",
		UmamiWebsiteId:  "",
	}

	workDir, _ := os.Getwd()
	rootDir := workDir

	for {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			rootDir = workDir
			break
		}
		rootDir = parent
	}

	envMap, err := godotenv.Read(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Println("No .env file found, using default configuration")
		return cfg
	}

	if port := envMap["PORT"]; strings.TrimSpace(port) != "" {
		cfg.Port = port
	}

	if env := envMap["ENVIRONMENT"]; strings.TrimSpace(env) != "" {
		cfg.Environment = env
	}

	if apiKey := envMap["API_KEY"]; strings.TrimSpace(apiKey) != "" {
		cfg.APIKey = apiKey
	}

	if sentryAuthToken := envMap["SENTRY_AUTH_TOKEN"]; strings.TrimSpace(sentryAuthToken) != "" {
		cfg.SentryAuthToken = sentryAuthToken
	}

	if umamiURL := envMap["UMAMI_URL"]; strings.TrimSpace(umamiURL) != "" {
		cfg.UmamiURL = umamiURL
	}

	if umamiAPIKey := envMap["UMAMI_API_KEY"]; strings.TrimSpace(umamiAPIKey) != "" {
		cfg.UmamiAPIKey = umamiAPIKey
	}

	if umamiWebsiteId := envMap["UMAMI_WEBSITE_ID"]; strings.TrimSpace(umamiWebsiteId) != "" {
		cfg.UmamiWebsiteId = umamiWebsiteId
	}

	return cfg
}
