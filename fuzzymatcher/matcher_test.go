package fuzzymatcher

import (
	"reflect"
	"testing"
)

func TestFuzzyMatchRanking(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		targets         []string
		expectedRanking []string
	}{
		{
			name:            "Basic Ranking",
			query:           "red apple",
			targets:         []string{"red delicious apple", "green apple pie", "apple cider", "red grapes", "big red balloon"},
			expectedRanking: []string{"red delicious apple", "green apple pie", "apple cider", "red grapes", "big red balloon"},
		},
		{
			name:            "Out of Order Ranking",
			query:           "apple red",
			targets:         []string{"red delicious apple", "red apple", "apple red sauce", "green apple"},
			expectedRanking: []string{"red apple", "red delicious apple", "apple red sauce", "green apple"},
		},
		{
			name:            "Partial Word Ranking",
			query:           "app red",
			targets:         []string{"red apple", "apply red paint", "red appetizer", "green application"},
			expectedRanking: []string{"red apple", "apply red paint", "red appetizer", "green application"},
		},
		{
			name:            "Case Insensitivity Ranking",
			query:           "RED APPLE",
			targets:         []string{"Red Delicious Apple", "GREEN APPLE PIE", "Apple cider", "red grape"},
			expectedRanking: []string{"Red Delicious Apple", "GREEN APPLE PIE", "Apple cider", "red grape"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FuzzyMatch(tt.query, tt.targets)

			// Check if the number of results matches the number of expected rankings
			if len(got) != len(tt.expectedRanking) {
				t.Errorf("FuzzyMatch() returned %d results, want %d", len(got), len(tt.expectedRanking))
				return
			}

			// Check if the ranking matches the expected ranking
			for i, expectedMatch := range tt.expectedRanking {
				if got[i].String != expectedMatch {
					t.Errorf("FuzzyMatch() ranking mismatch at position %d: got %s, want %s", i, got[i].String, expectedMatch)
				}
			}

			// Check if scores are in descending order
			for i := 1; i < len(got); i++ {
				if got[i-1].Score < got[i].Score {
					t.Errorf("Scores are not in descending order: %f > %f", got[i-1].Score, got[i].Score)
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		targets     []string
		expectedLen int
	}{
		{
			name:        "Empty Query",
			query:       "",
			targets:     []string{"apple", "banana", "cherry"},
			expectedLen: 0,
		},
		{
			name:        "No Matches",
			query:       "zebra",
			targets:     []string{"apple", "banana", "cherry"},
			expectedLen: 0,
		},
		{
			name:        "Single Character Query",
			query:       "a",
			targets:     []string{"apple", "banana", "cherry"},
			expectedLen: 2, // Expecting matches for "apple" and "banana"
		},
		{
			name:        "Special Characters",
			query:       "c++",
			targets:     []string{"c++ programming", "java coding", "python script"},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FuzzyMatch(tt.query, tt.targets)
			if len(got) != tt.expectedLen {
				t.Errorf("FuzzyMatch() returned %d results, want %d", len(got), tt.expectedLen)
			}
		})
	}
}

func TestConsistentRanking(t *testing.T) {
	query := "test query"
	targets := []string{"first match", "second match", "third match"}

	// Run the fuzzy match multiple times
	results1 := FuzzyMatch(query, targets)
	results2 := FuzzyMatch(query, targets)

	// Check if the rankings are consistent
	if !reflect.DeepEqual(results1, results2) {
		t.Errorf("Inconsistent rankings between multiple runs")
	}
}
