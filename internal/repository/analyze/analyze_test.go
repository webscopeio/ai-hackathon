package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
)

func TestAnalyze(t *testing.T) {
	// Create a context with timeout to prevent the test from running too long
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load configuration (including API key from environment variables)
	cfg := config.Load()

	if cfg.APIKey == "" {
		t.Skip("APIKey environment variable not set, skipping test")
	}

	client := llm.New(cfg)

	// Test parameters
	urlStr := "example.com"
	maxDepth := 2
	maxPathSegments := 2
	prompt := "The website is a small presentation page. It might contain a link to another page. Focus the test cases around this."

	// Call the Analyze function
	result, err := Analyze(ctx, urlStr, maxDepth, maxPathSegments, prompt, client)
	// Check for errors
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Print the results
	fmt.Println("Analysis Results:")
	fmt.Println("=================")

	// Print the number of links found
	fmt.Printf("Links found: %d\n", len(result.Links))

	// Print the first few links (if any)
	if len(result.Links) > 0 {
		fmt.Println("Sample links:")
		for i, link := range result.Links {
			if i >= 3 {
				break
			}
			fmt.Printf("  %d. %s\n", i+1, link)
		}
	}

	// Print the number of pages crawled
	fmt.Printf("Pages crawled: %d\n", len(result.Results))

	// Print the analysis
	fmt.Println("\nAnalysis:")
	fmt.Println(result.Analysis)

	// Optionally, print the full result as JSON
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("\nFull result (JSON):")
	fmt.Println(string(resultJSON))
}
