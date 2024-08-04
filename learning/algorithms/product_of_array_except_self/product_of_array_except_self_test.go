package productofarrayexceptself_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductOfArrayExceptSelf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{
			name: "Simple array",
			nums: []int{1, 2, 3, 4},
			want: []int{24, 12, 8, 6},
		},
		{
			name: "Array with negative numbers",
			nums: []int{-1, 1, 0, -3, 3},
			want: []int{0, 0, 9, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, productExceptSelf(tt.nums))
		})
	}
}

// Time O(N), we iterate linearly through the result array
// Space O(N), we consume only result array which equals input.
func productExceptSelf(nums []int) []int {
	res := make([]int, len(nums))
	prefix := 1

	for i := 0; i < len(res); i++ {
		res[i] = prefix
		prefix *= nums[i]
	}

	postfix := 1
	for i := len(res) - 1; i >= 0; i-- {
		res[i] *= postfix
		postfix *= nums[i]
	}

	return res
}
