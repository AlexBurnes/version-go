# Developer Workflow: Version CLI Utility

## Release Workflow

This document outlines the complete workflow for planning, developing, and releasing new versions of the Version CLI utility using buildfab for unified build management.

## Buildfab Ecosystem

This project is part of the **buildfab utilities and libraries** ecosystem, providing unified build management and development tools:

- **[buildfab](https://github.com/AlexBurnes/buildfab)** - Unified build management utility
- **[pre-push](https://github.com/AlexBurnes/pre-push)** - Git pre-push hook utility for project validation  
- **version** (this project) - Semantic version parsing and validation utility

The buildfab ecosystem provides a complete development workflow with build orchestration, project validation, and cross-platform testing capabilities.

## Buildfab Workflow

The project now uses **buildfab** for unified build management with the following stages:

### Available Buildfab Stages

```bash
# Pre-push validation (version checks, git status)
buildfab pre-push

# Full build (Conan, CMake, GoReleaser dry-run)
buildfab build

# Cross-platform testing (Docker-based)
buildfab test

# Complete release (build + test + GoReleaser release)
buildfab release
```

### Individual Actions

```bash
# Version and validation
buildfab version-check
buildfab version-greatest
buildfab version-module

# Build actions
buildfab check-conan
buildfab check-cmake
buildfab check-goreleaser
buildfab install-conan-deps
buildfab configure-cmake
buildfab build-binaries
buildfab install-binary

# Testing actions
buildfab check-binaries
buildfab test-linux-ubuntu
buildfab test-linux-debian
buildfab test-windows
buildfab test-darwin

# Release actions
buildfab create-installers
buildfab goreleaser-dry-run
buildfab goreleaser-release
```

### Prerequisites

- Go 1.22+ installed
- Git configured with proper remote
- GitHub token in `.env` file
- Conan and CMake for builds
- GoReleaser for releases
- **Buildfab** for unified build management
- **Pre-Push Utility** (recommended for developers)

### Pre-Push Hook Setup (Recommended)

For the best development experience, set up the pre-push hook to automatically run project checks:

#### Install Pre-Push Utility

```bash
# Install pre-push utility (same process as buildfab)
wget -O - https://github.com/AlexBurnes/pre-push/releases/latest/download/pre-push-linux-amd64-install.sh | sudo sh

# Or install locally
wget -O - https://github.com/AlexBurnes/pre-push/releases/latest/download/pre-push-linux-amd64-install.sh | INSTALL_DIR=./scripts sh
```

#### Setup Git Hook

```bash
# In project root, run pre-push to setup git hook
pre-push

# This will:
# - Install git pre-push hook
# - Configure project settings using .project.yml
# - Run pre-push stage checks automatically before each push
```

#### Usage

```bash
# Test pre-push checks without pushing
pre-push test

# Run with verbose output to see what's happening
export PRE_PUSH_VERBOSE=1
pre-push test

# Normal git push (will automatically run checks)
git push origin master
```

The pre-push hook will automatically run the project's `pre-push` stage (version checks, git status validation, etc.) before each push, keeping the project clean and preventing broken commits from being pushed.

### Phase 1: Planning & Development

#### 1. Plan Changes and Update Version

**Update VERSION file:**
```bash
# Option 1: Manual edit
echo "v0.5.8" > VERSION

# Option 2: Use version-bump script
./scripts/version-bump patch    # v0.5.9 -> v0.5.10
./scripts/version-bump minor    # v0.5.9 -> v0.6.0
./scripts/version-bump major    # v0.5.9 -> v1.0.0
```

**Update CHANGELOG.md:**
- Add new features, fixes, or changes
- Follow conventional commit format
- Update version numbers and dates

**Plan changes:**
- Review memory bank documents
- Identify what needs to be implemented
- Update project documentation if needed

#### 2. Development and Testing

**Make code changes:**
```bash
# Make your changes to the codebase
# Update documentation as needed
# Update memory bank documents
```

**Run tests:**
```bash
# Run all tests using buildfab (recommended)
buildfab test

# Run individual buildfab stages
buildfab pre-push    # Version checks and basic validation
buildfab build       # Full build with Conan, CMake, GoReleaser
buildfab test        # Cross-platform testing with Docker

# Run specific test suites manually
go test ./cmd/version/... -v
go test ./pkg/version/... -v
go test ./cmd/version/... -run TestPerformance -v
```

**Update documentation:**
- Update README.md if needed
- Update memory bank documents (activeContext.md, progress.md, etc.)
- Update package files if version numbers changed

### Phase 2: Pre-Release Validation

#### 3. Dry Run Validation

**Clean build artifacts:**
```bash
# Clean using buildfab
buildfab build  # This will clean and rebuild everything

# Or clean manually
rm -rf build/ bin/ dist/
```

**Run build and package dry run:**
```bash
# Use buildfab for comprehensive build testing
buildfab build       # Full build with Conan, CMake, GoReleaser
buildfab test        # Cross-platform testing
```

**Verify outputs:**
- Check that all platforms build successfully
- Verify package manager files are generated correctly
- Ensure no configuration errors
- Check that version numbers match VERSION file

### Phase 3: Release Execution

#### 4. Commit Changes

```bash
# Stage all changes
git add .

# Commit with conventional commit message
git commit -m "feat: add new features for v0.5.8"

# Or for fixes:
git commit -m "fix: resolve issue with version parsing"
```

#### 5. Set Git Tag

```bash
# Create tag from VERSION file
VERSION=$(cat VERSION)
git tag $VERSION

# Verify tag
git tag -l | grep $VERSION
```

#### 6. Upload to Git Repository

```bash
# Push commits
git push origin master

# Push tag
git push origin $VERSION
```

#### 7. Run Build and Package Release

```bash
# Run full release using buildfab
buildfab release

# This will:
# - Run Conan dependency management
# - Configure and build with CMake
# - Build all platform binaries
# - Run cross-platform tests
# - Create installers
# - Run GoReleaser release
```

### Phase 4: Post-Release Verification

#### 8. Verify Release

**Check GitHub release:**
- Visit: https://github.com/AlexBurnes/version-go/releases
- Verify all assets are uploaded
- Check that version matches VERSION file

**Test package managers:**
```bash
# Test Homebrew (macOS)
brew tap AlexBurnes/homebrew-tap
brew install version
version --version

# Test Scoop (Windows)
scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket
scoop install burnes/version
version --version
```

**Test direct downloads:**
- Download binaries from GitHub releases
- Verify they work on target platforms
- Check checksums

### Rollback Procedure

If something goes wrong during release:

```bash
# Delete local tag
git tag -d v0.5.8

# Delete remote tag
git push origin :refs/tags/v0.5.8

# Delete GitHub release manually via web interface
# Fix issues and retry
```

### Version Numbering

Follow semantic versioning (SemVer):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in backwards compatible manner
- **PATCH**: Backwards compatible bug fixes

Examples:
- `v1.0.0` - First stable release
- `v1.1.0` - New features added
- `v1.1.1` - Bug fixes only
- `v2.0.0` - Breaking changes

### Pre-Push Hook Integration

The project includes a pre-push hook that:
- Checks VERSION file exists and is valid
- Verifies version number format
- Ensures version is incremented from last tag
- Prevents accidental releases

### Memory Bank Maintenance

Always update memory bank documents:
- `activeContext.md` - Current work focus and recent changes
- `progress.md` - What works and what's left to build
- `systemPatterns.md` - Technical decisions and patterns
- `productContext.md` - Product goals and user experience
- `techContext.md` - Technologies and constraints

### Troubleshooting

**Common Issues:**

1. **GoReleaser dry run fails:**
   - Check `.goreleaser.yml` syntax
   - Verify all required files exist
   - Check GitHub token permissions

2. **Package manager files incorrect:**
   - Verify URLs in package files
   - Check version numbers match VERSION file
   - Ensure repository names are correct

3. **Tests fail:**
   - Run tests locally first
   - Check for race conditions with `-race` flag
   - Verify all dependencies are available

4. **Release fails:**
   - Check GitHub token has correct permissions
   - Verify repository access
   - Check for existing release with same version

### Best Practices

1. **Always test locally first**
2. **Never skip the dry run**
3. **Update documentation with changes**
4. **Use conventional commit messages**
5. **Verify package managers after release**
6. **Keep memory bank documents current**
7. **Test on multiple platforms when possible**

### Quick Reference

```bash
# Complete workflow
./scripts/version-bump patch
# ... make changes ...
go test ./... -v
./buildtools/build-and-package.sh clean
./buildtools/build-and-package.sh dry-run
git add .
VERSION=$(cat VERSION)
git commit -m "feat: add new features for $VERSION"
git tag $VERSION
git push origin master
git push origin $VERSION
./buildtools/build-and-package.sh release
```

### Helper Scripts

- **`scripts/version-bump [major|minor|patch]`**: Automatically increment version number
- **`scripts/version-check <version_tag>`**: Validate version tag format and increment
- **`buildtools/build-and-package.sh [command]`**: Main build and package management
- **`buildtools/build-goreleaser.sh [command]`**: GoReleaser-specific build management