# Tech Context: Version CLI Utility

## Technologies Used
- **✅ Go 1.22+**: Primary programming language (implemented)
- **✅ CMake**: Build system orchestrator (working)
- **✅ Conan**: Dependency management (including Go itself)
- **GoReleaser**: Cross-platform build and distribution (ready)
- **✅ Git**: Version control and tag integration (implemented)
- **✅ Bash**: Build scripts and installation scripts (working)
- **Makeself**: Linux self-extracting archive creation (ready)
- **Scoop**: Windows package management (ready)
- **tar.gz**: Linux distribution format (ready)
- **ZIP**: Windows distribution format (ready)

## Development Setup
- **Go Environment**: Managed via Conan (`conan install golang/<version>`)
- **Build Directory**: All build outputs go to `bin/` directory only
- **Source Structure**: 
  - Root - Go CLI source code (main.go, version.go, git.go)
  - `pkg/version/` - Reusable library package
  - `examples/basic/` - Library usage example
  - `bin/` - Built binaries
  - `docs/` - Documentation including library docs
  - `test/` - Test scripts and binaries
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
- **✅ No External Dependencies**: Avoided `go-git` and `cobra` as planned
- **Library Package**: Self-contained `pkg/version` package with no external dependencies
- **Build Dependencies**:
  - ✅ CMake (via Conan) - working
  - ✅ Go compiler (via Conan) - working
  - GoReleaser (via Conan or direct install) - ready
- **Packaging Dependencies**:
  - Makeself (for Linux self-extracting archives) - ready
  - Scoop (for Windows package management) - ready
  - Git (for bucket repository management) - working

## Tool Usage Patterns
- **✅ Development**: `go build`, `go test`, `go mod` for standard Go development (working)
- **✅ Local Build**: CMake orchestrates Conan → Go compilation (working)
- **Packaging Build**: `build-goreleaser.sh` script orchestrates GoReleaser → Conan hooks → Cross-compilation
- **✅ Testing**: `go test ./... -race` for comprehensive testing including library tests (all tests passing)
- **✅ Library Testing**: `go test ./pkg/version/... -v` for library-specific tests
- **✅ Library Usage**: `go run examples/basic/main.go` for library demonstration
- **Linux Packaging**: Makeself creates self-extracting .run archives (ready)
- **Windows Packaging**: Scoop manages package installation and updates (ready)
- **Distribution**: GoReleaser handles cross-platform builds and packaging (ready)
- **✅ Documentation**: Markdown files with Memory Bank integration and library documentation (complete)
- **✅ Version Control**: Git with semantic versioning tags for releases (working)