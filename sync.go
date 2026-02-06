package main

import (
	"fmt"

	"golang.org/x/mod/modfile"
)

// CompareVersions compares target versions against reference versions
// and returns a list of changes that need to be made
func CompareVersions(targetMod *modfile.File, refVersions VersionMap) []VersionChange {
	var changes []VersionChange

	for _, req := range targetMod.Require {
		if refVersion, exists := refVersions[req.Mod.Path]; exists {
			if req.Mod.Version != refVersion {
				changes = append(changes, VersionChange{
					Module:     req.Mod.Path,
					OldVersion: req.Mod.Version,
					NewVersion: refVersion,
				})
			}
		}
	}

	return changes
}

// ApplyVersionChanges applies the version changes to the target modfile
func ApplyVersionChanges(targetMod *modfile.File, changes []VersionChange) error {
	for _, change := range changes {
		if err := targetMod.AddRequire(change.Module, change.NewVersion); err != nil {
			return fmt.Errorf("failed to update %s: %w", change.Module, err)
		}
	}
	return nil
}

// SyncVersions is the main business logic function that syncs versions
// from reference to target, including the Go version
func SyncVersions(targetMod, referenceMod *modfile.File) (*SyncResult, error) {
	result := &SyncResult{}

	// Sync dependency versions
	refVersions := BuildVersionMap(referenceMod)
	depChanges := CompareVersions(targetMod, refVersions)

	if len(depChanges) > 0 {
		if err := ApplyVersionChanges(targetMod, depChanges); err != nil {
			return nil, err
		}
	}
	result.DependencyChanges = depChanges

	// Sync Go version
	var targetGoVersion, refGoVersion string
	if targetMod.Go != nil {
		targetGoVersion = targetMod.Go.Version
	}
	if referenceMod.Go != nil {
		refGoVersion = referenceMod.Go.Version
	}

	if refGoVersion != "" && targetGoVersion != refGoVersion {
		result.GoVersionChange = &GoVersionChange{
			OldVersion: targetGoVersion,
			NewVersion: refGoVersion,
		}
		if err := targetMod.AddGoStmt(refGoVersion); err != nil {
			return nil, fmt.Errorf("failed to update Go version: %w", err)
		}
	}

	return result, nil
}
