package example4

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
)

func tree() *utils.BST {
	return utils.NewBST(5).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(10)).
		InsertRecursively(utils.NewBST(8)).
		InsertRecursively(utils.NewBST(34))
}

func tree2() *utils.BST {
	return utils.NewBST(9).
		InsertRecursively(utils.NewBST(4)).
		InsertRecursively(utils.NewBST(17)).
		InsertRecursively(utils.NewBST(3)).
		InsertRecursively(utils.NewBST(6)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(7)).
		InsertRecursively(utils.NewBST(20)).
		InsertRecursively(utils.NewBST(23))
}

func TestFindClosestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bst  *utils.BST
		want int
	}{
		{
			name: "Case 0",
			bst:  tree(),
			want: 6,
		},
		{
			name: "Case 1",
			bst:  tree2(),
			want: 20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			depth := NodeDepthRecursively(tt.bst)
			assert.Equal(t, tt.want, depth)
			depth = NodeDepthIterative(tt.bst)
			assert.Equal(t, tt.want, depth)
			depth = NodeDepthIterative2(tt.bst)
			assert.Equal(t, tt.want, depth)
		})
	}
}
