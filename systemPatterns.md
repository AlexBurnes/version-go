# System Patterns: Version CLI Utility

## System Architecture
**✅ IMPLEMENTED WITH LIBRARY SUPPORT AND SELF-BUILDING** - The system follows a modular CLI architecture with clear separation of concerns, reusable library package, and self-building capabilities:

```
CLI Interface Layer ✅
├── Command Parser (flag package) ✅
├── Input/Output Handlers ✅
├── Error Handling & Exit Codes ✅
└── Configuration Management ✅
    ├── .project.yml File Parser ✅
    ├── Project Name Detection ✅
    ├── Module Name Detection ✅
    ├── Modules List Detection ✅
    ├── Git Fallback Integration ✅
    ├── CLI Options (--config, --git) ✅
    └── Test Configuration Files ✅

Library Package Layer ✅
├── pkg/version - Reusable Version Library ✅
│   ├── version.go - Core Version Parser & Logic ✅
│   │   ├── Version Parser (Custom Regex Engine) ✅
│   │   ├── Version Validator ✅
│   │   ├── Version Sorter ✅
│   │   ├── Type System (Release/Prerelease/Postrelease/Intermediate) ✅
│   │   └── Git Tag Conversion ✅
│   └── bump.go - Version Bumping Engine ✅
│       ├── BumpType Enum & Parsing ✅
│       ├── BumpResult Struct ✅
│       ├── Bump Rule Engine ✅
│       ├── Version State Detection ✅
│       ├── Smart Increment Logic ✅
│       └── Version Type Handlers ✅

Core Logic Layer ✅
├── CLI Command Handlers ✅
│   ├── Modules Command Handler ✅
│   └── Bump Command Handler ✅
├── Git Integration ✅
├── Library Integration ✅
└── Configuration Integration ✅
    ├── .project.yml Reader ✅
    ├── Project Info Provider ✅
    ├── Custom Config File Support ✅
    └── CLI Option Integration ✅

Utility Layer ✅
├── String Processing ✅
├── File I/O ✅
└── Platform Abstractions ✅

Build System Layer ✅
├── Main Entry: build-and-package.sh (orchestrates complete build flow) ✅
├── Local Development: Conan + CMake + bash scripts ✅
├── Packaging: GoReleaser + Conan hooks ✅
├── Cross-Platform: Automated builds for Linux/Windows/macOS ✅
├── Environment Setup: Automatic Go PATH configuration ✅
├── Environment Loading: .env file support for tokens and variables ✅
└── Self-Building: Version utility uses its own built binary for version detection ✅

Packaging Layer ✅
├── Linux: Makeself Self-Extracting Archives ✅
├── Windows: Scoop Package Manager ✅
├── macOS: Makeself Self-Extracting Archives ✅
└── Cross-Platform: GoReleaser Integration ✅
```

## Key Technical Decisions
- **✅ Configuration File Support**: .project.yml file for project name and module configuration with git fallback
- **✅ CLI Options**: --config FILE for custom config files, --git for forcing git detection
- **✅ YAML Parsing**: gopkg.in/yaml.v3 library for configuration file parsing
- **✅ Test Configuration**: Test files in test/ directory for different scenarios
- **✅ Modules Command**: New `modules` command to list all modules from .project.yml or single git module name
- **🔄 Version Bumping System**: Intelligent version increment with complex bump rules and state detection
- **✅ Self-Building System**: Version utility uses its own built binary for version detection during build process
- **✅ Bootstrap Process**: Initial build uses git describe, subsequent builds use built version utility
- **✅ Circular Dependency Resolution**: Eliminate dependency on git describe by using built version utility
- **✅ CLI Framework**: Implemented using Go's `flag` package for command parsing
- **✅ Library Package**: Refactored core functionality into reusable `pkg/version` package
- **✅ Grammar Engine**: Custom regex-based parser for extended version format support (avoiding semver library complexity)
- **✅ Git Integration**: Implemented using `os/exec` for git operations (standard library only)
- **✅ Build System**: CMake as orchestrator, Conan for dependency management, GoReleaser for distribution
- **✅ Local Development**: Conan + CMake + bash scripts for development and testing
- **✅ Packaging Build**: GoReleaser with Conan hooks for automated cross-platform builds
- **✅ Linux Packaging**: Makeself for self-extracting archives with professional installation experience
- **✅ Windows Packaging**: Scoop package manager for easy installation and updates
- **✅ macOS Packaging**: Makeself for self-extracting archives with professional installation experience
- **✅ Installer Naming**: Clean version numbers without SNAPSHOT and hex abbreviations
- **✅ Installation Approach**: No-sudo internal usage - users run with sudo if needed
- **✅ Build Script**: Automatic Go environment setup and .env file loading for seamless development
- **✅ Environment Management**: .env file support for GitHub tokens and other environment variables
- **✅ Git Remote Handling**: Automatic detection and configuration of git remotes for GoReleaser
- **✅ Custom Installation**: Support for custom installation directories via APP_DIR environment variable
- **✅ Testing**: Standard Go testing framework with comprehensive library tests, race detection, and performance testing
- **✅ Error Handling**: Structured error types with proper exit code mapping
- **✅ API Design**: Clean, well-documented public API with comprehensive examples

