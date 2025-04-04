package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/webscopeio/ai-hackathon/internal/models"
	"gopkg.in/yaml.v3"
)

const (
	configFileName = "user_config.yaml"
)

// Manager handles the operations for user configuration
type Manager struct {
	configPath string
	mutex      sync.RWMutex
}

// NewManager creates a new config manager
func NewManager() (*Manager, error) {
	// Get the root directory of the project
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	rootDir := workDir
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			rootDir = workDir
			break
		}
		rootDir = parent
	}

	// Store config file in the internal/config directory (same as config.go)
	configDir := filepath.Join(rootDir, "internal", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, configFileName)
	return &Manager{
		configPath: configPath,
	}, nil
}

// GetConfig retrieves the current user configuration
func (m *Manager) GetConfig() (*models.UserConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if config file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return &models.UserConfig{}, nil
	}

	// Read config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config models.UserConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the user configuration to a YAML file
func (m *Manager) SaveConfig(config *models.UserConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdateConfig updates only the provided fields in the configuration
func (m *Manager) UpdateConfig(updates *models.UserConfig) error {
	// Get current config
	current, err := m.GetConfig()
	if err != nil {
		return err
	}

	// Update only non-empty fields
	if updates.AnthropicApiKey != "" {
		current.AnthropicApiKey = updates.AnthropicApiKey
	}
	if updates.SentryApiKey != "" {
		current.SentryApiKey = updates.SentryApiKey
	}
	if updates.UmamiAPIKey != "" {
		current.UmamiAPIKey = updates.UmamiAPIKey
	}
	if updates.UmamiWebsiteId != "" {
		current.UmamiWebsiteId = updates.UmamiWebsiteId
	}
	if updates.TechSpecification != "" {
		current.TechSpecification = updates.TechSpecification
	}
	// Always update the product specification, even if empty
	current.ProductSpecification = updates.ProductSpecification

	// Save updated config
	return m.SaveConfig(current)
}
