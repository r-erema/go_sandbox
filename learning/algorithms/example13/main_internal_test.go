package example13

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertionSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "Case 0",
			str:  "redivider",
			want: true,
		},
		{
			name: "Case 1",
			str:  "abcdfcba",
			want: false,
		},
		{
			name: "Case 2",
			str:  "deified",
			want: true,
		},
		{
			name: "Case 3",
			str:  "level",
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, IsPalindrome(tt.str))
		})
	}
}
