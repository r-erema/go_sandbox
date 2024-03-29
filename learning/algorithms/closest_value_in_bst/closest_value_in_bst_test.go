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

/*
Average: O(log(n)) time | O(log(n)) space.
Worst: O(n) time | O(log(n)) space.
*/
func findClosestValueRecursively(bst *utils.BST, target float32) float32 {
	closestValue := bst.Value()

	return helperForRecursion(bst, target, closestValue)
}

func helperForRecursion(bst *utils.BST, target, closestValue float32) float32 {
	if bst == nil {
		return closestValue
	}

	if bst.Value() == target {
		return target
	}

	newDelta := math.Abs(float64(bst.Value() - target))
	oldDelta := math.Abs(float64(closestValue - target))

	if newDelta < oldDelta {
		closestValue = bst.Value()
	}

	if target > bst.Value() {
		closestValue = helperForRecursion(bst.Right(), target, closestValue)
	}

	if target < bst.Value() {
		closestValue = helperForRecursion(bst.Left(), target, closestValue)
	}

	return closestValue
}

/*
Average: O(log(n)) time | O(1) space.
Worst: O(n) time | O(1) space.
*/
func findClosestValueIteratively(bst *utils.BST, target float32) float32 {
	closestValue := bst.Value()

	return helperForIteration(bst, target, closestValue)
}

func helperForIteration(bst *utils.BST, target, closestValue float32) float32 {
	currentNode := bst
	for currentNode != nil {
		newDelta := math.Abs(float64(currentNode.Value() - target))
		oldDelta := math.Abs(float64(closestValue - target))

		if newDelta < oldDelta {
			closestValue = currentNode.Value()
		}

		if target > currentNode.Value() {
			currentNode = currentNode.Right()
		} else if target < currentNode.Value() {
			currentNode = currentNode.Left()
		} else {
			break
		}
	}

	return closestValue
}
