# Build Guide: Version CLI Utility

This document provides comprehensive instructions for building the `version` CLI utility using both local development builds and automated packaging builds.

## Overview

The project supports two main build approaches:

1. **Local Development Build**: Using Conan + CMake + bash scripts for development and testing
2. **Packaging Build**: Using GoReleaser with Conan hooks for automated cross-platform packaging

## Prerequisites

### Required Tools
- **Conan 2.x**: Package manager for dependencies
- **CMake 3.16+**: Build system orchestrator
- **Go 1.21+**: Go compiler (managed via Conan)
- **Git**: Version control and tag integration

### Optional Tools
- **golangci-lint**: Code linting (installed via Conan)
- **Makeself**: Linux self-extracting archives (for packaging)
- **Scoop**: Windows package manager (for packaging)

## Local Development Build

### Quick Start

```bash
# Navigate to buildtools directory
cd buildtools/

# Install dependencies and build
./build-conan.sh deps
./build-conan.sh build

# Run the application
./build-conan.sh run
```

### Detailed Build Process

#### 1. Install Dependencies

```bash
cd buildtools/
./build-conan.sh deps
```

This command:
- Installs Go 1.21.0 via Conan
- Installs CMake via Conan
- Sets up the build environment
- Creates necessary build directories

#### 2. Configure Build

```bash
./build-conan.sh configure
```

This command:
- Configures CMake with Conan toolchain
- Sets up cross-compilation targets
- Prepares build environment

#### 3. Build Binary

```bash
# Build for current platform
./build-conan.sh build

# Build for specific platform
./build-conan.sh build version-linux-amd64

# Build for all platforms
./build-conan.sh build-all
```

#### 4. Install Binary

```bash
# Install to system (requires sudo)
./build-conan.sh install /usr/local

# Install to project bin directory
./build-conan.sh install-current
```

### Available Build Targets

| Target | Description |
|--------|-------------|
| `version` | Build for current platform (static for Linux) |
| `version-all` | Build for all supported platforms |
| `version-linux-amd64` | Build for Linux/amd64 |
| `version-linux-amd64-static` | Build static binary for Linux/amd64 |
| `version-darwin-amd64` | Build for macOS/amd64 |
| `version-darwin-arm64` | Build for macOS/arm64 |
| `version-windows-amd64` | Build for Windows/amd64 |

### Build Commands

| Command | Description |
|---------|-------------|
| `deps` | Install dependencies with Conan |
| `build [target]` | Build binary using CMake + Conan |
| `build-conan` | Build using Conan directly |
| `build-all` | Build for all platforms |
| `install [path]` | Install binary (default: /usr/local) |
| `install-current` | Install current OS binary to project bin directory |
| `install-all` | Install all platform binaries to project bin directory |
| `clean [--clean-cache]` | Clean build artifacts |
| `configure` | Configure CMake with Conan |
| `profile` | Show Conan profile info |
| `packages` | Show available Conan packages |

## Packaging Build (GoReleaser)

### Overview

The packaging build uses GoReleaser with Conan hooks to:
1. Install build tools via Conan
2. Build cross-platform binaries
3. Create distribution packages
4. Prepare installation scripts

### Build Process

#### 1. Prerequisites

Ensure you have:
- GoReleaser installed (`go install github.com/goreleaser/goreleaser@latest`)
- Conan 2.x installed
- Git repository with proper tags

#### 2. Run GoReleaser

```bash
# Dry run (test configuration)
goreleaser release --snapshot --skip-publish

# Full release (requires git tag)
goreleaser release
```

#### 3. GoReleaser Hooks

The `.goreleaser.yml` configuration includes these hooks:

```yaml
before:
  hooks:
    # 1) Ensure Conan profiles exist
    - bash -lc 'conan profile detect --force'
    # 2) Install tool_requires and deploy binaries
    - bash -lc 'mkdir -p .toolchain && conan install . -g deploy --deployer-folder=.toolchain --build=missing -s build_type=Release'
    # 3) Normalize PATH for tools
    - bash -lc './buildtools/collect_bins.sh .toolchain'
    - go mod tidy
```

### Build Script for GoReleaser

Create `buildtools/build-goreleaser.sh`:

