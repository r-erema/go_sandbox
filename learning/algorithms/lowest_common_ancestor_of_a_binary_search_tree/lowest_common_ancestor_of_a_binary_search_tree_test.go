package lowestcommonancestorofabinarysearchtree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

// Time O(log n), we don't have to visit every node in the tree, but only 1 node in a level
// Space O(1), we don't allocate additional memory.
func TestLowestCommonAncestor(t *testing.T) {
	t.Parallel()

	commonTree := &TreeNode{
		Val: 6,
		Left: &TreeNode{
			Val:  2,
			Left: &TreeNode{Val: 0},
			Right: &TreeNode{
				Val:   4,
				Left:  &TreeNode{Val: 3},
				Right: &TreeNode{Val: 5},
			},
		},
		Right: &TreeNode{
			Val:   8,
			Left:  &TreeNode{Val: 7},
			Right: &TreeNode{Val: 9},
		},
	}

	tests := []struct {
		name string
		p, q *TreeNode
		want int
	}{
		{
			name: "lca is root node",
			p:    &TreeNode{Val: 2},
			q:    &TreeNode{Val: 8},
			want: 6,
		},
		{
			name: "lca is node itself",
			p:    &TreeNode{Val: 2},
			q:    &TreeNode{Val: 4},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, lowestCommonAncestor(commonTree, tt.p, tt.q).Val)
		})
	}
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	for root != nil {
		switch {
		case p.Val > root.Val && q.Val > root.Val:
			root = root.Right
		case p.Val < root.Val && q.Val < root.Val:
			root = root.Left
		default:
			return root
		}
	}

	return nil
}
