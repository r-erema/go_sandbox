package example9

/*
	Average, Worst: O(n log m) time | O(1) space
*/
func NLargestElements(array []int, n int) []int {
	result := make([]int, n)
	lastIndex := n - 1

	shift := func(shiftEndIndex int, elementToUpdate int) {
		for j := 0; j < shiftEndIndex; j++ {
			result[j] = result[j+1]
		}

		result[shiftEndIndex] = elementToUpdate
	}

	for _, element := range array {
		for i := lastIndex; i >= 0; i-- {
			if element > result[i] {
				shift(i, element)

				break
			}
		}
	}

	return result
}
