package main

import (
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
            expected: "0.3.0",
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

func isGitRepository() bool {
    cmd := exec.Command("git", "rev-parse", "--git-dir")
    err := cmd.Run()
    return err == nil
}