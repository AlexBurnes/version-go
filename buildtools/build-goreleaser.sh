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