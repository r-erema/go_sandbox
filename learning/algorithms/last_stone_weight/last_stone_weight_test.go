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

func lastStoneWeight(stones []int) int {
	heapifyMax(stones)

	for len(stones) > 1 {
		stone1 := popHeap(&stones)
		stone2 := popHeap(&stones)

		if stone1 == stone2 {
			continue
		}

		pushHeap(&stones, stone1-stone2)
	}

	if len(stones) == 1 {
		return stones[0]
	}

	return 0
}

func heapifyMax(arr []int) {
	for i := (len(arr) - 2) / 2; i >= 0; i-- {
		percolateDown(i, arr)
	}
}

func pushHeap(heap *[]int, val int) {
	*heap = append(*heap, val)
	percolateUp(len(*heap)-1, *heap)
}

func percolateUp(i int, heap []int) {
	for heap[i] > heap[(i-1)/2] {
		heap[(i-1)/2], heap[i] = heap[i], heap[(i-1)/2]
		i = (i - 1) / 2
	}
}

func popHeap(heap *[]int) int {
	res := (*heap)[0]
	(*heap)[0], *heap = (*heap)[len(*heap)-1], (*heap)[:len(*heap)-1]

	percolateDown(0, *heap)

	return res
}

func percolateDown(i int, heap []int) {
	for i*2+2 < len(heap) && (heap[i] < heap[i*2+1] || heap[i] < heap[i*2+2]) {
		leftChildGreaterThanRightChild := heap[i*2+1] > heap[i*2+2]

		if leftChildGreaterThanRightChild {
			heap[i], heap[i*2+1] = heap[i*2+1], heap[i]
			i = i*2 + 1

			continue
		}

		heap[i], heap[i*2+2] = heap[i*2+2], heap[i]
		i = i*2 + 2
	}

	if i*2+2 == len(heap) && heap[i] < heap[i*2+1] {
		heap[i], heap[i*2+1] = heap[i*2+1], heap[i]
	}
}
