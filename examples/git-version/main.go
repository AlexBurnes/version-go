package main

import (
    "fmt"
    "log"
    
    "github.com/AlexBurnes/version-go/pkg/version"
)

func main() {
    fmt.Println("=== Git Version Integration Example ===\n")
    
    // Get version from git repository
    gitVersion, err := version.GetVersion()
    if err != nil {
        // Handle different types of errors
        if version.IsGitNotFound(err) {
            fmt.Println("❌ Error: Git is not installed")
            fmt.Println("   Please install git to use this feature")
        } else if version.IsNotGitRepo(err) {
            fmt.Println("❌ Error: Not in a git repository")
            fmt.Println("   Please run this from within a git repository")
        } else if version.IsNoGitTags(err) {
            fmt.Println("❌ Error: No version tags found")
            fmt.Println("   Please create a version tag (e.g., v1.0.0) first")
        } else {
            log.Fatalf("❌ Unexpected error: %v", err)
        }
        return
    }
    
    fmt.Printf("✓ Git Version (without prefix): %s\n", gitVersion)
    
    // Get version with prefix
    gitVersionWithPrefix, err := version.GetVersionWithPrefix()
    if err != nil {
        log.Fatalf("Failed to get version with prefix: %v", err)
    }
    
    fmt.Printf("✓ Git Version (with prefix):    %s\n\n", gitVersionWithPrefix)
    
    // Parse the version to get detailed information
    v, err := version.Parse(gitVersion)
    if err != nil {
        log.Fatalf("Failed to parse version: %v", err)
    }
    
    // Display version details
    fmt.Println("Version Details:")
    fmt.Printf("  - Major:      %d\n", v.Major)
    fmt.Printf("  - Minor:      %d\n", v.Minor)
    fmt.Printf("  - Patch:      %d\n", v.Patch)
    fmt.Printf("  - Type:       %s\n", v.Type.String())
    fmt.Printf("  - Build Type: %s\n", v.Type.BuildType())
    
    if v.Prerelease != "" {
        fmt.Printf("  - Prerelease: %s\n", v.Prerelease)
    }
    if v.Postrelease != "" {
        fmt.Printf("  - Postrelease: %s\n", v.Postrelease)
    }
    if v.Intermediate != "" {
        fmt.Printf("  - Intermediate: %s\n", v.Intermediate)
    }
    
    fmt.Println("\n=== Example Complete ===")
}

