package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
	"github.com/webscopeio/ai-hackathon/internal/repository/analyzer"
	"github.com/webscopeio/ai-hackathon/internal/repository/gen_eval_loop"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func runPipeline(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		sendMessage(c, mt, string(message))
	}
}

func sendMessage(c *websocket.Conn, mt int, prompt string) {
	// Add context
	ctx := context.Background()

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

		IMPORTANT: Pass all the criteria into the get_final_criteria_tool. Format each criterion as follows:

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

	basePrompt += `\n\IMPORTANT: This is the input from the user. Use it to generate the test criteria:` + prompt

	analysis, err := analyzer.Analyze(ctx, cfg, client, "https://ai-hackathon-demo-delta.vercel.app/", basePrompt)
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

	c.WriteMessage(mt, []byte(fmt.Sprintf("\n[MAIN FLOW] Analyzer generated %d scenarios\n", len(criteria))))
	// print the criteria line by line
	logger.Debug("CRITERIA LENGTH: %d", len(criteria))
	for _, c := range criteria {
		logger.Debug("CRITERIA: %s", c)
	}

	noOfLoops := 6

	for i, c := range criteria {
		fmt.Printf("\n[MAIN FLOW] Generating test for scenario %d: %s\n", i, c)
		filename, err := gen_eval_loop.GenEvalLoop(ctx, client, &models.AnalyzerReturn{
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
}

func main() {
	http.HandleFunc("/ws", runPipeline)

	addr := fmt.Sprintf(":%s", "8080")
	fmt.Printf("Server starting on localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
