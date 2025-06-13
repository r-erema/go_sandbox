package palindromepartitioning_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPalindromePartitioning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want [][]string
	}{
		{
			name: "1 palindrome case",
			s:    "a",
			want: [][]string{{"a"}},
		},
		{
			name: "2 palindromes case",
			s:    "aab",
			want: [][]string{{"a", "a", "b"}, {"aa", "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, partition(tt.s))
		})
	}
}

// Time O(n∗2^n), in worst case when all substrings are palindromes e.g., "aaa"
// Space O(n∗2^n), due to storing all possible partitions (each up to O(n) size, O(2^n) of them).
// Recursion stack is O(n).
func partition(str string) [][]string {
	var (
		res  [][]string
		part []string
		dfs  func(i int)
	)

	dfs = func(i int) {
		if i >= len(str) {
			res = append(res, slices.Clone(part))
		}

		for j := i; j < len(str); j++ {
			if isPalindrome(str[i : j+1]) {
				part = append(part, str[i:j+1])
				dfs(j + 1)

				part = part[:len(part)-1]
			}
		}
	}

	dfs(0)

	return res
}

func isPalindrome(str string) bool {
	left, right := 0, len(str)-1

	for left <= right {
		if str[left] != str[right] {
			return false
		}

		left++
		right--
	}

	return true
}
