package lettercombinationsofaphonenumber_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		digits string
		want   []string
	}{
		{
			name:   "2 digits",
			digits: "12",
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

// Space O(?),.
func letterCombinations(digits string) []string {
	return nil
}
