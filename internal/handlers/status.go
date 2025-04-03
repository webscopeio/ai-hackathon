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

// TODO: delete once Demo'd
func Error(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.ErrorReturn{
		Error: "Houston is down (intentional)",
	})
}
