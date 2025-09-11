#!/usr/bin/env bash

# Build and release script that creates install scripts before running GoReleaser

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

# Get version from VERSION file or git
get_version() {
    local version=""
    if [[ -f "VERSION" ]]; then
        version=$(cat VERSION)
    else
        # Try to get version from git tag
        if command -v git >/dev/null 2>&1; then
            version=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
        else
            version="unknown"
        fi
    fi
    echo "$version"
}

# Build binaries using GoReleaser build only
build_binaries() {
    local mode="${1:-snapshot}"
    
    log_info "Building binaries using GoReleaser..."
    
    # Use goreleaser build command directly
    case "$mode" in
        "snapshot")
            goreleaser build --snapshot --clean
            ;;
        "release")
            goreleaser build --clean
            ;;
        "dry-run")
            goreleaser build --snapshot --clean
            ;;
        *)
            log_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    log_success "Binaries built successfully"
}

# Create install scripts
create_install_scripts() {
    local version="${1:-}"
    
    log_info "Creating install scripts..."
    
    if [[ -z "$version" ]]; then
        version=$(get_version)
    fi
    
    # Run the post-archive hook to create install scripts
    ./buildtools/post-archive-hook.sh "$version"
    
    log_success "Install scripts created"
}

# Run GoReleaser release with existing install scripts
run_goreleaser_release() {
    local mode="${1:-snapshot}"
    
    log_info "Running GoReleaser release with existing install scripts..."
    
    # Use the existing build-goreleaser.sh script but skip the post-archive hook
    # by temporarily renaming it
    if [[ -f "./buildtools/post-archive-hook.sh" ]]; then
        mv ./buildtools/post-archive-hook.sh ./buildtools/post-archive-hook.sh.bak
    fi
    
    case "$mode" in
        "snapshot")
            ./buildtools/build-goreleaser.sh snapshot
            ;;
        "release")
            ./buildtools/build-goreleaser.sh release
            ;;
        "dry-run")
            ./buildtools/build-goreleaser.sh dry-run
            ;;
        *)
            log_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    # Restore the post-archive hook
    if [[ -f "./buildtools/post-archive-hook.sh.bak" ]]; then
        mv ./buildtools/post-archive-hook.sh.bak ./buildtools/post-archive-hook.sh
    fi
    
    log_success "GoReleaser release completed"
}

# Clean build artifacts
clean_build() {
    log_info "Cleaning build artifacts..."
    rm -rf dist/
    log_success "Build artifacts cleaned"
}

# Show help
show_help() {
    echo "Usage: $0 [COMMAND] [MODE]"
    echo ""
    echo "Commands:"
    echo "  build-and-release   - Build binaries, create install scripts, then run GoReleaser release"
    echo "  build-only          - Build binaries only"
    echo "  install-scripts     - Create install scripts only"
    echo "  release-only        - Run GoReleaser release only (assumes binaries and install scripts exist)"
    echo "  clean               - Clean build artifacts"
    echo "  help                - Show this help"
    echo ""
    echo "Modes:"
    echo "  snapshot            - Build snapshot release (default)"
    echo "  release             - Build full release"
    echo "  dry-run             - Build dry run (no publishing)"
    echo ""
    echo "Examples:"
    echo "  $0 build-and-release snapshot    # Build snapshot with install scripts"
    echo "  $0 build-and-release release     # Build full release with install scripts"
    echo "  $0 build-only release            # Build binaries only for release"
    echo "  $0 install-scripts               # Create install scripts only"
    echo "  $0 release-only snapshot         # Run GoReleaser release only"
}

# Main script logic
main() {
    local command="${1:-help}"
    local mode="${2:-snapshot}"
    
    # Change to project root directory
    cd "$(dirname "$0")/.."
    
    case "$command" in
        "build-and-release")
            log_info "Starting build and release process (mode: $mode)..."
            clean_build
            build_binaries "$mode"
            create_install_scripts
            run_goreleaser_release "$mode"
            log_success "Build and release process completed successfully!"
            ;;
        "build-only")
            log_info "Building binaries only (mode: $mode)..."
            clean_build
            build_binaries "$mode"
            log_success "Binary build completed!"
            ;;
        "install-scripts")
            log_info "Creating install scripts only..."
            create_install_scripts
            log_success "Install scripts created!"
            ;;
        "release-only")
            log_info "Running GoReleaser release only (mode: $mode)..."
            run_goreleaser_release "$mode"
            log_success "GoReleaser release completed!"
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