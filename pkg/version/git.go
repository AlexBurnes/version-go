// Package version provides git integration functionality for retrieving version information.
// It provides functions to get version from git tags in a way that's safe for library consumers.
package version

import (
    "fmt"
    "os/exec"
    "strings"
)

// GitError represents different types of git-related errors
type GitError struct {
    Type    string // "not_found", "not_repo", "no_tags"
    Message string
}

func (e *GitError) Error() string {
    return e.Message
}

// IsGitNotFound returns true if the error indicates git is not available
func IsGitNotFound(err error) bool {
    if gitErr, ok := err.(*GitError); ok {
        return gitErr.Type == "not_found"
    }
    return false
}

// IsNotGitRepo returns true if the error indicates not in a git repository
func IsNotGitRepo(err error) bool {
    if gitErr, ok := err.(*GitError); ok {
        return gitErr.Type == "not_repo"
    }
    return false
}

// IsNoGitTags returns true if the error indicates no version tags found
func IsNoGitTags(err error) bool {
    if gitErr, ok := err.(*GitError); ok {
        return gitErr.Type == "no_tags"
    }
    return false
}

// checkGitAvailable verifies that git is installed and available
func checkGitAvailable() error {
    _, err := exec.LookPath("git")
    if err != nil {
        return &GitError{
            Type:    "not_found",
            Message: "git command is not available - please install git and ensure it's in your PATH",
        }
    }
    return nil
}

// runGitCommand executes a git command and returns its output
func runGitCommand(args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    output, err := cmd.Output()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
            return "", fmt.Errorf("git command failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
        }
        return "", fmt.Errorf("git command failed: %v", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// checkGitRepo verifies that the current directory is a git repository
func checkGitRepo() error {
    _, err := runGitCommand("rev-parse", "--git-dir")
    if err != nil {
        return &GitError{
            Type:    "not_repo",
            Message: "not a git repository - please run this from within a git repository",
        }
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

    output, err := runGitCommand("tag", "-l", "v[0-9]*")
    if err != nil {
        return err
    }
    
    if len(strings.TrimSpace(output)) == 0 {
        return &GitError{
            Type:    "no_tags",
            Message: "no version tags found - please create a version tag (e.g., v1.0.0)",
        }
    }
    return nil
}

// GetVersion returns the current project version from git tags.
// It retrieves the most recent version tag that matches the pattern v[0-9]*.
// The returned version string has the 'v' prefix removed and is converted
// from git tag format if necessary.
//
// Returns an error if:
//   - git is not available
//   - not in a git repository
//   - no version tags are found
//   - git command execution fails
//
// Example usage:
//
//	version, err := version.GetVersion()
//	if err != nil {
//	    if version.IsGitNotFound(err) {
//	        fmt.Println("Git is not installed")
//	    } else if version.IsNotGitRepo(err) {
//	        fmt.Println("Not in a git repository")
//	    } else if version.IsNoGitTags(err) {
//	        fmt.Println("No version tags found")
//	    } else {
//	        fmt.Printf("Error: %v\n", err)
//	    }
//	    return
//	}
//	fmt.Printf("Current version: %s\n", version)
func GetVersion() (string, error) {
    if err := checkGitTags(); err != nil {
        return "", err
    }

    output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", fmt.Errorf("failed to get version from git: %v", err)
    }
    
    versionStr := strings.TrimPrefix(output, "v")
    // Convert git tag format using the library
    versionStr = ConvertGitTag(versionStr)
    return versionStr, nil
}

// GetVersionWithPrefix returns the current project version from git tags with the 'v' prefix.
// This is the same as GetVersion but preserves the 'v' prefix in the version string.
// Git tag format is converted from x.y.z-(remainder) to x.y.z~(remainder) for consistency.
//
// Example usage:
//
//	version, err := version.GetVersionWithPrefix()
//	if err != nil {
//	    fmt.Printf("Error: %v\n", err)
//	    return
//	}
//	fmt.Printf("Current version: %s\n", version) // e.g., "v1.2.3~pre.1"
func GetVersionWithPrefix() (string, error) {
    if err := checkGitTags(); err != nil {
        return "", err
    }

    output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", fmt.Errorf("failed to get version from git: %v", err)
    }
    
    // Convert git tag format (e.g., v0.7.14-pre.1 -> v0.7.14~pre.1)
    versionStr := ConvertGitTag(output)
    return versionStr, nil
}

// GetRawTag returns the current git tag without any transformations.
// This function returns the exact tag string as it appears in git, without:
//   - Removing the 'v' prefix
//   - Converting '-' to '~' delimiter
//
// This is useful for applications that need the original git tag format.
//
// Returns an error if:
//   - git is not available
//   - not in a git repository
//   - no version tags are found
//   - git command execution fails
//
// Example usage:
//
//	tag, err := version.GetRawTag()
//	if err != nil {
//	    if version.IsGitNotFound(err) {
//	        fmt.Println("Git is not installed")
//	    } else if version.IsNotGitRepo(err) {
//	        fmt.Println("Not in a git repository")
//	    } else if version.IsNoGitTags(err) {
//	        fmt.Println("No version tags found")
//	    } else {
//	        fmt.Printf("Error: %v\n", err)
//	    }
//	    return
//	}
//	fmt.Printf("Current git tag: %s\n", tag) // e.g., "v1.2.3-pre.1" (exactly as in git)
func GetRawTag() (string, error) {
    if err := checkGitTags(); err != nil {
        return "", err
    }

    output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", fmt.Errorf("failed to get tag from git: %v", err)
    }
    
    // Return raw tag without any transformations
    return output, nil
}

// GetRawVersion returns the current git tag without transformations and without the 'v' prefix.
// This function returns the version string as it appears in git, but with the 'v' prefix removed.
// Unlike GetVersion(), it does NOT convert '-' to '~' delimiter.
//
// This is useful for applications that need the raw version format without any conversions.
//
// Returns an error if:
//   - git is not available
//   - not in a git repository
//   - no version tags are found
//   - git command execution fails
//
// Example usage:
//
//	version, err := version.GetRawVersion()
//	if err != nil {
//	    if version.IsGitNotFound(err) {
//	        fmt.Println("Git is not installed")
//	    } else if version.IsNotGitRepo(err) {
//	        fmt.Println("Not in a git repository")
//	    } else if version.IsNoGitTags(err) {
//	        fmt.Println("No version tags found")
//	    } else {
//	        fmt.Printf("Error: %v\n", err)
//	    }
//	    return
//	}
//	fmt.Printf("Current version: %s\n", version) // e.g., "1.2.3-pre.1" (without 'v', no conversion)
func GetRawVersion() (string, error) {
    if err := checkGitTags(); err != nil {
        return "", err
    }

    output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", fmt.Errorf("failed to get version from git: %v", err)
    }
    
    // Remove 'v' prefix but do NOT convert tag format
    versionStr := strings.TrimPrefix(output, "v")
    return versionStr, nil
}

