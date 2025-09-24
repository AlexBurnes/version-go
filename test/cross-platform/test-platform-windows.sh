#!/bin/bash

# Cross-platform test script for Windows (using Wine)
set -e

echo "=== Cross-Platform Platform Detection Test (Windows) ==="
echo "Testing Windows platform detection using Wine"
echo ""

# Test platform detection commands using Wine
echo "Testing platform detection commands..."

# Test platform command
echo "1. Testing 'platform' command:"
PLATFORM=$(wine version.exe platform 2>/dev/null | tr -d '\r')
echo "   Result: $PLATFORM"
if [ "$PLATFORM" != "windows" ]; then
    echo "   ❌ FAIL: Expected 'windows', got '$PLATFORM'"
    exit 1
fi
echo "   ✅ PASS"

# Test arch command
echo "2. Testing 'arch' command:"
ARCH=$(wine version.exe arch 2>/dev/null | tr -d '\r')
echo "   Result: $ARCH"
if [ "$ARCH" != "amd64" ] && [ "$ARCH" != "386" ]; then
    echo "   ❌ FAIL: Expected 'amd64' or '386', got '$ARCH'"
    exit 1
fi
echo "   ✅ PASS"

# Test os command
echo "3. Testing 'os' command:"
OS=$(wine version.exe os 2>/dev/null | tr -d '\r')
echo "   Result: $OS"
if [ "$OS" != "windows" ]; then
    echo "   ❌ FAIL: Expected 'windows', got '$OS'"
    exit 1
fi
echo "   ✅ PASS"

# Test os_version command
echo "4. Testing 'os_version' command:"
OS_VERSION=$(wine version.exe os_version 2>/dev/null | tr -d '\r')
echo "   Result: $OS_VERSION"
if [ -z "$OS_VERSION" ]; then
    echo "   ❌ FAIL: OS version is empty"
    exit 1
fi
echo "   ✅ PASS"

# Test cpu command
echo "5. Testing 'cpu' command:"
CPU=$(wine version.exe cpu 2>/dev/null | tr -d '\r')
echo "   Result: $CPU"
if ! [[ "$CPU" =~ ^[0-9]+$ ]] || [ "$CPU" -lt 1 ]; then
    echo "   ❌ FAIL: Expected positive integer, got '$CPU'"
    exit 1
fi
echo "   ✅ PASS"

echo ""
echo "=== All Windows Platform Detection Tests Passed! ==="
echo "Platform: $PLATFORM"
echo "Architecture: $ARCH"
echo "OS: $OS"
echo "OS Version: $OS_VERSION"
echo "CPU Cores: $CPU"
