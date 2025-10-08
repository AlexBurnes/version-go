# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.2] - 2025-10-08

### Fixed
- **GetVersionWithPrefix Git Tag Conversion**: Fixed `GetVersionWithPrefix()` to properly convert git tag format
  - Added `ConvertGitTag()` call to convert git tag delimiter from `-` to `~` for prerelease versions
  - Example: git tag `v0.7.14-pre.1` now correctly returns `v0.7.14~pre.1` (was returning `v0.7.14-pre.1`)
  - Ensures consistency with library version parsing expectations
  - Aligns with `GetVersion()` behavior which already performed this conversion

### Added
- **GetRawTag Function**: Added new library API method to retrieve raw git tag without transformations
  - New function `version.GetRawTag()` returns exact git tag as it appears in repository
  - No 'v' prefix removal, no delimiter conversion from `-` to `~`
  - Useful for applications that need the original git tag format for external integrations
  - Comprehensive test coverage in `pkg/version/git_test.go`
  - Full documentation with usage examples

### Summary
The library now provides three clear options for git tag retrieval:
1. **`GetVersion()`** - Returns version without 'v' prefix, converts `-` to `~` (e.g., `0.7.14~pre.1`)
2. **`GetVersionWithPrefix()`** - Returns version with 'v' prefix, converts `-` to `~` (e.g., `v0.7.14~pre.1`)
3. **`GetRawTag()`** - Returns exact git tag without transformations (e.g., `v0.7.14-pre.1`)

## [1.4.1] - 2025-10-08

### Fixed
- **GitHub Actions Test Failures**: Fixed cross-platform test workflow to properly fetch git tags
  - Updated `.github/workflows/cross-platform-test.yml` to use `fetch-depth: 0` in checkout actions
  - Resolves "No version tags found" errors in git-dependent tests on GitHub Actions
  - Ensures all git-related tests can access version tags during CI/CD execution
- **macOS Test Expectation**: Fixed platform test expectation for Darwin/macOS OS version detection
  - Updated `pkg/version/platform_test.go` to accept actual macOS version numbers (e.g., "15.6.1")
  - Removed incorrect expectation that OS version should start with "darwin" prefix
  - Test now correctly validates that `GetOSVersion()` returns actual version numbers on macOS
  - Aligns test expectations with correct implementation behavior from v1.2.5 bug fix

### Added
- **New Library API Methods**: Added GetVersionType and GetBuildTypeFromVersion methods to library package
  - New function `version.GetVersionType(versionStr string) (string, error)` returns version type with optional git fallback
  - New function `version.GetBuildTypeFromVersion(versionStr string) (string, error)` returns build type with optional git fallback
  - Both methods accept empty string to automatically retrieve version from git tags
  - Comprehensive test coverage with both version string and git integration testing
  - Full documentation in `docs/Library.md` with usage examples and API reference
  - Enables library consumers to get version type and build type information programmatically
  - Provides consistent API for pre-push project and other consumers needing version type/build type data

## [1.3.0] - 2025-10-01

### Added
- **Git Version Integration**: Added GetVersion() function to library package for retrieving version from git
  - New function `version.GetVersion()` retrieves current version from git tags without 'v' prefix
  - New function `version.GetVersionWithPrefix()` retrieves version with 'v' prefix preserved
  - Added `GitError` type for structured error handling with type checking helpers
  - Added helper functions `IsGitNotFound()`, `IsNotGitRepo()`, and `IsNoGitTags()` for error type detection
  - Comprehensive test coverage in `pkg/version/git_test.go`
  - Full documentation in `docs/Library.md` with usage examples
  - Example implementation in `examples/git-version/main.go`
  - Updated README.md with library usage examples demonstrating git integration
  - Enables library consumers to programmatically retrieve version information from git repositories

## [1.2.7] - 2025-10-01

### Fixed
- **Version Sort Order Bug**: Fixed incorrect version sorting where prerelease versions were considered greater than release versions
  - Corrected Type constant ordering in pkg/version/version.go to put TypePrerelease before TypeRelease
  - **Correct order** for versions with the same x.y.z is now: prerelease < release < postrelease < intermediate
  - Previously, v1.3.9-rc.9 was incorrectly identified as greater than v1.3.9
  - Now v1.3.9 is correctly identified as the greatest version among prereleases
  - Updated all tests to reflect correct version precedence
  - Updated documentation (Library.md, Project-specification.md, BNF-grammar.md) with correct precedence order
  - This fix aligns with the project specification where release versions should be greater than prerelease versions

