package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/iamhectorsosa/ai-hackathon/internal/handlers"
	"github.com/iamhectorsosa/ai-hackathon/internal/llm"
)

func New() *chi.Mux {
	return chi.NewRouter()
}

func RegisterRoutes(r *chi.Mux, llm *llm.Client) {
	r.Get("/status", handlers.Status)
	r.Get("/error", handlers.Error)

	r.Get("/posts", handlers.GetPosts)

	r.Post("/greet", handlers.Greet)

	r.Post("/ask", handlers.Ask(llm))

	r.Post("/crawl", handlers.Crawl)

	r.Post("/create-job", handlers.CreateJob(llm))

	r.Post("/generate-tests", handlers.GenerateTests(llm))
}
