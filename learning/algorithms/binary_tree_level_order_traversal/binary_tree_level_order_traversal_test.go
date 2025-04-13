package binary_tree_level_order_traversal_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func TestBinaryTreeLevelOrderTraversal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root *tree.Node
		want [][]int
	}{
		{
			name: "normal case",
			root: &tree.Node{
				Val:  3,
				Left: &tree.Node{Val: 9},
				Right: &tree.Node{
					Val:   20,
					Left:  &tree.Node{Val: 15},
					Right: &tree.Node{Val: 7},
				},
			},
			want: [][]int{
				{3},
				{9, 20},
				{15, 7},
			},
		},
		{
			name: "1 node",
			root: &tree.Node{
				Val: 1,
			},
			want: [][]int{
				{1},
			},
		},
		{
			name: "no nodes",
			root: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, levelOrder(tt.root))
		})
	}
}

// Time O(n), since we iterate input one time
// Space O(n), since we need an array containing elements from input tree.
func levelOrder(root *tree.Node) [][]int {
	if root == nil {
		return nil
	}

	var (
		res        [][]int
		levelQueue []*tree.Node
	)

	queue := [][]*tree.Node{{root}}

	for len(queue) > 0 {
		levelQueue, queue = queue[0], queue[1:]

		var (
			levelVals     []int
			newLevelQueue []*tree.Node
			node          *tree.Node
		)

		for len(levelQueue) > 0 {
			node, levelQueue = levelQueue[0], levelQueue[1:]
			levelVals = append(levelVals, node.Val)

			if node.Left != nil {
				newLevelQueue = append(newLevelQueue, node.Left)
			}

			if node.Right != nil {
				newLevelQueue = append(newLevelQueue, node.Right)
			}
		}

		res = append(res, levelVals)

		if len(newLevelQueue) > 0 {
			queue = append(queue, newLevelQueue)
		}
	}

	return res
}
