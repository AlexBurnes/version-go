package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// BumpType represents the type of version bump to perform
type BumpType int

const (
	BumpMajor BumpType = iota
	BumpMinor
	BumpPatch
	BumpPre
	BumpAlpha
	BumpBeta
	BumpRc
	BumpFix
	BumpNext
	BumpPost
	BumpFeat
	BumpSmart // Smart increment based on current version state
)

func (bt BumpType) String() string {
	switch bt {
	case BumpMajor:
		return "major"
	case BumpMinor:
		return "minor"
	case BumpPatch:
		return "patch"
	case BumpPre:
		return "pre"
	case BumpAlpha:
		return "alpha"
	case BumpBeta:
		return "beta"
	case BumpRc:
		return "rc"
	case BumpFix:
		return "fix"
	case BumpNext:
		return "next"
	case BumpPost:
		return "post"
	case BumpFeat:
		return "feat"
	case BumpSmart:
		return "smart"
	default:
		return "unknown"
	}
}

// BumpResult represents the result of a version bump operation
type BumpResult struct {
	OriginalVersion string
	BumpedVersion   string
	BumpType        BumpType
	AppliedRule     string
}

// Bump bumps a version according to the specified bump type
func Bump(versionStr string, bumpType BumpType) (*BumpResult, error) {
	version, err := Parse(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version '%s': %v", versionStr, err)
	}

	var bumpedVersion *Version
	var appliedRule string

	switch bumpType {
	case BumpSmart:
		bumpedVersion, appliedRule = bumpSmart(version)
	case BumpMajor:
		bumpedVersion, appliedRule = bumpMajor(version)
	case BumpMinor:
		bumpedVersion, appliedRule = bumpMinor(version)
	case BumpPatch:
		bumpedVersion, appliedRule = bumpPatch(version)
	case BumpPre:
		bumpedVersion, appliedRule = bumpPrerelease(version, "pre")
	case BumpAlpha:
		bumpedVersion, appliedRule = bumpPrerelease(version, "alpha")
	case BumpBeta:
		bumpedVersion, appliedRule = bumpPrerelease(version, "beta")
	case BumpRc:
		bumpedVersion, appliedRule = bumpPrerelease(version, "rc")
	case BumpFix:
		bumpedVersion, appliedRule = bumpPostrelease(version, "fix")
	case BumpNext:
		bumpedVersion, appliedRule = bumpPostrelease(version, "next")
	case BumpPost:
		bumpedVersion, appliedRule = bumpPostrelease(version, "post")
	case BumpFeat:
		bumpedVersion, appliedRule = bumpIntermediate(version, "feat")
	default:
		return nil, fmt.Errorf("unknown bump type: %v", bumpType)
	}

	return &BumpResult{
		OriginalVersion: version.Original,
		BumpedVersion:   bumpedVersion.String(),
		BumpType:        bumpType,
		AppliedRule:     appliedRule,
	}, nil
}

// bumpSmart performs intelligent version bumping based on current version state
func bumpSmart(version *Version) (*Version, string) {
	switch version.Type {
	case TypeRelease:
		// For release versions, increment patch
		bumped, _ := bumpPatch(version)
		return bumped, "increment patch for release version"
	case TypePrerelease:
		// For prerelease versions, increment the prerelease identifier
		bumped, _ := incrementPrerelease(version)
		return bumped, "increment prerelease identifier"
	case TypePostrelease:
		// For postrelease versions, increment the postrelease identifier
		bumped, _ := incrementPostrelease(version)
		return bumped, "increment postrelease identifier"
	case TypeIntermediate:
		// For intermediate versions, increment the intermediate identifier
		bumped, _ := incrementIntermediate(version)
		return bumped, "increment intermediate identifier"
	default:
		// Fallback to patch bump for invalid types
		bumped, _ := bumpPatch(version)
		return bumped, "fallback to patch bump for invalid version type"
	}
}

// bumpMajor increments the major version and resets minor and patch
func bumpMajor(version *Version) (*Version, string) {
	return &Version{
		Major:    version.Major + 1,
		Minor:    0,
		Patch:    0,
		Type:     TypeRelease,
		Original: fmt.Sprintf("%d.%d.%d", version.Major+1, 0, 0),
	}, "increment major version and reset minor/patch"
}

