package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func nodeFromTheHead() *LinkedListNode {
	node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
	node1 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
	node2 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
	node3 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

	node0.Next = node1
	node1.Prev = node0
	node1.Next = node2
	node2.Prev = node1
	node2.Next = node3
	node3.Prev = node2

	return node0
}

func nodeFromTheTail() *LinkedListNode {
	node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
	node1 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
	node2 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
	node3 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

	node0.Next = node1
	node1.Prev = node0
	node1.Next = node2
	node2.Prev = node1
	node2.Next = node3
	node3.Prev = node2

	return node3
}

func nodeFromTheMiddle() *LinkedListNode {
	node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
	node1 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
	node2 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
	node3 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

	node0.Next = node1
	node1.Prev = node0
	node1.Next = node2
	node2.Prev = node1
	node2.Next = node3
	node3.Prev = node2

	return node2
}

type testCaseBoolAssertion struct {
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
				{Value: 0},
				{Value: 1},
				{Value: 2},
				{Value: 3},
			},
			want: nodeFromTheHead(),
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			resultList := Build(testCase.nodes)
			assert.Equal(t, testCase.want, resultList)
		})
	}
}

func TestLinkedList_IsHead(t *testing.T) {
	t.Parallel()

	tests := []testCaseBoolAssertion{
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.list.IsHead())
		})
	}
}

func TestLinkedList_IsTail(t *testing.T) {
	t.Parallel()

	tests := []testCaseBoolAssertion{
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
		tt := tt

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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.listBuilder().Head().Value)
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.listBuilder().Tail().Value)
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
				{Value: -1},
				{Value: 4},
			},
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
				node1 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
				node2 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
				node3 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}
				node4 := &LinkedListNode{Value: -1, Prev: nil, Next: nil}
				node5 := &LinkedListNode{Value: 4, Prev: nil, Next: nil}

				node0.Next = node1
				node1.Prev = node0
				node1.Next = node2
				node2.Prev = node1
				node2.Next = node3
				node3.Prev = node2
				node3.Next = node4
				node4.Prev = node3
				node4.Next = node5
				node5.Prev = node4

				return node2
			},
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			resultList := testCase.listBuilder()
			for _, nodeToAppend := range testCase.nodesToAppend {
				resultList.Append(nodeToAppend)
			}
			wantedList := testCase.want()
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
				{Value: 8},
				{Value: -0.5},
			},
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{Value: -0.5, Prev: nil, Next: nil}
				node1 := &LinkedListNode{Value: 8, Prev: nil, Next: nil}
				node2 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
				node3 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
				node4 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
				node5 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

				node0.Next = node1
				node1.Prev = node0
				node1.Next = node2
				node2.Prev = node1
				node2.Next = node3
				node3.Prev = node2
				node3.Next = node4
				node4.Prev = node3
				node4.Next = node5
				node5.Prev = node4

				return node4
			},
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			resultList := testCase.listBuilder()
			for _, nodeToPrepend := range testCase.nodesToPrepend {
				resultList.Prepend(nodeToPrepend)
			}
			wantedList := testCase.want()
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.list.Search(tt.needle), tt.want)
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
				node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
				node1 := &LinkedListNode{Value: 1, Prev: nil, Next: nil}
				node2 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

				node0.Next = node1
				node1.Prev = node0
				node1.Next = node2
				node2.Prev = node1

				return node0
			},
		},
		{
			name:         "Case 1",
			list:         nodeFromTheTail(),
			nodeToRemove: 1,
			want: func() *LinkedListNode {
				node0 := &LinkedListNode{Value: 0, Prev: nil, Next: nil}
				node1 := &LinkedListNode{Value: 2, Prev: nil, Next: nil}
				node2 := &LinkedListNode{Value: 3, Prev: nil, Next: nil}

				node0.Next = node1
				node1.Prev = node0
				node1.Next = node2
				node2.Prev = node1

				return node2
			},
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.list.Remove(testCase.nodeToRemove)
			wantedList := testCase.want()
			assert.Equal(t, wantedList, testCase.list)
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
			nodeToInsert: &LinkedListNode{Value: 11, Prev: nil, Next: nil},
			insertAfter:  1,
			want: Build([]*LinkedListNode{
				{Value: 0},
				{Value: 1},
				{Value: 11},
				{Value: 2},
				{Value: 3},
			}),
		},
		{
			name:         "Case 1",
			list:         nodeFromTheHead(),
			nodeToInsert: &LinkedListNode{Value: 0.8, Prev: nil, Next: nil},
			insertAfter:  3,
			want: Build([]*LinkedListNode{
				{Value: 0},
				{Value: 1},
				{Value: 2},
				{Value: 3},
				{Value: 0.8},
			}),
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.list.InsertAfter(testCase.nodeToInsert, testCase.insertAfter)
			assert.Equal(t, testCase.want, testCase.list)
		})
	}
}
