package quick_sort_test

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
		{
			name:  "Case 3",
			array: []int{4, 3, 1, 2, 5, 9, 7, 6, 10},
			want:  []int{1, 2, 3, 4, 5, 6, 7, 9, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			quickSort(tt.array)
			assert.Equal(t, tt.want, tt.array)
		})
	}
}

// Time O(n*log(n)), on average, the pivot divides the array into two parts, but not necessarily equal
// Space O(log(n)), on average, due to recursion stack depth in balanced partitions.
func quickSort(arr []int) {
	if len(arr) == 0 {
		return
	}

	var i, j int

	for pivot := len(arr) - 1; j <= pivot; j++ {
		if arr[j] <= arr[pivot] {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	quickSort(arr[i:])
	quickSort(arr[:i-1])
}
