# Project Brief: Version CLI Utility

## Overview
A cross-platform command-line utility `version`, written in Go, that provides semantic version parsing, validation, and ordering. This tool replaces legacy bash scripts (`version`, `version-check`, `describe`) currently used in build pipelines and must support Linux, Windows, and macOS with a reproducible build/distribution process via GoReleaser.

## Core Requirements
- Semantic version parsing and validation against a custom grammar (BNF)
- Ordering and sorting of versions according to defined precedence rules
- Command set for integration in CI/CD and build environments (`check`, `check-greatest`, `type`, etc.)
- Cross-platform support (Linux, Windows, macOS) with static binaries
- Distribution via standard OS channels (Scoop for Windows, tar.gz + install.sh for Linux)
- Git integration for describing current project version from git tags

## Goals
- Replace legacy bash scripts with a robust, maintainable Go implementation
- Provide consistent version handling across different build environments
- Enable reproducible builds and distribution across multiple platforms
- Support extended grammar beyond SemVer 2.0 with prerelease, postrelease, and intermediate identifiers
- Maintain compatibility with existing build pipeline integrations

## Project Scope
**In Scope:**
- CLI tool with validation and sorting capabilities
- Custom BNF grammar supporting extended version formats
- Cross-platform builds with GoReleaser
- Integration with existing build scripts and CI/CD pipelines
- Comprehensive testing and documentation

**Out of Scope:**
- GUI interface or web-based version management
- Version storage or persistence beyond git integration
- Real-time version monitoring or alerting
- Integration with package managers beyond Scoop and tar.gz distribution