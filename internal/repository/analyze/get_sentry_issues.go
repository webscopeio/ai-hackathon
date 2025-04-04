package analyze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/webscopeio/ai-hackathon/internal/logger"
)

// SentryIssue represents a single issue from Sentry
type SentryIssue struct {
	ID                  string                 `json:"id"`
	ShortID             string                 `json:"shortId"`
	Title               string                 `json:"title"`
	Culprit             string                 `json:"culprit"`
	Level               string                 `json:"level"`
	Status              string                 `json:"status"`
	FirstSeen           string                 `json:"firstSeen"`
	LastSeen            string                 `json:"lastSeen"`
	Count               string                 `json:"count"`
	UserCount           int                    `json:"userCount"`
	Permalink           string                 `json:"permalink"`
	Type                string                 `json:"type"`
	Metadata            Metadata               `json:"metadata"`
	SubscriptionDetails *SubscriptionDetails   `json:"subscriptionDetails,omitempty"`
	Logger              *string                `json:"logger"`
	NumComments         int                    `json:"numComments"`
	IsPublic            bool                   `json:"isPublic"`
	HasSeen             bool                   `json:"hasSeen"`
	ShareID             *string                `json:"shareId"`
	IsSubscribed        bool                   `json:"isSubscribed"`
	IsBookmarked        bool                   `json:"isBookmarked"`
	Project             Project                `json:"project"`
	StatusDetails       map[string]interface{} `json:"statusDetails"`
	Stats               Stats                  `json:"stats"`
	Annotations         []string               `json:"annotations"`
	AssignedTo          interface{}            `json:"assignedTo"`
}

// Metadata contains additional information about the issue
type Metadata struct {
	Type     string `json:"type"`
	Message  string `json:"message,omitempty"`
	Title    string `json:"title,omitempty"`
	Value    string `json:"value,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// SubscriptionDetails contains subscription information for the issue
type SubscriptionDetails struct {
	Reason string `json:"reason"`
}

// Project contains information about the Sentry project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Stats contains statistical information about the issue
type Stats struct {
	TwentyFourHours [][2]float64 `json:"24h"`
}

// SentryEvent represents a single event from Sentry
type SentryEvent struct {
	EventID   string                   `json:"eventID"`
	ID        string                   `json:"id"`
	GroupID   string                   `json:"groupID"`
	Title     string                   `json:"title"`
	Message   string                   `json:"message"`
	Timestamp string                   `json:"dateCreated"`
	Tags      [][]string               `json:"tags"`
	Platform  string                   `json:"platform"`
	User      map[string]interface{}   `json:"user,omitempty"`
	Contexts  map[string]interface{}   `json:"contexts,omitempty"`
	Entries   []map[string]interface{} `json:"entries,omitempty"`
	Metadata  Metadata                 `json:"metadata"`
}

// SentryTagValue represents a single tag value with its metadata
type SentryTagValue struct {
	Value     string `json:"value"`
	Count     int    `json:"count"`
	LastSeen  string `json:"lastSeen"`
	FirstSeen string `json:"firstSeen"`
}

// SentryTagDetails represents the details of a specific tag for an issue
type SentryTagDetails struct {
	Key          string           `json:"key"`
	Name         string           `json:"name"`
	UniqueValues int              `json:"uniqueValues"`
	TotalValues  int              `json:"totalValues"`
	TopValues    []SentryTagValue `json:"topValues"`
}

// SentryTagValueSorted represents a tag value with its occurrence count, sorted by count
type SentryTagValueSorted struct {
	Value string
	Count int
}

// SentryAffectedPath represents a URL path affected by errors in Sentry
type SentryAffectedPath struct {
	Path  string
	Count int
}

// SentryIssuesResponse represents the response from Sentry's API
type SentryIssuesResponse struct {
	Issues []SentryIssue `json:"issues"`
}

// SentryIssueHash represents a hash for a specific issue
type SentryIssueHash struct {
	ID          string `json:"id"`
	LatestEvent string `json:"latestEvent"`
}

// SentryIssueWithHashes extends SentryIssue with its hashes
type SentryIssueWithHashes struct {
	SentryIssue
	Hashes []SentryIssueHash `json:"hashes"`
}

// GetSentryIssues retrieves issues from a Sentry project using the Sentry API
// It requires an auth token and organization/project slugs
func GetSentryIssues(authToken, orgSlug, projectSlug string) ([]SentryIssue, error) {
	logger.Debug("Getting Sentry issues for project: %s/%s", orgSlug, projectSlug)

	// Construct the API URL
	apiURL := fmt.Sprintf("https://sentry.io/api/0/projects/%s/%s/issues/", orgSlug, projectSlug)

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
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
	var issues []SentryIssue
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	logger.Debug("Retrieved %d Sentry issues", len(issues))
	return issues, nil
}

// GetSentryIssueTagDetails retrieves details for a specific tag of an issue
func GetSentryIssueTagDetails(authToken, orgSlug, issueID, tagKey string) (*SentryTagDetails, error) {
	logger.Debug("Getting tag details for Sentry issue: %s, tag: %s in organization: %s", issueID, tagKey, orgSlug)

	// Construct the API URL
	apiURL := fmt.Sprintf("https://sentry.io/api/0/organizations/%s/issues/%s/tags/%s/", orgSlug, issueID, tagKey)
	logger.Debug("Requesting URL: %s", apiURL)

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
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
	var tagDetails SentryTagDetails
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
func GetSentryIssueTagValuesSorted(authToken, orgSlug, issueID, tagKey string) ([]SentryTagValueSorted, error) {
	logger.Debug("Getting sorted tag values for Sentry issue: %s, tag: %s", issueID, tagKey)

	// First get the tag details
	tagDetails, err := GetSentryIssueTagDetails(authToken, orgSlug, issueID, tagKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag details: %w", err)
	}

	// Create a slice to hold the sorted values
	sortedValues := make([]SentryTagValueSorted, 0, len(tagDetails.TopValues))

	// Convert the tag values to our sorted format
	for _, value := range tagDetails.TopValues {
		sortedValues = append(sortedValues, SentryTagValueSorted{
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
func GetAffectedSentryPaths(authToken, orgSlug, projectSlug string) ([]SentryAffectedPath, error) {
	logger.Debug("Getting affected URL paths from Sentry for organization: %s, project: %s", orgSlug, projectSlug)
	
	// First get the issues
	issues, err := GetSentryIssues(authToken, orgSlug, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}
	
	logger.Debug("Retrieved %d issues, now fetching URL tag details for each", len(issues))
	if len(issues) == 0 {
		logger.Debug("No issues found to fetch URL paths for")
		return []SentryAffectedPath{}, nil
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
		tagDetails, err := GetSentryIssueTagDetails(authToken, orgSlug, issue.ID, "url")
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
	affectedPaths := make([]SentryAffectedPath, 0, len(pathCountMap))
	for path, count := range pathCountMap {
		affectedPaths = append(affectedPaths, SentryAffectedPath{
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
