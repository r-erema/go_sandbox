package findtheduplicatenumber_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindDuplicate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			name: "2 duplicates",
			nums: []int{1, 3, 4, 5, 2, 2},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, findDuplicate(tt.nums))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func findDuplicate(nums []int) int {
	slow, fast := nums[0], nums[nums[0]]

	for slow != fast {
		slow, fast = nums[slow], nums[nums[fast]]
	}

	for slow = 0; slow != fast; {
		slow, fast = nums[slow], nums[fast]
	}

	return slow
}
