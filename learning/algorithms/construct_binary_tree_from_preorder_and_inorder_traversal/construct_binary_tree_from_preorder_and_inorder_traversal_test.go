package constructbinarytreefrompreorderandinordertraversal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func TestBuildTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		preorder, inorder []int
		want              *TreeNode
	}{
		{
			name:     "Normal tree",
			preorder: []int{3, 9, 20, 15, 7},
			inorder:  []int{9, 3, 15, 20, 7},
			want: &TreeNode{
				Val:  3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
		},
		{
			name:     "Normal tree 2",
			preorder: []int{2, 8, 11, 7, 6, 12, 3, 5},
			inorder:  []int{11, 8, 6, 7, 2, 3, 5, 12},
			want: &TreeNode{
				Val: 2,
				Left: &TreeNode{
					Val:  8,
					Left: &TreeNode{Val: 11},
					Right: &TreeNode{
						Val:  7,
						Left: &TreeNode{Val: 6},
					},
				},
				Right: &TreeNode{
					Val: 12,
					Left: &TreeNode{
						Val:   3,
						Right: &TreeNode{Val: 5},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, buildTree(tt.preorder, tt.inorder))
		})
	}
}

// Time O(n), since each node in the tree is processed exactly once
// Space O(n), for the following reasons:
//   - hashmap d stores N key-value pairs representing each node's value and index in the inorder traversal
//   - the recursion stack may grow up to O(n) in the case of a skewed tree (where each node has only one child)
//   - the output structure, which is a binary tree that contains N TreeNode instances.
func buildTree(preorder, inorder []int) *TreeNode {
	inorderMap := make(map[int]int)
	for i := range inorder {
		inorderMap[inorder[i]] = i
	}

	var helper func(preorderStart, inorderStart, size int) *TreeNode
	helper = func(preorderStart, inorderStart, size int) *TreeNode {
		if size == 0 {
			return nil
		}

		root := &TreeNode{Val: preorder[preorderStart]}
		mid := inorderMap[root.Val]
		leftSubtreeSize := mid - inorderStart

		root.Left = helper(preorderStart+1, inorderStart, leftSubtreeSize)
		root.Right = helper(preorderStart+leftSubtreeSize+1, mid+1, size-1-leftSubtreeSize)

		return root
	}

	return helper(0, 0, len(inorder))
}
