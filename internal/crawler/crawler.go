package crawler

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func Crawl(ctx context.Context, urlStr string, maxDepth int, maxPathSegments int) ([]string, map[string]string, error) {
	if urlStr == "" {
		return nil, nil, errors.New("empty URL provided")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		if strings.Contains(parsedURL.String(), "localhost") || strings.Contains(parsedURL.String(), ":") {
			urlStr = "http://" + urlStr
		} else {
			urlStr = "https://" + urlStr
		}
		parsedURL, err = url.Parse(urlStr)
		if err != nil {
			return nil, nil, err
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

		mutex.Lock()
		results[url] = html
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
		return nil, nil, err
	}

	c.Wait()

	return links, results, nil
}
