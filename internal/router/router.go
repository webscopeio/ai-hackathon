package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/handlers"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func New() *chi.Mux {
	return chi.NewRouter()
}

func RegisterRoutes(r *chi.Mux, cfg *config.Config, llm *llm.Client) {
	r.Get("/status", handlers.Status)

	r.Post("/crawl", handlers.Crawl)

	// Configuration endpoints
	r.Get("/config", handlers.GetConfig())
	r.Post("/config", handlers.SaveConfig())

	// Analyze endpoints
	r.Post("/analyze", handlers.Analyze(cfg, llm))
}
