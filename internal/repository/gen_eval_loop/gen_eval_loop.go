package gen_eval_loop

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func GenEvalLoop(ctx context.Context, client *llm.Client, analyzerReturn *models.AnalyzerReturn) (string, error) {
	tempDir, testsDir, err := SetupTestEnvironment(ctx)
	if err != nil {
		return "", fmt.Errorf("SetupTestEnvironment failed: %w", err)
	}

	generatorMessages := []anthropic.MessageParam{}
	feedback := ""
	filename := ""
	testFileContent := []byte{}
	loopCount := 0

	for {
		if loopCount > 3 && false {
			return "", fmt.Errorf("Exceeded maximum number of loops")
		}
		filename, newMessages, err := generateTestFile(ctx, client, analyzerReturn, generatorMessages, feedback, string(testFileContent), testsDir)
		if err != nil {
			return "", fmt.Errorf("GenerateTestFile failed: %w", err)
		}
		generatorMessages = newMessages
		logger.Debug("Filename: %s", filename)

		var passed bool
		feedback, passed, err = evaluateTestFile(ctx, client, filename, tempDir, testsDir)
		logger.Debug("EVALUATOR feedback: %s", feedback)
		if err != nil {
			return "", fmt.Errorf("EvaluateTestFile failed: %w", err)
		}

		if passed {
			fmt.Printf("✅ Evaluator accepted the test file.\n")
			break
		}

		testPath := filepath.Join(testsDir, filename)
		testFileContent, err = os.ReadFile(testPath)
		if err != nil {
			return "", fmt.Errorf("ReadTestFile failed: %w", err)
		}

		feedback = `FEEDBACK: ` + feedback
		loopCount++
	}

	filePath := filepath.Join(testsDir, filename)

	return filePath, nil
}

// Tests generates test files based on a URL using the LLM client
// It also stores the generated test files in a temporary directory
func generateTestFile(ctx context.Context, client *llm.Client, analyzerReturn *models.AnalyzerReturn, prevMessages []anthropic.MessageParam, feedback string, testFileContent string, testsDir string) (string, []anthropic.MessageParam, error) {

	logger.Debug("Starting generateTestFile")

	if analyzerReturn == nil {
		return "", prevMessages, fmt.Errorf("analyzerReturn is nil")
	}

	var builder strings.Builder

	builder.WriteString("INPUTS: \n")
	builder.WriteString("TECHNICAL SPECIFICATION: ")
	builder.WriteString(analyzerReturn.TechSpec)
	builder.WriteString("\nCONTENT MAP (SEPARATED BY 2 NEWLINES): ")
	for url, content := range analyzerReturn.ContentMap {
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", url, content))
	}
	builder.WriteString("\nTEST CRITERIA: ")
	builder.WriteString(analyzerReturn.Criteria)
	builder.WriteString("\n---END PAGE---\n\n")

	context := builder.String()

	basePrompt := `You are a test engineer, your task is to write a focused end-to-end test suite written in TypeScript using Playwright Framework. You have been provided inputs from an analyzer.
The inputs are a technical specification (description) of the website, a map (directory) of the website's pages with urls as keys and html page content as values and a test criteria (scenario) for
the test you need to write. Focus on the provided criteria and tech spec. Your output should be one test suite file.

Important points:
- Focus on the provided criteria
- Do not add any other dependencies, only @playwright/test is allowed.
- The test file should be around 100 lines of code, the closer the better.
- Write consise test cases that won't fail instead of complex cases.

Format the tests following Playwright best practices with clear test descriptions and organized test suites.`

	basePrompt += `IMPORTANT: When providing the test file, ensure proper JSON formatting:
1. The "filename" field must be a string and cannot be empty
2. All string values, including code in the "content" field, must be enclosed in double quotes (").
3. Use \n for newlines, and try to use single quotes where possible (so we don't have to escape the quotes).
4. Never over-escape the quotes, only escape them when necessary.
5. Example of correct format:
   {
	"filename": "test.spec.ts",
	"content": "import { test } from '@playwright/test';\n\ntest('example', async () => {\n  // code\n});"
	"dependencies": ["@playwright/test"]
   }

Invalid formatting will cause errors in processing your response.
`
	if feedback != "" {
		basePrompt = `An Evaluation of the test file has been provided. Please revise the test file based on the feedback. Leave everything else the same. Mainly focus on fixing the failing tests. The resulting file should be no longer than 100 lines of code.
` + feedback
		basePrompt += `
TEST FILE CURRENT CONTENT:

` + testFileContent
	}

	// INFO: for a structured response the client requires tools, ref: https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
	tool, toolChoice := llm.GenerateTool[models.GenerateTestReturn]("get_generate_test_file_return", "Generate structured Playwright e2e test suite based on provided inputs. Return organized TypeScript code with proper test organization, assertions, and comments.")

	logger.Debug("GENERATOR Sending request to LLM with context length: %d characters", len(context))
	logger.Debug("GENERATOR feedback: %s", feedback)
	logger.Debug("GENERATOR Prompt:\n %s", basePrompt)
	logger.Debug("GENERATOR Previous messages:\n %v", prevMessages)
	rawResponse, err := client.GetStructuredCompletion(
		ctx,
		context,
		basePrompt,
		tool,
		toolChoice,
		prevMessages,
	)
	if err != nil {
		logger.Debug("LLM request failed: %v", err)
		return "", prevMessages, fmt.Errorf("couldn't process request: %w", err)
	}
	logger.Debug("Received response from LLM with length: %d characters", len(rawResponse))
	logger.Debug("GENERATOR Response:\n %s", string(rawResponse))

	newMessages := append(prevMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(string(rawResponse))))

	logger.Debug("Unmarshalling LLM response")
	var response models.GenerateTestReturn
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		logger.Debug("Primary unmarshal failed, trying fallback: %v", err)
		var interlayer struct {
			FileName     string   `json:"filename"`
			Content      string   `json:"content"`
			Dependencies []string `json:"dependencies"`
		}
		if err := json.Unmarshal(rawResponse, &interlayer); err != nil {
			logger.Debug("Interlayer unmarshal failed: %v", err)
			return "", newMessages, fmt.Errorf("couldn't process response: %w", err)
		}
		logger.Debug("Interlayer unmarshal successful")

		response = models.GenerateTestReturn{
			FileName:     interlayer.FileName,
			Content:      interlayer.Content,
			Dependencies: interlayer.Dependencies,
		}
	}

	logger.Debug("Validating response")
	if err := response.Validate(); err != nil {
		logger.Debug("Response validation failed: %v", err)
		return "", newMessages, fmt.Errorf("validation failed: %w", err)
	}
	logger.Debug("Response validation successful")

	logger.Debug("Write the test file to a temporary file")
	filePath := filepath.Join(testsDir, response.FileName)

	// Write the test file
	if err := os.WriteFile(filePath, []byte(response.Content), 0644); err != nil {
		return "", newMessages, fmt.Errorf("couldn't write test file: %w", err)
	}

	// TODO: we need to do somethign with the dependencies
	fmt.Println("Dependencies:")
	for _, dep := range response.Dependencies {
		fmt.Printf("  - %s\n", dep)
	}

	logger.Debug("GenerateTest completed successfully")
	return response.FileName, newMessages, nil
}

