package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Analyze(ctx context.Context, cfg *config.Config, client *llm.Client, urlStr string, prompt string) (*models.AnalyzerReturn, error) {
	userMessage := fmt.Sprintf("The website is: %s - %s", urlStr, prompt)
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
	}

	sitemapTool, _ := llm.GenerateTool[models.SitemapTool]("sitemap_tool", "This tool is able to get a website's sitemap using a base URL")
	getContentTool, _ := llm.GenerateTool[models.GetContentTool]("get_content_tool", "This tool is able to get the body content for a list of important URLs")
	sentryTool, _ := llm.GenerateTool[models.SentryTool]("get_sentry_tool", "This tool is able to get error information from Sentry for a specific project to give you a better context about the website")
	finalCriteriaTool, _ := llm.GenerateTool[models.FinalCriteriaTool]("get_final_criteria_tool", "This tool is able to get the final criteria for the analysis of the website from results of the other tools, run this always as the last step")

	toolParams := []anthropic.ToolParam{
		*sitemapTool,
		*getContentTool,
		*sentryTool,
		*finalCriteriaTool,
	}

	tools := make([]anthropic.ToolUnionParam, len(toolParams))
	for i, toolParam := range toolParams {
		tools[i] = anthropic.ToolUnionParam{OfTool: &toolParam}
	}

	var contentMap map[string]string

	fmt.Println("\n[ANALYZER] User Message: \n\n", userMessage)

	for {
		message, err := client.NewMessage(ctx, anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude3_5SonnetLatest,
			MaxTokens: 2048,
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
				fmt.Printf("\n[ANALYZER] Agent response: \n\n%s\n", block.Text)
			case anthropic.ToolUseBlock:
				inputJSON, _ := json.Marshal(block.Input)
				fmt.Printf("\n[ANALYZER] Tool call: \n\n%s\n", block.Name+": "+string(inputJSON))
			}
		}

		messages = append(messages, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, block := range message.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				var response any
				// fmt.Println("Running block", block.Name)
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

					result, err := GetContent(ctx, input.Urls)
					if err != nil {
						return nil, err
					}
					contentMap = result.Contents
					response = result
				case sentryTool.Name:
					input := models.SentryTool{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					response, err = GetSentryIssues(ctx, cfg, input.OrgSlug, input.ProjectSlug)
					if err != nil {
						return nil, fmt.Errorf("failed to get Sentry issues: %w", err)
					}
				case finalCriteriaTool.Name:
					logger.Debug("FROM ANALYZE: Final criteria tool raw: %s", variant.JSON.Input.Raw())
					input := models.FinalCriteriaTool{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					logger.Debug("FROM ANALYZE: Final contentMap: %s", input.ContentMap)

					return &models.AnalyzerReturn{
						TechSpec:   prompt,
						ContentMap: contentMap,
						Criteria:   input.Criteria,
					}, nil
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

	return nil, errors.New("no valid response from the model")
}
