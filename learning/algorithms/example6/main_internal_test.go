package example6

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, FibonacciRecursive(tt.n))
			assert.Equal(t, tt.want, FibonacciCache(tt.n))
			assert.Equal(t, tt.want, FibonacciIterative(tt.n))
		})
	}
}
