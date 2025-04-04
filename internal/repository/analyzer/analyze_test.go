package analyzer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func TestAnalyze(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	cfg := config.Load()
	llm := llm.New(cfg)

	res, err := Analyze(ctx, cfg, llm, "https://ai-hackathon-demo-delta.vercel.app/", "Check out the website, wonder how is it structured?. I am interested in the content of the most valuable pages to create the criteria to generate an E2E tests.")
	if err != nil {
		t.Fatalf("err=%v", err)
	}

	fmt.Printf("criteria=%v", res.Criteria)
	fmt.Printf("techSpec=%v", res.TechSpec)
	fmt.Printf("contentMap=%v", res.ContentMap)
}
