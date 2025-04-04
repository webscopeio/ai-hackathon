package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

// getUmamiAPIKey returns the API key for Umami API authentication
func getUmamiAPIKey(ctx context.Context, cfg *config.Config) (string, error) {
	logger.Debug("Getting Umami API key")

	// Check if API key is configured
	if cfg.UmamiAPIKey == "" {
		return "", fmt.Errorf("Umami API key not configured")
	}

	logger.Debug("Using Umami API key for authentication")
	return cfg.UmamiAPIKey, nil
}

// getUmamiSessions retrieves user sessions from the Umami API
func getUmamiSessions(ctx context.Context, cfg *config.Config, apiKey string, startDate, endDate time.Time) ([]models.UmamiSession, error) {
	logger.Debug("Getting Umami sessions for website ID: %s", cfg.UmamiWebsiteId)

	// Check if website ID is configured
	if cfg.UmamiWebsiteId == "" {
		return nil, fmt.Errorf("Umami website ID not configured")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/websites/%s/sessions", cfg.UmamiURL, cfg.UmamiWebsiteId)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with API key
	// Add API key header as per Umami Cloud API documentation
	req.Header.Add("x-umami-api-key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Add query parameters
	startAtValue := fmt.Sprintf("%d", startDate.Unix()*1000)
	endAtValue := fmt.Sprintf("%d", endDate.Unix()*1000)
	logger.Debug("Using startAt=%s (%s) and endAt=%s (%s) for sessions",
		startAtValue, startDate.Format(time.RFC3339),
		endAtValue, endDate.Format(time.RFC3339))

	query := url.Values{}
	query.Add("startAt", startAtValue)
	query.Add("endAt", endAtValue)
	req.URL.RawQuery = query.Encode()

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response
	var response models.UmamiSessionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	logger.Debug("Retrieved %d Umami sessions", len(response.Data))
	return response.Data, nil
}

// getSessionPageViews retrieves activity for a specific session
func getSessionPageViews(ctx context.Context, cfg *config.Config, apiKey, sessionId string) ([]models.UmamiSessionActivity, error) {
	logger.Debug("Getting activity for session ID: %s", sessionId)

	// Construct the API URL with the correct endpoint
	apiURL := fmt.Sprintf("%s/websites/%s/sessions/%s/activity", cfg.UmamiURL, cfg.UmamiWebsiteId, sessionId)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters - the activity endpoint requires startAt and endAt
	now := time.Now()
	// Get data for the last 90 days
	startDate := now.AddDate(0, 0, -90)
	
	startAtValue := fmt.Sprintf("%d", startDate.Unix()*1000)
	endAtValue := fmt.Sprintf("%d", now.Unix()*1000)
	logger.Debug("Using startAt=%s (%s) and endAt=%s (%s) for session activity",
		startAtValue, startDate.Format(time.RFC3339),
		endAtValue, now.Format(time.RFC3339))
	
	query := url.Values{}
	query.Add("startAt", startAtValue)
	query.Add("endAt", endAtValue)
	req.URL.RawQuery = query.Encode()

	// Add authorization header with API key
	// Add API key header as per Umami Cloud API documentation
	req.Header.Add("x-umami-api-key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response - the activity endpoint returns an array directly
	var activities []models.UmamiSessionActivity
	err = json.Unmarshal(body, &activities)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Log the total number of activities received from the API
	logger.Debug("Total activities received from API: %d", len(activities))

	// Print the first activity for debugging if available
	if len(activities) > 0 {
		activityJSON, _ := json.MarshalIndent(activities[0], "", "  ")
		logger.Debug("First activity sample: %s", string(activityJSON))
	}

	logger.Debug("Retrieved %d activities for session %s", len(activities), sessionId)
	return activities, nil
}

// extractPathFromURL extracts the path component from a URL
func extractPathFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr // Return the original string if parsing fails
	}

	// Return the path component, or "/" if empty
	if parsedURL.Path == "" {
		return "/"
	}

	return parsedURL.Path
}

// buildUserFlows constructs user flows from sessions and their activities
func buildUserFlows(sessions []models.UmamiSession, sessionActivities map[string][]models.UmamiSessionActivity) []models.UmamiUserFlow {
	var userFlows []models.UmamiUserFlow

	for _, session := range sessions {
		activities, ok := sessionActivities[session.ID]
		if !ok || len(activities) == 0 {
			continue // Skip sessions with no activities
		}

		// Sort activities by creation time
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].CreatedAt < activities[j].CreatedAt
		})

		// Extract paths from URLPath
		paths := make([]string, 0, len(activities))
		for _, activity := range activities {
			// The URLPath already contains the path, no need to extract it
			paths = append(paths, activity.URLPath)
		}

		// Create user flow
		userFlow := models.UmamiUserFlow{
			SessionID: session.ID,
			Path:      paths,
		}

		userFlows = append(userFlows, userFlow)
	}

	return userFlows
}

