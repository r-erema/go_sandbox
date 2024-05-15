package laststoneweight_test

import (
	"container/heap"
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
			name:   "simple stones",
			stones: []int{2, 7, 4, 1, 8, 1},
			want:   1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, lastStoneWeight(tt.stones))
		})
	}
}

type StonesHeap []int

func (sh *StonesHeap) Len() int           { return len(*sh) }
func (sh *StonesHeap) Less(i, j int) bool { return (*sh)[i] > (*sh)[j] }
func (sh *StonesHeap) Swap(i, j int)      { (*sh)[i], (*sh)[j] = (*sh)[j], (*sh)[i] }
func (sh *StonesHeap) Push(x any) {
	if xInt, ok := x.(int); ok {
		*sh = append(*sh, xInt)
	}
}

func (sh *StonesHeap) Pop() any {
	var popped any
	popped, *sh = (*sh)[len(*sh)-1], (*sh)[:len(*sh)-1]

	return popped
}

// Time O(N * logN), since we push and pop elements on each iteration
// Space O(N), since we create a heap from an input.
func lastStoneWeight(stones []int) int {
	stonesHeap := StonesHeap(stones)

	heap.Init(&stonesHeap)

	var first, second int

	for len(stonesHeap) > 1 {
		if popped, ok := heap.Pop(&stonesHeap).(int); ok {
			first = popped
		}

		if popped, ok := heap.Pop(&stonesHeap).(int); ok {
			second = popped
		}

		heap.Push(&stonesHeap, first-second)
	}

	return stonesHeap[0]
}
