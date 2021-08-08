package productsum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		array []interface{}
		want  int
	}{
		{
			name: "Case 0",
			array: []interface{}{
				5,
				2,
				[]interface{}{7, -1},
				3,
				[]interface{}{
					6,
					[]interface{}{-13, 8},
					4,
				},
			},
			want: 12,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, productSum(tt.array))
		})
	}
}

// Time O(n), since we walk trough the n 1 time
// Space O(1), we don't use any extra space.
func productSum(array []interface{}) int {
	return helper(array, 1)
}

func helper(array []interface{}, depth int) int {
	var sum int

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
