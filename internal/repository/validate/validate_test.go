package validate

import (
	"context"
	"testing"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func TestValidate(t *testing.T) {
	cfg := config.Load()
	client := llm.New(cfg)
	
	analysis, err := Validate(context.Background(), client, "temp")
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	// Verify the analysis structure
	if analysis.Failures == nil {
		t.Error("Expected Failures to be initialized, got nil")
	}
}
