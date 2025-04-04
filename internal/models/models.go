package models

import (
	"fmt"
	"strings"
)

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
	Url string `json:"url"`
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
