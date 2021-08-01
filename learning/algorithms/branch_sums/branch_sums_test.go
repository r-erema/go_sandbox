package branchsums_test

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
		want []float32
	}{
		{
			name: "Case 0",
			bst:  tree(),
			want: []float32{7, 23, 49},
		},
		{
			name: "Case 1",
			bst:  tree2(),
			want: []float32{16, 24, 26, 68, 71},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			sums := branchSums(testCase.bst)
			assert.Equal(t, testCase.want, sums)
		})
	}
}

/*
Worst: O(n) time | O(n) space.
*/
func branchSums(bst *utils.BST) []float32 {
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
