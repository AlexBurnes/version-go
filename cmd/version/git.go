package main

import (
    "fmt"
    "os/exec"
    "regexp"
    "strings"
    
    "github.com/AlexBurnes/version-go/pkg/version"
)

// Custom error types for git-related issues
type GitNotFoundError struct{}
type NotGitRepoError struct{}
type NoGitTagsError struct{}

func (e *GitNotFoundError) Error() string {
    return "git command is not available - please install git and ensure it's in your PATH"
}

func (e *NotGitRepoError) Error() string {
    return "not a git repository - please run this command from within a git repository"
}

func (e *NoGitTagsError) Error() string {
    return "no tag defined for project - please create a version tag (e.g., v1.0.0) before running this command"
}

// runCommand executes a command and returns its output
func runCommand(name string, args ...string) (string, error) {
    cmd := exec.Command(name, args...)
    
    if verboseFlag {
        printInfo("+ %s %s", name, strings.Join(args, " "))
    }
    
    output, err := cmd.Output()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
            return "", fmt.Errorf("command failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
        }
        return "", fmt.Errorf("command failed: %v", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// checkGitAvailable verifies that git is installed and available
func checkGitAvailable() error {
    _, err := exec.LookPath("git")
    if err != nil {
        return &GitNotFoundError{}
    }
    if verboseFlag {
        printInfo("+ which git")
    }
    return nil
}

// checkGitRepo verifies that the current directory is a git repository
func checkGitRepo() error {
    _, err := runCommand("git", "rev-parse", "--git-dir")
    if err != nil {
        return &NotGitRepoError{}
    }
    return nil
}

// checkGitTags verifies that the repository has at least one version tag
func checkGitTags() error {
    if err := checkGitAvailable(); err != nil {
        return err
    }
    if err := checkGitRepo(); err != nil {
        return err
    }

    output, err := runCommand("git", "tag", "-l", "v[0-9]*")
    if err != nil {
        return err
    }
    
    if len(strings.TrimSpace(output)) == 0 {
        return &NoGitTagsError{}
    }
    return nil
}

// runGitCommand executes a git command and returns its output
func runGitCommand(args ...string) (string, error) {
    if err := checkGitAvailable(); err != nil {
        return "", err
    }
    if err := checkGitRepo(); err != nil {
        return "", err
    }
    
    return runCommand("git", args...)
}

// getVersion returns the current project version from git tags
func getVersion() (string, error) {
    if err := checkGitTags(); err != nil {
        return "", err
    }

    output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", err
    }
    
    versionStr := strings.TrimPrefix(output, "v")
    // Convert git tag format using the library
    versionStr = version.ConvertGitTag(versionStr)
    return versionStr, nil
}

// getProject returns the project name from git remote
func getProject() (string, error) {
    output, err := runGitCommand("remote", "-v")
    if err != nil {
        return "", err
    }
    
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if strings.Contains(line, "fetch") {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                continue
            }
            remote := parts[1]
            remote = strings.TrimSuffix(remote, ".git")
            
            // Handle SSH URLs (git@github.com:user/repo)
            if strings.Contains(remote, ":") {
                remote = strings.Split(remote, ":")[1]
            } else {
                // Handle HTTPS URLs (https://github.com/user/repo)
                if strings.Contains(remote, "//") {
                    remote = strings.Split(remote, "//")[1]
                }
                parts := strings.SplitN(remote, "/", 2)
                if len(parts) > 1 {
                    remote = parts[1]
                }
            }
            
            // Convert slashes to dashes
            remote = strings.ReplaceAll(remote, "/", "-")
            
            // Remove any prefix matching --[^-]+-
            re := regexp.MustCompile(`^--[^-]+-`)
            remote = re.ReplaceAllString(remote, "")
            
            return remote, nil
        }
    }
    return "", fmt.Errorf("no git remote found - please add a remote to your repository")
}

// getModule returns the module name from git remote
func getModule() (string, error) {
    output, err := runGitCommand("remote", "-v")
    if err != nil {
        return "", err
    }
    
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if strings.Contains(line, "fetch") {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                continue
            }
            remote := parts[1]
            remote = strings.TrimSuffix(remote, ".git")
            
            // Handle SSH URLs (git@github.com:user/repo)
            if strings.Contains(remote, ":") {
                remote = strings.Split(remote, ":")[1]
            } else {
                // Handle HTTPS URLs (https://github.com/user/repo)
                if strings.Contains(remote, "//") {
                    remote = strings.Split(remote, "//")[1]
                }
                parts := strings.SplitN(remote, "/", 2)
                if len(parts) > 1 {
                    remote = parts[1]
                }
            }
            
            // Get the last component of the path
            parts = strings.Split(remote, "/")
            if len(parts) > 0 {
                return parts[len(parts)-1], nil
            }
        }
    }
    return "", fmt.Errorf("no git remote found - please add a remote to your repository")
}

// getRelease returns the release number (always 1 for now)
func getRelease() (string, error) {
    return "1", nil
}

// getFull returns the full project name-version-release
func getFull() (string, error) {
    version, err := getVersion()
    if err != nil {
        return "", err
    }
    
    project, err := getProject()
    if err != nil {
        return "", err
    }
    
    release, err := getRelease()
    if err != nil {
        return "", err
    }
    
    return fmt.Sprintf("%s-%s-%s", project, version, release), nil
}

// getGitTags returns all version tags from git repository
func getGitTags() ([]string, error) {
    output, err := runGitCommand("tag", "-l", "v[0-9]*")
    if err != nil {
        return nil, err
    }
    
    if output == "" {
        return []string{}, nil
    }
    
    return strings.Split(strings.TrimSpace(output), "\n"), nil
}

// checkGreatest checks if the given version is the greatest among all git tags
func checkGreatest(versionStr string) (string, error) {
    // Parse current version using the library
    currentVer, err := version.Parse(versionStr)
    if err != nil {
        return "", fmt.Errorf("invalid version format: %s", versionStr)
    }

    // Get all git tags
    tags, err := getGitTags()
    if err != nil {
        return "", fmt.Errorf("failed to get git tags: %v", err)
    }

    if len(tags) == 0 {
        return fmt.Sprintf("Version %s is the greatest (no other tags found)", versionStr), nil
    }

    // Compare with all valid version tags
    for _, tag := range tags {
        // Skip if tag is the same as current version
        if tag == versionStr || tag == "v"+versionStr {
            continue
        }

        // Try to parse the tag (skip invalid ones)
        tagVer, err := version.Parse(tag)
        if err != nil {
            printDebug("Skipping invalid tag: %s", tag)
            continue
        }

        // If we find a greater version, return false
        if version.Compare(tagVer, currentVer) > 0 {
            return "", fmt.Errorf("version %s is not the greatest among all tags", versionStr)
        }
    }

    return fmt.Sprintf("Version %s is the greatest among all tags", versionStr), nil
}