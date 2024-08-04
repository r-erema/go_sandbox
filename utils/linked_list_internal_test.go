package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func nodeFromTheHead() *LinkedListNode {
	node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
	node1 := &LinkedListNode{value: 1, prev: nil, next: nil}
	node2 := &LinkedListNode{value: 2, prev: nil, next: nil}
	node3 := &LinkedListNode{value: 3, prev: nil, next: nil}

	node0.next = node1
	node1.prev = node0
	node1.next = node2
	node2.prev = node1
	node2.next = node3
	node3.prev = node2

	return node0
}

func nodeFromTheTail() *LinkedListNode {
	node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
	node1 := &LinkedListNode{value: 1, prev: nil, next: nil}
	node2 := &LinkedListNode{value: 2, prev: nil, next: nil}
	node3 := &LinkedListNode{value: 3, prev: nil, next: nil}

	node0.next = node1
	node1.prev = node0
	node1.next = node2
	node2.prev = node1
	node2.next = node3
	node3.prev = node2

	return node3
}

func nodeFromTheMiddle() *LinkedListNode {
	node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
	node1 := &LinkedListNode{value: 1, prev: nil, next: nil}
	node2 := &LinkedListNode{value: 2, prev: nil, next: nil}
	node3 := &LinkedListNode{value: 3, prev: nil, next: nil}

	node0.next = node1
	node1.prev = node0
	node1.next = node2
	node2.prev = node1
	node2.next = node3
	node3.prev = node2

	return node2
}

type ttBoolAssertion struct {
	name string
	list *LinkedListNode
	want bool
}

func TestLinkedList_Build(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		nodes []*LinkedListNode
		want  *LinkedListNode
	}{
		{
			name: "Case 0",
			nodes: []*LinkedListNode{
				{value: 0},
				{value: 1},
				{value: 2},
				{value: 3},
			},
			want: nodeFromTheHead(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resultList := Build(tt.nodes)
			assert.Equal(t, tt.want, resultList)
		})
	}
}

func TestLinkedList_IsHead(t *testing.T) {
	t.Parallel()

	tests := []ttBoolAssertion{
		{
			name: "Head case",
			list: nodeFromTheHead(),
			want: true,
		},
		{
			name: "Not head case",
			list: nodeFromTheTail(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.list.IsHead())
		})
	}
}

func TestLinkedList_IsTail(t *testing.T) {
	t.Parallel()

	tests := []ttBoolAssertion{
		{
			name: "Tail case",
			list: nodeFromTheTail(),
			want: true,
		},
		{
			name: "Not tail case",
			list: nodeFromTheHead(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.list.IsTail())
		})
	}
}

func TestLinkedListNode_Head(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		listBuilder func() *LinkedListNode
		want        float32
	}{
		{
			name:        "Case 0",
			listBuilder: nodeFromTheMiddle,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.InEpsilon(t, tt.want, tt.listBuilder().Head().value, 0)
		})
	}
}

func TestLinkedListNode_Tail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		listBuilder func() *LinkedListNode
		want        float32
	}{
		{
			name:        "Case 0",
			listBuilder: nodeFromTheMiddle,
			want:        3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.InEpsilon(t, tt.want, tt.listBuilder().Tail().value, 0)
		})
	}
}

func TestLinkedListNode_Append(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		listBuilder, want func() *LinkedListNode
		nodesToAppend     []*LinkedListNode
	}{
		{
			name:        "Case 0",
			listBuilder: nodeFromTheMiddle,
			nodesToAppend: []*LinkedListNode{
				{value: -1},
				{value: 4},
			},
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
				node1 := &LinkedListNode{value: 1, prev: nil, next: nil}
				node2 := &LinkedListNode{value: 2, prev: nil, next: nil}
				node3 := &LinkedListNode{value: 3, prev: nil, next: nil}
				node4 := &LinkedListNode{value: -1, prev: nil, next: nil}
				node5 := &LinkedListNode{value: 4, prev: nil, next: nil}

				node0.next = node1
				node1.prev = node0
				node1.next = node2
				node2.prev = node1
				node2.next = node3
				node3.prev = node2
				node3.next = node4
				node4.prev = node3
				node4.next = node5
				node5.prev = node4

				return node2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resultList := tt.listBuilder()
			for _, nodeToAppend := range tt.nodesToAppend {
				resultList.Append(nodeToAppend)
			}

			wantedList := tt.want()
			assert.Equal(t, wantedList, resultList)
		})
	}
}

