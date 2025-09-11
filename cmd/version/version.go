package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    
    "github.com/AlexBurnes/version-go/pkg/version"
)

// checkVersion validates a version string using the library
func checkVersion(versionStr string) error {
    return version.Validate(versionStr)
}

// getVersionType returns the type of a version string using the library
func getVersionType(versionStr string) (string, error) {
    versionType, err := version.GetType(versionStr)
    if err != nil {
        return "", err
    }
    return versionType.String(), nil
}

// getBuildType returns the CMake build type for a version using the library
func getBuildType(versionStr string) (string, error) {
    return version.GetBuildType(versionStr)
}

// sortVersions sorts version strings from stdin using the library
func sortVersions() (string, error) {
    scanner := bufio.NewScanner(os.Stdin)
    var versions []string
    
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            // Split line by whitespace to handle multiple versions per line
            parts := strings.Fields(line)
            versions = append(versions, parts...)
        }
    }
    
    if err := scanner.Err(); err != nil {
        return "", fmt.Errorf("error reading input: %v", err)
    }
    
    if len(versions) == 0 {
        return "", nil
    }
    
    // Sort versions using the library
    sortedVersions, err := version.Sort(versions)
    if err != nil {
        return "", err
    }
    
    return strings.Join(sortedVersions, "\n"), nil
}