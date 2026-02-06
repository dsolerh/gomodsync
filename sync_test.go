package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name            string
		targetContent   string
		referenceMap    VersionMap
		expectedChanges []VersionChange
	}{
		{
			name: "versions need updating",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			referenceMap: VersionMap{
				"github.com/pkg/errors": "v0.9.2",
				"golang.org/x/text":     "v0.4.0",
			},
			expectedChanges: []VersionChange{
				{
					Module:     "github.com/pkg/errors",
					OldVersion: "v0.9.1",
					NewVersion: "v0.9.2",
				},
				{
					Module:     "golang.org/x/text",
					OldVersion: "v0.3.0",
					NewVersion: "v0.4.0",
				},
			},
		},
		{
			name: "no changes needed",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			referenceMap: VersionMap{
				"github.com/pkg/errors": "v0.9.1",
				"golang.org/x/text":     "v0.3.0",
			},
			expectedChanges: []VersionChange{},
		},
		{
			name: "partial overlap",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
	github.com/unique/dep v1.0.0
)`,
			referenceMap: VersionMap{
				"github.com/pkg/errors": "v0.9.2",
				"golang.org/x/text":     "v0.3.0",
			},
			expectedChanges: []VersionChange{
				{
					Module:     "github.com/pkg/errors",
					OldVersion: "v0.9.1",
					NewVersion: "v0.9.2",
				},
			},
		},
		{
			name: "reference has additional dependencies",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			referenceMap: VersionMap{
				"github.com/pkg/errors": "v0.9.1",
				"golang.org/x/text":     "v0.3.0",
			},
			expectedChanges: []VersionChange{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetMod, err := createTestModFile(tt.targetContent)
			require.NoError(t, err, "Failed to parse target modfile")

			changes := CompareVersions(targetMod, tt.referenceMap)

			assert.Equal(t, len(tt.expectedChanges), len(changes), "Number of changes mismatch")

			for i, expectedChange := range tt.expectedChanges {
				if assert.Less(t, i, len(changes), "Missing expected change: %+v", expectedChange) {
					actual := changes[i]
					assert.Equal(t, expectedChange.Module, actual.Module, "Change %d: module mismatch", i)
					assert.Equal(t, expectedChange.OldVersion, actual.OldVersion, "Change %d: old version mismatch", i)
					assert.Equal(t, expectedChange.NewVersion, actual.NewVersion, "Change %d: new version mismatch", i)
				}
			}
		})
	}
}

func TestApplyVersionChanges(t *testing.T) {
	tests := []struct {
		name             string
		targetContent    string
		changes          []VersionChange
		expectedVersions VersionMap
		expectError      bool
	}{
		{
			name: "apply single change",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			changes: []VersionChange{
				{
					Module:     "github.com/pkg/errors",
					OldVersion: "v0.9.1",
					NewVersion: "v0.9.2",
				},
			},
			expectedVersions: VersionMap{
				"github.com/pkg/errors": "v0.9.2",
			},
			expectError: false,
		},
		{
			name: "apply multiple changes",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			changes: []VersionChange{
				{
					Module:     "github.com/pkg/errors",
					OldVersion: "v0.9.1",
					NewVersion: "v0.9.2",
				},
				{
					Module:     "golang.org/x/text",
					OldVersion: "v0.3.0",
					NewVersion: "v0.4.0",
				},
			},
			expectedVersions: VersionMap{
				"github.com/pkg/errors": "v0.9.2",
				"golang.org/x/text":     "v0.4.0",
			},
			expectError: false,
		},
		{
			name: "no changes",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			changes: []VersionChange{},
			expectedVersions: VersionMap{
				"github.com/pkg/errors": "v0.9.1",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetMod, err := createTestModFile(tt.targetContent)
			require.NoError(t, err, "Failed to parse target modfile")

			err = ApplyVersionChanges(targetMod, tt.changes)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")

				resultVersions := BuildVersionMap(targetMod)
				for module, expectedVersion := range tt.expectedVersions {
					assert.Contains(t, resultVersions, module, "Module %s not found after applying changes", module)
					assert.Equal(t, expectedVersion, resultVersions[module], "Module %s: version mismatch", module)
				}
			}
		})
	}
}

func TestSyncVersions(t *testing.T) {
	tests := []struct {
		name                  string
		targetContent         string
		referenceContent      string
		expectedChanges       int
		expectGoVersionChange bool
		expectedOldGoVersion  string
		expectedNewGoVersion  string
		expectError           bool
	}{
		{
			name: "successful sync with changes",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			referenceContent: `module example.com/reference

go 1.21

require (
	github.com/pkg/errors v0.9.2
	golang.org/x/text v0.4.0
)`,
			expectedChanges:       2,
			expectGoVersionChange: false,
			expectError:           false,
		},
		{
			name: "no changes needed",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			referenceContent: `module example.com/reference

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			expectedChanges:       0,
			expectGoVersionChange: false,
			expectError:           false,
		},
		{
			name: "partial overlap sync",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
	github.com/unique/target v1.0.0
)`,
			referenceContent: `module example.com/reference

