package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/models"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		models.PostReturn{
			Posts: []models.Post{
				{
					ID:    1,
					Title: "Introduction to Go",
					Body:  "Go is a statically typed, compiled language designed at Google.",
				},
				{
					ID:    2,
					Title: "RESTful APIs with Chi",
					Body:  "Chi is a lightweight, idiomatic and composable router for building Go HTTP services.",
				},
				{
					ID:    3,
					Title: "Middleware in Go",
					Body:  "Middleware is a powerful concept that allows you to add functionality to your request handling pipeline.",
				},
				{
					ID:    4,
					Title: "JSON Serialization",
					Body:  "Go provides excellent support for encoding and decoding JSON with the encoding/json package.",
				},
				{
					ID:    5,
					Title: "Error Handling Patterns",
					Body:  "Proper error handling is crucial for building robust and maintainable applications.",
				},
			},
		},
	)
}
