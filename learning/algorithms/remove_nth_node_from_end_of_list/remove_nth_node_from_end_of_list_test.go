package removenthnodefromendoflist_test

import (
	"testing"

	linkedlist "github.com/r-erema/go_sendbox/utils/data_structure/linked_list"
	"github.com/stretchr/testify/assert"
)

func TestRemoveNthFromEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		head *linkedlist.Node
		n    int
		want *linkedlist.Node
	}{
		{
			name: "normal list",
			head: &linkedlist.Node{
				Val: 1,
				Next: &linkedlist.Node{
					Val: 2,
					Next: &linkedlist.Node{
						Val: 3,
						Next: &linkedlist.Node{
							Val:  4,
							Next: &linkedlist.Node{Val: 5},
						},
					},
				},
			},
			n: 2,
			want: &linkedlist.Node{
				Val: 1,
				Next: &linkedlist.Node{
					Val: 2,
					Next: &linkedlist.Node{
						Val: 3,
						Next: &linkedlist.Node{
							Val: 5,
						},
					},
				},
			},
		},
		{
			name: "1 node list",
			head: &linkedlist.Node{Val: 1},
			n:    1,
			want: nil,
		},
		{
			name: "2 nodes list, remove the last node",
			head: &linkedlist.Node{
				Val: 1,
				Next: &linkedlist.Node{
					Val: 2,
				},
			},
			n:    1,
			want: &linkedlist.Node{Val: 1},
		},
		{
			name: "2 nodes list, remove the 1 node",
			head: &linkedlist.Node{
				Val: 1,
				Next: &linkedlist.Node{
					Val: 2,
				},
			},
			n:    2,
			want: &linkedlist.Node{Val: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, removeNthFromEnd(tt.head, tt.n))
		})
	}
}

// Time O(n), since we need to iterate each node 1 time
// Time O(1), since we don't involve any additional data structure.
func removeNthFromEnd(head *linkedlist.Node, n int) *linkedlist.Node {
	dummy := &linkedlist.Node{Next: head}
	left, right := dummy, head

	for ; n > 0; n-- {
		right = right.Next
	}

	for right != nil {
		left, right = left.Next, right.Next
	}

	left.Next = left.Next.Next

	return dummy.Next
}
