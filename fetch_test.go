package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"http URL", "http://example.com/go.mod", true},
		{"https URL", "https://example.com/go.mod", true},
		{"local file path", "/path/to/go.mod", false},
		{"relative path", "./go.mod", false},
		{"filename only", "go.mod", false},
		{"ftp URL", "ftp://example.com/go.mod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFetchFromURL(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		responseStatus int
		expectError    bool
	}{
		{
			name:           "successful fetch",
			responseBody:   "module example.com/test\n\ngo 1.21\n",
			responseStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "404 not found",
			responseBody:   "Not Found",
			responseStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "500 internal server error",
			responseBody:   "Internal Server Error",
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Fetch from the test server
			data, err := fetchFromURL(server.URL)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.responseBody, string(data))
			}
		})
	}
}

func TestFetchReference_LocalFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-go.mod")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := "module example.com/test\n\ngo 1.21\n"
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	// Fetch the local file
	data, err := FetchReference(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestFetchReference_URL(t *testing.T) {
	expectedContent := "module example.com/test\n\ngo 1.21\n"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	// Fetch from URL
	data, err := FetchReference(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, string(data))
}

func TestFetchReference_NonExistentFile(t *testing.T) {
	_, err := FetchReference("/non/existent/file.mod")
	assert.Error(t, err)
}

func TestGetReferenceDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"URL", "https://example.com/go.mod", "https://example.com/go.mod"},
		{"local path", "/path/to/go.mod", "/path/to/go.mod"},
		{"relative path", "./go.mod", "./go.mod"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetReferenceDisplayName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
