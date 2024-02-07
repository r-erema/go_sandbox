package invertbinarytree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func TestInvertTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		root, want *TreeNode
	}{
		{
			name: "3-tier tree",
			root: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 1,
					},
					Right: &TreeNode{
						Val: 3,
					},
				},
				Right: &TreeNode{
					Val: 7,
					Left: &TreeNode{
						Val: 6,
					},
					Right: &TreeNode{
						Val: 9,
					},
				},
			},
			want: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 7,
					Left: &TreeNode{
						Val: 9,
					},
					Right: &TreeNode{
						Val: 6,
					},
				},
				Right: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 3,
					},
					Right: &TreeNode{
						Val: 1,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, invertTree(tt.root))
		})
	}
}

// Time O(n), we have to visit every node in the tree once to swap its left and right children
// Space O(n), due to the maximum depth of the recursion stack.
func invertTree(root *TreeNode) *TreeNode {
	if root != nil {
		root.Left, root.Right = invertTree(root.Right), invertTree(root.Left)
	}

	return root
}
