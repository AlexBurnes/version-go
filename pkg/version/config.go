// Package version provides semantic version parsing, validation, and ordering functionality.
// It supports extended version formats beyond standard SemVer 2.0, including prerelease,
// postrelease, and intermediate identifiers.
package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the structure of .project.yml file
type ProjectConfig struct {
	Project struct {
		Name    string   `yaml:"name"`
		Modules []string `yaml:"modules"`
	} `yaml:"project"`
}

// ConfigProvider provides project configuration information
type ConfigProvider struct {
	config *ProjectConfig
}

// NewConfigProvider creates a new configuration provider
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{}
}

// LoadProjectConfig loads .project.yml configuration from the project root
// Returns the configuration if found and valid, nil if not found, error if invalid
func (cp *ConfigProvider) LoadProjectConfig() (*ProjectConfig, error) {
	// Look for .project.yml in current directory and parent directories
	configPath, err := findProjectConfigFile()
	if err != nil {
		return nil, err
	}

	if configPath == "" {
		return nil, nil // No .project.yml found
	}

	// Read and parse the YAML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .project.yml: %v", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse .project.yml: %v", err)
	}

	// Validate the configuration
	if err := cp.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid .project.yml: %v", err)
	}

	cp.config = &config
	return &config, nil
}

// GetProjectName returns the project name from configuration or empty string if not available
func (cp *ConfigProvider) GetProjectName() string {
	if cp.config == nil {
		return ""
	}
	return cp.config.Project.Name
}

// GetModuleName returns the primary module name from configuration or empty string if not available
func (cp *ConfigProvider) GetModuleName() string {
	if cp.config == nil || len(cp.config.Project.Modules) == 0 {
		return ""
	}
	return cp.config.Project.Modules[0] // First module is primary
}

// GetAllModules returns all module names from configuration or empty slice if not available
func (cp *ConfigProvider) GetAllModules() []string {
	if cp.config == nil {
		return []string{}
	}
	return cp.config.Project.Modules
}

// HasConfig returns true if a valid configuration is loaded
func (cp *ConfigProvider) HasConfig() bool {
	return cp.config != nil
}

// findProjectConfigFile searches for .project.yml in current directory and parent directories
func findProjectConfigFile() (string, error) {
	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Search up the directory tree for .project.yml
	for {
		configPath := filepath.Join(dir, ".project.yml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Check if we've reached the root directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root directory
		}
		dir = parent
	}

	return "", nil // No .project.yml found
}

// validateConfig validates the loaded configuration
func (cp *ConfigProvider) validateConfig(config *ProjectConfig) error {
	if config.Project.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if len(config.Project.Modules) == 0 {
		return fmt.Errorf("at least one module is required")
	}

	// Validate module names (no empty strings)
	for i, module := range config.Project.Modules {
		if strings.TrimSpace(module) == "" {
			return fmt.Errorf("module %d cannot be empty", i+1)
		}
	}

	return nil
}

// GetProjectConfigFromFile loads configuration from a specific file path
// This is useful for testing or when you know the exact path
func GetProjectConfigFromFile(filePath string) (*ProjectConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate the configuration
	cp := &ConfigProvider{}
	if err := cp.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config file: %v", err)
	}

	return &config, nil
}