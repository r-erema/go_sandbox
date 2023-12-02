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
			name: "the original array was [1,2,3,4,5] rotated 3 times",
			nums: []int{3, 4, 5, 1, 2},
			want: 1,
		},
		{
			name: "the original array was [0,1,2,4,5,6,7] and it was rotated 4 times",
			nums: []int{4, 5, 6, 7, 0, 1, 2},
			want: 0,
		},
		{
			name: "the original array was [11,13,15,17] and it was rotated 4 times",
			nums: []int{11, 13, 15, 17},
			want: 11,
		},
		{
			name: "the original array was [12, 14, 1, 2, 3, 4, 5, 6, 8, 11] and it was rotated 2 times",
			nums: []int{12, 14, 1, 2, 3, 4, 5, 6, 8, 11},
			want: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, findMin(tt.nums))
		})
	}
}

// Time O(log*N), since the binary search
// Space O(1), since we don't involve any additional space.
func findMin(nums []int) int {
	left, right := 0, len(nums)-1
	result := nums[0]

	for left <= right {
		if nums[left] < nums[right] {
			result = min(result, nums[left])

			break
		}

		middle := (left + right) / 2

		result = min(result, nums[middle])

		if nums[middle] >= nums[left] {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
