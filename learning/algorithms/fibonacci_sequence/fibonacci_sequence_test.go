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
			want: 2,
		},
		{
			name: "Case 1",
			n:    10,
			want: 55,
		},
		{
			name: "Case 2",
			n:    0,
			want: 0,
		},
		{
			name: "Case 3",
			n:    30,
			want: 832040,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, fibonacciRecursive(tt.n))
			assert.Equal(t, tt.want, fibonacciCache(tt.n))
			assert.Equal(t, tt.want, fibonacciIterative(tt.n))
		})
	}
}

const fibPrev, fibNext = 1, 2

// Time O(N^2), since we recount numbers in recursion
// Space O(N), since we use stretch recursion stack.
func fibonacciRecursive(number int) int {
	if number <= 1 {
		return number
	}

	return fibonacciRecursive(number-fibPrev) + fibonacciRecursive(number-fibNext)
}

// Time O(N), since we iterate input one time and get counted result from cache
// Space O(N), since we involve map as a cache which is equal to input length.
func fibonacciCache(number int) int {
	if number <= 1 {
		return number
	}

	return cacheHelper(number, map[int]int{0: 0, 1: 1})
}

func cacheHelper(number int, cache map[int]int) int {
	if _, ok := cache[number]; !ok {
		cache[number] = cacheHelper(number-fibPrev, cache) + cacheHelper(number-fibNext, cache)
	}

	return cache[number]
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func fibonacciIterative(number int) int {
	if number <= 1 {
		return number
	}

	prev1, prev2 := 0, 1
	for i := 2; i <= number; i++ {
		prev1, prev2 = prev2, prev1+prev2
	}

	return prev2
}
