package debug

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Debug prints a debug message if debug mode is enabled
// Debug mode is enabled if the DEBUG_MODE environment variable is set to "true"
// or if it's defined as DEBUG_MODE=true in a .env file in the current directory
func Debug(format string, args ...interface{}) {
	// Check if debug mode is enabled
	if isDebugModeEnabled() {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// isDebugModeEnabled checks if debug mode is enabled via environment variable
// or .env file
func isDebugModeEnabled() bool {
	// First check environment variable
	debugMode := os.Getenv("DEBUG_MODE")
	if strings.ToLower(debugMode) == "true" {
		return true
	}

	// Then check .env file if it exists
	envMap, err := godotenv.Read()
	if err == nil {
		// Check if DEBUG_MODE is set to true in .env file
		if value, exists := envMap["DEBUG_MODE"]; exists && strings.ToLower(value) == "true" {
			return true
		}
	}

	return false
}
