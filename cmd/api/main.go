package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	// Create a cancellable context
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Create a WaitGroup to wait for sendMessage to complete
	var wg sync.WaitGroup

	// Channel to signal connection closure
	done := make(chan struct{})
	// Channel to control pause/resume
	pauseChan := make(chan bool)

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error reading message: %v", err)
				}
				return
			}
			log.Printf("recv: %s", message)

			// Handle pause/resume commands
			if string(message) == "PAUSE" {
				pauseChan <- true
				c.WriteMessage(mt, []byte("STATUS Paused"))
				continue
			} else if string(message) == "RESUME" {
				pauseChan <- false
				c.WriteMessage(mt, []byte("STATUS Resumed"))
				continue
			}

			// Start sendMessage in a goroutine
			wg.Add(1)
			go func(msg []byte) {
				defer wg.Done()
				sendMessage(ctx, c, mt, string(msg), pauseChan)
			}(message)
		}
	}()

	// Wait for either the connection to close or the context to be cancelled
	<-done
	cancel() // Cancel the context to stop any ongoing operations

	// Wait for any ongoing sendMessage operations to complete
	wg.Wait()
}

func sendMessage(ctx context.Context, c *websocket.Conn, mt int, url string, pauseChan chan bool) {
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

	// Check if context is cancelled before proceeding
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Function to check pause state
	checkPause := func() {
		for {
			select {
			case isPaused := <-pauseChan:
				if isPaused {
					// Wait for resume signal
					for p := range pauseChan {
						if !p {
							return
						}
					}
				}
				return
			default:
				return
			}
		}
	}

	analysis, err := analyzer.Analyze(ctx, cfg, client, url, basePrompt, c, mt)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	checkPause()

	if len(analysis.Criteria) == 0 {
		log.Println("Error: No test criteria were generated from the analysis")
		return
	}

	// Split criteria by double newlines
	criteria := strings.Split(analysis.Criteria, "\n\n")

	// Check context before sending message
	select {
	case <-ctx.Done():
		return
	default:
		if err := c.WriteMessage(mt, []byte(fmt.Sprintf("ANALYZER 'Generated %d scenarios, passing to Generator.'", len(criteria)))); err != nil {
			log.Printf("error writing message: %v", err)
			return
		}
	}

	checkPause()

	c.WriteMessage(mt, []byte(fmt.Sprintf("SCENARIOS %d", len(criteria))))

	logger.Debug("CRITERIA LENGTH: %d", len(criteria))
	for _, criterion := range criteria {
		logger.Debug("CRITERIA: %s", criterion)
	}

	noOfLoops := 6
	filenames := []string{}

	for i, criterion := range criteria {
		checkPause()

		// Check context before each iteration
		select {
		case <-ctx.Done():
			return
		default:
		}

		keys := make([]string, 0, len(analysis.ContentMap))
		for k := range analysis.ContentMap {
			keys = append(keys, k)
		}
		c.WriteMessage(mt, []byte(fmt.Sprintf("SCENARIO %s\n\nContent of these sites is passed to the generator: %s", criterion, strings.Join(keys, ", "))))

		fmt.Printf("\n[MAIN FLOW] Generating test for scenario %d: %s\n", i, criterion)
		c.WriteMessage(mt, []byte(fmt.Sprintf("GENERATOR 'Generating test for scenario %d: %s'", i, criterion)))
		filename, err := gen_eval_loop.GenEvalLoop(ctx, client, &models.AnalyzerReturn{
			TechSpec:   url,
			ContentMap: analysis.ContentMap,
			Criteria:   analysis.Criteria,
		}, i+1, noOfLoops, filenames, c, mt, pauseChan)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		checkPause()

		logger.Debug("[MAIN FLOW] Writing test file: %s\n", filepath.Base(filename))
		c.WriteMessage(mt, []byte(fmt.Sprintf("EVALUATOR 'Writing test file: %s'", filepath.Base(filename))))
		c.WriteMessage(mt, []byte(fmt.Sprintf("FILENAME %s", filepath.Base(filename))))
		testFileName := filepath.Base(filename)[7:]
		filenames = append(filenames, testFileName)
		logger.Debug("[MAIN FLOW] Updated filenames: %v", filenames)

		// copy the file to the current directory
		destPath := filepath.Join("./__generated__", filepath.Base(filename))
		fmt.Printf("\n[MAIN FLOW] Writing generated test file to %s\n", destPath)
		err = os.Rename(filename, destPath)
		if err != nil {
			log.Printf("Error copying file: %v\n", err)
			return
		}
	}
	c.WriteMessage(mt, []byte(fmt.Sprintf("FINAL_MESSAGE 'Tests are generated.'")))
}

func main() {
	http.HandleFunc("/ws", runPipeline)

	addr := fmt.Sprintf(":%s", "8080")
	fmt.Printf("Server starting on localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
