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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, minCostClimbingStairs(tt.cost))
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
