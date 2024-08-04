package subsets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{
			name: "8 subsets",
			nums: []int{1, 2, 3},
			want: [][]int{{}, {3}, {2}, {2, 3}, {1}, {1, 3}, {1, 2}, {1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, subsets(tt.nums))
		})
	}
}

// Time O(N * 2^N) it's a brute force way, there are 2^N subsets in total.
// Time O(N), since we allocate memory for each backtracking call and overwrite it on the next backtracking call,
// N is a number of calls.
func subsets(nums []int) [][]int {
	backtrack := func(num int, subs [][]int) [][]int {
		var newSubs [][]int

		for _, sub := range subs {
			newSub := make([]int, len(sub))
			copy(newSub, sub)
			sub = append(sub, num)
			newSubs = append(newSubs, newSub, sub)
		}

		return newSubs
	}

	subs := [][]int{{}}
	for i := range len(nums) {
		subs = backtrack(nums[i], subs)
	}

	return subs
}
