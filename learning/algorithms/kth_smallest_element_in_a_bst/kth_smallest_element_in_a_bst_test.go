package kthsmallestelementinabst_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/utils/data_structure/tree"
	"github.com/stretchr/testify/assert"
)

func TestKthSmallest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		root    *tree.Node
		k, want int
	}{
		{
			name: "Tree 1",
			root: &tree.Node{
				Val: 3,
				Left: &tree.Node{
					Val: 1,
					Right: &tree.Node{
						Val: 2,
					},
				},
				Right: &tree.Node{
					Val: 4,
				},
			},
			k:    1,
			want: 1,
		},
		{
			name: "Tree 2",
			root: &tree.Node{
				Val: 5,
				Left: &tree.Node{
					Val: 3,
					Left: &tree.Node{
						Val:  2,
						Left: &tree.Node{Val: 1},
					},
					Right: &tree.Node{
						Val: 4,
					},
				},
				Right: &tree.Node{
					Val: 6,
				},
			},

			k:    3,
			want: 3,
		},
		{
			name: "Tree 3",
			root: &tree.Node{
				Val: 3,
				Left: &tree.Node{
					Val: 1,
					Right: &tree.Node{
						Val: 2,
					},
				},
				Right: &tree.Node{
					Val: 4,
				},
			},

			k:    2,
			want: 2,
		},
		{
			name: "Tree 4",
			root: &tree.Node{
				Val: 5,
				Left: &tree.Node{
					Val: 3,
					Left: &tree.Node{
						Val: 1,
						Right: &tree.Node{
							Val: 2,
						},
					},
					Right: &tree.Node{
						Val: 4,
					},
				},
				Right: &tree.Node{
					Val: 6,
				},
			},

			k:    2,
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, kthSmallest(tt.root, tt.k))
		})
	}
}

// Time O(n), since we visit each node of the tree
// Space O(n), since the recursion stack grows as nodes count.
func kthSmallest(root *tree.Node, m int) int {
	var res int

	var inorderTraversal func(node *tree.Node)
	inorderTraversal = func(node *tree.Node) {
		if node == nil {
			return
		}

		inorderTraversal(node.Left)

		m--
		if m == 0 {
			res = node.Val

			return
		}

		inorderTraversal(node.Right)
	}

	inorderTraversal(root)

	return res
}
