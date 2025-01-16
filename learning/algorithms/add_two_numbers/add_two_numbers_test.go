package addtwonumbers_test

import (
	"testing"

	linkedlist "github.com/r-erema/go_sendbox/utils/data_structure/linked_list"
	"github.com/stretchr/testify/assert"
)

func TestAddTwoNumbers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		list1, list2, want *linkedlist.Node
	}{
		{
			name: "addition with moving to the next rank",
			list1: &linkedlist.Node{
				Val: 2,
				Next: &linkedlist.Node{
					Val: 4,
					Next: &linkedlist.Node{
						Val: 3,
					},
				},
			},
			list2: &linkedlist.Node{
				Val: 5,
				Next: &linkedlist.Node{
					Val: 6,
					Next: &linkedlist.Node{
						Val: 4,
					},
				},
			},
			want: &linkedlist.Node{
				Val: 7,
				Next: &linkedlist.Node{
					Val: 0,
					Next: &linkedlist.Node{
						Val: 8,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, addTwoNumbers(tt.list1, tt.list2))
		})
	}
}

// Time O(m+n), since we should iterate both lists
// Space O(1), sine we don't allocate any additional memory.
func addTwoNumbers(list1, list2 *linkedlist.Node) *linkedlist.Node {
	curr1, curr2, carry := list1, list2, 0

	for curr1 != nil && curr2 != nil {
		val1, val2 := curr1.Val, curr2.Val
		curr1.Val += val2
		curr2.Val += val1

		curr1.Val += carry
		curr2.Val += carry
		carry = 0

		if curr1.Val >= 10 {
			curr1.Val %= 10
			curr2.Val %= 10
			carry = 1
		}

		if curr1.Next == nil && curr2.Next == nil && carry == 1 {
			curr1.Next = &linkedlist.Node{Val: 1}

			return list1
		}

		curr1, curr2 = curr1.Next, curr2.Next
	}

	if list := finishList(curr1, list1, carry); list != nil {
		return list
	}

	if list := finishList(curr2, list2, carry); list != nil {
		return list
	}

	return list1
}

func finishList(curr, list *linkedlist.Node, carry int) *linkedlist.Node {
	for curr != nil {
		curr.Val += carry
		carry = 0

		if curr.Val >= 10 {
			curr.Val %= 10
			carry = 1
		}

		if curr.Next == nil {
			if carry == 1 {
				curr.Next = &linkedlist.Node{Val: 1}
			}

			return list
		}

		curr = curr.Next
	}

	return nil
}
