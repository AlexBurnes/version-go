package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestConfigCLIOptions(t *testing.T) {
	// Build the binary for testing
	binaryPath := filepath.Join("bin", "version")
	
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Change to project root
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
		description    string
	}{
		{
			name:           "default_behavior_with_project_yml",
			args:           []string{"project"},
			expectedOutput: "buildfab-version",
			expectedError:  false,
			description:    "Should use .project.yml when present (default behavior)",
		},
		{
			name:           "default_behavior_with_project_yml_module",
			args:           []string{"module"},
			expectedOutput: "version",
			expectedError:  false,
			description:    "Should use .project.yml module when present (default behavior)",
		},
		{
			name:           "git_flag_force_git_detection",
			args:           []string{"--git", "project"},
			expectedOutput: "AlexBurnes-version-go",
			expectedError:  false,
			description:    "Should force git detection when --git flag is used",
		},
		{
			name:           "git_flag_force_git_detection_module",
			args:           []string{"--git", "module"},
			expectedOutput: "version-go",
			expectedError:  false,
			description:    "Should force git detection for module when --git flag is used",
		},
		{
			name:           "config_flag_custom_file",
			args:           []string{"--config", "test/project_name.yml", "project"},
			expectedOutput: "test-project-name",
			expectedError:  false,
			description:    "Should use custom config file when --config flag is used",
		},
		{
			name:           "config_flag_custom_file_module",
			args:           []string{"--config", "test/project_name.yml", "module"},
			expectedOutput: "test-module",
			expectedError:  false,
			description:    "Should use custom config file for module when --config flag is used",
		},
		{
			name:           "config_flag_single_module",
			args:           []string{"--config", "test/project_module_single.yml", "module"},
			expectedOutput: "single-module",
			expectedError:  false,
			description:    "Should use single module from custom config file",
		},
		{
			name:           "config_flag_multiple_modules",
			args:           []string{"--config", "test/project_module_multiple.yml", "module"},
			expectedOutput: "primary-module",
			expectedError:  false,
			description:    "Should use first module from multiple modules config file",
		},
		{
			name:           "config_flag_multiple_modules_project",
			args:           []string{"--config", "test/project_module_multiple.yml", "project"},
			expectedOutput: "multi-module-project",
			expectedError:  false,
			description:    "Should use project name from multiple modules config file",
		},
		{
			name:           "config_flag_nonexistent_file",
			args:           []string{"--config", "test/nonexistent.yml", "project"},
			expectedOutput: "",
			expectedError:  true,
			description:    "Should error when config file doesn't exist",
		},
		{
			name:           "conflicting_flags",
			args:           []string{"--config", "test/project_name.yml", "--git", "project"},
			expectedOutput: "",
			expectedError:  true,
			description:    "Should error when both --config and --git flags are used",
		},
		{
			name:           "modules_command_default_behavior",
			args:           []string{"modules"},
			expectedOutput: "version",
			expectedError:  false,
			description:    "Should use .project.yml modules when present (default behavior)",
		},
		{
			name:           "modules_command_git_flag",
			args:           []string{"--git", "modules"},
			expectedOutput: "version-go",
			expectedError:  false,
			description:    "Should force git detection for modules when --git flag is used",
		},
		{
			name:           "modules_command_single_module",
			args:           []string{"--config", "test/project_module_single.yml", "modules"},
			expectedOutput: "single-module",
			expectedError:  false,
			description:    "Should use single module from custom config file",
		},
		{
			name:           "modules_command_multiple_modules",
			args:           []string{"--config", "test/project_module_multiple.yml", "modules"},
			expectedOutput: "primary-module\nsecondary-module\ntertiary-module",
			expectedError:  false,
			description:    "Should use all modules from multiple modules config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the binary with the specified arguments
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.Output()
			
			// Check for expected error
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", string(output))
				}
				return
			}
			
			// Check for unexpected error
			if err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, string(output))
				return
			}
			
			// Check output
			actualOutput := string(output)
			if tt.expectedOutput != "" {
				actualOutput = actualOutput[:len(actualOutput)-1] // Remove trailing newline
			}
			
			if actualOutput != tt.expectedOutput {
				t.Errorf("Expected output: %q, got: %q", tt.expectedOutput, actualOutput)
			}
		})
	}
}

func TestConfigCLIOptionsWithDebug(t *testing.T) {
	// Build the binary for testing
	binaryPath := filepath.Join("bin", "version")
	
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Change to project root
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	
	tests := []struct {
		name        string
		args        []string
		description string
	}{
		{
			name:        "debug_default_behavior",
			args:        []string{"--debug", "project"},
			description: "Should show debug output for default behavior",
		},
		{
			name:        "debug_git_flag",
			args:        []string{"--debug", "--git", "project"},
			description: "Should show debug output for git flag",
		},
		{
			name:        "debug_config_flag",
			args:        []string{"--debug", "--config", "test/project_name.yml", "project"},
			description: "Should show debug output for config flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the binary with the specified arguments
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			
			if err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, string(output))
				return
			}
			
			// Check that debug output is present
			outputStr := string(output)
			if !contains(outputStr, "#DEBUG") {
				t.Errorf("Expected debug output but got: %s", outputStr)
			}
		})
	}
}

func TestConfigFileValidation(t *testing.T) {
	// Change to project root directory for tests
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Change to project root
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	
	// Create invalid config files for testing
	invalidConfigs := []struct {
		filename string
		content  string
		desc     string
	}{
		{
			filename: "test/invalid_missing_name.yml",
			content: `project:
  modules:
    - "test-module"`,
			desc: "Missing project name",
		},
		{
			filename: "test/invalid_missing_modules.yml",
			content: `project:
  name: "test-project"`,
			desc: "Missing modules",
		},
		{
			filename: "test/invalid_empty_module.yml",
			content: `project:
  name: "test-project"
  modules:
    - ""
    - "valid-module"`,
			desc: "Empty module name",
		},
		{
			filename: "test/invalid_yaml.yml",
			content: `project:
  name: "test-project"
  modules:
    - "test-module"
invalid: yaml: content`,
			desc: "Invalid YAML syntax",
		},
	}

	// Create invalid config files
	for _, config := range invalidConfigs {
		err := os.WriteFile(config.filename, []byte(config.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid config file %s: %v", config.filename, err)
		}
	}
	defer func() {
		// Clean up invalid config files
		for _, config := range invalidConfigs {
			os.Remove(config.filename)
		}
	}()

	binaryPath := filepath.Join("bin", "version")

	for _, config := range invalidConfigs {
		t.Run(config.desc, func(t *testing.T) {
			// Test with --config flag (should fail)
			cmd := exec.Command(binaryPath, "--config", config.filename, "project")
			output, err := cmd.Output()
			
			if err == nil {
				t.Errorf("Expected error for invalid config file %s, but got success. Output: %s", config.filename, string(output))
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && contains(s[1:], substr)
}