## [1.2.6] - 2025-01-25

### Fixed
- **Static Build Configuration**: Fixed static build configuration for Linux and Darwin platforms
  - Added CGO_ENABLED=0 environment variable to GoReleaser configuration to disable CGO
  - Added -extldflags "-static" to ldflags to create static binaries
  - Ensures Linux and Darwin binaries are completely static with no external library dependencies
  - Improves binary portability and deployment across different Linux distributions and macOS versions

## [1.2.4] - 2025-01-25

### Fixed
- **macOS Version Detection**: Fixed getDarwinVersion() function to return actual macOS version numbers
  - Updated getDarwinVersion() to use `sw_vers -productVersion` command instead of returning "darwin"
  - Added proper error handling with fallback to MACOSX_DEPLOYMENT_TARGET environment variable
  - Now returns actual macOS version (e.g., "14.0", "15.0") instead of generic "darwin" string
  - Maintains backward compatibility with existing fallback mechanisms
  - Resolves GitHub issue #2: macOS OS Version Detection Returns "darwin" Instead of Numeric Version

## [1.2.3] - 2025-01-25

### Changed
- **Darwin Platform Naming**: Fixed inconsistent platform naming throughout the system
  - Fixed GoReleaser archive naming template to use darwin instead of macos
  - Updated installer script creation to use darwin instead of macos platform name
  - Renamed installer files from version-macos-* to version-darwin-*
  - Updated packaging directory from macos/ to darwin/
  - Updated all documentation to use darwin consistently
  - Updated test files and Docker files to use darwin naming
  - Updated build scripts and CMakeLists.txt to use darwin platform detection
  - Ensures consistent platform identification across all components
- **CI/CD Pipeline**: Removed cross-compilation tests from GitHub Actions workflows
  - Removed cross-compilation test step from cross-platform-test.yml workflow
  - Updated documentation to reflect focus on Docker-based cross-platform testing
  - Maintained Docker-based cross-platform tests for Ubuntu, Debian, and Windows (Wine)
  - Improved CI/CD reliability by removing complex cross-compilation dependencies

## [1.2.0] - 2025-01-24

### Added
- **Buildfab Migration**: Migrated build system from custom bash scripts to buildfab utility
  - Updated .project.yml with comprehensive build stages (pre-push, build, test, release)
  - Added buildfab actions for Conan dependency management and CMake configuration
  - Implemented unified build orchestration with dependency management
  - Added cross-platform testing stage with Docker integration
  - Added check-binaries action for simplified test stage verification
  - Added individual platform testing actions (test-linux-ubuntu, test-linux-debian, test-windows, test-macos)
  - Implemented parallel cross-platform testing using buildfab's native parallel execution
  - Each platform test action directly runs Docker build and Docker run commands
  - Preserved all existing functionality while gaining buildfab capabilities
  - Ready for buildfab-based development workflow with unified build management

### Changed
- **Build System**: Replaced custom build scripts with buildfab actions
  - Conan dependency management through buildfab actions
  - CMake configuration and build orchestration via buildfab
  - GoReleaser integration through buildfab actions
  - Cross-platform testing via Docker in buildfab test stage
  - Unified build management replacing multiple custom scripts
  - Simplified test stage to only check binaries exist in bin/ directory

