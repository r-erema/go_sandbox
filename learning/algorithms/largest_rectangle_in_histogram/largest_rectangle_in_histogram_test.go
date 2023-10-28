package largest_rectangle_in_histogram_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLargestRectangleInHistogram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		heights []int
		want    int
	}{
		{
			name:    "Normal plot",
			heights: []int{2, 1, 5, 6, 2, 3},
			want:    10,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, largestRectangleArea(tt.heights))
		})
	}
}

// Time O(n), since we iterate the input once
// Space O(n), since we use only stack structure which couldn't be greater then input.
func largestRectangleArea(heights []int) int {
	var (
		stack   = make([][2]int, 0)
		row     [2]int
		maxArea int
	)

	for i, height := range heights {
		start := i

		for len(stack) > 0 && height < stack[len(stack)-1][1] {
			row, stack = stack[len(stack)-1], stack[:len(stack)-1]
			maxArea = max(maxArea, row[1]*(i-row[0]))
			start = row[0]
		}

		stack = append(stack, [2]int{start, height})
	}

	for _, row = range stack {
		maxArea = max(maxArea, row[1]*(len(heights)-row[0]))
	}

	return maxArea
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}

	return n2
}
