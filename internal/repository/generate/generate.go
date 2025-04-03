package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/webscopeio/ai-hackathon/internal/crawler"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Generate() {
	fmt.Println("Hello from Generate")
}

// Tests generates test files based on a URL using the LLM client
func Tests(ctx context.Context, client *llm.Client, url string) (models.GenerateTestsReturn, error) {
	_, crawlerResults, err := crawler.Crawl(ctx, url, 6, 2)
	if err != nil {
		return models.GenerateTestsReturn{}, fmt.Errorf("couldn't access or process URL: %w", err)
	}

	var builder strings.Builder

	builder.WriteString("WEBSITE CONTENT: \n")
	for pageUrl, html := range crawlerResults {
		builder.WriteString("URL: ")
		builder.WriteString(pageUrl)
		builder.WriteString("\nCONTENT:\n")
		builder.WriteString(html)
		builder.WriteString("\n---END PAGE---\n\n")
	}

	context := builder.String()

	basePrompt := `Please analyze the website content provided and generate comprehensive end-to-end tests written in TypeScript using Playwright Framework. Focus on:
1. Critical user flows and main functionality
2. Navigation and routing
3. Error states and edge cases

Format the tests following Playwright best practices with clear test descriptions and organized test suites.`

	basePrompt += `IMPORTANT: When providing test files, ensure proper JSON formatting:
1. The "testFiles" field must be an array, not a string containing an array
2. All string values, including code in the "content" field, must be enclosed in double quotes (")
   and properly escaped, NOT backticks
3. Use proper JSON escaping for special characters in content: \n for newlines, \" for quotes
4. Example of correct format:
   "testFiles": [
     {
       "filename": "test.spec.ts",
       "content": "import { test } from '@playwright/test';\n\ntest('example', async () => {\n  // code\n});"
     }
   ]

Invalid formatting will cause errors in processing your response.
`

	prompt := fmt.Sprintf("%s %s", basePrompt, url)

	// INFO: for a structured response the client requires tools, ref: https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
	tool, toolChoice := llm.GenerateTool[models.GenerateTestsReturn]("get_generate_tests_return", "Generate structured Playwright e2e test scripts based on website analysis. Return organized TypeScript code with proper test organization, assertions, and comments.")

	rawResponse, err := client.GetStructuredCompletion(
		ctx,
		context,
		prompt,
		tool,
		toolChoice,
	)
	if err != nil {
		return models.GenerateTestsReturn{}, fmt.Errorf("couldn't process request: %w", err)
	}

	var response models.GenerateTestsReturn
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		var interlayer struct {
			TestFiles    string   `json:"testFiles"`
			Dependencies []string `json:"dependencies"`
		}
		if err := json.Unmarshal(rawResponse, &interlayer); err != nil {
			return models.GenerateTestsReturn{}, fmt.Errorf("couldn't process response: %w", err)
		}

		var testFiles []models.TestFile
		if err := json.Unmarshal([]byte(interlayer.TestFiles), &testFiles); err != nil {
			return models.GenerateTestsReturn{}, fmt.Errorf("couldn't process response: %w", err)
		}

		response = models.GenerateTestsReturn{
			Dependencies: interlayer.Dependencies,
			TestFiles:    testFiles,
		}
	}

	if err := response.Validate(); err != nil {
		return models.GenerateTestsReturn{}, fmt.Errorf("validation fail: %w", err)
	}

	return response, nil
}
