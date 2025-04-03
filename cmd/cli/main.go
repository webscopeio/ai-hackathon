package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/repository/analyze"
	"github.com/webscopeio/ai-hackathon/internal/repository/generate"
	"github.com/webscopeio/ai-hackathon/internal/repository/validate"
)

var (
	url string
)

var rootCmd = &cobra.Command{
	Use:   "testbuddy",
	Short: "TestBuddy CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := analyze.Analyze()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Playwright tests for a website",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Error: URL is required")
			return
		}

		// Initialize config and LLM client
		cfg := config.Load()
		client := llm.New(cfg)

		// Generate tests
		logger.Debug("Generating tests for URL: %s", url)
		response, err := generate.GenerateTests(context.Background(), client, url)
		if err != nil {
			fmt.Printf("Error generating tests: %v\n", err)
			return
		}

		// Create a temporary directory to store the test files
		logger.Debug("Creating temporary directory for test files")
		tempDir, err := os.MkdirTemp("", "playwright-tests-")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return
		}
		logger.Debug("Created temporary directory at: %s", tempDir)

		// Create a tests directory within the temporary directory
		testsDir := filepath.Join(tempDir, "tests")
		if err := os.MkdirAll(testsDir, 0755); err != nil {
			fmt.Printf("Error creating tests directory: %v\n", err)
			return
		}
		logger.Debug("Created tests directory at: %s", testsDir)

		// Store each test file in the tests directory
		logger.Debug("Writing %d test files to tests directory", len(response.TestFiles))
		for i, testFile := range response.TestFiles {
			filePath := filepath.Join(testsDir, testFile.Filename)
			logger.Debug("Writing test file %d: %s (content length: %d)", i+1, filePath, len(testFile.Content))
			if err := os.WriteFile(filePath, []byte(testFile.Content), 0644); err != nil {
				fmt.Printf("Error writing test file %s: %v\n", testFile.Filename, err)
				return
			}
			// Update the file path in the response
			response.TestFiles[i].FilePath = filePath
			logger.Debug("Successfully wrote test file: %s", filePath)
		}

		// TODO: we need to do somethign with the dependencies
		fmt.Printf("Successfully generated %d test files in directory: %s\n", len(response.TestFiles), tempDir)
		fmt.Println("Dependencies:")
		for _, dep := range response.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}

		analysis, err := validate.Validate(context.Background(), client, tempDir)
		if err != nil {
			fmt.Printf("Error validating tests: %v\n", err)
			return
		}
		
		// Print validation results
		if len(analysis.Failures) > 0 {
			fmt.Printf("❌ Test validation found %d failures:\n", len(analysis.Failures))
			for _, failure := range analysis.Failures {
				fmt.Printf("File: %s\nError: %s\n\n", failure.Filename, failure.Error)
			}
		} else {
			fmt.Printf("✅ All tests passed validation successfully!\n")
		}

	},
}

func init() {
	generateCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the website to generate tests for")
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
