package searchinrotatedsortedarray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchInRotatedSortedArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  []int
		target int
		want   int
	}{
		{
			name:   "target exists",
			input:  []int{4, 5, 6, 7, 0, 1, 2},
			target: 0,
			want:   4,
		},
		{
			name:   "target does not exist",
			input:  []int{4, 5, 6, 7, 0, 1, 2},
			target: 3,
			want:   -1,
		},
		{
			name:   "target does not exist, input with 1 value",
			input:  []int{1},
			target: 0,
			want:   -1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, search(tt.input, tt.target))
		})
	}
}

// Time O(log(N)), since it's binary search
// Space O(1), sine we don't allocate any additional memory.
func search(nums []int, target int) int {
	left, right := 0, len(nums)-1

	isLeftPortionSorted := func(mid int) bool {
		return nums[left] <= nums[mid]
	}

	handleLeftPortion := func(mid int) {
		targetInLeftPortion := nums[left] <= target && target < nums[mid]
		if targetInLeftPortion {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	handleRightPortion := func(mid int) {
		targetInRightPortion := nums[mid] < target && target <= nums[right]
		if targetInRightPortion {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	for left <= right {
		mid := (left + right) / 2

		if target == nums[mid] {
			return mid
		}

		if isLeftPortionSorted(mid) {
			handleLeftPortion(mid)
		} else {
			handleRightPortion(mid)
		}
	}

	return -1
}
