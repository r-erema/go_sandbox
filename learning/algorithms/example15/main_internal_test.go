package example15

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertionSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want string
	}{
		{
			name: "Case 0",
			arr:  []int{1, 4, 5, 2, 3, 9, 8, 11, 0},
			want: "0-5,8-9,11",
		},
		{
			name: "Case 1",
			arr:  []int{1, 3, 2},
			want: "1-3",
		},
		{
			name: "Case 2",
			arr:  []int{1, 4},
			want: "1,4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, CollapseArrayToRange(tt.arr))
		})
	}
}
