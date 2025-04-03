package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/iamhectorsosa/ai-hackathon/internal/config"
	"github.com/iamhectorsosa/ai-hackathon/internal/llm"
	"github.com/iamhectorsosa/ai-hackathon/internal/router"
)

func main() {
	cfg := config.Load()

	r := router.New()
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, time.Minute))

	llm := llm.New(cfg)
	router.RegisterRoutes(r, llm)

	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server starting on localhost%s in %s mode\n", addr, cfg.Environment)
	http.ListenAndServe(addr, r)
}
