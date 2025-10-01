package version

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected *Version
		hasError bool
	}{
		// Release versions
		{
			input: "1.2.3",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypeRelease,
				Original: "1.2.3",
			},
			hasError: false,
		},
		{
			input: "v1.2.3",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypeRelease,
				Original: "v1.2.3",
			},
			hasError: false,
		},
		{
			input: "0.0.0",
			expected: &Version{
				Major: 0, Minor: 0, Patch: 0,
				Type: TypeRelease,
				Original: "0.0.0",
			},
			hasError: false,
		},
		
		// Prerelease versions
		{
			input: "1.2.3-alpha",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypePrerelease,
				Prerelease: "~alpha",
				Original: "1.2.3-alpha",
			},
			hasError: false,
		},
		{
			input: "1.2.3~beta.1",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypePrerelease,
				Prerelease: "~beta.1",
				Original: "1.2.3~beta.1",
			},
			hasError: false,
		},
		{
			input: "v1.2.3-rc.1",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypePrerelease,
				Prerelease: "~rc.1",
				Original: "v1.2.3-rc.1",
			},
			hasError: false,
		},
		
		// Postrelease versions
		{
			input: "1.2.3.fix",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypePostrelease,
				Postrelease: ".fix",
				Original: "1.2.3.fix",
			},
			hasError: false,
		},
		{
			input: "1.2.3.post.1",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypePostrelease,
				Postrelease: ".post.1",
				Original: "1.2.3.post.1",
			},
			hasError: false,
		},
		
		// Intermediate versions
		{
			input: "1.2.3_feature",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypeIntermediate,
				Intermediate: "_feature",
				Original: "1.2.3_feature",
			},
			hasError: false,
		},
		{
			input: "1.2.3_exp.1",
			expected: &Version{
				Major: 1, Minor: 2, Patch: 3,
				Type: TypeIntermediate,
				Intermediate: "_exp.1",
				Original: "1.2.3_exp.1",
			},
			hasError: false,
		},
		
		// Invalid versions
		{
			input: "1.2",
			hasError: true,
		},
		{
			input: "1.2.3.4",
			hasError: true,
		},
		{
			input: "1.2.3-invalid",
			hasError: true,
		},
		{
			input: "invalid",
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Parse(test.input)
			
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", test.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
				return
			}
			
			if result.Major != test.expected.Major {
				t.Errorf("Major mismatch for %s: expected %d, got %d", test.input, test.expected.Major, result.Major)
			}
			if result.Minor != test.expected.Minor {
				t.Errorf("Minor mismatch for %s: expected %d, got %d", test.input, test.expected.Minor, result.Minor)
			}
			if result.Patch != test.expected.Patch {
				t.Errorf("Patch mismatch for %s: expected %d, got %d", test.input, test.expected.Patch, result.Patch)
			}
			if result.Type != test.expected.Type {
				t.Errorf("Type mismatch for %s: expected %v, got %v", test.input, test.expected.Type, result.Type)
			}
			if result.Prerelease != test.expected.Prerelease {
				t.Errorf("Prerelease mismatch for %s: expected %s, got %s", test.input, test.expected.Prerelease, result.Prerelease)
			}
			if result.Postrelease != test.expected.Postrelease {
				t.Errorf("Postrelease mismatch for %s: expected %s, got %s", test.input, test.expected.Postrelease, result.Postrelease)
			}
			if result.Intermediate != test.expected.Intermediate {
				t.Errorf("Intermediate mismatch for %s: expected %s, got %s", test.input, test.expected.Intermediate, result.Intermediate)
			}
		})
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		vt       Type
		expected string
	}{
		{TypeRelease, "release"},
		{TypePrerelease, "prerelease"},
		{TypePostrelease, "postrelease"},
		{TypeIntermediate, "intermediate"},
		{TypeInvalid, "invalid"},
	}

	for _, test := range tests {
		if test.vt.String() != test.expected {
			t.Errorf("Type.String() for %v: expected %s, got %s", test.vt, test.expected, test.vt.String())
		}
	}
}

