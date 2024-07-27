package indexer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFileIndex(t *testing.T) {
	index := NewFileIndex()
	if index == nil {
		t.Error("NewFileIndex() returned nil")
	}
	if len(index.Lines) != 0 {
		t.Errorf("NewFileIndex() returned an index with %d lines, expected 0", len(index.Lines))
	}
}

func TestIndexDirectory(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "indexer_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string][]string{
		"file1.txt": {"Line 1", "Line 2", "Line 3"},
		"file2.go":  {"package main", "func main() {", "}"},
		"file3.md":  {"# Header", "Content", "More content"},
	}

	for filename, lines := range testFiles {
		path := filepath.Join(tempDir, filename)
		err := os.WriteFile(path, []byte(joinLines(lines)), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Create and populate the index
	index := NewFileIndex()
	err = index.IndexDirectory(tempDir, []string{"*.txt", "*.go"})
	if err != nil {
		t.Fatalf("IndexDirectory failed: %v", err)
	}

	// Check the number of indexed lines
	expectedLines := len(testFiles["file1.txt"]) + len(testFiles["file2.go"])
	if len(index.Lines) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(index.Lines))
	}

	// Check if all lines from file1.txt and file2.go are indexed
	for _, file := range []string{"file1.txt", "file2.go"} {
		for i, expectedLine := range testFiles[file] {
			found := false
			for _, indexedLine := range index.Lines {
				if indexedLine.FilePath == filepath.Join(tempDir, file) &&
					indexedLine.LineNumber == i+1 &&
					indexedLine.Content == expectedLine {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected line not found: %s (file: %s, line: %d)", expectedLine, file, i+1)
			}
		}
	}

	// Check that file3.md was not indexed
	for _, indexedLine := range index.Lines {
		if indexedLine.FilePath == filepath.Join(tempDir, "file3.md") {
			t.Errorf("file3.md should not have been indexed")
		}
	}
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
