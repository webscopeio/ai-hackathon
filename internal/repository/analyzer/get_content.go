package analyzer

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func GetContent_OLD(ctx context.Context, urls []string) (*models.GetContentToolReturn, error) {
	if len(urls) == 0 {
		return nil, errors.New("empty URLs list provided")
	}

	validatedUrls := make([]string, 0, len(urls))
	for _, urlStr := range urls {
		if urlStr == "" {
			continue
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			continue
		}

		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			if strings.Contains(parsedURL.String(), "localhost") || strings.Contains(parsedURL.String(), ":") {
				urlStr = "http://" + urlStr
			} else {
				urlStr = "https://" + urlStr
			}
			parsedURL, err = url.Parse(urlStr)
			if err != nil {
				continue
			}
		}
		validatedUrls = append(validatedUrls, parsedURL.String())
	}

	if len(validatedUrls) == 0 {
		return nil, errors.New("no valid URLs provided")
	}

	results := make(map[string]string)
	var mutex sync.Mutex

	classRegex := regexp.MustCompile(`\s+class\s*=\s*"[^"]*"|\s+class\s*=\s*'[^']*'`)
	styleRegex := regexp.MustCompile(`\s+style\s*=\s*"[^"]*"|\s+style\s*=\s*'[^']*'`)
	scriptRegex := regexp.MustCompile(`<script\b[^>]*>[\s\S]*?</script>`)
	svgRegex := regexp.MustCompile(`<svg\b[^>]*>[\s\S]*?</svg>`)

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 24,
		Delay:       10 * time.Millisecond,
	})

	c.SetRequestTimeout(10 * time.Second)

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

		cleanedHtml := classRegex.ReplaceAllString(html, "")
		cleanedHtml = styleRegex.ReplaceAllString(cleanedHtml, "")
		cleanedHtml = scriptRegex.ReplaceAllString(cleanedHtml, "")
		cleanedHtml = svgRegex.ReplaceAllString(cleanedHtml, "")

		mutex.Lock()
		results[url] = cleanedHtml
		mutex.Unlock()
	})

	c.OnError(func(r *colly.Response, err error) {
		url := r.Request.URL.String()
		mutex.Lock()
		results[url] = ""
		mutex.Unlock()
	})

	for _, urlStr := range validatedUrls {
		if err := c.Visit(urlStr); err != nil {
			mutex.Lock()
			results[urlStr] = ""
			mutex.Unlock()
		}
	}

	c.Wait()

	if len(results) == 0 {
		return nil, errors.New("failed to fetch content from any URL")
	}

	fmt.Println("results are ready", len(results))

	return &models.GetContentToolReturn{
		Contents: results,
	}, nil
}

// GetContentCdp uses Chrome DevTools Protocol via chromedp to fetch content from URLs
// This is especially useful for SPAs and JavaScript-heavy websites
func GetContent(ctx context.Context, urls []string) (*models.GetContentToolReturn, error) {
	if len(urls) == 0 {
		return nil, errors.New("empty URLs list provided")
	}

	// Validate and normalize URLs
	validatedUrls := make([]string, 0, len(urls))
	for _, urlStr := range urls {
		if urlStr == "" {
			continue
		}

		// Ensure URL has a scheme
		if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
			if strings.Contains(urlStr, "localhost") || strings.Contains(urlStr, ":") {
				urlStr = "http://" + urlStr
			} else {
				urlStr = "https://" + urlStr
			}
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("Warning: Could not parse URL %s: %v\n", urlStr, err)
			continue
		}

		validatedUrls = append(validatedUrls, parsedURL.String())
	}

	if len(validatedUrls) == 0 {
		return nil, errors.New("no valid URLs provided")
	}

	// Map to store results
	results := make(map[string]string)
	var mutex sync.Mutex

	// Regular expressions for cleaning HTML
	cleaningRegexes := []*regexp.Regexp{
		regexp.MustCompile(`\s+class\s*=\s*"[^"]*"|\s+class\s*=\s*'[^']*'`),
		regexp.MustCompile(`\s+style\s*=\s*"[^"]*"|\s+style\s*=\s*'[^']*'`),
		regexp.MustCompile(`<script\b[^>]*>[\s\S]*?</script>`),
		regexp.MustCompile(`<svg\b[^>]*>[\s\S]*?</svg>`),
	}

	// Create a new browser context with chromedp
	browserCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Process URLs sequentially to avoid overwhelming the browser
	for _, urlStr := range validatedUrls {
		// Create a timeout context for each URL
		urlCtx, urlCancel := context.WithTimeout(browserCtx, 30*time.Second)
		defer urlCancel()
		
		// Variable to store the body content
		var bodyHTML string
		
		// Navigate to the URL and capture the body content
		err := chromedp.Run(urlCtx,
			// Navigate to the URL
			chromedp.Navigate(urlStr),
			
			// Wait for the page to be fully loaded
			chromedp.WaitReady("body", chromedp.ByQuery),
			
			// Optional: Wait some extra time for dynamic content
			chromedp.Sleep(1*time.Second),
			
			// Capture the HTML content
			chromedp.OuterHTML("body", &bodyHTML, chromedp.ByQuery),
		)		
		
		if err != nil {
			fmt.Printf("Error fetching %s with chromedp: %v\n", urlStr, err)
			
			// Store empty string for failed URLs
			mutex.Lock()
			results[urlStr] = ""
			mutex.Unlock()
			continue
		}
		
		// Clean the HTML content
		for _, regex := range cleaningRegexes {
			bodyHTML = regex.ReplaceAllString(bodyHTML, "")
		}
		
		// Store the result
		mutex.Lock()
		results[urlStr] = bodyHTML
		mutex.Unlock()
		
		fmt.Printf("Successfully fetched %s with chromedp\n", urlStr)
	}

	// Check if we got any results
	if len(results) == 0 {
		return nil, errors.New("failed to fetch content from any URL")
	}

	fmt.Printf("Results are ready, found data for %d URLs\n", len(results))

	return &models.GetContentToolReturn{
		Contents: results,
	}, nil
}
