package analyze

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Analyze(ctx context.Context, cfg *config.Config, client *llm.Client, urlStr string, prompt string) (*models.AnalyzerReturn, error) {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf("The website is: %s - %s", urlStr, prompt))),
	}

	sitemapTool, _ := llm.GenerateTool[models.SitemapTool]("sitemap_tool", "This tool is able to get a website's sitemap using a base URL")
	getContentTool, _ := llm.GenerateTool[models.GetContentTool]("get_content_tool", "This tool is able to get the body content for a list of important URLs")
	sentryTool, _ := llm.GenerateTool[models.SentryTool]("sentry_tool", "This tool is able to get error information from Sentry for a specific project")

	toolParams := []anthropic.ToolParam{
		*sitemapTool,
		*getContentTool,
		*sentryTool,
	}

	tools := make([]anthropic.ToolUnionParam, len(toolParams))
	for i, toolParam := range toolParams {
		tools[i] = anthropic.ToolUnionParam{OfTool: &toolParam}
	}

	for {
		message, err := client.NewMessage(ctx, anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude3_5HaikuLatest,
			MaxTokens: 1024,
			Messages:  messages,
			Tools:     tools,
		})
		if err != nil {
			return nil, err
		}

		// This is a debugging block
		for _, block := range message.Content {
			switch block := block.AsAny().(type) {
			case anthropic.TextBlock:
				fmt.Println(block.Text)
			case anthropic.ToolUseBlock:
				inputJSON, _ := json.Marshal(block.Input)
				fmt.Println(block.Name + ": " + string(inputJSON))
			}
		}

		messages = append(messages, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, block := range message.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				var response any
				switch block.Name {
				case sitemapTool.Name:
					input := models.SitemapTool{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					response, err = GetSitemap(ctx, input.BaseUrl)
					if err != nil {
						return nil, err
					}
				case getContentTool.Name:
					input := models.GetContentTool{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					response, err = GetContent(ctx, input.Urls)
					if err != nil {
						return nil, err
					}
				case sentryTool.Name:
					input := models.SentryTool{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					response, err = GetSentryIssues(ctx, cfg, input.OrgSlug, input.ProjectSlug)
					if err != nil {
						return nil, err
					}
				}

				b, err := json.Marshal(response)
				if err != nil {
					return nil, err
				}

				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, string(b), false))
			}
		}

		if len(toolResults) == 0 {
			break
		}

		messages = append(messages, anthropic.NewUserMessage(toolResults...))

	}
	
	// Create a simple analyzer return with dummy data
	// In a real implementation, this would process the LLM's final response
	return &models.AnalyzerReturn{
		TechSpec: "Website technology analysis",
		SiteMap: map[string]string{
			"home": urlStr,
		},
		Criteria: "Analysis criteria",
	}, nil
}
