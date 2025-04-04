package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/models"
	"github.com/webscopeio/ai-hackathon/internal/repository/config"
)

// GetConfig handles retrieving the user configuration
func GetConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create config manager
		configManager, err := config.NewManager()
		if err != nil {
			log.Printf("Error creating config manager: %v", err)
			http.Error(w, "Failed to initialize configuration manager", http.StatusInternalServerError)
			return
		}

		// Get current configuration
		userConfig, err := configManager.GetConfig()
		if err != nil {
			log.Printf("Error retrieving configuration: %v", err)
			http.Error(w, "Failed to retrieve configuration", http.StatusInternalServerError)
			return
		}

		// Return configuration as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userConfig); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// SaveConfig handles saving or updating the user configuration
func SaveConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var configUpdate models.UserConfig
		if err := json.NewDecoder(r.Body).Decode(&configUpdate); err != nil {
			log.Printf("Error parsing request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create config manager
		configManager, err := config.NewManager()
		if err != nil {
			log.Printf("Error creating config manager: %v", err)
			http.Error(w, "Failed to initialize configuration manager", http.StatusInternalServerError)
			return
		}

		// Update configuration
		if err := configManager.UpdateConfig(&configUpdate); err != nil {
			log.Printf("Error updating configuration: %v", err)
			http.Error(w, "Failed to update configuration", http.StatusInternalServerError)
			return
		}

		// Get updated configuration to return
		updatedConfig, err := configManager.GetConfig()
		if err != nil {
			log.Printf("Error retrieving updated configuration: %v", err)
			http.Error(w, "Failed to retrieve updated configuration", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Success bool              `json:"success"`
			Config  *models.UserConfig `json:"config"`
		}{
			Success: true,
			Config:  updatedConfig,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
