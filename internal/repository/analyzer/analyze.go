package analyze

import (
	"context"
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Analyze(ctx context.Context, client *llm.Client, urlStr string, prompt string) (*models.AnalyzerReturn, error) {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}

	// Call this more for more tools
	analyzerTool, _ := llm.GenerateTool[models.AnalyzerToolExample]("analyzer_tool_example", "This is a description for the analyzer")

	toolParams := []anthropic.ToolParam{
		*analyzerTool,
	}

	tools := make([]anthropic.ToolUnionParam, len(toolParams))
	for i, toolParam := range toolParams {
		tools[i] = anthropic.ToolUnionParam{OfTool: &toolParam}
	}

	for {
		message, err := client.NewMessage(ctx, anthropic.MessageNewParams{
			Model:     anthropic.ModelClaude_3_5_Sonnet_20240620,
			MaxTokens: 1024,
			Messages:  messages,
			Tools:     tools,
		})
		if err != nil {
			panic(err)
		}

		// I have no IDEA what is .ToParam
		messages = append(messages, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, block := range message.Content {
			switch variant := block.AsAny().(type) {
			case anthropic.ToolUseBlock:
				var response interface{}
				switch block.Name {
				case analyzerTool.Name:
					input := models.AnalyzerToolExample{}
					err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
					if err != nil {
						return nil, err
					}

					// Populate the response with the function return i.e., crawler
					// response = GetCoordinates(input.Location)
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
