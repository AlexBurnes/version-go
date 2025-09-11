#!/usr/bin/env bash

# Script to create all makeself installers for all platforms and architectures
# This should be run after GoReleaser has created the distribution archives

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

# Get version from git tag or environment
get_version() {
    local version="${1:-}"
    if [[ -z "$version" ]]; then
        # Try to get version from git tag
        if command -v git >/dev/null 2>&1; then
            version=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
        else
            version="unknown"
        fi
    fi
    echo "$version"
}

# Create installers for all platforms
create_all_installers() {
    local version="$1"
    
    log_info "Creating makeself installers for version: $version"
    
    # Platforms and architectures to build
    local platforms=("linux" "darwin")
    local arches=("amd64" "arm64")
    
    local created_installers=()
    
    for platform in "${platforms[@]}"; do
        for arch in "${arches[@]}"; do
            log_info "Creating installer for $platform-$arch"
            
            # Check if the source archive exists
            local archive_pattern="version_${version}_${platform}_${arch}"
            local archive_file=""
            
            if [[ "$platform" == "windows" ]]; then
                archive_file="${archive_pattern}.zip"
            else
                archive_file="${archive_pattern}.tar.gz"
            fi
            
            if [[ -f "dist/${archive_file}" ]]; then
                ./buildtools/create-makeself-installer.sh "$version" "$platform" "$arch"
                # Clean version for installer name
                local clean_version=$(echo "$version" | sed 's/-SNAPSHOT-[a-f0-9]*$//' | sed 's/-[a-f0-9]\{7,8\}$//')
                created_installers+=("version-${clean_version}-${platform}-${arch}-install.sh")
            else
                log_warning "Source archive not found: dist/${archive_file}"
            fi
        done
    done
    
    log_success "Created ${#created_installers[@]} installers:"
    for installer in "${created_installers[@]}"; do
        echo "  - dist/$installer"
    done
}

# Main function
main() {
    local version="${1:-}"
    version=$(get_version "$version")
    
    log_info "Creating makeself installers for version: $version"
    
    # Check if dist directory exists
    if [[ ! -d "dist" ]]; then
        log_error "dist directory not found. Run GoReleaser first."
        exit 1
    fi
    
    # Check if create-makeself-installer.sh exists
    if [[ ! -f "buildtools/create-makeself-installer.sh" ]]; then
        log_error "create-makeself-installer.sh not found"
        exit 1
    fi
    
    create_all_installers "$version"
    
    log_success "All installers created successfully!"
    log_info "Installers can be distributed via GitHub releases"
}

# Run main function with all arguments
main "$@"