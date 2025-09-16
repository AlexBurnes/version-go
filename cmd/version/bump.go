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
	if len(args) > 0 {
		// Version provided as argument
		return args[0], nil
	}

	// No version provided, use current git version
	version, err := getVersion()
	if err != nil {
		return "", fmt.Errorf("no version specified and failed to get current version: %v", err)
	}

	printDebug("Using current git version: %s", version)
	return version, nil
}

// getBumpType gets the bump type (from argument or default to smart)
func getBumpType(args []string) string {
	if len(args) > 1 {
		// Bump type provided as argument
		return args[1]
	}

	// Default to smart bump
	printDebug("No bump type specified, using smart bump")
	return "smart"
}

// validateBumpArgs validates the bump command arguments
func validateBumpArgs(args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("too many arguments - usage: bump [version] [type]")
	}

	// If version is provided, validate it
	if len(args) > 0 {
		if err := version.Validate(args[0]); err != nil {
			return fmt.Errorf("invalid version '%s': %v", args[0], err)
		}
	}

	// If bump type is provided, validate it
	if len(args) > 1 {
		_, err := version.ParseBumpType(args[1])
		if err != nil {
			return fmt.Errorf("invalid bump type '%s': %v", args[1], err)
		}
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