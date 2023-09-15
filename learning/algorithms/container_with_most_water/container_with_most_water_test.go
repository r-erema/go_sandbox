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

// Time O(n)
// Since we traverse each height
//
// Space O(1)
// We don't use any extra space for memorizing input, etc.
func maxArea(heights []int) int {
	maxLeftSideIndex, maxRightSideIndex, maxAreaResult := 0, len(heights)-1, 0

	for maxLeftSideIndex <= maxRightSideIndex {
		leftHeight, rightHeight := heights[maxLeftSideIndex], heights[maxRightSideIndex]
		height := min(leftHeight, rightHeight)
		width := maxRightSideIndex - maxLeftSideIndex
		currentArea := height * width

		if maxAreaResult < currentArea {
			maxAreaResult = currentArea
		}

		if leftHeight <= rightHeight {
			maxLeftSideIndex++
		} else {
			maxRightSideIndex--
		}
	}

	return maxAreaResult
}

func min(number1, number2 int) int {
	if number1 < number2 {
		return number1
	}

	return number2
}
