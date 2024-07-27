package indexer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type IndexedLine struct {
	FilePath   string
	LineNumber int
	Content    string
}

type FileIndex struct {
	Lines []IndexedLine
}

func NewFileIndex() *FileIndex {
	return &FileIndex{
		Lines: make([]IndexedLine, 0),
	}
}

func (fi *FileIndex) IndexDirectory(root string, patterns []string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDirectory(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if len(patterns) > 0 && !matchesAnyPattern(path, patterns) {
			return nil
		}
		return fi.indexFile(path)
	})
}

func (fi *FileIndex) indexFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // Increase buffer size to 1MB

	lineNumber := 1
	for scanner.Scan() {
		fi.Lines = append(fi.Lines, IndexedLine{
			FilePath:   path,
			LineNumber: lineNumber,
			Content:    scanner.Text(),
		})
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file %s: %v", path, err)
	}

	return nil
}

func matchesAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

func shouldSkipDirectory(path string) bool {
	// Skip hidden directories and some common directories we don't want to index
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}

	// Skip Go standard library and common build directories
	skippedDirs := []string{"go", "pkg", "node_modules", "vendor", "build", "dist", "env"}
	for _, dir := range skippedDirs {
		if base == dir {
			return true
		}
	}

	return false
}
