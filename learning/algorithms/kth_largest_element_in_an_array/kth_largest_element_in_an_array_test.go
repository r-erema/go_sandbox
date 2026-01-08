package kthlargestelementinanarray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKthLargest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{
			name: "ordinary input",
			nums: []int{3, 2, 1, 5, 6, 4},
			k:    2,
			want: 5,
		},
		{
			name: "input with duplicates",
			nums: []int{1, 2, 1, 1, 1, -1, -2},
			k:    1,
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, findKthLargest(tt.nums, tt.k))
		})
	}
}

// Time O(n) in average case, random pivots split the array roughly in half each time,
// leading to a geometric series of work: n + n/2 + n/4 + … = O(n).
// O(n^2) in worst case since the pivot is always the smallest or largest element,
// causing each partition to reduce the problem size by only one,
// leading to n + (n−1) + (n−2) + … + 1 = O(n²) total operations
//
// Space O(1), since we use only an initial array.
func findKthLargest(arr []int, k int) int {
	var pointer int

	for j, pivot := 0, len(arr)-1; j <= pivot; j++ {
		if arr[j] <= arr[pivot] {
			arr[pointer], arr[j] = arr[j], arr[pointer]
			pointer++
		}
	}

	if pointer == len(arr) {
		pointer--
	}

	// optimize for duplicates at the beginning of the array
	for len(arr) > 1 && arr[0] == arr[1] {
		arr = arr[1:]
		pointer--
	}

	if len(arr) == 1 {
		return arr[0]
	}

	// find in the right part
	if pointer <= len(arr)-k {
		return findKthLargest(arr[pointer:], k)
	}

	// find in the left part
	if k <= 1 {
		return arr[len(arr)-1]
	}

	k -= len(arr) - len(arr[:pointer])

	return findKthLargest(arr[:pointer], k)
}
