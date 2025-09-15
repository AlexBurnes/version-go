# Version Library Documentation

The `github.com/burnes/go-version/pkg/version` package provides a comprehensive library for parsing, validating, and sorting semantic versions with extended format support.

## Overview

This library extends standard SemVer 2.0 with support for:
- **Release versions**: `1.2.3`, `v1.2.3`
- **Prerelease versions**: `1.2.3-alpha`, `1.2.3~beta.1`, `1.2.3-rc.1`
- **Postrelease versions**: `1.2.3.fix`, `1.2.3.post.1`
- **Intermediate versions**: `1.2.3_feature`, `1.2.3_exp.1`

## Installation

```bash
go get github.com/AlexBurnes/version-go/pkg/version
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/AlexBurnes/version-go/pkg/version"
)

func main() {
    // Parse a version
    v, err := version.Parse("1.2.3-alpha.1")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Version: %s\n", v.String())
    fmt.Printf("Type: %s\n", v.Type.String())
    fmt.Printf("Build Type: %s\n", v.Type.BuildType())
    
    // Validate a version
    if version.IsValid("1.2.3") {
        fmt.Println("Valid version!")
    }
    
    // Sort versions
    versions := []string{"2.0.0", "1.2.3", "1.2.3-alpha"}
    sorted, err := version.Sort(versions)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Sorted: %v\n", sorted)
}
```

## API Reference

### Types

#### `Version`
Represents a parsed version with all its components.

```go
type Version struct {
    Major       int    // Major version number
    Minor       int    // Minor version number
    Patch       int    // Patch version number
    Type        Type   // Version type
    Prerelease  string // Prerelease identifier
    Postrelease string // Postrelease identifier
    Intermediate string // Intermediate identifier
    Original    string // Original version string
}
```

#### `Type`
Represents the type of a version.

```go
type Type int

const (
    TypeRelease Type = iota
    TypePrerelease
    TypePostrelease
    TypeIntermediate
    TypeInvalid
)
```

### Core Functions

#### `Parse(versionStr string) (*Version, error)`
Parses a version string and returns a `Version` struct.

```go
v, err := version.Parse("1.2.3-alpha.1")
if err != nil {
    log.Fatal(err)
}
```

#### `Validate(versionStr string) error`
Validates a version string without parsing it.

```go
err := version.Validate("1.2.3")
if err != nil {
    fmt.Printf("Invalid version: %v\n", err)
}
```

#### `IsValid(versionStr string) bool`
Checks if a version string is valid.

```go
if version.IsValid("1.2.3") {
    fmt.Println("Valid!")
}
```

#### `GetType(versionStr string) (Type, error)`
Returns the type of a version string.

```go
versionType, err := version.GetType("1.2.3-alpha")
if err != nil {
    log.Fatal(err)
}
fmt.Println(versionType.String()) // "prerelease"
```

#### `GetBuildType(versionStr string) (string, error)`
Returns the CMake build type for a version.

```go
buildType, err := version.GetBuildType("1.2.3")
if err != nil {
    log.Fatal(err)
}
fmt.Println(buildType) // "Release"
```

#### `Compare(a, b *Version) int`
Compares two versions for sorting. Returns -1 if a < b, 0 if a == b, 1 if a > b.

```go
v1, _ := version.Parse("1.2.3")
v2, _ := version.Parse("1.2.4")
result := version.Compare(v1, v2) // -1
```

#### `Sort(versions []string) ([]string, error)`
Sorts a list of version strings according to precedence rules.

```go
versions := []string{"2.0.0", "1.2.3", "1.2.3-alpha"}
sorted, err := version.Sort(versions)
if err != nil {
    log.Fatal(err)
}
// Result: ["1.2.3", "1.2.3-alpha", "2.0.0"]
```

#### `ConvertGitTag(tag string) string`
Converts git tag format from `x.y.z-(remainder)` to `x.y.z~(remainder)`.

```go
converted := version.ConvertGitTag("1.2.3-alpha")
fmt.Println(converted) // "1.2.3~alpha"
```

### Type Methods

#### `Type.String() string`
Returns the string representation of the version type.

```go
versionType := version.TypePrerelease
fmt.Println(versionType.String()) // "prerelease"
```

#### `Type.BuildType() string`
Returns the CMake build type for the version type.

```go
versionType := version.TypeRelease
fmt.Println(versionType.BuildType()) // "Release"
```

#### `Version.String() string`
Returns the original version string.

```go
v, _ := version.Parse("1.2.3")
fmt.Println(v.String()) // "1.2.3"
```

## Version Format Specification

### Release Versions
Standard semantic versions without additional identifiers.
- Format: `v?[0-9]+\.[0-9]+\.[0-9]+`
- Examples: `1.2.3`, `v1.2.3`, `0.0.0`

### Prerelease Versions
Versions with prerelease identifiers using `~` delimiter.
- Format: `v?[0-9]+\.[0-9]+\.[0-9]+~(alpha|beta|rc|pre)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*`
- Examples: `1.2.3~alpha`, `1.2.3~beta.1`, `1.2.3~rc.1`

### Postrelease Versions
Versions with postrelease identifiers using `.` delimiter.
- Format: `v?[0-9]+\.[0-9]+\.[0-9]+\.(fix|next|post)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*`
- Examples: `1.2.3.fix`, `1.2.3.post.1`, `1.2.3.next.1`

### Intermediate Versions
Versions with intermediate identifiers using `_` delimiter.
- Format: `v?[0-9]+\.[0-9]+\.[0-9]+_[a-zA-Z]+(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*`
- Examples: `1.2.3_feature`, `1.2.3_exp.1`, `1.2.3_dev.1`

## Sorting Rules

Versions are sorted according to the following precedence:

1. **Core version** (major.minor.patch) - numerical comparison
2. **Version type** - release < prerelease < postrelease < intermediate
3. **Type-specific identifiers** - alphanumeric comparison with numeric precedence

### Examples

```go
versions := []string{
    "2.0.0",           // Release
    "1.2.3",           // Release
    "1.2.3-alpha",     // Prerelease
    "1.2.3-beta.1",    // Prerelease
    "1.2.3-rc.1",      // Prerelease
    "1.2.3.fix",       // Postrelease
    "1.2.3_feature",   // Intermediate
    "1.2.4",           // Release
}

sorted, _ := version.Sort(versions)
// Result: ["1.2.3", "1.2.3-alpha", "1.2.3-beta.1", "1.2.3-rc.1", 
//          "1.2.3.fix", "1.2.3_feature", "1.2.4", "2.0.0"]
```

## Error Handling

All functions that can fail return an error as the last return value. Common errors include:

- `invalid version format: <version>` - Version string doesn't match any supported format
- Invalid input parameters

## Performance Considerations

- The library is optimized for performance with large version lists (10k+ versions)
- Parsing uses compiled regex patterns for efficiency
- Sorting uses Go's built-in sort algorithm with custom comparison function
- Memory usage is minimal with no unnecessary allocations

## Thread Safety

The library is thread-safe and can be used concurrently from multiple goroutines.

## Examples

See the `examples/basic/main.go` file for a comprehensive example demonstrating all library features.

## License

This library is part of the go-version project and is licensed under the Apache License, Version 2.0.