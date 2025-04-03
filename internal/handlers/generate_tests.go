package handlers

import (
	"fmt"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
	"github.com/webscopeio/ai-hackathon/internal/repository/generate"
)

func GenerateTests(client *llm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := decode[models.GenerateTestsArgs](r)
		if err != nil {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: fmt.Sprintf("Bad request, %v", err),
			})
			return
		}

		if args.Url == "" {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: "url is required",
			})
			return
		}

		response, err := generate.GenerateTests(r.Context(), client, args.Url)
		if err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Error generating tests: %v", err),
			})
			return
		}

		// TODO: Write the response to a directory

		encode(w, http.StatusOK, response)
	}
}
