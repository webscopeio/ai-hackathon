package models

// UserConfig represents the user configuration for the test generation
type UserConfig struct {
	AnthropicApiKey   string `json:"anthropicApiKey" yaml:"anthropicApiKey"`
	SentryApiKey      string `json:"sentryApiKey" yaml:"sentryApiKey"`
	TechSpecification string `json:"techSpecification" yaml:"techSpecification"`
}
