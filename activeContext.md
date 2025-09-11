# Active Context: Version CLI Utility

## Current Work Focus
Completed comprehensive build documentation covering both local development builds and automated packaging builds. Currently focused on documenting the complete build process using Conan + CMake for local builds and GoReleaser with Conan hooks for packaging builds.

## Recent Changes
- Created comprehensive build guide (docs/BUILD.md) covering both build approaches
- Documented local development build process using Conan + CMake + bash scripts
- Documented packaging build process using GoReleaser with Conan hooks
- Created build-goreleaser.sh script for automated packaging builds
- Documented GoReleaser hooks and Conan integration for cross-platform builds
- Updated memory bank with complete build process information

## Next Steps
- Review existing source code in src/old/ to understand current implementation
- Set up proper Go module structure in src/
- Create initial CLI framework and command structure
- Implement basic version parsing and validation logic
- Test both build processes with sample builds
- Validate GoReleaser configuration and Conan integration

## Active Decisions and Considerations
- **Language Choice**: Go is already decided and specified in project requirements
- **Build System**: CMake + Conan + bash scripts as specified in project.md
- **Distribution**: GoReleaser for cross-platform builds and distribution
- **Linux Packaging**: Makeself for self-extracting archives with professional installation experience
- **Windows Packaging**: Scoop package manager for easy installation and updates
- **Grammar**: Custom BNF grammar extending SemVer 2.0 with prerelease, postrelease, and intermediate identifiers
- **Compatibility**: Must maintain 100% compatibility with existing bash script interfaces

## Important Patterns and Preferences
- **Documentation First**: All decisions and patterns must be documented in Memory Bank
- **Cross-platform**: All code must work identically on Linux, Windows, and macOS
- **Minimal Dependencies**: Prefer standard library, avoid external dependencies unless justified
- **Test Coverage**: Comprehensive unit and integration tests required
- **Exit Codes**: POSIX-compliant exit codes (0 for success, >=1 for errors)

## Learnings and Project Insights
- Project has clear, well-defined requirements in project.md
- Existing bash scripts provide a reference implementation to maintain compatibility
- Custom grammar requirements mean standard SemVer libraries cannot be used directly
- Build system complexity requires careful setup with CMake, Conan, and GoReleaser integration