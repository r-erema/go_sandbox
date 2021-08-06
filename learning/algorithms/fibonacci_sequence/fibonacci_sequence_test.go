package fibonaccisequence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibonacci(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n, want int
	}{
		{
			name: "Case 0",
			n:    3,
			want: 1,
		},
		{
			name: "Case 1",
			n:    10,
			want: 34,
		},
		{
			name: "Case 2",
			n:    0,
			want: 0,
		},
		{
			name: "Case 3",
			n:    -4,
			want: 0,
		},
		{
			name: "Case 4",
			n:    30,
			want: 514229,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, fibonacciRecursive(testCase.n))
			assert.Equal(t, testCase.want, fibonacciCache(testCase.n))
			assert.Equal(t, testCase.want, fibonacciIterative(testCase.n))
		})
	}
}

const fibPrev, fibNext = 1, 2

/*
Average, Worst: O(n^2) time | O(n) space.
*/
func fibonacciRecursive(number int) int {
	if number < fibPrev {
		return 0
	}

	if number == fibNext {
		return 1
	}

	return fibonacciRecursive(number-fibPrev) + fibonacciRecursive(number-fibNext)
}

/*
Average, Worst: O(n) time | O(n) space.
*/
func fibonacciCache(number int) int {
	if number < fibPrev {
		return 0
	}

	return cacheHelper(number, map[int]int{1: 0, 2: 1})
}

func cacheHelper(number int, cache map[int]int) int {
	if _, ok := cache[number]; !ok {
		cache[number] = cacheHelper(number-fibPrev, cache) + cacheHelper(number-fibNext, cache)
	}

	return cache[number]
}

/*
Average, Worst: O(n) time | O(1) space.
*/
func fibonacciIterative(number int) int {
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
