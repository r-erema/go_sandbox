package longestsubstringwithoutrepeatingcharacters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "the same sequences go next to each other",
			input: "abcabcbb",
			want:  3,
		},
		{
			name:  "repeated chars are next to each other in the middle and in the end",
			input: "pwwkew",
			want:  3,
		},
		{
			name:  "repeated chars in the beginning and in the middle",
			input: "dvdf",
			want:  3,
		},
		{
			name:  "empty string",
			input: " ",
			want:  1,
		},
		{
			name:  "palindrome",
			input: "abba",
			want:  2,
		},
		{
			name:  "some repeated chars are next to each other in the middle and some in the start and in the end",
			input: "uqinntq",
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, lengthOfLongestSubstring(tt.input))
		})
	}
}

// Time O(n), since we should iterate all the input
// Space O(n), we may allocate memory in map equal to input chars count.
func lengthOfLongestSubstring(input string) int {
	var res, left, right int

	encounteredChars := make(map[byte]int)

	for left, right = 0, 0; right < len(input); right++ {
		if prevCharIndex, stopExpandWindow := encounteredChars[input[right]]; stopExpandWindow &&
			prevCharIndex >= left {
			res = max(res, right-left)
			left = prevCharIndex + 1
		}

		encounteredChars[input[right]] = right
	}

	return max(res, right-left)
}
