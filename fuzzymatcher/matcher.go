package fuzzymatcher

import (
	"sort"
	"strings"
	"unicode"
)

// Match represents a single match result
type Match struct {
	String       string   `json:"match"`
	Score        float64  `json:"score"`
	MatchedWords []string `json:"matched_words"`
}

// tokenize splits a string into lowercase words
func tokenize(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// isPrefix checks if a is a prefix of b
func isPrefix(a, b string) bool {
	return len(a) <= len(b) && b[:len(a)] == a
}

// matchWord checks if query word matches target word
func matchWord(query, target string) bool {
	if query == target {
		return true
	}
	return isPrefix(query, target)
}

func calculateScore(queryTokens, targetTokens []string, matchedIndices []int) float64 {
	// Base score: proportion of query words matched
	baseScore := float64(len(matchedIndices)) / float64(len(queryTokens))

	// Increase the importance of matching more words
	wordMatchBonus := baseScore * 0.5

	// Small bonus for order preservation (less important now)
	orderBonus := 0.0
	for i := 1; i < len(matchedIndices); i++ {
		if matchedIndices[i] > matchedIndices[i-1] {
			orderBonus += 0.01
		}
	}

	// Penalty for extra words (slightly reduced)
	extraWordsPenalty := float64(len(targetTokens)-len(queryTokens)) * 0.05

	// Combine scores
	totalScore := baseScore + wordMatchBonus + orderBonus - extraWordsPenalty

	return max(0.0, min(1.0, totalScore))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(i float64, f float64) float64 {
	if i > f {
		return i
	}
	return f
}

// FuzzyMatch performs fuzzy matching of query against targets
func FuzzyMatch(query string, targets []string) []Match {
	queryTokens := tokenize(query)
	var matches []Match

	for _, target := range targets {
		targetTokens := tokenize(target)
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
			matches = append(matches, Match{
				String:       target,
				Score:        score,
				MatchedWords: matchedWords,
			})
		}
	}

	// Sort matches by score (descending)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches
}
