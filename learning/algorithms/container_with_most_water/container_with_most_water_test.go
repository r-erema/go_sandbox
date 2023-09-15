package containerwithmostwater_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerWithMostWater(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		heights []int
		want    int
	}{
		{
			name:    "9 heights",
			heights: []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
			want:    49,
		},
		{
			name:    "2 heights",
			heights: []int{1, 1},
			want:    1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, maxArea(tt.heights))
		})
	}
}

// Time O(n), since we traverse each height
// Space O(1), We don't use any extra space for memorizing input, etc.
func maxArea(heights []int) int {
	var maxAreaResult int

	left, right := 0, len(heights)-1

	for left < right {
		minSide := min(heights[left], heights[right])
		maxAreaResult = max(maxAreaResult, minSide*(right-left))

		if heights[left] < heights[right] {
			left++
		} else {
			right--
		}
	}

	return maxAreaResult
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}

	return n2
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}

	return n2
}
