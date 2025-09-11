package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strconv"
    "strings"
)

// VersionType represents the type of a version
type VersionType int

const (
    VersionTypeRelease VersionType = iota
    VersionTypePrerelease
    VersionTypePostrelease
    VersionTypeIntermediate
    VersionTypeInvalid
)

func (vt VersionType) String() string {
    switch vt {
    case VersionTypeRelease:
        return "release"
    case VersionTypePrerelease:
        return "prerelease"
    case VersionTypePostrelease:
        return "postrelease"
    case VersionTypeIntermediate:
        return "intermediate"
    default:
        return "invalid"
    }
}

func (vt VersionType) BuildType() string {
    if vt == VersionTypeRelease {
        return "Release"
    }
    return "Debug"
}

// Version represents a parsed version
type Version struct {
    Major      int
    Minor      int
    Patch      int
    Type       VersionType
    Prerelease string
    Postrelease string
    Intermediate string
    Original   string
}

// Regex patterns for version parsing
var (
    versionRelease      = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
    versionPrerelease   = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\~(alpha|beta|rc|pre)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
    versionPostrelease  = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\.(fix|next|post)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
    versionIntermediate = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\_([a-zA-Z]+)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
)

// convertGitTag converts git tag format from x.y.z-(remainder) to x.y.z~(remainder)
func convertGitTag(tag string) string {
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

// ParseVersion parses a version string and returns a Version struct
func ParseVersion(versionStr string) (*Version, error) {
    versionStr = strings.TrimSpace(versionStr)
    
    // Convert git tag format if needed
    versionStr = convertGitTag(versionStr)
    
    // Try release version first
    if matches := versionRelease.FindStringSubmatch(versionStr); matches != nil {
        major, _ := strconv.Atoi(matches[1])
        minor, _ := strconv.Atoi(matches[2])
        patch, _ := strconv.Atoi(matches[3])
        
        return &Version{
            Major:    major,
            Minor:    minor,
            Patch:    patch,
            Type:     VersionTypeRelease,
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
            Type:       VersionTypePrerelease,
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
            Type:        VersionTypePostrelease,
            Postrelease: postrelease,
            Original:   versionStr,
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
            Type:         VersionTypeIntermediate,
            Intermediate: intermediate,
            Original:     versionStr,
        }, nil
    }
    
    return nil, fmt.Errorf("invalid version format: %s", versionStr)
}

// CheckVersion validates a version string
func checkVersion(versionStr string) error {
    _, err := ParseVersion(versionStr)
    return err
}

// GetVersionType returns the type of a version string
func getVersionType(versionStr string) (string, error) {
    version, err := ParseVersion(versionStr)
    if err != nil {
        return "", err
    }
    return version.Type.String(), nil
}

// GetBuildType returns the CMake build type for a version
func getBuildType(versionStr string) (string, error) {
    version, err := ParseVersion(versionStr)
    if err != nil {
        return "", err
    }
    return version.Type.BuildType(), nil
}

// CompareVersions compares two versions for sorting
func CompareVersions(a, b *Version) int {
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
    case VersionTypePrerelease:
        return compareIdentifiers(a.Prerelease, b.Prerelease)
    case VersionTypePostrelease:
        return compareIdentifiers(a.Postrelease, b.Postrelease)
    case VersionTypeIntermediate:
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

// SortVersions sorts a list of version strings
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
    
    // Parse all versions
    var parsedVersions []*Version
    for _, v := range versions {
        parsed, err := ParseVersion(v)
        if err != nil {
            return "", fmt.Errorf("invalid version '%s': %v", v, err)
        }
        parsedVersions = append(parsedVersions, parsed)
    }
    
    // Sort versions
    sort.Slice(parsedVersions, func(i, j int) bool {
        return CompareVersions(parsedVersions[i], parsedVersions[j]) < 0
    })
    
    // Build result
    var result []string
    for _, v := range parsedVersions {
        result = append(result, v.Original)
    }
    
    return strings.Join(result, "\n"), nil
}