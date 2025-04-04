package models

// UserConfig represents the user configuration for the test generation
type UserConfig struct {
	AnthropicApiKey       string `json:"anthropicApiKey" yaml:"anthropicApiKey"`
	SentryApiKey          string `json:"sentryApiKey" yaml:"sentryApiKey"`
	UmamiAPIKey           string `json:"umamiAPIKey" yaml:"umamiAPIKey"`
	UmamiWebsiteId        string `json:"umamiWebsiteId" yaml:"umamiWebsiteId"`
	TechSpecification     string `json:"techSpecification" yaml:"techSpecification"`
	ProductSpecification  string `json:"productSpecification" yaml:"productSpecification"`
}
