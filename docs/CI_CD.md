# CI/CD Workflow Documentation

## Overview

This project uses GitHub Actions for continuous integration and deployment, with Conan package management for Go toolchain dependencies. The CI/CD pipeline ensures consistent builds across multiple platforms and automated releases.

## Workflow Files

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `master`, `main`, or `develop` branches
- Pull requests to `master`, `main`, or `develop` branches

**Jobs:**
- **Test Job**: Runs on Ubuntu, Windows, and macOS
  - Installs Go 1.23.x
  - Installs Conan 2.0
  - Checks for golang package in Conan repositories
  - Creates golang package locally if not found
  - Installs dependencies with Conan
  - Runs Go tests with race detection
  - Runs golangci-lint for code quality

- **Build Job**: Runs on Ubuntu, Windows, and macOS
  - Similar setup to test job
  - Builds the project using Conan + CMake
  - Uploads build artifacts

- **Package Job**: Runs on Ubuntu only
  - Tests GoReleaser dry-run
  - Uploads package artifacts

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of version tags (`v*.*.*`)

**Jobs:**
- **Release Job**: Runs on Ubuntu
  - Installs Go 1.23.x and Conan 2.0
  - Checks for golang package in Conan repositories
  - Creates golang package locally if not found
  - Installs dependencies with Conan
  - Runs GoReleaser for cross-platform builds and distribution
  - Publishes to GitHub Releases, Scoop, and Homebrew

## Conan Package Management

### Golang Package Recipe

The project includes a local Conan recipe (`conanfile-golang.py`) for Go 1.23.0 that:

- Downloads Go toolchain from official Go website
- Supports Linux, Windows, and macOS
- Supports amd64 and arm64 architectures
- Automatically adds Go binaries to PATH
- Can be created locally when not available in remote repositories

### Package Check Process

Before building, the CI/CD pipeline:

1. **Checks for golang package** in Conan remote repositories
2. **If found**: Uses the remote package
3. **If not found**: Creates the package locally using `conanfile-golang.py`
4. **Installs dependencies** using the available golang package

### Local Development

For local development, you can:

```bash
# Check if golang package is available
./scripts/check-golang-conan

# Or manually create the package
conan create conanfile-golang.py --build=missing
```

## Build Scripts Integration

### Updated Build Scripts

Both `buildtools/build-conan.sh` and `buildtools/build-and-package.sh` now include:

- **Golang package check** before building
- **Automatic package creation** if not available
- **Error handling** for missing package recipe
- **Consistent behavior** across local and CI environments

### Build Process

1. **Check prerequisites** (Go, Conan, Git)
2. **Check for golang package** in Conan
3. **Create golang package locally** if needed
4. **Install dependencies** with Conan
5. **Build project** using CMake + Conan
6. **Package and distribute** using GoReleaser

## Platform Support

The CI/CD pipeline supports:

- **Linux**: Ubuntu (amd64, arm64)
- **Windows**: Windows Server (amd64, arm64)
- **macOS**: macOS (amd64, arm64)

## Dependencies

### Required Tools

- **Go 1.23.x**: Managed via Conan package
- **Conan 2.0**: Package manager for dependencies
- **CMake**: Build system (via Conan)
- **GoReleaser**: Cross-platform build and distribution

### Conan Packages

- **golang/1.23.0**: Go toolchain (created locally if not available)
- **cmake/[>=3.16]**: CMake build system

## Environment Variables

### GitHub Secrets

- `GITHUB_TOKEN`: For GitHub Releases and package manager updates
- `PAT_SCOOP`: (Optional) Personal access token for private Scoop bucket

### Local Development

- `.env` file support for local environment variables
- Automatic loading in build scripts

## Artifacts

### Build Artifacts

- **Binaries**: Platform-specific executables
- **Packages**: Cross-platform distribution packages
- **Checksums**: SHA256 checksums for verification

### Distribution

- **GitHub Releases**: Cross-platform binaries and packages
- **Scoop**: Windows package manager integration
- **Homebrew**: macOS package manager integration
- **Linux**: tar.gz archives with install scripts

## Troubleshooting

### Common Issues

1. **Golang package not found**: The CI will automatically create it locally
2. **Conan profile missing**: Automatically created with `conan profile detect --force`
3. **Build failures**: Check logs for specific error messages
4. **Permission issues**: Ensure proper GitHub token permissions

### Debug Commands

```bash
# Check Conan profile
conan profile show default

# Search for golang package
conan search golang --remote=all

# Create golang package locally
conan create conanfile-golang.py --build=missing

# Test GoReleaser dry-run
goreleaser release --snapshot --clean --skip=publish
```

## Best Practices

1. **Always test locally** before pushing
2. **Use semantic versioning** for tags
3. **Check CI status** before merging PRs
4. **Monitor build artifacts** for correct platform support
5. **Verify distribution** across all package managers

## Future Enhancements

- **Matrix builds** for multiple Go versions
- **Performance testing** in CI
- **Security scanning** integration
- **Dependency updates** automation
- **Build caching** optimization