package example13

/*
	Average, Worst: O(n) time | O(1) space
*/
func IsPalindrome(str string) bool {
	firstIndex, lastIndex := 0, len(str)-1

	for firstIndex < lastIndex {
		if str[firstIndex] != str[lastIndex] {
			return false
		}
		firstIndex++
		lastIndex--
	}

	return true
}
