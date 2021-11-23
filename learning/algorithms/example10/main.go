package example10

/*
Average, Worst: O(n^2) time | O(1) space.
*/
func InsertionSort(array *[]int) {
	for i := 0; i < len(*array); i++ {
		for j := 0; j < len((*array)[i:])-1; j++ {
			if (*array)[j] > (*array)[j+1] {
				(*array)[j], (*array)[j+1] = (*array)[j+1], (*array)[j]
			}
		}
	}
}
