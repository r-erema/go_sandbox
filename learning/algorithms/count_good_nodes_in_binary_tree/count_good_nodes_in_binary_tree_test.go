package countgoodnodesinbinarytree_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func TestGoodNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *tree.Node
		want int
	}{
		{
			name: "4 good nodes",
			root: &tree.Node{
				Val: 3,
				Left: &tree.Node{
					Val:  1,
					Left: &tree.Node{Val: 3},
				},
				Right: &tree.Node{
					Val:   4,
					Left:  &tree.Node{Val: 1},
					Right: &tree.Node{Val: 5},
				},
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, goodNodes(tt.root))
		})
	}
}

// Time O(n), since we iterate input one time
// Space O(n), due to the recursion stack.
func goodNodes(root *tree.Node) int {
	var (
		good int
		dfs  func(root *tree.Node, maxVal int)
	)

	dfs = func(root *tree.Node, maxVal int) {
		if root == nil {
			return
		}

		if root.Val >= maxVal {
			good++
			maxVal = root.Val
		}

		dfs(root.Left, maxVal)
		dfs(root.Right, maxVal)
	}
	dfs(root, root.Val)

	return good
}
