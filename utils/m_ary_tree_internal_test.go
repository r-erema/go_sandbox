package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func tree() *MAryTree {
	return &MAryTree{
		value: "A",
		children: []*MAryTree{
			{value: "B", children: []*MAryTree{
				{value: "E"},
				{value: "F", children: []*MAryTree{
					{value: "I"},
					{value: "J"},
				}},
			}},
			{value: "C"},
			{value: "D", children: []*MAryTree{
				{value: "G", children: []*MAryTree{
					{value: "K"},
				}},
				{value: "H"},
			}},
		},
	}
}

func TestKAryTree_Traverse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tree *MAryTree
		want []string
	}{
		{
			name: "Case 0",
			tree: tree(),
			want: []string{"A", "B", "E", "F", "I", "J", "C", "D", "G", "K", "H"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			values := tt.tree.Traverse()
			assert.Equal(t, tt.want, values)
		})
	}
}
