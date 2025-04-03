package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.StatusReturn{
		Status: "OK",
	})
}