func evaluateTestFile(ctx context.Context, client *llm.Client, filename string, tempDir string, testsDir string) (string, bool, error) {
	// List the provided test file
	filePath := filepath.Join(testsDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", false, fmt.Errorf("couldn't read test file: %w", err)
	}

	// Run pnpm test
	testCmd := exec.Command("pnpm", "test", filePath)
	testCmd.Dir = tempDir
	fmt.Printf("Running tests in %s...\n", tempDir)
	output, err := testCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Tests failed!\n")
		fmt.Printf("Error executing pnpm test: %v\n", err)
		fmt.Printf("Test output: %s\n", output)
		// but we don't want to return, we want to continue the loop
	} else {
		fmt.Printf("✅ Tests passed successfully!\n")
	}

	// Analyze the test output
	var builder strings.Builder

	builder.WriteString("INPUTS: \n")
	builder.WriteString("\nTEST FILE NAME: ")
	builder.WriteString(filename)
	builder.WriteString("\nTEST FILE CONTENTS: ")
	builder.WriteString(string(content))
	builder.WriteString("\nTEST OUTPUT: ")
	builder.WriteString(string(output))
	builder.WriteString("\n---END PAGE---\n\n")

	context := builder.String()

	basePrompt := `You are a test engineer, your task is to evaluate the test file and provide feedback on the test file.
Your feedback should be concise and to the point. You should provide feedback on the following:
- Focus mainly on fixing the failing tests.
- The only allowed dependency is @playwright/test, no other dependencies are allowed.
- Whether the test file is covering the provided criteria
- The length of the test file should be around 100 lines of code, the closer the better.
- Whether the test scope is too broad. If the test file is more than 100 lines of code, it is too broad, so suggest what tests to remove (prioritize removing the tests that are failing)


IMPORTANT: When providing the feedback, ensure proper JSON formatting:
1. The "passed" field must be a boolean and cannot be empty. You should return "true" if the test file is good enough and the file does not need more work and "false" otherwise.
   Remember, "passed" cannot be "true" if the test is not passing.
2. The "feedback" field must be a string and cannot be empty if "passed" is "false". Here, you should write your feedback on the test file.
3. Do not include newlines or any characters that would need to be escaped in the "feedback" field.
4. Example of correct format:
   {
	"passed": "false",
	"feedback": "The test file is not covering the footer interactions and is also failing on line 47. Try adjusting the selector for social media links to fix the test."
   }

Invalid formatting will cause errors in processing your response.
`

	// INFO: for a structured response the client requires tools, ref: https://docs.anthropic.com/en/docs/build-with-claude/tool-use/overview
	tool, toolChoice := llm.GenerateTool[models.EvaluationReturn]("get_generate_feedback_return", "")

	rawResponse, err := client.GetStructuredCompletion(
		ctx,
		context,
		basePrompt,
		tool,
		toolChoice,
		[]anthropic.MessageParam{},
	)
	if err != nil {
		return "", false, fmt.Errorf("couldn't process request: %w", err)
	}

	var response models.EvaluationReturn
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return "", false, fmt.Errorf("couldn't unmarshal response: %w", err)
	}

	if response.Passed {
		fmt.Printf("✅ Evaluator accepted the test file.\n")
		return "", true, nil
	}

	fmt.Printf("❌ Evaluator rejected the test file.\n")
	fmt.Printf("Feedback: %s\n", response.Feedback)

	return response.Feedback, false, nil
}

