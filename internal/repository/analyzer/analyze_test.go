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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cfg := config.Load()
	llm := llm.New(cfg)

	res, err := Analyze(ctx, cfg, llm, "https://ai-hackathon-demo-delta.vercel.app/", "Check out the website, wonder how is it structured?. I am interested in the content of the most valuable pages to create the criteria to generate an E2E tests. My orgSlug := \"webscopeio-pb\" and projectSlug := \"ai-hackathon-demo\" for Sentry, please check the errors in the last 14 days and include them in the analysis.")
	if err != nil {
		t.Fatalf("err=%v", err)
	}

	fmt.Printf("criteria=%v", res.Criteria)
	fmt.Printf("techSpec=%v", res.TechSpec)
	fmt.Printf("contentMap=%v", res.ContentMap)
}
