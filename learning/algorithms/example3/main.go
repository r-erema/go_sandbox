package example3

import (
	"math"

	"github.com/r-erema/go_sendbox/utils"
)

/*
	Average: O(log(n)) time | O(log(n)) space
	Worst: O(n) time | O(log(n)) space
*/
func FindClosestValueRecursively(bst *utils.BST, target float32) float32 {
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
	Average: O(log(n)) time | O(1) space
	Worst: O(n) time | O(1) space
*/
func FindClosestValueIteratively(bst *utils.BST, target float32) float32 {
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