### Removed
- **Unused Build Scripts**: Removed obsolete buildtools scripts after buildfab migration
  - Removed build-and-package.sh (replaced by buildfab build stage)
  - Removed build-conan.sh (replaced by buildfab install-conan-deps action)
  - Removed build-goreleaser.sh (replaced by buildfab goreleaser actions)
  - Removed collect_bins.sh (replaced by inline command in .goreleaser.yml)
  - Removed post-archive-hook.sh (no longer needed with buildfab)
  - Removed post-build-hook.sh (no longer needed with buildfab)
  - Cleaned up buildtools directory to only contain actively used scripts
  - Updated documentation to use buildfab commands instead of run-cross-platform-tests.sh script
  - Cross-platform testing now uses buildfab's native parallel execution capabilities
  - Removed simple-cross-platform-test.sh script (redundant with CMake cross-compilation)
  - CMake version-all target already handles cross-platform binary building
  - Updated Developer-workflow.md with complete buildfab workflow documentation
  - Added comprehensive buildfab stage and action documentation
  - Updated README.md with buildfab installation and usage instructions
  - Added buildfab installation commands for system and local installation
  - Documented .project.yml configuration and buildfab stages
  - Added pre-push hook setup documentation for developers
  - Documented pre-push utility installation and usage
  - Added automatic project validation before git push
  - Added buildfab ecosystem documentation to README and Developer Workflow
  - Documented project as part of buildfab utilities and libraries ecosystem

### Fixed
- **Install-Binary Action**: Fixed CMakeLists.txt install targets
  - Removed non-existent config/version.yaml file copy operations
  - Fixed install-current and install-all targets to only copy binaries
  - Fixed install-binary action failure by removing non-existent config file dependencies
  - Added check-binaries action for simplified test stage verification
  - Updated CMakeLists.txt to remove config/version.yaml file copy operations
  - Preserved all existing functionality while gaining buildfab orchestration capabilities
  - Ready for buildfab-based development workflow with unified build management

## [1.1.1] - 2025-01-24

### Added
- **Buildfab Migration**: Migrated build system from custom bash scripts to buildfab utility
  - Updated .project.yml with comprehensive build stages (pre-push, build, test, release)
  - Added buildfab actions for Conan dependency management and CMake configuration
  - Implemented unified build orchestration with dependency management
  - Added cross-platform testing stage with Docker integration
  - Added check-binaries action for simplified test stage verification
  - Added individual platform testing actions (test-linux-ubuntu, test-linux-debian, test-windows, test-macos)
  - Implemented parallel cross-platform testing using buildfab's native parallel execution
  - Each platform test action directly runs Docker build and Docker run commands
  - Preserved all existing functionality while gaining buildfab capabilities
  - Ready for buildfab-based development workflow with unified build management

### Changed
- **Build System**: Replaced custom build scripts with buildfab actions
  - Conan dependency management through buildfab actions
  - CMake configuration and build orchestration via buildfab
  - GoReleaser integration through buildfab actions
  - Cross-platform testing via Docker in buildfab test stage
  - Unified build management replacing multiple custom scripts
  - Simplified test stage to only check binaries exist in bin/ directory

### Removed
- **Unused Build Scripts**: Removed obsolete buildtools scripts after buildfab migration
  - Removed build-and-package.sh (replaced by buildfab build stage)
  - Removed build-conan.sh (replaced by buildfab install-conan-deps action)
  - Removed build-goreleaser.sh (replaced by buildfab goreleaser actions)
  - Removed collect_bins.sh (replaced by inline command in .goreleaser.yml)
  - Removed post-archive-hook.sh (no longer needed with buildfab)
  - Removed post-build-hook.sh (no longer needed with buildfab)
  - Cleaned up buildtools directory to only contain actively used scripts
  - Updated documentation to use buildfab commands instead of run-cross-platform-tests.sh script
  - Cross-platform testing now uses buildfab's native parallel execution capabilities
  - Removed simple-cross-platform-test.sh script (redundant with CMake cross-compilation)
  - CMake version-all target already handles cross-platform binary building
  - Updated Developer-workflow.md with complete buildfab workflow documentation
  - Added comprehensive buildfab stage and action documentation
  - Updated README.md with buildfab installation and usage instructions
  - Added buildfab installation commands for system and local installation
  - Documented .project.yml configuration and buildfab stages
  - Added pre-push hook setup documentation for developers
  - Documented pre-push utility installation and usage
  - Added automatic project validation before git push
  - Added buildfab ecosystem documentation to README and Developer Workflow
  - Documented project as part of buildfab utilities and libraries ecosystem

### Fixed
- **Install-Binary Action**: Fixed CMakeLists.txt install targets
  - Removed non-existent config/version.yaml file copy operations
  - Fixed install-current and install-all targets to only copy binaries
  - Eliminated config directory creation and file copying dependencies
  - Resolved buildfab install-binary action failure

## [1.1.1] - 2025-09-24

