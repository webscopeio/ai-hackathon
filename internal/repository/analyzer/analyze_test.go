package analyze

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func TestAnalyze(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg := config.Load()
	llm := llm.New(cfg)

	res, err := Analyze(ctx, llm, "example.com", "Check out the website")
	if err != nil {
		t.Fatalf("err=%v", err)
	}

	fmt.Printf("res=%v", res)
}
