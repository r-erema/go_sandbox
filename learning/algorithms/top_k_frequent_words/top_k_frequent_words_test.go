package topkfrequentwords_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopKFrequent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		words []string
		k     int
		want  []string
	}{
		{
			name:  "Two words",
			words: []string{"i", "love", "leetcode", "i", "love", "coding"},
			k:     2,
			want:  []string{"i", "love"},
		},
		{
			name:  "Three words",
			words: []string{"i", "love", "leetcode", "i", "love", "coding"},
			k:     3,
			want:  []string{"i", "love", "coding"},
		},
		{
			name:  "Four words",
			words: []string{"the", "day", "is", "sunny", "the", "the", "the", "sunny", "is", "is"},
			k:     4,
			want:  []string{"the", "is", "sunny", "day"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, topKFrequent(tt.words, tt.k))
		})
	}
}

// Time O(N*logN)
// O(NlogN) where N is the length of words.
// We count the frequency of each word in O(N) time,
// then we sort the given words in O(NlogN) time.
//
// Memory O(n)
// We have a hash table that contains a number of rows equal to the input.
func topKFrequent(words []string, k int) []string { //nolint: varnamelen
	wordsCount := make(map[string]int)

	var keys []string

	for _, word := range words {
		if _, ok := wordsCount[word]; !ok {
			keys = append(keys, word)
		}
		wordsCount[word]++
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if wordsCount[keys[i]] == wordsCount[keys[j]] {
			return keys[i] < keys[j]
		}

		return wordsCount[keys[i]] > wordsCount[keys[j]]
	})

	return keys[:k]
}