### Fixed
- **GitHub Actions Platform Tests**: Fixed build command syntax in cross-platform test workflow
  - Changed `go build -o bin/version cmd/version/*.go` to `go build -o bin/version ./cmd/version`
  - Fixed malformed import path error with invalid `*` character in build commands
  - Updated all cross-platform build commands to use proper Go build syntax
  - Verified platform detection commands work correctly on all target platforms
  - Cross-compilation tests now pass successfully for Linux, Windows, and macOS
- **Cross-Platform Docker Tests**: Fixed Docker build context issues in cross-platform test suite
  - Updated Docker build commands to run from project root with correct context
  - Fixed Dockerfile paths to use proper relative paths from project root
  - Updated all Dockerfiles (Ubuntu, Debian, Windows) to use correct file paths
  - Rebuilt all platform binaries with latest platform detection code
  - Verified all cross-platform tests pass successfully on Ubuntu, Debian, and Windows
  - **GitHub Actions Workflow**: Updated cross-platform test workflow to match local script fixes
  - Fixed GitHub Actions Docker build commands to run from project root instead of test directory
  - Ensured GitHub Actions uses same Docker build context as local scripts

## [1.1.0] - 2025-09-24

### Added
- **Platform Detection Commands**: Added new CLI commands for platform detection
  - `platform` - Returns current platform (GOOS value): `linux`
  - `arch` - Returns current architecture (GOARCH value): `amd64`
  - `os` - Returns current operating system (distribution name): `ubuntu`
  - `os_version` - Returns current OS version: `24.04`
  - `cpu` - Returns number of logical CPUs: `8`
- **Platform Detection Library API**: Enhanced `pkg/version` package with platform detection functions
  - `GetPlatform()` - Returns current platform name
  - `GetArch()` - Returns current architecture name
  - `GetOS()` - Returns current operating system name (distribution-specific on Linux)
  - `GetOSVersion()` - Returns current OS version
  - `GetNumCPU()` - Returns number of logical CPUs
  - `GetPlatformInfo()` - Returns comprehensive platform information struct
- **Cross-Platform Testing Infrastructure**: Comprehensive testing framework for platform detection
  - Docker-based testing for Ubuntu, Debian, and Windows (Wine)
  - Cross-compilation testing for all target platforms
  - GitHub Actions integration for automated testing
  - Test scripts for Linux, Windows, and macOS platforms
  - Comprehensive documentation and troubleshooting guides

### Changed
- **Platform Detection Logic**: Enhanced Linux OS detection to return actual distribution names
  - Linux platforms now return distribution name (e.g., `ubuntu`, `debian`) instead of generic `linux`
  - OS version detection improved to return actual version numbers (e.g., `24.04`, `12`)
  - Darwin platform detection keeps `darwin` naming (not converted to `macos`)
- **PlatformInfo Struct**: Added `NumCPU` field for CPU count information
- **CLI Help Text**: Updated help text to include new platform detection commands

### Technical Details
- **Library Functions**: All platform detection functions use Go's built-in runtime capabilities
- **Linux Detection**: Reads `/etc/os-release` to detect distribution name and version
- **Cross-Platform Support**: Works on Linux (amd64/arm64), Windows (amd64/arm64), macOS (amd64/arm64)
- **Testing Coverage**: Comprehensive testing across multiple Linux distributions and simulated Windows environment
- **Documentation**: Added detailed testing guides and cross-platform validation procedures

## [1.0.1] - 2025-09-22

### Documentation
- **Version Management Rules Enhancement**: Enhanced version management rules to ensure proper packaging updates
  - **VERSION File Updates**: Made VERSION file updates mandatory when bumping version
  - **Packaging Configuration**: Added requirements to update Windows Scoop and macOS Homebrew configurations
  - **Automated Workflow**: Updated `scripts/version-bump-with-file` to automatically update packaging files
  - **Rule Updates**: Enhanced `rule-versioning.mdc` and `rule-complete-changes.mdc` with packaging requirements
  - **Files Updated**:
    - `.cursor/rules/rule-versioning.mdc` - Added mandatory packaging updates section
    - `.cursor/rules/rule-complete-changes.mdc` - Added packaging files information to workflow
  - **Packaging Files**: When version changes, these files must be updated:
    - `VERSION` file with new version number
    - `packaging/windows/scoop-bucket/version.json` with new version and URLs
    - `packaging/macos/version.rb` with new URLs

