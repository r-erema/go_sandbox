package merge2nondecreasingarrays_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeArrays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		arr1, arr2 []int
		want       []int
	}{
		{
			name: "end on the first array",
			arr1: []int{-2, 3, 3, 22},
			arr2: []int{-5, 0},
			want: []int{-5, -2, 0, 3, 3, 22},
		},
		{
			name: "end on the second array",
			arr1: []int{-2, 3},
			arr2: []int{0, 8},
			want: []int{-2, 0, 3, 8},
		},
		{
			name: "first array is empty",
			arr1: []int{},
			arr2: []int{0},
			want: []int{0},
		},
		{
			name: "second array is empty",
			arr1: []int{-1},
			arr2: []int{},
			want: []int{-1},
		},
		{
			name: "both arrays are empty",
			arr1: []int{},
			arr2: []int{},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mergeArrays(tt.arr1, tt.arr2))
		})
	}
}

// Time O(m+n), where m and n are length of the arrays, and we need to iterate each of element
// Time O(m+n), since we create resulted array which equals m+n of lengths of input arrays.
func mergeArrays(arr1, arr2 []int) []int {
	pointer1, pointer2 := 0, 0
	result := make([]int, 0)

	for pointer1 < len(arr1) && pointer2 < len(arr2) {
		if arr1[pointer1] > arr2[pointer2] {
			result = append(result, arr2[pointer2])
			pointer2++
		} else {
			result = append(result, arr1[pointer1])
			pointer1++
		}
	}

	if len(arr1) > pointer1 {
		result = append(result, arr1[pointer1:]...)
	}

	if len(arr2) > pointer2 {
		result = append(result, arr2[pointer2:]...)
	}

	return result
}
