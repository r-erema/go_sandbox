package example8

/*
Average, Worst: O(log n) time | O(1) space.
*/
func BinarySearch(array []int, needle int) int {
	leftPointer, rightPointer := 0, len(array)-1

	for leftPointer <= rightPointer {
		cutPoint := (leftPointer + rightPointer) / 2
		potentialResult := array[cutPoint]

		if potentialResult > needle {
			rightPointer = cutPoint - 1
		} else if potentialResult < needle {
			leftPointer = cutPoint + 1
		} else {
			return cutPoint
		}
	}

	return -1
}
