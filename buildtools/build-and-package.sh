#!/usr/bin/env bash

# Build script using build-conan.sh + makeself + goreleaser (skip build)
# Strategy: build-conan -> makeself -> goreleaser (skip build, only package/publish)

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
    
    # Try to use built version utility first
    if [[ -f "scripts/version" ]]; then
        version=$(scripts/version version 2>/dev/null || echo "")
        if [[ -n "$version" ]]; then
            echo "$version"
            return
        fi
    fi
    
    # Fallback to VERSION file or git describe
    if [[ -f "VERSION" ]]; then
        version=$(cat VERSION)
    else
        version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    fi
    
    # Clean version (remove v prefix for some operations)
    echo "$version"
}

# Step 1: Build binaries using build-conan.sh
build_binaries() {
    log_info "Step 1: Building binaries using build-conan.sh..."
    
    # Clean previous builds
    rm -rf dist/
    mkdir -p dist/
    
    # Build for all platforms using build-conan.sh
    ./buildtools/build-conan.sh build-all
    
    # Store binaries in a backup location for GoReleaser hooks
    log_info "Storing binaries for GoReleaser hooks..."
    local version=$(get_version)
    local clean_version=$(echo "$version" | sed 's/^v//' | sed 's/-SNAPSHOT-[a-f0-9]*$//' | sed 's/-[a-f0-9]\{7,8\}$//' | sed 's/-dirty$//')
    
    # Create backup directory for binaries
    mkdir -p .goreleaser-binaries
    
    # Copy binaries with GoReleaser naming convention to backup location
    cp bin/version-linux-amd64 ".goreleaser-binaries/version_${clean_version}_linux_amd64"
    cp bin/version-linux-arm64 ".goreleaser-binaries/version_${clean_version}_linux_arm64"
    cp bin/version-darwin-amd64 ".goreleaser-binaries/version_${clean_version}_darwin_amd64"
    cp bin/version-darwin-arm64 ".goreleaser-binaries/version_${clean_version}_darwin_arm64"
    cp bin/version-windows-amd64.exe ".goreleaser-binaries/version_${clean_version}_windows_amd64.exe"
    
    log_success "Binaries stored for GoReleaser hooks"
    
    # Create archives in GoReleaser format for install scripts
    log_info "Creating archives in GoReleaser format..."
    
    # Ensure dist directory exists
    mkdir -p dist/
    
    # Get version for archive naming
    local version=$(get_version)
    local clean_version=$(echo "$version" | sed 's/^v//' | sed 's/-SNAPSHOT-[a-f0-9]*$//' | sed 's/-[a-f0-9]\{7,8\}$//' | sed 's/-dirty$//')
    
    # Create temporary directories for each platform
    local temp_dirs=()
    
    # Linux amd64
    local linux_amd64_dir=$(mktemp -d)
    temp_dirs+=("$linux_amd64_dir")
    cp bin/version-linux-amd64 "$linux_amd64_dir/version"
    cp LICENSE "$linux_amd64_dir/"
    cp README.md "$linux_amd64_dir/"
    tar -czf "dist/version_${clean_version}_linux_amd64.tar.gz" -C "$linux_amd64_dir" .
    
    # Linux arm64
    local linux_arm64_dir=$(mktemp -d)
    temp_dirs+=("$linux_arm64_dir")
    cp bin/version-linux-arm64 "$linux_arm64_dir/version"
    cp LICENSE "$linux_arm64_dir/"
    cp README.md "$linux_arm64_dir/"
    tar -czf "dist/version_${clean_version}_linux_arm64.tar.gz" -C "$linux_arm64_dir" .
    
    # Darwin amd64
    local darwin_amd64_dir=$(mktemp -d)
    temp_dirs+=("$darwin_amd64_dir")
    cp bin/version-darwin-amd64 "$darwin_amd64_dir/version"
    cp LICENSE "$darwin_amd64_dir/"
    cp README.md "$darwin_amd64_dir/"
    tar -czf "dist/version_${clean_version}_darwin_amd64.tar.gz" -C "$darwin_amd64_dir" .
    
    # Darwin arm64
    local darwin_arm64_dir=$(mktemp -d)
    temp_dirs+=("$darwin_arm64_dir")
    cp bin/version-darwin-arm64 "$darwin_arm64_dir/version"
    cp LICENSE "$darwin_arm64_dir/"
    cp README.md "$darwin_arm64_dir/"
    tar -czf "dist/version_${clean_version}_darwin_arm64.tar.gz" -C "$darwin_arm64_dir" .
    
    # Windows amd64
    local windows_amd64_dir=$(mktemp -d)
    temp_dirs+=("$windows_amd64_dir")
    cp bin/version-windows-amd64.exe "$windows_amd64_dir/version.exe"
    cp LICENSE "$windows_amd64_dir/"
    cp README.md "$windows_amd64_dir/"
    local dist_path=$(realpath dist)
    (cd "$windows_amd64_dir" && zip -r "$dist_path/version_${clean_version}_windows_amd64.zip" .)
    
    # Clean up temporary directories
    for dir in "${temp_dirs[@]}"; do
        rm -rf "$dir"
    done
    
    log_success "Archives created in dist/ directory"
}