func TestTypeBuildType(t *testing.T) {
	tests := []struct {
		vt       Type
		expected string
	}{
		{TypeRelease, "Release"},
		{TypePrerelease, "Debug"},
		{TypePostrelease, "Debug"},
		{TypeIntermediate, "Debug"},
		{TypeInvalid, "Debug"},
	}

	for _, test := range tests {
		if test.vt.BuildType() != test.expected {
			t.Errorf("Type.BuildType() for %v: expected %s, got %s", test.vt, test.expected, test.vt.BuildType())
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        *Version
		b        *Version
		expected int
	}{
		{
			name: "same release versions",
			a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			expected: 0,
		},
		{
			name: "different major versions",
			a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			b:    &Version{Major: 2, Minor: 2, Patch: 3, Type: TypeRelease},
			expected: -1,
		},
		{
			name: "different minor versions",
			a:    &Version{Major: 1, Minor: 1, Patch: 3, Type: TypeRelease},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			expected: -1,
		},
		{
			name: "different patch versions",
			a:    &Version{Major: 1, Minor: 2, Patch: 2, Type: TypeRelease},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			expected: -1,
		},
		{
			name: "release vs prerelease",
			a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeRelease},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypePrerelease, Prerelease: "~alpha"},
			expected: 1,
		},
		{
			name: "prerelease vs postrelease",
			a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypePrerelease, Prerelease: "~alpha"},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypePostrelease, Postrelease: ".fix"},
			expected: -2,
		},
		{
			name: "postrelease vs intermediate",
			a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypePostrelease, Postrelease: ".fix"},
			b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: TypeIntermediate, Intermediate: "_feature"},
			expected: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Compare(test.a, test.b)
			if result != test.expected {
				t.Errorf("Compare: expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestCompareIdentifiers(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"alpha", "beta", -1},
		{"beta", "alpha", 1},
		{"alpha", "alpha", 0},
		{"1", "2", -1},
		{"2", "1", 1},
		{"1", "1", 0},
		{"1", "alpha", -1},
		{"alpha", "1", 1},
		{"alpha.1", "alpha.2", -1},
		{"alpha.2", "alpha.1", 1},
		{"alpha.1", "alpha.1", 0},
		{"alpha.1", "beta.1", -1},
		{"beta.1", "alpha.1", 1},
	}

	for _, test := range tests {
		t.Run(test.a+"_vs_"+test.b, func(t *testing.T) {
			result := compareIdentifiers(test.a, test.b)
			if result != test.expected {
				t.Errorf("compareIdentifiers(%s, %s): expected %d, got %d", test.a, test.b, test.expected, result)
			}
		})
	}
}

func TestConvertGitTag(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3-alpha", "1.2.3~alpha"},
		{"v1.2.3-beta.1", "v1.2.3~beta.1"},
		{"1.2.3-rc.1", "1.2.3~rc.1"},
		{"1.2.3-pre.1", "1.2.3~pre.1"},
		{"1.2.3", "1.2.3"}, // No conversion needed
		{"1.2.3~alpha", "1.2.3~alpha"}, // Already converted
		{"1.2.3.fix", "1.2.3.fix"}, // Not a prerelease
		{"1.2.3_feature", "1.2.3_feature"}, // Not a prerelease
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ConvertGitTag(test.input)
			if result != test.expected {
				t.Errorf("ConvertGitTag(%s): expected %s, got %s", test.input, test.expected, result)
			}
		})
	}
}

