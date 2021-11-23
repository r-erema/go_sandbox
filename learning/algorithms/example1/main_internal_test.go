package example1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		numbers        []int
		targetSum      int
		expectedResult [2]int
	}{
		{
			name:           "Test case 0",
			numbers:        []int{3, 5, -4, 8, 11, -1, 6},
			targetSum:      10,
			expectedResult: [2]int{11, -1},
		},
		{
			name:           "Test case 1",
			numbers:        []int{-305, 3, 5, -4, 8, 312, 11, -1, 6, 7, 11, -4, 5, 8, 0, 0, 0, 9, -2, 2},
			targetSum:      7,
			expectedResult: [2]int{-305, 312},
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := sumLinear(testCase.numbers, testCase.targetSum)
			assert.ElementsMatch(t, testCase.expectedResult, result)
			result = sumHashTable(testCase.numbers, testCase.targetSum)
			assert.ElementsMatch(t, testCase.expectedResult, result)
			result = sumShiftingPointer(testCase.numbers, testCase.targetSum)
			assert.ElementsMatch(t, testCase.expectedResult, result)
		})
	}
}
