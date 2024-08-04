package topkfrequentelements_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopKFrequentElements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{
			name: "Normal array",
			nums: []int{1, 1, 1, 2, 2, 3},
			k:    2,
			want: []int{1, 2},
		},
		{
			name: "One element array",
			nums: []int{1},
			k:    1,
			want: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, topKFrequent(tt.nums, tt.k))
		})
	}
}

// Time O(n) since we have single loops
// Time O(N+M) since we have a map N containing counts of each number and slice M grouped numbers by count.
func topKFrequent(nums []int, frequency int) []int {
	numsCountMap := make(map[int]int, 0)

	for _, num := range nums {
		numsCountMap[num]++
	}

	countsArr := make([][]int, len(nums)+1)

	for num, count := range numsCountMap {
		countsArr[count] = append(countsArr[count], num)
	}

	var result []int
	for i := len(countsArr) - 1; i > 0; i-- {
		result = append(result, countsArr[i]...)

		if len(result) >= frequency {
			return result[:frequency]
		}
	}

	return result
}
