package nodedepth_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
)

func tree() *utils.BST {
	return utils.NewBST(5).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(10)).
		InsertRecursively(utils.NewBST(8)).
		InsertRecursively(utils.NewBST(34))
}

func tree2() *utils.BST {
	return utils.NewBST(9).
		InsertRecursively(utils.NewBST(4)).
		InsertRecursively(utils.NewBST(17)).
		InsertRecursively(utils.NewBST(3)).
		InsertRecursively(utils.NewBST(6)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(7)).
		InsertRecursively(utils.NewBST(20)).
		InsertRecursively(utils.NewBST(23))
}

func TestFindClosestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bst  *utils.BST
		want int
	}{
		{
			name: "Case 0",
			bst:  tree(),
			want: 6,
		},
		{
			name: "Case 1",
			bst:  tree2(),
			want: 20,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			depth := nodeDepthRecursively(testCase.bst)
			assert.Equal(t, testCase.want, depth)
			depth = nodeDepthIterative(testCase.bst)
			assert.Equal(t, testCase.want, depth)
			depth = nodeDepthIterative2(testCase.bst)
			assert.Equal(t, testCase.want, depth)
		})
	}
}

/*
Worst: O(n) time | O(h) space, h - height if the binary tree.
*/
func nodeDepthRecursively(bst *utils.BST) int {
	return recursionHelper(bst, 0)
}

func recursionHelper(node *utils.BST, currentDepth int) int {
	var depthLeft, depthRight int

	if node.Left() != nil {
		depthLeft = recursionHelper(node.Left(), currentDepth+1)
	}

	if node.Right() != nil {
		depthRight = recursionHelper(node.Right(), currentDepth+1)
	}

	return currentDepth + depthLeft + depthRight
}

/*
Worst: O(n) time | O(h) space, h - height if the binary tree.
*/
func nodeDepthIterative(bst *utils.BST) int {
	var depth int

	depthNodes := [][]*utils.BST{{bst}}

	for level := 0; len(depthNodes) > level; level++ {
		var (
			levelNodes      []*utils.BST
			levelNodesCount int
		)

		for _, node := range depthNodes[level] {
			if left := node.Left(); left != nil {
				levelNodes = append(levelNodes, left)
				levelNodesCount++
			}

			if right := node.Right(); right != nil {
				levelNodes = append(levelNodes, right)
				levelNodesCount++
			}
		}

		if len(levelNodes) > 0 {
			depthNodes = append(depthNodes, levelNodes)
			depth += (level + 1) * levelNodesCount
		}
	}

	return depth
}

/*
Worst: O(n) time | O(h) space, h - height if the binary tree.
*/
func nodeDepthIterative2(bst *utils.BST) int {
	type stackItem struct {
		node  *utils.BST
		depth int
	}

	var totalDepth int

	stack := []stackItem{{node: bst, depth: 0}}

	for len(stack) > 0 {
		item := stack[0]
		stack = stack[1:]
		node, depth := item.node, item.depth

		if node == nil {
			continue
		}

		stack = append(
			stack,
			stackItem{node: node.Left(), depth: depth + 1},
			stackItem{node: node.Right(), depth: depth + 1},
		)

		totalDepth += depth
	}

	return totalDepth
}
