package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// isURL checks if the given string is a URL
func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// FetchReference fetches the content from either a URL or local file path.
// If the reference is a URL (starts with http:// or https://), it downloads the content.
// Otherwise, it reads the content from the local file system.
func FetchReference(reference string) ([]byte, error) {
	if isURL(reference) {
		return fetchFromURL(reference)
	}
	return os.ReadFile(reference)
}

// fetchFromURL downloads content from a URL
func fetchFromURL(url string) ([]byte, error) {
	// #nosec G107 -- URL is user-provided via CLI flag, this is the intended functionality
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// GetReferenceDisplayName returns a display name for the reference.
// For URLs, it returns the URL itself. For file paths, it returns the path.
func GetReferenceDisplayName(reference string) string {
	if isURL(reference) {
		return reference
	}
	return reference
}
