package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Validate(ctx context.Context, client *llm.Client, tempDir string) (models.TestRunAnalysis, error) {
	// Get the absolute path to the src directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return models.TestRunAnalysis{}, err
	}

	templatePath := filepath.Join(currentDir, "internal/repository/validate/nodeTemplate")

	// Files to copy
	filesToCopy := []string{"tsconfig.json", "pnpm-lock.yaml", "package.json", "playwright.config.ts"}

	// Copy each file from template to temp directory
	for _, file := range filesToCopy {
		srcFile := filepath.Join(templatePath, file)
		dstFile := filepath.Join(tempDir, file)

		src, err := os.Open(srcFile)
		if err != nil {
			fmt.Printf("Error opening source file %s: %v\n", file, err)
			return models.TestRunAnalysis{}, err
		}
		defer src.Close()

		dst, err := os.Create(dstFile)
		if err != nil {
			fmt.Printf("Error creating destination file %s: %v\n", file, err)
			return models.TestRunAnalysis{}, err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			fmt.Printf("Error copying file %s: %v\n", file, err)
			return models.TestRunAnalysis{}, err
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
		return models.TestRunAnalysis{}, err
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
		return models.TestRunAnalysis{}, err
	}
	fmt.Printf("✅ Installation completed successfully!\n")

	// Run pnpm test
	testCmd := exec.Command("pnpm", "test")
	testCmd.Dir = tempDir
	fmt.Printf("Running tests in %s...\n", tempDir)
	output, err = testCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Tests failed!\n")
		fmt.Printf("Error executing pnpm test: %v\n", err)
		fmt.Printf("Test output: %s\n", output)

		// Even though tests failed, we want to analyze the output
		analysis, analyzeErr := analyzeOutput(ctx, client, output)
		if analyzeErr != nil {
			fmt.Printf("Error analyzing test output: %v\n", analyzeErr)
			return models.TestRunAnalysis{}, fmt.Errorf("test execution failed and output analysis failed: %w", analyzeErr)
		}
		return analysis, nil
	}

	fmt.Printf("✅ Tests passed successfully!\n")
	fmt.Printf("Test output: %s\n", output)

	// Even for successful runs, analyze the output to catch any warnings or informational messages
	return analyzeOutput(ctx, client, output)
}

func analyzeOutput(ctx context.Context, client *llm.Client, output []byte) (models.TestRunAnalysis, error) {
	analysis := models.TestRunAnalysis{}

	// Create a prompt for the LLM to analyze the test output
	prompt := `Please analyze the following Playwright test output and extract information about any test failures.
For each failure, provide:
1. The filename where the failure occurred
2. The specific error message or reason for failure

Format the response as a JSON object with an array of failures, where each failure has:
- filename: the test file where the failure occurred
- error: the error message or reason for failure

Test output:
` + string(output)

	// Create a tool for structured response
	tool, toolChoice := llm.GenerateTool[models.TestRunAnalysis]("get_test_analysis", "Analyze Playwright test output and return structured information about test failures")

	// Get structured completion from LLM
	rawResponse, err := client.GetStructuredCompletion(
		ctx,
		"", // No additional context needed
		prompt,
		tool,
		toolChoice,
	)
	if err != nil {
		return analysis, fmt.Errorf("couldn't process test output: %w", err)
	}

	// Parse the response
	if err := json.Unmarshal(rawResponse, &analysis); err != nil {
		return analysis, fmt.Errorf("couldn't parse LLM response: %w", err)
	}

	return analysis, nil
}
