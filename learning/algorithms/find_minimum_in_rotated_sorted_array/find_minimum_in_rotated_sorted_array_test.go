package findminimuminrotatedsortedarray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			name: "middle element in the array is searched",
			nums: []int{3, 4, 5, 1, 2},
			want: 1,
		},
		{
			name: "middle element in the longer array is searched",
			nums: []int{4, 5, 6, 7, 0, 1, 2},
			want: 0,
		},
		{
			name: "first element in the array is searched",
			nums: []int{11, 13, 15, 17},
			want: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, findMin(tt.nums))
		})
	}
}

// Time O(log*N), since the binary search
// Space O(1), since we don't involve any additional space.
func findMin(nums []int) int {
	left, right, res := 0, len(nums)-1, nums[0]

	for left <= right {
		mid := (left + right) / 2

		if nums[mid] >= res {
			left = mid + 1
		} else {
			right = mid - 1
		}

		res = min(res, nums[mid])
	}

	return res
}
