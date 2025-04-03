package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/webscopeio/ai-hackathon/internal/handlers"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func New() *chi.Mux {
	return chi.NewRouter()
}

func RegisterRoutes(r *chi.Mux, llm *llm.Client) {
	r.Get("/status", handlers.Status)

	r.Post("/crawl", handlers.Crawl)

	r.Post("/generate-tests", handlers.GenerateTests(llm))
}
