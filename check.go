package main

import "golang.org/x/mod/modfile"

// CheckVersions compares versions between target and reference
// and returns mismatches. If strict is true, it also reports
// dependencies that exist only in target.
func CheckVersions(targetMod, referenceMod *modfile.File, strict bool) *CheckResult {
	result := &CheckResult{}

	refVersions := BuildVersionMap(referenceMod)
	targetVersions := BuildVersionMap(targetMod)

	// Check for version mismatches and missing in reference
	for module, targetVersion := range targetVersions {
		if refVersion, exists := refVersions[module]; exists {
			// Module exists in both, check version
			if targetVersion != refVersion {
				result.DependencyMismatches = append(result.DependencyMismatches, VersionMismatch{
					Module:           module,
					TargetVersion:    targetVersion,
					ReferenceVersion: refVersion,
					OnlyInTarget:     false,
				})
			}
		} else if strict {
			// Module only exists in target, report if strict mode
			result.DependencyMismatches = append(result.DependencyMismatches, VersionMismatch{
				Module:           module,
				TargetVersion:    targetVersion,
				ReferenceVersion: "",
				OnlyInTarget:     true,
			})
		}
	}

	// Check Go version
	var targetGoVersion, refGoVersion string
	if targetMod.Go != nil {
		targetGoVersion = targetMod.Go.Version
	}
	if referenceMod.Go != nil {
		refGoVersion = referenceMod.Go.Version
	}

	if refGoVersion != "" && targetGoVersion != refGoVersion {
		result.GoVersionMismatch = &GoVersionMismatch{
			TargetVersion:    targetGoVersion,
			ReferenceVersion: refGoVersion,
		}
	}

	return result
}
