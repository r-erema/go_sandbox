package example4

import (
	"github.com/r-erema/go_sendbox/utils"
)

/*
	Worst: O(n) time | O(n) space
*/
func BranchSums(bst *utils.BST) []float32 {
	return helper(bst, 0, []float32{})
}

func helper(node *utils.BST, sum float32, sums []float32) []float32 {
	sum += node.Value()

	if node.IsLeaf() {
		return append(sums, sum)
	}

	if node.Left() != nil {
		sums = helper(node.Left(), sum, sums)
	}

	if node.Right() != nil {
		sums = helper(node.Right(), sum, sums)
	}

	return sums
}
