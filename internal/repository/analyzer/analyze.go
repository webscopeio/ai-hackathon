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

	// Call this more for more tools
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

		fmt.Println("new message is ready for processing")

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

		// I have no IDEA what is .ToParam
		messages = append(messages, message.ToParam())

		fmt.Println("These are the params", message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, block := range message.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				var response any
				switch block.Name {
				case sitemapTool.Name:
					fmt.Println("The Sitemap tool is engaged")
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
					fmt.Println("The Content Tool is engaged")
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

				fmt.Println("Marshalling response for", block.Name)
				b, err := json.Marshal(response)
				fmt.Printf("Message size: %d bytes\n", len(b))

				if err != nil {
					return nil, err
				}

				fmt.Println("Successfuly marshalled for", block.Name)

				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, string(b), false))

				fmt.Println("Ready tools results for next iteration")
			}
		}

		if len(toolResults) == 0 {
			fmt.Println("Breaking no results")
			break
		}

		messages = append(messages, anthropic.NewUserMessage(toolResults...))

		fmt.Println("Ready messages for next iteration")

	}
	return nil, nil
}
