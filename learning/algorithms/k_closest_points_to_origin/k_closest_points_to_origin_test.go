package k_closest_points_to_origin_test

import (
	"iter"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKClosest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		points [][]int
		k      int
		want   [][]int
	}{
		{
			name:   "3 points",
			points: [][]int{{0, 2}, {2, 0}, {2, 2}},
			k:      2,
			want:   [][]int{{0, 2}, {2, 0}},
		},
		{
			name:   "3 points",
			points: [][]int{{3, 3}, {5, -1}, {-2, 4}},
			k:      2,
			want:   [][]int{{3, 3}, {-2, 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, kClosest(tt.points, tt.k))
		})
	}
}

type point struct {
	coordinates        [2]int
	distanceFromOrigin float64
}

// Time O(N*log k), since we pop from and push to the heap which size is k
// Space O(k), as we allocate additional space of size k.
func kClosest(rawCoordinates [][]int, k int) [][]int {
	points := heapify(rawCoordinates[:k])

	// just for trying iterators, it's tha same as
	// rawCoordinates = rawCoordinates[k:]
	coordinates := func() iter.Seq[[]int] {
		return func(yield func([]int) bool) {
			for _, v := range rawCoordinates[k:] {
				if !yield(v) {
					return
				}
			}
		}
	}

	for c := range coordinates() {
		dist := euclideanDistance([2]int(c))

		if dist < points[0].distanceFromOrigin {
			pop(&points)
			push(&points, [2]int(c))
		}
	}

	res := make([][]int, k)
	for i := range points {
		res[i] = points[i].coordinates[:]
	}

	return res
}

func heapify(rawCoordinates [][]int) []point {
	heap := make([]point, 0, len(rawCoordinates))
	for i := range rawCoordinates {
		push(&heap, [2]int(rawCoordinates[i]))
	}

	return heap
}

func euclideanDistance(coordinates [2]int) float64 {
	originPoint := [2]int{0, 0}

	return math.Sqrt(
		math.Pow(
			float64(coordinates[0])-float64(originPoint[0]),
			2,
		) + math.Pow(
			float64(coordinates[1])-float64(originPoint[1]),
			2,
		),
	)
}

func push(heap *[]point, coordinates [2]int) {
	*heap = append(*heap, point{coordinates, euclideanDistance(coordinates)})

	percolateUpLastPoint(*heap)
}

func percolateUpLastPoint(heap []point) {
	parentIdx := (len(heap) - 2) / 2

	if parentIdx >= 0 && heap[len(heap)-1].distanceFromOrigin > heap[parentIdx].distanceFromOrigin {
		heap[len(heap)-1], heap[parentIdx] = heap[parentIdx], heap[len(heap)-1]
		percolateUpLastPoint(heap[:parentIdx+1])
	}
}

func pop(heap *[]point) [2]int {
	res := (*heap)[0].coordinates
	(*heap)[0], *heap = (*heap)[len(*heap)-1], (*heap)[:len(*heap)-1]

	percolateDownFirstPoint(*heap, 0)

	return res
}

func percolateDownFirstPoint(heap []point, i int) {
	leftChildIdx, rightChildIdx := i*2+1, i*2+2

	if len(heap) > leftChildIdx && heap[i].distanceFromOrigin < heap[leftChildIdx].distanceFromOrigin {
		heap[i], heap[leftChildIdx] = heap[leftChildIdx], heap[i]

		if len(heap) > rightChildIdx && heap[i].distanceFromOrigin < heap[rightChildIdx].distanceFromOrigin {
			heap[i], heap[rightChildIdx] = heap[rightChildIdx], heap[i]
			percolateDownFirstPoint(heap, rightChildIdx)
		}

		percolateDownFirstPoint(heap, leftChildIdx)

		return
	}

	if len(heap) > rightChildIdx && heap[i].distanceFromOrigin < heap[rightChildIdx].distanceFromOrigin {
		heap[i], heap[rightChildIdx] = heap[rightChildIdx], heap[i]

		percolateDownFirstPoint(heap, rightChildIdx)
	}
}
