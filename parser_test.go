package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

// Helper function to create a test modfile - shared across all test files
func createTestModFile(content string) (*modfile.File, error) {
	return modfile.Parse("test.mod", []byte(content), nil)
}

func TestParseGoMod(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid go.mod",
			content: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)`,
			expectError: false,
		},
		{
			name:        "invalid go.mod",
			content:     `this is not a valid go.mod file`,
			expectError: true,
		},
		{
			name: "minimal go.mod",
			content: `module example.com/test

go 1.21`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseGoMod("test.mod", []byte(tt.content))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildVersionMap(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected VersionMap
	}{
		{
			name: "simple dependencies",
			content: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
	golang.org/x/text v0.3.0
)`,
			expected: VersionMap{
				"github.com/pkg/errors": "v0.9.1",
				"golang.org/x/text":     "v0.3.0",
			},
		},
		{
			name: "with indirect dependencies",
			content: `module example.com/test

go 1.21

require (
	github.com/pkg/errors v0.9.1
)

require (
	golang.org/x/text v0.3.0 // indirect
)`,
			expected: VersionMap{
				"github.com/pkg/errors": "v0.9.1",
				"golang.org/x/text":     "v0.3.0",
			},
		},
		{
			name: "no dependencies",
			content: `module example.com/test

go 1.21`,
			expected: VersionMap{},
		},
		{
			name: "with replace dependencies",
			content: `module example.com/test

go 1.21

require package/something v0.0.0

replace package/something => ../../../package/src/go
`,
			expected: VersionMap{
				"package/something": "v0.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod, err := createTestModFile(tt.content)
			require.NoError(t, err, "Failed to parse modfile")

			result := BuildVersionMap(mod)

			assert.Equal(t, tt.expected, result)
		})
	}
}
