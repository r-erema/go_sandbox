package example6

const fibPrev, fibNext = 1, 2

/*
Average, Worst: O(n^2) time | O(n) space.
*/
func FibonacciRecursive(number int) int {
	if number < fibPrev {
		return 0
	}

	if number == fibNext {
		return 1
	}

	return FibonacciRecursive(number-fibPrev) + FibonacciRecursive(number-fibNext)
}

/*
Average, Worst: O(n) time | O(n) space.
*/
func FibonacciCache(n int) int {
	if n < fibPrev {
		return 0
	}

	return cacheHelper(n, map[int]int{1: 0, 2: 1})
}

func cacheHelper(n int, cache map[int]int) int {
	if _, ok := cache[n]; !ok {
		cache[n] = cacheHelper(n-fibPrev, cache) + cacheHelper(n-fibNext, cache)
	}

	return cache[n]
}

/*
Average, Worst: O(n) time | O(1) space.
*/
func FibonacciIterative(number int) int {
	if number < fibPrev {
		return 0
	}

	if number == fibNext {
		return 1
	}

	result, prev, current := 0, 0, 1
	for i := 2; i < number; i++ {
		result = prev + current
		prev = current
		current = result
	}

	return result
}
