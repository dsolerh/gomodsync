package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckVersions(t *testing.T) {
	tests := []struct {
		name                    string
		targetContent           string
		referenceContent        string
		strict                  bool
		expectedMismatches      int
		expectGoVersionMismatch bool
		expectedTargetGoVersion string
		expectedRefGoVersion    string
		checkMismatches         func(t *testing.T, mismatches []VersionMismatch)
	}{
		{
			name: "no mismatches",
			targetContent: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			referenceContent: `module example.com/reference

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			strict:                  false,
			expectedMismatches:      0,
			expectGoVersionMismatch: false,
			checkMismatches:         nil,
		},
		{
			name: "version mismatches",
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
			strict:                  false,
			expectedMismatches:      2,
			expectGoVersionMismatch: false,
			checkMismatches: func(t *testing.T, mismatches []VersionMismatch) {
				t.Helper()
				for _, m := range mismatches {
					assert.False(t, m.OnlyInTarget, "Expected no OnlyInTarget mismatches in non-strict mode")
				}
			},
		},
		{
			name: "target has extra deps - non-strict mode",
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
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			strict:                  false,
			expectedMismatches:      0,
			expectGoVersionMismatch: false,
			checkMismatches:         nil,
		},
		{
			name: "target has extra deps - strict mode",
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
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			strict:                  true,
			expectedMismatches:      1,
			expectGoVersionMismatch: false,
			checkMismatches: func(t *testing.T, mismatches []VersionMismatch) {
				t.Helper()
				found := false
				for _, m := range mismatches {
					if m.Module == "github.com/extra/dep" && m.OnlyInTarget {
						found = true
						assert.Equal(t, "v1.0.0", m.TargetVersion)
						assert.Empty(t, m.ReferenceVersion)
					}
				}
				assert.True(t, found, "Expected to find github.com/extra/dep in mismatches")
			},
		},
		{
			name: "mixed mismatches - strict mode",
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
	golang.org/x/text v0.3.0
)`,
			strict:                  true,
			expectedMismatches:      2,
			expectGoVersionMismatch: false,
			checkMismatches: func(t *testing.T, mismatches []VersionMismatch) {
				t.Helper()
				foundVersionMismatch := false
				foundOnlyInTarget := false
				for _, m := range mismatches {
					if m.Module == "github.com/pkg/errors" && !m.OnlyInTarget {
						foundVersionMismatch = true
					}
					if m.Module == "github.com/extra/dep" && m.OnlyInTarget {
						foundOnlyInTarget = true
					}
				}
				assert.True(t, foundVersionMismatch, "Expected to find version mismatch for github.com/pkg/errors")
				assert.True(t, foundOnlyInTarget, "Expected to find OnlyInTarget for github.com/extra/dep")
			},
		},
		{
			name: "go version mismatch only",
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
			strict:                  false,
			expectedMismatches:      0,
			expectGoVersionMismatch: true,
			expectedTargetGoVersion: "1.21",
			expectedRefGoVersion:    "1.22",
			checkMismatches:         nil,
		},
		{
			name: "go version and dependency mismatches",
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
			strict:                  false,
			expectedMismatches:      2,
			expectGoVersionMismatch: true,
			expectedTargetGoVersion: "1.21",
			expectedRefGoVersion:    "1.22",
			checkMismatches:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetMod, err := createTestModFile(tt.targetContent)
			require.NoError(t, err, "Failed to parse target modfile")

			referenceMod, err := createTestModFile(tt.referenceContent)
			require.NoError(t, err, "Failed to parse reference modfile")

			result := CheckVersions(targetMod, referenceMod, tt.strict)

			assert.Equal(t, tt.expectedMismatches, len(result.DependencyMismatches), "Unexpected number of mismatches")

			if tt.checkMismatches != nil {
				tt.checkMismatches(t, result.DependencyMismatches)
			}

			if tt.expectGoVersionMismatch {
				assert.NotNil(t, result.GoVersionMismatch, "Expected Go version mismatch but got none")
				if result.GoVersionMismatch != nil {
					assert.Equal(t, tt.expectedTargetGoVersion, result.GoVersionMismatch.TargetVersion, "Target Go version mismatch")
					assert.Equal(t, tt.expectedRefGoVersion, result.GoVersionMismatch.ReferenceVersion, "Reference Go version mismatch")
				}
			} else {
				assert.Nil(t, result.GoVersionMismatch, "Expected no Go version mismatch but got: %v", result.GoVersionMismatch)
			}
		})
	}
}
