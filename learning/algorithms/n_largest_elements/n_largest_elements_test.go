package nlargestelements_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNLargestElements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		array []int
		n     int
		want  []int
	}{
		{
			name:  "Case 0",
			array: []int{141, 1, 17, -7, -17, -27, 18, 541, 8, 7, 7},
			n:     3,
			want:  []int{18, 141, 541},
		},
		{
			name:  "Case 1",
			array: []int{141, 1, 17, -7, -17, -27, 18, 541, 8, 7, 7},
			n:     4,
			want:  []int{17, 18, 141, 541},
		},
		{
			name:  "Case 2",
			array: []int{141, 1, 17, -7, -17, -27, 18, 541, 8, 7, 7},
			n:     0,
			want:  []int{},
		},
		{
			name:  "Case 3",
			array: []int{141, 1, 17, -7, -17, -27, 18, 541, 8, 7, 7},
			n:     1,
			want:  []int{541},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, nLargestElements(tt.array, tt.n))
		})
	}
}

/*
Average, Worst: O(n log m) time | O(1) space.
*/
func nLargestElements(array []int, n int) []int {
	result := make([]int, n)
	lastIndex := n - 1

	shift := func(shiftEndIndex int, elementToUpdate int) {
		for j := range shiftEndIndex {
			result[j] = result[j+1]
		}

		result[shiftEndIndex] = elementToUpdate
	}

	for _, element := range array {
		for i := lastIndex; i >= 0; i-- {
			if element > result[i] {
				shift(i, element)

				break
			}
		}
	}

	return result
}