go 1.21

require (
	github.com/pkg/errors v0.9.2
	golang.org/x/text v0.4.0
	github.com/unique/reference v1.2.0
)`,
			expectedChanges:       2,
			expectGoVersionChange: false,
			expectError:           false,
		},
		{
			name: "target has more dependencies",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
	github.com/extra/dep v1.0.0
)`,
			referenceContent: `module example.com/reference

go 1.21

require (
	github.com/pkg/errors v0.9.2
)`,
			expectedChanges:       1,
			expectGoVersionChange: false,
			expectError:           false,
		},
		{
			name: "sync go version only",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			referenceContent: `module example.com/reference

go 1.22

require (
	github.com/pkg/errors v0.9.1
)`,
			expectedChanges:       0,
			expectGoVersionChange: true,
			expectedOldGoVersion:  "1.21",
			expectedNewGoVersion:  "1.22",
			expectError:           false,
		},
		{
			name: "sync go version and dependencies",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			referenceContent: `module example.com/reference

go 1.22

require (
	github.com/pkg/errors v0.9.2
	golang.org/x/text v0.4.0
)`,
			expectedChanges:       2,
			expectGoVersionChange: true,
			expectedOldGoVersion:  "1.21",
			expectedNewGoVersion:  "1.22",
			expectError:           false,
		},
		{
			name: "target missing go version",
			targetContent: `module example.com/test

require (
	github.com/pkg/errors v0.9.1
)`,
			referenceContent: `module example.com/reference

go 1.22

require (
	github.com/pkg/errors v0.9.1
)`,
			expectedChanges:       0,
			expectGoVersionChange: true,
			expectedOldGoVersion:  "",
			expectedNewGoVersion:  "1.22",
			expectError:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetMod, err := createTestModFile(tt.targetContent)
			require.NoError(t, err, "Failed to parse target modfile")

			referenceMod, err := createTestModFile(tt.referenceContent)
			require.NoError(t, err, "Failed to parse reference modfile")

			result, err := SyncVersions(targetMod, referenceMod)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}

			assert.Equal(t, tt.expectedChanges, len(result.DependencyChanges), "Dependency changes count mismatch")

			if tt.expectGoVersionChange {
				assert.NotNil(t, result.GoVersionChange, "Expected Go version change but got none")
				if result.GoVersionChange != nil {
					assert.Equal(t, tt.expectedOldGoVersion, result.GoVersionChange.OldVersion, "Old Go version mismatch")
					assert.Equal(t, tt.expectedNewGoVersion, result.GoVersionChange.NewVersion, "New Go version mismatch")
				}
			} else {
				assert.Nil(t, result.GoVersionChange, "Expected no Go version change but got: %v", result.GoVersionChange)
			}
		})
	}
}
