package reorderlist_test

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
		want  *ListNode
	}{
		{
			name: "Even count of nodes",
			input: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val: 2,
					Next: &ListNode{
						Val: 3,
						Next: &ListNode{
							Val:  4,
							Next: nil,
						},
					},
				},
			},
			want: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val: 4,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val:  3,
							Next: nil,
						},
					},
				},
			},
		},
		{
			name: "Not even count of nodes",
			input: &ListNode{
				Val: 11,
				Next: &ListNode{
					Val: 8,
					Next: &ListNode{
						Val: 5,
						Next: &ListNode{
							Val: -2,
							Next: &ListNode{
								Val:  0,
								Next: nil,
							},
						},
					},
				},
			},
			want: &ListNode{
				Val: 11,
				Next: &ListNode{
					Val: 0,
					Next: &ListNode{
						Val: 8,
						Next: &ListNode{
							Val: -2,
							Next: &ListNode{
								Val:  5,
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

			reorderList(tt.input)
			assert.Equal(t, tt.want, tt.input)
		})
	}
}

// Time O(n), since the iteration count depends on input linearly
// Space O(1), we don't use any extra space.
func reorderList(head *ListNode) {
	slow, fast := head, head.Next

	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
	}

	secondPart := slow.Next
	slow.Next = nil

	var secondReversedPart *ListNode

	for secondPart != nil {
		tmp := secondPart.Next
		secondPart.Next = secondReversedPart
		secondReversedPart = secondPart
		secondPart = tmp
	}

	curr := head
	for secondReversedPart != nil {
		tmp, tmp2 := curr.Next, secondReversedPart.Next
		curr.Next = secondReversedPart
		curr.Next.Next = tmp
		curr, secondReversedPart = curr.Next.Next, tmp2
	}
}
