package balancedbinarytree_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func TestBalancedTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{
			name: "Balanced tree",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{
					Val: 9,
				},
				Right: &TreeNode{
					Val: 20,
					Left: &TreeNode{
						Val: 15,
					},
					Right: &TreeNode{
						Val: 7,
					},
				},
			},
			want: true,
		},
		{
			name: "Non-balanced tree",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{
						Val: 3,
						Left: &TreeNode{
							Val: 4,
						},
						Right: &TreeNode{
							Val: 4,
						},
					},
					Right: &TreeNode{
						Val: 3,
					},
				},
				Right: &TreeNode{
					Val: 2,
				},
			},
			want: false,
		},
		{
			name: "Non-balanced tree, only 1 branch",
			root: &TreeNode{
				Val: 1,
				Right: &TreeNode{
					Val: 2,
					Right: &TreeNode{
						Val: 3,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isBalanced(tt.root))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(1), since we don't involve any additional data structure.
func isBalanced(root *TreeNode) bool {
	type res struct {
		isBalanced bool
		level      float64
	}

	var dfs func(root *TreeNode) res
	dfs = func(root *TreeNode) res {
		if root == nil {
			return res{true, 0}
		}

		left, right := dfs(root.Left), dfs(root.Right)
		balanced := left.isBalanced && right.isBalanced && math.Abs(left.level-right.level) <= 1

		return res{balanced, 1 + max(left.level, right.level)}
	}

	return dfs(root).isBalanced
}

func max(n1, n2 float64) float64 {
	if n1 > n2 {
		return n1
	}

	return n2
}