// SetupTestEnvironment creates a temporary directory and copies the config files to it
// returns the temp directory and the tests directory (which is just a subdirectory /tests in the temp directory)
func SetupTestEnvironment(ctx context.Context) (string, string, error) {
	// Create a temporary directory to store the test files
	tempDir, err := os.MkdirTemp("", "playwright-tests-")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		return "", "", fmt.Errorf("couldn't create temporary directory: %w", err)
	}

	logger.Debug("Created temporary directory at: %s", tempDir)

	// Create a tests directory within the temporary directory
	testsDir := filepath.Join(tempDir, "tests")
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		fmt.Printf("Error creating tests directory: %v\n", err)
		return tempDir, "", fmt.Errorf("couldn't create tests directory: %w", err)
	}
	logger.Debug("Created tests directory at: %s", testsDir)

	// Get the absolute path to the src directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return tempDir, testsDir, fmt.Errorf("couldn't get current directory: %w", err)
	}

	templatePath := filepath.Join(currentDir, "internal/repository/gen_eval_loop/nodeTemplate")

	// Files to copy
	filesToCopy := []string{"tsconfig.json", "pnpm-lock.yaml", "package.json", "playwright.config.ts"}

	// Copy each file from template to temp directory
	for _, file := range filesToCopy {
		srcFile := filepath.Join(templatePath, file)
		dstFile := filepath.Join(tempDir, file)

		src, err := os.Open(srcFile)
		if err != nil {
			fmt.Printf("Error opening source file %s: %v\n", file, err)
			return tempDir, testsDir, fmt.Errorf("couldn't open source file %s: %w", file, err)
		}
		defer src.Close()

		dst, err := os.Create(dstFile)
		if err != nil {
			fmt.Printf("Error creating destination file %s: %v\n", file, err)
			return tempDir, testsDir, fmt.Errorf("couldn't create destination file %s: %w", file, err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			fmt.Printf("Error copying file %s: %v\n", file, err)
			return tempDir, testsDir, fmt.Errorf("couldn't copy file %s: %w", file, err)
		}
	}

	// Run pnpm install
	installCmd := exec.Command("pnpm", "i")
	installCmd.Dir = tempDir
	fmt.Printf("Running pnpm install in %s...\n", tempDir)
	output, err := installCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Installation failed!\n")
		fmt.Printf("Error executing pnpm install: %v\n", err)
		fmt.Printf("Installation output: %s\n", output)
		return tempDir, testsDir, fmt.Errorf("couldn't execute pnpm install: %w", err)
	}
	fmt.Printf("✅ Installation completed successfully!\n")

	// Install browsers
	playwrightCmd := exec.Command("npx", "playwright", "install")
	playwrightCmd.Dir = tempDir
	fmt.Printf("Running npx playwright install in %s...\n", tempDir)
	output, err = playwrightCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Installation failed!\n")
		fmt.Printf("Error executing npx playwright install: %v\n", err)
		fmt.Printf("Installation output: %s\n", output)
		return tempDir, testsDir, fmt.Errorf("couldn't install playwright: %w", err)
	}
	fmt.Printf("✅ Installation completed successfully!\n")

	return tempDir, testsDir, nil
}
