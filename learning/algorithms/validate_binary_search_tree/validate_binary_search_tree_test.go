package validatebinarysearchtree_test

import (
	"math"
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func TestInvertTree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *tree.Node
		want bool
	}{
		{
			name: "valid tree",
			root: &tree.Node{
				Val: 2,
				Left: &tree.Node{
					Val: 1,
				},
				Right: &tree.Node{
					Val: 3,
				},
			},
			want: true,
		},
		{
			name: "invalid tree",
			root: &tree.Node{
				Val: 5,
				Left: &tree.Node{
					Val: 1,
				},
				Right: &tree.Node{
					Val: 4,
					Left: &tree.Node{
						Val: 3,
					},
					Right: &tree.Node{
						Val: 6,
					},
				},
			},
			want: false,
		},
		{
			name: "invalid tree",
			root: &tree.Node{
				Val: 1,
				Left: &tree.Node{
					Val: 1,
				},
			},
			want: false,
		},
		{
			name: "invalid tree",
			root: &tree.Node{
				Val: 5,
				Left: &tree.Node{
					Val: 4,
				},
				Right: &tree.Node{
					Val:   6,
					Left:  &tree.Node{Val: 3},
					Right: &tree.Node{Val: 7},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isValidBST(tt.root))
		})
	}
}

// Time O(n), since we iterate input one time
// Space O(n), due to the recursion stack, for a balanced tree, h is log(n), making the space complexity O(log(n)),
// in the worst case of a skewed tree, h is n, making the space complexity O(n).
func isValidBST(root *tree.Node) bool {
	var dfs func(root *tree.Node, minLimit, maxLimit int) bool
	dfs = func(root *tree.Node, minLimit, maxLimit int) bool {
		if root == nil {
			return true
		}

		if root.Val <= minLimit || root.Val >= maxLimit {
			return false
		}

		return dfs(root.Left, minLimit, root.Val) && dfs(root.Right, root.Val, maxLimit)
	}

	return dfs(root, math.MinInt64, math.MaxInt64)
}
