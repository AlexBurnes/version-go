# BNF Grammar Specification

This document defines the Backus-Naur Form (BNF) grammar for the version parsing system implemented in the `version` CLI utility. The grammar extends SemVer 2.0 to support additional version types used in build pipelines.

## Grammar Definition

### Core Grammar Rules

```
<version> ::= <version-core> | <prerelease-version> | <postrelease-version> | <intermediate-version>

<version-core> ::= <major> "." <minor> "." <patch>
<major> ::= <numeric-identifier>
<minor> ::= <numeric-identifier>
<patch> ::= <numeric-identifier>

<prerelease-version> ::= <version-core> "~" <prerelease-identifier>
<prerelease-identifier> ::= <prerelease-type> <prerelease-suffix>?
<prerelease-type> ::= "alpha" | "beta" | "rc" | "pre"
<prerelease-suffix> ::= "." <numeric-identifier> | "_" <alphanumeric-identifier> | <prerelease-suffix> "." <numeric-identifier> | <prerelease-suffix> "_" <alphanumeric-identifier>

<postrelease-version> ::= <version-core> "." <postrelease-identifier>
<postrelease-identifier> ::= <postrelease-type> <postrelease-suffix>?
<postrelease-type> ::= "fix" | "next" | "post"
<postrelease-suffix> ::= "." <numeric-identifier> | "_" <alphanumeric-identifier> | <postrelease-suffix> "." <numeric-identifier> | <postrelease-suffix> "_" <alphanumeric-identifier>

<intermediate-version> ::= <version-core> "_" <intermediate-identifier>
<intermediate-identifier> ::= <alphanumeric-identifier> <intermediate-suffix>?
<intermediate-suffix> ::= "." <numeric-identifier> | "_" <alphanumeric-identifier> | <intermediate-suffix> "." <numeric-identifier> | <intermediate-suffix> "_" <alphanumeric-identifier>

<numeric-identifier> ::= "0" | <positive-digit> | <positive-digit> <digit>*
<positive-digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
<digit> ::= "0" | <positive-digit>

<alphanumeric-identifier> ::= <letter> | <letter> <alphanumeric-char>*
<letter> ::= "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z" | "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z"
<alphanumeric-char> ::= <letter> | <digit>
```

### Optional Version Prefix

```
<version-with-prefix> ::= "v" <version> | <version>
```

## Version Types and Precedence

The grammar defines four distinct version types with the following precedence order (lowest to highest):

1. **Prerelease Versions** (`<prerelease-version>`)
   - Format: `major.minor.patch~type[suffix]` or `v major.minor.patch~type[suffix]`
   - Examples: `1.2.3~alpha`, `1.2.3~beta.1`, `1.2.3~rc.1_feature`
   - Precedence: Lowest (comes before release version)

2. **Release Versions** (`<version-core>`)
   - Format: `major.minor.patch` or `v major.minor.patch`
   - Examples: `1.2.3`, `v1.2.3`
   - Precedence: Second (greater than prerelease)

3. **Postrelease Versions** (`<postrelease-version>`)
   - Format: `major.minor.patch.type[suffix]` or `v major.minor.patch.type[suffix]`
   - Examples: `1.2.3.fix`, `1.2.3.next.1`, `1.2.3.post.1_feature`
   - Precedence: Third (greater than release)

4. **Intermediate Versions** (`<intermediate-version>`)
   - Format: `major.minor.patch_identifier[suffix]` or `v major.minor.patch_identifier[suffix]`
   - Examples: `1.2.3_feature`, `1.2.3_exp.1`, `1.2.3_dev.1_feature`
   - Precedence: Highest (greatest of all types)

## Implementation Details

### Regex Patterns

The grammar is implemented using the following regular expressions in `pkg/version/version.go`:

```go
// Release version: v?[0-9]+\.[0-9]+\.[0-9]+
versionRelease = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)

// Prerelease version: v?[0-9]+\.[0-9]+\.[0-9]+~(alpha|beta|rc|pre)(...)*
versionPrerelease = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\~(alpha|beta|rc|pre)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)

// Postrelease version: v?[0-9]+\.[0-9]+\.[0-9]+\.(fix|next|post)(...)*
versionPostrelease = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\.(fix|next|post)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)

// Intermediate version: v?[0-9]+\.[0-9]+\.[0-9]+_[a-zA-Z]+(...)*
versionIntermediate = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)\_([a-zA-Z]+)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
```

### Git Tag Conversion

The implementation includes automatic conversion of git tags that use `-` delimiter instead of `~` for prerelease versions:

- Input: `1.2.3-alpha.1` → Output: `1.2.3~alpha.1`
- Input: `v1.2.3-beta` → Output: `v1.2.3~beta`

This conversion is handled by the `ConvertGitTag()` function to maintain compatibility with existing git tag conventions.

## Sorting Rules

### Version Comparison Algorithm

1. **Core Version Comparison**: Compare `major.minor.patch` numerically
2. **Type Precedence**: Compare version types according to precedence order
3. **Identifier Comparison**: For same type, compare identifiers using:
   - Numeric identifiers: Compare numerically (`0 < 1 < 2 < ... < 10`)
   - Alphanumeric identifiers: Compare lexically (`a < b < ... < z < A < ... < Z`)
   - Mixed identifiers: Numbers come before letters

### Examples

```
Input versions:
1.2.3
1.2.3_feature
1.2.3~alpha
1.2.3.fix
2.0.0
1.2.3~beta.1
1.2.3~alpha.2

Sorted output:
1.2.3~alpha      (prerelease - lowest precedence)
1.2.3~alpha.2    (prerelease with suffix)
1.2.3~beta.1     (prerelease with different type)
1.2.3            (release - greater than prerelease)
1.2.3.fix        (postrelease - greater than release)
1.2.3_feature    (intermediate - highest precedence for same x.y.z)
2.0.0            (release with higher core version)
```

## Validation Rules

### Valid Version Strings

- `1.2.3` - Standard release version
- `v1.2.3` - Release version with prefix
- `1.2.3~alpha` - Prerelease version
- `1.2.3~beta.1` - Prerelease with numeric suffix
- `1.2.3~rc.1_feature` - Prerelease with mixed suffix
- `1.2.3.fix` - Postrelease version
- `1.2.3.next.1` - Postrelease with numeric suffix
- `1.2.3.post.1_feature` - Postrelease with mixed suffix
- `1.2.3_feature` - Intermediate version
- `1.2.3_exp.1` - Intermediate with numeric suffix
- `1.2.3_dev.1_feature` - Intermediate with mixed suffix

### Invalid Version Strings

- `1.2` - Missing patch version
- `1.2.3-alpha` - Missing tilde for prerelease
- `1.2.3.alpha` - Wrong delimiter for prerelease
- `1.2.3_` - Empty intermediate identifier
- `1.2.3~` - Empty prerelease identifier
- `1.2.3.` - Empty postrelease identifier
- `01.2.3` - Leading zeros in major version
- `1.02.3` - Leading zeros in minor version
- `1.2.03` - Leading zeros in patch version

## References

- [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
- [RPM Naming Conventions](http://ftp.rpm.org/max-rpm/ch-rpm-file-format.html)
- [Fedora Tilde Versioning](https://fedoraproject.org/wiki/PackagingDrafts/TildeVersioning)