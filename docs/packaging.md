# Packaging Documentation: Version CLI Utility

## Overview
This document describes the packaging and distribution strategy for the Version CLI utility across Linux and Windows platforms. The packaging system uses different approaches for each platform while maintaining consistency in user experience.

## Directory Structure

```
packaging/
├── linux/
│   ├── README.md              # Linux installation instructions
│   ├── install.sh             # Core installation script (used by both archives and installers)
│   └── installer-template.sh  # Template for simple installers
└── windows/
    └── scoop-bucket/
        ├── README.md          # Windows installation instructions
        └── version.json       # Scoop manifest for version management
```

## Cross-Platform Packaging

### Strategy
The packaging system uses different approaches for different platforms:

- **Linux/macOS**: Simple installer scripts that download archives directly from GitHub releases
- **Windows**: Scoop package manager for easy installation and updates

This approach provides:
- Simple and reliable installation
- No complex self-extracting archives
- Platform-appropriate installation methods
- Easy maintenance and debugging
- Support for both direct execution and pipe execution (Linux/macOS)

### Components

#### 1. Core Installation Script (`install.sh`)
- **Location**: `packaging/linux/install.sh`
- **Purpose**: Core installation logic used by both archives and installers
- **Usage**: `./install.sh [install_directory]`
- **Features**:
  - Supports both system-wide (`/usr/local/bin`) and user-local (`~/.local/bin`) installation
  - Handles permission management with sudo when available
  - Automatic binary detection (version or version.exe)
  - Clear error messages and validation
  - Works with both direct execution and makeself extraction

#### 2. Simple Installer Scripts
- **Template**: `packaging/linux/installer-template.sh`
- **Generator**: `buildtools/create-simple-installers.sh`
- **Purpose**: Create platform-specific installer scripts that download from GitHub releases
- **Features**:
  - Downloads archives directly from GitHub releases
  - Uses the core `install.sh` script for installation
  - Support for both direct execution and pipe execution
  - Cross-platform compatibility (Linux, macOS)
  - No duplication of installation logic

#### 3. Installation Methods

**Method 1: Using Simple Installer Scripts**
```bash
# Direct execution with custom directory
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh
chmod +x version-1.0.0-linux-amd64-install.sh
./version-1.0.0-linux-amd64-install.sh /usr/local/bin

# Pipe execution (uses default directory)
wget -O - https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh | sh

# Pipe execution with custom directory
INSTALL_DIR=/custom/path wget -O - https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh | sh
```

**Method 2: Manual Installation from Archive**
```bash
# Download and extract archive
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_linux_amd64.tar.gz
tar -xzf version_1.0.0_linux_amd64.tar.gz
cd version_1.0.0_linux_amd64

# Install using the included install.sh
./install.sh /usr/local/bin

# Or install to user directory
./install.sh ~/.local/bin
```

### Installer Configuration
The simple installer scripts (Linux/macOS only):
- **Download**: Archives directly from GitHub releases
- **Extraction**: Automatic tar.gz extraction
- **Installation**: Binary installation to specified directory
- **Compatibility**: Works with both bash and sh shells
- **Error Handling**: Clear error messages and validation
- **Platforms**: Linux (amd64/arm64) and macOS (amd64/arm64)

**Note**: Windows users should use Scoop for installation:
```bash
scoop install version
```

## Windows Packaging

### Strategy
Windows packaging uses **Scoop** package manager for easy installation and updates. This approach provides:

- One-command installation
- Automatic updates
- Dependency management
- Integration with Windows package ecosystem

### Components

#### 1. Scoop Manifest (`version.json`)
- Defines package metadata and download URLs
- Supports multiple architectures (amd64, arm64)
- Includes checksums for security verification
- Specifies binary location and execution

#### 2. Installation Process
```powershell
# Add the custom bucket
scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket

# Install the package
scoop install burnes/version

# Update the package
scoop update version
```

### Scoop Configuration
The Scoop manifest includes:
- **Version**: Semantic version for package management
- **Architecture Support**: Both 64-bit and ARM64 binaries
- **Checksums**: SHA256 verification for security
- **Binary Location**: Automatic PATH configuration
- **Metadata**: Description, homepage, and license information

## Build Integration

### GoReleaser Configuration
Both packaging methods integrate with GoReleaser for automated builds:

```yaml
# .goreleaser.yml (excerpt)
archives:
  - format: tar.gz
    files:
      - install.sh
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

scoop:
  bucket:
    owner: AlexBurnes
    name: scoop-bucket
  commit_author:
    name: GoReleaser Bot
    email: bot@example.com
```

### Build Process
1. **Compilation**: GoReleaser builds static binaries for all platforms
2. **Linux**: Creates tar.gz archives with makeself integration
3. **Windows**: Updates Scoop manifest with new version and checksums
4. **Release**: Publishes all artifacts to GitHub releases

## Security Considerations

### Integrity Verification
- **Linux**: Makeself provides built-in checksum verification
- **Windows**: Scoop validates SHA256 checksums before installation
- **Both**: GoReleaser generates and publishes checksums for all artifacts

### Installation Permissions
- **Linux**: Supports both system-wide and user-local installation
- **Windows**: Scoop manages permissions and PATH configuration
- **Both**: No elevated privileges required for user-local installation

## Maintenance

### Version Updates
- **Automated**: GoReleaser handles version updates for both platforms
- **Manual**: Scoop manifest requires manual updates for complex changes
- **Testing**: Both packaging methods support local testing before release

### Troubleshooting
- **Linux**: Check makeself script execution and permission issues
- **Windows**: Verify Scoop bucket configuration and network access
- **Both**: Validate checksums and binary compatibility

## Future Enhancements

### Planned Improvements
- **macOS**: Homebrew tap integration for macOS distribution
- **Docker**: Container images for CI/CD integration
- **Snap/Flatpak**: Additional Linux distribution methods
- **MSI**: Windows MSI installer for enterprise environments

### Extension Points
- **Custom Scripts**: Additional post-installation scripts
- **Dependencies**: Package dependency management
- **Configuration**: Default configuration file installation
- **Documentation**: Automatic documentation installation