package binarysearch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		array  []int
		needle int
		want   int
	}{
		{
			name:   "Case 0",
			array:  []int{0, 1, 21, 33, 45, 45, 61, 71, 72, 73},
			needle: 33,
			want:   3,
		},
		{
			name:   "Case 1",
			array:  []int{0},
			needle: 0,
			want:   0,
		},
		{
			name:   "Case 2",
			array:  []int{0},
			needle: 1,
			want:   -1,
		},
		{
			name:   "Case 3",
			array:  []int{},
			needle: 1,
			want:   -1,
		},
		{
			name:   "Case 4",
			array:  []int{3, 7, 55},
			needle: 3,
			want:   0,
		},
		{
			name:   "Case 5",
			array:  []int{3, 7, 55},
			needle: 55,
			want:   2,
		},
		{
			name:   "Case 6",
			array:  []int{0, 1, 21, 33, 45, 45, 61, 71, 72, 73},
			needle: 74,
			want:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, binarySearch(tt.array, tt.needle))
		})
	}
}

/*
Average, Worst: O(log n) time | O(1) space.
*/
func binarySearch(array []int, needle int) int {
	leftPointer, rightPointer := 0, len(array)-1

	for leftPointer <= rightPointer {
		cutPoint := (leftPointer + rightPointer) / 2
		potentialResult := array[cutPoint]

		if potentialResult > needle {
			rightPointer = cutPoint - 1
		} else if potentialResult < needle {
			leftPointer = cutPoint + 1
		} else {
			return cutPoint
		}
	}

	return -1
}
