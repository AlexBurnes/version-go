# System Patterns: Version CLI Utility

## System Architecture
The system follows a modular CLI architecture with clear separation of concerns:

```
CLI Interface Layer
├── Command Parser (cobra)
├── Input/Output Handlers
└── Error Handling & Exit Codes

Core Logic Layer
├── Version Parser (BNF Grammar Engine)
├── Version Validator
├── Version Sorter
└── Git Integration

Utility Layer
├── String Processing
├── File I/O
└── Platform Abstractions

Build System Layer
├── Local Development: Conan + CMake + bash scripts
├── Packaging: GoReleaser + Conan hooks
└── Cross-Platform: Automated builds for Linux/Windows/macOS

Packaging Layer
├── Linux: Makeself Self-Extracting Archives
├── Windows: Scoop Package Manager
└── Cross-Platform: GoReleaser Integration
```

## Key Technical Decisions
- **CLI Framework**: Use Go's package or `cobra` for command parsing
- **Grammar Engine**: Custom BNF parser for extended version format support
- **Git Integration**: Use `go-git` library for git tag operations
- **Build System**: CMake as orchestrator, Conan for dependency management, GoReleaser for distribution
- **Local Development**: Conan + CMake + bash scripts for development and testing
- **Packaging Build**: GoReleaser with Conan hooks for automated cross-platform builds
- **Linux Packaging**: Makeself for self-extracting archives with professional installation experience
- **Windows Packaging**: Scoop package manager for easy installation and updates
- **Testing**: Standard Go testing framework with comprehensive test coverage
- **Error Handling**: Structured error types with proper exit code mapping

## Design Patterns in Use
- **Command Pattern**: Each CLI command implemented as a separate handler
- **Strategy Pattern**: Different parsing strategies for different version types (release, prerelease, postrelease, intermediate)
- **Factory Pattern**: Version object creation based on input string analysis
- **Builder Pattern**: Complex version object construction with validation
- **Template Method**: Common validation and sorting logic with type-specific implementations

## Component Relationships
- **CLI Interface** → **Command Handlers** → **Core Logic** → **Utility Functions**
- **Version Parser** ← **BNF Grammar Rules** → **Version Objects**
- **Git Integration** → **Version Parser** → **Version Objects**
- **Build System** → **Go Compiler** → **Static Binaries** → **Distribution Packages**

## Critical Implementation Paths
1. **Version Parsing Pipeline**: Input validation → Grammar parsing → Object creation → Validation
2. **Sorting Algorithm**: Parse all versions → Categorize by type → Apply precedence rules → Sort within categories
3. **Git Integration**: Read git tags → Parse versions → Validate → Return appropriate version
4. **Local Build Pipeline**: Source code → Conan deps → CMake config → Go compilation → Static binary
5. **Packaging Build Pipeline**: GoReleaser → Conan hooks → Tool installation → Cross-compilation → Distribution packages
6. **Linux Distribution**: Binary + install.sh → Makeself → Self-extracting .run archive → User installation
7. **Windows Distribution**: Binary → ZIP archive → Scoop manifest → Scoop bucket → User installation
8. **Cross-Platform Distribution**: GoReleaser → Multiple packages → GitHub releases → Automated distribution