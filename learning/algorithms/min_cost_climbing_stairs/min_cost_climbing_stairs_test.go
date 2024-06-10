package mincostclimbingstairs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cost []int
		want int
	}{
		{
			name: "3 costs",
			cost: []int{10, 15, 20},
			want: 15,
		},
		{
			name: "10 costs",
			cost: []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1},
			want: 6,
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, minCostClimbingStairs(testCase.cost))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func minCostClimbingStairs(cost []int) int {
	for i := len(cost) - 3; i >= 0; i-- {
		cost[i] += min(cost[i+1], cost[i+2])
	}

	return min(cost[0], cost[1])
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}

	return n2
}
