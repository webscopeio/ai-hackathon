package analyze

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/webscopeio/ai-hackathon/internal/logger"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

// GetSitemap attempts to retrieve and parse a sitemap from a given URL
// It tries common sitemap locations if not explicitly provided
func GetSitemap(ctx context.Context, baseURL string) (*models.Sitemap, error) {
	logger.Debug("Getting sitemap for URL: %s", baseURL)
	
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	// Ensure the URL has a scheme
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}
	
	// Create a list of potential sitemap URLs to check
	sitemapURLs := []string{
		fmt.Sprintf("%s://%s/sitemap.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap_index.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap-index.xml", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap.php", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/sitemap", parsedURL.Scheme, parsedURL.Host),
	}
	
	// Try to find robots.txt first, which might contain sitemap URL
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	robotsSitemapURL := getSitemapFromRobots(ctx, robotsURL)
	if robotsSitemapURL != "" {
		sitemapURLs = append([]string{robotsSitemapURL}, sitemapURLs...)
	}
	
	// Try each potential sitemap URL
	for _, sitemapURL := range sitemapURLs {
		logger.Debug("Trying sitemap URL: %s", sitemapURL)
		
		// Try to get the sitemap
		sitemap, err := fetchSitemap(ctx, sitemapURL)
		if err == nil && sitemap != nil && len(sitemap.URLs) > 0 {
			logger.Debug("Found sitemap at %s with %d URLs", sitemapURL, len(sitemap.URLs))
			return sitemap, nil
		}
		
		// If we get a sitemap index instead, try to get the first sitemap from it
		sitemapIndex, err := fetchSitemapIndex(ctx, sitemapURL)
		if err == nil && sitemapIndex != nil && len(sitemapIndex.Sitemaps) > 0 {
			logger.Debug("Found sitemap index at %s with %d sitemaps", sitemapURL, len(sitemapIndex.Sitemaps))
			
			// Get the first sitemap from the index
			firstSitemapURL := sitemapIndex.Sitemaps[0].Loc
			sitemap, err := fetchSitemap(ctx, firstSitemapURL)
			if err == nil && sitemap != nil {
				logger.Debug("Found sitemap at %s with %d URLs", firstSitemapURL, len(sitemap.URLs))
				return sitemap, nil
			}
		}
	}
	
	return nil, fmt.Errorf("no sitemap found for %s", baseURL)
}

// fetchSitemap attempts to retrieve and parse a sitemap from a given URL
func fetchSitemap(ctx context.Context, sitemapURL string) (*models.Sitemap, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", sitemapURL, nil)
	if err != nil {
		return nil, err
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var sitemap models.Sitemap
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return nil, err
	}
	
	return &sitemap, nil
}

// fetchSitemapIndex attempts to retrieve and parse a sitemap index from a given URL
func fetchSitemapIndex(ctx context.Context, sitemapURL string) (*models.SitemapIndex, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", sitemapURL, nil)
	if err != nil {
		return nil, err
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var sitemapIndex models.SitemapIndex
	err = xml.Unmarshal(body, &sitemapIndex)
	if err != nil {
		return nil, err
	}
	
	return &sitemapIndex, nil
}

// getSitemapFromRobots tries to extract sitemap URL from robots.txt
func getSitemapFromRobots(ctx context.Context, robotsURL string) string {
	req, err := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	if err != nil {
		return ""
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	
	// Look for Sitemap: directive in robots.txt
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	
	return ""
}
