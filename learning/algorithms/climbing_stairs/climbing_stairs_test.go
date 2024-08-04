package climbingstairs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClimbingStairs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n, want int
	}{
		{
			name: "2 steps",
			n:    2,
			want: 2,
		},
		{
			name: "3 steps",
			n:    3,
			want: 3,
		},
		{
			name: "4 steps",
			n:    4,
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, climbStairs(tt.n))
		})
	}
}

// Time O(n), since we walk trough the n 1 time
// Space O(1), we don't use any extra space.
func climbStairs(n int) int {
	one, two := 1, 1
	for range n - 1 {
		one, two = one+two, one
	}

	return one
}