# Step 2: Create install scripts
create_install_scripts() {
    local version="${1:-}"
    log_info "Step 2: Creating install scripts..."
    
    # Create makeself installers in installers/ directory
    ./buildtools/create-all-installers.sh "$version" "installers"
    
    log_success "Install scripts created"
}

# Step 3: Publish with GoReleaser (builds binaries, then replaces with pre-built ones)
publish_release() {
    local mode="${1:-snapshot}"
    
    log_info "Step 3: Publishing release with GoReleaser..."
    
    case "$mode" in
        "snapshot")
            goreleaser release --snapshot --clean
            ;;
        "release")
            goreleaser release --clean
            ;;
        "dry-run")
            goreleaser release --snapshot --skip=publish --clean
            ;;
        *)
            log_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    log_success "Release published successfully"
}

# Clean build artifacts
clean_build() {
    log_info "Cleaning build artifacts..."
    rm -rf dist/
    rm -rf bin/
    rm -rf installers/
    rm -rf .goreleaser-binaries/
    rm -f .goreleaser-skip-build.yml
    log_success "Build artifacts cleaned"
}

# Show help
show_help() {
    echo "Usage: $0 [MODE]"
    echo ""
    echo "Modes:"
    echo "  snapshot          # Build snapshot with pre-built binaries"
    echo "  release           # Build full release with pre-built binaries"
    echo "  dry-run           # Build dry run with pre-built binaries"
    echo "  clean             # Clean build artifacts"
    echo "  help              # Show this help message"
    echo ""
    echo "Strategy:"
    echo "  1. Build binaries using build-conan.sh (static for darwin/linux)"
    echo "  2. Create makeself install scripts"
    echo "  3. Publish with GoReleaser (skip build, only package/publish)"
    echo ""
    echo "Examples:"
    echo "  $0 snapshot          # Build snapshot with install scripts"
    echo "  $0 release           # Build full release with install scripts"
    echo "  $0 dry-run           # Build dry run with install scripts"
    echo "  $0 clean             # Clean build artifacts"
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
            log_info "Starting build process with build-conan.sh (mode: $mode)..."
            
            # Get version
            local version=$(get_version)
            log_info "Using version: $version"
            
            # Step 1: Build binaries
            build_binaries
            
            # Step 2: Create install scripts
            create_install_scripts "$version"
            
            # Step 3: Publish release (skip build)
            if [[ "$mode" != "dry-run" ]]; then
                publish_release "$mode"
            else
                log_info "Dry-run mode: skipping publish step"
            fi
            
            log_success "Build process completed successfully!"
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