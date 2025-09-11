# System Patterns: Version CLI Utility

## System Architecture
**✅ IMPLEMENTED WITH LIBRARY SUPPORT** - The system follows a modular CLI architecture with clear separation of concerns and reusable library package:

```
CLI Interface Layer ✅
├── Command Parser (flag package) ✅
├── Input/Output Handlers ✅
└── Error Handling & Exit Codes ✅

Library Package Layer ✅
├── pkg/version - Reusable Version Library ✅
│   ├── Version Parser (Custom Regex Engine) ✅
│   ├── Version Validator ✅
│   ├── Version Sorter ✅
│   ├── Type System (Release/Prerelease/Postrelease/Intermediate) ✅
│   └── Git Tag Conversion ✅

Core Logic Layer ✅
├── CLI Command Handlers ✅
├── Git Integration ✅
└── Library Integration ✅

Utility Layer ✅
├── String Processing ✅
├── File I/O ✅
└── Platform Abstractions ✅

Build System Layer ✅
├── Local Development: Conan + CMake + bash scripts ✅
├── Packaging: GoReleaser + Conan hooks ✅
├── Cross-Platform: Automated builds for Linux/Windows/macOS ✅
├── Environment Setup: Automatic Go PATH configuration ✅
└── Environment Loading: .env file support for tokens and variables ✅

Packaging Layer ✅
├── Linux: Makeself Self-Extracting Archives ✅
├── Windows: Scoop Package Manager ✅
├── macOS: Makeself Self-Extracting Archives ✅
└── Cross-Platform: GoReleaser Integration ✅
```

## Key Technical Decisions
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
- **✅ Testing**: Standard Go testing framework with comprehensive library tests and race detection
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
4. **✅ Local Build Pipeline**: Source code → Conan deps → CMake config → Go compilation → Static binary
5. **✅ Makeself Installer Pipeline**: Binary + install script + docs → Package directory → Makeself compression → Self-extracting archive
6. **✅ Scoop Integration Pipeline**: GoReleaser build → Scoop manifest generation → Package manager distribution

## Packaging Implementation Details
### Makeself Self-Extracting Installers
- **Tool**: [makeself.sh](https://github.com/megastep/makeself) for creating self-extracting archives
- **Format**: Single `.sh` file containing compressed binary, install script, and documentation
- **Naming**: `version-{clean_version}-{platform}-{arch}-install.sh` (e.g., `version-0.5.2-linux-amd64-install.sh`)
- **Installation**: Users run `wget -O - URL | sh` or `wget -O - URL | sudo sh` for system installation
- **Features**: Professional header, integrity checking, no-sudo internal approach
- **Platforms**: Linux and macOS with platform-specific install scripts

### Scoop Package Manager Integration
- **Tool**: GoReleaser Scoop integration for Windows package management
- **Repository**: `https://github.com/AlexBurnes/scoop-bucket`
- **Installation**: `scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket && scoop install burnes/version`
- **Updates**: `scoop update burnes/version`
- **Features**: Automatic updates, dependency management, clean uninstallation

### Version Cleaning Logic
- **Input**: `0.5.2-SNAPSHOT-5bb31e3` or `1.0.0-abc1234`
- **Output**: `0.5.2` or `1.0.0`
- **Pattern**: Removes `-SNAPSHOT-*` and `-{7-8 hex chars}` suffixes
- **Purpose**: Clean, user-friendly installer names without development artifacts

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