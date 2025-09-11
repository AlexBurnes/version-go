# System Patterns: Version CLI Utility

## System Architecture
**✅ IMPLEMENTED WITH LIBRARY SUPPORT** - The system follows a modular CLI architecture with clear separation of concerns and reusable library package:

```
CLI Interface Layer ✅
├── Command Parser (flag package) ✅
├── Input/Output Handlers ✅
└── Error Handling & Exit Codes ✅

Library Package Layer ✅
├── pkg/version - Reusable Version Library ✅
│   ├── Version Parser (Custom Regex Engine) ✅
│   ├── Version Validator ✅
│   ├── Version Sorter ✅
│   ├── Type System (Release/Prerelease/Postrelease/Intermediate) ✅
│   └── Git Tag Conversion ✅

Core Logic Layer ✅
├── CLI Command Handlers ✅
├── Git Integration ✅
└── Library Integration ✅

Utility Layer ✅
├── String Processing ✅
├── File I/O ✅
└── Platform Abstractions ✅

Build System Layer ✅
├── Local Development: Conan + CMake + bash scripts ✅
├── Packaging: GoReleaser + Conan hooks ✅
└── Cross-Platform: Automated builds for Linux/Windows/macOS ✅

Packaging Layer (Ready)
├── Linux: Makeself Self-Extracting Archives
├── Windows: Scoop Package Manager
└── Cross-Platform: GoReleaser Integration
```

## Key Technical Decisions
- **✅ CLI Framework**: Implemented using Go's `flag` package for command parsing
- **✅ Library Package**: Refactored core functionality into reusable `pkg/version` package
- **✅ Grammar Engine**: Custom regex-based parser for extended version format support (avoiding semver library complexity)
- **✅ Git Integration**: Implemented using `os/exec` for git operations (standard library only)
- **✅ Build System**: CMake as orchestrator, Conan for dependency management, GoReleaser for distribution
- **✅ Local Development**: Conan + CMake + bash scripts for development and testing
- **✅ Packaging Build**: GoReleaser with Conan hooks for automated cross-platform builds
- **Linux Packaging**: Makeself for self-extracting archives with professional installation experience
- **Windows Packaging**: Scoop package manager for easy installation and updates
- **✅ Testing**: Standard Go testing framework with comprehensive library tests and race detection
- **✅ Error Handling**: Structured error types with proper exit code mapping
- **✅ API Design**: Clean, well-documented public API with comprehensive examples

## Design Patterns in Use
- **✅ Command Pattern**: Each CLI command implemented as a separate handler in main.go
- **✅ Strategy Pattern**: Different parsing strategies for different version types (release, prerelease, postrelease, intermediate)
- **✅ Factory Pattern**: Version object creation based on input string analysis in ParseVersion()
- **✅ Builder Pattern**: Complex version object construction with validation
- **✅ Template Method**: Common validation and sorting logic with type-specific implementations

## Component Relationships
- **CLI Interface** → **Command Handlers** → **Library Package** → **Core Logic**
- **Library Package** → **Version Parser** ← **BNF Grammar Rules** → **Version Objects**
- **Git Integration** → **Library Package** → **Version Objects**
- **Build System** → **Go Compiler** → **Static Binaries** → **Distribution Packages**
- **External Applications** → **Library Package** → **Version Objects**

## Critical Implementation Paths
1. **✅ Version Parsing Pipeline**: Input validation → Grammar parsing → Object creation → Validation
2. **✅ Sorting Algorithm**: Parse all versions → Categorize by type → Apply precedence rules → Sort within categories
3. **✅ Git Integration**: Read git tags → Parse versions → Validate → Return appropriate version
4. **✅ Local Build Pipeline**: Source code → Conan deps → CMake config → Go compilation → Static binary
5. **Packaging Build Pipeline**: GoReleaser → Conan hooks → Tool installation → Cross-compilation → Distribution packages
6. **Linux Distribution**: Binary + install.sh → Makeself → Self-extracting .run archive → User installation
7. **Windows Distribution**: Binary → ZIP archive → Scoop manifest → Scoop bucket → User installation
8. **Cross-Platform Distribution**: GoReleaser → Multiple packages → GitHub releases → Automated distribution