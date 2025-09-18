# New Features Specification

## Overview
This document specifies two new major features for the version CLI utility:
1. `.project.yml` configuration file support
2. `bump` command with intelligent version increment

## 1. .project.yml Configuration File Support

### Purpose
Allow projects to define their name and module names in a configuration file instead of relying solely on git for project information. This provides consistency across build utilities and better control over project naming.

### File Format
**Location**: Project root directory (`.project.yml`)

**Structure**:
```yaml
project:
  name: "project-name"
  modules:
    - "main-module-name"    # First is primary
    - "secondary-module"
    - "another-module"
```

### Behavior
- **Primary Source**: If `.project.yml` exists, use it for project and module information
- **Fallback**: If `.project.yml` doesn't exist, fall back to git-based detection
- **Integration**: Other build utilities can also use this file for consistent project naming
- **Commands Affected**: `project`, `module` commands will use this configuration

### Implementation Details
- YAML parsing using Go standard library or minimal YAML library
- File existence check in project root directory
- Graceful fallback to existing git-based detection
- Error handling for malformed YAML files

## 2. Bump Command Implementation

### Purpose
Provide intelligent version bumping functionality that understands the current version state and applies appropriate increment rules.

### Usage
```bash
version bump [version-type]
```

### Version Types
- `major` - Increment major version, reset minor=0, patch=0, clear other identifiers
- `minor` - Increment minor version, reset patch=0, clear other identifiers
- `patch`/`release` - Increment patch version, clear other identifiers
- `pre[release]` - Add/increment prerelease identifier
- `post[release]`/`next`/`fix` - Add/increment postrelease identifier
- `inter[mediate]`/`feat[ure]` - Add/increment intermediate identifier
- No argument - Smart increment of last version number

### Bump Rules Matrix

Based on the provided image, the following rules apply:

| Bump Command | Current State | Result |
|--------------|---------------|---------|
| (no argument) | major, basically patch | Patch +1 |
| (no argument) | minor, basically patch | Patch +1 |
| (no argument) | patch\|release | Patch +1 |
| (no argument) | release | clear other |
| (no argument) | pre | clear other |
| (no argument) | postrelease | Patch+1, clear other |
| (no argument) | intermediate | Patch+1, clear other |
| major | any | major +1, minor=0, patch=0, clean all |
| minor | major, basically patch | minor +1, patch=0 |
| minor | minor, basically patch | minor +1, patch=0 |
| minor | patch\|release | minor +1, patch=0 |
| minor | release | minor +1, patch=0, clean other |
| minor | pre | minor +1, patch=0, clean other |
| minor | postrelease | minor +1, patch=0, clean other |
| minor | intermediate | minor +1, patch=0, clean other |
| patch\|release | major, basically patch | Patch +1 |
| patch\|release | minor, basically patch | Patch +1 |
| patch\|release | patch\|release | Patch +1 |
| patch\|release | release | clear other |
| patch\|release | pre | clear other |
| patch\|release | postrelease | Patch+1 (if not pre in version), clear other |
| patch\|release | intermediate | Patch+1 (if not pre in version), clear other |
| pre[release] | patch\|release | Patch +1, add pre.1 |
| pre[release] | release | pre +1 |
| pre[release] | pre | pre +1 |
| pre[release] | postrelease | Patch+1 (if not pre in version), pre.1 (if not pre in version) |
| pre[release] | intermediate | Patch+1 (if not pre in version), pre.1 (if not pre in version) |
| post[release]\|next\|fix | patch\|release | add post.1 |
| post[release]\|next\|fix | release | add post.1 |
| post[release]\|next\|fix | pre | add post.1 |
| post[release]\|next\|fix | postrelease | Post +1 |
| post[release]\|next\|fix | intermediate | add post.1 |
| inter[mediate]\|feat[ure] | patch\|release | add feat.1 |
| inter[mediate]\|feat[ure] | release | add feat.1 |
| inter[mediate]\|feat[ure] | pre | add feat.1 |
| inter[mediate]\|feat[ure] | postrelease | add feat.1 |
| inter[mediate]\|feat[ure] | intermediate | feat +1 |

### Examples
- `v0.1.1` → `v0.1.2` (no argument)
- `v0.1.1~pre.1` → `v0.1.1~pre.2` (no argument)
- `v1.0.0` → `v2.0.0` (major)
- `v1.0.0` → `v1.1.0` (minor)
- `v1.0.0` → `v1.0.1` (patch)
- `v1.0.0` → `v1.0.1~pre.1` (pre)
- `v1.0.0` → `v1.0.1.post.1` (post)

### Constraints
- **Pre delimiter**: Pre must be only one after patch version, delimited by `~`
- **Infinite blocks**: Number of blocks after release/prerelease for postmediate/intermediate is infinite
- **Version state detection**: Must correctly identify current version state (major, minor, patch, release, pre, postrelease, intermediate)

### Implementation Details
- Version state detection algorithm
- Bump rule engine with state-based logic
- Smart increment when no version type specified
- Integration with existing version parsing and validation
- Error handling for invalid version states or bump types

## Integration Points

### Existing Commands
- `project` command will check `.project.yml` first, then fall back to git
- `module` command will check `.project.yml` first, then fall back to git
- New `bump` command added to CLI interface

### Library Package
- New functions in `pkg/version` for configuration file parsing
- New functions in `pkg/version` for version bumping logic
- Integration with existing version parsing and validation

### Build System
- No changes to existing build system
- New features are additive and don't affect existing functionality

## Testing Requirements

### .project.yml Feature
- Test with valid `.project.yml` file
- Test with invalid `.project.yml` file (malformed YAML)
- Test with missing `.project.yml` file (fallback to git)
- Test with multiple modules configuration
- Test integration with `project` and `module` commands

### Bump Command Feature
- Test all bump types with different version states
- Test smart increment (no argument) with different version states
- Test edge cases and error conditions
- Test integration with existing version parsing
- Test output format and validation

## Future Considerations
- Configuration file could be extended for other project settings
- Bump command could be extended with additional version types
- Integration with CI/CD systems for automated version bumping
- Support for custom bump rules or configuration