package mergesortedarray_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeSortedArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		nums1, nums2 []int
		m, n         int
		want         []int
	}{
		{
			name:  "2 normal arrays",
			nums1: []int{1, 2, 3, 0, 0, 0},
			m:     3,
			nums2: []int{2, 5, 6},
			n:     3,
			want:  []int{1, 2, 2, 3, 5, 6},
		},
		{
			name:  "second array is empty",
			nums1: []int{1},
			m:     1,
			nums2: []int{},
			n:     0,
			want:  []int{1},
		},
		{
			name:  "first array is empty",
			nums1: []int{0},
			m:     0,
			nums2: []int{1},
			n:     1,
			want:  []int{1},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			merge(tt.nums1, tt.m, tt.nums2, tt.n)
			assert.Equal(t, tt.want, tt.nums1)
		})
	}
}

// Time O(m+n), since we iterate once both arrays
//
// Space O(1), since we don't involve any extra space.
func merge(nums1 []int, m int, nums2 []int, n int) {
	for n > 0 {
		if m > 0 && nums1[m-1] > nums2[n-1] {
			nums1[n+m-1] = nums1[m-1]
			m--
		} else {
			nums1[n+m-1] = nums2[n-1]
			n--
		}
	}
}
