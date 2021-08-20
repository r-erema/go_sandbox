package example12

/*
	Average, Worst: O(n^2) time | O(1) space
*/
func SelectionSort(array *[]int) {
	a := *array
	startIndex, length := 0, len(a)

	for startIndex < length {
		minIndex := startIndex
		for i := startIndex; i < len(a); i++ {
			if a[minIndex] > a[i] {
				minIndex = i
			}
		}

		a[startIndex], a[minIndex] = a[minIndex], a[startIndex]
		startIndex++
	}
}
