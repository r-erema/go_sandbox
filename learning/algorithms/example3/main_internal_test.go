package example3

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils"
	"github.com/stretchr/testify/assert"
)

func bst1() *utils.BST {
	return utils.NewBST(5).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(10)).
		InsertRecursively(utils.NewBST(8)).
		InsertRecursively(utils.NewBST(34))
}

func bst2() *utils.BST {
	return utils.NewBST(9).
		InsertRecursively(utils.NewBST(4)).
		InsertRecursively(utils.NewBST(17)).
		InsertRecursively(utils.NewBST(3)).
		InsertRecursively(utils.NewBST(6)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(7)).
		InsertRecursively(utils.NewBST(20))
}

func bst3() *utils.BST {
	return utils.NewBST(10).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(15)).
		InsertRecursively(utils.NewBST(2)).
		InsertRecursively(utils.NewBST(5)).
		InsertRecursively(utils.NewBST(13)).
		InsertRecursively(utils.NewBST(22)).
		InsertRecursively(utils.NewBST(14))
}

func TestFindClosestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bst    *utils.BST
		target float32
		want   float32
	}{
		{
			name:   "case 0",
			bst:    bst1(),
			target: 20,
			want:   10,
		},
		{
			name:   "case 1",
			bst:    bst2(),
			target: 4,
			want:   4,
		},
		{
			name:   "case 2",
			bst:    bst2(),
			target: 18,
			want:   17,
		},
		{
			name:   "case 3",
			bst:    bst2(),
			target: 12,
			want:   9,
		},
		{
			name:   "case 4",
			bst:    bst3(),
			target: 12,
			want:   13,
		},
		{
			name:   "case 5",
			bst:    bst3(),
			target: 3.5,
			want:   5,
		},
		{
			name:   "case 6",
			bst:    bst3(),
			target: 22,
			want:   22,
		},
		{
			name:   "case 7",
			bst:    bst3(),
			target: 13.6,
			want:   14,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			closest := FindClosestValueRecursively(tt.bst, tt.target)
			assert.Equal(t, tt.want, closest)
			closest = FindClosestValueIteratively(tt.bst, tt.target)
			assert.Equal(t, tt.want, closest)
		})
	}
}
