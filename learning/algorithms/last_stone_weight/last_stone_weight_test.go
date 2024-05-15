package laststoneweight_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastStoneWeight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		stones []int
		want   int
	}{
		{
			name:   "6 stones",
			stones: []int{2, 7, 4, 1, 8, 1},
			want:   1,
		},
		{
			name:   "2 stones ascending ordering",
			stones: []int{1, 3},
			want:   2,
		},
		{
			name:   "2 stones descending ordering",
			stones: []int{3, 1},
			want:   2,
		},
		{
			name:   "4 stones",
			stones: []int{9, 3, 2, 10},
			want:   0,
		},
		{
			name:   "8 stones",
			stones: []int{10, 5, 4, 10, 3, 1, 7, 8},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, lastStoneWeight(tt.stones))
		})
	}
}

// Time O(N*logN), since we need to pop each element from the hea
// Space O(n), since the heap size equals the input.
func lastStoneWeight(stones []int) int {
	heap := heapify(stones)

	for len(heap) > 1 {
		stone1, stone2 := pop(&heap), pop(&heap)
		newStone := stone1 - stone2

		if newStone > 0 {
			push(newStone, &heap)
		}
	}

	if len(heap) == 1 {
		return heap[0]
	}

	return 0
}

func heapify(stones []int) []int {
	var heap []int

	for i := range stones {
		push(stones[i], &heap)
	}

	return heap
}

func push(node int, heap *[]int) {
	*heap = append(*heap, node)
	percolateUp(*heap)
}

func pop(heap *[]int) int {
	popped := (*heap)[0]
	(*heap)[0] = (*heap)[len(*heap)-1]
	*heap = (*heap)[:len(*heap)-1]
	percolateDown(*heap)

	return popped
}

func percolateUp(heap []int) {
	child1 := len(heap) - 1
	parent := (child1 - 1) / 2
	child2 := parent*2 + 1

	for heap[parent] < heap[child1] {
		if heap[child2] > heap[child1] {
			heap[child2], heap[child1] = heap[child1], heap[child2]
		} else {
			heap[parent], heap[child1] = heap[child1], heap[parent]
		}

		child1 = parent
		parent = (child1 - 1) / 2
		child2 = parent*2 + 1
	}
}

func percolateDown(heap []int) {
	parent, child1, child2 := 0, 1, 2

	for (child1 < len(heap) && heap[parent] < heap[child1]) || (child2 < len(heap) && heap[parent] < heap[child2]) {
		if child2 < len(heap) && heap[child1] < heap[child2] {
			heap[parent], heap[child2], parent = heap[child2], heap[parent], child2
		} else {
			heap[parent], heap[child1], parent = heap[child1], heap[parent], child1
		}

		child1, child2 = parent*2+1, (parent+1)*2
	}
}
