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
		name string
		root *TreeNode
		want int
	}{
		{
			name: "3-tier tree",
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
			want: 3,
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, maxDepthDFS(testCase.root))
			assert.Equal(t, testCase.want, maxDepthBFS(testCase.root))
			assert.Equal(t, testCase.want, maxDepthRecursive(testCase.root))
		})
	}
}

func maxDepthDFS(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return 1 + max(maxDepthDFS(root.Left), maxDepthDFS(root.Right))
}

func maxDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}

	queue := []*TreeNode{root}
	level := 0

	for len(queue) > 0 {
		queueLen := len(queue)
		for i := 0; i < queueLen; i++ {
			if queue[i].Left != nil {
				queue = append(queue, queue[i].Left)
			}

			if queue[i].Right != nil {
				queue = append(queue, queue[i].Right)
			}
		}

		queue = queue[queueLen:]
		level++
	}

	return level
}

func maxDepthRecursive(root *TreeNode) int {
	if root == nil {
		return 0
	}

	type stackRow struct {
		node  *TreeNode
		depth int
	}

	stack := []stackRow{
		{
			node:  root,
			depth: 1,
		},
	}

	var result int

	var row stackRow

	for len(stack) > 0 {
		row, stack = stack[0], stack[1:]

		if row.node != nil {
			result = max(result, row.depth)
			stack = append(
				stack,
				stackRow{node: row.node.Left, depth: row.depth + 1},
				stackRow{node: row.node.Right, depth: row.depth + 1},
			)
		}
	}

	return result
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}

	return n2
}
