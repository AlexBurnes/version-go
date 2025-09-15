# Linux Packaging Guide: Version CLI Utility

## Overview
This guide details the Linux packaging strategy using makeself for creating self-extracting archives. The approach provides a professional installation experience while maintaining compatibility across different Linux distributions.

## Makeself Integration

### What is Makeself?
[Makeself](https://github.com/megastep/makeself.git) is a shell script that generates self-extractable compressed tar archives. It creates a single executable file that contains:

- Compressed archive of the application files
- Installation script that runs automatically
- Integrity checking with checksums
- Progress indication during extraction
- Automatic cleanup after installation

### Why Makeself?
- **Single File Distribution**: Users only need to download one `.run` file
- **Cross-Distribution**: Works on any Linux distribution with bash
- **Professional Experience**: Progress bars, error handling, user confirmation
- **Security**: Built-in checksum verification
- **No Dependencies**: Pure shell script, no additional tools required

## File Structure

```
packaging/linux/
├── README.md              # User installation instructions
├── install.sh             # Core installation script
├── makeself-header.sh     # Makeself header (from upstream)
└── makeself.sh            # Makeself script (from upstream)
```

### install.sh
The core installation script that:
- Detects installation directory (system-wide or user-local)
- Installs the binary with proper permissions
- Provides feedback to the user
- Handles errors gracefully

Key features:
- **Flexible Installation**: Supports both `/usr/local/bin` and `~/.local/bin`
- **Permission Handling**: Uses sudo when available, falls back to user installation
- **Validation**: Checks for binary presence before installation
- **User Feedback**: Clear progress indication and error messages

### makeself.sh
The makeself script (downloaded from upstream) that creates self-extracting archives with:
- **Compression**: gzip compression for optimal file size
- **Checksums**: MD5 and SHA256 for integrity verification
- **Script Execution**: Automatic execution of install.sh after extraction
- **Progress**: Real-time progress indication during extraction
- **Cleanup**: Automatic removal of temporary files

## Installation Process

### For End Users

#### Method 1: Direct Download and Run
```bash
# Download the self-extracting installer
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_linux_amd64.run

# Make it executable
chmod +x version_1.0.0_linux_amd64.run

# Run the installer (system-wide installation)
./version_1.0.0_linux_amd64.run

# Or install to user directory
APP_DIR=$HOME/.local/bin ./version_1.0.0_linux_amd64.run
```

#### Method 2: Manual Installation (Advanced Users)
```bash
# Extract the tar.gz archive
tar xf version_1.0.0_linux_amd64.tar.gz
cd version_1.0.0_linux_amd64

# Run the installation script
./install.sh

# Or install to user directory
APP_DIR=$HOME/.local/bin ./install.sh
```

### Installation Directory Options

#### System-wide Installation (Default)
- **Directory**: `/usr/local/bin`
- **Permissions**: Requires sudo or root access
- **Access**: Available to all users
- **Use Case**: Production servers, shared systems

#### User-local Installation
- **Directory**: `~/.local/bin` (or custom via `APP_DIR`)
- **Permissions**: No elevated privileges required
- **Access**: Available only to the installing user
- **Use Case**: Development environments, personal systems

## Build Process

### Creating the Self-Extracting Archive

#### Prerequisites
```bash
# Clone makeself repository
git clone https://github.com/megastep/makeself.git
cd makeself

# Make makeself executable
chmod +x makeself.sh
```

#### Archive Creation
```bash
# Create the self-extracting archive
./makeself.sh \
  --gzip \
  --sha256 \
  --license "Apache License 2.0" \
  --help-header "Version CLI Utility Installer" \
  --header "packaging/linux/makeself-header.sh" \
  "packaging/linux" \
  "version_1.0.0_linux_amd64.run" \
  "Version CLI Utility v1.0.0" \
  "./install.sh"
```

#### Parameters Explained
- `--gzip`: Use gzip compression for smaller archives
- `--sha256`: Include SHA256 checksum for integrity verification
- `--license`: Display license information during installation
- `--help-header`: Custom help text for the installer
- `--header`: Custom header script (optional)
- `packaging/linux`: Source directory containing files to archive
- `version_1.0.0_linux_amd64.run`: Output filename
- `Version CLI Utility v1.0.0`: Archive description
- `./install.sh`: Command to run after extraction

### Integration with GoReleaser

The makeself integration works with GoReleaser through custom build hooks:

```yaml
# .goreleaser.yml
builds:
  - id: linux-amd64
    goos: linux
    goarch: amd64
    binary: version

archives:
  - id: linux-makeself
    format: tar.gz
    files:
      - install.sh
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

hooks:
  post: |
    # Create makeself archive after GoReleaser builds
    ./packaging/linux/makeself.sh \
      --gzip \
      --sha256 \
      --license "Apache License 2.0" \
      "{{ .Artifacts }}" \
      "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}.run" \
      "{{ .ProjectName }} {{ .Version }}" \
      "./install.sh"
```

## Security Features

### Integrity Verification
- **MD5 Checksum**: Fast verification for basic integrity
- **SHA256 Checksum**: Strong cryptographic verification
- **Automatic Validation**: Checksums verified before installation
- **Error Handling**: Installation fails if checksums don't match

### Permission Management
- **Minimal Privileges**: Only requests sudo when necessary
- **User Choice**: Supports both system-wide and user-local installation
- **Secure Defaults**: Binary installed with 755 permissions
- **Validation**: Verifies installation success

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Solution: Install to user directory
APP_DIR=$HOME/.local/bin ./version_1.0.0_linux_amd64.run
```

#### Binary Not Found
```bash
# Check if binary is in PATH
echo $PATH
which version

# Add user bin to PATH (add to ~/.bashrc or ~/.profile)
export PATH="$HOME/.local/bin:$PATH"
```

#### Checksum Verification Failed
```bash
# Re-download the file
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_linux_amd64.run

# Verify manually
sha256sum version_1.0.0_linux_amd64.run
```

### Debug Mode
```bash
# Run with verbose output
bash -x version_1.0.0_linux_amd64.run
```

## Maintenance

### Updating Makeself
```bash
# Update makeself from upstream
cd packaging/linux
git clone https://github.com/megastep/makeself.git temp-makeself
cp temp-makeself/makeself.sh .
cp temp-makeself/makeself-header.sh .
rm -rf temp-makeself
```

### Testing New Versions
```bash
# Test locally before release
./packaging/linux/makeself.sh \
  --gzip \
  --sha256 \
  "packaging/linux" \
  "test-installer.run" \
  "Test Version" \
  "./install.sh"

# Test the installer
./test-installer.run
```

## Best Practices

### Archive Optimization
- Use gzip compression for best compatibility
- Include only necessary files in the archive
- Test on multiple Linux distributions
- Verify checksums after creation

### User Experience
- Provide clear installation instructions
- Support both installation methods (system-wide and user-local)
- Include helpful error messages
- Test installation process thoroughly

### Security
- Always verify checksums before installation
- Use strong compression and checksum algorithms
- Test on clean systems to ensure no conflicts
- Document security considerations clearly