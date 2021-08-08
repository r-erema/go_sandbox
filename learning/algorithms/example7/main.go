package example7

/*
	Average, Worst: O(n) time | O(d) space, max depth of array
*/
func ProductSum(array []interface{}) int {
	return helper(array, 1)
}

func helper(array []interface{}, depth int) (sum int) {
	for _, element := range array {
		switch v := element.(type) {
		case []interface{}:
			sum += helper(v, depth+1)
		case int:
			sum += v
		}
	}

	return sum * depth
}
