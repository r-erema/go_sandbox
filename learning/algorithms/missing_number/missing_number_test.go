package missingnumber_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "number 2 is missed",
			input: []int{3, 0, 1},
			want:  2,
		},
		{
			name:  "number 1 is missed",
			input: []int{0},
			want:  1,
		},
		{
			name:  "number 8 is missed",
			input: []int{9, 6, 4, 2, 3, 5, 7, 0, 1},
			want:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, missingNumber(tt.input))
			assert.Equal(t, tt.want, missingNumber2(tt.input))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func missingNumber(nums []int) int {
	res := len(nums)

	for i := range nums {
		res += i - nums[i]
	}

	return res
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func missingNumber2(nums []int) int {
	var res int

	for i := range nums {
		res ^= (i + 1) ^ nums[i]
	}

	return res
}