```bash
#!/usr/bin/env bash

# Build script for GoReleaser packaging builds

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check GoReleaser
    if ! command -v goreleaser &> /dev/null; then
        log_error "GoReleaser is not installed. Install with: go install github.com/goreleaser/goreleaser@latest"
        exit 1
    fi
    
    # Check Conan
    if ! command -v conan &> /dev/null; then
        log_error "Conan is not installed. Install with: pip install conan"
        exit 1
    fi
    
    # Check Git
    if ! command -v git &> /dev/null; then
        log_error "Git is not installed"
        exit 1
    fi
    
    log_success "All prerequisites found"
}

# Setup build environment
setup_environment() {
    log_info "Setting up build environment..."
    
    # Create toolchain directory
    mkdir -p .toolchain
    
    # Detect Conan profile
    conan profile detect --force
    
    # Install tools via Conan
    conan install . -g deploy --deployer-folder=.toolchain --build=missing -s build_type=Release
    
    # Normalize PATH
    ./buildtools/collect_bins.sh .toolchain
    
    # Tidy Go modules
    go mod tidy
    
    log_success "Build environment ready"
}

# Run GoReleaser
run_goreleaser() {
    local mode="${1:-snapshot}"
    
    case "$mode" in
        "snapshot")
            log_info "Running GoReleaser snapshot build..."
            goreleaser release --snapshot --skip-publish
            ;;
        "release")
            log_info "Running GoReleaser release build..."
            goreleaser release
            ;;
        "dry-run")
            log_info "Running GoReleaser dry run..."
            goreleaser release --snapshot --skip-publish --rm-dist
            ;;
        *)
            log_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    log_success "GoReleaser build completed"
}

# Clean build artifacts
clean_build() {
    log_info "Cleaning build artifacts..."
    
    # Remove toolchain directory
    rm -rf .toolchain
    
    # Remove dist directory
    rm -rf dist/
    
    # Remove build artifacts
    rm -rf bin/
    
    log_success "Build artifacts cleaned"
}

# Show help
show_help() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  setup                 - Setup build environment"
    echo "  snapshot              - Run snapshot build (no publish)"
    echo "  release               - Run full release build"
    echo "  dry-run               - Run dry run build"
    echo "  clean                 - Clean build artifacts"
    echo "  help                  - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 setup              # Setup build environment"
    echo "  $0 snapshot           # Build snapshot"
    echo "  $0 release            # Build and publish release"
    echo "  $0 dry-run            # Test build without publishing"
}

# Main script logic
main() {
    local command="${1:-help}"
    
    case "$command" in
        "setup")
            check_prerequisites
            setup_environment
            ;;
        "snapshot")
            check_prerequisites
            setup_environment
            run_goreleaser "snapshot"
            ;;
        "release")
            check_prerequisites
            setup_environment
            run_goreleaser "release"
            ;;
        "dry-run")
            check_prerequisites
            setup_environment
            run_goreleaser "dry-run"
            ;;
        "clean")
            clean_build
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
```

## Build Artifacts

### Local Build Outputs

- **Binary Location**: `bin/` directory
- **Platform Binaries**: 
  - `version-linux-amd64` - Linux/amd64
  - `version-linux-amd64-static` - Static Linux/amd64
  - `version-darwin-amd64` - macOS/amd64
  - `version-darwin-arm64` - macOS/arm64
  - `version-windows-amd64.exe` - Windows/amd64

### Packaging Build Outputs

- **Linux**: `dist/version_<version>_linux_<arch>.tar.gz`
- **Windows**: `dist/version_<version>_windows_<arch>.zip`
- **macOS**: `dist/version_<version>_darwin_<arch>.tar.gz`
- **Checksums**: `dist/checksums.txt`

## Troubleshooting

### Common Issues

1. **Conan Profile Issues**
   ```bash
   conan profile detect --force
   ```

2. **Go Version Mismatch**
   ```bash
   conan install . --build=missing
   ```

3. **CMake Configuration Issues**
   ```bash
   rm -rf .build/
   ./build-conan.sh configure
   ```

4. **GoReleaser Hooks Fail**
   ```bash
   # Check Conan installation
   conan --version
   
   # Check GoReleaser installation
   goreleaser --version
   
   # Run hooks manually
   conan profile detect --force
   conan install . -g deploy --deployer-folder=.toolchain --build=missing
   ```

### Build Flags

- **Static Build**: `CGO_ENABLED=0` for maximum compatibility
- **Optimization**: `-s -w` ldflags for smaller binaries
- **Version Info**: Embedded via ldflags from git tags

## Development Workflow

1. **Local Development**:
   ```bash
   cd buildtools/
   ./build-conan.sh deps
   ./build-conan.sh build
   ./build-conan.sh run
   ```

2. **Testing**:
   ```bash
   ./build-conan.sh test
   ./build-conan.sh lint
   ```

3. **Packaging**:
   ```bash
   ./build-goreleaser.sh snapshot
   ```

4. **Release**:
   ```bash
   git tag v1.0.0
   ./build-goreleaser.sh release
   ```

## Integration with CI/CD

The build system is designed to work with:
- **GitHub Actions**: Use GoReleaser for automated releases
- **GitLab CI**: Use Conan + CMake for local builds
- **Jenkins**: Use build scripts for custom pipelines
- **Local Development**: Use build-conan.sh for development

## Performance Considerations

- **Static Binaries**: Linux builds use static linking for maximum compatibility
- **Cross-compilation**: All platforms built from single environment
- **Dependency Management**: Conan handles tool versions and dependencies
- **Build Caching**: Conan caches dependencies for faster subsequent builds