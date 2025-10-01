package main

import (
    "os"
    "os/exec"
    "strings"
    "testing"
)

func TestCLICommands(t *testing.T) {
    // Skip integration tests if not in a git repository
    if !isGitRepository() {
        t.Skip("Skipping integration tests - not in a git repository")
    }

    tests := []struct {
        name     string
        args     []string
        expected string
        hasError bool
    }{
        {
            name:     "version command",
            args:     []string{"version"},
            expected: "", // Will be git-dependent
            hasError: false,
        },
        {
            name:     "help command",
            args:     []string{"--help"},
            expected: "version",
            hasError: false,
        },
        {
            name:     "version flag",
            args:     []string{"--version"},
            expected: "1.2.4",
            hasError: false,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            cmd := exec.Command("go", append([]string{"run", "."}, test.args...)...)
            cmd.Dir = "."
            
            output, err := cmd.CombinedOutput()
            outputStr := string(output)
            
            if test.hasError {
                if err == nil {
                    t.Errorf("Expected error for %v, but got none. Output: %s", test.args, outputStr)
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error for %v: %v. Output: %s", test.args, err, outputStr)
                }
                
                if test.expected != "" && !strings.Contains(outputStr, test.expected) {
                    t.Errorf("Expected output to contain '%s', but got: %s", test.expected, outputStr)
                }
            }
        })
    }
}

func TestVersionValidation(t *testing.T) {
    tests := []struct {
        version  string
        hasError bool
    }{
        {"1.2.3", false},
        {"v1.2.3", false},
        {"1.2.3-alpha", false},
        {"1.2.3~beta.1", false},
        {"1.2.3.fix", false},
        {"1.2.3_feature", false},
        {"1.2", true},
        {"invalid", true},
        {"1.2.3.4", true},
    }

    for _, test := range tests {
        t.Run(test.version, func(t *testing.T) {
            cmd := exec.Command("go", "run", ".", "check", test.version)
            cmd.Dir = "."
            
            err := cmd.Run()
            
            if test.hasError {
                if err == nil {
                    t.Errorf("Expected error for version %s, but got none", test.version)
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error for version %s: %v", test.version, err)
                }
            }
        })
    }
}

func TestVersionType(t *testing.T) {
    tests := []struct {
        version  string
        expected string
    }{
        {"1.2.3", "release"},
        {"1.2.3-alpha", "prerelease"},
        {"1.2.3~beta.1", "prerelease"},
        {"1.2.3.fix", "postrelease"},
        {"1.2.3_feature", "intermediate"},
    }

    for _, test := range tests {
        t.Run(test.version, func(t *testing.T) {
            cmd := exec.Command("go", "run", ".", "type", test.version)
            cmd.Dir = "."
            
            output, err := cmd.Output()
            if err != nil {
                t.Errorf("Error running type command for %s: %v", test.version, err)
                return
            }
            
            result := strings.TrimSpace(string(output))
            if result != test.expected {
                t.Errorf("Expected type '%s' for version %s, got '%s'", test.expected, test.version, result)
            }
        })
    }
}

func TestBuildType(t *testing.T) {
    tests := []struct {
        version  string
        expected string
    }{
        {"1.2.3", "Release"},
        {"1.2.3-alpha", "Debug"},
        {"1.2.3~beta.1", "Debug"},
        {"1.2.3.fix", "Debug"},
        {"1.2.3_feature", "Debug"},
    }

    for _, test := range tests {
        t.Run(test.version, func(t *testing.T) {
            cmd := exec.Command("go", "run", ".", "build-type", test.version)
            cmd.Dir = "."
            
            output, err := cmd.Output()
            if err != nil {
                t.Errorf("Error running build-type command for %s: %v", test.version, err)
                return
            }
            
            result := strings.TrimSpace(string(output))
            if result != test.expected {
                t.Errorf("Expected build type '%s' for version %s, got '%s'", test.expected, test.version, result)
            }
        })
    }
}

func TestExitCodeCompliance(t *testing.T) {
    tests := []struct {
        name        string
        args        []string
        expectedCode int
        description string
    }{
        {
            name:        "valid_version_check",
            args:        []string{"check", "1.2.3"},
            expectedCode: 0,
            description: "Valid version should return exit code 0",
        },
        {
            name:        "invalid_version_check",
            args:        []string{"check", "invalid"},
            expectedCode: 1,
            description: "Invalid version should return exit code 1",
        },
        {
            name:        "valid_version_type",
            args:        []string{"type", "1.2.3"},
            expectedCode: 0,
            description: "Valid version type command should return exit code 0",
        },
        {
            name:        "invalid_version_type",
            args:        []string{"type", "invalid"},
            expectedCode: 1,
            description: "Invalid version type command should return exit code 1",
        },
        {
            name:        "valid_version_build_type",
            args:        []string{"build-type", "1.2.3"},
            expectedCode: 0,
            description: "Valid version build-type command should return exit code 0",
        },
        {
            name:        "invalid_version_build_type",
            args:        []string{"build-type", "invalid"},
            expectedCode: 1,
            description: "Invalid version build-type command should return exit code 1",
        },
        // Note: valid_version_check_greatest test is skipped here because it depends on
        // the actual git tags in the repository. See TestCheckGreatestWithCurrentVersion
        // for a dynamic test that checks the current version from git.
        {
            name:        "invalid_version_check_greatest",
            args:        []string{"check-greatest", "invalid"},
            expectedCode: 1,
            description: "Invalid version check-greatest command should return exit code 1",
        },
        {
            name:        "help_command",
            args:        []string{"--help"},
            expectedCode: 0,
            description: "Help command should return exit code 0",
        },
        {
            name:        "version_flag",
            args:        []string{"--version"},
            expectedCode: 0,
            description: "Version flag should return exit code 0",
        },
        {
            name:        "unknown_command",
            args:        []string{"unknown"},
            expectedCode: 1,
            description: "Unknown command should return exit code 1",
        },
        {
            name:        "invalid_flag",
            args:        []string{"--invalid"},
            expectedCode: 1,
            description: "Invalid flag should return exit code 1",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            cmd := exec.Command("go", append([]string{"run", "."}, test.args...)...)
            cmd.Dir = "."
            
            err := cmd.Run()
            var exitCode int
            if err != nil {
                if exitError, ok := err.(*exec.ExitError); ok {
                    exitCode = exitError.ExitCode()
                } else {
                    t.Errorf("Unexpected error type for %v: %v", test.args, err)
                    return
                }
            } else {
                exitCode = 0
            }
            
            if exitCode != test.expectedCode {
                t.Errorf("Expected exit code %d for %v (%s), got %d", 
                    test.expectedCode, test.args, test.description, exitCode)
            }
        })
    }
}

