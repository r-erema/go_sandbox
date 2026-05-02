package max_area_of_island_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxAreaOfIsland(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{
			name: "6 islands",
			grid: [][]int{
				{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
				{0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
				{0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, maxAreaOfIsland(tt.grid))
		})
	}
}

// Time O(N * M), where M and N are rows and columns, we need to visit each cell in a grid
// Space O(N * M), since we can store all cells in the visited hash map.

func maxAreaOfIsland(grid [][]int) int {
	visited := make(map[[2]int]struct{}, len(grid)*len(grid[0]))

	var maxArea int

	for i := range grid {
		for j := range grid[i] {
			if _, ok := visited[[2]int{i, j}]; ok || grid[i][j] == 0 {
				continue
			}

			visited[[2]int{i, j}] = struct{}{}

			currArea := bfs(grid, [2]int{i, j}, visited)

			maxArea = max(maxArea, currArea)
		}
	}

	return maxArea
}

func bfs(grid [][]int, startPoint [2]int, visited map[[2]int]struct{}) int {
	var (
		queue      = [][2]int{startPoint}
		curr       [2]int
		currArea   = 0
		directions = [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	)

	for len(queue) > 0 {
		currArea++

		curr, queue = queue[0], queue[1:]

		for i := range directions {
			directionPoint := [2]int{curr[0] + directions[i][0], curr[1] + directions[i][1]}

			if _, ok := visited[directionPoint]; ok {
				continue
			}

			pointValid := directionPoint[0] < len(grid) &&
				directionPoint[1] < len(grid[0]) &&
				directionPoint[0] >= 0 &&
				directionPoint[1] >= 0

			if pointValid && grid[directionPoint[0]][directionPoint[1]] == 1 {
				visited[directionPoint] = struct{}{}

				queue = append(queue, directionPoint)
			}
		}
	}

	return currArea
}
