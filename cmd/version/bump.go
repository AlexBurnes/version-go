package main

import (
	"fmt"

	"github.com/AlexBurnes/version-go/pkg/version"
)

// bumpVersion bumps a version according to the specified bump type
func bumpVersion(versionStr string, bumpTypeStr string) (string, error) {
	// Parse bump type
	bumpType, err := version.ParseBumpType(bumpTypeStr)
	if err != nil {
		return "", fmt.Errorf("invalid bump type '%s': %v", bumpTypeStr, err)
	}

	// Perform the bump operation
	result, err := version.Bump(versionStr, bumpType)
	if err != nil {
		return "", fmt.Errorf("failed to bump version '%s': %v", versionStr, err)
	}

	// Print debug information
	if debugFlag {
		printDebug("Bump operation: %s -> %s", result.OriginalVersion, result.BumpedVersion)
		printDebug("Applied rule: %s", result.AppliedRule)
	}

	return result.BumpedVersion, nil
}

// getBumpVersion gets the version to bump (from argument or current git version)
func getBumpVersion(args []string) (string, error) {
	if len(args) == 0 {
		// No arguments provided, use current git version
		version, err := getVersion()
		if err != nil {
			return "", fmt.Errorf("no version specified and failed to get current version: %v", err)
		}
		printDebug("Using current git version: %s", version)
		return version, nil
	}

	if len(args) == 1 {
		// One argument provided - check if it's a version or bump type
		arg := args[0]
		
		// First check if it's a valid version
		if err := version.Validate(arg); err == nil {
			// It's a valid version
			printDebug("Using provided version: %s", arg)
			return arg, nil
		}
		
		// Check if it's a valid bump type
		if _, err := version.ParseBumpType(arg); err == nil {
			// It's a valid bump type, so no version provided - use current git version
			version, err := getVersion()
			if err != nil {
				return "", fmt.Errorf("no version specified and failed to get current version: %v", err)
			}
			printDebug("Using current git version with bump type: %s", arg)
			return version, nil
		}
		
		// Neither a valid version nor bump type
		return "", fmt.Errorf("invalid argument '%s': must be a valid version or bump type", arg)
	}

	// Two or more arguments - first should be version
	versionStr := args[0]
	if err := version.Validate(versionStr); err != nil {
		return "", fmt.Errorf("invalid version '%s': %v", versionStr, err)
	}
	
	printDebug("Using provided version: %s", versionStr)
	return versionStr, nil
}

// getBumpType gets the bump type (from argument or default to smart)
func getBumpType(args []string) string {
	if len(args) == 0 {
		// No arguments provided, default to smart bump
		printDebug("No bump type specified, using smart bump")
		return "smart"
	}

	if len(args) == 1 {
		// One argument provided - check if it's a bump type
		arg := args[0]
		
		// Check if it's a valid bump type
		if _, err := version.ParseBumpType(arg); err == nil {
			// It's a valid bump type
			printDebug("Using provided bump type: %s", arg)
			return arg
		}
		
		// It's a version, default to smart bump
		printDebug("Version provided, using smart bump")
		return "smart"
	}

	// Two or more arguments - second should be bump type
	printDebug("Using provided bump type: %s", args[1])
	return args[1]
}

// validateBumpArgs validates the bump command arguments
func validateBumpArgs(args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("too many arguments - usage: bump [version] [type]")
	}

	if len(args) == 0 {
		// No arguments - both version and bump type will use defaults
		return nil
	}

	if len(args) == 1 {
		// One argument - could be version or bump type
		arg := args[0]
		
		// Check if it's a valid version
		if err := version.Validate(arg); err == nil {
			// It's a valid version
			return nil
		}
		
		// Check if it's a valid bump type
		if _, err := version.ParseBumpType(arg); err == nil {
			// It's a valid bump type
			return nil
		}
		
		// Neither valid version nor bump type
		return fmt.Errorf("invalid argument '%s': must be a valid version or bump type", arg)
	}

	// Two arguments - first should be version, second should be bump type
	if err := version.Validate(args[0]); err != nil {
		return fmt.Errorf("invalid version '%s': %v", args[0], err)
	}

	if _, err := version.ParseBumpType(args[1]); err != nil {
		return fmt.Errorf("invalid bump type '%s': %v", args[1], err)
	}

	return nil
}

// printBumpHelp prints help information for the bump command
func printBumpHelp() {
	fmt.Printf(`Bump command usage:
    version bump [version] [type]

Arguments:
    version    Version to bump (optional, uses current git version if not specified)
    type       Bump type (optional, defaults to 'smart')

Bump types:
    major      Increment major version and reset minor/patch (e.g., 1.2.3 -> 2.0.0)
    minor      Increment minor version and reset patch (e.g., 1.2.3 -> 1.3.0)
    patch      Increment patch version (e.g., 1.2.3 -> 1.2.4)
    pre        Convert to prerelease with pre.1 or increment prerelease identifier
    alpha      Convert to prerelease with alpha.1 or increment prerelease identifier
    beta       Convert to prerelease with beta.1 or increment prerelease identifier
    rc         Convert to prerelease with rc.1 or increment prerelease identifier
    fix        Convert to postrelease with fix.1 or increment postrelease identifier
    next       Convert to postrelease with next.1 or increment postrelease identifier
    post       Convert to postrelease with post.1 or increment postrelease identifier
    feat       Convert to intermediate with feat.1 or increment intermediate identifier
    smart      Intelligent bump based on current version type (default)

Examples:
    version bump                    # Smart bump current git version
    version bump 1.2.3             # Smart bump version 1.2.3
    version bump 1.2.3 major       # Major bump version 1.2.3
    version bump 1.2.3 alpha       # Convert to prerelease with alpha.1
    version bump 1.2.3~alpha.1     # Smart bump prerelease version
    version bump 1.2.3 fix         # Convert to postrelease with fix.1
    version bump 1.2.3 feat        # Convert to intermediate with feat.1

Build Script Usage:
    # Bump current git version and capture result (already silent by default)
    NEW_VERSION=$(scripts/version bump patch)
    echo "New version: $NEW_VERSION"
    
    # Smart bump current git version
    BUMPED=$(scripts/version bump)
    echo "Bumped version: $BUMPED"
    
    # Bump specific version with specific type
    NEXT_VERSION=$(scripts/version bump 1.2.3 minor)
    echo "Next version: $NEXT_VERSION"
    
    # Use in CMake or other build systems
    set(APP_VERSION "${NEW_VERSION}")
    message(STATUS "Building version: ${APP_VERSION}")
    
    # Use in shell scripts for version management
    CURRENT_VERSION=$(scripts/version version)
    NEW_VERSION=$(scripts/version bump patch)
    echo "Current: $CURRENT_VERSION -> New: $NEW_VERSION"
    
    # Direct usage with version utility
    scripts/version bump patch      # Basic patch bump
    scripts/version bump alpha      # Convert to prerelease
    scripts/version bump smart      # Intelligent bump

Smart bump behavior:
    - Release versions: increment patch
    - Prerelease versions: increment prerelease identifier
    - Postrelease versions: increment postrelease identifier
    - Intermediate versions: increment intermediate identifier
`)
}

// getBumpCommandHelp returns help text for the bump command
func getBumpCommandHelp() string {
	return `bump [version] [type]    bump version with specified type (smart, major, minor, patch, pre, alpha, beta, rc, fix, next, post, feat)`
}