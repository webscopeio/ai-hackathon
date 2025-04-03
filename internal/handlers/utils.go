package handlers

import (
	"encoding/json"
	"net/http"
)

func encode[T any](w http.ResponseWriter, statusCode int, payload T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}