### Documentation
- **Version Management Rules**: Enhanced version management rules to ensure proper packaging updates
  - **VERSION File Updates**: Made VERSION file updates mandatory when bumping version
  - **Packaging Configuration**: Added requirements to update Windows Scoop and macOS Homebrew configurations
  - **Automated Workflow**: Updated `scripts/version-bump-with-file` to automatically update packaging files
  - **Rule Updates**: Enhanced `rule-versioning.mdc` and `rule-complete-changes.mdc` with packaging requirements
  - **Files Updated**:
    - `.cursor/rules/rule-versioning.mdc` - Added mandatory packaging updates section
    - `.cursor/rules/rule-complete-changes.mdc` - Added packaging files information to workflow
  - **Packaging Files**: When version changes, these files must be updated:
    - `VERSION` file with new version number
    - `packaging/windows/scoop-bucket/version.json` with new version and URLs
    - `packaging/macos/version.rb` with new URLs

### Documentation
- **Documentation Reorganization**: Reorganized project documentation structure
  - **File Naming Convention**: Applied lowercase-first-word-with-dash naming convention to documentation files
  - **Directory Structure**: Moved all documentation files (except README and CHANGELOG) to docs/ directory
  - **Moved Files**:
    - `DEVELOPER_WORKFLOW.md` → `docs/Developer-workflow.md`
    - `IMPLEMENTATION_SUMMARY.md` → `docs/Implementation-summary.md`
    - `project.md` → `docs/Project-specification.md`
  - **Renamed Files in docs/ directory**:
    - `CI_CD.md` → `Deploy.md` (more descriptive name)
    - `LIBRARY.md` → `Library.md`
    - `NEW_FEATURES_SPECIFICATION.md` → `New-features-specification.md`
    - `BNF_GRAMMAR.md` → `BNF-grammar.md` (kept BNF uppercase as abbreviation)
    - `BUILD.md` → `Build.md`
    - `GORELEASER_HOOKS.md` → `Goreleaser-hooks.md`
  - **Reference Updates**: Updated all references throughout codebase to point to new file locations and names
  - **Memory Bank Files**: Kept memory bank documents in root directory as they are referenced by MCP server
  - **Consistency**: Ensured all documentation follows consistent naming and location patterns
  - **Documentation Rules**: Created comprehensive `rule-documents.mdc` with naming conventions and directory structure requirements

## [0.8.22] - 2025-09-17

### Added
- **DRY Refactoring for Build Scripts**: Extracted common functionality into reusable scripts
  - **GoReleaser Binary Backup Script**: `buildtools/create-goreleaser-backup.sh` for consistent binary backup
  - **GoReleaser Archive Creation Script**: `buildtools/create-goreleaser-archives.sh` for platform-specific archives
  - **GitHub Actions Integration**: Updated release workflow to use new scripts instead of inline code
  - **Build Script Integration**: Updated `build-and-package.sh` to use new scripts for consistency
  - **Error Handling**: Improved error handling with graceful fallbacks for missing binaries
  - **Version Detection**: Automatic version detection with fallback to VERSION file or git describe
  - **Cross-Platform Support**: Proper handling of Windows .exe extensions and Unix binary naming

### Changed
- **GitHub Actions Workflow**: Simplified release workflow by replacing inline steps with script calls
- **Build Process**: Streamlined build process by eliminating code duplication between GitHub Actions and build scripts
- **Script Organization**: Better separation of concerns with dedicated scripts for specific build tasks

## [0.8.21] - 2025-09-17

