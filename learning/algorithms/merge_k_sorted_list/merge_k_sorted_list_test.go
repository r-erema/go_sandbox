package mergeksortedlist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ListNode struct {
	Value int
	Next  *ListNode
}

func TestMergeKSortedLists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		lists []*ListNode
		want  *ListNode
	}{
		{
			name: "Case 0",
			lists: []*ListNode{
				{
					Value: 1,
					Next: &ListNode{
						Value: 4,
						Next: &ListNode{
							Value: 5,
							Next:  nil,
						},
					},
				},
				{
					Value: 1,
					Next: &ListNode{
						Value: 3,
						Next: &ListNode{
							Value: 4,
							Next:  nil,
						},
					},
				},
				{
					Value: 2,
					Next: &ListNode{
						Value: 6,
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
									Next: &ListNode{
										Value: 5,
										Next: &ListNode{
											Value: 6,
											Next:  nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mergeKSortedLists(tt.lists))
		})
	}
}

// Time O(N * logK) since we merge K lists with N numbers of nodes
// Space O(N) since mergedLists doesn't consume more memory than input.
func mergeKSortedLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	var list1, list2 *ListNode

	for len(lists) > 1 {
		mergedLists := make([]*ListNode, 0)

		for i := 0; i < len(lists); i += 2 {
			list1, list2 = lists[i], nil
			if i+1 < len(lists) {
				list2 = lists[i+1]
			}

			mergedLists = append(mergedLists, mergeTwoSortedLists(list1, list2))
		}

		lists = mergedLists
	}

	return lists[0]
}

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
