package main

// VersionChange represents a single version update
type VersionChange struct {
	Module     string
	OldVersion string
	NewVersion string
}

// SyncResult contains the results of a sync operation
type SyncResult struct {
	DependencyChanges []VersionChange
	GoVersionChange   *GoVersionChange
}

// GoVersionChange represents a Go version update
type GoVersionChange struct {
	OldVersion string
	NewVersion string
}

// VersionMismatch represents a version difference in check mode
type VersionMismatch struct {
	Module           string
	TargetVersion    string
	ReferenceVersion string
	OnlyInTarget     bool // true if module exists only in target
}

// CheckResult contains the results of a check operation
type CheckResult struct {
	DependencyMismatches []VersionMismatch
	GoVersionMismatch    *GoVersionMismatch
}

// GoVersionMismatch represents a Go version difference
type GoVersionMismatch struct {
	TargetVersion    string
	ReferenceVersion string
}

// VersionMap is a map of module paths to their versions
type VersionMap map[string]string
