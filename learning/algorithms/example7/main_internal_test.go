package example7

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
			array: append(make([]interface{}, 0),
				5,
				2,
				append(make([]interface{}, 0), 7, -1),
				3,
				append(make([]interface{}, 0),
					6,
					append(make([]interface{}, 0), -13, 8),
					4,
				),
			),
			want: 12,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, ProductSum(tt.array))
		})
	}
}