### Added
- **GitHub Actions CI/CD Pipeline (v0.8.14-v0.8.21)**: Comprehensive continuous integration and deployment workflow
  - **CI Workflow**: Organized into three phases: build → test → package
  - **Build Phase**: Pre-builds `bin/version` for testing, uses Conan for golang package management
  - **Test Phase**: Runs Go tests with race detection and linting, no Conan setup required
  - **Package Phase**: Builds all platforms, creates installer scripts, runs GoReleaser dry-run
  - **Conan Integration**: Automatic golang package management with local creation if not available
  - **Go Module Caching**: Added caching for Go modules across all phases
  - **Release Workflow**: Uses `buildtools/build-and-package.sh release` for consistent releases
  - **Package Management**: Conan 2.0 integration for Go toolchain dependencies
  - **Local Golang Recipe**: `conanfile-golang.py` for Go 1.23.0 with cross-platform support
  - **Automated Package Creation**: Creates golang package locally when not available in remote repositories
  - **Build Script Integration**: Updated all build scripts to include golang package checks
  - **Comprehensive Documentation**: Complete CI/CD workflow documentation in `docs/CI_CD.md`

- **Auto-Download Version Utility**: Enhanced bootstrap process with automatic download of latest version utility from GitHub releases
  - **Priority Order**: Built utility → Auto-download → Git describe fallback
  - **Platform Detection**: Automatic detection of platform (linux/macos) and architecture (amd64/arm64)
  - **Latest Release URLs**: Uses `/releases/latest/download/` URLs without version numbers for consistent downloads
  - **Simplified Download**: Uses `wget -O - url | INSTALL_DIR=./scripts sh` pipe approach for streamlined installation
  - **Cross-Platform Support**: Works on Linux and macOS with proper architecture mapping
  - **Error Handling**: Graceful fallback to git describe if download fails
  - **Build Script Integration**: Updated all build scripts (build-and-package.sh, build-conan.sh, check-version-status, CMakeLists.txt)
  - **Zero Setup**: New developers can build immediately without git tags or existing version utility
  - **CI/CD Friendly**: Works in clean CI environments without pre-existing version utility
- **Simplified Bump Workflow (v0.8.10)**: Streamlined bump functionality for build script usage
  - **Direct Usage**: Use `scripts/version bump` directly instead of wrapper scripts
  - **Removed Redundancy**: Eliminated `scripts/version-bump` wrapper script (only used in documentation examples)
  - **Updated References**: All documentation and scripts now use direct version utility commands
  - **Simplified Workflow**: Cleaner, more straightforward approach to version bumping in build scripts
  - **Comprehensive Support**: Full support for all bump types (major, minor, patch, pre, alpha, beta, rc, fix, next, post, feat, smart)

### Fixed
- **Git Remote Priority (v0.8.17)**: Fixed git remote priority to ensure consistent behavior across environments
  - **Origin Priority**: Modified `getProjectFromGit()` and `getModuleFromGit()` functions to prioritize `origin` remote
  - **Two-Pass Logic**: First pass looks for `origin` remote, second pass falls back to any other remote
  - **Consistent Behavior**: Ensures same git remote detection behavior in local and CI environments
  - **Test Updates**: Updated test expectations to match origin remote behavior (`AlexBurnes-version-go`)
  - **Cross-Platform**: Works consistently across different git remote configurations

### Removed
- **Redundant Scripts (v0.8.10)**: Removed unnecessary wrapper scripts
  - **scripts/version-bump**: Removed redundant wrapper script that only provided documentation examples
  - **Simplified Maintenance**: Reduced codebase complexity by eliminating unnecessary wrapper layer

### Changed
- **Build Script Integration (v0.8.10)**: Updated build script usage patterns
  - **Direct Commands**: Updated `scripts/check-version-status` to suggest `scripts/version bump` directly
  - **Documentation**: Updated all examples to use direct version utility commands
  - **Help Text**: Updated help text to reflect simplified workflow

## [0.8.9] - 2025-09-16

### Added
- **Bump Command**: New `bump` command for intelligent version incrementing
  - Support for major, minor, patch, and specific identifier bump types (pre, alpha, beta, rc, fix, next, post, feat)
  - Smart bump mode that automatically determines appropriate increment based on current version type
  - Comprehensive bump functionality in `pkg/version` library package
  - Detailed help text and usage examples for bump command
  - Full test coverage for all bump scenarios and edge cases
  - Support for complex version suffixes with proper numeric increment handling
  - Intuitive naming convention matching version identifier prefixes
  - **Code Organization**: Decomposed bump functionality into separate `pkg/version/bump.go` file for better maintainability
  - **Code Structure**: Improved separation of concerns with core version parsing in `version.go` and bump logic in `bump.go`

