package wordsearch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		board [][]byte
		word  string
		want  bool
	}{
		{
			name: "word 0",
			board: [][]byte{
				{'A', 'B'},
			},
			word: "AB",
			want: true,
		},
		{
			name: "word 0",
			board: [][]byte{
				{'A'},
				{'B'},
			},
			word: "AB",
			want: true,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'C', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word: "ESE",
			want: true,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'B', 'G'},
				{'G', 'S'},
				{'S', 'A'},
			},
			word: "BG",
			want: true,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'A', 'B', 'C', 'E'},
				{'S', 'F', 'E', 'S'},
				{'A', 'D', 'E', 'E'},
			},
			word: "ABCESEEEFS",
			want: true,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'A', 'A'},
				{'A', 'A'},
			},
			word: "AAAAA",
			want: false,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'A', 'A'},
				{'A', 'A'},
			},
			word: "AAAB",
			want: false,
		},
		{
			name: "word 1",
			board: [][]byte{
				{'a', 'a', 'b', 'a', 'a', 'b'},
				{'a', 'a', 'b', 'b', 'b', 'a'},
				{'a', 'a', 'a', 'a', 'b', 'a'},
				{'b', 'a', 'b', 'b', 'a', 'b'},
				{'a', 'b', 'b', 'a', 'b', 'a'},
				{'b', 'a', 'a', 'a', 'a', 'b'},
			},
			word: "bbbaabbbbbab",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, exist(tt.board, tt.word))
		})
	}
}

// Time O(m * 4^n), since we can start from any cell and explore all 4 directions
// Space O(n), since we can store the whole word in the visited map.
func exist(board [][]byte, word string) bool {
	if len(word) > len(board)*len(board[0]) {
		return false
	}

	var (
		res     bool
		dfs     func(i, dirX, dirY int)
		visited = make(map[[2]int]struct{})
	)

	dfs = func(i, dirX, dirY int) {
		if res || !isValidPosition(dirX, dirY, board) || word[i] != board[dirX][dirY] {
			return
		}

		key := [2]int{dirX, dirY}
		if _, exists := visited[key]; exists {
			return
		}

		visited[key] = struct{}{}
		defer delete(visited, key)

		if i == len(word)-1 {
			res = true

			return
		}

		neighbors := [][]int{
			{dirX, dirY + 1},
			{dirX - 1, dirY},
			{dirX, dirY - 1},
			{dirX + 1, dirY},
		}
		for _, neighbor := range neighbors {
			dfs(i+1, neighbor[0], neighbor[1])
		}
	}

	for x := range board {
		for y := range board[x] {
			dfs(0, x, y)
		}
	}

	return res
}

func isValidPosition(dirX, dirY int, board [][]byte) bool {
	return dirX >= 0 && dirX < len(board) && dirY >= 0 && dirY < len(board[0])
}