// findSignificantFlows identifies common patterns in user flows
func findSignificantFlows(userFlows []models.UmamiUserFlow, minPathLength, minFrequency int) []models.UmamiSignificantFlow {
	// Count occurrences of each path pattern
	pathCounts := make(map[string]int)
	pathArrays := make(map[string][]string)

	for _, flow := range userFlows {
		// Skip flows that are too short
		if len(flow.Path) < minPathLength {
			continue
		}

		// Consider all possible sub-paths of sufficient length
		for i := 0; i <= len(flow.Path)-minPathLength; i++ {
			for j := i + minPathLength; j <= len(flow.Path); j++ {
				subPath := flow.Path[i:j]
				pathKey := strings.Join(subPath, "|")
				pathCounts[pathKey]++
				pathArrays[pathKey] = subPath
			}
		}
	}

	// Convert to significant flows
	var significantFlows []models.UmamiSignificantFlow
	totalFlows := len(userFlows)

	for pathKey, count := range pathCounts {
		// Skip paths that don't meet the minimum frequency
		if count < minFrequency {
			continue
		}

		percentage := float64(count) / float64(totalFlows) * 100.0

		significantFlows = append(significantFlows, models.UmamiSignificantFlow{
			Path:       pathArrays[pathKey],
			Frequency:  count,
			Percentage: percentage,
		})
	}

	// Sort by frequency (descending)
	sort.Slice(significantFlows, func(i, j int) bool {
		return significantFlows[i].Frequency > significantFlows[j].Frequency
	})

	return significantFlows
}

// GetSignificantUserFlows retrieves user sessions from Umami and identifies significant user flows
func GetSignificantUserFlows(ctx context.Context, cfg *config.Config, daysBack int, minPathLength, minFrequency int) ([]models.UmamiSignificantFlow, error) {
	logger.Debug("Getting significant user flows from Umami for the last %d days", daysBack)

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -daysBack)

	// Get Umami API key
	apiKey, err := getUmamiAPIKey(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get Umami API key: %w", err)
	}

	// Get sessions
	sessions, err := getUmamiSessions(ctx, cfg, apiKey, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get Umami sessions: %w", err)
	}

	if len(sessions) == 0 {
		logger.Debug("No sessions found in the specified date range")
		return []models.UmamiSignificantFlow{}, nil
	}

	// Get activities for each session
	sessionActivities := make(map[string][]models.UmamiSessionActivity)
	for _, session := range sessions {
		activities, err := getSessionPageViews(ctx, cfg, apiKey, session.ID)
		if err != nil {
			logger.Debug("Failed to get activities for session %s: %v", session.ID, err)
			continue
		}
		sessionActivities[session.ID] = activities
	}

	// Build user flows
	userFlows := buildUserFlows(sessions, sessionActivities)
	logger.Debug("Built %d user flows from sessions", len(userFlows))

	// Find significant flows
	significantFlows := findSignificantFlows(userFlows, minPathLength, minFrequency)
	logger.Debug("Identified %d significant user flows", len(significantFlows))

	return significantFlows, nil
}

// GetUserPaths retrieves user sessions from Umami and returns user paths as lists of URLs
func GetUserPaths(ctx context.Context, cfg *config.Config, daysBack int) ([][]string, error) {
	logger.Debug("Getting user paths from Umami for the last %d days", daysBack)

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -daysBack)

	// Get Umami API key
	apiKey, err := getUmamiAPIKey(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get Umami API key: %w", err)
	}

	// Get sessions
	sessions, err := getUmamiSessions(ctx, cfg, apiKey, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get Umami sessions: %w", err)
	}

	if len(sessions) == 0 {
		logger.Debug("No sessions found in the specified date range")
		return [][]string{}, nil
	}

	// Get activities for each session
	sessionActivities := make(map[string][]models.UmamiSessionActivity)
	for _, session := range sessions {
		activities, err := getSessionPageViews(ctx, cfg, apiKey, session.ID)
		if err != nil {
			logger.Debug("Failed to get activities for session %s: %v", session.ID, err)
			continue
		}
		sessionActivities[session.ID] = activities
	}

	// Build user flows
	userFlows := buildUserFlows(sessions, sessionActivities)

	// Convert to simple path arrays
	paths := make([][]string, 0, len(userFlows))
	for _, flow := range userFlows {
		if len(flow.Path) > 0 {
			paths = append(paths, flow.Path)
		}
	}

	logger.Debug("Retrieved %d user paths", len(paths))
	return paths, nil
}
