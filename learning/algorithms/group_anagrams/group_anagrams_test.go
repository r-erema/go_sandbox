package groupanagrams_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupAnagrams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  [][]string
	}{
		{
			name:  "Case 0",
			input: []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			want:  [][]string{{"eat", "tea", "ate"}, {"tan", "nat"}, {"bat"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, groupAnagrams(tt.input))
		})
	}
}

// Time O(M * N) - we need to increment each symbol(M) in each word(N)
// Space O(N) - we need to persist each word(N) in hash map.
func groupAnagrams(strings []string) [][]string {
	anagramsMap := make(map[[26]byte][]string)

	for _, word := range strings {
		var wordASCIICount [26]byte

		for _, letterByte := range word {
			wordASCIICount[letterByte-'a']++
		}

		anagramsMap[wordASCIICount] = append(anagramsMap[wordASCIICount], word)
	}

	result := make([][]string, 0, len(anagramsMap))
	for _, group := range anagramsMap {
		result = append(result, group)
	}

	return result
}
