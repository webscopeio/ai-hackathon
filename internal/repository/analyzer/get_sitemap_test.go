package analyze

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSitemap(t *testing.T) {
	// Create a test server that serves a simple sitemap
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sitemap.xml" {
			// Return a simple sitemap
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2025-04-01</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://example.com/about</loc>
    <lastmod>2025-03-20</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://example.com/contact</loc>
    <lastmod>2025-03-15</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.5</priority>
  </url>
</urlset>`))
		} else if r.URL.Path == "/robots.txt" {
			// Return a robots.txt file with a sitemap directive
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`User-agent: *
Allow: /
Sitemap: http://example.com/sitemap.xml`))
		} else {
			// Return 404 for all other paths
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Test getting the sitemap
	sitemap, err := GetSitemap(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("Failed to get sitemap: %v", err)
	}

	// Check that we got the expected number of URLs
	if len(sitemap.URLs) != 3 {
		t.Errorf("Expected 3 URLs in sitemap, got %d", len(sitemap.URLs))
	}

	// Check the first URL
	if sitemap.URLs[0].Loc != "https://example.com/" {
		t.Errorf("Expected first URL to be https://example.com/, got %s", sitemap.URLs[0].Loc)
	}
	if sitemap.URLs[0].LastMod != "2025-04-01" {
		t.Errorf("Expected first URL lastmod to be 2025-04-01, got %s", sitemap.URLs[0].LastMod)
	}
	if sitemap.URLs[0].ChangeFreq != "daily" {
		t.Errorf("Expected first URL changefreq to be daily, got %s", sitemap.URLs[0].ChangeFreq)
	}
	if sitemap.URLs[0].Priority != 1.0 {
		t.Errorf("Expected first URL priority to be 1.0, got %f", sitemap.URLs[0].Priority)
	}
}

func TestGetSitemapIndex(t *testing.T) {
	// Create a handler that will be updated with the server URL after creation
	mux := http.NewServeMux()
	
	// Create a test server with the handler
	server := httptest.NewServer(mux)
	defer server.Close()
	
	// Now set up the handler with the server URL
	mux.HandleFunc("/sitemap_index.xml", func(w http.ResponseWriter, r *http.Request) {
		// Return a sitemap index
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		
		// Create the sitemap index XML with the server URL
		sitemapIndex := strings.ReplaceAll(`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <sitemap>
    <loc>SERVER_URL/sitemap1.xml</loc>
    <lastmod>2025-04-01</lastmod>
  </sitemap>
  <sitemap>
    <loc>SERVER_URL/sitemap2.xml</loc>
    <lastmod>2025-03-20</lastmod>
  </sitemap>
</sitemapindex>`, "SERVER_URL", server.URL)
		
		w.Write([]byte(sitemapIndex))
	})
	
	mux.HandleFunc("/sitemap1.xml", func(w http.ResponseWriter, r *http.Request) {
		// Return the first sitemap
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/page1</loc>
    <lastmod>2025-04-01</lastmod>
    <priority>0.9</priority>
  </url>
  <url>
    <loc>https://example.com/page2</loc>
    <lastmod>2025-03-28</lastmod>
    <priority>0.8</priority>
  </url>
</urlset>`))
	})
	
	// Handle 404 for all other paths
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Test getting the sitemap from the index
	sitemap, err := GetSitemap(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("Failed to get sitemap from index: %v", err)
	}

	// Check that we got the expected number of URLs
	if len(sitemap.URLs) != 2 {
		t.Errorf("Expected 2 URLs in sitemap, got %d", len(sitemap.URLs))
	}

	// Check the first URL
	if sitemap.URLs[0].Loc != "https://example.com/page1" {
		t.Errorf("Expected first URL to be https://example.com/page1, got %s", sitemap.URLs[0].Loc)
	}
	if sitemap.URLs[0].Priority != 0.9 {
		t.Errorf("Expected first URL priority to be 0.9, got %f", sitemap.URLs[0].Priority)
	}
}
