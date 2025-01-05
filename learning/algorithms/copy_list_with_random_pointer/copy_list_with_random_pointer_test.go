package copylistwithrandompointer_test

import (
	"testing"

	linkedlist "github.com/r-erema/go_sendbox/utils/data_structure/linked_list"
	"github.com/stretchr/testify/assert"
)

func TestMinEatingSpeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		buildList func() *linkedlist.Node
	}{
		{
			name: "2 nodes list",
			buildList: func() *linkedlist.Node {
				node1 := &linkedlist.Node{
					Val:  8,
					Next: nil,
					Prev: nil,
				}

				node2 := &linkedlist.Node{
					Val:  11,
					Next: nil,
					Prev: nil,
				}

				node1.Next = node2
				node1.Prev = node2

				node2.Next = nil
				node2.Prev = node1

				return node1
			},
		},
		{
			name: "5 nodes list",
			buildList: func() *linkedlist.Node {
				node1 := &linkedlist.Node{
					Val:  7,
					Next: nil,
					Prev: nil,
				}

				node2 := &linkedlist.Node{
					Val:  13,
					Next: nil,
					Prev: nil,
				}

				node3 := &linkedlist.Node{
					Val:  13,
					Next: nil,
					Prev: nil,
				}

				node4 := &linkedlist.Node{
					Val:  13,
					Next: nil,
					Prev: nil,
				}

				node5 := &linkedlist.Node{
					Val:  13,
					Next: nil,
					Prev: nil,
				}

				node1.Next = node2
				node1.Prev = nil

				node2.Next = node3
				node2.Prev = node1

				node3.Next = node4
				node3.Prev = node5

				node4.Next = node5
				node4.Prev = node3

				node5.Next = nil
				node5.Prev = node1

				return node1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			head := tt.buildList()
			copied := copyPrevList(head)
			assert.True(t, isClone(head, copied))
		})
	}
}

func isClone(originalList, clonedList *linkedlist.Node) bool {
	for originalList != nil && clonedList != nil {
		if originalList == clonedList || originalList.Val != clonedList.Val {
			return false
		}

		if randomNotValid(originalList, clonedList) {
			return false
		}

		originalList, clonedList = originalList.Next, clonedList.Next
	}

	return originalList == nil && clonedList == nil
}

func randomNotValid(originalList, clonedList *linkedlist.Node) bool {
	randomHaveSamePointer := originalList.Prev == clonedList.Prev && originalList.Prev != nil
	randomNotNilAndHaveDiffVals := originalList.Prev != nil && clonedList.Prev != nil && clonedList.Prev.Val != originalList.Prev.Val

	return randomHaveSamePointer || randomNotNilAndHaveDiffVals
}

// Time O(3n), since we should iterate all the input 3 times
// Space O(1), sine we don't allocate any additional memory.
func copyPrevList(head *linkedlist.Node) *linkedlist.Node {
	curr := head
	for curr != nil {
		tail := curr.Next
		curr.Next = &linkedlist.Node{Val: curr.Val, Next: tail, Prev: curr.Prev}
		curr = curr.Next.Next
	}

	curr = head
	for i := 0; curr != nil; i++ {
		if i%2 == 1 && curr.Prev != nil {
			curr.Prev = curr.Prev.Next
		}

		curr = curr.Next
	}

	dummy := &linkedlist.Node{}
	dummyCurr := dummy
	curr = head

	for i := 0; curr != nil; i++ {
		dummyCurr.Next = curr.Next

		if curr.Next != nil {
			curr.Next = curr.Next.Next
		}

		dummyCurr = dummyCurr.Next
		curr = curr.Next
	}

	return dummy.Next
}
