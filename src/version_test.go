package main

import (
    "testing"
)

func TestParseVersion(t *testing.T) {
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
                Type: VersionTypeRelease,
                Original: "1.2.3",
            },
            hasError: false,
        },
        {
            input: "v1.2.3",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypeRelease,
                Original: "v1.2.3",
            },
            hasError: false,
        },
        {
            input: "0.0.0",
            expected: &Version{
                Major: 0, Minor: 0, Patch: 0,
                Type: VersionTypeRelease,
                Original: "0.0.0",
            },
            hasError: false,
        },
        
        // Prerelease versions
        {
            input: "1.2.3-alpha",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypePrerelease,
                Prerelease: "~alpha",
                Original: "1.2.3-alpha",
            },
            hasError: false,
        },
        {
            input: "1.2.3~beta.1",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypePrerelease,
                Prerelease: "~beta.1",
                Original: "1.2.3~beta.1",
            },
            hasError: false,
        },
        {
            input: "v1.2.3-rc.1",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypePrerelease,
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
                Type: VersionTypePostrelease,
                Postrelease: ".fix",
                Original: "1.2.3.fix",
            },
            hasError: false,
        },
        {
            input: "1.2.3.post.1",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypePostrelease,
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
                Type: VersionTypeIntermediate,
                Intermediate: "_feature",
                Original: "1.2.3_feature",
            },
            hasError: false,
        },
        {
            input: "1.2.3_exp.1",
            expected: &Version{
                Major: 1, Minor: 2, Patch: 3,
                Type: VersionTypeIntermediate,
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
            result, err := ParseVersion(test.input)
            
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

func TestVersionTypeString(t *testing.T) {
    tests := []struct {
        vt       VersionType
        expected string
    }{
        {VersionTypeRelease, "release"},
        {VersionTypePrerelease, "prerelease"},
        {VersionTypePostrelease, "postrelease"},
        {VersionTypeIntermediate, "intermediate"},
        {VersionTypeInvalid, "invalid"},
    }

    for _, test := range tests {
        if test.vt.String() != test.expected {
            t.Errorf("VersionType.String() for %v: expected %s, got %s", test.vt, test.expected, test.vt.String())
        }
    }
}

func TestVersionTypeBuildType(t *testing.T) {
    tests := []struct {
        vt       VersionType
        expected string
    }{
        {VersionTypeRelease, "Release"},
        {VersionTypePrerelease, "Debug"},
        {VersionTypePostrelease, "Debug"},
        {VersionTypeIntermediate, "Debug"},
        {VersionTypeInvalid, "Debug"},
    }

    for _, test := range tests {
        if test.vt.BuildType() != test.expected {
            t.Errorf("VersionType.BuildType() for %v: expected %s, got %s", test.vt, test.expected, test.vt.BuildType())
        }
    }
}

func TestCompareVersions(t *testing.T) {
    tests := []struct {
        name     string
        a        *Version
        b        *Version
        expected int
    }{
        {
            name: "same release versions",
            a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            expected: 0,
        },
        {
            name: "different major versions",
            a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            b:    &Version{Major: 2, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            expected: -1,
        },
        {
            name: "different minor versions",
            a:    &Version{Major: 1, Minor: 1, Patch: 3, Type: VersionTypeRelease},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            expected: -1,
        },
        {
            name: "different patch versions",
            a:    &Version{Major: 1, Minor: 2, Patch: 2, Type: VersionTypeRelease},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            expected: -1,
        },
        {
            name: "release vs prerelease",
            a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeRelease},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypePrerelease, Prerelease: "-alpha"},
            expected: -1,
        },
        {
            name: "prerelease vs postrelease",
            a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypePrerelease, Prerelease: "-alpha"},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypePostrelease, Postrelease: ".fix"},
            expected: -1,
        },
        {
            name: "postrelease vs intermediate",
            a:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypePostrelease, Postrelease: ".fix"},
            b:    &Version{Major: 1, Minor: 2, Patch: 3, Type: VersionTypeIntermediate, Intermediate: "_feature"},
            expected: -1,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := CompareVersions(test.a, test.b)
            if result != test.expected {
                t.Errorf("CompareVersions: expected %d, got %d", test.expected, result)
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
            result := convertGitTag(test.input)
            if result != test.expected {
                t.Errorf("convertGitTag(%s): expected %s, got %s", test.input, test.expected, result)
            }
        })
    }
}

func TestSortVersions(t *testing.T) {
    // This test would require mocking stdin, so we'll test the parsing and comparison logic instead
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

    var parsedVersions []*Version
    for _, v := range versions {
        parsed, err := ParseVersion(v)
        if err != nil {
            t.Fatalf("Failed to parse version %s: %v", v, err)
        }
        parsedVersions = append(parsedVersions, parsed)
    }

    // Test that versions can be sorted
    // The expected order should be:
    // 1.2.3 (release)
    // 1.2.3-alpha (prerelease) -> converted to 1.2.3~alpha
    // 1.2.3-beta (prerelease) -> converted to 1.2.3~beta
    // 1.2.3-rc.1 (prerelease) -> converted to 1.2.3~rc.1
    // 1.2.3.fix (postrelease)
    // 1.2.3_feature (intermediate)
    // 1.2.4 (release)
    // 2.0.0 (release)

    // We can't easily test the full sorting without mocking stdin,
    // but we can test that the comparison function works correctly
    // The expected order should be: 1.2.3, 1.2.3-alpha, 1.2.3-beta, 1.2.3-rc.1, 1.2.3.fix, 1.2.3_feature, 1.2.4, 2.0.0
    expectedOrder := []string{"1.2.3", "1.2.3-alpha", "1.2.3-beta", "1.2.3-rc.1", "1.2.3.fix", "1.2.3_feature", "1.2.4", "2.0.0"}
    
    for i := 0; i < len(expectedOrder)-1; i++ {
        for j := i + 1; j < len(expectedOrder); j++ {
            // Find the versions in our parsed list
            var a, b *Version
            for _, v := range parsedVersions {
                if v.Original == expectedOrder[i] {
                    a = v
                }
                if v.Original == expectedOrder[j] {
                    b = v
                }
            }
            
            if a != nil && b != nil {
                result := CompareVersions(a, b)
                if result > 0 {
                    t.Errorf("Version %s should come before %s, but comparison returned %d", 
                        a.Original, b.Original, result)
                }
            }
        }
    }
}