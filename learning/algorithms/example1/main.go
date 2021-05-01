package example1

import "sort"

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
