package analyzer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/webscopeio/ai-hackathon/internal/config"
)

func TestGetUmamiSessions(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load config from environment variables
	cfg := config.Load()

	// Skip test if Umami API key is not configured
	if cfg.UmamiAPIKey == "" {
		t.Skip("Skipping test as Umami API key is not configured in .env file")
	}

	// Skip test if Umami website ID is not configured
	if cfg.UmamiWebsiteId == "" {
		t.Skip("Skipping test as Umami website ID is not configured in .env file")
	}

	// Calculate date range for the last 7 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// Get the API key
	apiKey, err := getUmamiAPIKey(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Failed to get Umami API key: %v", err)
	}

	// Call the function with context, config, and date range
	sessions, err := getUmamiSessions(context.Background(), cfg, apiKey, startDate, endDate)
	if err != nil {
		t.Logf("Note: This test requires valid Umami credentials to pass")
		t.Fatalf("Error getting Umami sessions: %v", err)
	}

	// Print the results
	t.Logf("Successfully retrieved %d Umami sessions", len(sessions))

	// Print details of up to 3 sessions
	for i, session := range sessions {
		if i >= 3 {
			break
		}

		sessionJSON, _ := json.MarshalIndent(session, "", "  ")
		t.Logf("Session %d:\n%s", i+1, string(sessionJSON))

		// Get activities for this session
		activities, err := getSessionPageViews(context.Background(), cfg, apiKey, session.ID)
		if err != nil {
			t.Logf("Error getting activities for session %s: %v", session.ID, err)
			continue
		}

		t.Logf("Retrieved %d activities for session %d", len(activities), i+1)

		// Print details of up to 3 activities
		for j, activity := range activities {
			if j >= 3 {
				break
			}

			activityJSON, _ := json.MarshalIndent(activity, "", "  ")
			t.Logf("Activity %d.%d:\n%s", i+1, j+1, string(activityJSON))
		}
	}
}

func TestGetUserPaths(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load config from environment variables
	cfg := config.Load()

	// Skip test if Umami API key is not configured
	if cfg.UmamiAPIKey == "" {
		t.Skip("Skipping test as Umami API key is not configured in .env file")
	}

	// Skip test if Umami website ID is not configured
	if cfg.UmamiWebsiteId == "" {
		t.Skip("Skipping test as Umami website ID is not configured in .env file")
	}

	// Call the function with context, config, and days back
	userPaths, err := GetUserPaths(context.Background(), cfg, 7) // Get paths for the last 7 days
	if err != nil {
		t.Logf("Note: This test requires valid Umami credentials to pass")
		t.Fatalf("Error getting user paths: %v", err)
	}

	// Print the results
	t.Logf("Successfully retrieved %d user paths", len(userPaths))

	// Print details of up to 5 user paths
	for i, path := range userPaths {
		if i >= 5 {
			break
		}

		t.Logf("Path %d: %v", i+1, path)
	}
}

func TestGetSignificantUserFlows(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load config from environment variables
	cfg := config.Load()

	// Skip test if Umami API key is not configured
	if cfg.UmamiAPIKey == "" {
		t.Skip("Skipping test as Umami API key is not configured in .env file")
	}

	// Skip test if Umami website ID is not configured
	if cfg.UmamiWebsiteId == "" {
		t.Skip("Skipping test as Umami website ID is not configured in .env file")
	}

	// Call the function with context, config, and parameters
	significantFlows, err := GetSignificantUserFlows(context.Background(), cfg, 7, 2, 2) // 7 days back, min path length 2, min frequency 2
	if err != nil {
		t.Logf("Note: This test requires valid Umami credentials to pass")
		t.Fatalf("Error getting significant user flows: %v", err)
	}

	// Print the results
	t.Logf("Successfully retrieved %d significant user flows", len(significantFlows))

	// Print details of all significant flows
	for i, flow := range significantFlows {
		flowJSON, _ := json.MarshalIndent(flow, "", "  ")
		t.Logf("Flow %d (frequency: %d):\n%s", i+1, flow.Frequency, string(flowJSON))
	}
}
