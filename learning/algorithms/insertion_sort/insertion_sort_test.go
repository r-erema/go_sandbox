package insertionsort_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertionSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		array []int
		want  []int
	}{
		{
			name:  "Case 0",
			array: []int{7, 0, 4},
			want:  []int{0, 4, 7},
		},
		{
			name:  "Case 1",
			array: []int{8, 5, 2, 9, 5, 6, 3},
			want:  []int{2, 3, 5, 5, 6, 8, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			insertionSort(tt.array)
			assert.Equal(t, tt.want, tt.array)
		})
	}
}

// Time O(n^2), since it's possible that we iterate the input twice for the each number
// Space O(1), since we don't allocate any additional memory.
func insertionSort(array []int) {
	for i := 1; i < len(array); i++ {
		for j := range i {
			if array[i] < array[j] {
				array[i], array[j] = array[j], array[i]
			}
		}
	}
}
