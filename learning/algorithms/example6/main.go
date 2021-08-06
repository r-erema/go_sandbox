package example6

const fibPrev, fibNext = 1, 2

/*
	Average, Worst: O(n^2) time | O(n) space
*/
func FibonacciRecursive(n int) int {
	if n < fibPrev {
		return 0
	}

	if n == fibNext {
		return 1
	}

	return FibonacciRecursive(n-fibPrev) + FibonacciRecursive(n-fibNext)
}

/*
	Average, Worst: O(n) time | O(n) space
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
	Average, Worst: O(n) time | O(1) space
*/
func FibonacciIterative(n int) int {
	if n < fibPrev {
		return 0
	}

	if n == fibNext {
		return 1
	}

	result, prev, current := 0, 0, 1
	for i := 2; i < n; i++ {
		result = prev + current
		prev = current
		current = result
	}

	return result
}
