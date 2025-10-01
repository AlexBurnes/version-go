// Package version provides semantic version parsing, validation, and ordering functionality.
// It supports extended version formats beyond standard SemVer 2.0, including prerelease,
// postrelease, and intermediate identifiers.
package version

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Type represents the type of a version
type Type int

const (
	TypePrerelease Type = iota
	TypeRelease
	TypePostrelease
	TypeIntermediate
	TypeInvalid
)

func (t Type) String() string {
	switch t {
	case TypeRelease:
		return "release"
	case TypePrerelease:
		return "prerelease"
	case TypePostrelease:
		return "postrelease"
	case TypeIntermediate:
		return "intermediate"
	default:
		return "invalid"
	}
}

// BuildType returns the CMake build type for this version type
func (t Type) BuildType() string {
	if t == TypeRelease {
		return "Release"
	}
	return "Debug"
}

// Version represents a parsed version with all its components
type Version struct {
	Major       int    // Major version number
	Minor       int    // Minor version number
	Patch       int    // Patch version number
	Type        Type   // Version type (release, prerelease, etc.)
	Prerelease  string // Prerelease identifier (e.g., "~alpha.1")
	Postrelease string // Postrelease identifier (e.g., ".fix.1")
	Intermediate string // Intermediate identifier (e.g., "_feature.1")
	Original    string // Original version string
}

// Regex patterns for version parsing
var (
	versionRelease      = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	versionPrerelease   = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\~(alpha|beta|rc|pre)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
	versionPostrelease  = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\.(fix|next|post)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
	versionIntermediate = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\_([a-zA-Z]+)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
)

// ConvertGitTag converts git tag format from x.y.z-(remainder) to x.y.z~(remainder)
// This is useful for handling git tags that use - delimiter instead of ~ for prerelease versions
func ConvertGitTag(tag string) string {
	// Check if this looks like a prerelease tag with - delimiter
	if strings.Contains(tag, "-") && !strings.Contains(tag, "~") {
		// Find the first occurrence of - after the version number
		// Pattern: v?x.y.z-remainder
		re := regexp.MustCompile(`^(v?[0-9]+\.[0-9]+\.[0-9]+)-(.+)$`)
		if matches := re.FindStringSubmatch(tag); matches != nil {
			return matches[1] + "~" + matches[2]
		}
	}
	return tag
}

// Parse parses a version string and returns a Version struct
// It supports release, prerelease, postrelease, and intermediate version formats
func Parse(versionStr string) (*Version, error) {
	versionStr = strings.TrimSpace(versionStr)
	
	// Convert git tag format if needed
	versionStr = ConvertGitTag(versionStr)
	
	// Try release version first
	if matches := versionRelease.FindStringSubmatch(versionStr); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		
		return &Version{
			Major:    major,
			Minor:    minor,
			Patch:    patch,
			Type:     TypeRelease,
			Original: versionStr,
		}, nil
	}
	
	// Try prerelease version
	if matches := versionPrerelease.FindStringSubmatch(versionStr); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		prerelease := "~" + matches[4] + matches[5] + matches[6]
		
		return &Version{
			Major:      major,
			Minor:      minor,
			Patch:      patch,
			Type:       TypePrerelease,
			Prerelease: prerelease,
			Original:   versionStr,
		}, nil
	}
	
	// Try postrelease version
	if matches := versionPostrelease.FindStringSubmatch(versionStr); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		postrelease := "." + matches[4] + matches[5]
		
		return &Version{
			Major:       major,
			Minor:       minor,
			Patch:       patch,
			Type:        TypePostrelease,
			Postrelease: postrelease,
			Original:    versionStr,
		}, nil
	}
	
	// Try intermediate version
	if matches := versionIntermediate.FindStringSubmatch(versionStr); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		intermediate := "_" + matches[4] + matches[5]
		
		return &Version{
			Major:        major,
			Minor:        minor,
			Patch:        patch,
			Type:         TypeIntermediate,
			Intermediate: intermediate,
			Original:     versionStr,
		}, nil
	}
	
	return nil, fmt.Errorf("invalid version format: %s", versionStr)
}