// bumpMinor increments the minor version and resets patch
func bumpMinor(version *Version) (*Version, string) {
	return &Version{
		Major:    version.Major,
		Minor:    version.Minor + 1,
		Patch:    0,
		Type:     TypeRelease,
		Original: fmt.Sprintf("%d.%d.%d", version.Major, version.Minor+1, 0),
	}, "increment minor version and reset patch"
}

// bumpPatch increments the patch version
func bumpPatch(version *Version) (*Version, string) {
	return &Version{
		Major:    version.Major,
		Minor:    version.Minor,
		Patch:    version.Patch + 1,
		Type:     TypeRelease,
		Original: fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch+1),
	}, "increment patch version"
}

// bumpPrerelease increments the prerelease version
func bumpPrerelease(version *Version, identifier string) (*Version, string) {
	if version.Type == TypePrerelease {
		// Increment existing prerelease
		bumped, _ := incrementPrerelease(version)
		return bumped, "increment existing prerelease identifier"
	}
	// Convert to prerelease
	return &Version{
		Major:      version.Major,
		Minor:      version.Minor,
		Patch:      version.Patch,
		Type:       TypePrerelease,
		Prerelease: "~" + identifier + ".1",
		Original:   fmt.Sprintf("%d.%d.%d~%s.1", version.Major, version.Minor, version.Patch, identifier),
	}, "convert to prerelease with " + identifier + ".1"
}

// bumpPostrelease increments the postrelease version
func bumpPostrelease(version *Version, identifier string) (*Version, string) {
	if version.Type == TypePostrelease {
		// Increment existing postrelease
		bumped, _ := incrementPostrelease(version)
		return bumped, "increment existing postrelease identifier"
	}
	// Convert to postrelease
	return &Version{
		Major:       version.Major,
		Minor:       version.Minor,
		Patch:       version.Patch,
		Type:        TypePostrelease,
		Postrelease: "." + identifier + ".1",
		Original:    fmt.Sprintf("%d.%d.%d.%s.1", version.Major, version.Minor, version.Patch, identifier),
	}, "convert to postrelease with " + identifier + ".1"
}

// bumpIntermediate increments the intermediate version
func bumpIntermediate(version *Version, identifier string) (*Version, string) {
	if version.Type == TypeIntermediate {
		// Increment existing intermediate
		bumped, _ := incrementIntermediate(version)
		return bumped, "increment existing intermediate identifier"
	}
	// Convert to intermediate
	return &Version{
		Major:        version.Major,
		Minor:        version.Minor,
		Patch:        version.Patch,
		Type:         TypeIntermediate,
		Intermediate: "_" + identifier + ".1",
		Original:     fmt.Sprintf("%d.%d.%d_%s.1", version.Major, version.Minor, version.Patch, identifier),
	}, "convert to intermediate with " + identifier + ".1"
}

// incrementPrerelease increments an existing prerelease version
func incrementPrerelease(version *Version) (*Version, string) {
	// Parse the prerelease identifier to find the last numeric part
	identifier := version.Prerelease[1:] // Remove the ~ prefix
	
	// Find the last numeric part and increment it
	// Look for patterns like .1, .1.2, _1, _1.2, etc.
	re := regexp.MustCompile(`([._])([0-9]+)([^0-9]*)$`)
	matches := re.FindStringSubmatch(identifier)
	
	if matches != nil {
		// Found a numeric suffix, increment it
		num, _ := strconv.Atoi(matches[2])
		newNum := strconv.Itoa(num + 1)
		newIdentifier := re.ReplaceAllString(identifier, matches[1]+newNum+matches[3])
		newPrerelease := "~" + newIdentifier
		
		return &Version{
			Major:      version.Major,
			Minor:      version.Minor,
			Patch:      version.Patch,
			Type:       TypePrerelease,
			Prerelease: newPrerelease,
			Original:   fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newPrerelease),
		}, "increment prerelease numeric identifier"
	}
	
	// No numeric suffix found, add .1
	newPrerelease := "~" + identifier + ".1"
	return &Version{
		Major:      version.Major,
		Minor:      version.Minor,
		Patch:      version.Patch,
		Type:       TypePrerelease,
		Prerelease: newPrerelease,
		Original:   fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newPrerelease),
	}, "add prerelease numeric identifier"
}