func TestLinkedListNode_Prepend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		listBuilder, want func() *LinkedListNode
		nodesToPrepend    []*LinkedListNode
	}{
		{
			name:        "Case 0",
			listBuilder: nodeFromTheMiddle,
			nodesToPrepend: []*LinkedListNode{
				{value: 8},
				{value: -0.5},
			},
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{value: -0.5, prev: nil, next: nil}
				node1 := &LinkedListNode{value: 8, prev: nil, next: nil}
				node2 := &LinkedListNode{value: 0, prev: nil, next: nil}
				node3 := &LinkedListNode{value: 1, prev: nil, next: nil}
				node4 := &LinkedListNode{value: 2, prev: nil, next: nil}
				node5 := &LinkedListNode{value: 3, prev: nil, next: nil}

				node0.next = node1
				node1.prev = node0
				node1.next = node2
				node2.prev = node1
				node2.next = node3
				node3.prev = node2
				node3.next = node4
				node4.prev = node3
				node4.next = node5
				node5.prev = node4

				return node4
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resultList := tt.listBuilder()
			for _, nodeToPrepend := range tt.nodesToPrepend {
				resultList.Prepend(nodeToPrepend)
			}

			wantedList := tt.want()
			assert.Equal(t, wantedList, resultList)
		})
	}
}

func TestLinkedListNode_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		list   *LinkedListNode
		needle float32
		want   *LinkedListNode
	}{
		{
			name:   "Case 0",
			list:   nodeFromTheMiddle(),
			needle: 0,
			want:   nodeFromTheMiddle().Head(),
		},
		{
			name:   "Case 1",
			list:   nodeFromTheMiddle(),
			needle: 3,
			want:   nodeFromTheMiddle().Tail(),
		},
		{
			name:   "Case 2",
			list:   nodeFromTheHead(),
			needle: 3,
			want:   nodeFromTheMiddle().Tail(),
		},
		{
			name:   "Case 3",
			list:   nodeFromTheTail(),
			needle: 0,
			want:   nodeFromTheMiddle().Head(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := tt.list.Search(tt.needle)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLinkedListNode_Remove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		list         *LinkedListNode
		nodeToRemove float32
		want         func() *LinkedListNode
	}{
		{
			name:         "Case 0",
			list:         nodeFromTheHead(),
			nodeToRemove: 2,
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
				node1 := &LinkedListNode{value: 1, prev: nil, next: nil}
				node2 := &LinkedListNode{value: 3, prev: nil, next: nil}

				node0.next = node1
				node1.prev = node0
				node1.next = node2
				node2.prev = node1

				return node0
			},
		},
		{
			name:         "Case 1",
			list:         nodeFromTheTail(),
			nodeToRemove: 1,
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{value: 0, prev: nil, next: nil}
				node1 := &LinkedListNode{value: 2, prev: nil, next: nil}
				node2 := &LinkedListNode{value: 3, prev: nil, next: nil}

				node0.next = node1
				node1.prev = node0
				node1.next = node2
				node2.prev = node1

				return node2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.list.Remove(tt.nodeToRemove)
			wantedList := tt.want()
			assert.Equal(t, wantedList, tt.list)
		})
	}
}

func TestLinkedListNode_InsertAfter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		list         *LinkedListNode
		nodeToInsert *LinkedListNode
		insertAfter  float32
		want         *LinkedListNode
	}{
		{
			name:         "Case 0",
			list:         nodeFromTheHead(),
			nodeToInsert: &LinkedListNode{value: 11, prev: nil, next: nil},
			insertAfter:  1,
			want: Build([]*LinkedListNode{
				{value: 0},
				{value: 1},
				{value: 11},
				{value: 2},
				{value: 3},
			}),
		},
		{
			name:         "Case 1",
			list:         nodeFromTheHead(),
			nodeToInsert: &LinkedListNode{value: 0.8, prev: nil, next: nil},
			insertAfter:  3,
			want: Build([]*LinkedListNode{
				{value: 0},
				{value: 1},
				{value: 2},
				{value: 3},
				{value: 0.8},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.list.InsertAfter(tt.nodeToInsert, tt.insertAfter)
			assert.Equal(t, tt.want, tt.list)
		})
	}
}