### Changed
- **Code Organization**: Improved code structure and maintainability
  - Decomposed bump functionality from `pkg/version/version.go` to separate `pkg/version/bump.go` file
  - Better separation of concerns between core version parsing and bump logic
  - Cleaner file organization with focused responsibilities
  - Maintained full backward compatibility and functionality

### Fixed
- **Bump Command Argument Parsing**: Fixed bump command argument parsing logic
  - Now correctly handles single argument as either version or bump type
  - `version bump patch` now works correctly (uses current git version with patch bump)
  - `version bump 1.2.3` now works correctly (uses provided version with smart bump)
  - `version bump 1.2.3 patch` now works correctly (uses provided version with patch bump)
  - Improved argument validation to accept valid versions or bump types as single argument
  - Enhanced error messages to clearly indicate whether argument should be version or bump type
- **Version Status Script**: Fixed version comparison logic in `scripts/check-version-status`
  - Now correctly normalizes version strings by removing 'v' prefix before comparison
  - Prevents false positive "ready for changes" when VERSION file matches current git tag
  - Ensures accurate version status reporting for development workflow

## [0.8.7] - 2025-09-16

### Changed
- **GoReleaser Artifact Naming**: Updated artifact naming to remove version prefix for better download script compatibility
  - Archive names changed from `version-{version}-{os}-{arch}` to `version-{os}-{arch}` format
  - Installer script names changed from `version-{version}-{os}-{arch}-install.sh` to `version-{os}-{arch}-install.sh` format
  - This enables consistent download URLs for latest release scripts
  - Updated .goreleaser.yml archive name template to remove version component
  - Updated installer creation scripts to use new naming convention
  - Updated README.md installation examples to use new artifact names
  - GoReleaser picks up installer scripts directly from installers/ directory without copying to dist/

### Documentation
- **GitHub Download URLs**: Updated all documentation with correct download URLs for latest release scripts
  - Updated README.md installation examples to use correct GitHub URLs without version numbers
  - Updated packaging documentation with correct download instructions
  - Updated memory bank files to reflect all artifact naming and documentation changes
  - Enables reliable download scripts for latest releases across all platforms

### Added
- **Platform Naming Consistency**: Fixed inconsistent platform naming throughout the system
  - Reverted GoReleaser archive naming to use "darwin" consistently with Go's platform naming
  - Updated README.md installation examples to reference darwin artifacts
  - Updated installer creation scripts to use darwin platform name consistently
  - Maintained darwin naming throughout the system for consistency with Go's platform conventions
  - This ensures consistent platform identification across all components

## [0.8.6] - 2025-09-16

### Added
- **Modules Command**: Added new `modules` command to list all modules from .project.yml configuration file
  - Returns all modules from .project.yml if defined, otherwise falls back to single git module name
  - Supports --config and --git flags for configuration control
  - Outputs modules as newline-separated list for easy parsing
  - Integrated with existing configuration system and CLI framework
  - Added comprehensive tests for single and multiple module scenarios

## [0.8.3] - 2025-09-15

### Fixed
- **Install Command Documentation**: Fixed install command documentation format to show correct `INSTALL_DIR=install_dir` format instead of `APP_DIR` for simple installers
  - Updated README.md both Linux and macOS installation sections
  - Enhanced build script output to show both argument and environment variable formats
  - Ensured documentation consistency across all installation methods
  - Verified installer template correctly handles `INSTALL_DIR` environment variable

### Documentation
- **README.md**: Updated installation commands to use correct `INSTALL_DIR=install_dir` format
- **buildtools/create-simple-installers.sh**: Enhanced help output to show both installation formats

## [0.8.2] - 2025-09-15

### Fixed
- **GoReleaser Old Installers Issue**: Fixed GoReleaser publishing old version installers by cleaning installers directory before creating new ones
  - Added cleanup step in create_install_scripts() function to remove old installers
  - Enhanced clean_build function to be more explicit about cleaning installers
  - Verified build process now only publishes current version installers
  - Tested dry-run process to confirm fix works correctly

### Technical Details
- Build process now follows: Create install scripts → Clean installers directory → Create new installers → Run GoReleaser
- Only current version installers are published, eliminating confusion from mixed version installers

## [0.8.1] - 2025-09-15

