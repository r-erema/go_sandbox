package breadthfirstsearch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		edges          [][]int
		vertices, want []int
	}{
		{
			name: "8 nodes graph",
			edges: [][]int{
				0: {1, 2, 3},
				1: {0, 8},
				2: {0, 4},
				3: {0, 6},
				4: {2, 5, 7},
				5: {4},
				6: {3},
				7: {4},
				8: {1},
			},
			vertices: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:     []int{1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name: "2 unconnected graphs",
			edges: [][]int{
				0: {1, 2},
				1: {0, 2},
				2: {0, 1},
				3: {4, 5},
				4: {3, 5},
				5: {3, 4},
			},
			vertices: []int{0, 0, 0, 0, 0, 0},
			want:     []int{1, 1, 1, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bfs := func(nodeIndex int) {
				queue := []int{nodeIndex}

				for len(queue) > 0 {
					nodeIndex, queue = queue[0], queue[1:]
					tt.vertices[nodeIndex] = 1

					if nodeIndex < len(tt.edges) {
						for _, edge := range tt.edges[nodeIndex] {
							if tt.vertices[edge] == 0 {
								queue = append(queue, edge)
							}
						}
					}
				}
			}

			bfs(0)
			assert.Equal(t, tt.want, tt.vertices)
		})
	}
}
