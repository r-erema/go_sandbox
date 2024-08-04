package containsduplicate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsDuplicate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want bool
	}{
		{
			name: "2 duplicates exist",
			arr:  []int{1, 2, 3, 1},
			want: true,
		},
		{
			name: "Duplicates do not exist",
			arr:  []int{1, 2, 3, 4},
			want: false,
		},
		{
			name: "Multiple duplicates exist",
			arr:  []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, containsDuplicate(tt.arr))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(N), since we involve map.
func containsDuplicate(nums []int) bool {
	visitedNumbers := make(map[int]struct{})

	for _, num := range nums {
		if _, ok := visitedNumbers[num]; ok {
			return true
		}

		visitedNumbers[num] = struct{}{}
	}

	return false
}