// incrementPostrelease increments an existing postrelease version
func incrementPostrelease(version *Version) (*Version, string) {
	// Parse the postrelease identifier to find the last numeric part
	identifier := version.Postrelease[1:] // Remove the . prefix
	
	// Find the last numeric part and increment it
	// Look for patterns like .1, .1.2, _1, _1.2, etc.
	re := regexp.MustCompile(`([._])([0-9]+)([^0-9]*)$`)
	matches := re.FindStringSubmatch(identifier)
	
	if matches != nil {
		// Found a numeric suffix, increment it
		num, _ := strconv.Atoi(matches[2])
		newNum := strconv.Itoa(num + 1)
		newIdentifier := re.ReplaceAllString(identifier, matches[1]+newNum+matches[3])
		newPostrelease := "." + newIdentifier
		
		return &Version{
			Major:       version.Major,
			Minor:       version.Minor,
			Patch:       version.Patch,
			Type:        TypePostrelease,
			Postrelease: newPostrelease,
			Original:    fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newPostrelease),
		}, "increment postrelease numeric identifier"
	}
	
	// No numeric suffix found, add .1
	newPostrelease := "." + identifier + ".1"
	return &Version{
		Major:       version.Major,
		Minor:       version.Minor,
		Patch:       version.Patch,
		Type:        TypePostrelease,
		Postrelease: newPostrelease,
		Original:    fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newPostrelease),
	}, "add postrelease numeric identifier"
}

// incrementIntermediate increments an existing intermediate version
func incrementIntermediate(version *Version) (*Version, string) {
	// Parse the intermediate identifier to find the last numeric part
	identifier := version.Intermediate[1:] // Remove the _ prefix
	
	// Find the last numeric part and increment it
	// Look for patterns like .1, .1.2, _1, _1.2, etc.
	re := regexp.MustCompile(`([._])([0-9]+)([^0-9]*)$`)
	matches := re.FindStringSubmatch(identifier)
	
	if matches != nil {
		// Found a numeric suffix, increment it
		num, _ := strconv.Atoi(matches[2])
		newNum := strconv.Itoa(num + 1)
		newIdentifier := re.ReplaceAllString(identifier, matches[1]+newNum+matches[3])
		newIntermediate := "_" + newIdentifier
		
		return &Version{
			Major:        version.Major,
			Minor:        version.Minor,
			Patch:        version.Patch,
			Type:         TypeIntermediate,
			Intermediate: newIntermediate,
			Original:     fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newIntermediate),
		}, "increment intermediate numeric identifier"
	}
	
	// No numeric suffix found, add .1
	newIntermediate := "_" + identifier + ".1"
	return &Version{
		Major:        version.Major,
		Minor:        version.Minor,
		Patch:        version.Patch,
		Type:         TypeIntermediate,
		Intermediate: newIntermediate,
		Original:     fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Patch, newIntermediate),
	}, "add intermediate numeric identifier"
}

// ParseBumpType parses a bump type string
func ParseBumpType(bumpTypeStr string) (BumpType, error) {
	switch strings.ToLower(bumpTypeStr) {
	case "major":
		return BumpMajor, nil
	case "minor":
		return BumpMinor, nil
	case "patch":
		return BumpPatch, nil
	case "pre":
		return BumpPre, nil
	case "alpha":
		return BumpAlpha, nil
	case "beta":
		return BumpBeta, nil
	case "rc":
		return BumpRc, nil
	case "fix":
		return BumpFix, nil
	case "next":
		return BumpNext, nil
	case "post":
		return BumpPost, nil
	case "feat":
		return BumpFeat, nil
	case "smart":
		return BumpSmart, nil
	default:
		return BumpSmart, fmt.Errorf("unknown bump type: %s", bumpTypeStr)
	}
}