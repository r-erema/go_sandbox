package twosumiiinputarrayissorted_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		numbers []int
		target  int
		want    []int
	}{
		{
			name:    "Simple array",
			numbers: []int{1, 2, 10, 17},
			target:  12,
			want:    []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, twoSum(tt.numbers, tt.target))
		})
	}
}

// Time O(n), since in the worst case we don't exceed iterations count more then input
// Space O(1), since we don't involve an additional data structure.
func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1

	for left <= right {
		sum := numbers[left] + numbers[right]
		if sum > target {
			right--
		} else if sum < target {
			left++
		} else {
			return []int{left + 1, right + 1}
		}
	}

	return nil
}
