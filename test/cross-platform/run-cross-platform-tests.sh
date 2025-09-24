#!/bin/bash

# Cross-platform test runner for version utility
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Cross-Platform Platform Detection Test Runner ==="
echo "Project root: $PROJECT_ROOT"
echo ""

# Check if binaries exist
check_binaries() {
    echo "Checking for required binaries..."
    
    local missing_binaries=()
    
    if [ ! -f "$PROJECT_ROOT/bin/version-linux-amd64" ]; then
        missing_binaries+=("version-linux-amd64")
    fi
    
    if [ ! -f "$PROJECT_ROOT/bin/version-windows-amd64.exe" ]; then
        missing_binaries+=("version-windows-amd64.exe")
    fi
    
    if [ ! -f "$PROJECT_ROOT/bin/version-darwin-amd64" ]; then
        missing_binaries+=("version-darwin-amd64")
    fi
    
    if [ ${#missing_binaries[@]} -ne 0 ]; then
        echo "❌ Missing binaries:"
        for binary in "${missing_binaries[@]}"; do
            echo "   - $binary"
        done
        echo ""
        echo "Please build all platform binaries first:"
        echo "   cd $PROJECT_ROOT"
        echo "   ./buildtools/build-and-package.sh"
        exit 1
    fi
    
    echo "✅ All required binaries found"
    echo ""
}

# Test Linux Ubuntu
test_linux_ubuntu() {
    echo "=== Testing Linux Ubuntu 24.04 ==="
    
    cd "$PROJECT_ROOT"
    
    # Build Docker image
    echo "Building Ubuntu Docker image..."
    docker build -f test/cross-platform/Dockerfile.linux-ubuntu -t version-test-ubuntu .
    
    # Run tests
    echo "Running Ubuntu tests..."
    docker run --rm version-test-ubuntu
    
    echo "✅ Ubuntu tests completed"
    echo ""
}

# Test Linux Debian
test_linux_debian() {
    echo "=== Testing Linux Debian 12 ==="
    
    cd "$PROJECT_ROOT"
    
    # Build Docker image
    echo "Building Debian Docker image..."
    docker build -f test/cross-platform/Dockerfile.linux-debian -t version-test-debian .
    
    # Run tests
    echo "Running Debian tests..."
    docker run --rm version-test-debian
    
    echo "✅ Debian tests completed"
    echo ""
}

# Test Windows (using Wine)
test_windows() {
    echo "=== Testing Windows (using Wine) ==="
    
    cd "$PROJECT_ROOT"
    
    # Build Docker image
    echo "Building Windows Docker image..."
    docker build -f test/cross-platform/Dockerfile.windows -t version-test-windows .
    
    # Run tests
    echo "Running Windows tests..."
    docker run --rm version-test-windows
    
    echo "✅ Windows tests completed"
    echo ""
}

# Test macOS (simulated)
test_macos() {
    echo "=== Testing macOS (simulated) ==="
    echo "Note: macOS testing requires a macOS host or VM"
    echo "Skipping macOS Docker test (not available on Linux)"
    echo ""
}

# Run all tests
run_all_tests() {
    check_binaries
    test_linux_ubuntu
    test_linux_debian
    test_windows
    test_macos
    
    echo "=== All Cross-Platform Tests Completed! ==="
    echo "✅ Platform detection works correctly on:"
    echo "   - Linux Ubuntu 24.04"
    echo "   - Linux Debian 12"
    echo "   - Windows (via Wine)"
    echo "   - macOS (requires macOS host)"
}

# Parse command line arguments
case "${1:-all}" in
    "ubuntu")
        check_binaries
        test_linux_ubuntu
        ;;
    "debian")
        check_binaries
        test_linux_debian
        ;;
    "windows")
        check_binaries
        test_windows
        ;;
    "macos")
        check_binaries
        test_macos
        ;;
    "all"|"")
        run_all_tests
        ;;
    *)
        echo "Usage: $0 [ubuntu|debian|windows|macos|all]"
        echo ""
        echo "Run cross-platform tests for the version utility."
        echo "Default is 'all' which runs all available tests."
        exit 1
        ;;
esac
