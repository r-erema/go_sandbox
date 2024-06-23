package countingbits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountBits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number int
		want   []int
	}{
		{
			name:   "number 2",
			number: 2,
			want:   []int{0, 1, 1},
		},
		{
			name:   "number 5",
			number: 5,
			want:   []int{0, 1, 1, 2, 1, 2},
		},
		{
			name:   "number 8",
			number: 8,
			want:   []int{0, 1, 1, 2, 1, 2, 2, 3, 1},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, countBits(testCase.number))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we involve array with lengths equals input number.
func countBits(number int) []int {
	dynProg := make([]int, number+1)

	for offset, i := 1, 1; i <= number; i++ {
		if offset*2 == i {
			offset = i
		}

		dynProg[i] = 1 + dynProg[i-offset]
	}

	return dynProg
}
