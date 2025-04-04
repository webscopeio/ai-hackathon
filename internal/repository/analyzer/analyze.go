package analyze

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Analyze(ctx context.Context, client *llm.Client, urlStr string, prompt string) (*models.AnalyzerReturn, error) {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf("The website is: %s - %s", urlStr, prompt))),
	}

	sitemapTool, _ := llm.GenerateTool[models.SitemapTool]("sitemap_tool", "This tool is able to get a website's sitemap using a base URL")
	getContentTool, _ := llm.GenerateTool[models.GetContentTool]("get_content_tool", "This tool is able to get the body content for a list of important URLs")

	toolParams := []anthropic.ToolParam{
		*sitemapTool,
		*getContentTool,
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

					response, err = GetSitemap(input.BaseUrl)
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
	return nil, nil
}
