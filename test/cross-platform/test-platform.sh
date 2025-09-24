#!/bin/bash

# Cross-platform test script for Linux distributions
set -e

echo "=== Cross-Platform Platform Detection Test ==="
echo "Testing on: $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
echo ""

# Test platform detection commands
echo "Testing platform detection commands..."

# Test platform command
echo "1. Testing 'platform' command:"
PLATFORM=$(./version platform)
echo "   Result: $PLATFORM"
if [ "$PLATFORM" != "linux" ]; then
    echo "   ❌ FAIL: Expected 'linux', got '$PLATFORM'"
    exit 1
fi
echo "   ✅ PASS"

# Test arch command
echo "2. Testing 'arch' command:"
ARCH=$(./version arch)
echo "   Result: $ARCH"
if [ "$ARCH" != "amd64" ] && [ "$ARCH" != "arm64" ]; then
    echo "   ❌ FAIL: Expected 'amd64' or 'arm64', got '$ARCH'"
    exit 1
fi
echo "   ✅ PASS"

# Test os command
echo "3. Testing 'os' command:"
OS=$(./version os)
echo "   Result: $OS"
if [ "$OS" != "ubuntu" ] && [ "$OS" != "debian" ] && [ "$OS" != "linux" ]; then
    echo "   ❌ FAIL: Expected 'ubuntu', 'debian', or 'linux', got '$OS'"
    exit 1
fi
echo "   ✅ PASS"

# Test os_version command
echo "4. Testing 'os_version' command:"
OS_VERSION=$(./version os_version)
echo "   Result: $OS_VERSION"
if [ -z "$OS_VERSION" ]; then
    echo "   ❌ FAIL: OS version is empty"
    exit 1
fi
echo "   ✅ PASS"

# Test cpu command
echo "5. Testing 'cpu' command:"
CPU=$(./version cpu)
echo "   Result: $CPU"
if ! [[ "$CPU" =~ ^[0-9]+$ ]] || [ "$CPU" -lt 1 ]; then
    echo "   ❌ FAIL: Expected positive integer, got '$CPU'"
    exit 1
fi
echo "   ✅ PASS"

echo ""
echo "=== All Platform Detection Tests Passed! ==="
echo "Platform: $PLATFORM"
echo "Architecture: $ARCH"
echo "OS: $OS"
echo "OS Version: $OS_VERSION"
echo "CPU Cores: $CPU"
