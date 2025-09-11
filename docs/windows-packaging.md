# Windows Packaging Guide: Version CLI Utility

## Overview
This guide details the Windows packaging strategy using Scoop package manager for easy installation and updates. Scoop provides a command-line package manager for Windows that simplifies software installation and management.

## Scoop Integration

### What is Scoop?
[Scoop](https://scoop.sh/) is a command-line installer for Windows that:

- Installs programs from the command line with minimal input
- Manages dependencies automatically
- Provides easy updates and uninstallation
- Integrates with Windows PATH automatically
- Supports multiple package sources and buckets

### Why Scoop?
- **Simple Installation**: One command to install the application
- **Automatic Updates**: Easy updates with `scoop update`
- **Dependency Management**: Handles PATH and environment setup
- **Windows Integration**: Works seamlessly with Windows command line
- **Version Management**: Easy switching between versions

## File Structure

```
packaging/windows/
└── scoop-bucket/
    ├── README.md          # Installation instructions
    └── version.json       # Scoop manifest
```

### version.json
The Scoop manifest that defines:
- Package metadata (name, version, description)
- Download URLs for different architectures
- Checksums for security verification
- Binary location and execution details

## Installation Process

### For End Users

#### Prerequisites
```powershell
# Install Scoop (if not already installed)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
```

#### Installation
```powershell
# Add the custom bucket
scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket

# Install the package
scoop install burnes/version

# Verify installation
version --version
```

#### Updates
```powershell
# Update the package
scoop update version

# Update all packages
scoop update
```

#### Uninstallation
```powershell
# Remove the package
scoop uninstall version
```

### Installation Directory Options

#### Default Installation
- **Directory**: `~/scoop/apps/version/current/`
- **Binary**: `~/scoop/apps/version/current/version.exe`
- **PATH**: Automatically added to user PATH
- **Access**: Available system-wide for the user

#### Global Installation (Administrator)
```powershell
# Install globally (requires admin privileges)
scoop install -g burnes/version
```

## Scoop Manifest Configuration

### version.json Structure
```json
{
  "version": "1.0.0",
  "description": "git describe-like CLI",
  "homepage": "https://github.com/AlexBurnes/version-go",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_windows_amd64.zip",
      "hash": "sha256-REPLACE_FROM_checksums.txt"
    },
    "arm64": {
      "url": "https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_windows_arm64.zip",
      "hash": "sha256-REPLACE_FROM_checksums.txt"
    }
  },
  "bin": ["version.exe"]
}
```

### Field Descriptions

#### Basic Metadata
- **version**: Semantic version of the package
- **description**: Short description of the package
- **homepage**: Project homepage URL
- **license**: Software license (MIT, GPL, etc.)

#### Architecture Support
- **64bit**: AMD64/x86_64 architecture support
- **arm64**: ARM64 architecture support
- **url**: Download URL for the specific architecture
- **hash**: SHA256 checksum for integrity verification

#### Binary Configuration
- **bin**: Array of executable files to add to PATH
- **shortcuts**: Optional desktop shortcuts
- **persist**: Directories to persist during updates

## Build Process

### Creating the Package

#### Prerequisites
```powershell
# Install required tools
scoop install git
scoop install 7zip
```

#### Manual Package Creation
```powershell
# Clone the bucket repository
git clone https://github.com/AlexBurnes/scoop-bucket.git
cd scoop-bucket

# Create or update the manifest
# Edit version.json with new version information

# Test the package locally
scoop install ./version.json

# Commit and push changes
git add version.json
git commit -m "Update version to 1.0.0"
git push origin main
```

### Integration with GoReleaser

GoReleaser can automatically update Scoop manifests:

```yaml
# .goreleaser.yml
scoop:
  bucket:
    owner: AlexBurnes
    name: scoop-bucket
  commit_author:
    name: GoReleaser Bot
    email: bot@example.com
  commit_msg_template: "Scoop: Update {{ .ProjectName }} to {{ .Tag }}"
  folder: "packaging/windows/scoop-bucket"
  homepage: "https://github.com/AlexBurnes/version-go"
  description: "git describe-like CLI"
  license: "MIT"
```

### Automated Build Process
1. **Build**: GoReleaser compiles Windows binaries
2. **Package**: Creates ZIP archives for each architecture
3. **Checksums**: Generates SHA256 checksums
4. **Manifest**: Updates version.json with new URLs and checksums
5. **Commit**: Pushes changes to the bucket repository
6. **Release**: Publishes GitHub release with all artifacts

## Security Features

### Integrity Verification
- **SHA256 Checksums**: Strong cryptographic verification
- **Automatic Validation**: Scoop verifies checksums before installation
- **Secure Downloads**: HTTPS-only downloads from GitHub releases
- **Error Handling**: Installation fails if checksums don't match

### Permission Management
- **User Scope**: Default installation in user directory
- **No Admin Required**: Standard installation doesn't need elevated privileges
- **Global Option**: Optional system-wide installation with admin rights
- **PATH Management**: Automatic PATH configuration

## Troubleshooting

### Common Issues

#### Scoop Not Found
```powershell
# Install Scoop
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

# Refresh environment
refreshenv
```

#### Bucket Not Found
```powershell
# Add the bucket
scoop bucket add burnes https://github.com/AlexBurnes/scoop-bucket

# List available buckets
scoop bucket list
```

#### Installation Fails
```powershell
# Check for errors
scoop install burnes/version --verbose

# Clear cache and retry
scoop cache rm version
scoop install burnes/version
```

#### Binary Not Found
```powershell
# Check PATH
$env:PATH

# Refresh environment
refreshenv

# Check installation
scoop list version
```

### Debug Mode
```powershell
# Install with verbose output
scoop install burnes/version --verbose

# Check installation details
scoop info version

# View installation directory
scoop which version
```

## Maintenance

### Updating the Package

#### Manual Update
```powershell
# Update version.json with new version
# Update URLs and checksums
# Commit and push changes
```

#### Automated Update (GoReleaser)
```powershell
# Tag new release
git tag v1.0.1
git push origin v1.0.1

# GoReleaser automatically updates the manifest
```

### Testing New Versions
```powershell
# Test locally before release
scoop install ./version.json

# Test functionality
version --version
version --help

# Uninstall test version
scoop uninstall version
```

## Best Practices

### Manifest Management
- Keep version.json up to date with releases
- Use semantic versioning consistently
- Include proper checksums for all architectures
- Test on both AMD64 and ARM64 systems

### User Experience
- Provide clear installation instructions
- Document prerequisites and dependencies
- Include helpful error messages
- Test installation process thoroughly

### Security
- Always verify checksums before installation
- Use HTTPS for all downloads
- Keep manifests in version control
- Test on clean systems to ensure no conflicts

## Advanced Configuration

### Custom Installation Directory
```json
{
  "architecture": {
    "64bit": {
      "url": "...",
      "hash": "...",
      "extract_dir": "custom-dir"
    }
  }
}
```

### Desktop Shortcuts
```json
{
  "shortcuts": [
    ["version.exe", "Version CLI"]
  ]
}
```

### Persistent Data
```json
{
  "persist": [
    "config",
    "data"
  ]
}
```

### Pre/Post Install Scripts
```json
{
  "pre_install": "echo 'Installing Version CLI...'",
  "post_install": "echo 'Installation complete!'"
}
```