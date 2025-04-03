package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Greet(w http.ResponseWriter, r *http.Request) {
	var args models.GreetArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorReturn{Error: "Bad request"})
		return
	}

	if args.Message == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorReturn{Error: "Message is required"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.GreetReturn{
		Message: args.Message + " is OK!",
	}
	json.NewEncoder(w).Encode(response)
}
