package version

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	platform := GetPlatform()
	
	// Platform should match runtime.GOOS
	if platform != runtime.GOOS {
		t.Errorf("GetPlatform() = %v, want %v", platform, runtime.GOOS)
	}
	
	// Platform should be one of the supported platforms
	supportedPlatforms := GetSupportedPlatforms()
	found := false
	for _, supported := range supportedPlatforms {
		if platform == supported {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetPlatform() returned unsupported platform: %v", platform)
	}
}

func TestGetArch(t *testing.T) {
	arch := GetArch()
	
	// Arch should match runtime.GOARCH
	if arch != runtime.GOARCH {
		t.Errorf("GetArch() = %v, want %v", arch, runtime.GOARCH)
	}
	
	// Arch should be one of the supported architectures
	supportedArchs := GetSupportedArchs()
	found := false
	for _, supported := range supportedArchs {
		if arch == supported {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetArch() returned unsupported architecture: %v", arch)
	}
}

func TestGetOS(t *testing.T) {
	os := GetOS()
	
	// OS should be user-friendly format
	switch runtime.GOOS {
	case "linux":
		// For linux, OS should be the distribution name (e.g., "ubuntu", "debian") or "linux" as fallback
		if os != "linux" && os != "ubuntu" && os != "debian" && os != "centos" && os != "fedora" {
			t.Errorf("GetOS() = %v, should be 'linux' or a known distribution name", os)
		}
	case "windows":
		if os != "windows" {
			t.Errorf("GetOS() = %v, want windows", os)
		}
	case "darwin":
		if os != "darwin" {
			t.Errorf("GetOS() = %v, want darwin", os)
		}
	default:
		if os != runtime.GOOS {
			t.Errorf("GetOS() = %v, want %v", os, runtime.GOOS)
		}
	}
}

func TestGetOSVersion(t *testing.T) {
	osVersion := GetOSVersion()
	
	// OSVersion should not be empty
	if osVersion == "" {
		t.Errorf("GetOSVersion() returned empty string")
	}
	
	// OSVersion should contain expected values based on platform
	// For linux, it should be just the version number (e.g., "24.04", "12") or "linux" as fallback
	switch runtime.GOOS {
	case "linux":
		if osVersion != "linux" && !strings.Contains(osVersion, ".") && !strings.Contains(osVersion, "-") {
			// Should be a version number like "24.04" or "12" or fallback to "linux"
			if osVersion != "linux" {
				t.Errorf("GetOSVersion() = %v, should be 'linux' or a version number", osVersion)
			}
		}
	case "windows":
		if !strings.HasPrefix(osVersion, "windows") {
			t.Errorf("GetOSVersion() = %v, should start with 'windows'", osVersion)
		}
	case "darwin":
		// On Darwin/macOS, GetOSVersion() should return actual version number (e.g., "15.6.1")
		// or "darwin" as fallback, not necessarily starting with "darwin"
		if osVersion == "" {
			t.Errorf("GetOSVersion() returned empty string")
		}
	}
}

func TestGetPlatformInfo(t *testing.T) {
	info := GetPlatformInfo()
	
	// All fields should be populated
	if info.Platform == "" {
		t.Errorf("PlatformInfo.Platform is empty")
	}
	if info.Arch == "" {
		t.Errorf("PlatformInfo.Arch is empty")
	}
	if info.OS == "" {
		t.Errorf("PlatformInfo.OS is empty")
	}
	if info.OSVersion == "" {
		t.Errorf("PlatformInfo.OSVersion is empty")
	}
	if info.NumCPU < 1 {
		t.Errorf("PlatformInfo.NumCPU is invalid: %d", info.NumCPU)
	}
	
	// Platform should match runtime.GOOS
	if info.Platform != runtime.GOOS {
		t.Errorf("PlatformInfo.Platform = %v, want %v", info.Platform, runtime.GOOS)
	}
	
	// Arch should match runtime.GOARCH
	if info.Arch != runtime.GOARCH {
		t.Errorf("PlatformInfo.Arch = %v, want %v", info.Arch, runtime.GOARCH)
	}
	
	// OS should match GetOS()
	if info.OS != GetOS() {
		t.Errorf("PlatformInfo.OS = %v, want %v", info.OS, GetOS())
	}
	
	// OSVersion should match GetOSVersion()
	if info.OSVersion != GetOSVersion() {
		t.Errorf("PlatformInfo.OSVersion = %v, want %v", info.OSVersion, GetOSVersion())
	}
	
	// NumCPU should match GetNumCPU()
	if info.NumCPU != GetNumCPU() {
		t.Errorf("PlatformInfo.NumCPU = %v, want %v", info.NumCPU, GetNumCPU())
	}
}

func TestPlatformInfoString(t *testing.T) {
	info := GetPlatformInfo()
	str := info.String()
	
	// String should contain platform and arch
	if !strings.Contains(str, info.Platform) {
		t.Errorf("PlatformInfo.String() = %v, should contain platform %v", str, info.Platform)
	}
	if !strings.Contains(str, info.Arch) {
		t.Errorf("PlatformInfo.String() = %v, should contain arch %v", str, info.Arch)
	}
	if !strings.Contains(str, info.OS) {
		t.Errorf("PlatformInfo.String() = %v, should contain OS %v", str, info.OS)
	}
	if !strings.Contains(str, info.OSVersion) {
		t.Errorf("PlatformInfo.String() = %v, should contain OSVersion %v", str, info.OSVersion)
	}
	if !strings.Contains(str, fmt.Sprintf("%d", info.NumCPU)) {
		t.Errorf("PlatformInfo.String() = %v, should contain NumCPU %v", str, info.NumCPU)
	}
}

func TestFormatPlatformArch(t *testing.T) {
	info := GetPlatformInfo()
	formatted := info.FormatPlatformArch()
	
	expected := info.Platform + "-" + info.Arch
	if formatted != expected {
		t.Errorf("FormatPlatformArch() = %v, want %v", formatted, expected)
	}
}

func TestFormatOSArch(t *testing.T) {
	info := GetPlatformInfo()
	formatted := info.FormatOSArch()
	
	expected := info.OS + "-" + info.Arch
	if formatted != expected {
		t.Errorf("FormatOSArch() = %v, want %v", formatted, expected)
	}
}

func TestPlatformDetectionMethods(t *testing.T) {
	info := GetPlatformInfo()
	
	// Test platform detection methods
	switch runtime.GOOS {
	case "linux":
		if !info.IsLinux() {
			t.Errorf("IsLinux() should return true for linux platform")
		}
		if info.IsWindows() || info.IsMacOS() {
			t.Errorf("IsWindows() or IsMacOS() should return false for linux platform")
		}
	case "windows":
		if !info.IsWindows() {
			t.Errorf("IsWindows() should return true for windows platform")
		}
		if info.IsLinux() || info.IsMacOS() {
			t.Errorf("IsLinux() or IsMacOS() should return false for windows platform")
		}
	case "darwin":
		if !info.IsDarwin() {
			t.Errorf("IsDarwin() should return true for darwin platform")
		}
		if info.IsLinux() || info.IsWindows() {
			t.Errorf("IsLinux() or IsWindows() should return false for darwin platform")
		}
	}
}

func TestArchDetectionMethods(t *testing.T) {
	info := GetPlatformInfo()
	
	// Test architecture detection methods
	switch runtime.GOARCH {
	case "amd64":
		if !info.IsAMD64() {
			t.Errorf("IsAMD64() should return true for amd64 architecture")
		}
		if info.IsARM64() || info.Is386() {
			t.Errorf("IsARM64() or Is386() should return false for amd64 architecture")
		}
	case "arm64":
		if !info.IsARM64() {
			t.Errorf("IsARM64() should return true for arm64 architecture")
		}
		if info.IsAMD64() || info.Is386() {
			t.Errorf("IsAMD64() or Is386() should return false for arm64 architecture")
		}
	case "386":
		if !info.Is386() {
			t.Errorf("Is386() should return true for 386 architecture")
		}
		if info.IsAMD64() || info.IsARM64() {
			t.Errorf("IsAMD64() or IsARM64() should return false for 386 architecture")
		}
	}
}

func TestValidatePlatform(t *testing.T) {
	tests := []struct {
		platform string
		wantErr  bool
	}{
		{"linux", false},
		{"windows", false},
		{"darwin", false},
		{"Linux", false}, // case insensitive
		{"WINDOWS", false},
		{"Darwin", false},
		{"invalid", true},
		{"", true},
		{"macos", true}, // should be "darwin"
	}
	
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			err := ValidatePlatform(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlatform(%v) error = %v, wantErr %v", tt.platform, err, tt.wantErr)
			}
		})
	}
}

