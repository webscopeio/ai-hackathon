package main

import (
	"fmt"
	"os"
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

		analysis, err := analyzer.Analyze(cmd.Context(), cfg, client, "https://jakub.kr/", `You are an analyst, your task is to analyze the provided website and generate a list of criteria that can be used by another agent to generate E2E tests.

The criteria should be:
- Short and concise
- It must be a list of criteria separated by 2 newlines
- Cover the most important parts of the website
- Be easy to understand
- Be easy to test
- Focus on simple tests that are easy to write, we can interate later with more complex tests


IMPORTANT: The critera must be a list of criteria always separated by 2 newlines.
`)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		logger.Debug("CRITERIA FROM THE ANALYZER: %v", analysis.Criteria)
		logger.Debug("LENGTH OF THE CRITERIA: %d", len(analysis.Criteria))

		if len(analysis.Criteria) == 0 {
			fmt.Println("Error: No test criteria were generated from the analysis")
			return
		}

		// Split criteria by double newlines
		criteria := strings.Split(analysis.Criteria, "\n\n")

		// print the criteria line by line
		fmt.Println("CRITERIA LENGTH: ", len(criteria))
		for _, c := range criteria {
			fmt.Println("CRITERIA: ", c)
		}

		for i, c := range criteria {
			fmt.Printf("CRITERIA %d: %s\n", i, c)
			filename, err := gen_eval_loop.GenEvalLoop(cmd.Context(), client, &models.AnalyzerReturn{
				TechSpec:   analysis.TechSpec,
				ContentMap: analysis.ContentMap,
				Criteria:   analysis.Criteria,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("Generated test file: %s\n", filename)
			// copy the file to the current directory
			err = os.Rename(filename, "./__generated__")
			if err != nil {
				fmt.Printf("Error copying file: %v\n", err)
				return
			}
			fmt.Printf("Copied test file to current directory: %s\n", "./__generated__")
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
