package mergetwosortedlists_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	ListNode struct {
		Value int
		Next  *ListNode
	}
	tt struct {
		name  string
		list1 *ListNode
		list2 *ListNode
		want  *ListNode
	}
)

// https://leetcode.com/problems/merge-two-sorted-lists/
func TestMerge2SortedLists(t *testing.T) {
	t.Parallel()

	tests := []tt{
		normalListstt(),
		listIsNil(),
		listHasNegativeNumber(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mergeTwoSortedLists(tt.list1, tt.list2))
		})
	}
}

// Time O(n+m)
// n = number of nodes in list1
// m = number of nodes in list2
//
// Space O(1)
// we have a constant space, since we are just shifting the pointers.
func mergeTwoSortedLists(list1, list2 *ListNode) *ListNode {
	dummyNode := &ListNode{}
	tail := dummyNode

	for list1 != nil && list2 != nil {
		if list1.Value < list2.Value {
			tail.Next = list1
			list1 = list1.Next
		} else {
			tail.Next = list2
			list2 = list2.Next
		}

		tail = tail.Next
	}

	if list1 != nil {
		tail.Next = list1
	}

	if list2 != nil {
		tail.Next = list2
	}

	return dummyNode.Next
}

func normalListstt() tt {
	return tt{
		name: "Normal lists",
		list1: &ListNode{
			Value: 1,
			Next: &ListNode{
				Value: 2,
				Next: &ListNode{
					Value: 4,
					Next:  nil,
				},
			},
		},
		list2: &ListNode{
			Value: 1,
			Next: &ListNode{
				Value: 3,
				Next: &ListNode{
					Value: 4,
					Next:  nil,
				},
			},
		},
		want: &ListNode{
			Value: 1,
			Next: &ListNode{
				Value: 1,
				Next: &ListNode{
					Value: 2,
					Next: &ListNode{
						Value: 3,
						Next: &ListNode{
							Value: 4,
							Next: &ListNode{
								Value: 4,
								Next:  nil,
							},
						},
					},
				},
			},
		},
	}
}

func listIsNil() tt {
	return tt{
		name: "List is nil",
		list1: &ListNode{
			Value: 1,
			Next:  nil,
		},
		list2: nil,
		want: &ListNode{
			Value: 1,
			Next:  nil,
		},
	}
}

func listHasNegativeNumber() tt {
	return tt{
		name: "List has a negative number",
		list1: &ListNode{
			Value: -9,
			Next: &ListNode{
				Value: 3,
				Next:  nil,
			},
		},
		list2: &ListNode{
			Value: 5,
			Next: &ListNode{
				Value: 7,
				Next:  nil,
			},
		},
		want: &ListNode{
			Value: -9,
			Next: &ListNode{
				Value: 3,
				Next: &ListNode{
					Value: 5,
					Next: &ListNode{
						Value: 7,
						Next:  nil,
					},
				},
			},
		},
	}
}
