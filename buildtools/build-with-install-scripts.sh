#!/usr/bin/env bash

# Build script that creates install scripts before running GoReleaser release

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

# Load environment variables from .env file
load_environment() {
    if [[ -f ".env" ]]; then
        log_info "Loading environment variables from .env file..."
        # Source the .env file, ignoring comments and empty lines
        set -a  # automatically export all variables
        source .env
        set +a  # disable automatic export
        log_success "Environment variables loaded from .env"
    else
        log_warning "No .env file found, using system environment variables only"
    fi
}

# Setup Go environment
setup_go_environment() {
    log_info "Setting up Go environment..."
    
    # Get Go's GOPATH and add bin directory to PATH
    local go_bin_dir
    go_bin_dir=$(go env GOPATH)/bin
    
    if [[ -d "$go_bin_dir" ]]; then
        export PATH="$go_bin_dir:$PATH"
        log_info "Added Go bin directory to PATH: $go_bin_dir"
    else
        log_warning "Go bin directory not found: $go_bin_dir"
    fi
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

# Step 1: Build binaries using GoReleaser build only
build_binaries() {
    local mode="${1:-snapshot}"
    
    log_info "Step 1: Building binaries using GoReleaser..."
    
    # Use the existing build-goreleaser.sh script but modify it to skip the post-archive hook
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
    
    log_success "Binaries built successfully"
}

# Step 2: Create install scripts
create_install_scripts() {
    local version="${1:-}"
    
    log_info "Step 2: Creating install scripts..."
    
    if [[ -z "$version" ]]; then
        version=$(get_version)
    fi
    
    # Run the post-archive hook to create install scripts
    ./buildtools/post-archive-hook.sh "$version"
    
    log_success "Install scripts created"
}

# Step 3: Run GoReleaser release with existing install scripts
run_goreleaser_release() {
    local mode="${1:-snapshot}"
    
    log_info "Step 3: Running GoReleaser release with existing install scripts..."
    
    # Backup install scripts before GoReleaser cleans the dist directory
    log_info "Backing up install scripts..."
    mkdir -p /tmp/version-install-scripts
    cp dist/version-*-install.sh /tmp/version-install-scripts/ 2>/dev/null || true
    
    # Use goreleaser with --clean flag
    case "$mode" in
        "snapshot")
            goreleaser release --snapshot --skip-publish --clean
            ;;
        "release")
            goreleaser release --clean
            ;;
        "dry-run")
            goreleaser release --snapshot --skip-publish --clean
            ;;
        *)
            log_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    # Restore install scripts to dist directory
    log_info "Restoring install scripts..."
    cp /tmp/version-install-scripts/version-*-install.sh dist/ 2>/dev/null || true
    rm -rf /tmp/version-install-scripts
    
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
    echo "Usage: $0 [MODE]"
    echo ""
    echo "This script builds binaries, creates install scripts, then runs GoReleaser release"
    echo "with the install scripts already present for the extra_files to find."
    echo ""
    echo "Modes:"
    echo "  snapshot            - Build snapshot release (default)"
    echo "  release             - Build full release"
    echo "  dry-run             - Build dry run (no publishing)"
    echo "  clean               - Clean build artifacts"
    echo "  help                - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 snapshot         # Build snapshot with install scripts"
    echo "  $0 release          # Build full release with install scripts"
    echo "  $0 dry-run          # Build dry run with install scripts"
    echo "  $0 clean            # Clean build artifacts"
}

# Main script logic
main() {
    local mode="${1:-snapshot}"
    
    # Change to project root directory
    cd "$(dirname "$0")/.."
    
    # Load environment and setup
    load_environment
    setup_go_environment
    check_prerequisites
    
    case "$mode" in
        "clean")
            clean_build
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        "snapshot"|"release"|"dry-run")
            log_info "Starting build process with install scripts (mode: $mode)..."
            
            # Step 1: Build binaries
            build_binaries "$mode"
            
            # Step 2: Create install scripts
            create_install_scripts
            
            # Step 3: Run GoReleaser release (which will find the install scripts)
            run_goreleaser_release "$mode"
            
            log_success "Build and release process completed successfully!"
            log_info "Install scripts should now be included in the GitHub release"
            ;;
        *)
            log_error "Unknown mode: $mode"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"