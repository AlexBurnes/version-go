#!/usr/bin/env bash

# Script to create self-extracting installers using makeself.sh
# Usage: ./create-makeself-installer.sh <version> <platform> <arch>

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
    
    # Check if makeself.sh exists
    if [[ ! -f "buildtools/makeself.sh" ]]; then
        log_info "Downloading makeself.sh..."
        curl -L -o buildtools/makeself.sh https://raw.githubusercontent.com/megastep/makeself/master/makeself.sh
        chmod +x buildtools/makeself.sh
    fi
    
    # Check if dist directory exists
    if [[ ! -d "dist" ]]; then
        log_error "dist directory not found. Run GoReleaser first."
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Create installer package directory
create_package_dir() {
    local version="$1"
    local platform="$2"
    local arch="$3"
    
    local package_dir="dist/installer-package-${platform}-${arch}"
    local binary_name="version"
    if [[ "$platform" == "windows" ]]; then
        binary_name="version.exe"
    fi
    
    log_info "Creating package directory: $package_dir"
    
    # Create package directory
    mkdir -p "$package_dir"
    
    # Copy binary
    local binary_path="dist/version_${version}_${platform}_${arch}"
    if [[ "$platform" == "windows" ]]; then
        binary_path="${binary_path}.zip"
        # Extract binary from zip
        unzip -q "$binary_path" -d "$package_dir"
    else
        binary_path="${binary_path}.tar.gz"
        # Extract binary from tar.gz
        tar -xzf "$binary_path" -C "$package_dir"
    fi
    
    # Copy install script
    if [[ "$platform" == "linux" ]]; then
        # Use install script without sudo - user should run with sudo if needed
        cat > "$package_dir/install.sh" << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/usr/local/bin}"
SELF=$(readlink -f "$0" || echo "$0")
BASE_DIR=$(dirname "$SELF")

echo "[*] Target bin dir: $APP_DIR"

install_one() {
  local src="$1"
  local dst="$APP_DIR/$(basename "$src")"
  
  # Create target directory if it doesn't exist
  mkdir -p "$APP_DIR"
  
  # Install without sudo - user should run with sudo if needed
  install -m 0755 "$src" "$dst"
  echo "  + $(basename "$src") -> $dst"
}

if [[ -f "$BASE_DIR/version" ]]; then
  install_one "$BASE_DIR/version"
else
  echo "version binary not found near install.sh"; exit 1
fi

echo "[✓] Installation completed successfully"
echo "Please restart your terminal or run: source ~/.bashrc"
EOF
        chmod +x "$package_dir/install.sh"
    elif [[ "$platform" == "darwin" ]]; then
        # Create macOS install script without sudo
        cat > "$package_dir/install.sh" << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/usr/local/bin}"
SELF=$(readlink -f "$0" || echo "$0")
BASE_DIR=$(dirname "$SELF")

echo "[*] Installing version CLI to: $APP_DIR"

install_one() {
  local src="$1"
  local dst="$APP_DIR/$(basename "$src")"
  
  # Create target directory if it doesn't exist
  mkdir -p "$APP_DIR"
  
  # Install without sudo - user should run with sudo if needed
  install -m 0755 "$src" "$dst"
  echo "  + $(basename "$src") -> $dst"
}

if [[ -f "$BASE_DIR/version" ]]; then
  install_one "$BASE_DIR/version"
else
  echo "version binary not found near install.sh"; exit 1
fi

echo "[✓] Installation completed successfully"
echo "Please restart your terminal or run: source ~/.bashrc"
EOF
        chmod +x "$package_dir/install.sh"
    fi
    
    # Copy documentation
    cp LICENSE "$package_dir/"
    cp README.md "$package_dir/"
    
    log_success "Package directory created: $package_dir"
    echo "$package_dir"
}

# Create self-extracting installer
create_installer() {
    local version="$1"
    local platform="$2"
    local arch="$3"
    local package_dir="$4"
    
    local installer_name="version-${version}-${platform}-${arch}-install.sh"
    local installer_path="dist/${installer_name}"
    
    log_info "Creating self-extracting installer: $installer_name"
    
    # Create makeself installer
    if [[ ! -d "$package_dir" ]]; then
        log_error "Package directory does not exist: $package_dir"
        exit 1
    fi
    
    # Change to package directory and run makeself from there
    cd "$package_dir"
    
    # Use the working makeself approach with header
    ../../buildtools/makeself.sh \
        --header "../../packaging/linux/makeself-header.sh" \
        . \
        "../../$installer_path" \
        "version ${version} - git describe CLI" \
        "./install.sh"
    
    cd - > /dev/null
    
    chmod +x "$installer_path"
    
    log_success "Self-extracting installer created: $installer_path"
    echo "$installer_path" >&2
    echo "$installer_path"
}

# Main function
main() {
    local version="${1:-}"
    local platform="${2:-}"
    local arch="${3:-}"
    
    if [[ -z "$version" || -z "$platform" || -z "$arch" ]]; then
        echo "Usage: $0 <version> <platform> <arch>"
        echo "Example: $0 0.5.2 linux amd64"
        exit 1
    fi
    
    log_info "Creating makeself installer for version $version, platform $platform, arch $arch"
    
    check_prerequisites
    
    log_info "Creating package directory..."
    local package_dir
    package_dir=$(create_package_dir "$version" "$platform" "$arch")
    
    log_info "Creating self-extracting installer..."
    local installer_path
    installer_path=$(create_installer "$version" "$platform" "$arch" "$package_dir")
    
    log_success "Installer created successfully: $installer_path"
    log_info "Users can install with: wget -O - https://github.com/AlexBurnes/version-go/releases/download/v${version}/$(basename "$installer_path") | sh"
}

# Run main function with all arguments
main "$@"