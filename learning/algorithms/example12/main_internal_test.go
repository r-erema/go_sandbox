package example12

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertionSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		array []int
		want  []int
	}{
		{
			name:  "Case 0",
			array: []int{7, 0, 4},
			want:  []int{0, 4, 7},
		},
		{
			name:  "Case 1",
			array: []int{8, 5, 2, 9, 5, 6, 3},
			want:  []int{2, 3, 5, 5, 6, 8, 9},
		},
		{
			name:  "Case 2",
			array: []int{1, -1},
			want:  []int{-1, 1},
		},
		{
			name:  "Case 3",
			array: []int{7},
			want:  []int{7},
		},
		{
			name:  "Case 4",
			array: []int{8, 6, 4, 2, 0, -1, -2, -3, -4, -5},
			want:  []int{-5, -4, -3, -2, -1, 0, 2, 4, 6, 8},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			SelectionSort(&tt.array)
			assert.Equal(t, tt.want, tt.array)
		})
	}
}
