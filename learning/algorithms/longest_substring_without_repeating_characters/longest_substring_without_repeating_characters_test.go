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
			input: "qrsvbspk",
			want:  5,
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
func lengthOfLongestSubstring(str string) int {
	var left, res int

	chars := make(map[byte]int)

	for right := range str {
		if lastPosition, ok := chars[str[right]]; ok && lastPosition >= left {
			left = lastPosition + 1
		}

		chars[str[right]] = right
		res = max(res, right-left+1)
	}

	return res
}
