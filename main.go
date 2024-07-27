package main

import (
	"bufio"
	"chhabra.com/wuzzer/fuzzymatcher"
	"chhabra.com/wuzzer/indexer"
	"fmt"
	"os"
	"strings"
)

const bold = "\033[1m"
const reset = "\033[0m"

func printBoldMatches(line string, matchedWords []string) {
	lowercaseLine := strings.ToLower(line)
	lastIndex := 0

	for _, word := range matchedWords {
		index := strings.Index(lowercaseLine[lastIndex:], strings.ToLower(word))
		if index != -1 {
			index += lastIndex
			fmt.Print(line[lastIndex:index])
			fmt.Print(bold + line[index:index+len(word)] + reset)
			lastIndex = index + len(word)
		}
	}

	if lastIndex < len(line) {
		fmt.Print(line[lastIndex:])
	}
	fmt.Println()
}

func main() {
	rootDir := "/Users/chhabra/temple/wuzzer"

	patterns := []string{"*.go", "*.txt", "*.md"}

	index := indexer.NewFileIndex()
	err := index.IndexDirectory(rootDir, patterns)
	if err != nil {
		fmt.Printf("Error indexing directory: %v\n", err)
		return
	}

	fmt.Printf("Indexed %d lines\n", len(index.Lines))

	for {
		fmt.Println("\nEnter a search query (or 'quit' to exit):")
		query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		query = strings.TrimSpace(query)

		if query == "quit" {
			break
		}

		matches := fuzzymatcher.FuzzyMatchIndexed(query, index)

		fmt.Printf("\nSearch results for query '%s':\n\n", query)
		for i, match := range matches {
			if i >= 10 {
				break // Limit to top 10 results
			}
			fmt.Printf("%d. (Score: %.2f) %s:%d\n", i+1, match.Score, match.IndexedLine.FilePath, match.IndexedLine.LineNumber)
			printBoldMatches(match.IndexedLine.Content, match.MatchedWords)
		}
	}
}
