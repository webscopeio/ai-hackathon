package models

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Sitemap represents the XML structure of a sitemap
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents a URL entry in a sitemap
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// SitemapIndex represents the XML structure of a sitemap index
type SitemapIndex struct {
	XMLName  xml.Name        `xml:"sitemapindex"`
	Sitemaps []SitemapEntry `xml:"sitemap"`
}

// SitemapEntry represents a sitemap entry in a sitemap index
type SitemapEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

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

type StatusReturn struct {
	Status string `json:"status"`
}

type CrawlArgs struct {
	Url             string `json:"url"`
	MaxDepth        int    `json:"maxDepth,omitempty"`
	MaxPathSegments int    `json:"maxPathSegments,omitempty"`
}

type CrawlReturn struct {
	Links   []string          `json:"links"`
	Results map[string]string `json:"results"`
}

type GenerateTestsArgs struct {
	Url              string `json:"url"`
	TechSpecification string `json:"techSpecification,omitempty"`
}

type GenerateTestReturn struct {
	FileName     string   `json:"filename" jsonschema_description:"Name of the test file (e.g., 'login.spec.ts')"`
	Content      string   `json:"content" jsonschema_description:"Complete content of the test file"`
	Dependencies []string `json:"dependencies" jsonschema_description:"NPM packages required for the test file"`
}

type TestFile struct {
	Filename string `json:"filename" jsonschema_description:"Name of the test file (e.g., 'login.spec.ts')"`
	Content  string `json:"content" jsonschema_description:"Complete content of the test file"`
	FilePath string `json:"-" jsonschema_description:"Absolute path to the temporary file where the test content is stored"`
}

type Failure struct {
	Filename string `json:"filename" jsonschema_description:"Name of the test file where the failure occurred"`
	Error    string `json:"error" jsonschema_description:"Error message or reason for failure"`
}

type SitemapTool struct {
	BaseUrl string `json:"baseUrl" jsonschema_description:"The base URL needed to get a website's sitemap"`
}

type GetContentTool struct {
	Urls []string `json:"urls" jsonschema_description:"Array of the URLs which content should be retrieved"`
}

type GetContentToolReturn struct {
	Contents map[string]string `json:"contents"`
}

type SentryTool struct {
	OrgSlug     string `json:"orgSlug" jsonschema_description:"The Sentry organization slug"`
	ProjectSlug string `json:"projectSlug" jsonschema_description:"The Sentry project slug"`
}

type FinalCriteriaTool struct {
	Criteria     string `json:"criteria" jsonschema_description:"The criteria to be used for the generation of the E2E tests"`
	TechSpec string            `json:"techSpec" jsonschema_description:"The technical specification of the website"`
	ContentMap  map[string]string `json:"contentMap" jsonschema_description:"Map of URLs to their HTML content"`
}

type AnalyzerToolExample struct {
	Greeting string `json:"greeting" jsonschema_description:"This is just a friendly greeting"`
}

type AnalyzerReturn struct {
	TechSpec string            `json:"techSpec"`
	SiteMap  map[string]string `json:"siteMap"`
	Criteria string            `json:"criteria"`
}

type EvaluationReturn struct {
	Passed   bool   `json:"passed" jsonschema_description:"Whether the test file is good enough or needs more work"`
	Feedback string `json:"feedback" jsonschema_description:"Feedback on the test file"`
}

func (r *GenerateTestReturn) Validate() error {
	var missingFields []string

	if r.FileName == "" {
		missingFields = append(missingFields, "filename")
	}
	if r.Content == "" {
		missingFields = append(missingFields, "content")
	}
	if len(r.Dependencies) == 0 {
		missingFields = append(missingFields, "dependencies")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

type AnalysisReturn struct {
	Analysis string            `json:"analysis" jsonschema_description:"Analysis of the website and potential test cases"`
	Links    []string          `json:"links" jsonschema_description:"List of URLs that were crawled"`
	Results  map[string]string `json:"results" jsonschema_description:"Map of URLs to their HTML content"`
}

func (r *AnalysisReturn) Validate() error {
	if r.Analysis == "" {
		return fmt.Errorf("required field: analysis")
	}
	return nil
}

type ErrorReturn struct {
	Error string `json:"error"`
}
