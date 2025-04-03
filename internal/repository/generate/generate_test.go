package generate

import (
	"context"
	"testing"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func TestGenerate(t *testing.T) {
	cfg := config.Load()

	client := llm.New(cfg)
	_, err := GenerateTests(context.Background(), client, "https://example.com")
	if err != nil {
		t.Errorf("GenerateTests failed: %v", err)
	}
}
