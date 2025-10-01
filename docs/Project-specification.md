# Project: Version CLI Utility

## 1. Purpose & Scope
The purpose of this project is to design and implement a cross-platform command-line utility `version`, written in Go, that provides semantic version parsing, validation, and ordering.  
The tool replaces legacy bash scripts (`version`, `version-check`, `describe`) currently used in build pipelines.  
It must support Linux, Windows, and macOS and provide a reproducible build/distribution process via GoReleaser (binary archives, Scoop bucket, Linux tar.gz installer).

The scope includes:
- Semantic version parsing and validation against a custom grammar (BNF).
- Ordering and sorting of versions according to defined precedence rules.
- Command set for integration in CI/CD and build environments (e.g., `check`, `check-greatest`, `type`, etc.).
- Distribution via standard OS channels (Scoop for Windows, tar.gz + install.sh for Linux).

## 2. Features
- Validate a version string passed as an argument.
- Sort a list of version strings provided through `stdin`.
- Support extended grammar beyond SemVer 2.0: prerelease (`~alpha`, `~beta`, `~rc`, etc.), postrelease (`.fix`, `.post`, `.next`), intermediate identifiers (`_feature`, `_exp`).
- Provide exit codes consistent with POSIX conventions: `0` for valid, `>=1` for invalid.
- Git integration: allow describing the current project version from git tags.
- Command set inspired by original bash tool (`version`, `release`, `full`, `check`, `check-greatest`, `type`, `build-type`).
- Cross-platform builds with static binaries.
- Distribution through GoReleaser: Scoop manifests (Windows) and tar.gz installers (Linux).

## 3. Functional Requirements
### 3.1 Input / Output
- **Validation mode**:  
  - Input: version string as CLI argument.  
  - Output: no output, exit code `0` if valid, `1` if invalid.
- **Sorting mode**:  
  - Input: list of version strings via stdin.  
  - Output: sorted list to stdout, according to grammar rules.

### 3.2 Grammar & Ordering
The grammar follows a Backus–Naur Form (BNF) specification as defined in [BNF-grammar.md](docs/BNF-grammar.md) with the following main categories:
- `<version core>` ::= `<major> "." <minor> "." <patch>`
- Extended forms:
  - `<version core> "~" <pre-release>`
  - `<version core> "." <post-release>`
  - `<version core> "_" <intermediate-release>`

Precedence order for version categories (lowest to highest):
1. prerelease: `^v?[0-9]+\.[0-9]+\.[0-9]+~(alpha|beta|pre|rc)(...)$`
2. release: `^v?[0-9]+\.[0-9]+\.[0-9]+$`
3. postrelease: `^v?[0-9]+\.[0-9]+\.[0-9]+\.(fix|next|post)(...)$`
4. intermediate: `^v?[0-9]+\.[0-9]+\.[0-9]+_[a-zA-Z]+(...)$`

Sorting rules:
- Numeric identifiers are compared numerically (`0 < 1 < 2 < ... < 10`).
- Alphanumeric identifiers are compared lexically (`a < b < ... < z < A < ... < Z`, with correct handling of numeric suffixes like `a1 < a2 < a10`).

### 3.3 Command Set
- `project` → print project name.
- `module` → print module name.
- `version` → print project version.
- `release` → print project release number.
- `full` → print full project name-version-release.
- `check [ver]` → validate version (default: current git version).
- `check-greatest [ver]` → check if version is the greatest among tags.
- `type [ver]` → print semantic version type (release, prerelease, postrelease, intermediate).
- `build-type [ver]` → print CMake build type (Release/Debug).

### 3.4 Return Codes
- `0` — successful execution, valid version.
- `1` or higher — error, invalid input, or failure.

## 4. Non-Functional Requirements
- Implemented in Go (>=1.22).
- Cross-platform: Linux (amd64/arm64), Windows (amd64/arm64), macOS (amd64/arm64).
- No CGO; static binaries reproducible via GoReleaser.
- Minimal dependencies; only standard library unless justified.
- Performance: must handle large lists (10k+ versions) in sorting without significant slowdown.

## 5. Delivery & Distribution
- **GoReleaser** as primary build/release tool.
- **Windows**: distributed through Scoop bucket manifests.
- **Linux**: delivered as `.tar.gz` archive with `install.sh`.
- **macOS**: Homebrew tap (future extension).
- Checksums published for all artifacts.
- GitHub Actions for automated release on `v*.*.*` tags.

## 6. Documentation
- `README.md`: usage, examples, installation.
- `.cursor/rules.md`: coding rules, style, test requirements (also synced with memory-base MPC).
- `CHANGELOG.md`: release history and changelog.
- ADRs stored in `/docs/adr` and mirrored in MPC.

## 7. Testing
- Unit tests for:
  - Validation of correct/incorrect version strings.
  - Sorting correctness across mixed categories.
  - Edge cases: leading zeros, mixed alphanumeric identifiers.
- Integration tests for CLI commands.
- CI: run `go test ./... -race`.
- Provide test data sets, including rpm-style filenames.

## 8. Directory structure
- bin/ - builded binaries, install binaries only in this directory
- cmd/version/ - CLI source code (following Go conventions)
- pkg/version/ - reusable library package
- docs/ - all documentation for project
  - BNF-grammar.md - BNF grammar specification for version parsing
- test/ - test scripts and test binaries
- buildtools/ - build scripts (build-goreleaser.sh, etc.)
- VERSION - current version number for centralized management
- docs/Developer-workflow.md - complete developer workflow documentation
- README.md - short readme, purpose, how to build, links to other docs
- docs/Project-specification.md - main documentation for development (this file)
- CHANGELOG.md - changelog
- LICENSE - project license

## 9. Build
- Project build must be automated via **bash build scripts** in `buildtools/` directory.
- Use **CMake** as the top-level build orchestrator to compile and install Go binaries.
- All build tools required must be managed via **Conan** in integration with CMake.  
  - Even Go itself must be installed via Conan (e.g., `conan install golang/<version>`).
- Build outputs must be placed in `bin/` directory only.
- Support reproducible builds across Linux, Windows, and macOS.
- Main build scripts: `build-and-package.sh`, `build-goreleaser.sh`, `build-conan.sh`, `create-all-installers.sh`

## 10. References
- Semantic Versioning 2.0: https://semver.org/
- RPM naming conventions:  
  - http://ftp.rpm.org/max-rpm/ch-rpm-file-format.html  
  - https://fedoraproject.org/wiki/PackagingDrafts/TildeVersioning
- Ragel FSM Compiler: https://www.colm.net/open-source/ragel/
- Build scripts (legacy bash versions): https://github.com/AlexBurnes/build_scripts.git
- Related Go libraries:  
  - https://github.com/Masterminds/semver  
  - https://github.com/Masterminds/vert  
  *(not directly usable due to grammar differences)*
