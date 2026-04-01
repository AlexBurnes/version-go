// Package version provides git integration functionality for retrieving version information.
// It provides functions to get version from git tags in a way that's safe for library consumers.
package version

import (
    "fmt"
    "os/exec"
    "regexp"
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
func runGitCommand(dir string, args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    if (dir != "") {
        cmd.Dir = dir
    }
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
func checkGitRepo(dir string) error {
    _, err := runGitCommand(dir, "rev-parse", "--git-dir")
    if err != nil {
        return &GitError{
            Type:    "not_repo",
            Message: "not a git repository - please run this from within a git repository",
        }
    }
    return nil
}

// checkGitTags verifies that the repository has at least one version tag
func checkGitTags(dir string) error {
    if err := checkGitAvailable(); err != nil {
        return err
    }
    if err := checkGitRepo(dir); err != nil {
        return err
    }

    output, err := runGitCommand(dir, "tag", "-l", "v[0-9]*")
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
func GetVersion(dir string) (string, error) {
    if err := checkGitTags(dir); err != nil {
        return "", err
    }

    output, err := runGitCommand(dir, "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
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
func GetVersionWithPrefix(dir string) (string, error) {
    if err := checkGitTags(dir); err != nil {
        return "", err
    }

    output, err := runGitCommand(dir, "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
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
func GetRawTag(dir string) (string, error) {
    if err := checkGitTags(dir); err != nil {
        return "", err
    }

    output, err := runGitCommand(dir, "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
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
func GetRawVersion(dir string) (string, error) {
    if err := checkGitTags(dir); err != nil {
        return "", err
    }

    output, err := runGitCommand(dir, "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
    if err != nil {
        return "", fmt.Errorf("failed to get version from git: %v", err)
    }
    
    // Remove 'v' prefix but do NOT convert tag format
    versionStr := strings.TrimPrefix(output, "v")
    return versionStr, nil
}


// remoteURLPath returns the path component of the fetch URL for the best available
// remote (origin preferred, any other remote as fallback).
func remoteURLPath(dir string) (string, error) {
    if err := checkGitAvailable(); err != nil {
        return "", err
    }
    if err := checkGitRepo(dir); err != nil {
        return "", err
    }
    output, err := runGitCommand(dir, "remote", "-v")
    if err != nil {
        return "", fmt.Errorf("failed to get git remote: %v", err)
    }
    lines := strings.Split(output, "\n")
    for _, wantOrigin := range []bool{true, false} {
        for _, line := range lines {
            if !strings.Contains(line, "fetch") {
                continue
            }
            if wantOrigin && !strings.HasPrefix(strings.TrimSpace(line), "origin") {
                continue
            }
            parts := strings.Fields(line)
            if len(parts) < 2 {
                continue
            }
            return parseRemoteURL(parts[1]), nil
        }
    }
    return "", fmt.Errorf("no git remote found - please add a remote to your repository")
}

// parseRemoteURL extracts the path component (user/repo) from a git remote URL.
// Handles SSH (git@host:path) and HTTPS (https://host/path) formats.
func parseRemoteURL(url string) string {
    url = strings.TrimSuffix(url, ".git")
    if strings.Contains(url, ":") && !strings.HasPrefix(url, "http") {
        // SSH: git@github.com:user/repo
        return strings.SplitN(url, ":", 2)[1]
    }
    if strings.Contains(url, "//") {
        // HTTPS: https://github.com/user/repo
        url = strings.SplitN(url, "//", 2)[1]
        parts := strings.SplitN(url, "/", 2)
        if len(parts) > 1 {
            return parts[1]
        }
    }
    return url
}

// GetProjectFromGit returns the project name derived solely from the git remote URL
// for the repository at dir (empty = current directory), without consulting .project.yml.
// Slashes in the remote path are converted to dashes and leading "--[^-]+-" prefixes stripped.
func GetProjectFromGit(dir string) (string, error) {
    urlPath, err := remoteURLPath(dir)
    if err != nil {
        return "", err
    }
    name := strings.ReplaceAll(urlPath, "/", "-")
    re := regexp.MustCompile(`^--[^-]+-`)
    return re.ReplaceAllString(name, ""), nil
}

// GetProject returns the project name for the repository at dir (empty = current directory).
// It reads .project.yml first (searching upward from dir); falls back to GetProjectFromGit.
func GetProject(dir string) (string, error) {
    config, _ := loadConfigFrom(dir)
    if config != nil && config.Project.Name != "" {
        return config.Project.Name, nil
    }
    return GetProjectFromGit(dir)
}

// GetModuleFromGit returns the primary module name derived solely from the git remote URL
// for the repository at dir (empty = current directory), without consulting .project.yml.
// The module name is the last path component of the remote URL.
func GetModuleFromGit(dir string) (string, error) {
    urlPath, err := remoteURLPath(dir)
    if err != nil {
        return "", err
    }
    parts := strings.Split(urlPath, "/")
    return parts[len(parts)-1], nil
}

// GetModule returns the primary module name for the repository at dir (empty = current directory).
// It reads the first module from .project.yml; falls back to GetModuleFromGit.
func GetModule(dir string) (string, error) {
    config, _ := loadConfigFrom(dir)
    if config != nil && len(config.Project.Modules) > 0 {
        return config.Project.Modules[0], nil
    }
    return GetModuleFromGit(dir)
}

// GetModules returns all module names for the repository at dir (empty = current directory).
// It reads all modules from .project.yml; falls back to a single-element slice containing
// the repository name from the git remote URL.
func GetModules(dir string) ([]string, error) {
    config, _ := loadConfigFrom(dir)
    if config != nil && len(config.Project.Modules) > 0 {
        return config.Project.Modules, nil
    }
    urlPath, err := remoteURLPath(dir)
    if err != nil {
        return nil, err
    }
    parts := strings.Split(urlPath, "/")
    return []string{parts[len(parts)-1]}, nil
}

// GetRelease returns the release number. Release is always 1.
func GetRelease() string {
    return "1"
}

// GetFull returns the full project identifier as "project-version-release" for the
// repository at dir (empty = current directory).
func GetFull(dir string) (string, error) {
    project, err := GetProject(dir)
    if err != nil {
        return "", err
    }
    ver, err := GetVersion(dir)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s-%s-%s", project, ver, GetRelease()), nil
}

// CheckGreatest returns nil if versionStr is the greatest version tag reachable from
// the current branch in the repository at dir (empty = current directory).
// Returns an error naming the greater tag if a higher version exists.
func CheckGreatest(versionStr string, dir string) error {
    currentVer, err := Parse(versionStr)
    if err != nil {
        return fmt.Errorf("invalid version format: %s", versionStr)
    }
    if err := checkGitTags(dir); err != nil {
        return err
    }
    branch, err := runGitCommand(dir, "branch", "--show-current")
    if err != nil {
        return fmt.Errorf("failed to get current branch: %v", err)
    }
    output, err := runGitCommand(dir, "tag", "-l", "v[0-9]*", "--merged", branch)
    if err != nil {
        return err
    }
    if strings.TrimSpace(output) == "" {
        return nil
    }
    var greatestVer *Version
    var greatestTag string
    for _, tag := range strings.Split(strings.TrimSpace(output), "\n") {
        if tag == versionStr || tag == "v"+versionStr {
            continue
        }
        tagVer, err := Parse(tag)
        if err != nil {
            continue
        }
        if greatestVer == nil || Compare(tagVer, greatestVer) > 0 {
            greatestVer = tagVer
            greatestTag = tag
        }
    }
    if greatestVer != nil && Compare(greatestVer, currentVer) > 0 {
        return fmt.Errorf("version %s is not the greatest: tag %s is greater", versionStr, greatestTag)
    }
    return nil
}

