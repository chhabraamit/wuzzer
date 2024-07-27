package fuzzymatcher

import (
	"chhabra.com/wuzzer/indexer"
	"sort"
)

type IndexedMatch struct {
	IndexedLine  indexer.IndexedLine
	Score        float64
	MatchedWords []string
}

func FuzzyMatchIndexed(query string, index *indexer.FileIndex) []IndexedMatch {
	queryTokens := tokenize(query)
	var matches []IndexedMatch

	for _, line := range index.Lines {
		targetTokens := tokenize(line.Content)
		matchedIndices := make([]int, 0, len(queryTokens))
		matchedWords := make([]string, 0, len(queryTokens))

		for _, queryToken := range queryTokens {
			for j, targetToken := range targetTokens {
				if matchWord(queryToken, targetToken) {
					matchedIndices = append(matchedIndices, j)
					matchedWords = append(matchedWords, targetToken)
					break
				}
			}
		}

		if len(matchedIndices) > 0 {
			score := calculateScore(queryTokens, targetTokens, matchedIndices)
			matches = append(matches, IndexedMatch{
				IndexedLine:  line,
				Score:        score,
				MatchedWords: matchedWords,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches
}
