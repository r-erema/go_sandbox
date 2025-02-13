package binarytreerightsideview_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func TestInvertTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *tree.Node
		want []int
	}{
		{
			name: "see 3 nodes",
			root: &tree.Node{
				Val: 1,
				Left: &tree.Node{
					Val:   2,
					Right: &tree.Node{Val: 5},
				},
				Right: &tree.Node{
					Val:   3,
					Right: &tree.Node{Val: 4},
				},
			},
			want: []int{1, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, rightSideView(tt.root))
		})
	}
}

// Time O(n), since we iterate input one time
// Space O(n), due to the recursion stack.
func rightSideView(root *tree.Node) []int {
	var res []int

	var preorderTraversal func(root *tree.Node, level int)
	preorderTraversal = func(root *tree.Node, level int) {
		if root == nil {
			return
		}

		if len(res) == level {
			res = append(res, root.Val)
		}

		preorderTraversal(root.Right, level+1)
		preorderTraversal(root.Left, level+1)
	}

	preorderTraversal(root, 0)

	return res
}
