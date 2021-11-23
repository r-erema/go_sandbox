package example5

import (
	"github.com/r-erema/go_sendbox/utils"
)

/*
Worst: O(n) time | O(h) space, h - height if the binary tree.
*/
func NodeDepthRecursively(bst *utils.BST) int {
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
func NodeDepthIterative(bst *utils.BST) int {
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
func NodeDepthIterative2(bst *utils.BST) int {
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
