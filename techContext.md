# Tech Context: Version CLI Utility

## Technologies Used
- **Go 1.22+**: Primary programming language
- **CMake**: Build system orchestrator
- **Conan**: Dependency management (including Go itself)
- **GoReleaser**: Cross-platform build and distribution
- **Git**: Version control and tag integration
- **Bash**: Build scripts and installation scripts
- **Makeself**: Linux self-extracting archive creation
- **Scoop**: Windows package management
- **tar.gz**: Linux distribution format
- **ZIP**: Windows distribution format

## Development Setup
- **Go Environment**: Managed via Conan (`conan install golang/<version>`)
- **Build Directory**: All build outputs go to `bin/` directory only
- **Source Structure**: 
  - `src/` - Go source code
  - `bin/` - Built binaries
  - `docs/` - Documentation
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
- **Standard Library Only**: `fmt`, `os`, `strings`, `sort`, `regexp`
- **Potential External**:
  - `go-git` for git integration (if standard library insufficient)
  - `cobra` for advanced CLI features (replace `flag` package)
- **Build Dependencies**:
  - CMake (via Conan)
  - Go compiler (via Conan)
  - GoReleaser (via Conan or direct install)
- **Packaging Dependencies**:
  - Makeself (for Linux self-extracting archives)
  - Scoop (for Windows package management)
  - Git (for bucket repository management)

## Tool Usage Patterns
- **Development**: `go build`, `go test`, `go mod` for standard Go development
- **Local Build**: `build-conan.sh` script orchestrates Conan → CMake → Go compilation
- **Packaging Build**: `build-goreleaser.sh` script orchestrates GoReleaser → Conan hooks → Cross-compilation
- **Testing**: `go test ./... -race` for comprehensive testing
- **Linux Packaging**: Makeself creates self-extracting .run archives
- **Windows Packaging**: Scoop manages package installation and updates
- **Distribution**: GoReleaser handles cross-platform builds and packaging
- **Documentation**: Markdown files with Memory Bank integration
- **Version Control**: Git with semantic versioning tags for releases