## Design Patterns in Use
- **✅ Command Pattern**: Each CLI command implemented as a separate handler in main.go
- **✅ Strategy Pattern**: Different parsing strategies for different version types (release, prerelease, postrelease, intermediate)
- **✅ Factory Pattern**: Version object creation based on input string analysis in ParseVersion()
- **✅ Builder Pattern**: Complex version object construction with validation
- **✅ Template Method**: Common validation and sorting logic with type-specific implementations

## Component Relationships
- **CLI Interface** → **Command Handlers** → **Library Package** → **Core Logic**
- **Library Package** → **Version Parser** ← **BNF Grammar Rules** → **Version Objects**
- **Git Integration** → **Library Package** → **Version Objects**
- **Build System** → **Go Compiler** → **Static Binaries** → **Distribution Packages**
- **External Applications** → **Library Package** → **Version Objects**

## Critical Implementation Paths
1. **✅ Version Parsing Pipeline**: Input validation → Grammar parsing → Object creation → Validation
2. **✅ Sorting Algorithm**: Parse all versions → Categorize by type → Apply precedence rules → Sort within categories
3. **✅ Git Integration**: Read git tags → Parse versions → Validate → Return appropriate version
4. **✅ Configuration Pipeline**: Check .project.yml → Parse YAML → Extract project/module info → Fallback to git if missing
5. **✅ Version Bumping Pipeline**: Parse current version → Detect version state → Apply bump rules → Generate new version
6. **✅ Self-Building Pipeline**: Bootstrap build (git describe) → Build version utility → Use built utility for subsequent builds
7. **✅ Local Build Pipeline**: Source code → Conan deps → CMake config → Go compilation → Static binary
8. **✅ Makeself Installer Pipeline**: Binary + install script + docs → Package directory → Makeself compression → Self-extracting archive
9. **✅ Scoop Integration Pipeline**: GoReleaser build → Scoop manifest generation → Package manager distribution

## Self-Building Implementation Pattern

### Bootstrap Process
1. **Initial Build**: Use `git describe` for version detection during first build
2. **Build Version Utility**: Compile version utility binary to `scripts/version` directory
3. **Subsequent Builds**: Use built version utility for all version operations

### Version Detection Strategy
- **CMakeLists.txt**: Use `scripts/version version` instead of `git describe`
- **Build Scripts**: Use `scripts/version version` for version detection
- **Pre-Push Hook**: Use `scripts/version check` for version validation
- **Version Checking**: Use `scripts/version check-greatest` for version comparison

### Circular Dependency Resolution
- **Problem**: Project needs git describe to build, but version utility replaces git describe
- **Solution**: Bootstrap process uses git describe initially, then switches to built version utility
- **Benefits**: Eliminates external git dependency, uses own version utility for consistency

## Packaging Implementation Details
### Simple Installer Scripts
- **Tool**: Custom shell scripts that download archives from GitHub releases
- **Format**: Single `.sh` file that downloads, extracts, and installs version utility
- **Naming**: `version-{platform}-{arch}-install.sh` (e.g., `version-linux-amd64-install.sh`)
- **Installation**: Users run `wget -O - URL | sh` or `wget -O - URL | sudo sh` for system installation
- **Custom Directory**: Users can specify custom installation directory using `INSTALL_DIR=install_dir` environment variable
- **Documentation Format**: All installation documentation shows correct `INSTALL_DIR=install_dir` format for simple installers
- **Features**: Downloads from GitHub releases, uses install.sh from archive, no-sudo internal approach
- **Platforms**: Linux and macOS with platform-specific install scripts
- **Cleanup Process**: Old installers are removed before creating new ones to prevent version mixing

### Installer Environment Variables
- **Simple Installers**: Use `INSTALL_DIR` environment variable for custom installation directory
  - Format: `wget -O - URL | INSTALL_DIR=/custom/path sh`
  - Fallback: Uses first argument if `INSTALL_DIR` not set
  - Default: `/usr/local/bin` if neither environment variable nor argument provided
- **Makeself Installers**: Use `APP_DIR` environment variable for custom installation directory
  - Format: `APP_DIR=/custom/path ./installer.run`
  - Used for self-extracting archives with professional branding
  - Different from simple installers to avoid conflicts

### GoReleaser Artifact Naming (v0.8.7)
- **Archive Naming**: Removed version prefix from artifact names for better download script compatibility
  - Format: `version-{os}-{arch}` instead of `version-{version}-{os}-{arch}`
  - Enables consistent download URLs for latest release scripts
  - Updated .goreleaser.yml archive name template with conditional logic
