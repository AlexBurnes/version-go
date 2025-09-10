# Active Context: Version CLI Utility

## Current Work Focus
Setting up the Memory Bank system and establishing the foundational project structure. Currently in the initial project setup phase, creating comprehensive documentation to guide future development work.

## Recent Changes
- Created Memory Bank files based on project.md and cursor rules
- Established project brief with core requirements and scope
- Defined product context and user experience goals
- Set up memory bank instructions for consistent documentation practices

## Next Steps
- Complete the remaining Memory Bank files (systemPatterns.md, techContext.md, progress.md)
- Review existing source code in src/old/ to understand current implementation
- Set up proper Go module structure in src/
- Create initial CLI framework and command structure
- Implement basic version parsing and validation logic

## Active Decisions and Considerations
- **Language Choice**: Go is already decided and specified in project requirements
- **Build System**: CMake + Conan + bash scripts as specified in project.md
- **Distribution**: GoReleaser for cross-platform builds and distribution
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