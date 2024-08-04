package happynumber_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHappyNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int
		want  bool
	}{
		{
			name:  "number 19",
			input: 19,
			want:  true,
		},
		{
			name:  "number 19",
			input: 2,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isHappy(tt.input))
		})
	}
}

// Time O(log(n)) where n is the input number, as the sequence of transformations will be at most O(log n) before repeating or reaching 1
// Space O(n), since we use map.
func isHappy(n int) bool {
	nums := make(map[int]struct{})
	for n != 1 {
		if _, ok := nums[n]; ok {
			return false
		}

		nums[n] = struct{}{}

		sum := 0

		for n != 0 {
			digit := n % 10
			sum += digit * digit
			n /= 10
		}

		n = sum
	}

	return true
}
