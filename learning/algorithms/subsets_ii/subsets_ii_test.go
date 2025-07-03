package subsetsii_test

import (
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsetsII(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "2 subsets",
			nums: []int{0},
			want: [][]int{{}, {0}},
		},
		{
			name: "4 subsets",
			nums: []int{1, 2},
			want: [][]int{{}, {1}, {1, 2}, {2}},
		},
		{
			name: "6 subsets",
			nums: []int{1, 2, 2},
			want: [][]int{{}, {1}, {1, 2}, {1, 2, 2}, {2}, {2, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, subsetsWithDup(tt.nums))
		})
	}
}

// Time O(n * 2^n), in the worst case, we explore all possible subsets - O(2^n), and for each subset,
// we perform O(n) copy operations
// Time O(n * 2^n), O(2^n) space for the output list, and the recursion depth is O(n)

func subsetsWithDup(nums []int) [][]int {
	res := [][]int{{}}

	var (
		curr []int
		dfs  func(int)
	)

	sort.Ints(nums)

	dfs = func(i int) {
		for ; i < len(nums); i++ {
			curr = append(curr, nums[i])

			res = append(res, slices.Clone(curr))

			dfs(i + 1)

			curr = curr[:len(curr)-1]

			for i+1 < len(nums) && nums[i] == nums[i+1] {
				i++
			}
		}
	}

	dfs(0)

	return res
}