func TestValidateArch(t *testing.T) {
	tests := []struct {
		arch    string
		wantErr bool
	}{
		{"amd64", false},
		{"arm64", false},
		{"386", false},
		{"AMD64", false}, // case insensitive
		{"ARM64", false},
		{"386", false},
		{"invalid", true},
		{"", true},
		{"x86_64", true}, // should be "amd64"
	}
	
	for _, tt := range tests {
		t.Run(tt.arch, func(t *testing.T) {
			err := ValidateArch(tt.arch)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArch(%v) error = %v, wantErr %v", tt.arch, err, tt.wantErr)
			}
		})
	}
}

func TestGetSupportedPlatforms(t *testing.T) {
	platforms := GetSupportedPlatforms()
	
	expectedPlatforms := []string{"linux", "windows", "darwin"}
	if len(platforms) != len(expectedPlatforms) {
		t.Errorf("GetSupportedPlatforms() length = %v, want %v", len(platforms), len(expectedPlatforms))
	}
	
	for _, expected := range expectedPlatforms {
		found := false
		for _, platform := range platforms {
			if platform == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetSupportedPlatforms() missing platform: %v", expected)
		}
	}
}

func TestGetSupportedArchs(t *testing.T) {
	archs := GetSupportedArchs()
	
	expectedArchs := []string{"amd64", "arm64", "386"}
	if len(archs) != len(expectedArchs) {
		t.Errorf("GetSupportedArchs() length = %v, want %v", len(archs), len(expectedArchs))
	}
	
	for _, expected := range expectedArchs {
		found := false
		for _, arch := range archs {
			if arch == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetSupportedArchs() missing architecture: %v", expected)
		}
	}
}