// Validate checks if a version string is valid
func Validate(versionStr string) error {
	_, err := Parse(versionStr)
	return err
}

// GetType returns the type of a version string
func GetType(versionStr string) (Type, error) {
	version, err := Parse(versionStr)
	if err != nil {
		return TypeInvalid, err
	}
	return version.Type, nil
}

// GetBuildType returns the CMake build type for a version
func GetBuildType(versionStr string) (string, error) {
	version, err := Parse(versionStr)
	if err != nil {
		return "", err
	}
	return version.Type.BuildType(), nil
}

// Compare compares two versions for sorting
// Returns -1 if a < b, 0 if a == b, 1 if a > b
func Compare(a, b *Version) int {
	// First compare major.minor.patch
	if a.Major != b.Major {
		return a.Major - b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor - b.Minor
	}
	if a.Patch != b.Patch {
		return a.Patch - b.Patch
	}
	
	// For same core version, compare by type precedence
	if a.Type != b.Type {
		return int(a.Type) - int(b.Type)
	}
	
	// For same type, compare by type-specific identifiers
	switch a.Type {
	case TypePrerelease:
		return compareIdentifiers(a.Prerelease, b.Prerelease)
	case TypePostrelease:
		return compareIdentifiers(a.Postrelease, b.Postrelease)
	case TypeIntermediate:
		return compareIdentifiers(a.Intermediate, b.Intermediate)
	default:
		return 0
	}
}

// compareIdentifiers compares alphanumeric identifiers
func compareIdentifiers(a, b string) int {
	// Split identifiers by dots and underscores
	aParts := splitIdentifier(a)
	bParts := splitIdentifier(b)
	
	// Compare each part
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	
	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}
		
		result := compareIdentifierPart(aPart, bPart)
		if result != 0 {
			return result
		}
	}
	
	return 0
}

// splitIdentifier splits an identifier by dots and underscores
func splitIdentifier(identifier string) []string {
	// Replace underscores with dots for consistent splitting
	normalized := strings.ReplaceAll(identifier, "_", ".")
	return strings.Split(normalized, ".")
}

// compareIdentifierPart compares two identifier parts
func compareIdentifierPart(a, b string) int {
	// Handle empty parts
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}
	
	// Try to parse as numbers
	aNum, aErr := strconv.Atoi(a)
	bNum, bErr := strconv.Atoi(b)
	
	// If both are numbers, compare numerically
	if aErr == nil && bErr == nil {
		return aNum - bNum
	}
	
	// If one is a number and the other isn't, number comes first
	if aErr == nil && bErr != nil {
		return -1
	}
	if aErr != nil && bErr == nil {
		return 1
	}
	
	// Both are strings, compare lexically
	return strings.Compare(a, b)
}

// Sort sorts a list of version strings according to precedence rules
// Returns the sorted versions as a slice of strings
func Sort(versions []string) ([]string, error) {
	if len(versions) == 0 {
		return []string{}, nil
	}
	
	// Parse all versions
	var parsedVersions []*Version
	for _, v := range versions {
		parsed, err := Parse(v)
		if err != nil {
			return nil, fmt.Errorf("invalid version '%s': %v", v, err)
		}
		parsedVersions = append(parsedVersions, parsed)
	}
	
	// Sort versions
	sort.Slice(parsedVersions, func(i, j int) bool {
		return Compare(parsedVersions[i], parsedVersions[j]) < 0
	})
	
	// Build result
	var result []string
	for _, v := range parsedVersions {
		result = append(result, v.Original)
	}
	
	return result, nil
}

// IsValid checks if a version string is valid
func IsValid(versionStr string) bool {
	return Validate(versionStr) == nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.Original
}