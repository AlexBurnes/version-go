# Tech Context: Version CLI Utility

## Technologies Used
- **✅ Go 1.22+**: Primary programming language (implemented)
- **✅ CMake**: Build system orchestrator (working)
- **✅ Conan**: Dependency management (including Go itself)
- **✅ GoReleaser**: Cross-platform build and distribution (verified working)
- **✅ Git**: Version control and tag integration (implemented)
- **✅ Bash**: Build scripts and installation scripts (working)
- **✅ Makeself**: Linux and macOS self-extracting archive creation (implemented)
- **✅ Scoop**: Windows package management (implemented)
- **✅ tar.gz**: Linux distribution format (working)
- **✅ ZIP**: Windows distribution format (working)
- **✅ YAML**: Configuration file parsing (gopkg.in/yaml.v3)

## Development Setup
- **Go Environment**: Managed via Conan (`conan install golang/<version>`)
- **Build Directory**: All build outputs go to `bin/` directory only
- **Source Structure**: 
  - Root - Go CLI source code (main.go, version.go, git.go)
  - `pkg/version/` - Reusable library package
  - `examples/basic/` - Library usage example
  - `bin/` - Built binaries
  - `docs/` - Documentation including library docs
  - `test/` - Test configuration files and test scripts
- **Cross-platform**: Support for Linux (amd64/arm64), Windows (amd64/arm64), macOS (amd64/arm64)

## Technical Constraints
- **No CGO**: Static binaries only, no C dependencies
- **Minimal Dependencies**: Prefer standard library, external deps only when justified
- **Performance**: Handle 10k+ versions in sorting without significant slowdown
- **Memory**: Efficient parsing and sorting for large version lists
- **Compatibility**: Must work identically across all target platforms
- **Reproducible Builds**: All builds must be reproducible via GoReleaser

## Dependencies
- **✅ Standard Library Only**: `fmt`, `os`, `strings`, `sort`, `regexp`, `flag`, `os/exec`, `bufio` (implemented)
- **✅ Minimal External Dependencies**: Added `gopkg.in/yaml.v3` for configuration file parsing
- **Library Package**: Self-contained `pkg/version` package with minimal external dependencies
- **Build Dependencies**:
  - ✅ CMake (via Conan) - working
  - ✅ Go compiler (via Conan) - working
  - GoReleaser (via Conan or direct install) - ready
- **Packaging Dependencies**:
  - ✅ Makeself (for Linux and macOS self-extracting archives) - implemented
  - ✅ Scoop (for Windows package management) - implemented
  - ✅ Git (for bucket repository management) - working

## Tool Usage Patterns
- **✅ Development**: `go build`, `go test`, `go mod` for standard Go development (working)
- **✅ Local Build**: CMake orchestrates Conan → Go compilation (working)
- **✅ Main Build**: `build-and-package.sh` orchestrates complete build flow (build-conan → makeself → goreleaser)
- **✅ Packaging Build**: `build-goreleaser.sh` script orchestrates GoReleaser → Conan hooks → Cross-compilation (verified)
- **✅ Testing**: `go test ./... -race` for comprehensive testing including library tests (all tests passing)
- **✅ Library Testing**: `go test ./pkg/version/... -v` for library-specific tests
- **✅ Library Usage**: `go run examples/basic/main.go` for library demonstration
- **✅ Linux Packaging**: Makeself creates self-extracting .sh archives (implemented)
- **✅ macOS Packaging**: Makeself creates self-extracting .sh archives (implemented)
- **✅ Windows Packaging**: Scoop manages package installation and updates (implemented)
- **✅ Distribution**: GoReleaser handles cross-platform builds and packaging (verified working)
- **✅ Documentation**: Markdown files with Memory Bank integration and library documentation (complete)
- **✅ Version Control**: Git with semantic versioning tags for releases (working)

## CLI Commands
- **✅ Core Commands**: `project`, `module`, `modules`, `version`, `release`, `full` (implemented)
- **✅ Validation Commands**: `check`, `check-greatest` (implemented)
- **✅ Type Commands**: `type`, `build-type` (implemented)
- **✅ Utility Commands**: `sort` (implemented)
- **✅ Modules Command**: New `modules` command returns all modules from .project.yml or single git module name (v0.8.6)
- **✅ Configuration Support**: All commands support --config and --git flags for configuration control

## Packaging Tools and Scripts
- **✅ create-simple-installers.sh**: Creates simple installer scripts that download from GitHub releases
- **✅ create-all-installers.sh**: Batch creation of installers for all platforms
- **✅ installer-template.sh**: Template for simple installer scripts with platform-specific variables
- **✅ install.sh**: Core installation script included in release archives for manual installation
- **✅ GoReleaser Scoop Integration**: Automated Scoop manifest generation for Windows
- **✅ Clean Naming Logic**: Version cleaning to remove SNAPSHOT and hex suffixes
- **✅ No-Sudo Approach**: Install scripts don't use sudo internally - users run with sudo if needed
- **✅ Installer Cleanup Process**: Old installers are removed before creating new ones to prevent version mixing