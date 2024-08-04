package merge_sort_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{
			name: "simple array",
			arr:  []int{5, 3},
			want: []int{3, 5},
		},
		{
			name: "simple array",
			arr:  []int{5, 3, 2, 1},
			want: []int{1, 2, 3, 5},
		},
		{
			name: "array has not unique values",
			arr:  []int{5, 1, 1, 2, 0, 0},
			want: []int{0, 0, 1, 1, 2, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mergeSort(tt.arr))
		})
	}
}

// Time O(n*log(n)), since we divide array to sub-arrays until sub-array become 1 length,
// and then merge each sub array in sorted order
// Space O(n), since we use sub-arrays that sum is equal an input.
func mergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	subArr1 := mergeSort(arr[:len(arr)/2])
	subArr2 := mergeSort(arr[len(arr)/2:])

	return merge(subArr1, subArr2)
}

func merge(subArr1, subArr2 []int) []int {
	result := make([]int, 0)
	subArr1pointer, subArr2pointer := 0, 0

	for subArr1pointer < len(subArr1) && subArr2pointer < len(subArr2) {
		if subArr1[subArr1pointer] < subArr2[subArr2pointer] {
			result = append(result, subArr1[subArr1pointer])
			subArr1pointer++
		} else {
			result = append(result, subArr2[subArr2pointer])
			subArr2pointer++
		}
	}

	if subArr1pointer < len(subArr1) {
		result = append(result, subArr1[subArr1pointer:]...)
	}

	if subArr2pointer < len(subArr2) {
		result = append(result, subArr2[subArr2pointer:]...)
	}

	return result
}
