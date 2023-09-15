package dailytemperatures_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDailyTemperatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		temperatures []int
		want         []int
	}{
		{
			name:         "Random temperatures",
			temperatures: []int{73, 74, 75, 71, 69, 72, 76, 73},
			want:         []int{1, 1, 4, 2, 1, 1, 0, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, dailyTemperatures(tt.temperatures))
		})
	}
}

// Time: O(n), since we handle each element only within a stack which can't be greater than the input
// Memory: O(n), since we involve a stack that can't be greater than the input.
func dailyTemperatures(temperatures []int) []int {
	res := make([]int, len(temperatures))
	stack := make([][2]int, 0)

	var stackIndex int

	for i, temperature := range temperatures {
		for len(stack) > 0 && temperature > stack[len(stack)-1][1] {
			stackIndex, stack = stack[len(stack)-1][0], stack[:len(stack)-1]
			res[stackIndex] = i - stackIndex
		}

		stack = append(stack, [2]int{i, temperature})
	}

	return res
}
