package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/webscopeio/ai-hackathon/internal/crawler"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func GenerateTests(client *llm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := decode[models.GenerateTestsArgs](r)
		if err != nil {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: fmt.Sprintf("Bad request, %v", err),
			})
			return
		}

		if args.Url == "" {
			encode(w, http.StatusBadRequest, models.ErrorReturn{
				Error: "url is required",
			})
			return
		}

		_, crawlerResults, err := crawler.Crawl(r.Context(), args.Url, 6, 2)
		if err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't access or process URL, %v", err),
			})
			return
		}

		var builder strings.Builder

		builder.WriteString("WEBSITE CONTENT: \n")
		for url, html := range crawlerResults {
			builder.WriteString("URL: ")
			builder.WriteString(url)
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

		prompt := fmt.Sprintf("%s %s", basePrompt, args.Url)

		// INFO: for a structured response the client requires tools, ref: https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
		tool, toolChoice := llm.GenerateTool[models.GenerateTestsReturn]("get_generate_tests_return", "Generate structured Playwright e2e test scripts based on website analysis. Return organized TypeScript code with proper test organization, assertions, and comments.")

		rawResponse, err := client.GetStructuredCompletion(
			r.Context(),
			context,
			prompt,
			tool,
			toolChoice,
		)
		if err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Couldn't process request, %v", err),
			})
			return
		}

		var response models.GenerateTestsReturn
		if err := json.Unmarshal(rawResponse, &response); err != nil {
			var interlayer struct {
				TestFiles    string   `json:"testFiles"`
				Dependencies []string `json:"dependencies"`
			}
			if err := json.Unmarshal(rawResponse, &interlayer); err != nil {
				encode(w, http.StatusInternalServerError, models.ErrorReturn{
					Error: fmt.Sprintf("Couldn't process response, %v", err),
				})
				return
			}

			var testFiles []models.TestFile
			if err := json.Unmarshal([]byte(interlayer.TestFiles), &testFiles); err != nil {
				encode(w, http.StatusInternalServerError, models.ErrorReturn{
					Error: fmt.Sprintf("Couldn't process response, %v", err),
				})
				return
			}

			response = models.GenerateTestsReturn{
				Dependencies: interlayer.Dependencies,
				TestFiles:    testFiles,
			}
		}

		if err := response.Validate(); err != nil {
			encode(w, http.StatusInternalServerError, models.ErrorReturn{
				Error: fmt.Sprintf("Validation fail, %v", err),
			})
			return
		}

		// TODO: Write the response to a directory

		encode(w, http.StatusOK, response)
	}
}
