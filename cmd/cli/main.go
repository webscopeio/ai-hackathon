package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
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

		analysis := models.AnalyzerReturn{
			TechSpec: "A simple one site website that should serve as an example of a website",
			SiteMap: map[string]string{
				"https://example.com": `<body><div><h1>Example Domain</h1><p>This domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.</p><p><a href="https://www.iana.org/domains/example">More information...</a></p></div></body>`,
			},
			Criteria: "Test Criteria:\n1. Verify that the main heading displays 'Example Domain'\n2. Check if the informational paragraph about domain usage is present\n3. Ensure the 'More information' link points to iana.org and is clickable\n4. Validate that all text content is properly rendered",
		}

		filename, err := gen_eval_loop.GenEvalLoop(cmd.Context(), client, &analysis)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Generated test file: %s\n", filename)
		// copy the file to the current directory
		err = os.Rename(filename, "./generate.spec.ts")
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return
		}
		fmt.Printf("Copied test file to current directory: %s\n", "./__generated__")
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
