package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectConfig(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected *ProjectConfig
		wantErr  bool
	}{
		{
			name: "valid config",
			yaml: `project:
  name: "test-project"
  modules:
    - "main-module"
    - "secondary-module"`,
			expected: &ProjectConfig{
				Project: struct {
					Name    string   `yaml:"name"`
					Modules []string `yaml:"modules"`
				}{
					Name:    "test-project",
					Modules: []string{"main-module", "secondary-module"},
				},
			},
			wantErr: false,
		},
		{
			name: "minimal config",
			yaml: `project:
  name: "simple-project"
  modules:
    - "single-module"`,
			expected: &ProjectConfig{
				Project: struct {
					Name    string   `yaml:"name"`
					Modules []string `yaml:"modules"`
				}{
					Name:    "simple-project",
					Modules: []string{"single-module"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing project name",
			yaml: `project:
  modules:
    - "test-module"`,
			wantErr: true,
		},
		{
			name: "missing modules",
			yaml: `project:
  name: "test-project"`,
			wantErr: true,
		},
		{
			name: "empty module name",
			yaml: `project:
  name: "test-project"
  modules:
    - ""
    - "valid-module"`,
			wantErr: true,
		},
		{
			name: "invalid yaml",
			yaml: `project:
  name: "test-project"
  modules:
    - "test-module"
invalid: yaml: content`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test-project-*.yml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write YAML content
			if _, err := tmpFile.WriteString(tt.yaml); err != nil {
				t.Fatalf("Failed to write YAML content: %v", err)
			}
			tmpFile.Close()

			// Test loading from file
			config, err := GetProjectConfigFromFile(tmpFile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && config != nil {
				if config.Project.Name != tt.expected.Project.Name {
					t.Errorf("Project name = %v, want %v", config.Project.Name, tt.expected.Project.Name)
				}
				if len(config.Project.Modules) != len(tt.expected.Project.Modules) {
					t.Errorf("Modules length = %v, want %v", len(config.Project.Modules), len(tt.expected.Project.Modules))
				}
				for i, module := range config.Project.Modules {
					if module != tt.expected.Project.Modules[i] {
						t.Errorf("Module[%d] = %v, want %v", i, module, tt.expected.Project.Modules[i])
					}
				}
			}
		})
	}
}

func TestConfigProvider(t *testing.T) {
	cp := NewConfigProvider()

	// Test initial state
	if cp.HasConfig() {
		t.Error("Expected no config initially")
	}
	if cp.GetProjectName() != "" {
		t.Error("Expected empty project name initially")
	}
	if cp.GetModuleName() != "" {
		t.Error("Expected empty module name initially")
	}
	if len(cp.GetAllModules()) != 0 {
		t.Error("Expected empty modules initially")
	}
}

func TestFindProjectConfigFile(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test-project-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create .project.yml in tempDir
	configFile := filepath.Join(tempDir, ".project.yml")
	configContent := `project:
  name: "test-project"
  modules:
    - "test-module"`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change to subdirectory and test finding config
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to subdir: %v", err)
	}

	// Test finding config from subdirectory
	foundPath, err := findProjectConfigFile()
	if err != nil {
		t.Fatalf("findProjectConfigFile() error = %v", err)
	}
	if foundPath != configFile {
		t.Errorf("findProjectConfigFile() = %v, want %v", foundPath, configFile)
	}

	// Test not finding config when none exists
	if err := os.Remove(configFile); err != nil {
		t.Fatalf("Failed to remove config file: %v", err)
	}

	foundPath, err = findProjectConfigFile()
	if err != nil {
		t.Fatalf("findProjectConfigFile() error = %v", err)
	}
	if foundPath != "" {
		t.Errorf("findProjectConfigFile() = %v, want empty string", foundPath)
	}
}

func TestConfigProviderIntegration(t *testing.T) {
	// Create a temporary directory with .project.yml
	tempDir, err := os.MkdirTemp("", "test-project-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .project.yml
	configFile := filepath.Join(tempDir, ".project.yml")
	configContent := `project:
  name: "integration-test"
  modules:
    - "primary-module"
    - "secondary-module"
    - "tertiary-module"`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test loading configuration
	cp := NewConfigProvider()
	config, err := cp.LoadProjectConfig()
	if err != nil {
		t.Fatalf("LoadProjectConfig() error = %v", err)
	}
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	// Test configuration values
	if !cp.HasConfig() {
		t.Error("Expected config to be loaded")
	}
	if cp.GetProjectName() != "integration-test" {
		t.Errorf("GetProjectName() = %v, want integration-test", cp.GetProjectName())
	}
	if cp.GetModuleName() != "primary-module" {
		t.Errorf("GetModuleName() = %v, want primary-module", cp.GetModuleName())
	}
	allModules := cp.GetAllModules()
	expectedModules := []string{"primary-module", "secondary-module", "tertiary-module"}
	if len(allModules) != len(expectedModules) {
		t.Errorf("GetAllModules() length = %v, want %v", len(allModules), len(expectedModules))
	}
	for i, module := range allModules {
		if module != expectedModules[i] {
			t.Errorf("GetAllModules()[%d] = %v, want %v", i, module, expectedModules[i])
		}
	}
}