package combinationsumii_test

import (
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinationSum2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		candidates []int
		target     int
		want       [][]int
	}{
		{
			name:       "1 result sets",
			candidates: []int{1},
			target:     1,
			want:       [][]int{{1}},
		},
		{
			name:       "1 result sets",
			candidates: []int{1, 2},
			target:     2,
			want:       [][]int{{2}},
		},
		{
			name:       "2 result sets",
			candidates: []int{1, 2, 3, 2},
			target:     3,
			want:       [][]int{{1, 2}, {3}},
		},
		{
			name:       "2 result sets",
			candidates: []int{2, 5, 2, 1, 2},
			target:     5,
			want:       [][]int{{1, 2, 2}, {5}},
		},
		{
			name:       "2 result sets",
			candidates: []int{1, 2, 0},
			target:     2,
			want:       [][]int{{0, 2}, {2}},
		},
		{
			name:       "4 result sets",
			candidates: []int{10, 1, 2, 7, 6, 1, 5},
			target:     8,
			want:       [][]int{{1, 1, 6}, {1, 2, 5}, {1, 7}, {2, 6}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, combinationSum2(tt.candidates, tt.target))
		})
	}
}

// Time O(n*2^n), due to the number of all possible subsets 2^n, and their possible max length n
// Space O(n), due to the recursion stack depth

func combinationSum2(candidates []int, target int) [][]int {
	var (
		result  [][]int
		currSet []int
		total   int
		dfs     func(startIndex int)
	)

	sort.Ints(candidates)

	dfs = func(startIndex int) {
		for i := startIndex; i < len(candidates); i++ {
			currSet = append(currSet, candidates[i])
			total += candidates[i]

			if total < target {
				dfs(i + 1)
			}

			if total == target {
				result = append(result, slices.Clone(currSet))
			}

			total -= currSet[len(currSet)-1]
			currSet = currSet[:len(currSet)-1]

			for i+1 < len(candidates) && candidates[i+1] == candidates[i] {
				i++
			}
		}
	}

	dfs(0)

	return result
}
