package twonumbersum_test

import (
	"sort"
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

func sumLinear(numbers []int, targetSum int) [2]int {
	for i, number1 := range numbers {
		for _, number2 := range numbers[i+1:] {
			if number1+number2 == targetSum {
				return [2]int{number1, number2}
			}
		}
	}

	return [2]int{}
}

func sumHashTable(numbers []int, targetSum int) [2]int {
	hashTable := make(map[int]struct{})

	for _, number := range numbers {
		neededNumber := targetSum - number
		if _, ok := hashTable[neededNumber]; ok {
			return [2]int{number, neededNumber}
		}

		hashTable[number] = struct{}{}
	}

	return [2]int{}
}

func sumShiftingPointer(numbers []int, targetSum int) [2]int {
	leftPointer, rightPointer := 0, len(numbers)-1
	sort.Ints(numbers)

	for leftPointer < rightPointer {
		currentSum := numbers[leftPointer] + numbers[rightPointer]

		if currentSum == targetSum {
			return [2]int{numbers[leftPointer], numbers[rightPointer]}
		}

		if currentSum < targetSum {
			leftPointer++
		}

		if currentSum > targetSum {
			rightPointer--
		}
	}

	return [2]int{}
}
