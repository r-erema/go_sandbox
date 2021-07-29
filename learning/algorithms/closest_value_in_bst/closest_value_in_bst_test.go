package closestvalueinbst_test

import (
	"math"
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
)

func bst1() *utils.BST {
	return utils.NewBST(5).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(10)).
		InsertRecursively(utils.NewBST(8)).
		InsertRecursively(utils.NewBST(34))
}

func bst2() *utils.BST {
	return utils.NewBST(9).
		InsertRecursively(utils.NewBST(4)).
		InsertRecursively(utils.NewBST(17)).
		InsertRecursively(utils.NewBST(3)).
		InsertRecursively(utils.NewBST(6)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(7)).
		InsertRecursively(utils.NewBST(20))
}

func bst3() *utils.BST {
	return utils.NewBST(10).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(15)).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(13)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(14))
}

func TestFindClosestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bst    *utils.BST
		target float32
		want   float32
	}{
		{
			name:   "case 0",
			bst:    bst1(),
			target: 20,
			want:   10,
		},
		{
			name:   "case 1",
			bst:    bst2(),
			target: 4,
			want:   4,
		},
		{
			name:   "case 2",
			bst:    bst2(),
			target: 18,
			want:   17,
		},
		{
			name:   "case 3",
			bst:    bst2(),
			target: 12,
			want:   9,
		},
		{
			name:   "case 4",
			bst:    bst3(),
			target: 12,
			want:   13,
		},
		{
			name:   "case 5",
			bst:    bst3(),
			target: 3.5,
			want:   5,
		},
		{
			name:   "case 6",
			bst:    bst3(),
			target: 22,
			want:   22,
		},
		{
			name:   "case 7",
			bst:    bst3(),
			target: 13.6,
			want:   14,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			closest := findClosestValueRecursively(testCase.bst, testCase.target)
			assert.Equal(t, testCase.want, closest)
			closest = findClosestValueIteratively(testCase.bst, testCase.target)
			assert.Equal(t, testCase.want, closest)
		})
	}
}

func findClosestValueRecursively(bst *utils.BST, target float32) float32 {
	closest := bst.Value()
	currentDelta := float32(math.Abs(float64(target - closest)))

	var walkBST func(node *utils.BST)
	walkBST = func(node *utils.BST) {
		if node == nil {
			return
		}

		if float32(math.Abs(float64(target-node.Value()))) < currentDelta {
			closest = node.Value()
			currentDelta = float32(math.Abs(float64(target - closest)))
		}

		walkBST(node.Left())
		walkBST(node.Right())
	}

	walkBST(bst)

	return closest
}

func findClosestValueIteratively(bst *utils.BST, target float32) float32 {
	closest := bst.Value()
	queue := []*utils.BST{bst}

	currentDelta := float32(math.Abs(float64(target - closest)))

	for len(queue) > 0 {
		bst, queue = queue[0], queue[1:]

		if float32(math.Abs(float64(target-bst.Value()))) < currentDelta {
			closest = bst.Value()
			currentDelta = float32(math.Abs(float64(target - closest)))
		}

		if bst.Left() != nil {
			queue = append(queue, bst.Left())
		}

		if bst.Right() != nil {
			queue = append(queue, bst.Right())
		}
	}

	return closest
}