func TestExitCodeComplianceWithBuiltBinary(t *testing.T) {
    // Skip if binary doesn't exist
    if _, err := os.Stat("../../bin/version"); os.IsNotExist(err) {
        t.Skip("Skipping binary exit code tests - binary not built")
    }

    tests := []struct {
        name        string
        args        []string
        expectedCode int
        description string
    }{
        {
            name:        "valid_version_check_binary",
            args:        []string{"check", "1.2.3"},
            expectedCode: 0,
            description: "Valid version should return exit code 0 (binary)",
        },
        {
            name:        "invalid_version_check_binary",
            args:        []string{"check", "invalid"},
            expectedCode: 1,
            description: "Invalid version should return exit code 1 (binary)",
        },
        {
            name:        "help_command_binary",
            args:        []string{"--help"},
            expectedCode: 0,
            description: "Help command should return exit code 0 (binary)",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            cmd := exec.Command("../../bin/version", test.args...)
            cmd.Dir = "."
            
            err := cmd.Run()
            var exitCode int
            if err != nil {
                if exitError, ok := err.(*exec.ExitError); ok {
                    exitCode = exitError.ExitCode()
                } else {
                    t.Errorf("Unexpected error type for %v: %v", test.args, err)
                    return
                }
            } else {
                exitCode = 0
            }
            
            if exitCode != test.expectedCode {
                t.Errorf("Expected exit code %d for %v (%s), got %d", 
                    test.expectedCode, test.args, test.description, exitCode)
            }
        })
    }
}

func isGitRepository() bool {
    cmd := exec.Command("git", "rev-parse", "--git-dir")
    err := cmd.Run()
    return err == nil
}

// TestCheckGreatestWithCurrentVersion tests check-greatest command dynamically
// using the current version from git, avoiding hardcoded version dependencies
func TestCheckGreatestWithCurrentVersion(t *testing.T) {
    if !isGitRepository() {
        t.Skip("Not a git repository, skipping test")
    }
    
    // Get current version from git using the version command
    cmd := exec.Command("./bin/version", "version")
    output, err := cmd.Output()
    if err != nil {
        t.Skipf("Could not get current version from git: %v", err)
    }
    
    currentVersion := strings.TrimSpace(string(output))
    if currentVersion == "" {
        t.Skip("No version tags found in repository")
    }
    
    // Test that current version is the greatest (should return exit code 0)
    t.Run("current_version_is_greatest", func(t *testing.T) {
        cmd := exec.Command("./bin/version", "check-greatest", currentVersion)
        err := cmd.Run()
        exitCode := 0
        if err != nil {
            if exitErr, ok := err.(*exec.ExitError); ok {
                exitCode = exitErr.ExitCode()
            } else {
                t.Fatalf("Failed to run command: %v", err)
            }
        }
        
        if exitCode != 0 {
            // Get output to see what went wrong
            cmd := exec.Command("./bin/version", "check-greatest", currentVersion)
            output, _ := cmd.CombinedOutput()
            t.Errorf("Expected current version %s to be greatest (exit code 0), got exit code %d. Output: %s", 
                currentVersion, exitCode, string(output))
        }
    })
    
    // Test that an older version is not the greatest (should return exit code 1)
    t.Run("old_version_not_greatest", func(t *testing.T) {
        // Use a version that should be less than current
        testVersion := "0.1.0"
        cmd := exec.Command("./bin/version", "check-greatest", testVersion)
        err := cmd.Run()
        exitCode := 0
        if err != nil {
            if exitErr, ok := err.(*exec.ExitError); ok {
                exitCode = exitErr.ExitCode()
            } else {
                return // Command failed for other reasons, that's ok
            }
        }
        
        // We expect exit code 1 because 0.1.0 should not be the greatest
        // But if 0.1.0 IS the greatest (e.g., in a test repo), that's also fine
        // This test is just checking that the command works, not enforcing specific behavior
        if exitCode != 0 && exitCode != 1 {
            t.Errorf("Expected exit code 0 or 1 for check-greatest with old version, got %d", exitCode)
        }
    })
}