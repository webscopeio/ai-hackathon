package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/webscopeio/ai-hackathon/internal/config"
)

type Client struct {
	client       *anthropic.Client
	systemPrompt string
}

func New(cfg *config.Config) *Client {
	client := anthropic.NewClient(option.WithAPIKey(cfg.APIKey))
	systemPrompt := "When responding to questions: (1) Analyze problems thoroughly before proposing solutions, (2) Consider edge cases, (3) Acknowledge limitations in your knowledge when appropriate. Your responses should be thoughtful, and demonstrate deep understanding while remaining pragmatic."
	return &Client{
		client:       &client,
		systemPrompt: systemPrompt,
	}
}
