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
	Lines             []IndexedLine
	ScannedDirs       []string
	SkippedDirs       []string
	TotalFilesScanned int
}

func NewFileIndex() *FileIndex {
	return &FileIndex{
		Lines:       make([]IndexedLine, 0),
		ScannedDirs: make([]string, 0),
		SkippedDirs: make([]string, 0),
	}
}

func (fi *FileIndex) IndexDirectory(root string, patterns []string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDirectory(path) {
				fi.SkippedDirs = append(fi.SkippedDirs, path)
				return filepath.SkipDir
			}
			fi.ScannedDirs = append(fi.ScannedDirs, path)
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

	fi.TotalFilesScanned++
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
	skippedDirs := []string{"go", "pkg", "node_modules", "vendor", "build", "dist"}
	for _, dir := range skippedDirs {
		if base == dir {
			return true
		}
	}

	return false
}

func (fi *FileIndex) PrintStats() {
	fmt.Printf("Total files scanned: %d\n", fi.TotalFilesScanned)
	fmt.Printf("Total lines indexed: %d\n", len(fi.Lines))
	fmt.Printf("Scanned directories: %d\n", len(fi.ScannedDirs))
	fmt.Printf("Skipped directories: %d\n", len(fi.SkippedDirs))

	fmt.Println("\nTop 10 scanned directories:")
	for i, dir := range fi.ScannedDirs {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s\n", dir)
	}

	fmt.Println("\nSkipped directories:")
	for _, dir := range fi.SkippedDirs {
		fmt.Printf("  %s\n", dir)
	}
}
