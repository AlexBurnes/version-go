#!/bin/bash

# Cross-compilation test for platform detection
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Cross-Compilation Platform Detection Test ==="
echo "Project root: $PROJECT_ROOT"
echo ""

cd "$PROJECT_ROOT"

# Test function
test_platform_binary() {
    local platform=$1
    local arch=$2
    local binary_name=$3
    local expected_platform=$4
    local expected_os=$5
    
    echo "Testing $platform-$arch binary..."
    
    # Build binary
    echo "  Building binary..."
    GOOS=$platform GOARCH=$arch go build -o "test/cross-platform/$binary_name" cmd/version/*.go
    
    # Test platform detection (if binary can run on current platform)
    if [ "$platform" = "linux" ] && [ "$arch" = "amd64" ]; then
        echo "  Testing platform detection..."
        
        PLATFORM=$(./test/cross-platform/$binary_name platform)
        ARCH=$(./test/cross-platform/$binary_name arch)
        OS=$(./test/cross-platform/$binary_name os)
        OS_VERSION=$(./test/cross-platform/$binary_name os_version)
        CPU=$(./test/cross-platform/$binary_name cpu)
        
        echo "    Platform: $PLATFORM (expected: $expected_platform)"
        echo "    Arch: $ARCH (expected: $arch)"
        echo "    OS: $OS (expected: $expected_os)"
        echo "    OS Version: $OS_VERSION"
        echo "    CPU: $CPU"
        
        if [ "$PLATFORM" != "$expected_platform" ]; then
            echo "    ❌ FAIL: Platform mismatch"
            return 1
        fi
        
        if [ "$ARCH" != "$arch" ]; then
            echo "    ❌ FAIL: Architecture mismatch"
            return 1
        fi
        
        echo "    ✅ PASS"
    else
        echo "  Binary built successfully (runtime testing requires target platform/architecture)"
    fi
    
    echo ""
}

# Create test directory
mkdir -p test/cross-platform

# Test different platforms
echo "Building and testing cross-compiled binaries..."
echo ""

test_platform_binary "linux" "amd64" "version-linux-amd64" "linux" "ubuntu"
test_platform_binary "linux" "arm64" "version-linux-arm64" "linux" "ubuntu"
test_platform_binary "windows" "amd64" "version-windows-amd64.exe" "windows" "windows"
test_platform_binary "windows" "arm64" "version-windows-arm64.exe" "windows" "windows"
test_platform_binary "darwin" "amd64" "version-darwin-amd64" "darwin" "darwin"
test_platform_binary "darwin" "arm64" "version-darwin-arm64" "darwin" "darwin"

echo "=== Cross-Compilation Test Results ==="
echo "✅ Successfully built binaries for:"
echo "   - Linux amd64/arm64"
echo "   - Windows amd64/arm64" 
echo "   - Darwin amd64/arm64"
echo ""
echo "Note: Runtime testing requires actual target platforms."
echo "Use Docker tests for comprehensive platform testing."

# Cleanup
echo "Cleaning up test binaries..."
rm -f test/cross-platform/version-*

echo "✅ Cross-compilation test completed"