- **Installer Script Naming**: Removed version prefix from installer script names
  - Format: `version-{os}-{arch}-install.sh` instead of `version-{version}-{os}-{arch}-install.sh`
  - Enables consistent download URLs for latest release scripts
  - Updated installer creation scripts to use new naming convention
  - GoReleaser picks up installer scripts directly from installers/ directory
- **Platform Naming**: Renamed darwin to macos in user-facing documentation
  - GoReleaser goos field maintains "darwin" for compatibility
  - Archive output names use "macos" for clearer platform identification
  - Updated installer creation scripts to use macos platform name
  - Updated README.md installation examples to reference macos artifacts
- **GitHub Download URLs**: Updated all documentation with correct download URLs
  - Latest release scripts can use consistent URLs without version numbers
  - Updated README.md, packaging docs, and other documentation
  - Enables reliable download scripts for latest releases

### Scoop Package Manager Integration
- **Tool**: GoReleaser Scoop integration for Windows package management
- **Repository**: `https://github.com/AlexBurnes/scoop-bucket`
- **Installation**: `scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket && scoop install burnes/version`
- **Updates**: `scoop update burnes/version`
- **Features**: Automatic updates, dependency management, clean uninstallation

### Homebrew Tap Integration
- **Tool**: GoReleaser Homebrew integration for macOS package management
- **Repository**: `https://github.com/AlexBurnes/homebrew-tap`
- **Installation**: `brew tap AlexBurnes/homebrew-tap https://github.com/AlexBurnes/homebrew-tap && brew install version`
- **Updates**: `brew update && brew upgrade version`
- **Features**: Automatic updates, dependency management, clean uninstallation

### Version Cleaning Logic
- **Input**: `0.5.2-SNAPSHOT-5bb31e3` or `1.0.0-abc1234`
- **Output**: `0.5.2` or `1.0.0`
- **Pattern**: Removes `-SNAPSHOT-*` and `-{7-8 hex chars}` suffixes
- **Purpose**: Clean, user-friendly installer names without development artifacts

### Installer Cleanup Process
- **Problem**: GoReleaser glob patterns match all installers in directory, including old versions
- **Solution**: Clean installers directory before creating new ones in create_install_scripts()
- **Implementation**: `rm -rf installers/` before `./buildtools/create-all-installers.sh`
- **Benefit**: Ensures only current version installers are published
- **Verification**: Dry-run testing confirms only current version installers are present

## Build Script Enhancements (v0.5.4)

### Automatic Go Environment Setup
- **Function**: `setup_go_environment()` in `buildtools/build-goreleaser.sh`
- **Purpose**: Automatically detect and add Go's bin directory to PATH
- **Implementation**: Uses `go env GOPATH` to find Go workspace and adds `$GOPATH/bin` to PATH
- **Benefit**: No manual PATH configuration required for GoReleaser and other Go tools

### Environment Variable Loading
- **Function**: `load_environment()` in `buildtools/build-goreleaser.sh`
- **Purpose**: Load environment variables from `.env` file for GitHub tokens and other secrets
- **Implementation**: Uses `set -a` to automatically export variables from `.env` file
- **Benefit**: Secure token management without hardcoding secrets in scripts

### Git Remote Detection
- **Function**: Enhanced `run_goreleaser()` in `buildtools/build-goreleaser.sh`
- **Purpose**: Automatically detect and validate git remotes for GoReleaser publishing
- **Implementation**: Checks if origin remote exists and validates it's GitHub for publishing
- **Benefit**: Prevents publishing to wrong repositories and provides clear error messages

### Installation Documentation
- **Platforms**: Linux, Windows, macOS with comprehensive installation guides
- **Windows**: Complete Scoop setup instructions with PowerShell commands
- **Linux/macOS**: Multiple installation options (user, system-wide, custom directory)
- **Custom Installation**: Support for `APP_DIR` environment variable for custom paths
- **Benefit**: Clear, step-by-step instructions for all user types and preferences

### Performance Testing System
- **Large Version Lists**: Tests with 10,000+ versions for performance validation
- **Benchmark Testing**: Automated benchmarks for sorting and validation operations
- **Test Data Generation**: Realistic version pattern generation (release, prerelease, postrelease, intermediate)
- **Performance Assertions**: Time-based limits (5-second limit for 10k versions)
- **Binary Testing**: Performance testing with both source and compiled binary
- **Metrics Collection**: Duration measurement and performance logging

### Developer Workflow System
- **VERSION File**: Centralized version management with single source of truth
- **Pre-Push Hook**: Automated version validation and increment checking
- **Release Process**: Standardized workflow from planning to publication
- **Documentation**: Complete developer workflow with troubleshooting guide
- **Validation**: Prevents accidental releases and ensures version consistency
- **Automation**: Reduces human error in release process