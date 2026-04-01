package main

import (
	"fmt"
	"strings"

	"github.com/AlexBurnes/version-go/pkg/version"
)

// getVersion returns the current project version from git tags.
func getVersion() (string, error) {
	return version.GetVersion("")
}

// getProject returns the project name from .project.yml or git remote.
func getProject() (string, error) {
	if gitFlag {
		printDebug("--git flag set, forcing git-based detection")
		return version.GetProjectFromGit("")
	}
	if configFile != "" {
		printDebug("--config flag set, using custom config file: %s", configFile)
		return getProjectFromCustomConfig(configFile)
	}
	return version.GetProject("")
}

// getProjectFromCustomConfig returns the project name from a specific config file.
func getProjectFromCustomConfig(configPath string) (string, error) {
	config, err := version.GetProjectConfigFromFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to load config file %s: %v", configPath, err)
	}
	if config.Project.Name == "" {
		return "", fmt.Errorf("project name not specified in config file %s", configPath)
	}
	printDebug("Using project name from config file %s: %s", configPath, config.Project.Name)
	return config.Project.Name, nil
}

// getModule returns the module name from .project.yml or git remote.
func getModule() (string, error) {
	if gitFlag {
		printDebug("--git flag set, forcing git-based detection")
		return version.GetModuleFromGit("")
	}
	if configFile != "" {
		printDebug("--config flag set, using custom config file: %s", configFile)
		return getModuleFromCustomConfig(configFile)
	}
	return version.GetModule("")
}

// getModuleFromCustomConfig returns the primary module name from a specific config file.
func getModuleFromCustomConfig(configPath string) (string, error) {
	config, err := version.GetProjectConfigFromFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to load config file %s: %v", configPath, err)
	}
	if len(config.Project.Modules) == 0 {
		return "", fmt.Errorf("no modules specified in config file %s", configPath)
	}
	printDebug("Using module name from config file %s: %s", configPath, config.Project.Modules[0])
	return config.Project.Modules[0], nil
}

// getModules returns all module names from .project.yml or a single module from git remote.
func getModules() (string, error) {
	if gitFlag {
		printDebug("--git flag set, forcing git-based detection")
		mod, err := version.GetModuleFromGit("")
		if err != nil {
			return "", err
		}
		return mod, nil
	}
	if configFile != "" {
		printDebug("--config flag set, using custom config file: %s", configFile)
		return getModulesFromCustomConfig(configFile)
	}
	mods, err := version.GetModules("")
	if err != nil {
		return "", err
	}
	return strings.Join(mods, "\n"), nil
}

// getModulesFromCustomConfig returns all module names from a specific config file.
func getModulesFromCustomConfig(configPath string) (string, error) {
	config, err := version.GetProjectConfigFromFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to load config file %s: %v", configPath, err)
	}
	if len(config.Project.Modules) == 0 {
		return "", fmt.Errorf("no modules specified in config file %s", configPath)
	}
	printDebug("Using modules from config file %s: %v", configPath, config.Project.Modules)
	return strings.Join(config.Project.Modules, "\n"), nil
}

// getRelease returns the release number.
func getRelease() (string, error) {
	return version.GetRelease(), nil
}

// getFull returns the full project name-version-release.
func getFull() (string, error) {
	return version.GetFull("")
}

// checkGreatest checks if the given version is the greatest among tags on the current branch.
func checkGreatest(versionStr string) (string, error) {
	if err := version.CheckGreatest(versionStr, ""); err != nil {
		return "", err
	}
	return fmt.Sprintf("Version %s is the greatest among tags on current branch", versionStr), nil
}

// getPlatform returns the current platform (GOOS value).
func getPlatform() (string, error) {
	return version.GetPlatform(), nil
}

// getArch returns the current architecture (GOARCH value).
func getArch() (string, error) {
	return version.GetArch(), nil
}

// getOS returns the current operating system (user-friendly format).
func getOS() (string, error) {
	return version.GetOS(), nil
}

// getOSVersion returns the current operating system version.
func getOSVersion() (string, error) {
	return version.GetOSVersion(), nil
}

// getCPU returns the number of logical CPUs.
func getCPU() (string, error) {
	return fmt.Sprintf("%d", version.GetNumCPU()), nil
}
