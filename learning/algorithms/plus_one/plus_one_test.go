package plusone_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlusOne(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input, want []int
	}{
		{
			name:  "no increasing the next digit",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 4},
		},
		{
			name:  "increasing the next digit",
			input: []int{4, 3, 9, 9},
			want:  []int{4, 4, 0, 0},
		},
		{
			name:  "adding one more digit",
			input: []int{9, 9},
			want:  []int{1, 0, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, plusOne(tt.input))
		})
	}
}

// Time O(n), since we iterate each element of input
// Space O(1), since we don't involve any additional data structure.
func plusOne(digits []int) []int {
	digits[len(digits)-1]++
	for i := len(digits) - 1; digits[i] == 10; i-- {
		digits[i] = 0
		if i-1 < 0 {
			return append([]int{1}, digits...)
		}

		digits[i-1]++
	}

	return digits
}
