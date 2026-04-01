# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

`version-go` is a Go-based semantic version parsing, validation, sorting, and bumping utility — a cross-platform replacement for legacy bash version scripts. It ships as both a standalone CLI binary and a reusable Go library. It is used by buildfab but works independently of it.

Module path: `github.com/AlexBurnes/version-go`

## Commands

### Build
```bash
# Compile the CLI binary (requires Conan + CMake)
cmake --preset conan-release
cmake --build --preset conan-release --target version

# Or via buildfab (includes linting and GoReleaser dry-run)
buildfab build
```

### Test
```bash
# All unit and integration tests
go test ./... -v -race

# Run a single test
go test ./pkg/version/... -run TestGetVersion -v

# Cross-platform Docker tests
buildfab test
```

### Lint
```bash
golangci-lint run
```

## Architecture

The project has two layers:

**`pkg/version/` — library** (import-safe, used by external consumers):
- `version.go` — core `Version` struct, `Parse()`, `Sort()`, `ConvertGitTag()`; supports extended SemVer: release (`1.2.3`), prerelease (`1.2.3~beta.1`), postrelease (`1.2.3.fix`), intermediate (`1.2.3_feature`)
- `git.go` — `GetVersion()`, `GetVersionWithPrefix()`, `GetRawTag()`, `GetRawVersion()`, `GetProjectName()`, `GetModuleName()`; wraps `git` subprocess calls; exposes typed `GitError` with `IsGitNotFound`, `IsNotGitRepo`, `IsNoGitTags` predicates
- `bump.go` — `Bump()` with `BumpType` enum (smart, major, minor, patch, alpha, beta, rc, fix, next, post, feat)
- `platform.go` — OS/arch detection helpers
- `config.go` — `.project.yml` parser via `gopkg.in/yaml.v3`

**`cmd/version/` — CLI**:
- `main.go` — command dispatcher, flag parsing, color-coded help
- `version.go` — validation, sorting, type-detection commands
- `git.go` — git-backed commands exposed to CLI
- `bump.go` — bump commands

## Version Format

The project uses an extended BNF grammar beyond standard SemVer:
- `-` separates prerelease (`1.2.3-alpha.1`)
- `~` is the internal canonical separator (git tags use `-`, library converts to `~`)
- `.` separates postrelease (`1.2.3.fix`)
- `_` separates intermediate/dev (`1.2.3_feature`)

`ConvertGitTag()` in `pkg/version/version.go` performs the `-` → `~` conversion.

## Key Design Constraints

- The `pkg/version` library must remain importable without pulling in CLI concerns.
- Git operations use subprocess calls (`exec.Command("git", ...)`) — no `go-git` library.
- Minimal dependencies: only `gopkg.in/yaml.v3` plus stdlib.
- Binaries are statically linked for cross-platform distribution (Linux amd64/arm64, macOS amd64/arm64, Windows amd64).

## Release Process

Releases are triggered by pushing a `v*.*.*` tag, which runs `.github/workflows/release.yml`. GoReleaser handles multi-platform builds, archives, and publishing to GitHub Releases and the Scoop bucket. Do not manually edit `.goreleaser.yml` archives without testing via `goreleaser release --snapshot --clean`.

## Coding Standards

### Go Style
- Imports ordered: standard library, third-party, local application.
- Every exported symbol must have a GoDoc comment: `// FunctionName …`.
- Package names: lowercase, short, no `util`/`helpers` buckets.
- Functions should stay under ~60 lines; extract helpers when logic branches.
- No panics in normal control flow.
- No global mutable state; pass dependencies explicitly.
- Prefer `errors.Is`/`errors.As` and sentinel errors when branching by error kind; wrap with context: `fmt.Errorf("opening %q: %w", path, err)`.
- CLI errors go to **stderr**; machine-parsable output goes to **stdout**.
- Cross-platform paths: always use `filepath.*`, never hardcode `/` or `\`.
- Discover external tools (e.g. `git`) via `exec.LookPath`; allow override via env/flag.
- CGO disabled; builds must be reproducible with `-trimpath -ldflags="-s -w"`.

### Testing
- Table-driven tests with `t.Run` subtests; no sleeps or timing-based assertions.
- Run with `-race` in CI: `go test ./... -race`.
- Target ≥70% coverage in parsing/logic packages.
- Golden tests for stable textual outputs — fixtures in `testdata/`.
- Do not rely on network/filesystem unless explicitly marked as integration test.

## Version Management

The `VERSION` file is the single source of truth for the current version.

**Before making any changes**, check if `VERSION` matches the latest git tag. If they match, bump the version first:
```bash
scripts/version-bump-with-file patch   # bug fixes — updates VERSION + packaging files
scripts/version-bump-with-file minor   # new features
scripts/version-bump-with-file major   # breaking changes
```

This script updates `VERSION`, `packaging/windows/scoop-bucket/version.json`, and `packaging/macos/version.rb`.

**Spec source of truth**: `docs/Project-specification.md`. Do not implement features not defined there.

## Git & Commit Conventions

### Commit message format
```
<type>(<scope>): <vX.Y.Z> <description>

<body — what and why, not how; bullet points; wrap at 72 chars>

<footer: Closes #123 / BREAKING CHANGE: …>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`, `perf`, `revert`.
Scopes (lowercase, hyphenated): `api`, `library`, `tests`, `docs`, `build`, `ci`.
Description: imperative mood, lowercase, no trailing period, ≤50 chars.

Example:
```
feat(api): v1.5.0, add GetRawVersion function

- Expose raw git describe output without tag conversion
- Add unit tests for empty repo and detached HEAD cases

Closes #42
```

### Pushing releases
Always push commits and tags together:
```bash
git push origin master --tags
```

### Retagging (when a tag push fails)
Do **not** bump the version. Instead:
```bash
git tag -d vX.Y.Z
git add .
git commit -m "chore(retag): retag vX.Y.Z after push failure

- Remove and recreate tag due to push failure
- No version bump"
git tag vX.Y.Z
git push origin master --tags
```

Tag only from green CI. Never force-push tags.

## Documentation Rules

- `CHANGELOG.md` **must** be updated for every code or documentation change.
- Use `## [X.Y.Z] - YYYY-MM-DD` sections, newest first; historical dates from `git log -1 --format="%ai" <tag>`.
- Doc files in `docs/` use `First-word-second-word.md` naming (e.g. `Developer-workflow.md`, `BNF-grammar.md`). `README.md` and `CHANGELOG.md` are exceptions.
- `docs/Project-specification.md` is the single source of truth for behavior and spec; update it when behavior changes.
