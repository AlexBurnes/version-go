# Cross-Platform Platform Detection Tests

This directory contains comprehensive cross-platform testing for the version utility's platform detection functionality.

## Overview

The platform detection feature provides commands to detect:
- `platform` - Current platform (GOOS value)
- `arch` - Current architecture (GOARCH value) 
- `os` - Current operating system (distribution name)
- `os_version` - Current OS version
- `cpu` - Number of logical CPUs

## Test Strategies

### 1. Docker-based Testing (Recommended)

Uses Docker containers to test on different Linux distributions and simulate Windows using Wine.

#### Prerequisites
- Docker installed
- All platform binaries built (`version-linux-amd64`, `version-windows-amd64.exe`, `version-darwin-amd64`)

#### Running Tests

```bash
# Run all platform tests using buildfab (recommended)
buildfab test

# Run individual platform tests using buildfab
buildfab test-linux-ubuntu
buildfab test-linux-debian
buildfab test-windows
buildfab test-macos

# Legacy script-based testing (still available)
./test/cross-platform/run-cross-platform-tests.sh
```

#### Supported Platforms
- **Ubuntu 24.04** - Tests Linux distribution detection
- **Debian 12** - Tests Linux distribution detection
- **Windows** - Tests Windows platform detection using Wine

### 2. Cross-Compilation Testing

Cross-compilation testing has been removed from the CI/CD pipeline. The project now focuses on cross-platform testing using Docker containers for better reliability and consistency.

### 3. GitHub Actions Testing

Automated testing on GitHub Actions runners:
- Ubuntu latest
- Windows latest  
- macOS latest

## Expected Results

### Linux (Ubuntu/Debian)
```
Platform: linux
Architecture: amd64
OS: ubuntu (or debian)
OS Version: 24.04 (or 12)
CPU: 8
```

### Windows
```
Platform: windows
Architecture: amd64
OS: windows
OS Version: windows
CPU: 8
```

### macOS
```
Platform: darwin
Architecture: amd64 (or arm64)
OS: darwin
OS Version: darwin
CPU: 8
```

## Test Files

### Dockerfiles
- `Dockerfile.linux-ubuntu` - Ubuntu 24.04 test environment
- `Dockerfile.linux-debian` - Debian 12 test environment
- `Dockerfile.windows` - Windows test environment using Wine
- `Dockerfile.macos` - macOS test environment (requires macOS host)

### Test Scripts
- `test-platform.sh` - Linux platform detection tests
- `test-platform-windows.sh` - Windows platform detection tests
- `test-platform-macos.sh` - macOS platform detection tests
- `run-cross-platform-tests.sh` - Main test runner

## Building Test Binaries

Before running tests, ensure all platform binaries are built:

```bash
# Build all platform binaries
./buildtools/build-and-package.sh

# Or build individually
go build -o bin/version-linux-amd64 cmd/version/*.go
GOOS=windows GOARCH=amd64 go build -o bin/version-windows-amd64.exe cmd/version/*.go
GOOS=darwin GOARCH=amd64 go build -o bin/version-darwin-amd64 cmd/version/*.go
```

## Manual Testing

For manual testing on different platforms:

```bash
# Test platform detection
./bin/version platform
./bin/version arch
./bin/version os
./bin/version os_version
./bin/version cpu

# Test comprehensive platform info
./bin/version --help | grep -E "(platform|arch|os|cpu)"
```

## Troubleshooting

### Docker Tests Fail
- Ensure Docker is running
- Check that all required binaries exist in `bin/` directory
- Verify Docker has sufficient resources

### Cross-Compilation Fails
- Ensure Go toolchain is properly installed
- Check that target platforms are supported
- Verify no CGO dependencies (should use pure Go)

### Platform Detection Incorrect
- Check `/etc/os-release` file on Linux systems
- Verify runtime.GOOS and runtime.GOARCH values
- Test on actual target platform if possible

## CI/CD Integration

The cross-platform tests are integrated into GitHub Actions:

```yaml
# Runs on every push/PR
name: Cross-Platform Platform Detection Tests
```

Tests run on:
- Ubuntu latest
- Windows latest
- macOS latest

## Contributing

When adding new platform detection features:

1. Add tests to appropriate test scripts
2. Update Dockerfiles if new dependencies needed
3. Update expected results in this README
4. Ensure GitHub Actions tests pass

## Limitations

- **macOS Testing**: Requires macOS host or VM (Docker on Linux cannot run macOS)
- **Windows Testing**: Uses Wine simulation, may not catch all Windows-specific issues
- **ARM Testing**: Limited ARM platform availability in CI/CD

For production deployment, test on actual target platforms when possible.
