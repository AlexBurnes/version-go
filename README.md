# Version CLI Utility

A cross-platform command-line utility written in Go that provides semantic version parsing, validation, and ordering. This tool replaces legacy bash scripts (`version`, `version-check`, `describe`) currently used in build pipelines and supports Linux, Windows, and macOS with a reproducible build/distribution process.

**Version 1.3.0** - New git integration feature with GetVersion() function for library consumers!

## Buildfab Ecosystem

This project is part of the **buildfab utilities and libraries** ecosystem, providing unified build management and development tools:

- **[buildfab](https://github.com/AlexBurnes/buildfab)** - Unified build management utility
- **[pre-push](https://github.com/AlexBurnes/pre-push)** - Git pre-push hook utility for project validation
- **version** (this project) - Semantic version parsing and validation utility

The buildfab ecosystem provides a complete development workflow with build orchestration, project validation, and cross-platform testing capabilities.

## Features

- **Semantic Version Parsing**: Custom BNF grammar engine supporting extended version formats beyond SemVer 2.0 (see [BNF Grammar](docs/BNF-grammar.md))
- **Version Validation**: Validate version strings with detailed error reporting
- **Version Sorting**: Sort version lists according to defined precedence rules
- **Version Bumping**: Intelligent version incrementing with smart mode and specific bump types
- **Git Integration**: Extract version information from git tags and remotes
- **Cross-Platform**: Static binaries for Linux, Windows, and macOS
- **Colored Output**: Terminal-friendly colored output with `--no-color` support
- **CMake Integration**: Build system integration with CMake and Conan
- **Library Package**: Reusable Go library (`pkg/version`) for other utilities

## Supported Version Formats

The tool supports an extended grammar beyond SemVer 2.0 as defined in the [BNF Grammar Specification](docs/BNF-grammar.md):

### Release Versions
- `1.2.3` - Standard semantic version
- `v1.2.3` - Version with 'v' prefix

### Prerelease Versions
- `1.2.3-alpha` - Alpha prerelease
- `1.2.3~beta.1` - Beta prerelease with tilde separator
- `1.2.3-rc.1` - Release candidate

### Postrelease Versions
- `1.2.3.fix` - Fix postrelease
- `1.2.3.post.1` - Post release with version
- `1.2.3.next` - Next release

### Intermediate Versions
- `1.2.3_feature` - Feature intermediate release
- `1.2.3_exp.1` - Experimental release

## Installation

### Linux (Self-extracting installer)

#### Quick Install (User Directory)
```bash
# Download and run the self-extracting installer to ~/.local/bin
# For x86_64/amd64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-amd64-install.sh | sh

# For ARM64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-arm64-install.sh | sh
```
**Note**: Default installation directory is `/usr/local/bin` (system-wide). For user-only installation, the installer will use `~/.local/bin` if `/usr/local/bin` is not writable.

#### System-wide Install (Requires sudo)
```bash
# Install to /usr/local/bin (system-wide)
# For x86_64/amd64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-amd64-install.sh | sudo sh

# For ARM64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-arm64-install.sh | sudo sh
```

#### Custom Directory Install
```bash
# Install to custom directory (e.g., /opt/version)
# For x86_64/amd64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-amd64-install.sh | INSTALL_DIR=/opt/version sh

# For ARM64 systems:
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-linux-arm64-install.sh | INSTALL_DIR=/opt/version sh

# Add to PATH if needed
echo 'export PATH="/opt/version:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Windows (Scoop Package Manager)

#### First-time Scoop Setup
If you don't have Scoop installed, follow these steps:

1. **Open PowerShell as Administrator** (version 5.1 or later) and from the PS C:\> prompt, run:
2. **Set execution policy** (if needed):
```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```
3. **Install Scoop**:
```powershell
   Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```
   or shorten version
```powershell
   iwr -useb get.scoop.sh | iex
```
4. **Restart PowerShell** and verify installation:
```powershell
   scoop --version
```

#### Install Version CLI
```powershell
# Add the bucket (if not already added)
scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket

# Install version
scoop install burnes/version

# Update version
scoop update burnes/version
```

#### Alternative: Manual Installation
If you prefer not to use Scoop:
1. Download the latest Windows binary from [Releases](https://github.com/AlexBurnes/version-go/releases)
2. Extract to a directory in your PATH (e.g., `C:\Program Files\version\`)
3. Add the directory to your system PATH

### macOS (Self-extracting installer)

#### Quick Install (User Directory)
```bash
# Download and run the self-extracting installer to ~/.local/bin
# For Intel Macs (x86_64/amd64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-amd64-install.sh | sh

# For Apple Silicon Macs (ARM64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-arm64-install.sh | sh
```
**Note**: Default installation directory is `/usr/local/bin` (system-wide). For user-only installation, the installer will use `~/.local/bin` if `/usr/local/bin` is not writable.

#### System-wide Install (Requires sudo)
```bash
# Install to /usr/local/bin (system-wide)
# For Intel Macs (x86_64/amd64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-amd64-install.sh | sudo sh

# For Apple Silicon Macs (ARM64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-arm64-install.sh | sudo sh
```

#### Custom Directory Install
```bash
# Install to custom directory (e.g., /opt/version)
# For Intel Macs (x86_64/amd64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-amd64-install.sh | INSTALL_DIR=/opt/version sh

# For Apple Silicon Macs (ARM64):
wget -O - https://github.com/AlexBurnes/version-go/releases/latest/download/version-darwin-arm64-install.sh | INSTALL_DIR=/opt/version sh

# Add to PATH if needed
echo 'export PATH="/opt/version:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### From Source
```bash
# Clone the repository
git clone https://github.com/burnes/go-version.git
cd go-version

# Build using CMake
mkdir build && cd build
cmake .. && make version

# Or build directly with Go
cd src
go build -o ../bin/version .
```

## Usage

### Basic Commands

```bash
# Show help
version --help

# Show version
version --version

# Validate a version string
version check 1.2.3
version check 1.2.3-alpha

# Get version type
version type 1.2.3-alpha
# Output: Pre release

# Get CMake build type
version build-type 1.2.3-alpha
# Output: Debug

# Sort versions from stdin
echo "1.2.3 1.2.4 1.2.3-alpha 2.0.0" | version sort
# Output:
# 1.2.3
# 1.2.3-alpha
# 1.2.4
# 2.0.0
```

### Git Integration

```bash
# Get current project version from git tags
version version

# Get project name from git remote
version project

# Get module name from git remote
version module

# Get full project-version-release
version full

# Check if current version is greatest among all tags
version check-greatest
```

### Version Bumping

```bash
# Smart bump (automatically determines appropriate increment)
version bump                    # Uses current git version
version bump 1.2.3             # Bump specific version

# Specific bump types
version bump 1.2.3 major       # 1.2.3 -> 2.0.0
version bump 1.2.3 minor       # 1.2.3 -> 1.3.0
version bump 1.2.3 patch       # 1.2.3 -> 1.2.4
version bump 1.2.3 alpha       # 1.2.3 -> 1.2.3~alpha.1
version bump 1.2.3 beta        # 1.2.3 -> 1.2.3~beta.1
version bump 1.2.3 rc          # 1.2.3 -> 1.2.3~rc.1
version bump 1.2.3 fix         # 1.2.3 -> 1.2.3.fix.1
version bump 1.2.3 next        # 1.2.3 -> 1.2.3.next.1
version bump 1.2.3 post        # 1.2.3 -> 1.2.3.post.1
version bump 1.2.3 feat        # 1.2.3 -> 1.2.3_feat.1

# Increment existing version types
version bump 1.2.3~alpha.1     # 1.2.3~alpha.1 -> 1.2.3~alpha.2
version bump 1.2.3.fix.1       # 1.2.3.fix.1 -> 1.2.3.fix.2
version bump 1.2.3_feat.1      # 1.2.3_feat.1 -> 1.2.3_feat.2

# Get help for bump command
version bump --help
```

### Project Configuration (.project.yml)

The utility supports a `.project.yml` configuration file for consistent project naming across build utilities. If present, it takes precedence over git-based detection.

**File Format** (`.project.yml` in project root):
```yaml
project:
  name: "my-project"
  modules:
    - "primary-module"    # First is primary
    - "secondary-module"
    - "another-module"
```

**Behavior**:
- If `.project.yml` exists and is valid, use it for project and module names
- If `.project.yml` doesn't exist or is invalid, fall back to git-based detection
- Other build utilities can also use this file for consistent project naming

**Example**:
```bash
# With .project.yml present
version project
# Output: my-project

version module
# Output: primary-module

# Without .project.yml (fallback to git)
version project
# Output: user-repository-name
```

### Options

- `-h, --help` - Print help and exit
- `-V, --version` - Print version and exit
- `-v, --verbose` - Verbose output
- `-d, --debug` - Debug output
- `--no-color` - Disable colored output

## Development

### Prerequisites

- Go 1.22 or later
- CMake 3.16 or later
- Git (for version extraction)
- **Buildfab** (for unified build management)

### Installing Buildfab

The project uses [buildfab](https://github.com/AlexBurnes/buildfab) for unified build management. Install buildfab before building the project:

#### System Installation (Recommended)
```bash
# Install buildfab system-wide
wget -O - https://github.com/AlexBurnes/buildfab/releases/latest/download/buildfab-linux-amd64-install.sh | sudo sh
```

#### Local Installation (Development)
```bash
# Install buildfab locally to ./scripts directory
wget -O - https://github.com/AlexBurnes/buildfab/releases/latest/download/buildfab-linux-amd64-install.sh | INSTALL_DIR=./scripts sh
```

### Buildfab Configuration

The project uses `.project.yml` for buildfab configuration with the following stages and actions:

- **Stages**: `pre-push`, `build`, `test`, `release`
- **Actions**: Conan dependency management, CMake configuration, cross-platform testing, GoReleaser integration
- **Parallel Execution**: Cross-platform tests run simultaneously using Docker

See [Developer Workflow](docs/Developer-workflow.md) for complete buildfab usage documentation.

### Pre-Push Hook Setup (Recommended for Developers)

For the best development experience, set up the pre-push hook to automatically run project checks before pushing commits:

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
# - Configure project settings
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

### Building

```bash
# Clone repository
git clone https://github.com/burnes/go-version.git
cd go-version

# Build using buildfab (recommended)
buildfab build       # Full build with Conan, CMake, GoReleaser
buildfab test        # Cross-platform testing with Docker

# Or build individual stages
buildfab pre-push    # Version checks and validation
buildfab build       # Conan dependencies, CMake build, GoReleaser dry-run
buildfab test        # Cross-platform testing
buildfab release     # Complete release process

# Legacy CMake build (still available)
mkdir build && cd build
cmake .. && make version
make test
make test-coverage
make format

# Lint code
make lint
```

### Project Structure

```
version/
├── src/                    # Go source code
│   ├── main.go            # CLI interface
│   ├── version.go         # Version parsing logic
│   ├── git.go            # Git integration
│   └── *_test.go         # Test files
├── bin/                   # Built binaries
├── build/                 # CMake build directory
├── docs/                  # Documentation
├── packaging/             # Distribution packages
└── CMakeLists.txt         # CMake configuration
```

## Version Precedence

Versions are sorted according to the following precedence order:

1. **Release versions** (highest priority)
2. **Prerelease versions** (alpha, beta, rc, pre)
3. **Postrelease versions** (fix, next, post)
4. **Intermediate versions** (feature, experimental)

Within each category, versions are sorted by:
- Major version number
- Minor version number  
- Patch version number
- Type-specific identifiers (alphanumeric comparison)

## Exit Codes

- `0` - Success, valid version
- `1` - Error, invalid input or failure
- `2` - System error

## Project Structure

This project follows standard Go conventions:

```
version-go/
├── cmd/version/          # CLI executable source code
│   ├── main.go          # CLI entry point
│   ├── git.go           # Git integration
│   └── version.go       # Version functions
├── pkg/version/         # Reusable library package
│   ├── version.go       # Library implementation
│   └── version_test.go  # Library tests
├── examples/basic/      # Example usage
└── docs/               # Documentation
```

## Library Usage

The version utility also provides a reusable Go library package for other utilities:

```go
package main

import (
    "fmt"
    "log"
    "github.com/AlexBurnes/version-go/pkg/version"
)

func main() {
    // Get version from git repository
    gitVersion, err := version.GetVersion()
    if err != nil {
        if version.IsGitNotFound(err) {
            fmt.Println("Git is not installed")
        } else if version.IsNotGitRepo(err) {
            fmt.Println("Not in a git repository")
        } else if version.IsNoGitTags(err) {
            fmt.Println("No version tags found")
        } else {
            log.Fatal(err)
        }
    } else {
        fmt.Printf("Git Version: %s\n", gitVersion)
    }
    
    // Parse a version
    v, err := version.Parse("1.2.3-alpha.1")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Version: %s\n", v.String())
    fmt.Printf("Type: %s\n", v.Type.String())
    
    // Sort versions
    versions := []string{"2.0.0", "1.2.3", "1.2.3-alpha"}
    sorted, err := version.Sort(versions)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Sorted: %v\n", sorted)
}
```

See [Library Documentation](docs/Library.md) for complete API reference and examples.

## Examples

### CI/CD Integration

```bash
# Validate version in build script
if ! version check $VERSION; then
    echo "Invalid version: $VERSION"
    exit 1
fi

# Get build type for CMake
BUILD_TYPE=$(version build-type $VERSION)
cmake -DCMAKE_BUILD_TYPE=$BUILD_TYPE ..

# Check if this is the greatest version
if version check-greatest $VERSION; then
    echo "This is the latest version"
fi
```

### Version Sorting

```bash
# Sort a list of versions
versions="1.2.3 1.2.4 1.2.3-alpha 2.0.0 1.2.3.fix"
echo $versions | version sort
```

## Developer Workflow

This project includes a complete developer workflow with automated version management:

### Quick Start
```bash
# Bump version (patch/minor/major)
./scripts/version-bump patch

# Make changes, test, and release
go test ./... -v
./buildtools/build-and-package.sh dry-run
git add . && git commit -m "feat: changes for $(cat VERSION)"
git tag $(cat VERSION)
git push origin master && git push origin $(cat VERSION)
./buildtools/build-and-package.sh release
```

### Key Features
- **VERSION file**: Centralized version management
- **Pre-push hooks**: Automated version validation
- **Helper scripts**: `version-bump` and `version-check`
- **Complete workflow**: From planning to publication

See [docs/Developer-workflow.md](docs/Developer-workflow.md) for detailed instructions.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the Apache License, Version 2.0 - see the LICENSE file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## Support

- Issues: [GitHub Issues](https://github.com/burnes/go-version/issues)
- Documentation: [Project Wiki](https://github.com/burnes/go-version/wiki)