func TestGetSupportedPlatformArchs(t *testing.T) {
	combinations := GetSupportedPlatformArchs()
	
	// Should have combinations for all platforms and architectures
	// except darwin-386 which is not supported
	expectedCount := 3*3 - 1 // 3 platforms * 3 archs - 1 unsupported combination
	if len(combinations) != expectedCount {
		t.Errorf("GetSupportedPlatformArchs() length = %v, want %v", len(combinations), expectedCount)
	}
	
	// Should not contain darwin-386
	for _, combo := range combinations {
		if combo == "darwin-386" {
			t.Errorf("GetSupportedPlatformArchs() should not contain darwin-386")
		}
	}
	
	// Should contain expected combinations
	expectedCombinations := []string{
		"linux-amd64", "linux-arm64", "linux-386",
		"windows-amd64", "windows-arm64", "windows-386",
		"darwin-amd64", "darwin-arm64",
	}
	
	for _, expected := range expectedCombinations {
		found := false
		for _, combo := range combinations {
			if combo == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetSupportedPlatformArchs() missing combination: %v", expected)
		}
	}
}

func TestGetNumCPU(t *testing.T) {
	numCPU := GetNumCPU()
	
	// Should have at least 1 CPU
	if numCPU < 1 {
		t.Errorf("GetNumCPU() = %v, should be >= 1", numCPU)
	}
}
