package linkedlistcycle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func TestReverseLinkedList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ListNode
		pos   int
		want  bool
	}{
		{
			name: "Simple list",
			input: func() *ListNode {
				node1 := &ListNode{
					Val:  3,
					Next: nil,
				}
				node2 := &ListNode{
					Val:  2,
					Next: nil,
				}
				node3 := &ListNode{
					Val:  0,
					Next: nil,
				}
				node4 := &ListNode{
					Val:  -4,
					Next: nil,
				}

				node1.Next = node2
				node2.Next = node3
				node3.Next = node4
				node4.Next = node2

				return node1
			}(),
			pos:  1,
			want: true,
		},
		{
			name:  "Empty list",
			input: nil,
			pos:   -1,
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, hasCycle(tt.input))
		})
	}
}

// Time O(n), since we should iterate all the input
// Space O(1), we don't allocate additional memory.
func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}

	return false
}
