package analyze

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/webscopeio/ai-hackathon/internal/llm"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func Analyze(ctx context.Context, urlStr string, maxDepth int, maxPathSegments int, prompt string, client *llm.Client) (*models.AnalysisReturn, error) {
	if urlStr == "" {
		return nil, errors.New("empty URL provided")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		if strings.Contains(parsedURL.String(), "localhost") || strings.Contains(parsedURL.String(), ":") {
			urlStr = "http://" + urlStr
		} else {
			urlStr = "https://" + urlStr
		}
		parsedURL, err = url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
	}

	basePathSegments := 0
	trimmedBasePath := strings.Trim(parsedURL.Path, "/")
	if trimmedBasePath != "" {
		basePathSegments = strings.Count(trimmedBasePath, "/") + 1
	}

	results := make(map[string]string)
	var links []string
	var mutex sync.Mutex

	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 24,
		Delay:       10 * time.Millisecond,
	})

	c.SetRequestTimeout(10 * time.Second)
	c.AllowedDomains = []string{parsedURL.Host}

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.OnRequest(func(r *colly.Request) {
				r.Abort()
			})
		case <-done:
			return
		}
	}()
	defer close(done)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		html, err := e.DOM.Html()
		if err != nil {
			return
		}

		// Regex to remove CSS classes while preserving tags and IDs
		// This matches class="anything" or class='anything' patterns
		classRegex := regexp.MustCompile(`\s+class\s*=\s*"[^"]*"|\s+class\s*=\s*'[^']*'`)
		cleanedHTML := classRegex.ReplaceAllString(html, "")

		mutex.Lock()
		results[url] = cleanedHTML
		links = append(links, url)
		mutex.Unlock()
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		parsedLink, err := url.Parse(link)
		if err != nil {
			return
		}

		if parsedLink.Fragment != "" {
			return
		}

		if !parsedLink.IsAbs() {
			parsedLink = e.Request.URL.ResolveReference(parsedLink)
		}

		if parsedLink.Host != parsedURL.Host &&
			!strings.HasPrefix(parsedLink.Path, parsedURL.Path) {
			return
		}

		linkPathSegments := 0
		trimmedPath := strings.Trim(parsedLink.Path, "/")
		if trimmedPath != "" {
			linkPathSegments = strings.Count(trimmedPath, "/") + 1
		}

		relativePathSegments := linkPathSegments - basePathSegments

		if maxPathSegments != 0 && relativePathSegments > maxPathSegments {
			return
		}

		e.Request.Visit(e.Attr("href"))
	})

	if err := c.Visit(urlStr); err != nil {
		return nil, err
	}

	c.Wait()

	// If no client is provided, skip the analysis
	if client == nil {
		return &models.AnalysisReturn{
			Links:   links,
			Results: results,
		}, nil
	}

	// Build context from crawled results
	var builder strings.Builder
	builder.WriteString("WEBSITE CONTENT: \n")
	for url, html := range results {
		builder.WriteString("URL: ")
		builder.WriteString(url)
		builder.WriteString("\nCONTENT:\n")
		builder.WriteString(html)
		builder.WriteString("\n---END PAGE---\n\n")
	}
	crawlContext := builder.String()

	// Create analysis prompt
	analysisPrompt := fmt.Sprintf(`Analyze the following website content and description to identify potential test cases for E2E testing:

WEBSITE CONTENT:
%s

WEBSITE DESCRIPTION:
%s

Based on the website content and description, provide a comprehensive analysis that includes:
1. Key user flows that should be tested
2. Important UI components that need validation
3. Potential edge cases and error states
4. Recommendations for comprehensive test coverage

Your response should be a detailed text analysis with clear, actionable test case suggestions.`, crawlContext, prompt)

	// Get text completion from LLM
	analysisText, err := client.GetCompletion(ctx, analysisPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	// Create analysis return object
	analysis := &models.AnalysisReturn{
		Analysis: analysisText,
		Links:    links,
		Results:  results,
	}

	// Validate the response
	if err := analysis.Validate(); err != nil {
		return nil, fmt.Errorf("invalid analysis: %w", err)
	}

	return analysis, nil
}
