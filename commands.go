package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func syncCommand(args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	targetFile := fs.String("target", "", "Path to the target go.mod file to be modified")
	referenceFile := fs.String("reference", "", "Path or URL to the reference go.mod file with desired versions")
	dryRun := fs.Bool("dry-run", false, "Show changes without modifying the target file")
	verbose := fs.Bool("verbose", false, "Show detailed changes")

	// no need to check for errors cause the flagset is configured with flag.ExitOnError
	_ = fs.Parse(args)

	if *targetFile == "" || *referenceFile == "" {
		fmt.Println("Usage: gomodsync sync -target <target-go.mod> -reference <reference-go.mod|URL> [-dry-run] [-verbose]")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		os.Exit(1)
	}

	// Read and parse the target file
	targetData, err := os.ReadFile(*targetFile)
	if err != nil {
		log.Fatalf("Failed to read target file: %v", err)
	}

	// Get original file permissions to preserve them
	targetInfo, err := os.Stat(*targetFile)
	if err != nil {
		log.Fatalf("Failed to stat target file: %v", err)
	}
	targetPerms := targetInfo.Mode().Perm()

	targetMod, err := ParseGoMod(*targetFile, targetData)
	if err != nil {
		log.Fatalf("Failed to parse target file: %v", err)
	}

	// Fetch and parse the reference (from URL or local path)
	referenceData, err := FetchReference(*referenceFile)
	if err != nil {
		log.Fatalf("Failed to fetch reference: %v", err)
	}

	referenceName := GetReferenceDisplayName(*referenceFile)
	referenceMod, err := ParseGoMod(referenceName, referenceData)
	if err != nil {
		log.Fatalf("Failed to parse reference: %v", err)
	}

	// Sync versions
	result, err := SyncVersions(targetMod, referenceMod)
	if err != nil {
		log.Fatalf("Failed to sync versions: %v", err)
	}

	totalChanges := len(result.DependencyChanges)
	if result.GoVersionChange != nil {
		totalChanges++
	}

	if totalChanges == 0 {
		fmt.Println("✓ No version differences found. Target file is already in sync.")
		return
	}

	// Print changes if verbose
	if *verbose {
		fmt.Printf("Changes to be made:\n\n")

		if result.GoVersionChange != nil {
			fmt.Printf("  go: %s -> %s\n", result.GoVersionChange.OldVersion, result.GoVersionChange.NewVersion)
		}

		for _, change := range result.DependencyChanges {
			fmt.Printf("  %s: %s -> %s\n", change.Module, change.OldVersion, change.NewVersion)
		}
		fmt.Println()
	}

	if *dryRun {
		fmt.Printf("Dry-run mode: %d change(s) identified but not applied.\n", totalChanges)
		if *verbose {
			fmt.Println("\nPreview of updated go.mod:")
			fmt.Println("---")
			formatted, err := targetMod.Format()
			if err != nil {
				log.Fatalf("Failed to format target file: %v", err)
			}
			fmt.Println(string(formatted))
		}
		return
	}

	// Format and write the updated target file
	formatted, err := targetMod.Format()
	if err != nil {
		log.Fatalf("Failed to format target file: %v", err)
	}

	// Write with original file permissions
	if err := os.WriteFile(*targetFile, formatted, targetPerms); err != nil {
		log.Fatalf("Failed to write target file: %v", err)
	}

	fmt.Printf("✓ Successfully updated %s (%d change(s) applied)\n", *targetFile, totalChanges)
}

func checkCommand(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	targetFile := fs.String("target", "", "Path to the target go.mod file to check")
	referenceFile := fs.String("reference", "", "Path or URL to the reference go.mod file with desired versions")
	strict := fs.Bool("strict", false, "Fail if target has dependencies not in reference")
	verbose := fs.Bool("verbose", false, "Show detailed version mismatches")

	// no need to check for errors cause the flagset is configured with flag.ExitOnError
	_ = fs.Parse(args)

	if *targetFile == "" || *referenceFile == "" {
		fmt.Println("Usage: gomodsync check -target <target-go.mod> -reference <reference-go.mod|URL> [-strict] [-verbose]")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		os.Exit(1)
	}

	// Read and parse the target file
	targetData, err := os.ReadFile(*targetFile)
	if err != nil {
		log.Fatalf("Failed to read target file: %v", err)
	}

	targetMod, err := ParseGoMod(*targetFile, targetData)
	if err != nil {
		log.Fatalf("Failed to parse target file: %v", err)
	}

	// Fetch and parse the reference (from URL or local path)
	referenceData, err := FetchReference(*referenceFile)
	if err != nil {
		log.Fatalf("Failed to fetch reference: %v", err)
	}

	referenceName := GetReferenceDisplayName(*referenceFile)
	referenceMod, err := ParseGoMod(referenceName, referenceData)
	if err != nil {
		log.Fatalf("Failed to parse reference: %v", err)
	}

	// Check versions
	result := CheckVersions(targetMod, referenceMod, *strict)

	totalMismatches := len(result.DependencyMismatches)
	if result.GoVersionMismatch != nil {
		totalMismatches++
	}

	if totalMismatches == 0 {
		fmt.Println("✓ All versions match (dependencies and Go version)!")
		os.Exit(0)
	}

	// Print summary or detailed mismatches based on verbose flag
	if *verbose {
		fmt.Printf("✗ Found %d version mismatch(es):\n\n", totalMismatches)

		if result.GoVersionMismatch != nil {
			fmt.Printf("  go: %s != %s\n", result.GoVersionMismatch.TargetVersion, result.GoVersionMismatch.ReferenceVersion)
		}

		for _, mismatch := range result.DependencyMismatches {
			if mismatch.OnlyInTarget {
				fmt.Printf("  %s: %s (not in reference)\n", mismatch.Module, mismatch.TargetVersion)
			} else {
				fmt.Printf("  %s: %s != %s\n", mismatch.Module, mismatch.TargetVersion, mismatch.ReferenceVersion)
			}
		}
	} else {
		fmt.Printf("✗ Version check failed: %d mismatch(es) found\n", totalMismatches)
		fmt.Println("Run with -verbose to see details")
	}

	os.Exit(1)
}
