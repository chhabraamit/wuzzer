package main

import (
	"bufio"
	"chhabra.com/wuzzer/fuzzymatcher"
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
	fmt.Println("Enter the search query:")
	query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	query = strings.TrimSpace(query)

	fmt.Println("Enter the webpage content (type 'EOF' on a new line when finished):")
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOF" {
			break
		}
		lines = append(lines, line)
	}

	matches := fuzzymatcher.FuzzyMatch(query, lines)

	fmt.Printf("\nSearch results for query '%s':\n\n", query)
	for i, match := range matches {
		fmt.Printf("%d. (Score: %.2f) ", i+1, match.Score)
		printBoldMatches(match.String, match.MatchedWords)
	}
}
