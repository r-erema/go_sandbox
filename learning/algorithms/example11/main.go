package example11

/*
	Average, Worst: O(n^2) time | O(1) space
*/
func BubbleSort(array *[]int) {
	a := *array
	boundary := len(a)

	for boundary > 1 {
		for i := 0; i+1 < boundary; i++ {
			if a[i] > a[i+1] {
				a[i], a[i+1] = a[i+1], a[i]
			}
		}
		boundary--
	}
}
