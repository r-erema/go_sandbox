package combinationsum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinationSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		candidates []int
		target     int
		want       [][]int
	}{
		{
			name:       "basic case",
			candidates: []int{2, 3, 6, 7},
			target:     7,
			want:       [][]int{{2, 2, 3}, {7}},
		},
		{
			name:       "basic case 2",
			candidates: []int{2, 4, 8},
			target:     8,
			want:       [][]int{{2, 2, 2, 2}, {2, 2, 4}, {4, 4}, {8}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, combinationSum(tt.candidates, tt.target))
		})
	}
}

// Time O(2^n), or more precisely O(len(candidates)^[target/min(candidates)], since the recursion tree can have a depth
// of target/min_candidate and each level can branch up to n ways
// Space O(2^n), or more precisely O(2^[target/min(candidates)]), due to storage of all combinations

func combinationSum(candidates []int, target int) [][]int {
	var (
		res [][]int
		set []int
	)

	var backtrack func(start, total int)
	backtrack = func(start, total int) {
		if total > target {
			return
		}

		if total == target {
			newSet := make([]int, len(set))
			copy(newSet, set)
			res = append(res, newSet)

			return
		}

		for i := start; i < len(candidates); i++ {
			set = append(set, candidates[i])
			backtrack(i, total+candidates[i])

			set = set[:len(set)-1]
		}
	}

	backtrack(0, 0)

	return res
}
