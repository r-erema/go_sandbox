package validpalindrome_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecodeStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Valid palindrome",
			input: "A man, a plan, a canal: Panama",
			want:  true,
		},
		{
			name:  "Not palindrome",
			input: "race a car",
			want:  false,
		},
		{
			name:  "Palindrome: no alphanumerical symbols",
			input: " .",
			want:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isPalindrome(tt.input))
		})
	}
}

// Time O(n), since we move consequently and narrow a search area
// Time O(1), since we don't involve any additional data structure.
func isPalindrome(str string) bool {
	left, right := 0, len(str)-1

	for left < right {
		if !isAlphanumeric(str[left]) {
			left++

			continue
		}

		if !isAlphanumeric(str[right]) {
			right--

			continue
		}

		if toLower(str[left]) != toLower(str[right]) {
			return false
		}

		left++
		right--
	}

	return true
}

func isAlphanumeric(symbol byte) bool {
	return symbol >= 'a' && symbol <= 'z' || symbol >= 'A' && symbol <= 'Z' || symbol >= '0' && symbol <= '9'
}

func toLower(symbol byte) byte {
	if symbol >= 'A' && symbol <= 'Z' {
		return symbol + 32
	}

	return symbol
}
