# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.1] - 2024-12-19

### Fixed
- **Conan Build Script**: Fixed critical issues in `buildtools/build-conan.sh` for local builds
  - Fixed version detection to use `git describe` instead of non-existent `./scripts/describe`
  - Corrected source directory references from `src/` to `cmd/version/` in CMakeLists.txt
  - Fixed Conan file path issues when running from `buildtools/` directory
  - Resolved CMake preset path issues for proper build configuration
  - Fixed build directory path from `../../bin` to `bin` for correct binary placement
  - Updated static build target to use correct working directory (`cmd/version/`)

### Changed
- **Build System**: Improved Conan integration for reliable local builds
  - All Conan commands now properly reference project root directory
  - CMake configuration works correctly from any directory
  - Cross-platform builds now generate binaries in correct project `bin/` directory

### Technical Details
- Conan build script now fully functional for local development
- Cross-platform builds working: Linux/amd64, macOS/amd64, macOS/arm64, Windows/amd64
- Static Linux binary builds successfully (2.1MB vs 3.1MB regular binary)
- All build targets accessible via `./build-conan.sh` commands
- Dependencies managed via Conan: Go 1.21.0, CMake 3.31.8

## [0.5.0] - 2024-12-19

### Added
- **Clean Project Structure**: Reorganized project to follow Go conventions with `cmd/` directory
- **Improved Organization**: CLI source code now in `cmd/version/` directory
- **Better Separation**: Clear separation between CLI code and library code
- **Standard Go Layout**: Follows standard Go project structure for better maintainability

### Changed
- **Project Structure**: Moved CLI files from root to `cmd/version/` directory
  - `main.go` → `cmd/version/main.go`
  - `git.go` → `cmd/version/git.go`
  - `version.go` → `cmd/version/version.go`
  - `version_test.go` → `cmd/version/version_test.go`
  - `integration_test.go` → `cmd/version/integration_test.go`
- **Build Commands**: Updated build commands to use `./cmd/version` path
- **Repository Organization**: Cleaner root directory with only configuration and documentation files

### Technical Details
- CLI executable: `go build -o bin/version ./cmd/version`
- Library package: `pkg/version/` (unchanged)
- Example code: `examples/basic/` (unchanged)
- All functionality remains identical, only project structure improved
- Maintains 100% backward compatibility with existing CLI interface

## [0.4.0] - 2024-12-19

### Added
- **Library Package**: Refactored core functionality into reusable `pkg/version` package
- **Library API**: Comprehensive public API with exported types and functions
- **Library Documentation**: Complete API reference and usage examples in `docs/LIBRARY.md`
- **Library Example**: Working example in `examples/basic/` demonstrating library usage
- **Repository Migration**: Moved from `github.com/burnes/go-version` to `github.com/AlexBurnes/version-go`

### Changed
- **CLI Architecture**: CLI now uses the library package internally while maintaining full compatibility
- **Import Paths**: Updated all import paths to use new repository URL
- **Documentation**: Updated README.md with library usage section and new repository URLs

### Technical Details
- Library package provides: `Parse()`, `Validate()`, `Sort()`, `Compare()`, `GetType()`, `GetBuildType()`
- Exported types: `Version` struct, `Type` enum
- Maintains 100% backward compatibility with existing CLI interface
- All tests passing for both CLI and library functionality

## [0.3.0] - 2024-12-19

### Added
- Initial implementation of version CLI utility
- Custom regex-based version parser supporting extended grammar
- Colored output system with terminal detection
- Complete CLI command set (project, module, version, release, full, check, check-greatest, type, build-type, sort)
- Git integration for version extraction
- Version sorting and comparison algorithms
- CMake build system integration
- Comprehensive test suite with 25.7% coverage
- Cross-platform support (Linux, Windows, macOS)

### Changed
- Replaced legacy bash scripts with Go implementation
- Updated build system to use CMake + Conan + GoReleaser

### Fixed
- Fixed version parsing for postrelease and intermediate formats
- Fixed comparison logic for proper version sorting
- Fixed CMake test targets to run from correct directory

## [0.3.0] - 2024-09-11

### Added
- **Core Version Parser**: Custom regex-based parser avoiding semver library complexity
  - Release versions: `1.2.3`, `v1.2.3`
  - Prerelease versions: `1.2.3-alpha`, `1.2.3~beta.1`, `1.2.3-rc.1`
  - Postrelease versions: `1.2.3.fix`, `1.2.3.post.1`, `1.2.3.next`
  - Intermediate versions: `1.2.3_feature`, `1.2.3_exp.1`

- **Colored Output System**: Terminal-friendly output based on bash script patterns
  - Error messages: Red with bold "ERROR" prefix
  - Success messages: Green with bold "SUCCESS" prefix
  - Debug messages: Yellow with bold "#DEBUG" prefix
  - Warning messages: Purple with "WARNING" prefix
  - Info messages: Green with "INFO" prefix
  - `--no-color` flag support for non-terminal output

- **CLI Commands**:
  - `project` - Print project name from git remote
  - `module` - Print module name from git remote
  - `version` - Print project version from git tags
  - `release` - Print project release number
  - `full` - Print full project name-version-release
  - `check [version]` - Validate version string (uses current git version if not specified)
  - `check-greatest [version]` - Check if version is greatest among all tags
  - `type [version]` - Print version type (release, prerelease, postrelease, intermediate)
  - `build-type [version]` - Print CMake build type (Release/Debug) based on version type
  - `sort` - Sort version strings from stdin

- **Version Sorting & Comparison**:
  - Correct precedence order: release → prerelease → postrelease → intermediate
  - Numeric identifiers compared numerically
  - Alphanumeric identifiers compared lexically
  - Handles mixed alphanumeric identifiers correctly

- **Git Integration**:
  - Reads git tags for version information
  - Handles both SSH and HTTPS remote URLs
  - Extracts project and module names from git remotes
  - Supports version validation against git tags

- **Build System**:
  - Go module structure in `src/` directory
  - CMake integration with cross-platform builds
  - Static binary builds for maximum compatibility
  - Test integration through CMake targets

- **Testing**:
  - Unit tests for version parsing and validation
  - Integration tests for CLI commands
  - Test coverage reporting
  - Race detection enabled

### Technical Details
- **Language**: Go 1.22+ (no CGO, static binaries)
- **Dependencies**: Standard library only (no external dependencies)
- **Platforms**: Linux (amd64/arm64), Windows (amd64/arm64), macOS (amd64/arm64)
- **Build System**: CMake + Conan + GoReleaser
- **Exit Codes**: POSIX compliant (0 for success, ≥1 for errors)

### Performance
- Efficient parsing and sorting for large version lists (10k+ versions)
- Minimal memory footprint with standard library only
- Fast startup time with static binaries

## [0.2.0] - Legacy

### Added
- Initial bash script implementation
- Basic version checking functionality
- Git integration for version extraction

### Changed
- Migrated from bash to Go implementation
- Enhanced version grammar support
- Improved error handling and user experience

## [0.1.0] - Legacy

### Added
- Initial project setup
- Basic version parsing
- Simple CLI interface

---

## Migration Guide

### From Bash Scripts

The Go implementation maintains compatibility with existing bash script interfaces:

```bash
# Old bash script
./version project

# New Go implementation  
./version project
```

### New Features

The Go implementation adds several new features not available in bash scripts:

- Enhanced version grammar support
- Colored terminal output
- Better error messages
- Cross-platform compatibility
- Improved performance

### Breaking Changes

- None - maintains full backward compatibility with existing bash script interfaces