func TestSort(t *testing.T) {
	versions := []string{
		"1.2.3",
		"1.2.3-alpha",
		"1.2.3-beta",
		"1.2.3-rc.1",
		"1.2.3.fix",
		"1.2.3_feature",
		"2.0.0",
		"1.2.4",
	}

	result, err := Sort(versions)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	expected := []string{
		"1.2.3~alpha",
		"1.2.3~beta",
		"1.2.3~rc.1",
		"1.2.3",
		"1.2.3.fix",
		"1.2.3_feature",
		"1.2.4",
		"2.0.0",
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d versions, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Position %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestSortReleaseVsPrerelease(t *testing.T) {
	// Test the specific scenario: v1.3.9 should be greater than v1.3.9-rc.9
	versions := []string{
		"v1.3.9",
		"v1.3.9-pre.1",
		"v1.3.9-pre.2",
		"v1.3.9-rc.3",
		"v1.3.9-rc.4",
		"v1.3.9-rc.5",
		"v1.3.9-rc.6",
		"v1.3.9-rc.7",
		"v1.3.9-rc.8",
		"v1.3.9-rc.9",
	}

	result, err := Sort(versions)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	// Expected order: all prereleases first, then the release version
	expected := []string{
		"v1.3.9~pre.1",
		"v1.3.9~pre.2",
		"v1.3.9~rc.3",
		"v1.3.9~rc.4",
		"v1.3.9~rc.5",
		"v1.3.9~rc.6",
		"v1.3.9~rc.7",
		"v1.3.9~rc.8",
		"v1.3.9~rc.9",
		"v1.3.9",
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d versions, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Position %d: expected %s, got %s", i, expected[i], v)
		}
	}

	// Verify that v1.3.9 is the greatest
	greatest := result[len(result)-1]
	if greatest != "v1.3.9" {
		t.Errorf("Expected v1.3.9 to be the greatest version, got %s", greatest)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"1.2.3", false},
		{"v1.2.3", false},
		{"1.2.3-alpha", false},
		{"1.2.3~beta.1", false},
		{"1.2.3.fix", false},
		{"1.2.3_feature", false},
		{"1.2", true},
		{"1.2.3.4", true},
		{"invalid", true},
		{"", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			err := Validate(test.input)
			if test.hasError && err == nil {
				t.Errorf("Expected error for input %s, but got none", test.input)
			}
			if !test.hasError && err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1.2.3", true},
		{"v1.2.3", true},
		{"1.2.3-alpha", true},
		{"1.2.3~beta.1", true},
		{"1.2.3.fix", true},
		{"1.2.3_feature", true},
		{"1.2", false},
		{"1.2.3.4", false},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsValid(test.input)
			if result != test.expected {
				t.Errorf("IsValid(%s): expected %v, got %v", test.input, test.expected, result)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
		hasError bool
	}{
		{"1.2.3", TypeRelease, false},
		{"1.2.3-alpha", TypePrerelease, false},
		{"1.2.3~beta.1", TypePrerelease, false},
		{"1.2.3.fix", TypePostrelease, false},
		{"1.2.3_feature", TypeIntermediate, false},
		{"1.2", TypeInvalid, true},
		{"invalid", TypeInvalid, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := GetType(test.input)
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", test.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
				return
			}
			if result != test.expected {
				t.Errorf("GetType(%s): expected %v, got %v", test.input, test.expected, result)
			}
		})
	}
}

func TestGetBuildType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"1.2.3", "Release", false},
		{"1.2.3-alpha", "Debug", false},
		{"1.2.3~beta.1", "Debug", false},
		{"1.2.3.fix", "Debug", false},
		{"1.2.3_feature", "Debug", false},
		{"1.2", "", true},
		{"invalid", "", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := GetBuildType(test.input)
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", test.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
				return
			}
			if result != test.expected {
				t.Errorf("GetBuildType(%s): expected %s, got %s", test.input, test.expected, result)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	v := &Version{
		Major:    1,
		Minor:    2,
		Patch:    3,
		Type:     TypeRelease,
		Original: "1.2.3",
	}
	
	expected := "1.2.3"
	if v.String() != expected {
		t.Errorf("Version.String(): expected %s, got %s", expected, v.String())
	}
}

func TestSortEmpty(t *testing.T) {
	result, err := Sort([]string{})
	if err != nil {
		t.Errorf("Sort with empty slice failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestSortSingleVersion(t *testing.T) {
	versions := []string{"1.2.3"}
	result, err := Sort(versions)
	if err != nil {
		t.Errorf("Sort with single version failed: %v", err)
	}
	if len(result) != 1 || result[0] != "1.2.3" {
		t.Errorf("Expected [1.2.3], got %v", result)
	}
}