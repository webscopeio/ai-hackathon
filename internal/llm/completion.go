package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

func (c *Client) GetCompletion(ctx context.Context, prompt string) (string, error) {
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5HaikuLatest,
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: c.systemPrompt,
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, content := range message.Content {
		if content.Type == "text" {
			sb.WriteString(content.Text)
		}
	}

	return sb.String(), nil
}

func (c *Client) GetStructuredCompletion(
	ctx context.Context,
	context string,
	prompt string,
	tool *anthropic.ToolParam,
	toolChoice *anthropic.ToolChoiceToolParam,
) ([]byte, error) {
	systemBlocks := []anthropic.TextBlockParam{
		{
			Type: "text",
			Text: c.systemPrompt,
		},
		{
			Type: "text",
			Text: "In this environment you have access to a set of tools you can use to answer the user's request. You should use JSON format. Specifications are available in JSONSchema format.",
		},
	}

	if context != "" {
		systemBlocks = append(systemBlocks, anthropic.TextBlockParam{
			Type: "text",
			Text: context,
			CacheControl: anthropic.CacheControlEphemeralParam{
				Type: "ephemeral",
			},
		})
	}

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.ModelClaude3_5SonnetLatest,
		// INFO: tools typically require more tokens
		MaxTokens: 2400,
		System:    systemBlocks,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		Tools: []anthropic.ToolUnionParam{
			{
				OfTool: tool,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfToolChoiceTool: toolChoice,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			if block.Name == tool.Name {
				return []byte(variant.JSON.Input.Raw()), nil
			}
		}
	}

	return nil, fmt.Errorf("message content missing")
}
