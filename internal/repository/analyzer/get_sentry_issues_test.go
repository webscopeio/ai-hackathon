package analyze

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/webscopeio/ai-hackathon/internal/config"
)

func TestGetSentryIssues(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load Sentry auth token from environment variables
	cfg := config.Load()

	// Skip test if Sentry auth token is not configured
	if cfg.SentryAuthToken == "" {
		t.Skip("Skipping test as Sentry auth token is not configured in .env file")
	}

	// Replace with your actual Sentry credentials
	orgSlug := "webscopeio-pb"
	projectSlug := "ai-hackathon-demo"

	// Call the function with context and config
	issues, err := GetSentryIssues(context.Background(), cfg, orgSlug, projectSlug)
	if err != nil {
		t.Logf("Note: This test requires valid Sentry credentials to pass")
		t.Logf("Error getting Sentry issues: %v", err)
		return
	}

	// Print the results
	t.Logf("Successfully retrieved %d Sentry issues", len(issues))

	// Print details of up to 3 issues
	for i, issue := range issues {
		if i >= 3 {
			break
		}

		issueJSON, _ := json.MarshalIndent(issue, "", "  ")
		t.Logf("Issue %d:\n%s", i+1, string(issueJSON))
	}
}

func TestGetSentryIssuesWithTagDetails(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load Sentry auth token from environment variables
	cfg := config.Load()

	// Skip test if Sentry auth token is not configured
	if cfg.SentryAuthToken == "" {
		t.Skip("Skipping test as Sentry auth token is not configured in .env file")
	}

	// Replace with your actual Sentry credentials
	orgSlug := "webscopeio-pb"
	projectSlug := "ai-hackathon-demo"

	// Tag key to retrieve details for
	tagKey := "url"

	// Call the function to get issues
	issues, err := GetSentryIssues(context.Background(), cfg, orgSlug, projectSlug)
	if err != nil {
		t.Logf("Note: This test requires valid Sentry credentials to pass")
		t.Logf("Error getting Sentry issues: %v", err)
		return
	}

	// Print the results
	t.Logf("Successfully retrieved %d Sentry issues", len(issues))

	// Limit to first 20 issues or less
	maxIssues := 20
	if len(issues) < maxIssues {
		maxIssues = len(issues)
	}

	// For each issue, get tag details
	for i, issue := range issues[:maxIssues] {
		t.Logf("\nIssue %d/%d: %s - %s", i+1, maxIssues, issue.ShortID, issue.Title)

		// Get details for the tag
		tagDetails, err := GetSentryIssueTagDetails(context.Background(), cfg, orgSlug, issue.ID, tagKey)
		if err != nil {
			t.Logf("  Error getting tag details for %s: %v", tagKey, err)
			continue
		}

		t.Logf("  Tag: %s (%s), Unique Values: %d, Total Values: %d",
			tagDetails.Key, tagDetails.Name, tagDetails.UniqueValues, tagDetails.TotalValues)

		// Print up to 2 top values for this tag
		for j, value := range tagDetails.TopValues {
			if j >= 2 {
				break
			}
			t.Logf("    Value: %s, Count: %d, First Seen: %s",
				value.Value, value.Count, value.FirstSeen)
		}
	}
}

func TestGetSentryIssueTagValuesSorted(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load Sentry auth token from environment variables
	cfg := config.Load()

	// Skip test if Sentry auth token is not configured
	if cfg.SentryAuthToken == "" {
		t.Skip("Skipping test as Sentry auth token is not configured in .env file")
	}

	// Use hardcoded value for organization
	orgSlug := "webscopeio-pb"

	// This should be a valid issue ID from your Sentry project
	issueID := "35963534" // Same as in other tests

	// Tag key to test
	tagKey := "url"

	// Call the function with context and config
	sortedValues, err := GetSentryIssueTagValuesSorted(context.Background(), cfg, orgSlug, issueID, tagKey)
	if err != nil {
		t.Logf("Note: This test requires valid Sentry credentials, issue ID, and tag key to pass")
		t.Logf("Error getting sorted tag values: %v", err)
		return
	}

	// Print the results
	t.Logf("Successfully retrieved %d sorted tag values for issue %s, tag: %s",
		len(sortedValues), issueID, tagKey)

	// Print all values in sorted order
	t.Logf("Values sorted by occurrence count (highest first):")
	for i, value := range sortedValues {
		t.Logf("%d. %s: %d occurrences", i+1, value.Value, value.Count)
	}

	// Verify the sorting is correct
	if len(sortedValues) > 1 {
		for i := 0; i < len(sortedValues)-1; i++ {
			if sortedValues[i].Count < sortedValues[i+1].Count {
				t.Errorf("Values are not properly sorted: %d occurrences followed by %d occurrences",
					sortedValues[i].Count, sortedValues[i+1].Count)
			}
		}
	}
}

func TestGetAffectedSentryPaths(t *testing.T) {
	// Skip this test in short mode as it requires real API access
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires real API access")
	}

	// Load Sentry auth token from environment variables
	cfg := config.Load()

	// Skip test if Sentry auth token is not configured
	if cfg.SentryAuthToken == "" {
		t.Skip("Skipping test as Sentry auth token is not configured in .env file")
	}

	// Use hardcoded values for organization and project
	orgSlug := "webscopeio-pb"
	projectSlug := "ai-hackathon-demo"

	// Call the function with context and config
	affectedPaths, err := GetAffectedSentryPaths(context.Background(), cfg, orgSlug, projectSlug)
	if err != nil {
		t.Logf("Note: This test requires valid Sentry credentials to pass")
		t.Logf("Error getting affected Sentry paths: %v", err)
		return
	}

	// Print the results
	t.Logf("Successfully retrieved %d affected URL paths", len(affectedPaths))

	// Print the top 10 paths or all if less than 10
	maxPaths := 10
	if len(affectedPaths) < maxPaths {
		maxPaths = len(affectedPaths)
	}

	t.Logf("Top %d affected URL paths (sorted by occurrence count):", maxPaths)
	for i, path := range affectedPaths[:maxPaths] {
		t.Logf("%d. %s: %d occurrences", i+1, path.Path, path.Count)
	}

	// Verify the sorting is correct
	if len(affectedPaths) > 1 {
		for i := 0; i < len(affectedPaths)-1; i++ {
			if affectedPaths[i].Count < affectedPaths[i+1].Count {
				t.Errorf("Paths are not properly sorted: %d occurrences followed by %d occurrences",
					affectedPaths[i].Count, affectedPaths[i+1].Count)
			}
		}
	}
}
