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

	"github.com/gocolly/colly/v2"
	"github.com/webscopeio/ai-hackathon/internal/models"
)

func GetContent(ctx context.Context, urls []string) (*models.GetContentToolReturn, error) {
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
