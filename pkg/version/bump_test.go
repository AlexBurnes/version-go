package version

import (
	"testing"
)

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release version", "1.2.3", "2.0.0"},
		{"prerelease version", "1.2.3~alpha.1", "2.0.0"},
		{"postrelease version", "1.2.3.fix.1", "2.0.0"},
		{"intermediate version", "1.2.3_feature.1", "2.0.0"},
		{"version with v prefix", "v1.2.3", "2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpMajor)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpMajor(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpMinor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release version", "1.2.3", "1.3.0"},
		{"prerelease version", "1.2.3~alpha.1", "1.3.0"},
		{"postrelease version", "1.2.3.fix.1", "1.3.0"},
		{"intermediate version", "1.2.3_feature.1", "1.3.0"},
		{"version with v prefix", "v1.2.3", "1.3.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpMinor)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpMinor(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release version", "1.2.3", "1.2.4"},
		{"prerelease version", "1.2.3~alpha.1", "1.2.4"},
		{"postrelease version", "1.2.3.fix.1", "1.2.4"},
		{"intermediate version", "1.2.3_feature.1", "1.2.4"},
		{"version with v prefix", "v1.2.3", "1.2.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpPatch)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpPatch(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release to prerelease", "1.2.3", "1.2.3~alpha.1"},
		{"increment prerelease", "1.2.3~alpha.1", "1.2.3~alpha.2"},
		{"increment prerelease with complex suffix", "1.2.3~alpha.1_feature", "1.2.3~alpha_feature.1"},
		{"increment prerelease with multiple numbers", "1.2.3~alpha.1.2", "1.2.3~alpha.3"},
		{"prerelease with v prefix", "v1.2.3", "1.2.3~alpha.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpAlpha)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpAlpha(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpFix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release to postrelease", "1.2.3", "1.2.3.fix.1"},
		{"increment postrelease", "1.2.3.fix.1", "1.2.3.fix.2"},
		{"increment postrelease with complex suffix", "1.2.3.fix.1_feature", "1.2.3.fix_feature.1"},
		{"increment postrelease with multiple numbers", "1.2.3.fix.1.2", "1.2.3.fix.3"},
		{"postrelease with v prefix", "v1.2.3", "1.2.3.fix.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpFix)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpFix(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpFeat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release to intermediate", "1.2.3", "1.2.3_feat.1"},
		{"increment intermediate", "1.2.3_feat.1", "1.2.3_feat.2"},
		{"increment intermediate with complex suffix", "1.2.3_feat.1_dev", "1.2.3_feat_dev.1"},
		{"increment intermediate with multiple numbers", "1.2.3_feat.1.2", "1.2.3_feat.3"},
		{"intermediate with v prefix", "v1.2.3", "1.2.3_feat.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpFeat)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpFeat(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestBumpSmart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"release version", "1.2.3", "1.2.4"},
		{"prerelease version", "1.2.3~alpha.1", "1.2.3~alpha.2"},
		{"postrelease version", "1.2.3.fix.1", "1.2.3.fix.2"},
		{"intermediate version", "1.2.3_feature.1", "1.2.3_feature.2"},
		{"version with v prefix", "v1.2.3", "1.2.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bump(tt.input, BumpSmart)
			if err != nil {
				t.Fatalf("Bump failed: %v", err)
			}
			if result.BumpedVersion != tt.expected {
				t.Errorf("BumpSmart(%s) = %s, want %s", tt.input, result.BumpedVersion, tt.expected)
			}
		})
	}
}

func TestParseBumpType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BumpType
		hasError bool
	}{
		{"major", "major", BumpMajor, false},
		{"minor", "minor", BumpMinor, false},
		{"patch", "patch", BumpPatch, false},
		{"pre", "pre", BumpPre, false},
		{"alpha", "alpha", BumpAlpha, false},
		{"beta", "beta", BumpBeta, false},
		{"rc", "rc", BumpRc, false},
		{"fix", "fix", BumpFix, false},
		{"next", "next", BumpNext, false},
		{"post", "post", BumpPost, false},
		{"feat", "feat", BumpFeat, false},
		{"smart", "smart", BumpSmart, false},
		{"uppercase", "MAJOR", BumpMajor, false},
		{"mixed case", "Major", BumpMajor, false},
		{"invalid type", "invalid", BumpSmart, true},
		{"empty string", "", BumpSmart, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBumpType(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ParseBumpType(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseBumpType(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseBumpType(%s) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestBumpResult(t *testing.T) {
	result, err := Bump("1.2.3", BumpMajor)
	if err != nil {
		t.Fatalf("Bump failed: %v", err)
	}

	if result.OriginalVersion != "1.2.3" {
		t.Errorf("OriginalVersion = %s, want 1.2.3", result.OriginalVersion)
	}
	if result.BumpedVersion != "2.0.0" {
		t.Errorf("BumpedVersion = %s, want 2.0.0", result.BumpedVersion)
	}
	if result.BumpType != BumpMajor {
		t.Errorf("BumpType = %v, want BumpMajor", result.BumpType)
	}
	if result.AppliedRule == "" {
		t.Error("AppliedRule should not be empty")
	}
}

func TestBumpInvalidVersion(t *testing.T) {
	_, err := Bump("invalid-version", BumpMajor)
	if err == nil {
		t.Error("Bump with invalid version should return error")
	}
}

func TestBumpTypeString(t *testing.T) {
	tests := []struct {
		bumpType BumpType
		expected string
	}{
		{BumpMajor, "major"},
		{BumpMinor, "minor"},
		{BumpPatch, "patch"},
		{BumpPre, "pre"},
		{BumpAlpha, "alpha"},
		{BumpBeta, "beta"},
		{BumpRc, "rc"},
		{BumpFix, "fix"},
		{BumpNext, "next"},
		{BumpPost, "post"},
		{BumpFeat, "feat"},
		{BumpSmart, "smart"},
		{BumpType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.bumpType.String()
			if result != tt.expected {
				t.Errorf("BumpType.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}