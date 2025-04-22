package heap_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func testHeap() *tree.Node {
	return &tree.Node{
		Val: 14,
		Left: &tree.Node{
			Val: 19,
			Left: &tree.Node{
				Val: 21,
				Left: &tree.Node{
					Val: 65,
				},
				Right: &tree.Node{
					Val: 30,
				},
			},
			Right: &tree.Node{
				Val: 26,
			},
		},
		Right: &tree.Node{
			Val: 16,
			Left: &tree.Node{
				Val: 19,
			},
			Right: &tree.Node{
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

	expectedHeap := &tree.Node{
		Val: 14,
		Left: &tree.Node{
			Val: 17,
			Left: &tree.Node{
				Val: 21,
				Left: &tree.Node{
					Val: 65,
				},
				Right: &tree.Node{
					Val: 30,
				},
			},
			Right: &tree.Node{
				Val: 19,
				Left: &tree.Node{
					Val: 26,
				},
			},
		},
		Right: &tree.Node{
			Val: 16,
			Left: &tree.Node{
				Val: 19,
			},
			Right: &tree.Node{
				Val: 68,
			},
		},
	}

	heap := testHeap()
	push(heap, 17)
	assert.Equal(t, expectedHeap, heap)
}

func TestPopHeap(t *testing.T) {
	t.Parallel()

	expectedHeap := &tree.Node{
		Val: 16,
		Left: &tree.Node{
			Val: 19,
			Left: &tree.Node{
				Val: 21,
				Left: &tree.Node{
					Val: 65,
				},
			},
			Right: &tree.Node{
				Val: 26,
			},
		},
		Right: &tree.Node{
			Val: 19,
			Left: &tree.Node{
				Val: 30,
			},
			Right: &tree.Node{
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

	sourceHeap := &tree.Node{
		Val: 50,
		Left: &tree.Node{
			Val: 80,
			Left: &tree.Node{
				Val: 30,
				Left: &tree.Node{
					Val: 90,
				},
				Right: &tree.Node{
					Val: 60,
				},
			},
			Right: &tree.Node{
				Val: 10,
			},
		},
		Right: &tree.Node{
			Val: 40,
			Left: &tree.Node{
				Val: 70,
			},
			Right: &tree.Node{
				Val: 20,
			},
		},
	}

	expectedHeap := &tree.Node{
		Val: 10,
		Left: &tree.Node{
			Val: 30,
			Left: &tree.Node{
				Val: 50,
				Left: &tree.Node{
					Val: 90,
				},
				Right: &tree.Node{
					Val: 60,
				},
			},
			Right: &tree.Node{
				Val: 80,
			},
		},
		Right: &tree.Node{
			Val: 20,
			Left: &tree.Node{
				Val: 70,
			},
			Right: &tree.Node{
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
func push(heap *tree.Node, val int) /**tree.Node */ {
	arrHeap := treeToArray(heap)

	arrHeap = append(arrHeap, val)
	i := len(arrHeap) - 1

	percolateUp(i, arrHeap)

	*heap = *arrayToTree(arrHeap)
}

// Time O(log(N)), the new root has to be swapped with its child on each level
// until it reaches the bottom level of the heap
// Space O(N), since we need an array containing elements from input heap.
func pop(heapTree *tree.Node) int {
	heap := treeToArray(heapTree)
	heapPtr := &heap

	res := heap[0]
	heap[0], *heapPtr = heap[len(heap)-1], heap[:len(heap)-1]

	percolateDown(0, heap)

	*heapTree = *arrayToTree(heap)

	return res
}

// Time O(N), since it iterates through the elements from the bottom level upwards,
// potentially shifting down elements multiple times to maintain the heap property
// Space O(N), since we need an array containing elements from input tree.
func heapify(tree *tree.Node) {
	arrTree := treeToArray(tree)

	for i := (len(arrTree) - 2) / 2; i >= 0; i-- {
		percolateDown(i, arrTree)
	}

	*tree = *arrayToTree(arrTree)
}

func percolateUp(i int, heap []int) {
	for heap[i] < heap[(i-1)/2] {
		heap[(i-1)/2], heap[i] = heap[i], heap[(i-1)/2]
		i = (i - 1) / 2
	}
}

func percolateDown(i int, heap []int) {
	if len(heap) == 2 && heap[0] > heap[1] {
		heap[0], heap[1] = heap[1], heap[0]

		return
	}

	for i*2+1 < len(heap) && (heap[i] > heap[i*2+1] || heap[i] > heap[i*2+2]) {
		leftChildGreaterRightChild := heap[i*2+1] > heap[i*2+2]

		if leftChildGreaterRightChild {
			heap[i], heap[i*2+2] = heap[i*2+2], heap[i]
			i = i*2 + 2

			continue
		}

		heap[i], heap[i*2+1] = heap[i*2+1], heap[i]
		i = i*2 + 1
	}
}

func treeToArray(heap *tree.Node) []int {
	var arr []int

	queue := []*tree.Node{heap}

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

func arrayToTree(arr []int) *tree.Node {
	var helper func(node *tree.Node, i int)
	helper = func(node *tree.Node, i int) {
		if i*2+1 < len(arr) {
			node.Left = &tree.Node{Val: arr[i*2+1]}
			helper(node.Left, i*2+1)
		}

		if i*2+2 < len(arr) {
			node.Right = &tree.Node{Val: arr[i*2+2]}
			helper(node.Right, i*2+2)
		}
	}

	heap := &tree.Node{Val: arr[0]}
	helper(heap, 0)

	return heap
}
