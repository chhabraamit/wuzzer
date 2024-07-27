package indexer

import (
	"bufio"
	"os"
	"path/filepath"
)

type IndexedLine struct {
	FilePath   string
	LineNumber int
	Content    string
}

qtype FileIndex struct {
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
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		fi.Lines = append(fi.Lines, IndexedLine{
			FilePath:   path,
			LineNumber: lineNumber,
			Content:    scanner.Text(),
		})
		lineNumber++
	}
	return scanner.Err()
}

func matchesAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}
