package validsudoku_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidSudoku(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		board [][]byte
		want  bool
	}{
		{
			name: "Valid",
			board: [][]byte{
				{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
				{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
				{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
				{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
				{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
				{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
				{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
				{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
				{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
			},
			want: true,
		},
		{
			name: "Not valid due to the first sub-square",
			board: [][]byte{
				{'8', '3', '.', '.', '7', '.', '.', '.', '.'},
				{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
				{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
				{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
				{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
				{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
				{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
				{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
				{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
			},
			want: false,
		},
		{
			name: "Not valid due to the column with index 3",
			board: [][]byte{
				{'.', '.', '4', '.', '.', '.', '6', '3', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'5', '.', '.', '.', '.', '.', '.', '9', '.'},
				{'.', '.', '.', '5', '6', '.', '.', '.', '.'},
				{'4', '.', '3', '.', '.', '.', '.', '.', '1'},
				{'.', '.', '.', '7', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '5', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
				{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isValidSudoku(tt.board))
		})
	}
}

// Time O(1) or O(9^2) we iterate consequently on each cell
// Time O(1) we allocate memory for maps equals to input.
func isValidSudoku(board [][]byte) bool {
	rows := make([]map[byte]struct{}, len(board))
	cols := make([]map[byte]struct{}, len(board))
	subSquares := make(map[[2]int]map[byte]struct{}, len(board)*len(board[0])/3)

	for i := range board {
		for j := range len(board[0]) {
			if board[i][j] == '.' {
				continue
			}

			_, duplicateInRow := rows[i][board[i][j]]
			_, duplicateInCol := cols[j][board[i][j]]
			_, duplicateInSubSquare := subSquares[[2]int{i / 3, j / 3}][board[i][j]]

			if duplicateInRow || duplicateInCol || duplicateInSubSquare {
				return false
			}

			if rows[i] == nil {
				rows[i] = make(map[byte]struct{})
			}

			if cols[j] == nil {
				cols[j] = make(map[byte]struct{})
			}

			if subSquares[[2]int{i / 3, j / 3}] == nil {
				subSquares[[2]int{i / 3, j / 3}] = make(map[byte]struct{})
			}

			rows[i][board[i][j]] = struct{}{}
			cols[j][board[i][j]] = struct{}{}
			subSquares[[2]int{i / 3, j / 3}][board[i][j]] = struct{}{}
		}
	}

	return true
}
