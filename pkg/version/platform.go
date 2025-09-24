// Package version provides platform detection functionality for cross-platform compatibility.
// It detects the current platform, architecture, operating system, and OS version using Go's built-in runtime capabilities.
package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// PlatformInfo represents information about the current platform
type PlatformInfo struct {
	Platform  string // e.g., "linux", "windows", "darwin"
	Arch      string // e.g., "amd64", "arm64", "386"
	OS        string // e.g., "ubuntu", "windows", "darwin"
	OSVersion string // e.g., "24.04", "windows", "darwin"
	NumCPU    int    // number of logical CPUs
}

// GetPlatform returns the current platform name (GOOS value)
func GetPlatform() string {
	return runtime.GOOS
}

// GetArch returns the current architecture name (GOARCH value)
func GetArch() string {
	return runtime.GOARCH
}

// GetOS returns the current operating system name (user-friendly format)
func GetOS() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxDistribution()
	case "windows":
		return "windows"
	case "darwin":
		return "darwin"
	default:
		return runtime.GOOS
	}
}

// GetOSVersion returns the current operating system version
// This is a simplified implementation that returns basic version info
func GetOSVersion() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxVersion()
	case "windows":
		return getWindowsVersion()
	case "darwin":
		return getDarwinVersion()
	default:
		return "unknown"
	}
}

// GetPlatformInfo returns comprehensive platform information
func GetPlatformInfo() *PlatformInfo {
	return &PlatformInfo{
		Platform:  GetPlatform(),
		Arch:      GetArch(),
		OS:        GetOS(),
		OSVersion: GetOSVersion(),
		NumCPU:    GetNumCPU(),
	}
}

// GetNumCPU returns the number of logical CPUs usable by the current process
func GetNumCPU() int {
	return runtime.NumCPU()
}

// getLinuxVersion attempts to detect Linux distribution and version
func getLinuxVersion() string {
	// Try to read from /etc/os-release first
	if version := readOSReleaseVersion(); version != "" {
		return version
	}
	
	// Fallback to generic linux
	return "linux"
}

// getWindowsVersion attempts to detect Windows version
func getWindowsVersion() string {
	// Use Go's built-in capabilities to detect Windows version
	// runtime.GOOS already gives us "windows", but we can get more specific info
	
	// For now, return generic windows version
	// In a full implementation, this could use syscalls to get actual version
	return "windows"
}

// getDarwinVersion attempts to detect Darwin version
func getDarwinVersion() string {
	// Use Go's built-in capabilities to detect Darwin version
	// runtime.GOOS already gives us "darwin", but we can get more specific info
	
	// Check for macOS version using environment variables
	if version := os.Getenv("MACOSX_DEPLOYMENT_TARGET"); version != "" {
		return version
	}
	
	// For now, return generic darwin version
	// In a full implementation, this could use syscalls to get actual version
	return "darwin"
}

// getLinuxDistribution returns the Linux distribution name (e.g., "ubuntu", "debian")
func getLinuxDistribution() string {
	// Use Go's built-in file reading capabilities
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "linux"
	}
	
	// Parse the os-release file
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID=") {
			name := strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			return strings.ToLower(name)
		}
	}
	
	return "linux"
}

// readOSReleaseVersion reads Linux distribution version from /etc/os-release
func readOSReleaseVersion() string {
	// Use Go's built-in file reading capabilities
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	
	// Parse the os-release file
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "VERSION_ID=") {
			version := strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			return version
		}
	}
	
	return ""
}

// readOSRelease reads Linux distribution information from /etc/os-release (legacy function)
func readOSRelease() string {
	// Use Go's built-in file reading capabilities
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	
	// Parse the os-release file
	lines := strings.Split(string(content), "\n")
	var name, version string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}
	
	if name != "" && version != "" {
		return fmt.Sprintf("%s-%s", strings.ToLower(name), version)
	} else if name != "" {
		return strings.ToLower(name)
	}
	
	return ""
}

// String returns a string representation of PlatformInfo
func (p *PlatformInfo) String() string {
	return fmt.Sprintf("%s-%s (%s %s, %d CPUs)", p.Platform, p.Arch, p.OS, p.OSVersion, p.NumCPU)
}

// FormatPlatformArch returns platform-arch format (e.g., "linux-amd64")
func (p *PlatformInfo) FormatPlatformArch() string {
	return fmt.Sprintf("%s-%s", p.Platform, p.Arch)
}

// FormatOSArch returns OS-arch format (e.g., "linux-amd64", "macos-arm64")
func (p *PlatformInfo) FormatOSArch() string {
	return fmt.Sprintf("%s-%s", p.OS, p.Arch)
}

// IsLinux returns true if the platform is Linux
func (p *PlatformInfo) IsLinux() bool {
	return p.Platform == "linux"
}

// IsWindows returns true if the platform is Windows
func (p *PlatformInfo) IsWindows() bool {
	return p.Platform == "windows"
}

// IsDarwin returns true if the platform is Darwin
func (p *PlatformInfo) IsDarwin() bool {
	return p.Platform == "darwin"
}

// IsMacOS returns true if the platform is macOS (alias for IsDarwin for compatibility)
func (p *PlatformInfo) IsMacOS() bool {
	return p.Platform == "darwin"
}

// IsAMD64 returns true if the architecture is amd64
func (p *PlatformInfo) IsAMD64() bool {
	return p.Arch == "amd64"
}

// IsARM64 returns true if the architecture is arm64
func (p *PlatformInfo) IsARM64() bool {
	return p.Arch == "arm64"
}

// Is386 returns true if the architecture is 386
func (p *PlatformInfo) Is386() bool {
	return p.Arch == "386"
}

// ValidatePlatform validates that a platform string is supported
func ValidatePlatform(platform string) error {
	supportedPlatforms := []string{"linux", "windows", "darwin"}
	platform = strings.ToLower(platform)
	
	for _, supported := range supportedPlatforms {
		if platform == supported {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported platform: %s (supported: %s)", platform, strings.Join(supportedPlatforms, ", "))
}

// ValidateArch validates that an architecture string is supported
func ValidateArch(arch string) error {
	supportedArchs := []string{"amd64", "arm64", "386"}
	arch = strings.ToLower(arch)
	
	for _, supported := range supportedArchs {
		if arch == supported {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported architecture: %s (supported: %s)", arch, strings.Join(supportedArchs, ", "))
}

// GetSupportedPlatforms returns a list of supported platforms
func GetSupportedPlatforms() []string {
	return []string{"linux", "windows", "darwin"}
}

// GetSupportedArchs returns a list of supported architectures
func GetSupportedArchs() []string {
	return []string{"amd64", "arm64", "386"}
}

// GetSupportedPlatformArchs returns all supported platform-architecture combinations
func GetSupportedPlatformArchs() []string {
	platforms := GetSupportedPlatforms()
	archs := GetSupportedArchs()
	
	var combinations []string
	for _, platform := range platforms {
		for _, arch := range archs {
			// Skip unsupported combinations
			if platform == "darwin" && arch == "386" {
				continue // macOS doesn't support 386
			}
			combinations = append(combinations, fmt.Sprintf("%s-%s", platform, arch))
		}
	}
	
	return combinations
}
