package utils

type LinkedListNode struct {
	Value      float32
	Prev, Next *LinkedListNode
}

func Build(nodes []*LinkedListNode) *LinkedListNode {
	var startNode *LinkedListNode

	iterationThreshold := len(nodes) - 1

	for index := range nodes {
		var currentNode, prevNode, nextNode *LinkedListNode
		currentNode = nodes[index]

		if index != 0 {
			prevNode = nodes[index-1]
			prevNode.Next = currentNode
		} else {
			startNode = nodes[index]
		}

		if index < iterationThreshold {
			nextNode = nodes[index+1]
			nextNode.Prev = currentNode
		}

		currentNode.Next = nextNode
		currentNode.Prev = prevNode
	}

	return startNode
}

func (node LinkedListNode) IsHead() bool {
	return node.Prev == nil
}

func (node *LinkedListNode) IsTail() bool {
	return node.Next == nil
}

func (node *LinkedListNode) Head() *LinkedListNode {
	currentNode := node

	for {
		if currentNode.Prev == nil {
			return currentNode
		}

		currentNode = currentNode.Prev
	}
}

func (node *LinkedListNode) Tail() *LinkedListNode {
	currentNode := node

	for {
		if currentNode.Next == nil {
			return currentNode
		}

		currentNode = currentNode.Next
	}
}

func (node *LinkedListNode) Append(nodeToAppend *LinkedListNode) {
	tail := node.Tail()
	tail.Next = nodeToAppend
	nodeToAppend.Prev = tail
}

func (node *LinkedListNode) Prepend(nodeToPrepend *LinkedListNode) {
	head := node.Head()
	head.Prev = nodeToPrepend
	nodeToPrepend.Next = head
}

func (node *LinkedListNode) Search(needle float32) *LinkedListNode {
	currentNode := node

	for {
		if currentNode.Value == needle {
			return currentNode
		}

		if currentNode.Next == nil {
			break
		}

		currentNode = currentNode.Next
	}

	if currentNode = node.Prev; currentNode != nil {
		for {
			if currentNode.Value == needle {
				return currentNode
			}

			if currentNode.Prev == nil {
				break
			}

			currentNode = currentNode.Prev
		}
	}

	return nil
}

func (node *LinkedListNode) Remove(needle float32) {
	nodeToRemove := node.Search(needle)
	if nodeToRemove == nil {
		return
	}

	prev, next := nodeToRemove.Prev, nodeToRemove.Next

	if prev != nil {
		prev.Next = next
	}

	if next != nil {
		next.Prev = prev
	}
}

func (node *LinkedListNode) InsertAfter(nodeToInsert *LinkedListNode, mountPoint float32) {
	startNode := node.Search(mountPoint)
	if startNode == nil {
		return
	}

	next := startNode.Next

	startNode.Next = nodeToInsert
	nodeToInsert.Prev = startNode
	nodeToInsert.Next = next

	if next != nil {
		next.Prev = nodeToInsert
	}
}
