package longestrepeatingcharacterreplacement_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharacterReplacement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		k     int
		want  int
	}{
		{
			name:  "Equal count of 2 chars",
			input: "ABAB",
			k:     2,
			want:  4,
		},
		{
			name:  "Count of 1 char is bigger then other",
			input: "AABABBA",
			k:     1,
			want:  4,
		},
		{
			name:  "All chars the same",
			input: "AAAA",
			k:     0,
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, characterReplacement(tt.input, tt.k))
		})
	}
}

// Time O(n), since we iterate input one time
// Space O(1), since we use only additional fixed array with 26 elements.
func characterReplacement(str string, m int) int {
	var (
		res, left, right, maxFrequentChar int
		charsCount                        [26]int
	)

	for right = range str {
		charIndexRight := str[right] - 65
		charsCount[charIndexRight]++

		maxFrequentChar = max(maxFrequentChar, charsCount[charIndexRight])

		for (right - left + 1 - maxFrequentChar) > m {
			charIndexLeft := str[left] - 65
			charsCount[charIndexLeft]--
			left++
		}

		res = max(res, right-left+1)
	}

	return res
}
