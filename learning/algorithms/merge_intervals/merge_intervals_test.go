package mergeintervals_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeIntervals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		intervals [][]int
		want      [][]int
	}{
		{
			name:      "One interval",
			intervals: [][]int{{1, 4}, {4, 5}},
			want:      [][]int{{1, 5}},
		},
		{
			name:      "Two intervals",
			intervals: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			want:      [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "Three intervals",
			intervals: [][]int{{9, 14}, {1, 2}, {4, 6}, {2, 2}, {8, 9}, {2, 2}, {8, 10}, {11, 15}},
			want:      [][]int{{1, 2}, {4, 6}, {8, 15}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mergeIntervals(tt.intervals))
		})
	}
}

// Time O(NlogN)
// Other than the sort invocation,
// we do a simple linear scan of the list,
// so the runtime is dominated by the (NlogN) complexity of sorting
//
// Space O(logN)
// If we can sort intervals in place,
// we do not need more than constant additional space,
// although the sorting itself takes O(logN) space.
func mergeIntervals(intervals [][]int) [][]int {
	sort.SliceStable(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	output := [][]int{intervals[0]}

	for _, interval := range intervals[1:] {
		start, end := interval[0], interval[1]

		lastEnd := output[len(output)-1][1]

		if start <= lastEnd {
			output[len(output)-1][1] = max(lastEnd, end)
		} else {
			output = append(output, []int{start, end})
		}
	}

	return output
}

func max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}

	return num2
}
