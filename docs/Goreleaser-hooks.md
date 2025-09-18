# GoReleaser Hooks and Conan Integration

This document explains the GoReleaser hooks configuration and how it integrates with Conan for automated cross-platform builds.

## Overview

The GoReleaser configuration uses Conan hooks to:
1. Install build tools (Go, CMake, etc.) via Conan
2. Deploy tools to a local `.toolchain` directory
3. Normalize the PATH for consistent tool access
4. Enable cross-platform builds from a single environment

## Hook Configuration

### Before Hooks

The `.goreleaser.yml` file includes these `before` hooks:

```yaml
before:
  hooks:
    # 1) Ensure Conan profiles exist (idempotent)
    - bash -lc 'conan profile detect --force'
    # 2) Install tool_requires from conanfile.py and DEPLOY their binaries into ./.toolchain
    - bash -lc 'mkdir -p .toolchain && conan install . -g deploy --deployer-folder=.toolchain --build=missing -s build_type=Release'
    # 3) Normalize PATH: collect all */bin under .toolchain into .toolchain/bin for a stable PATH entry
    - bash -lc './buildtools/collect_bins.sh .toolchain'
    - go mod tidy
```

### Hook Breakdown

#### 1. Conan Profile Detection
```bash
conan profile detect --force
```
- **Purpose**: Ensures Conan has a valid profile for the current system
- **Idempotent**: Safe to run multiple times
- **Output**: Creates or updates the default Conan profile

#### 2. Tool Installation and Deployment
```bash
mkdir -p .toolchain && conan install . -g deploy --deployer-folder=.toolchain --build=missing -s build_type=Release
```
- **Purpose**: Installs build tools specified in `conanfile.py`
- **Deploy Generator**: Copies tool binaries to `.toolchain` directory
- **Build Missing**: Builds tools from source if not in cache
- **Release Build**: Uses Release build type for optimal performance

#### 3. PATH Normalization
```bash
./buildtools/collect_bins.sh .toolchain
```
- **Purpose**: Collects all `*/bin` directories under `.toolchain` into `.toolchain/bin`
- **Result**: Single directory containing all tool executables
- **PATH**: Tools become available via `PATH={{ .Env.PATH }}:{{ .ProjectDir }}/.toolchain/bin`

#### 4. Go Module Tidy
```bash
go mod tidy
```
- **Purpose**: Ensures Go modules are clean and up-to-date
- **Dependencies**: Removes unused dependencies
- **Version**: Updates to latest compatible versions

## Conan Integration

### conanfile.py Configuration

The `conanfile.py` specifies build requirements:

```python
def build_requirements(self):
    # Go as build requirement
    self.tool_requires("golang/1.21.0")
    
    # CMake as build requirement
    self.tool_requires("cmake/[>=3.16]")
```

### Tool Deployment

The deploy generator copies tool binaries to the specified folder:

```python
def generate(self):
    tc = CMakeToolchain(self)
    tc.generate()
    
    deps = CMakeDeps(self)
    deps.generate()
```

### Build Environment

The GoReleaser build configuration exposes the toolchain:

```yaml
builds:
  - id: version
    main: ./cmd/version
    binary: version
    env:
      # Expose Conan-deployed tools (go, cmake, ninja, etc.) to the GoReleaser build
      - PATH={{ .Env.PATH }}:{{ .ProjectDir }}/.toolchain/bin
      - CGO_ENABLED=0
    goos: [linux, windows, darwin]
    goarch: [amd64, arm64]
```

## collect_bins.sh Script

The `collect_bins.sh` script normalizes the toolchain structure:

```bash
#!/usr/bin/env bash
# Collect all */bin directories under .toolchain into .toolchain/bin

set -euo pipefail

TOOLCHAIN_DIR="${1:-.toolchain}"

if [ ! -d "$TOOLCHAIN_DIR" ]; then
    echo "Error: Toolchain directory $TOOLCHAIN_DIR does not exist"
    exit 1
fi

# Create target bin directory
mkdir -p "$TOOLCHAIN_DIR/bin"

# Find all bin directories and copy executables
find "$TOOLCHAIN_DIR" -name "bin" -type d | while read -r bin_dir; do
    if [ "$bin_dir" != "$TOOLCHAIN_DIR/bin" ]; then
        echo "Copying executables from $bin_dir to $TOOLCHAIN_DIR/bin"
        cp -f "$bin_dir"/* "$TOOLCHAIN_DIR/bin/" 2>/dev/null || true
    fi
done

echo "Toolchain normalized: $TOOLCHAIN_DIR/bin"
```

## Build Process Flow

1. **Profile Detection**: Conan detects system profile
2. **Tool Installation**: Conan installs Go, CMake, and other tools
3. **Tool Deployment**: Tools are copied to `.toolchain` directory
4. **PATH Normalization**: All tool executables are collected in `.toolchain/bin`
5. **Go Module Tidy**: Go modules are cleaned and updated
6. **Cross-Platform Build**: GoReleaser builds for all target platforms
7. **Package Creation**: Archives and packages are created for each platform

## Benefits

### Reproducible Builds
- **Consistent Tools**: Same tool versions across all environments
- **Isolated Environment**: Tools don't interfere with system installations
- **Cacheable**: Conan caches tools for faster subsequent builds

### Cross-Platform Support
- **Single Environment**: Build all platforms from one machine
- **Tool Management**: Conan handles tool installation and versioning
- **Path Isolation**: Tools are isolated in `.toolchain` directory

### CI/CD Integration
- **Docker Friendly**: Works well in containerized environments
- **Cache Efficient**: Conan caches can be shared across builds
- **Tool Updates**: Easy to update tool versions via `conanfile.py`

## Troubleshooting

### Common Issues

1. **Conan Profile Missing**
   ```bash
   conan profile detect --force
   ```

2. **Tool Installation Fails**
   ```bash
   conan install . -g deploy --deployer-folder=.toolchain --build=missing -v
   ```

3. **PATH Issues**
   ```bash
   # Check toolchain structure
   ls -la .toolchain/bin/
   
   # Verify PATH
   echo $PATH
   ```

4. **Go Module Issues**
   ```bash
   go mod tidy
   go mod download
   ```

### Debug Commands

```bash
# Check Conan profile
conan profile show default

# List installed packages
conan list "*"

# Check toolchain structure
find .toolchain -name "bin" -type d

# Verify tool availability
.toolchain/bin/go version
.toolchain/bin/cmake --version
```

## Configuration Files

### .goreleaser.yml
- GoReleaser configuration with Conan hooks
- Cross-platform build settings
- Archive and package creation

### conanfile.py
- Conan package definition
- Build requirements specification
- Tool deployment configuration

### collect_bins.sh
- PATH normalization script
- Tool collection utility
- Build environment setup

## Integration with Build Scripts

The hooks work with both build scripts:

- **build-conan.sh**: Local development builds
- **build-goreleaser.sh**: Packaging builds with GoReleaser

Both scripts use the same Conan integration patterns for consistency.