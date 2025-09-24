#!/bin/bash

# Simple cross-platform test that builds and validates binaries
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Simple Cross-Platform Platform Detection Test ==="
echo "Project root: $PROJECT_ROOT"
echo ""

cd "$PROJECT_ROOT"

# Create test directory
mkdir -p test/cross-platform

# Test function
test_platform_build() {
    local platform=$1
    local arch=$2
    local binary_name=$3
    
    echo "Testing $platform-$arch binary build..."
    
    # Build binary
    echo "  Building binary..."
    if GOOS=$platform GOARCH=$arch go build -o "test/cross-platform/$binary_name" cmd/version/*.go; then
        echo "  ✅ Build successful"
        
        # Check binary exists and is executable
        if [ -f "test/cross-platform/$binary_name" ]; then
            echo "  ✅ Binary created: $binary_name"
            
            # Get binary size
            local size=$(ls -lh "test/cross-platform/$binary_name" | awk '{print $5}')
            echo "  📦 Binary size: $size"
        else
            echo "  ❌ FAIL: Binary not created"
            return 1
        fi
    else
        echo "  ❌ FAIL: Build failed"
        return 1
    fi
    
    echo ""
}

# Test different platforms
echo "Building cross-compiled binaries..."
echo ""

test_platform_build "linux" "amd64" "version-linux-amd64"
test_platform_build "linux" "arm64" "version-linux-arm64"
test_platform_build "windows" "amd64" "version-windows-amd64.exe"
test_platform_build "windows" "arm64" "version-windows-arm64.exe"
test_platform_build "darwin" "amd64" "version-darwin-amd64"
test_platform_build "darwin" "arm64" "version-darwin-arm64"

# Test the current platform binary
echo "Testing current platform binary (Linux AMD64)..."
if [ -f "test/cross-platform/version-linux-amd64" ]; then
    echo "Running platform detection tests..."
    
    PLATFORM=$(./test/cross-platform/version-linux-amd64 platform)
    ARCH=$(./test/cross-platform/version-linux-amd64 arch)
    OS=$(./test/cross-platform/version-linux-amd64 os)
    OS_VERSION=$(./test/cross-platform/version-linux-amd64 os_version)
    CPU=$(./test/cross-platform/version-linux-amd64 cpu)
    
    echo "Results:"
    echo "  Platform: $PLATFORM"
    echo "  Architecture: $ARCH"
    echo "  OS: $OS"
    echo "  OS Version: $OS_VERSION"
    echo "  CPU: $CPU"
    
    # Validate results
    if [ "$PLATFORM" = "linux" ] && [ "$ARCH" = "amd64" ] && [ -n "$OS" ] && [ -n "$OS_VERSION" ] && [ -n "$CPU" ]; then
        echo "  ✅ All platform detection tests passed"
    else
        echo "  ❌ Platform detection validation failed"
        exit 1
    fi
else
    echo "  ❌ Linux AMD64 binary not found"
    exit 1
fi

echo ""
echo "=== Cross-Platform Build Test Results ==="
echo "✅ Successfully built binaries for:"
echo "   - Linux amd64/arm64"
echo "   - Windows amd64/arm64" 
echo "   - Darwin amd64/arm64"
echo ""
echo "✅ Platform detection works on current platform (Linux AMD64)"
echo ""
echo "📝 For comprehensive testing on other platforms:"
echo "   1. Use Docker tests: ./run-cross-platform-tests.sh"
echo "   2. Test on actual target platforms"
echo "   3. Use GitHub Actions for automated testing"

# List all built binaries
echo ""
echo "Built binaries:"
ls -lh test/cross-platform/version-*

echo ""
echo "✅ Simple cross-platform test completed successfully"
