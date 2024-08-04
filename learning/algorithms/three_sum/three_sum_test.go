package threesum_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test3Sum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "Two sets",
			nums: []int{-1, 0, 1, 2, -1, -4},
			want: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, threeSum(tt.nums))
		})
	}
}

// Time O(N^2), since we need to iterate an input for each number
// Time O(1) or O(n) depends on sorting algorithm.
func threeSum(nums []int) [][]int { //nolint:cyclop
	sort.Ints(nums)

	var res [][]int

	for i := range len(nums) - 2 {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		left, right := i+1, len(nums)-1

		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum > 0 {
				right--
			} else if sum < 0 {
				left++
			} else {
				res = append(res, []int{nums[i], nums[left], nums[right]})
				left++
				right--

				for left < right && nums[left] == nums[left-1] {
					left++
				}

				for left < right && nums[right] == nums[right+1] {
					right--
				}
			}
		}
	}

	return res
}
