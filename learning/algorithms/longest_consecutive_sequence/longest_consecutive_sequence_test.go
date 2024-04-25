package longestconsecutivesequence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestConsecutive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			name: "6 nums",
			nums: []int{100, 4, 200, 1, 3, 2},
			want: 4,
		},
		{
			name: "10 nums",
			nums: []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1},
			want: 9,
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, longestConsecutive(testCase.nums))
		})
	}
}

func longestConsecutive(nums []int) int {
	numsMap := make(map[int]struct{})
	for i := range nums {
		numsMap[nums[i]] = struct{}{}
	}

	var result int

	for num := range numsMap {
		if _, ok := numsMap[num-1]; !ok {
			i := 1

			for {
				if _, ok := numsMap[num+i]; !ok {
					break
				}
				i++
			}

			result = max(i, result)
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
