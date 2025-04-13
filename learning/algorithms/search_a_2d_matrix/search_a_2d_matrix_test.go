package searcha2dmatrix_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		{
			name:   "target exists in matrix",
			matrix: [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}},
			target: 3,
			want:   true,
		},
		{
			name:   "target does not exist in matrix",
			matrix: [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}},
			target: 13,
			want:   false,
		},
		{
			name:   "target does not exist in 1 element matrix",
			matrix: [][]int{{1}},
			target: 1,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, searchMatrix(tt.matrix, tt.target))
		})
	}
}

// Time O(log(N*M)), since we apply binary search to matrix rows and then apply binary search in the found row
// Space O(1), we don't involve any additional data structures.
func searchMatrix(matrix [][]int, target int) bool {
	var row int

	top, bottom := 0, len(matrix)-1

search:
	for top <= bottom {
		row = (top + bottom) / 2

		switch {
		case target > matrix[row][len(matrix[0])-1]:
			top = row + 1
		case target < matrix[row][0]:
			bottom = row - 1
		default:
			break search
		}
	}

	left, right := 0, len(matrix[row])-1
	for left <= right {
		pointer := (left + right) / 2

		switch cmp := matrix[row][pointer]; {
		case cmp < target:
			left = pointer + 1
		case cmp > target:
			right = pointer - 1
		default:
			return true
		}
	}

	return false
}
