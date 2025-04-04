package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
	"github.com/webscopeio/ai-hackathon/internal/repository/analyzer"
	"github.com/webscopeio/ai-hackathon/internal/repository/gen_eval_loop"
)

var url string

var rootCmd = &cobra.Command{
	Use:   "testbuddy",
	Short: "TestBuddy CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Playwright tests for a website",
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize config and LLM client
		cfg := config.Load()
		client := llm.New(cfg)

		basePrompt := `You are a test planning expert. Your task is to analyze the provided website and generate EXACTLY 4 specific test criteria that can be used by another agent to generate E2E tests.

		The criteria should:
		1. Cover the core functionality of the application
		2. Focus on different user journeys, I am interested in the content of the most valuable pages
		3. Include both happy path and edge case scenarios
		4. Be specific enough to be implemented as end-to-end tests
		5. Be short, concise and easy to understand
		6. Focus on simple tests that are easy to write (we can iterate later with more complex tests)

		IMPORTANT: Format each criterion as follows:

		CRITERION #1:
		TITLE: [Short descriptive title]
		SCENARIO: [Clear description of what should be tested]
		EXPECTED: [Expected outcome or behavior]

		(Repeat for CRITERION #2, #3, and #4)

		Each criterion must be separated by 2 newlines for proper parsing.

		Example:
		CRITERION #1:
		TITLE: User Login Authentication
		SCENARIO: Verify a registered user can successfully log in with valid credentials
		EXPECTED: User should be authenticated and redirected to their personalized dashboard


		CRITERION #2:
		TITLE: Product Search Functionality
		SCENARIO: Verify users can search for products and get relevant results
		EXPECTED: Search results page should display matching products with correct information`

		websiteDescription := "Check out the website, wonder how is it structured?. I am interested in the content of the most valuable pages to create the criteria to generate an E2E tests. My orgSlug := \"webscopeio-pb\" and projectSlug := \"ai-hackathon-demo\" for Sentry, please check the errors in the last 14 days and include them in the analysis."
		basePrompt += `\n\IMPORTANT: You are analyzing the following website:` + websiteDescription

		analysis, err := analyzer.Analyze(cmd.Context(), cfg, client, "https://ai-hackathon-demo-delta.vercel.app/", basePrompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(analysis.Criteria) == 0 {
			fmt.Println("Error: No test criteria were generated from the analysis")
			return
		}

		// Split criteria by double newlines
		criteria := strings.Split(analysis.Criteria, "\n\n")

		fmt.Printf("\n[MAIN FLOW] Analyzer generated %d scenarios\n", len(criteria))
		// print the criteria line by line
		logger.Debug("CRITERIA LENGTH: %d", len(criteria))
		for _, c := range criteria {
			logger.Debug("CRITERIA: %s", c)
		}

		noOfLoops := 6

		for i, c := range criteria {
			fmt.Printf("\n[MAIN FLOW] Generating test for scenario %d: %s\n", i, c)
			filename, err := gen_eval_loop.GenEvalLoop(cmd.Context(), client, &models.AnalyzerReturn{
				TechSpec:   analysis.TechSpec,
				ContentMap: analysis.ContentMap,
				Criteria:   analysis.Criteria,
			}, i+1, noOfLoops)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			logger.Debug("[MAIN FLOW] Writing test file: %s\n", filepath.Base(filename))

			// copy the file to the current directory
			destPath := filepath.Join("./__generated__", filepath.Base(filename))
			fmt.Printf("\n[MAIN FLOW] Writing generated test file to %s\n", destPath)
			err = os.Rename(filename, destPath)
			if err != nil {
				fmt.Printf("Error copying file: %v\n", err)
				return
			}
		}

		return

	},
}

func init() {
	generateCmd.Flags()
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
