package handlers

import (
	"fmt"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
	"github.com/webscopeio/ai-hackathon/internal/repository/analyzer"
)

func Analyze(cfg *config.Config, client *llm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := decode[models.AnalyzerArgs](r)
		if err != nil {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: fmt.Sprintf("Bad request, %v", err),
			})
			return
		}

		res, err := analyzer.Analyze(r.Context(), cfg, client, args.Url, args.Prompt)
		if err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't analyze website, %v", err),
			})
			return
		}

		encode(w, http.StatusOK, res)
	}
}
