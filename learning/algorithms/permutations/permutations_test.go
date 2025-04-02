package permutations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "2 subsets",
			nums: []int{0, 1},
			want: [][]int{{0, 1}, {1, 0}},
		},
		{
			name: "3 subsets",
			nums: []int{1, 2, 3},
			want: [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, permute(tt.nums))
		})
	}
}

// Time O(n * n!), the total number of permutations generated is n!, since we're generating all possible permutations of nums,
// additionally, creating the newSet is O(n) per permutation (due to the copy operation)
// Space O(n * n!), the output contains n! permutations, each of length n. Thus, the output itself occupies O(n * n!) space.
func permute(nums []int) [][]int {
	res := [][]int{{}}

	var set []int

	for len(res[0]) < len(nums) {
		set, res = res[0], res[1:]

		for _, n := range nums {
			if !contains(n, set) {
				newSet := make([]int, len(set)+1)
				copy(newSet, set)
				newSet[len(newSet)-1] = n

				res = append(res, newSet)
			}
		}
	}

	return res
}

func contains(needle int, nums []int) bool {
	for _, num := range nums {
		if num == needle {
			return true
		}
	}

	return false
}
