package diameterofbinarytree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func TestDiameterOfBinaryTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *TreeNode
		want int
	}{
		{
			name: "Diameter 3",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 4,
					},
					Right: &TreeNode{
						Val: 5,
					},
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			want: 3,
		},
		{
			name: "Diameter 2",
			root: &TreeNode{
				Val: 2,
				Left: &TreeNode{
					Val: 3,
					Left: &TreeNode{
						Val: 1,
					},
					Right: nil,
				},
				Right: nil,
			},
			want: 2,
		},
		{
			name: "Diameter 1",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, diameterOfBinaryTree(tt.root))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(N), due to the maximum depth of the recursion stack.
func diameterOfBinaryTree(root *TreeNode) int {
	var diameter int

	var dfs func(root *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}

		left, right := dfs(node.Left), dfs(node.Right)
		diameter = max(diameter, left+right)

		return 1 + max(left, right)
	}
	dfs(root)

	return diameter
}
