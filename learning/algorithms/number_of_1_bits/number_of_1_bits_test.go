package numberof1bits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHammingWeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input, want int
	}{
		{
			name:  "3 ones",
			input: 11,
			want:  3,
		},
		{
			name:  "1 one",
			input: 128,
			want:  1,
		},
		{
			name:  "30 ones",
			input: 2147483645,
			want:  30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, hammingWeight(tt.input))
			assert.Equal(t, tt.want, hammingWeight2(tt.input))
			assert.Equal(t, tt.want, hammingWeight3(tt.input))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func hammingWeight(number int) int {
	var res int

	for number > 0 {
		res += number % 2
		number >>= 1
	}

	return res
}

func hammingWeight2(number int) int {
	var res int

	for number > 0 {
		number &= number - 1
		res++
	}

	return res
}

func hammingWeight3(number int) int {
	var res int

	for number > 0 {
		if number&1 == 1 {
			res++
		}

		number >>= 1
	}

	return res
}
