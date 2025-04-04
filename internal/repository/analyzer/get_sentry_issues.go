package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/webscopeio/ai-hackathon/internal/config"
	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

// GetSentryIssues retrieves issues from a Sentry project using the Sentry API
// It requires organization/project slugs and gets the auth token from config
func GetSentryIssues(ctx context.Context, cfg *config.Config, orgSlug, projectSlug string) ([]models.SentryIssue, error) {
	logger.Debug("Getting Sentry issues for project: %s/%s", orgSlug, projectSlug)

	// Get auth token from config
	authToken := cfg.SentryAuthToken
	if authToken == "" {
		return nil, fmt.Errorf("sentry auth token not configured")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("https://sentry.io/api/0/projects/%s/%s/issues/", orgSlug, projectSlug)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")

	// Add query parameters if needed
	query := url.Values{}
	query.Add("statsPeriod", "14d")     // Get issues from the last 14 days
	query.Add("query", "is:unresolved") // Only get unresolved issues
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
	var issues []models.SentryIssue
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	logger.Debug("Retrieved %d Sentry issues", len(issues))
	return issues, nil
}

// GetSentryIssueTagDetails retrieves details for a specific tag of an issue
func GetSentryIssueTagDetails(ctx context.Context, cfg *config.Config, orgSlug, issueID, tagKey string) (*models.SentryTagDetails, error) {
	logger.Debug("Getting tag details for Sentry issue: %s, tag: %s in organization: %s", issueID, tagKey, orgSlug)

	// Get auth token from config
	authToken := cfg.SentryAuthToken
	if authToken == "" {
		return nil, fmt.Errorf("sentry auth token not configured")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("https://sentry.io/api/0/organizations/%s/issues/%s/tags/%s/", orgSlug, issueID, tagKey)
	logger.Debug("Requesting URL: %s", apiURL)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	logger.Debug("Sending request to Sentry API for issue tag details...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	logger.Debug("Received response with status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Debug("Error response body: %s", string(body))
		return nil, fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response body for debugging
	logger.Debug("Response body: %s", string(body))

	// Unmarshal the JSON response
	var tagDetails models.SentryTagDetails
	err = json.Unmarshal(body, &tagDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	logger.Debug("Retrieved tag details for issue %s, tag: %s with %d unique values",
		issueID, tagKey, tagDetails.UniqueValues)
	return &tagDetails, nil
}

// GetSentryIssueTagValuesSorted retrieves all unique values for a specific tag of an issue,
// sorted by occurrence count (most frequent first)
func GetSentryIssueTagValuesSorted(ctx context.Context, cfg *config.Config, orgSlug, issueID, tagKey string) ([]models.SentryTagValueSorted, error) {
	logger.Debug("Getting sorted tag values for Sentry issue: %s, tag: %s", issueID, tagKey)

	// First get the tag details
	tagDetails, err := GetSentryIssueTagDetails(ctx, cfg, orgSlug, issueID, tagKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag details: %w", err)
	}

	// Create a slice to hold the sorted values
	sortedValues := make([]models.SentryTagValueSorted, 0, len(tagDetails.TopValues))

	// Convert the tag values to our sorted format
	for _, value := range tagDetails.TopValues {
		sortedValues = append(sortedValues, models.SentryTagValueSorted{
			Value: value.Value,
			Count: value.Count,
		})
	}

	// The values should already be sorted by count in the API response,
	// but we'll sort them again just to be sure
	sort.Slice(sortedValues, func(i, j int) bool {
		return sortedValues[i].Count > sortedValues[j].Count
	})

	logger.Debug("Retrieved and sorted %d tag values for issue %s, tag: %s",
		len(sortedValues), issueID, tagKey)
	return sortedValues, nil
}

// GetAffectedSentryPaths retrieves a list of URL paths affected by errors from Sentry
// The paths are sorted by occurrence count (most frequent first)
func GetAffectedSentryPaths(ctx context.Context, cfg *config.Config, orgSlug, projectSlug string) ([]models.SentryAffectedPath, error) {
	logger.Debug("Getting affected URL paths from Sentry for organization: %s, project: %s", orgSlug, projectSlug)

	// First get the issues
	issues, err := GetSentryIssues(ctx, cfg, orgSlug, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	logger.Debug("Retrieved %d issues, now fetching URL tag details for each", len(issues))
	if len(issues) == 0 {
		logger.Debug("No issues found to fetch URL paths for")
		return []models.SentryAffectedPath{}, nil
	}

	// Limit to first 20 issues or less
	maxIssues := 20
	if len(issues) < maxIssues {
		maxIssues = len(issues)
	}

	// Map to aggregate URL paths and their counts across all issues
	pathCountMap := make(map[string]int)

	// For each issue, get URL tag details
	for i, issue := range issues[:maxIssues] {
		logger.Debug("Processing issue %d/%d: ID=%s, ShortID=%s", i+1, maxIssues, issue.ID, issue.ShortID)

		// Get URL tag details for this issue
		tagDetails, err := GetSentryIssueTagDetails(ctx, cfg, orgSlug, issue.ID, "url")
		if err != nil {
			logger.Debug("Failed to get URL tag details for issue %s: %v", issue.ID, err)
			// Continue with the next issue even if this one fails
			continue
		}

		logger.Debug("Retrieved %d URL values for issue %s", len(tagDetails.TopValues), issue.ID)

		// Add the URL values to our aggregate map
		for _, value := range tagDetails.TopValues {
			pathCountMap[value.Value] += value.Count
		}
	}

	// Convert the map to a slice for sorting
	affectedPaths := make([]models.SentryAffectedPath, 0, len(pathCountMap))
	for path, count := range pathCountMap {
		affectedPaths = append(affectedPaths, models.SentryAffectedPath{
			Path:  path,
			Count: count,
		})
	}

	// Sort the paths by count in descending order
	sort.Slice(affectedPaths, func(i, j int) bool {
		return affectedPaths[i].Count > affectedPaths[j].Count
	})

	logger.Debug("Retrieved and aggregated %d affected URL paths across %d issues",
		len(affectedPaths), maxIssues)
	return affectedPaths, nil
}
