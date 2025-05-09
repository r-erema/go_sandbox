package lettercombinationsofaphonenumber_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetterCombinations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		digits string
		want   []string
	}{
		{
			name:   "2 digits",
			digits: "23",
			want:   []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, letterCombinations(tt.digits))
		})
	}
}

// Time O(n * 4^n), where n is the length of the input string,
// since there are O(4^n) combinations, and each takes O(n) time to build

// Space O(4^n), because we store all the possible combinations of letters.
func letterCombinations(digits string) []string {
	if digits == "" {
		return nil
	}

	digitsToChar := map[byte][]string{
		'2': {"a", "b", "c"},
		'3': {"d", "e", "f"},
		'4': {"g", "h", "i"},
		'5': {"j", "k", "l"},
		'6': {"m", "n", "o"},
		'7': {"p", "q", "r", "s"},
		'8': {"t", "u", "v"},
		'9': {"w", "x", "y", "z"},
	}

	var (
		backtrack func(i int, curStr string)
		res       []string
	)

	backtrack = func(i int, curStr string) {
		if len(curStr) == len(digits) {
			res = append(res, curStr)

			return
		}

		for _, letter := range digitsToChar[digits[i]] {
			backtrack(i+1, curStr+letter)
		}
	}

	backtrack(0, "")

	return res
}
