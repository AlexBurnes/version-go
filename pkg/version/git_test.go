package version

import (
    "testing"
)

func TestCheckGitAvailable(t *testing.T) {
    // This test will only pass if git is installed
    err := checkGitAvailable()
    if err != nil {
        // Skip test if git is not available
        t.Skipf("Git is not available: %v", err)
    }
}

func TestCheckGitRepo(t *testing.T) {
    // This test assumes we're running in a git repository
    err := checkGitRepo()
    if err != nil {
        // Skip test if not in a git repository
        t.Skipf("Not in a git repository: %v", err)
    }
}

func TestGetVersion(t *testing.T) {
    // This test assumes we're running in a git repository with version tags
    version, err := GetVersion()
    if err != nil {
        if IsGitNotFound(err) {
            t.Skipf("Git is not available: %v", err)
            return
        }
        if IsNotGitRepo(err) {
            t.Skipf("Not in a git repository: %v", err)
            return
        }
        if IsNoGitTags(err) {
            t.Skipf("No version tags found: %v", err)
            return
        }
        t.Fatalf("GetVersion() failed: %v", err)
    }

    // Verify that the version is not empty
    if version == "" {
        t.Error("GetVersion() returned empty version")
    }

    // Verify that the version can be parsed
    _, err = Parse(version)
    if err != nil {
        t.Errorf("GetVersion() returned invalid version format: %s, error: %v", version, err)
    }

    t.Logf("Current git version: %s", version)
}

func TestGetVersionWithPrefix(t *testing.T) {
    // This test assumes we're running in a git repository with version tags
    version, err := GetVersionWithPrefix()
    if err != nil {
        if IsGitNotFound(err) {
            t.Skipf("Git is not available: %v", err)
            return
        }
        if IsNotGitRepo(err) {
            t.Skipf("Not in a git repository: %v", err)
            return
        }
        if IsNoGitTags(err) {
            t.Skipf("No version tags found: %v", err)
            return
        }
        t.Fatalf("GetVersionWithPrefix() failed: %v", err)
    }

    // Verify that the version is not empty
    if version == "" {
        t.Error("GetVersionWithPrefix() returned empty version")
    }

    // Verify that the version starts with 'v'
    if len(version) == 0 || version[0] != 'v' {
        t.Errorf("GetVersionWithPrefix() should return version with 'v' prefix, got: %s", version)
    }

    t.Logf("Current git version with prefix: %s", version)
}

func TestGetRawTag(t *testing.T) {
    // This test assumes we're running in a git repository with version tags
    tag, err := GetRawTag()
    if err != nil {
        if IsGitNotFound(err) {
            t.Skipf("Git is not available: %v", err)
            return
        }
        if IsNotGitRepo(err) {
            t.Skipf("Not in a git repository: %v", err)
            return
        }
        if IsNoGitTags(err) {
            t.Skipf("No version tags found: %v", err)
            return
        }
        t.Fatalf("GetRawTag() failed: %v", err)
    }

    // Verify that the tag is not empty
    if tag == "" {
        t.Error("GetRawTag() returned empty tag")
    }

    // Verify that the tag starts with 'v'
    if len(tag) == 0 || tag[0] != 'v' {
        t.Errorf("GetRawTag() should return tag with 'v' prefix, got: %s", tag)
    }

    t.Logf("Current git raw tag: %s", tag)
}

func TestGitErrorTypes(t *testing.T) {
    tests := []struct {
        name     string
        err      error
        isNotFound bool
        isNotRepo  bool
        isNoTags   bool
    }{
        {
            name: "git not found error",
            err: &GitError{
                Type:    "not_found",
                Message: "git not found",
            },
            isNotFound: true,
            isNotRepo:  false,
            isNoTags:   false,
        },
        {
            name: "not git repo error",
            err: &GitError{
                Type:    "not_repo",
                Message: "not a git repository",
            },
            isNotFound: false,
            isNotRepo:  true,
            isNoTags:   false,
        },
        {
            name: "no git tags error",
            err: &GitError{
                Type:    "no_tags",
                Message: "no tags found",
            },
            isNotFound: false,
            isNotRepo:  false,
            isNoTags:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := IsGitNotFound(tt.err); got != tt.isNotFound {
                t.Errorf("IsGitNotFound() = %v, want %v", got, tt.isNotFound)
            }
            if got := IsNotGitRepo(tt.err); got != tt.isNotRepo {
                t.Errorf("IsNotGitRepo() = %v, want %v", got, tt.isNotRepo)
            }
            if got := IsNoGitTags(tt.err); got != tt.isNoTags {
                t.Errorf("IsNoGitTags() = %v, want %v", got, tt.isNoTags)
            }
        })
    }
}

func TestGitErrorMessage(t *testing.T) {
    err := &GitError{
        Type:    "not_found",
        Message: "test error message",
    }

    if err.Error() != "test error message" {
        t.Errorf("GitError.Error() = %v, want %v", err.Error(), "test error message")
    }
}

// TestGetVersionIntegration is an integration test that requires:
// - git to be installed
// - running in a git repository
// - at least one version tag (v[0-9]*) to exist
func TestGetVersionIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Try to get version
    version, err := GetVersion()
    if err != nil {
        if IsGitNotFound(err) {
            t.Skip("Git is not installed, skipping integration test")
            return
        }
        if IsNotGitRepo(err) {
            t.Skip("Not in a git repository, skipping integration test")
            return
        }
        if IsNoGitTags(err) {
            t.Skip("No version tags found, skipping integration test")
            return
        }
        t.Fatalf("Unexpected error: %v", err)
    }

    // Verify version format
    if version == "" {
        t.Fatal("GetVersion returned empty string")
    }

    // Parse and validate the version
    v, err := Parse(version)
    if err != nil {
        t.Fatalf("GetVersion returned invalid version %q: %v", version, err)
    }

    t.Logf("Successfully retrieved version from git: %s (type: %s)", version, v.Type)
}