### Fixed
- **Packaging and Release Clean Function**: Fixed clean_build function in build-and-package.sh to properly clean installers/ directory
  - Added .goreleaser-binaries/ directory to clean function for complete cleanup
  - Verified clean function removes all build artifacts (dist/, bin/, installers/, .goreleaser-binaries/)
  - Ensured installers directory is empty before release as required

### Technical Details
- Complete build cleanup now removes all build artifacts before new builds
- Prevents old installer files from interfering with new releases

## [0.8.0] - 2025-09-15

### Added
- **Project Configuration Support**: Added `.project.yml` configuration file support for project name and module configuration
  - Configuration file takes precedence over git-based detection when present
  - Supports multiple modules with first module as primary
  - Graceful fallback to git-based detection when configuration is missing or invalid
  - YAML parsing with comprehensive validation and error handling
  - Updated `project` and `module` commands to use configuration when available
  - Added debug output to show configuration source (file vs git fallback)

### Changed
- **CLI Options**: Added `--config FILE` and `--git` options for configuration control
  - `--config FILE`: Specify custom configuration file path
  - `--git`: Force git-based detection instead of configuration file
- **Command Behavior**: `project` and `module` commands now use configuration when available

### Technical Details
- Added `gopkg.in/yaml.v3` dependency for YAML parsing
- Configuration validation with comprehensive error handling
- Test configuration files in `test/` directory for different scenarios

## [0.7.0] - 2025-09-15

### Added
- **Self-Building Version Utility**: Implemented self-building capability where version utility uses its own built binary
  - Version utility now uses its own built binary for version detection during build process
  - Bootstrap process uses git describe initially, then switches to built version utility
  - Eliminates external git dependency for version detection

### Changed
- **CMakeLists.txt**: Updated to use built version utility instead of git describe
- **Build Scripts**: Modified all build scripts to use built version utility for version detection
- **Pre-Push Hook**: Updated to use built version utility for version checking

### Technical Details
- Circular dependency resolution: Initial build uses git describe, subsequent builds use built version utility
- Version detection strategy: Use built version utility in scripts/ directory for all version operations
- Self-building process tested and verified to work correctly

## [0.6.1] - 2025-09-15

### Fixed
- **License Inconsistencies**: Fixed license inconsistencies across all documentation to reference Apache 2.0 License
  - Updated all documentation to consistently reference Apache 2.0 License
  - Fixed license mentions in README.md, packaging docs, and library docs

### Added
- **BNF Grammar Documentation**: Created comprehensive BNF grammar specification document (docs/BNF_GRAMMAR.md)
  - Complete grammar specification matching the implementation
  - Detailed explanation of version format rules and precedence
  - Examples for all supported version types

### Documentation
- **BNF Grammar References**: Updated project documentation to reference BNF grammar specification
- **License Consistency**: All documentation now consistently references Apache 2.0 License

## [0.6.0] - 2025-09-12

### Changed
- **Build Script Rename**: Renamed `buildtools/build-with-conan.sh` to `buildtools/build-and-package.sh` for better semantic clarity
  - The new name better describes the script's purpose (builds binaries AND creates packages)
  - Updated all documentation references to use the new script name
  - Enhanced BUILD.md with main build script section explaining the primary entry point
  - Updated developer workflow documentation with new script name
  - Updated project.md main build scripts list

### Documentation
- **BUILD.md**: Added "Main Build Script" section explaining `build-and-package.sh` as the primary entry point
- **README.md**: Updated developer workflow examples to use new script name
- **DEVELOPER_WORKFLOW.md**: Updated all references throughout the workflow documentation
- **project.md**: Updated main build scripts list to include the new script name

## [0.5.2] - 2025-09-11

### Verified
- **GoReleaser Integration**: Cross-platform build and distribution system fully functional
- **Conan Hooks**: GoReleaser hooks properly install and deploy build tools via Conan
- **Distribution Artifacts**: All platform archives and binaries generated correctly
- **Conan Profile Detection**: Respects existing profiles without overwriting
- **Apache 2.0 License**: Properly integrated for distribution

## [0.5.1] - 2025-09-11

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

## [0.5.0] - 2025-09-11

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

## [0.4.0] - 2025-09-11

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

## [0.3.0] - 2025-09-11

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