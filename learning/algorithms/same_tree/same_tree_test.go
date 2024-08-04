package sametree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func TestSameTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tree1, tree2 *TreeNode
		want         bool
	}{
		{
			name: "same tree",
			tree1: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			tree2: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			want: true,
		},
		{
			name: "not same tree with 2 nodes",
			tree1: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
				},
			},
			tree2: &TreeNode{
				Val: 1,
				Right: &TreeNode{
					Val: 2,
				},
			},
			want: false,
		},
		{
			name: "not same tree with 3 nodes",
			tree1: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
				},
				Right: &TreeNode{
					Val: 1,
				},
			},
			tree2: &TreeNode{
				Val: 1,
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

			assert.Equal(t, tt.want, isSameTree(tt.tree1, tt.tree2))
		})
	}
}

// Time O(p+q), since we should iterate recursively through the both trees
// Space O(1), we don't allocate additional memory.
func isSameTree(tree1, tree2 *TreeNode) bool {
	if tree1 == nil && tree2 == nil {
		return true
	}

	if tree1 == nil || tree2 == nil || tree1.Val != tree2.Val {
		return false
	}

	return isSameTree(tree1.Left, tree2.Left) && isSameTree(tree1.Right, tree2.Right)
}
