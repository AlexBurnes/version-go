package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const VERSION = "0.1.0"

var (
	useRelease bool
	debug      bool
	verbose    bool
)

// Custom error types for specific git-related issues
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

func init() {
	flag.BoolVar(&useRelease, "release", false, "use commit number as release number")
	flag.BoolVar(&useRelease, "r", false, "use commit number as release number (shorthand)")
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.BoolVar(&debug, "d", false, "enable debug output (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&verbose, "v", false, "verbose output (shorthand)")
}

// runCommand executes a command and returns its output, with verbose logging if enabled
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	
	// In verbose mode, print the command being executed
	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s %s\n", name, strings.Join(args, " "))
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
	if verbose {
		fmt.Fprintf(os.Stderr, "+ which git\n")
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
	// First check if git is available and if we're in a git repo
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
	
	if len(strings.TrimSpace(string(output))) == 0 {
		return &NoGitTagsError{}
	}
	return nil
}

// runGitCommand executes a git command and returns its output
func runGitCommand(args ...string) (string, error) {
	// First check if git is available and if we're in a git repo
	if err := checkGitAvailable(); err != nil {
		return "", err
	}
	if err := checkGitRepo(); err != nil {
		return "", err
	}
	
	return runCommand("git", args...)
}

func getVersion() (string, error) {
	if err := checkGitTags(); err != nil {
		return "", err
	}

	output, err := runGitCommand("describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
	if err != nil {
		return "", err
	}
	
	version := strings.TrimPrefix(output, "v")
	version = regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+)\-`).ReplaceAllString(version, "${1}~")
	return version, nil
}

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

func getRelease() (string, error) {
	if !useRelease {
		return "1", nil
	}

	// Get the current tag
	tag, err := runCommand("git", "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags")
	if err != nil {
		return "", err
	}

	// Count commits since that tag
	output, err := runCommand("git", "rev-list", tag + "..HEAD", "--count")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

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

func printHelp() {
	fmt.Printf(`version %s - describe project version and release using git describe command

Usage:
    version [-h|--help] [-v|--version] [-r|--release] project|version|release|full

Options:
    -h, --help      print this help and exit
    -V, --version   print script version and exit
    -v, --verbose   verbose output
    -d, --debug     debug output
    -r, --release   use commit number as release number, default is no and release is 1

Commands:
    project         print project name
    module          print module name
    version         print project version
    release         print project release
    full            print full project name-version-release
`, VERSION)
}

func printError(err error) {
	if debug {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	switch err.(type) {
	case *GitNotFoundError, *NotGitRepoError, *NoGitTagsError:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func main() {
	versionFlag := flag.Bool("V", false, "print version")
	helpFlag := flag.Bool("h", false, "print help")
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if *versionFlag {
		fmt.Println(VERSION)
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		printHelp()
		os.Exit(1)
	}

	var err error
	var result string

	switch args[0] {
	case "version":
		result, err = getVersion()
	case "project":
		result, err = getProject()
	case "module":
		result, err = getModule()
	case "release":
		result, err = getRelease()
	case "full":
		result, err = getFull()
	default:
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		printError(err)
		os.Exit(1)
	}

	fmt.Println(result)
} 