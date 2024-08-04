package reverselinkedlist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	ListNode struct {
		Val  int
		Next *ListNode
	}
)

func TestReverseLinkedList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ListNode
		want  *ListNode
	}{
		{
			name: "Simple list",
			input: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val: 2,
					Next: &ListNode{
						Val: 3,
						Next: &ListNode{
							Val: 4,
							Next: &ListNode{
								Val:  5,
								Next: nil,
							},
						},
					},
				},
			},
			want: &ListNode{
				Val: 5,
				Next: &ListNode{
					Val: 4,
					Next: &ListNode{
						Val: 3,
						Next: &ListNode{
							Val: 2,
							Next: &ListNode{
								Val:  1,
								Next: nil,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, reverseList(tt.input))
		})
	}
}

// Time O(n), since we should iterate all the input
// Space O(1), we allocate a new linked list, but we reduce an input linked list.
func reverseList(head *ListNode) *ListNode {
	dummyNode := &ListNode{}

	for head != nil {
		dummyNode.Next = &ListNode{Val: head.Val, Next: dummyNode.Next}
		head = head.Next
	}

	return dummyNode.Next
}
