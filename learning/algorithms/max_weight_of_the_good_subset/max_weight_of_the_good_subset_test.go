package maxweightofthegoodsubset_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxWeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{
			name: "subset [1 3 3 3 3]",
			arr:  []int{3, 3, 3, 1, 3, 7, 1},
			want: 13,
		},
		{
			name: "subset [7 15]",
			arr:  []int{1, 7, 3, 15, 2, 5, 2, 1, 4},
			want: 22,
		},
		{
			name: "subset [3 4 5 6 7]",
			arr:  []int{6, 2, 5, 1, 7, 4, 3},
			want: 25,
		},
		{
			name: "subset []",
			arr:  []int{},
			want: 0,
		},
		{
			name: "subset [4 11]",
			arr:  []int{4, 11},
			want: 15,
		},
		{
			name: "subset [4]",
			arr:  []int{4},
			want: 4,
		},
		{
			name: "subset [1 1 1 1 1 1 1 1 1 1 1 1 1 2]",
			arr:  []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 5, 6},
			want: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, maxWeight(tt.arr))
		})
	}
}

// Time O(N^2 * log(N)), since we have a nested loop
// (we need search an appropriate sum pair for each biggest number in iteration) and also we have sorting
// Space O(1), since we don't use any additional space.
func maxWeight(arr []int) int {
	if len(arr) == 1 {
		return arr[0]
	}

	sort.Ints(arr)

	var result int

	for i := len(arr) - 1; i > 0; i-- {
		examineElem := arr[i]
		leftPointer, rightPointer := i-1, i
		intermediateSum := arr[rightPointer]

		for arr[leftPointer]+arr[rightPointer] >= examineElem {
			intermediateSum += arr[leftPointer]
			leftPointer--
			rightPointer--

			if leftPointer < 0 {
				break
			}
		}

		result = max(result, intermediateSum)
	}

	return result
}
