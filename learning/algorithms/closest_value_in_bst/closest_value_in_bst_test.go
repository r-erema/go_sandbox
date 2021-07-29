package closestvalueinbst_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         float64
	Left, Right *TreeNode
}

func TestFindClosestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bst    *TreeNode
		target float64
		want   float64
	}{
		{
			name: "case 0",
			bst: &TreeNode{
				Val:  5,
				Left: &TreeNode{Val: 2},
				Right: &TreeNode{
					Val:   10,
					Left:  &TreeNode{Val: 8},
					Right: &TreeNode{Val: 34},
				},
			},
			target: 20,
			want:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			closest := closestValueInBST(tt.bst, tt.target)
			assert.InEpsilon(t, tt.want, closest, 0)
		})
	}
}

// Time O(log(N)), since we explore a particular subtree
// Space O(1), since we don't involve any additional data structure.
func closestValueInBST(root *TreeNode, target float64) float64 {
	closest := root.Val

	for root != nil {
		if math.Abs(root.Val-target) < math.Abs(closest-target) {
			closest = root.Val
		}

		if target < root.Val {
			root = root.Left
		} else if target > root.Val {
			root = root.Right
		} else {
			break
		}
	}

	return closest
}
