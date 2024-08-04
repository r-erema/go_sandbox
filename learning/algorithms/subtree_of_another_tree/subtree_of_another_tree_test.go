package subtreeofanothertree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func TestIsSubtree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		root, subRoot *TreeNode
		want          bool
	}{
		{
			name: "There is subtree",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{
					Val: 4,
					Left: &TreeNode{
						Val: 1,
					},
					Right: &TreeNode{
						Val: 2,
					},
				},
				Right: &TreeNode{
					Val: 5,
				},
			},
			subRoot: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 1,
				},
				Right: &TreeNode{
					Val: 2,
				},
			},
			want: true,
		},
		{
			name: "There is no subtree",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{
					Val: 4,
					Left: &TreeNode{
						Val: 1,
					},
					Right: &TreeNode{
						Val: 2,
						Left: &TreeNode{
							Val: 0,
						},
					},
				},
				Right: &TreeNode{
					Val: 5,
				},
			},
			subRoot: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 1,
				},
				Right: &TreeNode{
					Val: 2,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isSubtree(tt.root, tt.subRoot))
		})
	}
}

// Time O(n*m), where n, m sizes of both trees
// Space O(1), we don't allocate additional memory.
func isSubtree(root, subRoot *TreeNode) bool {
	if subRoot == nil {
		return true
	}

	if root == nil {
		return false
	}

	if sameTree(root, subRoot) {
		return true
	}

	return isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
}

func sameTree(tree1, tree2 *TreeNode) bool {
	if tree1 == nil && tree2 == nil {
		return true
	}

	if tree1 != nil && tree2 != nil && tree1.Val == tree2.Val {
		return sameTree(tree1.Left, tree2.Left) && sameTree(tree1.Right, tree2.Right)
	}

	return false
}
