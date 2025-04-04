package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/webscopeio/ai-hackathon/internal/config"
)

type Client struct {
	client       *anthropic.Client
	systemPrompt string
	apiKey       string
}

func New(cfg *config.Config) *Client {
	apiKey := cfg.APIKey
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	systemPrompt := "When responding to questions: (1) Analyze problems thoroughly before proposing solutions, (2) Consider edge cases, (3) Acknowledge limitations in your knowledge when appropriate. Your responses should be thoughtful, and demonstrate deep understanding while remaining pragmatic."
	return &Client{
		client:       &client,
		systemPrompt: systemPrompt,
		apiKey:       apiKey,
	}
}

// UpdateAPIKey updates the API key and recreates the client with the new key
func (c *Client) UpdateAPIKey(apiKey string) {
	c.apiKey = apiKey
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	c.client = &client
}
