package heap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func testHeap() *TreeNode {
	return &TreeNode{
		Val: 14,
		Left: &TreeNode{
			Val: 19,
			Left: &TreeNode{
				Val: 21,
				Left: &TreeNode{
					Val: 65,
				},
				Right: &TreeNode{
					Val: 30,
				},
			},
			Right: &TreeNode{
				Val: 26,
			},
		},
		Right: &TreeNode{
			Val: 16,
			Left: &TreeNode{
				Val: 19,
			},
			Right: &TreeNode{
				Val: 68,
			},
		},
	}
}

func TestHeapToArrayAndArrayToHeap(t *testing.T) {
	t.Parallel()

	arr := treeToArray(testHeap())
	assert.Equal(t, []int{14, 19, 16, 21, 26, 19, 68, 65, 30}, arr)

	heap := arrayToTree(arr)
	assert.Equal(t, testHeap(), heap)
}

func TestPushToHeap(t *testing.T) {
	t.Parallel()

	expectedHeap := &TreeNode{
		Val: 14,
		Left: &TreeNode{
			Val: 17,
			Left: &TreeNode{
				Val: 21,
				Left: &TreeNode{
					Val: 65,
				},
				Right: &TreeNode{
					Val: 30,
				},
			},
			Right: &TreeNode{
				Val: 19,
				Left: &TreeNode{
					Val: 26,
				},
			},
		},
		Right: &TreeNode{
			Val: 16,
			Left: &TreeNode{
				Val: 19,
			},
			Right: &TreeNode{
				Val: 68,
			},
		},
	}

	resultHeap := push(testHeap(), 17)
	assert.Equal(t, expectedHeap, resultHeap)
}

func TestPopHeap(t *testing.T) {
	t.Parallel()

	expectedHeap := &TreeNode{
		Val: 16,
		Left: &TreeNode{
			Val: 19,
			Left: &TreeNode{
				Val: 21,
				Left: &TreeNode{
					Val: 65,
				},
			},
			Right: &TreeNode{
				Val: 26,
			},
		},
		Right: &TreeNode{
			Val: 19,
			Left: &TreeNode{
				Val: 30,
			},
			Right: &TreeNode{
				Val: 68,
			},
		},
	}

	clonedTestHeap := *testHeap()
	popped := pop(&clonedTestHeap)
	assert.Equal(t, 14, popped)
	assert.Equal(t, expectedHeap, &clonedTestHeap)
}

func TestHeapify(t *testing.T) {
	t.Parallel()

	sourceHeap := &TreeNode{
		Val: 50,
		Left: &TreeNode{
			Val: 80,
			Left: &TreeNode{
				Val: 30,
				Left: &TreeNode{
					Val: 90,
				},
				Right: &TreeNode{
					Val: 60,
				},
			},
			Right: &TreeNode{
				Val: 10,
			},
		},
		Right: &TreeNode{
			Val: 40,
			Left: &TreeNode{
				Val: 70,
			},
			Right: &TreeNode{
				Val: 20,
			},
		},
	}

	expectedHeap := &TreeNode{
		Val: 10,
		Left: &TreeNode{
			Val: 30,
			Left: &TreeNode{
				Val: 50,
				Left: &TreeNode{
					Val: 90,
				},
				Right: &TreeNode{
					Val: 60,
				},
			},
			Right: &TreeNode{
				Val: 80,
			},
		},
		Right: &TreeNode{
			Val: 20,
			Left: &TreeNode{
				Val: 70,
			},
			Right: &TreeNode{
				Val: 40,
			},
		},
	}

	heapify(sourceHeap)
	assert.Equal(t, expectedHeap, sourceHeap)
}

// Time O(log(N)), the number of operations required depends only on the number of levels
// the new element must rise to satisfy the heap property
// Space O(N), since we need an array containing elements from input heap.
func push(heap *TreeNode, val int) *TreeNode {
	arrHeap := leadingZeroArray(treeToArray(heap))

	arrHeap = append(arrHeap, val)
	i := len(arrHeap) - 1

	for arrHeap[i] < arrHeap[i/2] {
		arrHeap[i], arrHeap[i/2] = arrHeap[i/2], arrHeap[i]
		i /= 2
	}

	return arrayToTree(arrHeap[1:])
}

// Time O(log(N)), the new root has to be swapped with its child on each level until it reaches the bottom level of the heap
// Space O(N), since we need an array containing elements from input heap.
func pop(heap *TreeNode) int {
	arrHeap := leadingZeroArray(treeToArray(heap))

	res := arrHeap[1]
	arrHeap[1], arrHeap = arrHeap[len(arrHeap)-1], arrHeap[:len(arrHeap)-1]
	i := 1

	percolateDown(i, arrHeap)

	*heap = *arrayToTree(arrHeap[1:])

	return res
}

// Time O(N), since it iterates through the elements from the bottom level upwards,
// potentially sifting down elements multiple times to maintain the heap property
// Space O(N), since we need an array containing elements from input tree.
func heapify(tree *TreeNode) {
	arrTree := leadingZeroArray(treeToArray(tree))

	cur := (len(arrTree) - 1) / 2

	for cur > 0 {
		i := cur
		percolateDown(i, arrTree)

		cur--
	}

	*tree = *arrayToTree(arrTree[1:])
}

func percolateDown(i int, arr []int) {
	leftChildExists := func() bool {
		return 2*i < len(arr)
	}

	rightChildExists := func() bool {
		return 2*i+1 < len(arr)
	}

	rightChildLessThanLeftChild := func() bool {
		return arr[2*i+1] < arr[2*i]
	}

	parentGreaterThanRightChild := func() bool {
		return arr[i] > arr[2*i+1]
	}

	parentGreaterThanLeftChild := func() bool {
		return arr[i] > arr[2*i]
	}

	for leftChildExists() {
		if rightChildExists() && rightChildLessThanLeftChild() && parentGreaterThanRightChild() {
			arr[i], arr[2*i+1] = arr[2*i+1], arr[i]
			i = 2*i + 1
		} else if parentGreaterThanLeftChild() {
			arr[i], arr[2*i] = arr[2*i], arr[i]
			i *= 2
		} else {
			break
		}
	}
}

func treeToArray(heap *TreeNode) []int {
	var arr []int

	queue := []*TreeNode{heap}

	for len(queue) > 0 {
		heap, queue = queue[0], queue[1:]
		arr = append(arr, heap.Val)

		if heap.Left != nil {
			queue = append(queue, heap.Left)
		}

		if heap.Right != nil {
			queue = append(queue, heap.Right)
		}
	}

	return arr
}

func arrayToTree(arr []int) *TreeNode {
	arr = leadingZeroArray(arr)

	var helper func(node *TreeNode, i int)
	helper = func(node *TreeNode, i int) {
		if i*2 < len(arr) {
			node.Left = &TreeNode{Val: arr[i*2]}
			helper(node.Left, i*2)
		}

		if i*2+1 < len(arr) {
			node.Right = &TreeNode{Val: arr[i*2+1]}
			helper(node.Right, i*2+1)
		}
	}

	heap := &TreeNode{Val: arr[1]}
	helper(heap, 1)

	return heap
}

func leadingZeroArray(arr []int) []int {
	updatedArr := make([]int, len(arr)+1)
	for i := range arr {
		updatedArr[i+1] = arr[i]
	}

	return updatedArr
}
