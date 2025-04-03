package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/webscopeio/ai-hackathon/internal/crawler"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Crawl(w http.ResponseWriter, r *http.Request) {
	var args models.CrawlArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorReturn{Error: "Bad request"})
		return
	}

	if args.Url == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorReturn{Error: "URL and Depth are required"})
		return
	}

	links, results, err := crawler.Crawl(r.Context(), args.Url, args.MaxDepth, args.MaxPathSegments)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorReturn{Error: fmt.Sprintf("Unable to crawl, %s", err.Error())})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.CrawlReturn{
		Links:   links,
		Results: results,
	}
	json.NewEncoder(w).Encode(response)
}
