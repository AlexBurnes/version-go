package main

import (
    "fmt"
    "log"
    
    "github.com/AlexBurnes/version-go/pkg/version"
)

func main() {
	// Example 1: Parse and validate a version
	fmt.Println("=== Version Parsing Example ===")
	
	versionStr := "1.2.3-alpha.1"
	v, err := version.Parse(versionStr)
	if err != nil {
		log.Fatalf("Failed to parse version: %v", err)
	}
	
	fmt.Printf("Version: %s\n", v.String())
	fmt.Printf("Major: %d, Minor: %d, Patch: %d\n", v.Major, v.Minor, v.Patch)
	fmt.Printf("Type: %s\n", v.Type.String())
	fmt.Printf("Build Type: %s\n", v.Type.BuildType())
	
	// Example 2: Check if a version is valid
	fmt.Println("\n=== Version Validation Example ===")
	
	validVersions := []string{"1.2.3", "v2.0.0-beta.1", "1.0.0.fix", "1.0.0_feature"}
	invalidVersions := []string{"1.2", "invalid", "1.2.3.4.5"}
	
	fmt.Println("Valid versions:")
	for _, v := range validVersions {
		if version.IsValid(v) {
			fmt.Printf("  ✓ %s\n", v)
		} else {
			fmt.Printf("  ✗ %s\n", v)
		}
	}
	
	fmt.Println("Invalid versions:")
	for _, v := range invalidVersions {
		if !version.IsValid(v) {
			fmt.Printf("  ✗ %s (correctly identified as invalid)\n", v)
		} else {
			fmt.Printf("  ✓ %s (incorrectly identified as valid)\n", v)
		}
	}
	
	// Example 3: Get version type and build type
	fmt.Println("\n=== Version Type Example ===")
	
	versions := []string{"1.2.3", "1.2.3-alpha", "1.2.3.fix", "1.2.3_feature"}
	
	for _, v := range versions {
		versionType, err := version.GetType(v)
		if err != nil {
			fmt.Printf("Error getting type for %s: %v\n", v, err)
			continue
		}
		
		buildType, err := version.GetBuildType(v)
		if err != nil {
			fmt.Printf("Error getting build type for %s: %v\n", v, err)
			continue
		}
		
		fmt.Printf("Version: %s | Type: %s | Build Type: %s\n", v, versionType.String(), buildType)
	}
	
	// Example 4: Sort versions
	fmt.Println("\n=== Version Sorting Example ===")
	
	versionsToSort := []string{
		"2.0.0",
		"1.2.3",
		"1.2.3-alpha",
		"1.2.3-beta.1",
		"1.2.3-rc.1",
		"1.2.3.fix",
		"1.2.3_feature",
		"1.2.4",
		"1.0.0",
	}
	
	fmt.Println("Original order:")
	for i, v := range versionsToSort {
		fmt.Printf("  %d. %s\n", i+1, v)
	}
	
	sortedVersions, err := version.Sort(versionsToSort)
	if err != nil {
		log.Fatalf("Failed to sort versions: %v", err)
	}
	
	fmt.Println("Sorted order:")
	for i, v := range sortedVersions {
		fmt.Printf("  %d. %s\n", i+1, v)
	}
	
	// Example 5: Compare versions
	fmt.Println("\n=== Version Comparison Example ===")
	
	v1, _ := version.Parse("1.2.3")
	v2, _ := version.Parse("1.2.4")
	v3, _ := version.Parse("1.2.3-alpha")
	
	comparisons := []struct {
		a, b *version.Version
		desc string
	}{
		{v1, v2, "1.2.3 vs 1.2.4"},
		{v2, v1, "1.2.4 vs 1.2.3"},
		{v1, v3, "1.2.3 vs 1.2.3-alpha"},
		{v3, v1, "1.2.3-alpha vs 1.2.3"},
		{v1, v1, "1.2.3 vs 1.2.3"},
	}
	
	for _, comp := range comparisons {
		result := version.Compare(comp.a, comp.b)
		var relation string
		switch {
		case result < 0:
			relation = "less than"
		case result > 0:
			relation = "greater than"
		default:
			relation = "equal to"
		}
		fmt.Printf("  %s is %s %s (result: %d)\n", comp.a.String(), relation, comp.b.String(), result)
	}
}